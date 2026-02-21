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

func (s *Service) GetCommonLogs(ctx context.Context, startDate, endDate, keywords string) ([]models.CommonLog, int, error) {
	return s.repo.GetCommonLogs(ctx, startDate, endDate, keywords)
}

func (s *Service) GetPunishmentLogs(ctx context.Context, startDate, endDate, keywords string) ([]models.PunishmentLog, int, error) {
	return s.repo.GetPunishmentLogs(ctx, startDate, endDate, keywords)
}

func (s *Service) GetTicketsLog(ctx context.Context, startDate, endDate, keywords string) ([]models.TicketsLog, int, error) {
	return s.repo.GetTicketsLog(ctx, startDate, endDate, keywords)
}

func (s *Service) Log(
	ctx context.Context,
	loggerType models.LoggerType,
	adminID *uint64,
	userID *uint64,
	action interface{},
	additional ...string,
) {
	var addInfo *string
	if len(additional) > 0 {
		addInfo = &additional[0]
	}

	base := models.BaseLog{
		AdditionalInfo: addInfo,
		Date:           time.Now(),
	}

	act, ok := action.(models.Action)
	if !ok {
		slog.Error("invalid action type")
		return
	}
	base.Action = act

	switch loggerType {

	case models.PunishmentLogger:
		s.repo.SavePunishmentLog(ctx, &models.PunishmentLog{
			BaseLog: base,
			AdminID: *adminID,
			UserID:  userID,
		})

	case models.CommonLogger:
		s.repo.SaveCommonLog(ctx, &models.CommonLog{
			BaseLog: base,
			AdminID: *adminID,
			UserID:  userID,
		})

	case models.TicketLogger:
		s.repo.SaveTicketsLog(ctx, &models.TicketsLog{
			BaseLog: base,
			AdminID: *adminID,
			UserID:  userID,
		})

	default:
		slog.Error("unknown logger type")
	}
}
