package notifications

import (
	"context"
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

func (s *Service) SendNotification(ctx context.Context, userID uint64, notifyType models.NotificationsType, title, text string, senderID *uint64) {
	s.repo.SendNotification(
		ctx,
		models.Notification{
			Title:     title,
			UserID:    userID,
			SenderID:  senderID,
			Text:      text,
			CreatedAt: time.Now(),
		},
	)
}

func (s *Service) PopulateNotifications(ctx context.Context, userID uint64) ([]models.Notification, error) {
	notifications, err := s.repo.PopulateNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}

	return notifications, err
}

func (s *Service) ReadNotifications(ctx context.Context, userID uint64) []models.Notification {
	notifications := s.repo.ReadNotifications(ctx, userID)

	return notifications
}
