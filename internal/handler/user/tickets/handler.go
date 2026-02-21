package tickets

import (
	"strconv"

	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/user/tickets"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *tickets.Service
}

func NewHandler(s *tickets.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) SearchTickets(ctx *fiber.Ctx) error {
	uVal := ctx.Locals("user")
	_, ok := uVal.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	tickets, _, err := h.service.SearchTickets(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(tickets)
}

func (h *Handler) PopulateTicket(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	uVal := c.Locals("user")
	u, ok := uVal.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	ticket, err := h.service.PopulateTicket(c.UserContext(), uint64(id), u)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ticket)
}

func (h *Handler) PopulateAllUserTickets(ctx *fiber.Ctx) error {
	uVal := ctx.Locals("user")
	u, ok := uVal.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	tickets, err := h.service.PopulateAllUserTickets(ctx.UserContext(), u.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(tickets)
}

func (h *Handler) CreateTicket(ctx *fiber.Ctx) error {
	uVal := ctx.Locals("user")
	u, ok := uVal.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var input models.TicketCreationInput
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}

	ticket, err := h.service.CreateTicket(ctx.UserContext(), *u, input.Title, input.Category, input.FirstMessage)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(ticket)
}

func (h *Handler) CreateTicketMessage(ctx *fiber.Ctx) error {
	uVal := ctx.Locals("user")
	u, ok := uVal.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var input models.TicketMessageCreationInput
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}

	message, err := h.service.CreateTicketMessage(ctx.UserContext(), &input.Ticket, u, input.Content)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(message)
}

func (h *Handler) PopulateTicketMessages(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	uVal := ctx.Locals("user")
	u, ok := uVal.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	messages, err := h.service.PopulateTicketMessages(ctx.UserContext(), uint64(id), u)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(messages)
}

func (h *Handler) TicketAssignment(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	uVal := ctx.Locals("user")
	u, ok := uVal.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	ticket, err := h.service.TicketAssignment(ctx.UserContext(), uint64(id), u)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(ticket)
}

func (h *Handler) CloseTicket(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	uVal := ctx.Locals("user")
	u, ok := uVal.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	ticket, err := h.service.CloseTicket(ctx.UserContext(), uint64(id), *u)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(ticket)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/tickets")
	groupAdmin := router.Group("/admin/tickets")

	group.Post("/create", h.CreateTicket)                                               // Create ticket
	group.Get("/populate/:id", h.PopulateTicket)                                        // Populate one ticket
	group.Get("/mytickets", h.PopulateAllUserTickets)                                   // Populate all tickets created by selected user
	group.Get("/populate/messages/:id", h.PopulateTicketMessages)                       // Populate ticket messages
	group.Post("/send", h.CreateTicketMessage)                                          // Create message in a ticket
	group.Patch("/close/:id", h.CloseTicket)                                            // Close ticket
	groupAdmin.Get("/populate", middleware.RequireFlag("STAFF"), h.SearchTickets)       // Populate all tickets existed in database
	groupAdmin.Post("/assign/:id", middleware.RequireFlag("STAFF"), h.TicketAssignment) // Assign an admin to the ticket
}
