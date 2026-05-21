package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

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

func newTestSyncRouter(t *testing.T, db *gorm.DB) http.Handler {
	t.Helper()

	syncHandler := NewSyncHandler(service.NewSyncService(db))
	return NewRouter("secret-token", syncHandler.HandleSync)
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
