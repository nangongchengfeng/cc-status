package repository

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cc-status/server/internal/model/entity"
)

func TestOpenDatabaseCreatesParentDirAndFile(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "server", "data", "server.db")
	db, err := OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("OpenDatabase() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	if _, err := os.Stat(filepath.Dir(dbPath)); err != nil {
		t.Fatalf("expected parent dir to exist: %v", err)
	}
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected sqlite file to exist: %v", err)
	}
}

func TestInitializeSchemaCreatesCoreTablesAndColumns(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "server.db")
	db, err := OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("OpenDatabase() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	if err := InitializeSchema(db); err != nil {
		t.Fatalf("InitializeSchema() returned error: %v", err)
	}

	for _, tableName := range []string{"usage_reports", "model_pricing"} {
		var count int
		if err := sqlDB.QueryRow(
			`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`,
			tableName,
		).Scan(&count); err != nil {
			t.Fatalf("query table %s returned error: %v", tableName, err)
		}
		if count != 1 {
			t.Fatalf("expected table %s to exist once, got %d", tableName, count)
		}
	}

	requiredColumns := map[string][]string{
		"usage_reports": {"client_id", "request_id", "pricing_source", "inserted_at"},
		"model_pricing": {"model_id", "is_placeholder", "updated_at"},
	}
	for tableName, columns := range requiredColumns {
		existing, err := tableColumns(sqlDB, tableName)
		if err != nil {
			t.Fatalf("tableColumns(%s) returned error: %v", tableName, err)
		}
		for _, column := range columns {
			if !contains(existing, column) {
				t.Fatalf("expected column %s on %s, got %v", column, tableName, existing)
			}
		}
	}

	indexes, err := tableIndexes(sqlDB, "usage_reports")
	if err != nil {
		t.Fatalf("tableIndexes() returned error: %v", err)
	}
	if !contains(indexes, "idx_usage_reports_client_request") {
		t.Fatalf("expected unique index idx_usage_reports_client_request, got %v", indexes)
	}
}

func TestInitializeSchemaSeedsPricingIdempotently(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "server.db")
	db, err := OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("OpenDatabase() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	if err := InitializeSchema(db); err != nil {
		t.Fatalf("first InitializeSchema() returned error: %v", err)
	}
	if err := InitializeSchema(db); err != nil {
		t.Fatalf("second InitializeSchema() returned error: %v", err)
	}

	var placeholderCount int
	if err := sqlDB.QueryRow(
		`SELECT COUNT(*) FROM model_pricing WHERE is_placeholder = 1`,
	).Scan(&placeholderCount); err != nil {
		t.Fatalf("count placeholder pricing returned error: %v", err)
	}
	if placeholderCount != 1 {
		t.Fatalf("expected one placeholder pricing, got %d", placeholderCount)
	}

	var knownModelCount int
	if err := sqlDB.QueryRow(
		`SELECT COUNT(*) FROM model_pricing WHERE model_id LIKE 'claude-%'`,
	).Scan(&knownModelCount); err != nil {
		t.Fatalf("count seeded claude models returned error: %v", err)
	}
	if knownModelCount == 0 {
		t.Fatal("expected seeded claude models")
	}
}

func TestInitializeSchemaEnforcesClientRequestUniqueness(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "server.db")
	db, err := OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("OpenDatabase() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	if err := InitializeSchema(db); err != nil {
		t.Fatalf("InitializeSchema() returned error: %v", err)
	}

	report := entity.UsageReport{
		ClientID:             "client-1",
		RequestID:            "session:msg-1",
		AppType:              "claude",
		Model:                "claude-opus-4-1",
		PricingSource:        "exact",
		CreatedAtUnix:        1743840000,
		DataSource:           "session_log",
		InputCostUSD:         "0",
		OutputCostUSD:        "0",
		CacheReadCostUSD:     "0",
		CacheCreationCostUSD: "0",
		TotalCostUSD:         "0",
	}
	if err := db.Create(&report).Error; err != nil {
		t.Fatalf("first Create() returned error: %v", err)
	}

	duplicate := report
	duplicate.ID = 0
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("expected duplicate client/request insert to fail")
	}
}

func tableColumns(db *sql.DB, tableName string) ([]string, error) {
	rows, err := db.Query(`PRAGMA table_info(` + tableName + `)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var (
			cid        int
			name       string
			columnType string
			notNull    int
			defaultVal sql.NullString
			primaryKey int
		)
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultVal, &primaryKey); err != nil {
			return nil, err
		}
		columns = append(columns, name)
	}
	return columns, rows.Err()
}

func tableIndexes(db *sql.DB, tableName string) ([]string, error) {
	rows, err := db.Query(`PRAGMA index_list(` + tableName + `)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []string
	for rows.Next() {
		var (
			seq     int
			name    string
			unique  int
			origin  string
			partial int
		)
		if err := rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			return nil, err
		}
		indexes = append(indexes, name)
	}
	return indexes, rows.Err()
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if strings.EqualFold(value, expected) {
			return true
		}
	}
	return false
}
