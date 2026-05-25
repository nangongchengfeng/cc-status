package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 测试前准备：创建临时的测试 index.html
func setupTestDist(t *testing.T) (cleanup func()) {
	t.Helper()

	// 创建 internal/ui/dist 目录
	distDir := "internal/ui/dist"
	if err := os.MkdirAll(distDir, 0755); err != nil {
		t.Fatalf("failed to create dist dir: %v", err)
	}

	// 创建测试用的 index.html
	indexPath := filepath.Join(distDir, "index.html")
	indexContent := "<html><body>Test Index</body></html>"
	if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
		t.Fatalf("failed to write test index.html: %v", err)
	}

	// 返回清理函数
	return func() {
		os.RemoveAll("internal/ui")
	}
}

func TestUIRouteServesIndexHTML(t *testing.T) {
	t.Parallel()

	cleanup := setupTestDist(t)
	defer cleanup()

	router := NewRouter("secret-token", nil, nil, nil, nil)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "Test Index") {
		t.Fatalf("expected index.html content, got %q", recorder.Body.String())
	}
}

func TestUIRouteSPAFallback(t *testing.T) {
	t.Parallel()

	cleanup := setupTestDist(t)
	defer cleanup()

	router := NewRouter("secret-token", nil, nil, nil, nil)
	request := httptest.NewRequest(http.MethodGet, "/some/spa/route", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "Test Index") {
		t.Fatalf("expected index.html for SPA route, got %q", recorder.Body.String())
	}
}

func TestAPIRouteNotBlockedByUI(t *testing.T) {
	t.Parallel()

	cleanup := setupTestDist(t)
	defer cleanup()

	router := NewRouter("secret-token", nil, nil, nil, nil)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	request.Header.Set("Authorization", "Bearer secret-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "pong") {
		t.Fatalf("expected ping response, got %q", recorder.Body.String())
	}
}

func TestHealthzNotBlockedByUI(t *testing.T) {
	t.Parallel()

	cleanup := setupTestDist(t)
	defer cleanup()

	router := NewRouter("secret-token", nil, nil, nil, nil)
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}
