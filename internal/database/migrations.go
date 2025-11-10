package database

import (
	"github.com/ethandiaz/yad-back/internal/domain"
	"gorm.io/gorm"
)

// RunMigrations automatically runs all database migrations
func RunMigrations(db *gorm.DB) error {
	// Auto migrate all domain models
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.AuthToken{},
		&domain.Minyan{},
		&domain.Location{},
		&domain.UserSettings{},
		&domain.Notification{},
		&domain.Place{},
	); err != nil {
		return err
	}

	return nil
}
