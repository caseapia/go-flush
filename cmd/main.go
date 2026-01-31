package main

import (
	"log"

	"github.com/caseapia/goproject-flush/config"
	loggerhandler "github.com/caseapia/goproject-flush/internal/handler/logger"
	userHandler "github.com/caseapia/goproject-flush/internal/handler/user"
	loggerService "github.com/caseapia/goproject-flush/internal/service/logger"
	userService "github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.LoadEnv()
	config.Connect()

	if err := config.Connect(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	app := fiber.New()

	// репозитории
	userRepo := config.NewUserRepository()
	loggerRepo := config.NewLoggerRepository()

	// сервисы
	loggerSrv := loggerService.NewLoggerService(loggerRepo)    // сервис логов
	userSrv := userService.NewUserService(userRepo, loggerSrv) // сервис пользователей с логированием

	userHandler := userHandler.NewUserHandler(userSrv)
	loggerHandler := loggerhandler.NewLoggerHandler(loggerSrv)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	config.SetupRoutes(app, userHandler, loggerHandler)

	log.Fatal(app.Listen(":8080"))
}
