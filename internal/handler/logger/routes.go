package loggerhandler

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router, h *LoggerHandler) {
	logs := app.Group("/logs")

	logs.Get("/", h.GetLogs)
}
