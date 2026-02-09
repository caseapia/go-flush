package adminUser

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) EditFlags(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var Input struct {
		Flags []string `json:"flags"`
	}

	c.BodyParser(&Input)

	u, err := h.service.EditFlags(c.UserContext(), uint64(id), Input.Flags)

	if err != nil {
		return err
	}

	return c.JSON(u)
}
