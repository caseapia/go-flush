package adminRanks

import "github.com/gofiber/fiber/v2"

func (r *Handler) CreateRank(c *fiber.Ctx) error {
	var input struct {
		Name  string   `json:"name"`
		Color string   `json:"color"`
		Flags []string `json:"flags"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	rank, err := r.service.CreateRank(c, 0, input.Name, input.Color, input.Flags)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(rank)
}
