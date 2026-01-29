package loggerservice

import (
	repository "github.com/caseapia/goproject-flush/internal/repository/logger"
)

type LoggerService struct {
	repo *repository.LoggerRepository
}

func NewLoggerService(r *repository.LoggerRepository) *LoggerService {
	return &LoggerService{repo: r}
}
