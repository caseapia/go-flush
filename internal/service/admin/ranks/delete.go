package adminRanks

import (
	"strconv"

	"github.com/caseapia/goproject-flush/internal/models/logger"
	adminError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/admin"
	"github.com/gofiber/fiber/v2"
)

func (s *RanksService) DeleteRank(ctx *fiber.Ctx, id int) (bool, error) {
	r, err := s.GetByID(ctx.UserContext(), id)
	if err != nil {
		return false, err
	}
	if r == nil {
		return false, adminError.RankNotExists()
	}

	addInfo := "with ID: " + strconv.FormatInt(r.ID, 10) + ", with name: " + r.Name

	_ = s.logger.Log(ctx.UserContext(), logger.CommonLogger, 0, nil, logger.DeleteRank, addInfo)

	if err := s.ranksRepo.DeleteRank(ctx.UserContext(), r); err != nil {
		return false, err
	}

	return true, nil
}
