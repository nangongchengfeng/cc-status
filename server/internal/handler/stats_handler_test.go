package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"cc-status/server/internal/model/entity"
	"cc-status/server/internal/repository"
	"cc-status/server/internal/service"

	"gorm.io/gorm"
)

func TestStatsOverviewRouteReturnsAggregatedOverview(t *testing.T) {
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
		{ClientID: "client-b", RequestID: "req-1", AppType: "claude", Model: "zeta", InputTokens: 50, OutputTokens: 50, PricingSource: "exact", CreatedAtUnix: 1743840000, DataSource: "session_log", TotalCostUSD: "10", InputCostUSD: "4", OutputCostUSD: "6"},
		{ClientID: "client-a", RequestID: "req-2", AppType: "claude", Model: "alpha", InputTokens: 30, OutputTokens: 70, PricingSource: "exact", CreatedAtUnix: 1743843600, DataSource: "session_log", TotalCostUSD: "20", InputCostUSD: "8", OutputCostUSD: "12"},
		{ClientID: "client-c", RequestID: "req-3", AppType: "claude", Model: "beta", InputTokens: 20, OutputTokens: 10, CacheReadTokens: 5, CacheCreationTokens: 5, PricingSource: "prefix", CreatedAtUnix: 1743847200, DataSource: "session_log", TotalCostUSD: "5", InputCostUSD: "2", OutputCostUSD: "2", CacheReadCostUSD: "0.5", CacheCreationCostUSD: "0.5"},
		{ClientID: "client-a", RequestID: "req-4", AppType: "claude", Model: "gamma", InputTokens: 5, OutputTokens: 5, PricingSource: "default", CreatedAtUnix: 1743850800, DataSource: "session_log", TotalCostUSD: "1", InputCostUSD: "0.4", OutputCostUSD: "0.6"},
		{ClientID: "client-d", RequestID: "req-5", AppType: "claude", Model: "delta", InputTokens: 7, OutputTokens: 8, PricingSource: "exact", CreatedAtUnix: 1743854400, DataSource: "session_log", TotalCostUSD: "12", InputCostUSD: "4", OutputCostUSD: "8"},
		{ClientID: "client-e", RequestID: "req-6", AppType: "claude", Model: "epsilon", InputTokens: 1, OutputTokens: 1, PricingSource: "exact", CreatedAtUnix: 1743858000, DataSource: "session_log", TotalCostUSD: "0.5", InputCostUSD: "0.2", OutputCostUSD: "0.3"},
		{ClientID: "client-f", RequestID: "req-7", AppType: "claude", Model: "eta", InputTokens: 1, OutputTokens: 0, PricingSource: "exact", CreatedAtUnix: 1743861600, DataSource: "session_log", TotalCostUSD: "12", InputCostUSD: "12"},
	})

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/stats/overview", nil)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Data struct {
			TotalTokens   int64  `json:"total_tokens"`
			TotalRequests int64  `json:"total_requests"`
			ActiveClients int64  `json:"active_clients"`
			TotalCostUSD  string `json:"total_cost_usd"`
			TopModels     []struct {
				Model  string `json:"model"`
				Tokens int64  `json:"tokens"`
			} `json:"top_models"`
			TopClients []struct {
				ClientID     string `json:"client_id"`
				TotalCostUSD string `json:"total_cost_usd"`
			} `json:"top_clients"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if response.Data.TotalTokens != 268 {
		t.Fatalf("unexpected total tokens: %d", response.Data.TotalTokens)
	}
	if response.Data.TotalRequests != 7 {
		t.Fatalf("unexpected total requests: %d", response.Data.TotalRequests)
	}
	if response.Data.ActiveClients != 6 {
		t.Fatalf("unexpected active clients: %d", response.Data.ActiveClients)
	}
	assertDecimalEqual(t, response.Data.TotalCostUSD, "60.5")

	if len(response.Data.TopModels) != 5 {
		t.Fatalf("expected top 5 models, got %d", len(response.Data.TopModels))
	}
	expectedModels := []struct {
		model  string
		tokens int64
	}{
		{model: "alpha", tokens: 100},
		{model: "zeta", tokens: 100},
		{model: "beta", tokens: 40},
		{model: "delta", tokens: 15},
		{model: "gamma", tokens: 10},
	}
	for index, expected := range expectedModels {
		actual := response.Data.TopModels[index]
		if actual.Model != expected.model || actual.Tokens != expected.tokens {
			t.Fatalf("unexpected top model at %d: %+v", index, actual)
		}
	}

	if len(response.Data.TopClients) != 5 {
		t.Fatalf("expected top 5 clients, got %d", len(response.Data.TopClients))
	}
	expectedClients := []struct {
		clientID string
		total    string
	}{
		{clientID: "client-a", total: "21"},
		{clientID: "client-d", total: "12"},
		{clientID: "client-f", total: "12"},
		{clientID: "client-b", total: "10"},
		{clientID: "client-c", total: "5"},
	}
	for index, expected := range expectedClients {
		actual := response.Data.TopClients[index]
		if actual.ClientID != expected.clientID {
			t.Fatalf("unexpected top client at %d: %+v", index, actual)
		}
		assertDecimalEqual(t, actual.TotalCostUSD, expected.total)
	}
}

func TestStatsOverviewRouteRequiresAuth(t *testing.T) {
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

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/stats/overview", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d with body %s", recorder.Code, recorder.Body.String())
	}
}

func TestStatsTrendRouteReturnsShanghaiBucketsWithZeroFill(t *testing.T) {
	t.Parallel()

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("LoadLocation() returned error: %v", err)
	}

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
			ClientID:      "client-1",
			RequestID:     "trend-1",
			AppType:       "claude",
			Model:         "alpha",
			InputTokens:   10,
			OutputTokens:  5,
			PricingSource: "exact",
			CreatedAtUnix: time.Date(2026, 5, 21, 10, 15, 0, 0, location).Unix(),
			DataSource:    "session_log",
			TotalCostUSD:  "1.5",
		},
		{
			ClientID:      "client-2",
			RequestID:     "trend-2",
			AppType:       "claude",
			Model:         "beta",
			InputTokens:   20,
			OutputTokens:  10,
			PricingSource: "exact",
			CreatedAtUnix: time.Date(2026, 5, 21, 12, 5, 0, 0, location).Unix(),
			DataSource:    "session_log",
			TotalCostUSD:  "2.25",
		},
	})

	startAt := time.Date(2026, 5, 21, 10, 0, 0, 0, location).Unix()
	endAt := time.Date(2026, 5, 21, 12, 0, 0, 0, location).Unix()

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/stats/trend?interval=hour&start_at="+formatUnix(startAt)+"&end_at="+formatUnix(endAt),
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
			Bucket        string `json:"bucket"`
			TotalTokens   int64  `json:"total_tokens"`
			TotalRequests int64  `json:"total_requests"`
			TotalCostUSD  string `json:"total_cost_usd"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if len(response.Data) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(response.Data))
	}

	expected := []struct {
		bucket   string
		tokens   int64
		requests int64
		cost     string
	}{
		{bucket: "2026-05-21T10:00:00+08:00", tokens: 15, requests: 1, cost: "1.5"},
		{bucket: "2026-05-21T11:00:00+08:00", tokens: 0, requests: 0, cost: "0"},
		{bucket: "2026-05-21T12:00:00+08:00", tokens: 30, requests: 1, cost: "2.25"},
	}
	for index, expectedItem := range expected {
		actual := response.Data[index]
		if actual.Bucket != expectedItem.bucket || actual.TotalTokens != expectedItem.tokens || actual.TotalRequests != expectedItem.requests {
			t.Fatalf("unexpected trend bucket at %d: %+v", index, actual)
		}
		assertDecimalEqual(t, actual.TotalCostUSD, expectedItem.cost)
	}
}

