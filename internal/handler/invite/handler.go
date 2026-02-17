package invite

import (
	"strconv"

	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/invite"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type Handler struct {
	service *invite.Service
}

func NewHandler(s *invite.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) GetInviteCodes(c *fiber.Ctx) error {
	val := c.Locals("user")
	_, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	invites, err := h.service.GetInviteCodes(c.UserContext())
	if err != nil {
		slog.WithData(slog.M{
			"error": err,
		}).Error("error when fetch invite codes")
		return err
	}

	return c.JSON(invites)
}

func (h *Handler) CreateInvite(c *fiber.Ctx) error {
	val := c.Locals("user")
	user, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	newInvite, err := h.service.CreateInvite(c.Context(), user.ID)
	if err != nil {
		slog.WithData(slog.M{
			"error": err,
		}).Error("error when invitation code creation")
		return err
	}

	return c.JSON(newInvite)
}

func (h *Handler) DeleteInvite(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	err := h.service.DeleteInvite(c.Context(), sender.ID, uint64(id))
	if err != nil {
		slog.WithData(slog.M{
			"error": err,
		}).Error("error when delete invitation codes")
		return err
	}

	return c.JSON(true)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/admin/invite")

	group.Get("/list", middleware.RequireFlag("ADMIN"), h.GetInviteCodes)
	group.Post("/create", middleware.RequireFlag("ADMIN"), h.CreateInvite)
	group.Delete("/delete/:id", middleware.RequireFlag("LEAD"), h.DeleteInvite)
}
