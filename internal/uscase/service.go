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
	UpdateUserByLogin(login string, updatedUser *dto.USERS) error
	DeleteUserByLogin(login string) error
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

	token, err := s.tokenService.GenerateToken(user.USERNAME, user.EMAIL, user.LOGIN, user.ISADMIN)
	if err != nil {
		return "", err
	}

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

func (s *UserService) UpdateUserByLogin(login, token string, updatedUser *dto.USERS) error {
	claims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		return err
	}
	tokenIsAdmin, isAdminOk := claims["isadmin"].(bool)
	tokenLogin, loginOk := claims["login"].(string)

	if !loginOk {
		return errors.New("invalid token: login not found")
	}

	if !isAdminOk || !tokenIsAdmin {
		if tokenLogin != login {
			return errors.New("token is not authorized for this user")
		}
	}

	existingUser, err := s.userRepo.GetUserByLogin(login)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser == nil {
		return errors.New("user with this login does not exist")
	}

	return s.userRepo.UpdateUserByLogin(login, updatedUser)
}

func (s *UserService) DeleteUserByLogin(login, token string) error {
	claims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		return err
	}
	tokenIsAdmin, isAdminOk := claims["isadmin"].(bool)

	if !isAdminOk || !tokenIsAdmin {
		return errors.New("token is not have admin's root")
	}

	existingUser, err := s.userRepo.GetUserByLogin(login)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser == nil {
		return errors.New("user with this login does not exist")
	}

	return s.userRepo.DeleteUserByLogin(login)
}
