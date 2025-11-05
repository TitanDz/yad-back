package handler

import (
	"gorm.io/gorm"
)

type Handlers struct {
	db       *gorm.DB
	User     *UserHandler
	Auth     *AuthHandler
	Minyan   *MinyanHandler
	Location *LocationHandler
	Settings *UserSettingsHandler
}

func NewHandlers(db *gorm.DB) *Handlers {
	return &Handlers{
		db:       db,
		User:     NewUserHandler(db),
		Auth:     NewAuthHandler(db),
		Minyan:   NewMinyanHandler(db),
		Location: NewLocationHandler(db),
		Settings: NewUserSettingsHandler(db),
	}
}
