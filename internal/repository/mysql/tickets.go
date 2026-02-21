package mysql

import (
	"context"
	"database/sql"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gookit/slog"
)

func (r *Repository) SearchTickets(ctx context.Context) ([]models.Ticket, int, error) {
	var tickets []models.Ticket

	query := r.db.NewSelect().
		Model(&tickets).
		Order("created_at ASC").
		Relation("Author").
		Relation("Handler").
		Limit(COLUMNS_LIMIT)

	err := query.Scan(ctx)
	return tickets, COLUMNS_LIMIT, err
}

func (r *Repository) SearchTicketByID(ctx context.Context, ticketID uint64) (*models.Ticket, *[]models.TicketMessage, error) {
	var ticket models.Ticket
	var messages []models.TicketMessage

	err := r.db.NewSelect().
		Model(&ticket).
		Where("ticket.id = ?", ticketID).
		Relation("Author").
		Relation("Handler").
		Scan(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = r.db.NewSelect().
		Model(&messages).
		Where("ticket_id = ?", ticketID).
		Relation("Author").
		Scan(ctx)
	if err != nil {
		return nil, nil, err
	}

	return &ticket, &messages, nil
}

func (r *Repository) PopulateTicket(ctx context.Context, ticketID uint64) (*models.Ticket, error) {
	t := new(models.Ticket)

	err := r.db.NewSelect().
		Model(t).
		Relation("Author").
		Relation("Handler").
		Where("ticket.id = ?", ticketID).
		Limit(1).
		Scan(ctx)

	return t, err
}

func (r *Repository) TicketAssignment(ctx context.Context, ticketID uint64, userID uint64) (*models.Ticket, error) {
	var ticket models.Ticket

	res, err := r.db.NewUpdate().
		Model((*models.Ticket)(nil)).
		Set("handling_by = ?", userID).
		Where("id = ?", ticketID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, sql.ErrNoRows
	}

	err = r.db.NewSelect().
		Model(&ticket).
		Relation("Author").
		Relation("Handler").
		Where("?TableAlias.id = ?", ticketID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *Repository) TicketUnassignment(ctx context.Context) {

}

func (r *Repository) PopulateAllUserTickets(ctx context.Context, userID uint64) ([]models.Ticket, error) {
	var tickets []models.Ticket

	err := r.db.NewSelect().
		Model(&tickets).
		Relation("Author").
		Relation("Handler").
		Where("author_id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *Repository) CreateTicket(ctx context.Context, entry models.Ticket) (*models.Ticket, error) {
	_, err := r.db.NewInsert().
		Model(&entry).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("failed to create ticket!")

		return nil, err
	}

	return &entry, nil
}

func (r *Repository) CreateTicketMessage(ctx context.Context, entry models.TicketMessage, handlerID *uint64) (*models.TicketMessage, error) {
	_, err := r.db.NewInsert().
		Model(&entry).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err, "entry": entry}).Error("failed to send message in ticket")

		return nil, err
	}

	var newStatus models.Status
	if handlerID != nil && entry.AuthorID == *handlerID {
		newStatus = models.Open
	} else {
		newStatus = models.Pending
	}

	_, err = r.db.NewUpdate().
		Model((*models.Ticket)(nil)).
		Set("status = ?", newStatus).
		Where("id = ?", entry.TicketID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (r *Repository) PopulateTicketMessages(ctx context.Context, ticketID uint64) ([]models.TicketMessage, error) {
	var messages []models.TicketMessage

	query := r.db.NewSelect().
		Model(&messages).
		Relation("Author").
		Where("ticket_id = ?", ticketID)

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func (r *Repository) CloseTicket(ctx context.Context, ticketID uint64) (*models.Ticket, error) {
	var ticket models.Ticket

	res, err := r.db.NewUpdate().
		Model((*models.Ticket)(nil)).
		Set("status = ?", models.Closed).
		Where("id = ?", ticketID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, sql.ErrNoRows
	}

	err = r.db.NewSelect().
		Model(&ticket).
		Relation("Author").
		Relation("Handler").
		Where("?TableAlias.id = ?", ticketID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}
