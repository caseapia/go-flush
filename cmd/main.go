package main

import (
	"log"

	"github.com/caseapia/goproject-flush/config"
	adminUserHandler "github.com/caseapia/goproject-flush/internal/handler/admin/user"
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/caseapia/goproject-flush/internal/module/admin"
	"github.com/caseapia/goproject-flush/internal/module/logger"
	"github.com/caseapia/goproject-flush/internal/module/user"
	adminUserRepoPkg "github.com/caseapia/goproject-flush/internal/repository/admin/user"
	adminUserService "github.com/caseapia/goproject-flush/internal/service/admin/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gookit/slog"
)

func main() {
	// ---------------- Logger ----------------
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})
	slog.SetFormatter(slog.NewJSONFormatter())

	// ---------------- DB ----------------
	config.LoadEnv()
	db := config.Connect()
	if db == nil {
		log.Fatal("Failed to connect to DB")
	}

	// ---------------- Fiber App ----------------
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.AppErrorHandler,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,https://fe-go-flush.vercel.app",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// ---------------- Modules ----------------
	loggerM := logger.NewLoggerModule(db)

	adminM := admin.NewAdminModule(db, nil, nil, loggerM.Service)

	userRepo := config.NewUserRepository()
	adminUserRepo := adminUserRepoPkg.NewAdminUserRepository(db)

	adminUserSrv := adminUserService.NewAdminUserService(
		userRepo,
		adminM.RanksService,
		loggerM.Service,
		adminUserRepo,
	)

	userHandler := adminUserHandler.NewAdminUserHandler(adminUserSrv)

	adminM.UserHandler = userHandler

	userM := user.NewUserModule(db, loggerM.Service, adminM.RanksService)
	userM.Service.SetRanksService(adminM.RanksService)

	// ---------------- Routes ----------------
	config.SetupRoutes(app, userM, loggerM, adminM)

	// ---------------- Start Server ----------------
	log.Fatal(app.Listen(":8080"))
}
