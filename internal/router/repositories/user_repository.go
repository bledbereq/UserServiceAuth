package repositories

import (
	models "UserServiceAuth/storage"

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

func (r *UserRepository) CreateUser(user *models.USERS) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserByLogin(login string) (*models.USERS, error) {
	var user models.USERS
	if err := r.db.Where("login = ?", login).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
