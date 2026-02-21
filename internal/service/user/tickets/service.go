package tickets

import (
	"context"
	"fmt"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/caseapia/goproject-flush/internal/service/user/notifications"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type Service struct {
	repo   mysql.Repository
	notify notifications.Service
	logger logger.Service
}

func NewService(r mysql.Repository, n notifications.Service, l logger.Service) *Service {
	return &Service{
		repo:   r,
		notify: n,
		logger: l,
	}
}

func (s *Service) SearchTickets(ctx context.Context) ([]models.Ticket, int, error) {
	tickets, columns, err := s.repo.SearchTickets(ctx)
	if err != nil {
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return tickets, columns, err
}

func (s *Service) PopulateTicket(ctx context.Context, ticketID uint64, user *models.User) (*models.Ticket, error) {
	ticket, err := s.repo.PopulateTicket(ctx, ticketID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	rank, err := s.repo.SearchRankByID(ctx, user.StaffRank)

	isAuthor := ticket.Author.ID == user.ID
	isStaff := user.UserHasFlag("STAFF")
	rankIsStaff := rank.HasFlag("STAFF")

	if !isAuthor && !(isStaff || rankIsStaff) {
		return nil, fiber.NewError(fiber.StatusForbidden, "you have no access to this ticket")
	}

	return ticket, nil
}

func (s *Service) PopulateAllUserTickets(ctx context.Context, userID uint64) ([]models.Ticket, error) {
	tickets, err := s.repo.PopulateAllUserTickets(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return tickets, nil
}

func (s *Service) CreateTicket(ctx context.Context, user models.User, title, category, message string) (*models.Ticket, error) {
	ticket, err := s.repo.CreateTicket(ctx, models.Ticket{
		AuthorID:  user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     title,
		Category:  category,
		HandledBy: nil,
		Status:    models.Pending,
		Priority:  models.Low,
	})
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var handlerID *uint64
	if ticket.Handler != nil {
		handlerID = &ticket.Handler.ID
	}

	if len(title) > 255 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "title too long")
	}

	if ticket != nil {
		_, mErr := s.repo.CreateTicketMessage(ctx, models.TicketMessage{
			TicketID:  ticket.ID,
			AuthorID:  user.ID,
			CreatedAt: time.Now(),
			Content:   message,
		}, handlerID)
		if mErr != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, mErr.Error())
		}
	}

	return ticket, nil
}

func (s *Service) CreateTicketMessage(ctx context.Context, ticket *models.Ticket, user *models.User, content string) (*models.TicketMessage, error) {
	rank, err := s.repo.SearchRankByID(ctx, user.StaffRank)
	if err != nil {
		return nil, err
	}

	isAuthor := ticket.Author.ID == user.ID
	isStaff := user.UserHasFlag("STAFF")
	rankIsStaff := rank.HasFlag("STAFF")
	isStaffManagement := user.UserHasFlag("STAFFMANAGEMENT")
	rankIsStaffManagement := rank.HasFlag("STAFFMANAGEMENT")
	isHandler := ticket.Handler != nil && ticket.Handler.ID == user.ID

	if !isAuthor && !(isStaff || rankIsStaff) {
		slog.WithData(slog.M{
			"isAuthor":    isAuthor,
			"isStaff":     isStaff,
			"rankIsStaff": rankIsStaff,
			"ticket":      ticket,
			"user.ID":     user.ID,
		}).Error("error when sending message")
		return nil, fiber.NewError(fiber.StatusForbidden, "you have no access to this ticket")
	}
	if ticket.Status == models.Closed && !(isStaffManagement || rankIsStaffManagement) {
		return nil, fiber.NewError(fiber.StatusNotAcceptable, "the ticket is closed and you can't answer here")
	}

	if !isAuthor && !isHandler {
		slog.WithData(slog.M{
			"isAuthor": isAuthor,
			"ticket":   ticket,
			"user":     user,
		}).Error("message was not sent")

		return nil, fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("this ticket already handled by %s and you cannot answer here", ticket.Handler.Name))
	}

	var handlerID *uint64
	if ticket.Handler != nil {
		handlerID = &ticket.Handler.ID
	}

	message, err := s.repo.CreateTicketMessage(ctx, models.TicketMessage{
		TicketID:  ticket.ID,
		AuthorID:  user.ID,
		CreatedAt: time.Now(),
		Content:   content,
	}, handlerID)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *Service) PopulateTicketMessages(ctx context.Context, ticketID uint64, user *models.User) (*[]models.TicketMessage, error) {
	rank, err := s.repo.SearchRankByID(ctx, user.StaffRank)
	if err != nil {
		return nil, err
	}

	ticket, messages, err := s.repo.SearchTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	isAuthor := ticket.Author.ID == user.ID
	isStaff := user.UserHasFlag("STAFF")
	rankIsStaff := rank.HasFlag("STAFF")

	if !isAuthor && !(isStaff || rankIsStaff) {
		return nil, fiber.NewError(fiber.StatusForbidden, "you have no access to this ticket")
	}

	return messages, err
}

func (s *Service) TicketAssignment(ctx context.Context, ticketID uint64, user *models.User) (*models.Ticket, error) {
	ticket, err := s.repo.TicketAssignment(ctx, ticketID, user.ID)
	if err != nil {
		return nil, err
	}

	s.notify.SendNotification(ctx, ticket.Author.ID, models.Success, "Your ticket has been updated", fmt.Sprintf("You have a new staff member, assigned for your ticket #%v", ticket.ID), &user.ID)
	s.logger.Log(ctx, models.TicketLogger, &user.ID, &ticket.Author.ID, models.AssignedToTicket, fmt.Sprintf("ID: %v | Title: %v", ticketID, ticket.Title))

	return ticket, err
}

func (s *Service) CloseTicket(ctx context.Context, ticketID uint64, user models.User) (*models.Ticket, error) {
	rank, err := s.repo.SearchRankByID(ctx, user.StaffRank)
	if err != nil {
		return nil, err
	}

	ticket, _, err := s.repo.SearchTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	isAuthor := ticket.Author.ID == user.ID
	isStaffManagement := user.UserHasFlag("STAFFMANAGEMENT")
	rankIsStaffManagement := rank.HasFlag("STAFFMANAGEMENT")
	isHandler := ticket.Handler != nil && ticket.Handler.ID == user.ID

	if !((isAuthor || isHandler) || isStaffManagement || rankIsStaffManagement) {
		return nil, fiber.NewError(fiber.StatusForbidden, "you are not allowed to use this action")
	}
	if ticket.Status == models.Closed {
		return nil, fiber.NewError(fiber.StatusBadRequest, "ticket already closed")
	}

	ticket, err = s.repo.CloseTicket(ctx, ticket.ID)
	if err != nil {
		return nil, err
	}

	if !isAuthor {
		s.notify.SendNotification(ctx, ticket.Author.ID, models.Success, "Your ticket has been updated", fmt.Sprintf("Your ticket #%v was closed by an admin", ticket.ID), &user.ID)
		s.logger.Log(ctx, models.TicketLogger, &user.ID, &ticket.Author.ID, models.CloseTicket, fmt.Sprintf("ID: %v | Title: %v", ticketID, ticket.Title))
	}

	return ticket, nil
}
