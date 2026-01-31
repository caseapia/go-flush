package loggerhandler

import service "github.com/caseapia/goproject-flush/internal/service/logger"

type LoggerHandler struct {
	service *service.LoggerService
}

func NewLoggerHandler(s *service.LoggerService) *LoggerHandler {
	return &LoggerHandler{service: s}
}
