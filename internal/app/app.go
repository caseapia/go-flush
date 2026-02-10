package app

import (
	"github.com/caseapia/goproject-flush/config"
	database "github.com/caseapia/goproject-flush/internal/db"
	adminUserHandler "github.com/caseapia/goproject-flush/internal/handler/admin/user"
	"github.com/caseapia/goproject-flush/internal/middleware"
	adminModule "github.com/caseapia/goproject-flush/internal/module/admin"
	loggerModule "github.com/caseapia/goproject-flush/internal/module/logger"
	userModule "github.com/caseapia/goproject-flush/internal/module/user"
	adminRanksRepo "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"
	adminUserRepo "github.com/caseapia/goproject-flush/internal/repository/admin/user"
	loggerRepo "github.com/caseapia/goproject-flush/internal/repository/logger"
	userRepo "github.com/caseapia/goproject-flush/internal/repository/user"
	adminRanksService "github.com/caseapia/goproject-flush/internal/service/admin/ranks"
	adminUserService "github.com/caseapia/goproject-flush/internal/service/admin/user"
	loggerService "github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gookit/slog"
)

func NewApp() (*fiber.App, error) {
	config.LoadEnv()

	dbs, err := database.NewDatabases()
	if err != nil {
		return nil, err
	}

	slog.Configure(func(l *slog.SugaredLogger) {
		f := l.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})
	slog.SetFormatter(slog.NewJSONFormatter())

	userRepo := userRepo.NewUserRepository(dbs.Main)
	adminUserRepo := adminUserRepo.NewAdminUserRepository(dbs.Main)
	ranksRepo := adminRanksRepo.NewRanksRepository(dbs.Main)
	loggerRepo := loggerRepo.NewLoggerRepository(dbs.Logs)

	loggerSrv := loggerService.NewLoggerService(loggerRepo, userRepo)

	userRankSetter := adminRanksService.NewUserRankSetter(
		userRepo,
		ranksRepo,
		loggerSrv,
	)

	ranksSrv := adminRanksService.NewRanksService(
		ranksRepo,
		userRankSetter,
		loggerSrv,
	)

	adminUserSrv := adminUserService.NewAdminUserService(
		userRepo,
		ranksSrv,
		loggerSrv,
		adminUserRepo,
	)

	adminUserHandler := adminUserHandler.NewAdminUserHandler(adminUserSrv)

	userM := userModule.NewUserModule(dbs.Main, loggerSrv, ranksSrv)
	adminM := adminModule.NewAdminModule(dbs.Main, ranksSrv, adminUserHandler, loggerSrv)
	loggerM := loggerModule.NewLoggerModule(dbs.Logs, userRepo)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.AppErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,https://fe-go-flush.vercel.app",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	api := app.Group("/api")

	userM.RegisterRoutes(api)
	adminM.RegisterRoutes(api)
	loggerM.RegisterRoutes(api)

	return app, nil
}