func TestStatsTrendRouteSupportsDayInterval(t *testing.T) {
	t.Parallel()

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("LoadLocation() returned error: %v", err)
	}

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
			ClientID:      "client-1",
			RequestID:     "trend-day-1",
			AppType:       "claude",
			Model:         "alpha",
			InputTokens:   8,
			OutputTokens:  2,
			PricingSource: "exact",
			CreatedAtUnix: time.Date(2026, 5, 20, 9, 0, 0, 0, location).Unix(),
			DataSource:    "session_log",
			TotalCostUSD:  "1.2",
		},
		{
			ClientID:      "client-2",
			RequestID:     "trend-day-2",
			AppType:       "claude",
			Model:         "beta",
			InputTokens:   4,
			OutputTokens:  1,
			PricingSource: "exact",
			CreatedAtUnix: time.Date(2026, 5, 22, 20, 0, 0, 0, location).Unix(),
			DataSource:    "session_log",
			TotalCostUSD:  "0.8",
		},
	})

	startAt := time.Date(2026, 5, 20, 0, 0, 0, 0, location).Unix()
	endAt := time.Date(2026, 5, 22, 0, 0, 0, 0, location).Unix()
	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/stats/trend?interval=day&start_at="+formatUnix(startAt)+"&end_at="+formatUnix(endAt),
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
			Bucket        string `json:"bucket"`
			TotalTokens   int64  `json:"total_tokens"`
			TotalRequests int64  `json:"total_requests"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if len(response.Data) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(response.Data))
	}
	if response.Data[0].Bucket != "2026-05-20T00:00:00+08:00" || response.Data[0].TotalTokens != 10 {
		t.Fatalf("unexpected first day bucket: %+v", response.Data[0])
	}
	if response.Data[1].Bucket != "2026-05-21T00:00:00+08:00" || response.Data[1].TotalTokens != 0 {
		t.Fatalf("unexpected zero-filled day bucket: %+v", response.Data[1])
	}
	if response.Data[2].Bucket != "2026-05-22T00:00:00+08:00" || response.Data[2].TotalTokens != 5 {
		t.Fatalf("unexpected last day bucket: %+v", response.Data[2])
	}
}

func TestStatsTrendRouteRejectsInvalidInterval(t *testing.T) {
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

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/stats/trend?interval=week", nil)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d with body %s", recorder.Code, recorder.Body.String())
	}
}

func TestStatsDashboardRouteReturnsContractSkeleton(t *testing.T) {
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

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("LoadLocation() returned error: %v", err)
	}

	startAt := time.Date(2026, 5, 21, 0, 0, 0, 0, location).Unix()
	endAt := time.Date(2026, 5, 21, 23, 0, 0, 0, location).Unix()

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/stats/dashboard?interval=day&start_at="+formatUnix(startAt)+"&end_at="+formatUnix(endAt),
		nil,
	)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Data struct {
			Overview struct {
				TotalTokens   int64  `json:"total_tokens"`
				TotalCostUSD  string `json:"total_cost_usd"`
				TotalRequests int64  `json:"total_requests"`
				ActiveClients int64  `json:"active_clients"`
			} `json:"overview"`
			Trend []struct {
				Bucket              string `json:"bucket"`
				InputTokens         int64  `json:"input_tokens"`
				OutputTokens        int64  `json:"output_tokens"`
				CacheReadTokens     int64  `json:"cache_read_tokens"`
				CacheCreationTokens int64  `json:"cache_creation_tokens"`
				TotalRequests       int64  `json:"total_requests"`
				TotalCostUSD        string `json:"total_cost_usd"`
			} `json:"trend"`
			TopModels []struct {
				Model       string `json:"model"`
				DisplayName string `json:"display_name"`
				TotalTokens int64  `json:"total_tokens"`
			} `json:"top_models"`
			TopClients []struct {
				ClientID     string `json:"client_id"`
				TotalCostUSD string `json:"total_cost_usd"`
			} `json:"top_clients"`
			CacheAnalysis struct {
				SavedCostUSD         string `json:"saved_cost_usd"`
				CacheReadCostUSD     string `json:"cache_read_cost_usd"`
				CacheCreationCostUSD string `json:"cache_creation_cost_usd"`
			} `json:"cache_analysis"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if response.Data.Overview.TotalTokens != 0 || response.Data.Overview.TotalRequests != 0 || response.Data.Overview.ActiveClients != 0 {
		t.Fatalf("unexpected overview skeleton: %+v", response.Data.Overview)
	}
	assertDecimalEqual(t, response.Data.Overview.TotalCostUSD, "0")
	if len(response.Data.Trend) != 0 || len(response.Data.TopModels) != 0 || len(response.Data.TopClients) != 0 {
		t.Fatalf("unexpected non-empty dashboard slices: %+v", response.Data)
	}
	assertDecimalEqual(t, response.Data.CacheAnalysis.SavedCostUSD, "0")
	assertDecimalEqual(t, response.Data.CacheAnalysis.CacheReadCostUSD, "0")
	assertDecimalEqual(t, response.Data.CacheAnalysis.CacheCreationCostUSD, "0")
}

