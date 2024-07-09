package uscase

import (
	dto "UserServiceAuth/storage"
	"errors"
	"time"

	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *dto.USERS) error
	GetUserByLogin(login string) (*dto.USERS, error)
	UpdateUserByID(id uint, updatedUser *dto.USERS) error
	GetUserByID(id uint) (*dto.USERS, error)
	SaveToken(token *dto.TOKENS) error
}

type UserService struct {
	userRepo     IUserRepository
	tokenService *TokenService
}

func NewUserService(userRepo IUserRepository, tokenService *TokenService) *UserService {
	return &UserService{userRepo: userRepo, tokenService: tokenService}
}

func (s *UserService) RegisterUser(user *dto.USERS) error {
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

	// Сохранение токена в БД
	expTime := time.Now().Add(time.Hour * 1).Unix()
	tokenRecord := &dto.TOKENS{
		USERID:       user.USERID,
		ACCESSTOCKEN: token,
		EXP:          expTime,
		TIMECREATE:   time.Now().Unix(),
	}

	if err := s.userRepo.SaveToken(tokenRecord); err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) UpdateUserByID(id uint, updatedUser *dto.USERS) error {
	existingID, err := s.userRepo.GetUserByID(id)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingID == nil {
		return errors.New("user with this id not exists")
	}

	return s.userRepo.UpdateUserByID(id, updatedUser)
}
