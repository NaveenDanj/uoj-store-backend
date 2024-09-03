package db

import (
	"fmt"
	"peer-store/config"
	"peer-store/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDatabase() (*gorm.DB, error) {
	// Connect to SQLite database (or create it if it doesn't exist)
	db, err := gorm.Open(sqlite.Open(config.CONFIG.DatabaseName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the models
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
