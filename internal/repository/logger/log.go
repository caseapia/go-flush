package logger

import (
	"context"
	"errors"
	"time"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	"github.com/gookit/slog"
)

func (l *LoggerRepository) Log(
	ctx context.Context,
	loggerType loggermodule.LoggerType,
	entry interface{},
) error {
	switch e := entry.(type) {
	case *loggermodule.CommonLog:
		e.Date = time.Now()
	case *loggermodule.PunishmentLog:
		e.Date = time.Now()
	default:
		return errors.New("unsupported log entry type")
	}

	_, err := l.db.NewInsert().
		Model(entry).
		Exec(ctx)
	if err != nil {
		slog.Error("failed to insert action log:", err)
		return err
	}

	slog.WithData(slog.M{
		"loggerType": loggerType,
		"entryData":  entry,
	}).Debugf("log inserted successfully")

	return nil
}
