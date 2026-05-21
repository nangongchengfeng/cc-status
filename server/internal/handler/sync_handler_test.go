package handler

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"cc-status/server/internal/model/entity"
	"cc-status/server/internal/repository"
	"cc-status/server/internal/service"

	"gorm.io/gorm"
)

type legacySyncResponse struct {
	Code           int    `json:"code"`
	Message        string `json:"message"`
	AcceptedCount  int    `json:"accepted_count"`
	DuplicateCount int    `json:"duplicate_count"`
}

func TestSyncRoutePersistsBatchAndReturnsLegacyPayload(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "server.db")
	db, err := repository.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("OpenDatabase() returned error: %v", err)
	}
	if err := repository.InitializeSchema(db); err != nil {
		t.Fatalf("InitializeSchema() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	router := newTestSyncRouter(t, db)

	body := map[string]any{
		"client_id": "client-1",
		"reports": []map[string]any{
			{
				"request_id":            "session:msg-1",
				"app_type":              "claude",
				"model":                 "claude-sonnet-4-0",
				"input_tokens":          10,
				"output_tokens":         20,
				"cache_read_tokens":     0,
				"cache_creation_tokens": 0,
				"created_at":            int64(1743840000),
				"session_id":            "session-1",
				"data_source":           "session_log",
			},
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Marshal() returned error: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/sync", bytes.NewReader(payload))
	request.Header.Set("Authorization", "Bearer secret-token")
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response legacySyncResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if response.Code != 0 || response.Message != "success" {
		t.Fatalf("unexpected legacy response: %+v", response)
	}
	if response.AcceptedCount != 1 || response.DuplicateCount != 0 {
		t.Fatalf("unexpected counters: %+v", response)
	}

	var storedCount int64
	if err := db.Table("usage_reports").Where("client_id = ? AND request_id = ?", "client-1", "session:msg-1").Count(&storedCount).Error; err != nil {
		t.Fatalf("count stored report returned error: %v", err)
	}
	if storedCount != 1 {
		t.Fatalf("expected one stored report, got %d", storedCount)
	}
}

func TestSyncRouteRejectsInvalidPayloads(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		body map[string]any
	}{
		{
			name: "empty reports",
			body: map[string]any{
				"client_id": "client-1",
				"reports":   []map[string]any{},
			},
		},
		{
			name: "invalid data source",
			body: map[string]any{
				"client_id": "client-1",
				"reports": []map[string]any{
					{
						"request_id":            "session:msg-1",
						"app_type":              "claude",
						"model":                 "claude-sonnet-4-0",
						"input_tokens":          10,
						"output_tokens":         20,
						"cache_read_tokens":     0,
						"cache_creation_tokens": 0,
						"created_at":            int64(1743840000),
						"session_id":            "session-1",
						"data_source":           "other_source",
					},
				},
			},
		},
		{
			name: "missing request id",
			body: map[string]any{
				"client_id": "client-1",
				"reports": []map[string]any{
					{
						"app_type":              "claude",
						"model":                 "claude-sonnet-4-0",
						"input_tokens":          10,
						"output_tokens":         20,
						"cache_read_tokens":     0,
						"cache_creation_tokens": 0,
						"created_at":            int64(1743840000),
						"session_id":            "session-1",
						"data_source":           "session_log",
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			dbPath := filepath.Join(t.TempDir(), "server.db")
			db, err := repository.OpenDatabase(dbPath)
			if err != nil {
				t.Fatalf("OpenDatabase() returned error: %v", err)
			}
			if err := repository.InitializeSchema(db); err != nil {
				t.Fatalf("InitializeSchema() returned error: %v", err)
			}

			sqlDB, err := db.DB()
			if err != nil {
				t.Fatalf("DB() returned error: %v", err)
			}
			t.Cleanup(func() {
				_ = sqlDB.Close()
			})

			router := newTestSyncRouter(t, db)
			recorder := performSyncRequest(t, router, testCase.body)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d with body %s", recorder.Code, recorder.Body.String())
			}
			if !strings.Contains(recorder.Body.String(), "INVALID_REQUEST") {
				t.Fatalf("expected INVALID_REQUEST, got %s", recorder.Body.String())
			}
		})
	}
}

func TestSyncRouteCountsBatchAndStoredDuplicates(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "server.db")
	db, err := repository.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("OpenDatabase() returned error: %v", err)
	}
	if err := repository.InitializeSchema(db); err != nil {
		t.Fatalf("InitializeSchema() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	router := newTestSyncRouter(t, db)
	requestBody := map[string]any{
		"client_id": "client-1",
		"reports": []map[string]any{
			{
				"request_id":            "session:msg-dup",
				"app_type":              "claude",
				"model":                 "claude-sonnet-4-0",
				"input_tokens":          10,
				"output_tokens":         20,
				"cache_read_tokens":     0,
				"cache_creation_tokens": 0,
				"created_at":            int64(1743840000),
				"session_id":            "session-1",
				"data_source":           "session_log",
			},
			{
				"request_id":            "session:msg-dup",
				"app_type":              "claude",
				"model":                 "claude-sonnet-4-0",
				"input_tokens":          11,
				"output_tokens":         21,
				"cache_read_tokens":     0,
				"cache_creation_tokens": 0,
				"created_at":            int64(1743840001),
				"session_id":            "session-1",
				"data_source":           "session_log",
			},
		},
	}

	firstResponse := decodeSyncResponse(t, performSyncRequest(t, router, requestBody))
	if firstResponse.AcceptedCount != 1 || firstResponse.DuplicateCount != 1 {
		t.Fatalf("unexpected first duplicate counters: %+v", firstResponse)
	}

	secondResponse := decodeSyncResponse(t, performSyncRequest(t, router, requestBody))
	if secondResponse.AcceptedCount != 0 || secondResponse.DuplicateCount != 2 {
		t.Fatalf("unexpected stored duplicate counters: %+v", secondResponse)
	}

	var storedCount int64
	if err := db.Table("usage_reports").Where("client_id = ? AND request_id = ?", "client-1", "session:msg-dup").Count(&storedCount).Error; err != nil {
		t.Fatalf("count stored duplicate report returned error: %v", err)
	}
	if storedCount != 1 {
		t.Fatalf("expected one stored duplicate report, got %d", storedCount)
	}
}

func TestSyncRouteCalculatesCostsWithExactPricing(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "server.db")
	db, err := repository.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("OpenDatabase() returned error: %v", err)
	}
	if err := repository.InitializeSchema(db); err != nil {
		t.Fatalf("InitializeSchema() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	router := newTestSyncRouter(t, db)
	response := decodeSyncResponse(t, performSyncRequest(t, router, map[string]any{
		"client_id": "client-1",
		"reports": []map[string]any{
			{
				"request_id":            "session:exact-cost",
				"app_type":              "claude",
				"model":                 "CLAUDE-SONNET-4-0",
				"input_tokens":          int64(1_000_000),
				"output_tokens":         int64(2_000_000),
				"cache_read_tokens":     int64(500_000),
				"cache_creation_tokens": int64(250_000),
				"created_at":            int64(1743840000),
				"session_id":            "session-1",
				"data_source":           "session_log",
			},
		},
	}))
	if response.AcceptedCount != 1 || response.DuplicateCount != 0 {
		t.Fatalf("unexpected sync response: %+v", response)
	}

	var stored entity.UsageReport
	if err := db.Where("client_id = ? AND request_id = ?", "client-1", "session:exact-cost").First(&stored).Error; err != nil {
		t.Fatalf("load stored report returned error: %v", err)
	}

	if stored.Model != "claude-sonnet-4-0" {
		t.Fatalf("expected normalized model, got %q", stored.Model)
	}
	if stored.PricingSource != "exact" {
		t.Fatalf("expected exact pricing source, got %q", stored.PricingSource)
	}
	assertDecimalEqual(t, stored.InputCostUSD, "3.0000000000")
	assertDecimalEqual(t, stored.OutputCostUSD, "30.0000000000")
	assertDecimalEqual(t, stored.CacheReadCostUSD, "0.1500000000")
	assertDecimalEqual(t, stored.CacheCreationCostUSD, "0.9375000000")
	assertDecimalEqual(t, stored.TotalCostUSD, "34.0875000000")
}

func TestSyncRouteUsesLongestPrefixPricing(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "server.db")
	db, err := repository.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("OpenDatabase() returned error: %v", err)
	}
	if err := repository.InitializeSchema(db); err != nil {
		t.Fatalf("InitializeSchema() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	for _, pricing := range []entity.ModelPricing{
		{
			ModelID:                     "claude",
			DisplayName:                 "Claude Family",
			InputCostPerMillion:         "9",
			OutputCostPerMillion:        "19",
			CacheReadCostPerMillion:     "0",
			CacheCreationCostPerMillion: "0",
		},
		{
			ModelID:                     "claude-sonnet",
			DisplayName:                 "Claude Sonnet Family",
			InputCostPerMillion:         "5",
			OutputCostPerMillion:        "17",
			CacheReadCostPerMillion:     "0",
			CacheCreationCostPerMillion: "0",
		},
	} {
		if err := db.Create(&pricing).Error; err != nil {
			t.Fatalf("seed prefix pricing returned error: %v", err)
		}
	}

	router := newTestSyncRouter(t, db)
	response := decodeSyncResponse(t, performSyncRequest(t, router, map[string]any{
		"client_id": "client-1",
		"reports": []map[string]any{
			{
				"request_id":            "session:prefix-cost",
				"app_type":              "claude",
				"model":                 "Claude-Sonnet-4-2026",
				"input_tokens":          int64(1_000_000),
				"output_tokens":         int64(1_000_000),
				"cache_read_tokens":     int64(0),
				"cache_creation_tokens": int64(0),
				"created_at":            int64(1743840000),
				"session_id":            "session-1",
				"data_source":           "session_log",
			},
		},
	}))
	if response.AcceptedCount != 1 || response.DuplicateCount != 0 {
		t.Fatalf("unexpected sync response: %+v", response)
	}

	var stored entity.UsageReport
	if err := db.Where("client_id = ? AND request_id = ?", "client-1", "session:prefix-cost").First(&stored).Error; err != nil {
		t.Fatalf("load stored report returned error: %v", err)
	}

	if stored.PricingSource != "prefix" {
		t.Fatalf("expected prefix pricing source, got %q", stored.PricingSource)
	}
	assertDecimalEqual(t, stored.InputCostUSD, "5")
	assertDecimalEqual(t, stored.OutputCostUSD, "17")
	assertDecimalEqual(t, stored.TotalCostUSD, "22")
}

func TestSyncRouteFallsBackToDefaultPricing(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "server.db")
	db, err := repository.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("OpenDatabase() returned error: %v", err)
	}
	if err := repository.InitializeSchema(db); err != nil {
		t.Fatalf("InitializeSchema() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	router := newTestSyncRouter(t, db)
	response := decodeSyncResponse(t, performSyncRequest(t, router, map[string]any{
		"client_id": "client-1",
		"reports": []map[string]any{
			{
				"request_id":            "session:default-cost",
				"app_type":              "claude",
				"model":                 "Unknown-New-Model",
				"input_tokens":          int64(1_000_000),
				"output_tokens":         int64(1_000_000),
				"cache_read_tokens":     int64(0),
				"cache_creation_tokens": int64(0),
				"created_at":            int64(1743840000),
				"session_id":            "session-1",
				"data_source":           "session_log",
			},
		},
	}))
	if response.AcceptedCount != 1 || response.DuplicateCount != 0 {
		t.Fatalf("unexpected sync response: %+v", response)
	}

	var stored entity.UsageReport
	if err := db.Where("client_id = ? AND request_id = ?", "client-1", "session:default-cost").First(&stored).Error; err != nil {
		t.Fatalf("load stored report returned error: %v", err)
	}

	if stored.Model != "unknown-new-model" {
		t.Fatalf("expected normalized unknown model, got %q", stored.Model)
	}
	if stored.PricingSource != "default" {
		t.Fatalf("expected default pricing source, got %q", stored.PricingSource)
	}
	assertDecimalEqual(t, stored.InputCostUSD, "0.657")
	assertDecimalEqual(t, stored.OutputCostUSD, "3.429")
	assertDecimalEqual(t, stored.TotalCostUSD, "4.086")
}

func newTestSyncRouter(t *testing.T, db *gorm.DB) http.Handler {
	t.Helper()

	syncHandler := NewSyncHandler(service.NewSyncService(db))
	return NewRouter("secret-token", syncHandler.HandleSync, nil, nil)
}

func performSyncRequest(t *testing.T, router http.Handler, body map[string]any) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Marshal() returned error: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/sync", bytes.NewReader(payload))
	request.Header.Set("Authorization", "Bearer secret-token")
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	return recorder
}

func decodeSyncResponse(t *testing.T, recorder *httptest.ResponseRecorder) legacySyncResponse {
	t.Helper()

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response legacySyncResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}
	return response
}

func assertDecimalEqual(t *testing.T, actual string, expected string) {
	t.Helper()

	actualValue := new(big.Rat)
	if _, ok := actualValue.SetString(actual); !ok {
		t.Fatalf("invalid actual decimal %q", actual)
	}

	expectedValue := new(big.Rat)
	if _, ok := expectedValue.SetString(expected); !ok {
		t.Fatalf("invalid expected decimal %q", expected)
	}

	if actualValue.Cmp(expectedValue) != 0 {
		t.Fatalf("unexpected decimal value: actual=%q expected=%q", actual, expected)
	}
}
