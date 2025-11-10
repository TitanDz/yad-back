package database

import (
	"fmt"
	"log"

	"github.com/ethandiaz/yad-back/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Try MySQL first
	db, err := initMySQL(cfg)
	if err == nil {
		log.Println("✓ Connected to MySQL database")
		return db, nil
	}

	// Log MySQL error
	log.Printf("⚠ MySQL connection failed: %v\n", err)
	log.Println("ℹ Falling back to SQLite for local development...")

	// Fall back to SQLite
	db, err = initSQLite()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database (MySQL and SQLite both failed): %w", err)
	}

	log.Println("✓ Connected to SQLite database (dev mode)")
	return db, nil
}

func initMySQL(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initSQLite() (*gorm.DB, error) {
	// Use SQLite for local development
	dbFile := "yad.db"
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	return db, nil
}
