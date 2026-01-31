package loggerservice

import (
	"context"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
)

func (l *LoggerService) Log(
	ctx context.Context,
	adminID uint64,
	userID uint64,
	action loggermodule.LoggerAction,
	additional ...string,
) error {
	var addInfo *string
	if len(additional) > 0 {
		addInfo = &additional[0]
	}

	logEntry := loggermodule.ActionLog{
		AdminID:        adminID,
		UserID:         userID,
		Action:         action,
		AdditionalInfo: addInfo,
	}

	return l.repo.Log(ctx, &logEntry)
}
