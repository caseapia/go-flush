package logger

import (
	"github.com/caseapia/goproject-flush/internal/handler/logger"
	loggerhandler "github.com/caseapia/goproject-flush/internal/handler/logger"
	loggerrepo "github.com/caseapia/goproject-flush/internal/repository/logger"
	userrepo "github.com/caseapia/goproject-flush/internal/repository/user"
	loggerservice "github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type LoggerModule struct {
	Handler *logger.Handler
	Service *loggerservice.LoggerService
}

func NewLoggerModule(
	logsDB *bun.DB,
	userRepo *userrepo.UserRepository,
) *LoggerModule {

	lrepo := loggerrepo.NewLoggerRepository(logsDB)
	srv := loggerservice.NewLoggerService(lrepo, userRepo)
	h := logger.NewHandler(srv)

	return &LoggerModule{
		Handler: h,
		Service: srv,
	}
}

func (m *LoggerModule) RegisterRoutes(app fiber.Router) {
	loggerhandler.RegisterRoutes(app, m.Handler)
}
