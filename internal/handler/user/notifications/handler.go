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

	notifications, err := s.service.PopulateNotifications(c.UserContext(), u.ID, u.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(notifications)
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

func (s *Handler) ClearNotifications(c *fiber.Ctx) error {
	uVal := c.Locals("user")
	u, ok := uVal.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	notifications, _ := s.service.ClearNotifications(c.Context(), u.ID)

	return c.Status(fiber.StatusOK).JSON(notifications)
}

func (s *Handler) RemoveOwnNotification(c *fiber.Ctx) error {
	uVal := c.Locals("user")
	sender, ok := uVal.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var input models.RemoveNotificationsInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	isDeleted, err := s.service.RemoveNotification(c.UserContext(), sender.ID, sender.ID, input.NotifyID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(isDeleted)
}

// ! Admin actions
func (s *Handler) SendNotification(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	uVal := c.Locals("user")
	sender, ok := uVal.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var input models.SendNotificationInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	s.service.SendNotification(c.UserContext(), uint64(id), input.Type, input.Title, input.Text, &sender.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (s *Handler) PopulateUserNotifications(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	uVal := c.Locals("user")
	sender, ok := uVal.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	notifications, err := s.service.PopulateNotifications(c.UserContext(), uint64(id), sender.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(notifications)
}

func (s *Handler) RemoveNotification(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	uVal := c.Locals("user")
	sender, ok := uVal.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var input models.RemoveNotificationsInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	isDeleted, err := s.service.RemoveNotification(c.UserContext(), uint64(id), sender.ID, input.NotifyID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(isDeleted)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/notifications")
	groupAdmin := router.Group("/admin/notifications")

	group.Get("/populate", h.PopulateNotifications)
	group.Post("/read", h.ReadNotifications)
	group.Delete("/clear", h.ClearNotifications)
	group.Delete("/remove", h.RemoveOwnNotification)
	groupAdmin.Post("/send/:id", middleware.RequireFlag("ADMIN"), h.SendNotification)
	groupAdmin.Get("/populate/:id", middleware.RequireFlag("ADMIN"), h.PopulateUserNotifications)
	groupAdmin.Delete("/remove/:id", middleware.RequireFlag("SENIOR"), h.RemoveNotification)
}
