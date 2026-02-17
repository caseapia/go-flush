package ranks

import (
	"strings"

	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/ranks"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type Handler struct {
	service *ranks.Service
}

func NewHandler(s *ranks.Service) *Handler {
	return &Handler{service: s}
}

func (r *Handler) GetRanksList(c *fiber.Ctx) error {
	ranks, err := r.service.SearchAllRanks(c)
	if err != nil {
		slog.WithData(slog.M{
			"error": err.Error(),
		}).Debug("Error fetching ranks")

		return &fiber.Error{Code: 500, Message: err.Error()}
	}

	return c.JSON(ranks)
}

func (r *Handler) CreateRank(c *fiber.Ctx) error {
	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var input models.CreateRankBody

	if err := c.BodyParser(&input); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	for i, flag := range input.Flags {
		input.Flags[i] = strings.ToUpper(flag)
	}

	rank, err := r.service.CreateRank(c, sender.ID, input.Name, input.Color, input.Flags)

	if err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	return c.Status(fiber.StatusCreated).JSON(rank)
}

func (r *Handler) DeleteRank(c *fiber.Ctx) error {
	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	rankID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	IsSuccess, err := r.service.DeleteRank(c, sender.ID, rankID)
	if err != nil {
		return err
	}

	return c.JSON(IsSuccess)
}

func (h *Handler) EditRank(c *fiber.Ctx) error {
	rankID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid rank id")
	}

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var input models.RankStructure
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	input.ID = rankID

	r, err := h.service.EditRank(c.UserContext(), sender.ID, &input)
	if err != nil {
		return err
	}

	return c.JSON(r)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/admin/ranks")

	group.Get("/", h.GetRanksList)
	group.Post("/create", middleware.RequireFlag("STAFFMANAGEMENT"), h.CreateRank)
	group.Delete("/delete/:id", middleware.RequireFlag("MANAGER"), h.DeleteRank)
	group.Patch("/edit/:id", middleware.RequireFlag("MANAGER"), h.EditRank)
}
