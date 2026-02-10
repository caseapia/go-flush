package logger

import (
	"context"
	"fmt"
	"time"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
)

func (l *LoggerService) Log(
	ctx context.Context,
	loggerType loggermodule.LoggerType,
	adminID uint64,
	userID *uint64,
	action interface{},
	additional ...string,
) error {
	var addInfo *string
	if len(additional) > 0 {
		addInfo = &additional[0]
	}

	var u *usermodel.User
	if userID != nil {
		var err error
		u, err = l.uRepo.GetByID(ctx, *userID)
		if err != nil {
			u = nil
		}
	}

	base := loggermodule.BaseLog{
		AdminID:        adminID,
		AdminName:      "",
		UserID:         userID,
		UserName:       nil,
		AdditionalInfo: addInfo,
		Date:           time.Now(),
	}

	if userID != nil && u != nil {
		base.UserName = &u.Name
	}

	switch loggerType {
	case loggermodule.PunishmentLogger:
		p, ok := action.(loggermodule.UserPunishment)
		if !ok {
			return fmt.Errorf("expected UserPunishment for PunishmentLogger, got %T", action)
		}
		return l.repo.Log(ctx, loggerType, &loggermodule.PunishmentLog{
			BaseLog: base,
			Action:  p,
		})

	case loggermodule.CommonLogger:
		c, ok := action.(loggermodule.CommonAction)
		if !ok {
			return fmt.Errorf("expected CommonAction for CommonLogger, got %T", action)
		}
		return l.repo.Log(ctx, loggerType, &loggermodule.CommonLog{
			BaseLog: base,
			Action:  c,
		})

	default:
		return fmt.Errorf("unknown loggerType: %v", loggerType)
	}
}
