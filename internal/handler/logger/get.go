package loggerhandler

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func (l *LoggerHandler) GetLogs(c *fiber.Ctx) error {
	logs, err := l.service.GetLogs(c.UserContext())

	if err != nil {
		log.Println("Error getting logs:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch logs"})
	}

	return c.JSON(logs)
}
