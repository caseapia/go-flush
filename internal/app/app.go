package app

import (
	"time"

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
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
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

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"code":    code,
				"message": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,https://fe-go-flush.vercel.app,http://localhost:8080",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Cache-Control",
	}))

	app.Get("/api/ping", func(c *fiber.Ctx) error {
		v, err := mem.VirtualMemory()
		if err != nil {
			return err
		}

		cpuPercent, err := cpu.Percent(time.Millisecond*100, false)
		if err != nil {
			return err
		}

		uptime, err := host.Uptime()
		if err != nil {
			return err
		}

		var cpuUsage float64
		if len(cpuPercent) > 0 {
			cpuUsage = cpuPercent[0]
		}

		return c.JSON(fiber.Map{
			"status":    "pong",
			"timestamp": time.Now().Unix(),
			"system": fiber.Map{
				"cpu":    cpuUsage,                    // cpu loading
				"ram":    v.UsedPercent,               // ram loading
				"ram_gb": v.Used / 1024 / 1024 / 1024, // gb usage
				"uptime": uptime,                      // server uptime
			},
		})
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
	f := slog.NewJSONFormatter()
	f.PrettyPrint = true
	f.TimeFormat = "02/01/2006 15:04:05.000"
	slog.SetFormatter(f)
}
