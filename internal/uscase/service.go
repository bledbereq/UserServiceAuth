package uscase

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
	userRepo     IUserRepository
	tokenService *TokenService
}

func NewUserService(userRepo IUserRepository, tokenService *TokenService) *UserService {
	return &UserService{userRepo: userRepo, tokenService: tokenService}
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

func (s *UserService) AuthenticateUser(login, password string) (string, error) {
	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil {
		return "", err
	}
	if user.PASSWORD != password {
		return "", errors.New("invalid login or password")
	}

	token, err := s.tokenService.GenerateToken(user.USERID, user.USERNAME, user.EMAIL)
	if err != nil {
		return "", err
	}

	return token, nil
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
