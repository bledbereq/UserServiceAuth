package storage

import (
	"fmt"
	"log"
	"time"

	"UserServiceAuth/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(cfg *config.Config) *gorm.DB {
	const maxAttempts = 10
	const delay = 5 * time.Second

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName)

	var err error
	for attempts := 0; attempts < maxAttempts; attempts++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err == nil {
			sqlDB, err := db.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					break
				}
			}
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", attempts+1, maxAttempts, err)
		time.Sleep(delay)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxAttempts, err)
	}

	err = db.AutoMigrate(&TOKENS{}, &USERS{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	return db
}