func TestStatsDashboardRouteRejectsInvalidQuery(t *testing.T) {
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

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/stats/dashboard?interval=week&start_at=1743840000&end_at=1743926400",
		nil,
	)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d with body %s", recorder.Code, recorder.Body.String())
	}
}

func TestStatsDashboardRouteRequiresExplicitRange(t *testing.T) {
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

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/stats/dashboard?interval=day&end_at=1743926400",
		nil,
	)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d with body %s", recorder.Code, recorder.Body.String())
	}
}

func TestStatsDashboardRouteRequiresAuth(t *testing.T) {
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

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/stats/dashboard?interval=day&start_at=1743840000&end_at=1743926400",
		nil,
	)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d with body %s", recorder.Code, recorder.Body.String())
	}
}

func TestStatsDashboardRouteReturnsOverviewAndUnifiedTrend(t *testing.T) {
	t.Parallel()

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("LoadLocation() returned error: %v", err)
	}

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
			ClientID:             "client-1",
			RequestID:            "dashboard-1",
			AppType:              "claude",
			Model:                "alpha",
			InputTokens:          10,
			OutputTokens:         5,
			PricingSource:        "exact",
			CreatedAtUnix:        time.Date(2026, 5, 21, 10, 15, 0, 0, location).Unix(),
			DataSource:           "session_log",
			TotalCostUSD:         "1.5",
			InputCostUSD:         "0.6",
			OutputCostUSD:        "0.9",
			CacheReadCostUSD:     "0",
			CacheCreationCostUSD: "0",
		},
		{
			ClientID:             "client-2",
			RequestID:            "dashboard-2",
			AppType:              "claude",
			Model:                "beta",
			InputTokens:          20,
			OutputTokens:         10,
			CacheReadTokens:      7,
			CacheCreationTokens:  3,
			PricingSource:        "exact",
			CreatedAtUnix:        time.Date(2026, 5, 21, 12, 5, 0, 0, location).Unix(),
			DataSource:           "session_log",
			TotalCostUSD:         "2.25",
			InputCostUSD:         "1",
			OutputCostUSD:        "1",
			CacheReadCostUSD:     "0.15",
			CacheCreationCostUSD: "0.10",
		},
		{
			ClientID:             "client-2",
			RequestID:            "dashboard-3",
			AppType:              "claude",
			Model:                "beta",
			InputTokens:          2,
			OutputTokens:         3,
			CacheReadTokens:      0,
			CacheCreationTokens:  1,
			PricingSource:        "exact",
			CreatedAtUnix:        time.Date(2026, 5, 21, 12, 30, 0, 0, location).Unix(),
			DataSource:           "session_log",
			TotalCostUSD:         "0.75",
			InputCostUSD:         "0.25",
			OutputCostUSD:        "0.4",
			CacheReadCostUSD:     "0",
			CacheCreationCostUSD: "0.1",
		},
		{
			ClientID:      "client-3",
			RequestID:     "dashboard-outside",
			AppType:       "claude",
			Model:         "gamma",
			InputTokens:   99,
			OutputTokens:  1,
			PricingSource: "exact",
			CreatedAtUnix: time.Date(2026, 5, 21, 13, 5, 0, 0, location).Unix(),
			DataSource:    "session_log",
			TotalCostUSD:  "9.9",
			InputCostUSD:  "4.5",
			OutputCostUSD: "5.4",
		},
	})

	startAt := time.Date(2026, 5, 21, 10, 0, 0, 0, location).Unix()
	endAt := time.Date(2026, 5, 21, 12, 0, 0, 0, location).Unix()

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/stats/dashboard?interval=hour&start_at="+formatUnix(startAt)+"&end_at="+formatUnix(endAt),
		nil,
	)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Data struct {
			Overview struct {
				TotalTokens   int64  `json:"total_tokens"`
				TotalCostUSD  string `json:"total_cost_usd"`
				TotalRequests int64  `json:"total_requests"`
				ActiveClients int64  `json:"active_clients"`
			} `json:"overview"`
			Trend []struct {
				Bucket              string `json:"bucket"`
				InputTokens         int64  `json:"input_tokens"`
				OutputTokens        int64  `json:"output_tokens"`
				CacheReadTokens     int64  `json:"cache_read_tokens"`
				CacheCreationTokens int64  `json:"cache_creation_tokens"`
				TotalRequests       int64  `json:"total_requests"`
				TotalCostUSD        string `json:"total_cost_usd"`
			} `json:"trend"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if response.Data.Overview.TotalTokens != 61 {
		t.Fatalf("unexpected total tokens: %d", response.Data.Overview.TotalTokens)
	}
	if response.Data.Overview.TotalRequests != 3 {
		t.Fatalf("unexpected total requests: %d", response.Data.Overview.TotalRequests)
	}
	if response.Data.Overview.ActiveClients != 2 {
		t.Fatalf("unexpected active clients: %d", response.Data.Overview.ActiveClients)
	}
	assertDecimalEqual(t, response.Data.Overview.TotalCostUSD, "4.5")

	if len(response.Data.Trend) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(response.Data.Trend))
	}

	expected := []struct {
		bucket        string
		inputTokens   int64
		outputTokens  int64
		cacheRead     int64
		cacheCreation int64
		totalRequests int64
		totalCostUSD  string
	}{
		{bucket: "2026-05-21T10:00:00+08:00", inputTokens: 10, outputTokens: 5, cacheRead: 0, cacheCreation: 0, totalRequests: 1, totalCostUSD: "1.5"},
		{bucket: "2026-05-21T11:00:00+08:00", inputTokens: 0, outputTokens: 0, cacheRead: 0, cacheCreation: 0, totalRequests: 0, totalCostUSD: "0"},
		{bucket: "2026-05-21T12:00:00+08:00", inputTokens: 22, outputTokens: 13, cacheRead: 7, cacheCreation: 4, totalRequests: 2, totalCostUSD: "3"},
	}
	for index, expectedPoint := range expected {
		actual := response.Data.Trend[index]
		if actual.Bucket != expectedPoint.bucket ||
			actual.InputTokens != expectedPoint.inputTokens ||
			actual.OutputTokens != expectedPoint.outputTokens ||
			actual.CacheReadTokens != expectedPoint.cacheRead ||
			actual.CacheCreationTokens != expectedPoint.cacheCreation ||
			actual.TotalRequests != expectedPoint.totalRequests {
			t.Fatalf("unexpected dashboard trend at %d: %+v", index, actual)
		}
		assertDecimalEqual(t, actual.TotalCostUSD, expectedPoint.totalCostUSD)
	}
}

func TestStatsDashboardRouteSupportsDayInterval(t *testing.T) {
	t.Parallel()

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("LoadLocation() returned error: %v", err)
	}

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
			ClientID:      "client-1",
			RequestID:     "dashboard-day-1",
			AppType:       "claude",
			Model:         "alpha",
			InputTokens:   8,
			OutputTokens:  2,
			PricingSource: "exact",
			CreatedAtUnix: time.Date(2026, 5, 20, 9, 0, 0, 0, location).Unix(),
			DataSource:    "session_log",
			TotalCostUSD:  "1.2",
		},
		{
			ClientID:      "client-2",
			RequestID:     "dashboard-day-2",
			AppType:       "claude",
			Model:         "beta",
			InputTokens:   4,
			OutputTokens:  1,
			PricingSource: "exact",
			CreatedAtUnix: time.Date(2026, 5, 22, 20, 0, 0, 0, location).Unix(),
			DataSource:    "session_log",
			TotalCostUSD:  "0.8",
		},
	})

	startAt := time.Date(2026, 5, 20, 0, 0, 0, 0, location).Unix()
	endAt := time.Date(2026, 5, 22, 0, 0, 0, 0, location).Unix()

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/stats/dashboard?interval=day&start_at="+formatUnix(startAt)+"&end_at="+formatUnix(endAt),
		nil,
	)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Data struct {
			Overview struct {
				TotalTokens int64 `json:"total_tokens"`
			} `json:"overview"`
			Trend []struct {
				Bucket       string `json:"bucket"`
				InputTokens  int64  `json:"input_tokens"`
				OutputTokens int64  `json:"output_tokens"`
			} `json:"trend"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if response.Data.Overview.TotalTokens != 15 {
		t.Fatalf("unexpected total tokens: %d", response.Data.Overview.TotalTokens)
	}
	if len(response.Data.Trend) != 3 {
		t.Fatalf("expected 3 day buckets, got %d", len(response.Data.Trend))
	}
	if response.Data.Trend[0].Bucket != "2026-05-20T00:00:00+08:00" || response.Data.Trend[0].InputTokens != 8 {
		t.Fatalf("unexpected first day bucket: %+v", response.Data.Trend[0])
	}
	if response.Data.Trend[1].Bucket != "2026-05-21T00:00:00+08:00" || response.Data.Trend[1].InputTokens != 0 || response.Data.Trend[1].OutputTokens != 0 {
		t.Fatalf("unexpected zero-filled day bucket: %+v", response.Data.Trend[1])
	}
	if response.Data.Trend[2].Bucket != "2026-05-22T00:00:00+08:00" || response.Data.Trend[2].InputTokens != 4 {
		t.Fatalf("unexpected last day bucket: %+v", response.Data.Trend[2])
	}
}

