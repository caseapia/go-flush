package loggerrepository

import (
	"context"
	"time"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
)

func (l *LoggerRepository) Log(
	ctx context.Context,
	entry *loggermodule.ActionLog,
) error {
	entry.CreatedAt = time.Now()

	_, err := l.db.NewInsert().
		Model(entry).
		Exec(ctx)

	return err
}
