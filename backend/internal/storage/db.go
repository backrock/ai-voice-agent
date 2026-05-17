package storage

import (
	"os"
	"path/filepath"

	"github.com/backrock/ai-voice-agent/backend/internal/config"
	"github.com/backrock/ai-voice-agent/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	switch dbConfig.Type {
	case "sqlite":
		return initSQLite(dbConfig.SQLite.Path)
	default:
		return initSQLite("./data/app.db")
	}
}

func initSQLite(path string) (*gorm.DB, error) {
	// 创建目录
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&models.Session{},
		&models.Message{},
		&models.Provider{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
