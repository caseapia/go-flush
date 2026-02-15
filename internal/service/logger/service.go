package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
)

type Service struct {
	repo mysql.Repository
}

func NewService(r mysql.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetCommonLogs(ctx context.Context) ([]models.LogDTO, error) {
	return s.repo.GetCommonLogs(ctx)
}

func (s *Service) GetPunishmentLogs(ctx context.Context) ([]models.LogDTO, error) {
	return s.repo.GetPunishmentLogs(ctx)
}

func (s *Service) Log(
	ctx context.Context,
	loggerType models.LoggerType,
	adminID uint64,
	userID *uint64,
	action interface{},
	additional ...string,
) error {
	var addInfo *string
	if len(additional) > 0 {
		addInfo = &additional[0]
	}

	base := models.BaseLog{
		AdminID:        adminID,
		UserID:         userID,
		AdditionalInfo: addInfo,
		Date:           time.Now(),
	}

	switch loggerType {
	case models.PunishmentLogger:
		act, ok := action.(models.Action)
		if !ok {
			return fmt.Errorf("expected models.Action for PunishmentLogger, got %T", action)
		}
		base.Action = act

		return s.repo.SavePunishmentLog(ctx, &models.PunishmentLog{
			BaseLog: base,
		})

	case models.CommonLogger:
		act, ok := action.(models.Action)
		if !ok {
			return fmt.Errorf("expected models.Action for CommonLogger, got %T", action)
		}
		base.Action = act
		return s.repo.SaveCommonLog(ctx, &models.CommonLog{
			BaseLog: base,
		})

	default:
		return fmt.Errorf("unknown logger type: %s", loggerType)
	}
}
