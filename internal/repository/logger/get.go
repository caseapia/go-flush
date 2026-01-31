package loggerrepository

import (
	"context"

	models "github.com/caseapia/goproject-flush/internal/models/logger"
)

func (l *LoggerRepository) GetLogs(ctx context.Context) ([]models.ActionLog, error) {
	var logs []models.ActionLog

	err := l.db.NewSelect().
		Model(&logs).
		Scan(ctx)

	return logs, err
}
