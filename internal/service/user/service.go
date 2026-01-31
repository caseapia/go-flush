package service

import (
	"errors"

	repository "github.com/caseapia/goproject-flush/internal/repository/user"
	loggerservice "github.com/caseapia/goproject-flush/internal/service/logger"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserBanned        = errors.New("user banned")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotBanned     = errors.New("user is not banned")
	ErrInvalidUserName   = errors.New("invalid user name")
)

type UserService struct {
	repo   *repository.UserRepository
	logger *loggerservice.LoggerService
}

func NewUserService(r *repository.UserRepository, l *loggerservice.LoggerService) *UserService {
	return &UserService{
		repo:   r,
		logger: l,
	}
}
