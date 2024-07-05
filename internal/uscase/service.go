package service

import (
	models "UserServiceAuth/storage"
	"errors"

	"gorm.io/gorm"
)

// IUserRepository определяет интерфейс для работы с хранилищем пользователей
type IUserRepository interface {
	CreateUser(user *models.USERS) error
	GetUserByLogin(login string) (*models.USERS, error)
}

// UserService представляет сервис для работы с пользователями
type UserService struct {
	userRepo IUserRepository
}

// NewUserService создает новый UserService
func NewUserService(userRepo IUserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// RegisterUser регистрирует нового пользователя
func (s *UserService) RegisterUser(user *models.USERS) error {
	existingUser, err := s.userRepo.GetUserByLogin(user.LOGIN)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("user with this login already exists")
	}
	return s.userRepo.CreateUser(user)
}

// AuthenticateUser аутентифицирует пользователя
func (s *UserService) AuthenticateUser(login, password string) (*models.USERS, error) {
	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}
	if user.PASSWORD != password {
		return nil, errors.New("invalid login or password")
	}
	return user, nil
}
