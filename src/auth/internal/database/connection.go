package database

import (
	"os"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func NewDB() *gorm.DB {
	var db *gorm.DB
	dsn := os.Getenv("DB")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database")
	}
	return db
}

