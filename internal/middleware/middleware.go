package middleware

import (
	"fmt"

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

		adminRankID := user.StaffRank
		developerRankID := user.DeveloperRank

		adminRank, aRankErr := rankSrv.SearchRankByID(c, adminRankID)
		if aRankErr != nil {
			return aRankErr
		}
		developerRank, devRankErr := rankSrv.SearchRankByID(c, developerRankID)
		if devRankErr != nil {
			return devRankErr
		}

		c.Locals("rank", []*models.RankStructure{adminRank, developerRank})

		return c.Next()
	}
}

func RequireFlag(flags ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		val := c.Locals("rank")
		userVal := c.Locals("user")

		ranks, ok := val.([]*models.RankStructure)
		if !ok || len(ranks) == 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		user, ok := userVal.(*models.User)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		userFlags := user.Flags

		for _, requiredFlag := range flags {
			for _, rank := range ranks {
				if rank.HasFlag(requiredFlag) {
					return c.Next()
				}
			}
			for _, userFlag := range *userFlags {
				if userFlag == requiredFlag || userFlag == "MANAGER" {
					return c.Next()
				}
			}
		}

		slog.WithData(slog.M{
			"required_flags": flags,
			"rank":           ranks,
			"user":           user,
		}).Errorf("action stopped because it must have flags: %v", flags)

		return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("forbidden. required flags: %s", flags))
	}
}

func UpdateLastLogin(repo *mysql.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		user := c.Locals("user")
		if user != nil {
			if u, ok := user.(*models.User); ok && u != nil {
				if updateErr := repo.UpdateLastLogin(c, u.ID); updateErr != nil {
					slog.WithData(slog.M{
						"userID": u.ID,
						"error":  updateErr,
					}).Warn("Failed to update last_login")
				}
			}
		}

		return err
	}
}
