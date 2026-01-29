package config

import (
	userHandler "github.com/caseapia/goproject-flush/internal/handler/user"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userH *userHandler.UserHandler) {
	api := app.Group("/api")
	userHandler.RegisterRoutes(api, userH)
}
