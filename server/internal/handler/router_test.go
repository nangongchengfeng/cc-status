package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouterExposesHealthzWithoutAuth(t *testing.T) {
	t.Parallel()

	router := NewRouter("secret-token", nil, nil, nil, nil)
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "ok") {
		t.Fatalf("expected health body to contain ok, got %q", recorder.Body.String())
	}
}

func TestRouterProtectsAPIV1Routes(t *testing.T) {
	t.Parallel()

	router := NewRouter("secret-token", nil, nil, nil, nil)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}
