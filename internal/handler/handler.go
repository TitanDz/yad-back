package handler

import (
"gorm.io/gorm"
)

type Handlers struct {
	db    *gorm.DB
	User  *UserHandler
}

func NewHandlers(db *gorm.DB) *Handlers {
	return &Handlers{
		db:   db,
		User: NewUserHandler(db),
	}
}
