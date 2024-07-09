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
func (r *UserRepository) GetUserByID(id uint) (*models.USERS, error) {
	var user models.USERS
	if err := r.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserByLogin(login string, updatedUser *models.USERS) error {
	return r.db.Model(&models.USERS{}).Where("login = ?", login).Updates(updatedUser).Error
}

func (r *UserRepository) SaveToken(token *models.TOKENS) error {
	return r.db.Create(token).Error
}
