package adminRanks

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/caseapia/goproject-flush/internal/models/admin/ranks"
	"github.com/caseapia/goproject-flush/internal/models/logger"
	adminError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/admin"
	userError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"
	"github.com/gofiber/fiber/v2"
)

func (s *RanksService) CreateRank(ctx *fiber.Ctx, adminID uint64, rankName string, rankColor string, rankFlags []string) (*ranks.RankStructure, error) {
	existing, err := s.ranksRepo.GetByName(ctx.UserContext(), rankName)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if existing != nil {
		return nil, adminError.RankAlreadyExists()
	}

	if rankName == "" || len(rankName) < 3 || len(rankName) > 30 {
		return nil, userError.UserInvalidUsername()
	}

	rank := &ranks.RankStructure{Name: rankName, Color: rankColor, Flags: rankFlags}

	if err := s.ranksRepo.Create(ctx.UserContext(), rank); err != nil {
		return nil, err
	}

	addInfo := "with name: " + rankName + ", with color: " + rankColor + "with flags: " + strings.Join(rankFlags, ", ")

	_ = s.logger.Log(ctx.UserContext(), "common", adminID, nil, logger.CreateRank, addInfo)

	return rank, nil
}
