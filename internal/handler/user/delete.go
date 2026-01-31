package handler

import (
	"strconv"

	service "github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	IsDeleted, err := h.service.DeleteUser(c.UserContext(), uint64(id))

	if err != nil {
		status := fiber.StatusNotFound

		if err == service.ErrUserBanned {
			status = fiber.StatusForbidden
		}

		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(IsDeleted)
}

func (h *UserHandler) RestoreUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	IsRestored, err := h.service.RestoreUser(c.UserContext(), uint64(id))

	if err != nil {
		status := fiber.StatusNotFound

		if err == service.ErrUserBanned {
			status = fiber.StatusForbidden
		}

		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(IsRestored)
}
