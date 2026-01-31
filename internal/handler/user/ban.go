package handler

import (
	"strconv"

	service "github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) BanUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var input struct {
		Reason string `json:"reason"`
	}

	c.BodyParser(&input)

	user, err := h.service.BanUser(c.UserContext(), uint64(0), uint64(id), input.Reason)

	if err != nil {
		status := fiber.StatusNotFound

		if err == service.ErrUserBanned {
			status = fiber.StatusForbidden
		}

		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (h *UserHandler) UnbanUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	user, err := h.service.UnbanUser(c.UserContext(), uint64(0), uint64(id))

	if err != nil {
		status := fiber.StatusNotFound
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}
