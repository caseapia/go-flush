package auth

import (
	"strings"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/auth"
	"github.com/caseapia/goproject-flush/internal/service/invite"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type Handler struct {
	authService   *auth.Service
	inviteService *invite.Service
}

func NewHandler(auth *auth.Service, invite *invite.Service) *Handler {
	return &Handler{authService: auth, inviteService: invite}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var body struct {
		Login      string `json:"login"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		InviteCode string `json:"inviteCode"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	invite, err := h.inviteService.GetInviteByID(c.UserContext(), body.InviteCode)
	if err != nil || invite.Used {
		return fiber.NewError(fiber.StatusBadRequest, "invite code is invalid or already used")
	}

	user, err := h.authService.Register(
		c.Context(),
		body.Login,
		body.InviteCode,
		body.Email,
		body.Password,
		c.IP(),
	)

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				if strings.Contains(mysqlErr.Message, "users.name") {
					return fiber.NewError(fiber.StatusConflict, "login already exists")
				}
				if strings.Contains(mysqlErr.Message, "users.email") {
					return fiber.NewError(fiber.StatusConflict, "email already exists")
				}
				return fiber.NewError(fiber.StatusConflict, "duplicate entry")
			}
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = h.inviteService.UseInvite(c.UserContext(), body.InviteCode, user.ID)
	if err != nil {
		slog.Error("Failed to mark invite as used", "error", err, "code", body.InviteCode)
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	slog.WithData(slog.M{
		"login": user.Name,
		"id":    user.ID,
	}).Debug("User successfully registered")

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var body struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	access, refresh, err := h.authService.Login(c.Context(), body.Login, body.Password, c.Get("User-Agent"), c.IP())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	slog.WithData(slog.M{
		"user": body.Login,
	}).Debug("User successfully logged in")

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Strict",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	return c.JSON(fiber.Map{
		"accessToken":  access,
		"refreshToken": refresh,
	})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	val := c.Locals("user")
	user, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	sessionIDVal := c.Locals("session_id")

	sessionID, ok := sessionIDVal.(string)
	if !ok {
		return &fiber.Error{Code: 401, Message: "invalid session"}
	}

	status := h.authService.Logout(c.Context(), sessionID)

	slog.WithData(slog.M{
		"user": user.ID,
	}).Debug("User logouted successfully")

	return c.JSON(status)
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "refresh token is missing")
	}

	newAccess, newRefresh, err := h.authService.Refresh(
		c.Context(),
		refreshToken,
		c.Get("User-Agent"),
		c.IP(),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefresh,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	return c.JSON(fiber.Map{
		"accessToken":  newAccess,
		"refreshToken": newRefresh,
	})
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/auth")

	group.Post("/refresh", h.Refresh)
	group.Post("/register", h.Register)
	group.Post("/login", h.Login)
}

func (h *Handler) RegisterPrivateRoute(router fiber.Router) {
	group := router.Group("/auth")

	group.Delete("/logout", h.Logout)
}
