package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	repo   mysql.Repository
	logger logger.Service
}

func NewService(r mysql.Repository, l logger.Service) *Service {
	return &Service{repo: r, logger: l}
}

func (s *Service) SendNotification(ctx context.Context, userID uint64, notifyType models.NotificationsType, title, text string, senderID *uint64) {
	s.repo.SendNotification(
		ctx,
		models.Notification{
			Title:     title,
			UserID:    userID,
			SenderID:  senderID,
			Text:      text,
			Type:      notifyType,
			CreatedAt: time.Now(),
		},
	)

	addInfo := fmt.Sprintf("Title: %s | Type: %s | Text: %s", title, notifyType, text)
	s.logger.Log(ctx, models.CommonLogger, *senderID, &userID, models.SendNotification, addInfo)
}

func (s *Service) PopulateNotifications(ctx context.Context, userID uint64, senderID uint64) ([]models.Notification, error) {
	notifications, err := s.repo.PopulateNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}

	if userID != senderID {
		s.logger.Log(ctx, models.CommonLogger, senderID, &userID, models.LookupNotifications)
	}

	return notifications, err
}

func (s *Service) ReadNotifications(ctx context.Context, userID uint64) []models.Notification {
	notifications := s.repo.ReadNotifications(ctx, userID)

	return notifications
}

func (s *Service) RemoveNotification(ctx context.Context, userID, senderID, notifyID uint64) (bool, error) {
	isDeleted, err := s.repo.RemoveNotification(ctx, userID, notifyID)
	if err != nil {
		return false, err
	}

	if userID != senderID {
		s.logger.Log(ctx, models.CommonLogger, senderID, &userID, models.DeleteNotification)
	}

	return isDeleted, nil
}

func (s *Service) ClearNotifications(ctx context.Context, userID uint64) ([]models.Notification, error) {
	notifications, err := s.repo.ClearNotifications(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return notifications, nil
}
