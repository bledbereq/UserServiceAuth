package repositories

import (
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	router := &UserRepository{
		db: db,
	}
	return router
}
