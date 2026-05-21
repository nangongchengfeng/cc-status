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

func seedUsageReports(t *testing.T, db *gorm.DB, reports []entity.UsageReport) {
	t.Helper()

	for _, report := range reports {
		if err := db.Create(&report).Error; err != nil {
			t.Fatalf("Create() returned error: %v", err)
		}
	}
}

func newTestStatsRouter(t *testing.T, db *gorm.DB) http.Handler {
	t.Helper()

	syncHandler := NewSyncHandler(service.NewSyncService(db))
	modelPricingHandler := NewModelPricingHandler(service.NewModelPricingService(db))
	statsHandler := NewStatsHandler(service.NewStatsService(db))
	return NewRouter("secret-token", syncHandler.HandleSync, modelPricingHandler, statsHandler)
}
