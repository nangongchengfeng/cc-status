package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"cc-status/server/internal/model/entity"
	"cc-status/server/internal/repository"
	"cc-status/server/internal/service"

	"gorm.io/gorm"
)

func TestLogsRouteReturnsDefaultPaginationAndOrdering(t *testing.T) {
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

	seedUsageReports(t, db, []entity.UsageReport{
		{
			ClientID:             "client-a",
			RequestID:            "req-old",
			AppType:              "claude",
			Model:                "alpha",
			InputTokens:          1,
			OutputTokens:         2,
			CacheReadTokens:      3,
			CacheCreationTokens:  4,
			InputCostUSD:         "1.1",
			OutputCostUSD:        "2.2",
			CacheReadCostUSD:     "0.3",
			CacheCreationCostUSD: "0.4",
			TotalCostUSD:         "4.0",
			SessionID:            "session-old",
			PricingSource:        "exact",
			CreatedAtUnix:        1743840000,
			DataSource:           "session_log",
		},
		{
			ClientID:             "client-b",
			RequestID:            "req-same-1",
			AppType:              "claude",
			Model:                "beta",
			InputTokens:          10,
			OutputTokens:         20,
			CacheReadTokens:      0,
			CacheCreationTokens:  0,
			InputCostUSD:         "1",
			OutputCostUSD:        "2",
			CacheReadCostUSD:     "0",
			CacheCreationCostUSD: "0",
			TotalCostUSD:         "3",
			SessionID:            "session-same-1",
			PricingSource:        "prefix",
			CreatedAtUnix:        1743850000,
			DataSource:           "session_log",
		},
		{
			ClientID:             "client-c",
			RequestID:            "req-same-2",
			AppType:              "claude",
			Model:                "gamma",
			InputTokens:          100,
			OutputTokens:         200,
			CacheReadTokens:      1,
			CacheCreationTokens:  1,
			InputCostUSD:         "4",
			OutputCostUSD:        "5",
			CacheReadCostUSD:     "0.1",
			CacheCreationCostUSD: "0.2",
			TotalCostUSD:         "9.3",
			SessionID:            "session-same-2",
			PricingSource:        "default",
			CreatedAtUnix:        1743850000,
			DataSource:           "session_log",
		},
	})

	router := newTestLogsRouter(t, db)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/logs", nil)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Data []struct {
			RequestID            string `json:"request_id"`
			InputTokens          int64  `json:"input_tokens"`
			OutputTokens         int64  `json:"output_tokens"`
			CacheReadTokens      int64  `json:"cache_read_tokens"`
			CacheCreationTokens  int64  `json:"cache_creation_tokens"`
			InputCostUSD         string `json:"input_cost_usd"`
			OutputCostUSD        string `json:"output_cost_usd"`
			CacheReadCostUSD     string `json:"cache_read_cost_usd"`
			CacheCreationCostUSD string `json:"cache_creation_cost_usd"`
			TotalCostUSD         string `json:"total_cost_usd"`
			PricingSource        string `json:"pricing_source"`
		} `json:"data"`
		Total  int64 `json:"total"`
		Offset int   `json:"offset"`
		Limit  int   `json:"limit"`
		Page   int   `json:"page"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if response.Total != 3 || response.Offset != 0 || response.Limit != 20 || response.Page != 1 {
		t.Fatalf("unexpected pagination metadata: %+v", response)
	}
	if len(response.Data) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(response.Data))
	}
	if response.Data[0].RequestID != "req-same-2" || response.Data[1].RequestID != "req-same-1" || response.Data[2].RequestID != "req-old" {
		t.Fatalf("unexpected ordering: %+v", response.Data)
	}
	if response.Data[0].InputTokens != 100 || response.Data[0].PricingSource != "default" {
		t.Fatalf("expected full fields on first row: %+v", response.Data[0])
	}
	assertDecimalEqual(t, response.Data[0].TotalCostUSD, "9.3")
}

func TestLogsRouteAppliesFiltersAndOffsetPagination(t *testing.T) {
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

	seedUsageReports(t, db, []entity.UsageReport{
		{ClientID: "client-a", RequestID: "req-1", AppType: "claude", Model: "alpha", InputTokens: 1, OutputTokens: 1, PricingSource: "exact", CreatedAtUnix: 1743840000, DataSource: "session_log", TotalCostUSD: "1"},
		{ClientID: "client-a", RequestID: "req-2", AppType: "claude", Model: "alpha", InputTokens: 2, OutputTokens: 2, PricingSource: "exact", CreatedAtUnix: 1743843600, DataSource: "session_log", TotalCostUSD: "2"},
		{ClientID: "client-a", RequestID: "req-3", AppType: "claude", Model: "alpha", InputTokens: 3, OutputTokens: 3, PricingSource: "exact", CreatedAtUnix: 1743847200, DataSource: "session_log", TotalCostUSD: "3"},
		{ClientID: "client-b", RequestID: "req-4", AppType: "claude", Model: "beta", InputTokens: 4, OutputTokens: 4, PricingSource: "exact", CreatedAtUnix: 1743850800, DataSource: "session_log", TotalCostUSD: "4"},
	})

	router := newTestLogsRouter(t, db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/logs?client_id=client-a&model=ALPHA&start_time=1743840000&end_time=1743847200&offset=1&limit=1",
		nil,
	)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Data []struct {
			RequestID string `json:"request_id"`
		} `json:"data"`
		Total  int64 `json:"total"`
		Offset int   `json:"offset"`
		Limit  int   `json:"limit"`
		Page   int   `json:"page"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if response.Total != 3 || response.Offset != 1 || response.Limit != 1 || response.Page != 2 {
		t.Fatalf("unexpected pagination metadata: %+v", response)
	}
	if len(response.Data) != 1 || response.Data[0].RequestID != "req-2" {
		t.Fatalf("unexpected filtered rows: %+v", response.Data)
	}
}

func TestLogsRouteFiltersByRequestIDAndCapsLimit(t *testing.T) {
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

	seedUsageReports(t, db, []entity.UsageReport{
		{ClientID: "client-a", RequestID: "req-target", AppType: "claude", Model: "alpha", InputTokens: 1, OutputTokens: 1, PricingSource: "exact", CreatedAtUnix: 1743840000, DataSource: "session_log", TotalCostUSD: "1"},
		{ClientID: "client-b", RequestID: "req-other", AppType: "claude", Model: "beta", InputTokens: 2, OutputTokens: 2, PricingSource: "exact", CreatedAtUnix: 1743843600, DataSource: "session_log", TotalCostUSD: "2"},
	})

	router := newTestLogsRouter(t, db)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/logs?request_id=req-target&limit=999", nil)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Data []struct {
			RequestID string `json:"request_id"`
		} `json:"data"`
		Total int64 `json:"total"`
		Limit int   `json:"limit"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if response.Total != 1 || response.Limit != 100 {
		t.Fatalf("unexpected filtered metadata: %+v", response)
	}
	if len(response.Data) != 1 || response.Data[0].RequestID != "req-target" {
		t.Fatalf("unexpected request_id filter result: %+v", response.Data)
	}
}

func newTestLogsRouter(t *testing.T, db *gorm.DB) http.Handler {
	t.Helper()

	syncHandler := NewSyncHandler(service.NewSyncService(db))
	modelPricingHandler := NewModelPricingHandler(service.NewModelPricingService(db))
	statsHandler := NewStatsHandler(service.NewStatsService(db))
	logsHandler := NewLogsHandler(service.NewLogsService(db))
	return NewRouter("secret-token", syncHandler.HandleSync, modelPricingHandler, statsHandler, logsHandler)
}
