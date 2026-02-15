package user

import (
	"strconv"
	"time"

	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/ranks"
	"github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type Handler struct {
	service *user.Service
	rank    *ranks.Service
}

func NewUserHandler(s *user.Service, r *ranks.Service) *Handler {
	return &Handler{service: s, rank: r}
}

func (h *Handler) SearchAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetUsersList(c.UserContext())
	if err != nil {
		slog.WithData(slog.M{
			"e": err.Error(),
		}).Debug("Error fetching users")

		return &fiber.Error{Code: 500, Message: err.Error()}
	}

	return c.JSON(users)
}

func (h *Handler) GetOwnAccount(c *fiber.Ctx) error {
	val := c.Locals("user")
	u, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	user, err := h.service.GetOwnAccount(c.UserContext(), u.ID)
	if err != nil {
		slog.WithData(slog.M{
			"e": err,
		}).Debug("Error get user account")

		return &fiber.Error{Code: 500, Message: err.Error()}
	}

	return c.JSON(user)
}

// ! Admin actions
func (h *Handler) SearchUserByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	rank, err := h.rank.SearchRankByID(c, sender.StaffRank)

	if !rank.HasFlag("ADMIN") && sender.ID != uint64(id) {
		return &fiber.Error{Code: 401, Message: "no access"}
	}

	u, err := h.service.SearchUser(c.UserContext(), sender.ID, uint64(id))
	if err != nil {
		return err
	}

	return c.JSON(u)
}

func (h *Handler) GetUserPrivate(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user ID")
	}

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	user, err := h.service.SearchUser(c.Context(), sender.ID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return c.JSON(user.GetPrivateData())
}

func (h *Handler) BanUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var Body struct {
		UnbanDate time.Time `json:"unbanDate"`
		Reason    string    `json:"reason"`
	}

	c.BodyParser(&Body)

	ban, err := h.service.BanUser(c.UserContext(), admin.ID, uint64(id), Body.UnbanDate, Body.Reason)
	if err != nil {
		return err
	}

	return c.JSON(ban)
}

func (h *Handler) UnbanUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	unban, err := h.service.UnbanUser(c.UserContext(), admin.ID, uint64(id))
	if err != nil {
		return err
	}

	return c.JSON(unban)
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var Body struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&Body); err != nil {
		return &fiber.Error{Code: 400, Message: "invalid request"}
	}

	newUser, err := h.service.CreateUser(c, admin.ID, Body.Name)
	if err != nil {
		return err
	}

	return c.JSON(newUser)
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	id, _ := strconv.Atoi(c.Params("id"))

	deleted, err := h.service.DeleteUser(c.UserContext(), admin.ID, uint64(id))
	if err != nil {
		return err
	}

	return c.JSON(deleted)
}

func (h *Handler) RestoreUser(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	id, _ := strconv.Atoi(c.Params("id"))

	restored, err := h.service.RestoreUser(c.UserContext(), admin.ID, uint64(id))
	if err != nil {
		return err
	}

	return c.JSON(restored)
}

func (h *Handler) SetStaffRank(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	userID, err := c.ParamsInt("id")
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	var input struct {
		Status int `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	u, err := h.service.SetStaffRank(
		c.Context(),
		admin.ID,
		uint64(userID),
		input.Status,
	)
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return &fiber.Error{Code: 500, Message: err.Error()}
	}

	return c.Status(fiber.StatusOK).JSON(u)
}

func (h *Handler) SetDeveloperRank(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	userID, err := c.ParamsInt("id")
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	var input struct {
		Status int `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	u, err := h.service.SetDeveloperRank(
		c.Context(),
		admin.ID,
		uint64(userID),
		input.Status,
	)
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(u)
}

func (h *Handler) ChangeUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var input struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if input.Name != nil {
		if len(*input.Name) <= 1 {
			return fiber.NewError(fiber.StatusBadRequest, "new nickname is too short")
		}
	}

	u, err := h.service.ChangeUser(c.UserContext(), sender.ID, uint64(userID), input.Name, input.Email)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(u)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/user")
	groupAdmin := router.Group("/admin/user")

	group.Get("/all", h.SearchAllUsers)
	group.Get("/account", h.GetOwnAccount)
	group.Get("/:id", h.SearchUserByID)

	groupAdmin.Put("/create", middleware.RequireRankFlag("SENIOR"), h.CreateUser)
	groupAdmin.Patch("/ban/:id", middleware.RequireRankFlag("ADMIN"), h.BanUser)
	groupAdmin.Delete("/unban/:id", middleware.RequireRankFlag("ADMIN"), h.UnbanUser)
	groupAdmin.Delete("/delete/:id", middleware.RequireRankFlag("SENIOR"), h.DeleteUser)
	groupAdmin.Put("/restore/:id", middleware.RequireRankFlag("MANAGER"), h.RestoreUser)
	groupAdmin.Patch("/rank/staff/:id", middleware.RequireRankFlag("STAFFMANAGEMENT"), h.SetStaffRank)
	groupAdmin.Patch("/rank/developer/:id", middleware.RequireRankFlag("STAFFMANAGEMENT"), h.SetDeveloperRank)
	groupAdmin.Get("/:id", middleware.RequireRankFlag("ADMIN"), h.GetUserPrivate)
	groupAdmin.Patch("/edit/:id", middleware.RequireRankFlag("SENIOR"), h.ChangeUser)
}
