package config

import (
	loggerHandler "github.com/caseapia/goproject-flush/internal/handler/logger"
	userHandler "github.com/caseapia/goproject-flush/internal/handler/user"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userH *userHandler.UserHandler, loggerH *loggerHandler.LoggerHandler) {
	api := app.Group("/api")
	userHandler.RegisterRoutes(api, userH)
	loggerHandler.RegisterRoutes(api, loggerH)
}
