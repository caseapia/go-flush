package auth

import (
	"strings"

	"github.com/caseapia/goproject-flush/internal/service/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

func AuthMiddleware(authSrv *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing Authorization header FROM AUTH MIDDLEWARE")
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid Authorization header format")
		}

		token := parts[1]

		user, claims, err := authSrv.ParseJWT(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		if user == nil || claims == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token data")
		}

		if user.ActiveBanID != nil {
			slog.WithData(slog.M{
				"user":        user,
				"activeBanID": user.ActiveBanID,
			}).Error("user action stopped due to active ban")
			return fiber.NewError(fiber.StatusForbidden, "you're not allowed to perform this action due to active ban")
		}

		c.Locals("user", user)
		c.Locals("session_id", claims.SessionID)
		return c.Next()
	}
}
