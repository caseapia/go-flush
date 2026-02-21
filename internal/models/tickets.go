package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Status string
type Priority string

const (
	Open     Status = "Waiting for user"
	Closed   Status = "Closed"
	Pending  Status = "Waiting for staff"
	Resolved Status = "Resolved"
)

const (
	Low      Priority = "Low"
	Medium   Priority = "Medium"
	High     Priority = "High"
	Critical Priority = "Critical"
)

// Ticket
type Ticket struct {
	bun.BaseModel `bun:"table:tickets"`

	ID        uint64    `bun:"id,pk,autoincrement" json:"id"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"createdAt"`
	Status    Status    `bun:"status,notnull" json:"status"`
	Priority  Priority  `bun:"priority" json:"prioriy"`
	AuthorID  uint64    `bun:"author_id,notnull" json:"-"`
	HandledBy *uint64   `bun:"handling_by" json:"-"`
	Title     string    `bun:"title,notnull" json:"title"`
	Category  string    `bun:"category,notnull" json:"category"`
	UpdatedAt time.Time `bun:"updated_at" json:"updatedAt"`

	Author  *TicketAuthor  `bun:"rel:belongs-to,join:author_id=id" json:"author"`
	Handler *TicketHandler `bun:"rel:belongs-to,join:handling_by=id" json:"handler"`
}

type TicketCreationInput struct {
	Title        string `bun:"title,notnull" json:"title"`
	Category     string `bun:"category,notnull" json:"category"`
	FirstMessage string `json:"message"`
}

// Ticket message
type TicketMessage struct {
	bun.BaseModel `bun:"table:tickets_messages"`

	ID        uint64    `bun:"id,pk,autoincrement" json:"id"`
	TicketID  uint64    `bun:"ticket_id,notnull" json:"ticketID"`
	AuthorID  uint64    `bun:"author_id,notnull" json:"-"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"createdAt"`
	Content   string    `bun:"content,notnull" json:"content"`

	Author *TicketAuthor `bun:"rel:belongs-to,join:author_id=id" json:"author"`
}

type TicketMessageCreationInput struct {
	Ticket   Ticket `json:"ticket"`
	AuthorID uint64 `bun:"author_id,notnull" json:"author_id"`
	Content  string `bun:"content,notnull" json:"content"`
}

type TicketAuthor struct {
	bun.BaseModel `bun:"table:users"`

	ID            uint64     `bun:"id,pk" json:"id"`
	Name          string     `bun:"name" json:"name"`
	LastLogin     *time.Time `bun:"last_login" json:"lastLogin"`
	StaffRank     int        `bun:"staff_rank" json:"staffRank,omitempty"`
	DeveloperRank int        `bun:"developer_rank" json:"developerRank,omitempty"`
}

type TicketHandler struct {
	bun.BaseModel `bun:"table:users"`

	ID            uint64 `bun:"id,pk" json:"id"`
	Name          string `bun:"name" json:"name"`
	StaffRank     int    `bun:"staff_rank" json:"staffRank"`
	DeveloperRank int    `bun:"developer_rank" json:"developerRank"`
}
