package notifications

import (
	"strconv"

	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/user/notifications"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *notifications.Service
}

func NewHandler(s *notifications.Service) *Handler {
	return &Handler{service: s}
}

func (s *Handler) PopulateNotifications(c *fiber.Ctx) error {
	uVal := c.Locals("user")
	u, ok := uVal.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	notifications, err := s.service.PopulateNotifications(c.UserContext(), u.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(notifications)
}

func (s *Handler) SendNotification(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	uVal := c.Locals("user")
	sender, ok := uVal.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var input models.SendNotificationInput

	s.service.SendNotification(c.UserContext(), uint64(id), input.Type, input.Title, input.Text, &sender.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (s *Handler) ReadNotifications(c *fiber.Ctx) error {
	uVal := c.Locals("user")
	u, ok := uVal.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	notifications := s.service.ReadNotifications(c.Context(), u.ID)

	return c.Status(fiber.StatusOK).JSON(notifications)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/notifications")
	groupAdmin := router.Group("/admin/notifications")

	group.Get("/populate", h.PopulateNotifications)
	group.Post("/read", h.ReadNotifications)
	groupAdmin.Post("/send/:id", middleware.RequireFlag("ADMIN"), h.SendNotification)
}
