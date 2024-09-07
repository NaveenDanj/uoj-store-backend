package db

import (
	"log"
	"peer-store/config"
	"peer-store/models"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func Setup() {
	once.Do(func() { // Ensure the instance is created only once
		db, err := gorm.Open(sqlite.Open(config.CONFIG.DatabaseName), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		// Migrate the schema if needed
		db.AutoMigrate(&models.User{})
		db.AutoMigrate(&models.AccessToken{})
		db.AutoMigrate(&models.File{})

		instance = db
	})
}

func GetDB() *gorm.DB {
	if instance == nil {
		log.Fatal("Database is not initialized. Please call db.Setup() first.")
	}
	return instance
}
