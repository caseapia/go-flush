package logger

import (
	"context"

	models "github.com/caseapia/goproject-flush/internal/models/logger"
	"github.com/gofiber/fiber/v2"
)

func (s *LoggerService) GetLogs(ctx context.Context, logType string) ([]models.BaseLog, error) {
	switch logType {
	case "all":
		logs, err := s.repo.GetAllLogs(ctx)
		return logs, err
	case "common":
		commonLogs, err := s.repo.GetCommonLogs(ctx)
		if err != nil {
			return nil, err
		}
		res := make([]models.BaseLog, len(commonLogs))
		for i, l := range commonLogs {
			res[i] = l.BaseLog
		}
		return res, nil
	case "punish":
		punishLogs, err := s.repo.GetPunishmentLogs(ctx)
		if err != nil {
			return nil, err
		}
		res := make([]models.BaseLog, len(punishLogs))
		for i, l := range punishLogs {
			res[i] = l.BaseLog
		}
		return res, nil
	default:
		return nil, fiber.NewError(404, "log type was not founded")
	}
}
