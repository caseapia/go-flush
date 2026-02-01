package userhandler

import userservice "github.com/caseapia/goproject-flush/internal/service/user"

type UserHandler struct {
	service *userservice.UserService
}

func NewUserHandler(s *userservice.UserService) *UserHandler {
	return &UserHandler{service: s}
}
