package app

import (
	"github.com/caseapia/goproject-flush/config"
	database "github.com/caseapia/goproject-flush/internal/db"
	"github.com/caseapia/goproject-flush/internal/handler/auth"
	"github.com/caseapia/goproject-flush/internal/handler/invite"
	"github.com/caseapia/goproject-flush/internal/handler/logger"
	"github.com/caseapia/goproject-flush/internal/handler/ranks"
	"github.com/caseapia/goproject-flush/internal/handler/user"
	"github.com/caseapia/goproject-flush/internal/middleware"
	mysqlRepo "github.com/caseapia/goproject-flush/internal/repository/mysql"
	authService "github.com/caseapia/goproject-flush/internal/service/auth"
	inviteService "github.com/caseapia/goproject-flush/internal/service/invite"
	loggerService "github.com/caseapia/goproject-flush/internal/service/logger"
	ranksService "github.com/caseapia/goproject-flush/internal/service/ranks"
	userService "github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gookit/slog"
)

func NewApp() (*fiber.App, error) {
	config.LoadEnv()
	setupLogger()

	dbs, err := database.NewDatabases()
	if err != nil {
		return nil, err
	}
	mainRepo := mysqlRepo.NewRepository(dbs.Main)
	logsRepo := mysqlRepo.NewRepository(dbs.Logs)

	loggerSrv := loggerService.NewService(*logsRepo)
	ranksSrv := ranksService.NewService(mainRepo, loggerSrv)
	userSrv := userService.NewService(mainRepo, loggerSrv)
	inviteSrv := inviteService.NewService(mainRepo, *loggerSrv)
	authSrv := authService.NewService(*mainRepo)

	authHandler := auth.NewHandler(authSrv, inviteSrv)
	userHandler := user.NewUserHandler(userSrv, ranksSrv)
	inviteHandler := invite.NewHandler(inviteSrv)
	loggerHandler := logger.NewHandler(loggerSrv)
	ranksHandler := ranks.NewHandler(ranksSrv)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,https://fe-go-flush.vercel.app,http://localhost:8080",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Cache-Control",
	}))

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	api := app.Group("/api")

	public := api.Group("/public")
	authHandler.RegisterRoutes(public)

	private := api.Group("/private")
	private.Use(auth.AuthMiddleware(authSrv))
	private.Use(middleware.UpdateLastLogin(mainRepo))
	private.Use(middleware.LoadRank(ranksSrv))

	authHandler.RegisterPrivateRoute(private)

	userHandler.RegisterRoutes(private)
	inviteHandler.RegisterRoutes(private)
	loggerHandler.RegisterRoutes(private)
	ranksHandler.RegisterRoutes(private)

	return app, nil
}

func setupLogger() {
	slog.Configure(func(l *slog.SugaredLogger) {
		if f, ok := l.Formatter.(*slog.TextFormatter); ok {
			f.EnableColor = true
		}
	})
	slog.SetFormatter(slog.NewJSONFormatter())
}
