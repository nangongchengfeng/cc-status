package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"cc-status/server/internal/model/entity"
	"cc-status/server/internal/repository"
	"cc-status/server/internal/service"

	"gorm.io/gorm"
)

func TestModelPricingRouteListsAllPricings(t *testing.T) {
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

	customPricing := entity.ModelPricing{
		ModelID:                     "claude-custom-1",
		DisplayName:                 "Claude Custom 1",
		InputCostPerMillion:         "1",
		OutputCostPerMillion:        "2",
		CacheReadCostPerMillion:     "0",
		CacheCreationCostPerMillion: "0",
	}
	if err := db.Create(&customPricing).Error; err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	router := newTestModelPricingRouter(t, db)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/model-pricings", nil)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Data []struct {
			ModelID       string `json:"model_id"`
			IsPlaceholder bool   `json:"is_placeholder"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}
	if len(response.Data) == 0 {
		t.Fatal("expected non-empty pricing list")
	}

	var foundCustom bool
	var foundPlaceholder bool
	for _, item := range response.Data {
		if item.ModelID == "claude-custom-1" {
			foundCustom = true
		}
		if item.ModelID == "__default__" && item.IsPlaceholder {
			foundPlaceholder = true
		}
	}
	if !foundCustom {
		t.Fatalf("expected custom pricing in response: %+v", response.Data)
	}
	if !foundPlaceholder {
		t.Fatalf("expected placeholder pricing in response: %+v", response.Data)
	}
}

func TestModelPricingRouteCreatesPricing(t *testing.T) {
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

	router := newTestModelPricingRouter(t, db)
	body := []byte(`{"model_id":"Claude-Custom-2","display_name":"Claude Custom 2","input_cost_per_million":"2","output_cost_per_million":"8","cache_read_cost_per_million":"0.2","cache_creation_cost_per_million":"1.5","is_placeholder":false}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/model-pricings", bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer secret-token")
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var stored entity.ModelPricing
	if err := db.Where("model_id = ?", "claude-custom-2").First(&stored).Error; err != nil {
		t.Fatalf("load stored pricing returned error: %v", err)
	}
	if stored.DisplayName != "Claude Custom 2" {
		t.Fatalf("unexpected display name: %q", stored.DisplayName)
	}
	if stored.IsPlaceholder {
		t.Fatal("expected normal pricing")
	}
}

func TestModelPricingRouteRejectsSecondPlaceholderOnCreate(t *testing.T) {
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

	router := newTestModelPricingRouter(t, db)
	body := []byte(`{"model_id":"__fallback__","display_name":"Another Default","input_cost_per_million":"1","output_cost_per_million":"2","cache_read_cost_per_million":"0","cache_creation_cost_per_million":"0","is_placeholder":true}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/model-pricings", bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer secret-token")
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d with body %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "CONFLICT") {
		t.Fatalf("expected conflict body, got %s", recorder.Body.String())
	}
}

func TestModelPricingRouteUpdatesPricingWithFullPayload(t *testing.T) {
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

	existing := entity.ModelPricing{
		ModelID:                     "claude-edit-me",
		DisplayName:                 "Before Update",
		InputCostPerMillion:         "1",
		OutputCostPerMillion:        "2",
		CacheReadCostPerMillion:     "0.1",
		CacheCreationCostPerMillion: "0.2",
		IsPlaceholder:               false,
	}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	router := newTestModelPricingRouter(t, db)
	body := []byte(`{"model_id":"CLAUDE-EDIT-ME","display_name":"After Update","input_cost_per_million":"6","output_cost_per_million":"18","cache_read_cost_per_million":"0.6","cache_creation_cost_per_million":"2.5","is_placeholder":false}`)
	request := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/model-pricings/"+strconv.FormatUint(uint64(existing.ID), 10),
		bytes.NewReader(body),
	)
	request.Header.Set("Authorization", "Bearer secret-token")
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var updated entity.ModelPricing
	if err := db.First(&updated, existing.ID).Error; err != nil {
		t.Fatalf("load updated pricing returned error: %v", err)
	}
	if updated.ModelID != "claude-edit-me" {
		t.Fatalf("expected normalized model id, got %q", updated.ModelID)
	}
	if updated.DisplayName != "After Update" {
		t.Fatalf("unexpected display name: %q", updated.DisplayName)
	}
	if updated.InputCostPerMillion != "6" || updated.OutputCostPerMillion != "18" {
		t.Fatalf("unexpected updated costs: %+v", updated)
	}
}

func TestModelPricingRouteRejectsPartialUpdatePayload(t *testing.T) {
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

	existing := entity.ModelPricing{
		ModelID:                     "claude-partial-update",
		DisplayName:                 "Before Update",
		InputCostPerMillion:         "1",
		OutputCostPerMillion:        "2",
		CacheReadCostPerMillion:     "0.1",
		CacheCreationCostPerMillion: "0.2",
	}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	router := newTestModelPricingRouter(t, db)
	body := []byte(`{"model_id":"claude-partial-update","display_name":"Only Partial"}`)
	request := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/model-pricings/"+strconv.FormatUint(uint64(existing.ID), 10),
		bytes.NewReader(body),
	)
	request.Header.Set("Authorization", "Bearer secret-token")
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d with body %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "INVALID_REQUEST") {
		t.Fatalf("expected invalid request body, got %s", recorder.Body.String())
	}
}

func TestModelPricingRouteRejectsPlaceholderConflictOnUpdate(t *testing.T) {
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

	existing := entity.ModelPricing{
		ModelID:                     "claude-update-conflict",
		DisplayName:                 "Conflict Candidate",
		InputCostPerMillion:         "1",
		OutputCostPerMillion:        "2",
		CacheReadCostPerMillion:     "0.1",
		CacheCreationCostPerMillion: "0.2",
		IsPlaceholder:               false,
	}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	router := newTestModelPricingRouter(t, db)
	body := []byte(`{"model_id":"claude-update-conflict","display_name":"Conflict Candidate","input_cost_per_million":"1","output_cost_per_million":"2","cache_read_cost_per_million":"0.1","cache_creation_cost_per_million":"0.2","is_placeholder":true}`)
	request := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/model-pricings/"+strconv.FormatUint(uint64(existing.ID), 10),
		bytes.NewReader(body),
	)
	request.Header.Set("Authorization", "Bearer secret-token")
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d with body %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "CONFLICT") {
		t.Fatalf("expected conflict body, got %s", recorder.Body.String())
	}
}

func newTestModelPricingRouter(t *testing.T, db *gorm.DB) http.Handler {
	t.Helper()

	syncHandler := NewSyncHandler(service.NewSyncService(db))
	modelPricingHandler := NewModelPricingHandler(service.NewModelPricingService(db))
	return NewRouter("secret-token", syncHandler.HandleSync, modelPricingHandler)
}
