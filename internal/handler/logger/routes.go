package logger

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router, h *Handler) {
	logs := app.Group("/logs")

	logs.Get("/:type", h.GetLogs)
}
