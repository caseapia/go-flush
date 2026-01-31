package config

import (
	database "github.com/caseapia/goproject-flush/internal/db"
	loggerrepo "github.com/caseapia/goproject-flush/internal/repository/logger"
	userrepo "github.com/caseapia/goproject-flush/internal/repository/user"
)

func Connect() error {
	return database.Connect()
}

func NewUserRepository() *userrepo.UserRepository {
	return userrepo.NewUserRepository(database.DB)
}

func NewLoggerRepository() *loggerrepo.LoggerRepository {
	return loggerrepo.NewLoggerRepository(database.DB)
}
