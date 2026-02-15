package logger

import (
	"context"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/gookit/slog"
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
) {
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
			slog.WithData(slog.M{
				"action": action,
			}).Error("expected models.Action for PunishmentLogger")
		}
		base.Action = act

		s.repo.SavePunishmentLog(ctx, &models.PunishmentLog{
			BaseLog: base,
		})

	case models.CommonLogger:
		act, ok := action.(models.Action)
		if !ok {
			slog.WithData(slog.M{
				"action": action,
			}).Error("expected models.Action for CommonLogger")
		}
		base.Action = act
		s.repo.SaveCommonLog(ctx, &models.CommonLog{
			BaseLog: base,
		})

	default:
		slog.WithData(slog.M{
			"loggerType": loggerType,
		}).Error("unknown logger type")
	}
}
