package loggerservice

import (
	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
)

func (l *LoggerService) Log(adminID uint, userID uint, action loggermodule.LoggerAction, additional ...string) error {
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

	return l.repo.Log(&logEntry)
}
