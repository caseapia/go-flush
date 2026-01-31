package loggerservice

import (
	"context"

	models "github.com/caseapia/goproject-flush/internal/models/logger"
)

func (s *LoggerService) GetLogs(ctx context.Context) ([]models.ActionLog, error) {
	logs, err := s.repo.GetLogs(ctx)

	if err != nil {
		return nil, err
	}

	return logs, nil
}
