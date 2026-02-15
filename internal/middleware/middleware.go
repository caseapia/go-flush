package middleware

import (
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/ranks"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

func LoadRank(rankSrv *ranks.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		val := c.Locals("user")
		user, ok := val.(*models.User)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		rankID := user.StaffRank

		rank, err := rankSrv.SearchRankByID(c, rankID)
		if err != nil {
			return err
		}

		c.Locals("rank", rank)

		return c.Next()
	}
}

func RequireRankFlag(flags ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		val := c.Locals("rank")
		rank, ok := val.(*models.RankStructure)
		if !ok {
			return &fiber.Error{Code: 401, Message: "unauthorized"}
		}

		for _, flag := range flags {
			if rank.HasFlag(flag) {
				return c.Next()
			}
		}

		slog.WithData(slog.M{
			"required_flags": flags,
			"rank":           rank,
		}).Errorf("action stopped because it must have flags: %v", flags)

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":          "forbidden",
			"required_flags": flags,
		})
	}
}

func UpdateLastLogin(repo *mysql.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		user := c.Locals("user")
		if user != nil {
			if u, ok := user.(*models.User); ok && u != nil {
				if updateErr := repo.UpdateLastLogin(c, u.ID); updateErr != nil {
					slog.Warn("Failed to update last_login", "userID", u.ID, "error", updateErr)
				}
			}
		}

		return err
	}
}
