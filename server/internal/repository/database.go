package repository

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// OpenDatabase 打开 SQLite 数据库，并确保父目录存在。
func OpenDatabase(path string) (*gorm.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	return gorm.Open(sqlite.Open(path), &gorm.Config{})
}
