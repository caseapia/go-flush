package handler

import service "github.com/caseapia/goproject-flush/internal/service/user"

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}