func TestStatsDashboardRouteReturnsTopModelsWithDisplayName(t *testing.T) {
	t.Parallel()

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("LoadLocation() returned error: %v", err)
	}

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

	seedModelPricings(t, db, []entity.ModelPricing{
		{
			ModelID:                     "alpha",
			DisplayName:                 "Alpha Display",
			InputCostPerMillion:         "1",
			OutputCostPerMillion:        "1",
			CacheReadCostPerMillion:     "0.1",
			CacheCreationCostPerMillion: "0.2",
		},
		{
			ModelID:                     "gamma",
			DisplayName:                 "Gamma Display",
			InputCostPerMillion:         "1",
			OutputCostPerMillion:        "1",
			CacheReadCostPerMillion:     "0.1",
			CacheCreationCostPerMillion: "0.2",
		},
	})

	seedUsageReports(t, db, []entity.UsageReport{
		{
			ClientID:            "client-1",
			RequestID:           "model-rank-1",
			AppType:             "claude",
			Model:               "alpha",
			InputTokens:         30,
			OutputTokens:        10,
			CacheReadTokens:     5,
			CacheCreationTokens: 5,
			PricingSource:       "exact",
			CreatedAtUnix:       time.Date(2026, 5, 21, 10, 0, 0, 0, location).Unix(),
			DataSource:          "session_log",
			TotalCostUSD:        "2",
		},
		{
			ClientID:            "client-2",
			RequestID:           "model-rank-2",
			AppType:             "claude",
			Model:               "beta",
			InputTokens:         20,
			OutputTokens:        10,
			CacheReadTokens:     10,
			CacheCreationTokens: 5,
			PricingSource:       "exact",
			CreatedAtUnix:       time.Date(2026, 5, 21, 11, 0, 0, 0, location).Unix(),
			DataSource:          "session_log",
			TotalCostUSD:        "1.5",
		},
		{
			ClientID:            "client-3",
			RequestID:           "model-rank-3",
			AppType:             "claude",
			Model:               "gamma",
			InputTokens:         20,
			OutputTokens:        10,
			CacheReadTokens:     10,
			CacheCreationTokens: 5,
			PricingSource:       "exact",
			CreatedAtUnix:       time.Date(2026, 5, 21, 12, 0, 0, 0, location).Unix(),
			DataSource:          "session_log",
			TotalCostUSD:        "1.4",
		},
	})

	startAt := time.Date(2026, 5, 21, 10, 0, 0, 0, location).Unix()
	endAt := time.Date(2026, 5, 21, 12, 0, 0, 0, location).Unix()

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/stats/dashboard?interval=hour&start_at="+formatUnix(startAt)+"&end_at="+formatUnix(endAt),
		nil,
	)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Data struct {
			TopModels []struct {
				Model       string `json:"model"`
				DisplayName string `json:"display_name"`
				TotalTokens int64  `json:"total_tokens"`
			} `json:"top_models"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if len(response.Data.TopModels) != 3 {
		t.Fatalf("expected 3 top models, got %d", len(response.Data.TopModels))
	}

	expected := []struct {
		model       string
		displayName string
		totalTokens int64
	}{
		{model: "alpha", displayName: "Alpha Display", totalTokens: 50},
		{model: "beta", displayName: "", totalTokens: 45},
		{model: "gamma", displayName: "Gamma Display", totalTokens: 45},
	}
	for index, expectedModel := range expected {
		actual := response.Data.TopModels[index]
		if actual.Model != expectedModel.model || actual.DisplayName != expectedModel.displayName || actual.TotalTokens != expectedModel.totalTokens {
			t.Fatalf("unexpected top model at %d: %+v", index, actual)
		}
	}
}

func TestStatsDashboardRouteTopModelsUseStableTieBreak(t *testing.T) {
	t.Parallel()

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("LoadLocation() returned error: %v", err)
	}

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
			ClientID:      "client-1",
			RequestID:     "model-tie-1",
			AppType:       "claude",
			Model:         "zeta",
			InputTokens:   20,
			OutputTokens:  10,
			PricingSource: "exact",
			CreatedAtUnix: time.Date(2026, 5, 21, 10, 0, 0, 0, location).Unix(),
			DataSource:    "session_log",
			TotalCostUSD:  "1",
		},
		{
			ClientID:      "client-2",
			RequestID:     "model-tie-2",
			AppType:       "claude",
			Model:         "alpha",
			InputTokens:   15,
			OutputTokens:  15,
			PricingSource: "exact",
			CreatedAtUnix: time.Date(2026, 5, 21, 11, 0, 0, 0, location).Unix(),
			DataSource:    "session_log",
			TotalCostUSD:  "1",
		},
	})

	startAt := time.Date(2026, 5, 21, 10, 0, 0, 0, location).Unix()
	endAt := time.Date(2026, 5, 21, 11, 0, 0, 0, location).Unix()

	router := newTestStatsRouter(t, db)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/stats/dashboard?interval=hour&start_at="+formatUnix(startAt)+"&end_at="+formatUnix(endAt),
		nil,
	)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Data struct {
			TopModels []struct {
				Model       string `json:"model"`
				TotalTokens int64  `json:"total_tokens"`
			} `json:"top_models"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}

	if len(response.Data.TopModels) != 2 {
		t.Fatalf("expected 2 top models, got %d", len(response.Data.TopModels))
	}
	if response.Data.TopModels[0].Model != "alpha" || response.Data.TopModels[1].Model != "zeta" {
		t.Fatalf("unexpected stable order: %+v", response.Data.TopModels)
	}
}

func seedUsageReports(t *testing.T, db *gorm.DB, reports []entity.UsageReport) {
	t.Helper()

	for _, report := range reports {
		if err := db.Create(&report).Error; err != nil {
			t.Fatalf("Create() returned error: %v", err)
		}
	}
}

func seedModelPricings(t *testing.T, db *gorm.DB, pricings []entity.ModelPricing) {
	t.Helper()

	for _, pricing := range pricings {
		if err := db.Create(&pricing).Error; err != nil {
			t.Fatalf("Create() returned error: %v", err)
		}
	}
}

func newTestStatsRouter(t *testing.T, db *gorm.DB) http.Handler {
	t.Helper()

	syncHandler := NewSyncHandler(service.NewSyncService(db))
	modelPricingHandler := NewModelPricingHandler(service.NewModelPricingService(db))
	statsHandler := NewStatsHandler(service.NewStatsService(db))
	return NewRouter("secret-token", syncHandler.HandleSync, modelPricingHandler, statsHandler, nil)
}

func formatUnix(value int64) string {
	return strconv.FormatInt(value, 10)
}
