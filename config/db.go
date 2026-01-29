package config

import (
	"github.com/caseapia/goproject-flush/internal/db"
	loggerrepo "github.com/caseapia/goproject-flush/internal/repository/logger"
	userrepo "github.com/caseapia/goproject-flush/internal/repository/user"
)

func Connect() {
	db.Connect()
}

func NewUserRepository() *userrepo.UserRepository {
	return userrepo.NewUserRepository(db.DB)
}

func NewLoggerRepository() *loggerrepo.LoggerRepository {
	return loggerrepo.NewLoggerRepository(db.DB)
}
