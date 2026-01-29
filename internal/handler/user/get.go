package handler

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	user, err := h.service.GetUser(uint(id))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (h *UserHandler) GetUsersList(c *fiber.Ctx) error {
	users, err := h.service.GetUsersList()

	if err != nil {
		log.Println("Error fetching users:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch users"})
	}

	return c.JSON(users)
}
