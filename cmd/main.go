package main

import (
	"log"

	"github.com/caseapia/goproject-flush/config"
	userHandler "github.com/caseapia/goproject-flush/internal/handler/user"
	loggerService "github.com/caseapia/goproject-flush/internal/service/logger"
	userService "github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Connect()

	app := fiber.New()

	// репозитории
	userRepo := config.NewUserRepository()
	loggerRepo := config.NewLoggerRepository()

	// сервисы
	loggerSrv := loggerService.NewLoggerService(loggerRepo)    // сервис логов
	userSrv := userService.NewUserService(userRepo, loggerSrv) // сервис пользователей с логированием

	userHandler := userHandler.NewUserHandler(userSrv)

	config.SetupRoutes(app, userHandler)

	log.Fatal(app.Listen(":8080"))
}
