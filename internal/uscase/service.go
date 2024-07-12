package service

import (
	models "UserServiceAuth/storage"
	"errors"

	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *models.USERS) error
	GetUserByLogin(login string) (*models.USERS, error)
	UpdateUserByID(id uint, updatedUser *models.USERS) error
	GetUserByID(id uint) (*models.USERS, error)
}

type UserService struct {
	userRepo IUserRepository
}

func NewUserService(userRepo IUserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

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

func (s *UserService) UpdateUserByID(id uint, updatedUser *models.USERS) error {
	existingID, err := s.userRepo.GetUserByID(id)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingID == nil {
		return errors.New("user with this id not exists")
	}

	return s.userRepo.UpdateUserByID(id, updatedUser)
}
