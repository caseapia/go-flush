package models

import (
	"time"

	"github.com/uptrace/bun"
)

type LoggerType string
type Action string

const (
	CommonLogger     = "common"
	PunishmentLogger = "punish"
	TicketLogger     = "tickets"
)

const (
	// ! Punishments
	Ban   Action = "has banned"
	Unban Action = "has unbanned"

	// ! Common actions
	CreateRank             Action = "has created rank"
	SearchByUsername       Action = "searched by username"
	SearchByUserID         Action = "searched by user ID"
	SearchLogs             Action = "searched logs"
	SetStaffRank           Action = "has set admin rank"
	SetDeveloperRank       Action = "has set developer rank"
	RestoreUser            Action = "has restored"
	Create                 Action = "has created user"
	ChangeFlags            Action = "has changed flags"
	DeleteRank             Action = "has delete rank"
	SoftDelete             Action = "has soft-deleted"
	HardDelete             Action = "has hard-deleted"
	TriedToDeleteManager   Action = "has tried to delete manager's account and action has stopped"
	CreateInvite           Action = "has created invite code"
	DeleteInvite           Action = "has deleted invite code"
	ChangeUserData         Action = "has changed user's data"
	ChangeUserPassword     Action = "has changed user's password"
	ResetUserSensetiveData Action = "has reset user IPs and last seen information"
	EditRank               Action = "has edited rank"
	LookupNotifications    Action = "lookup user notifications"
	SendNotification       Action = "send notify"
	DeleteNotification     Action = "has deleted notification"
	AssignedToTicket       Action = "has assigned to the ticket"
	CloseTicket            Action = "has closed ticket"
)

type BaseLog struct {
	ID             uint64    `bun:"id,pk,autoincrement" json:"id"`
	Date           time.Time `bun:"date,notnull" json:"date"`
	Action         Action    `bun:"action,notnull" json:"action"`
	AdditionalInfo *string   `bun:"additional_information" json:"additionalInfo,omitempty"`
}

type CommonLog struct {
	bun.BaseModel `bun:"table:admin_common"`
	BaseLog
	AdminID uint64  `bun:"admin_id,notnull" json:"-"`
	UserID  *uint64 `bun:"user_id" json:"-"`
	Limit   int     `bun:"-" json:"limit"`

	User  *User `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	Admin *User `bun:"rel:belongs-to,join:admin_id=id" json:"admin"`
}
type PunishmentLog struct {
	bun.BaseModel `bun:"table:admin_punishments"`
	BaseLog
	AdminID uint64  `bun:"admin_id,notnull" json:"-"`
	UserID  *uint64 `bun:"user_id" json:"-"`
	Limit   int     `bun:"-" json:"limit"`

	User  *User `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	Admin *User `bun:"rel:belongs-to,join:admin_id=id" json:"admin"`
}

type TicketsLog struct {
	bun.BaseModel `bun:"table:tickets_log"`
	BaseLog
	AdminID uint64  `bun:"admin_id,notnull" json:"-"`
	UserID  *uint64 `bun:"user_id" json:"-"`

	Limit int   `bun:"-" json:"limit"`
	User  *User `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	Admin *User `bun:"rel:belongs-to,join:admin_id=id" json:"admin"`
}

type LogPopulate struct {
	StartDate string     `json:"dateStart"`
	EndDate   string     `json:"dateEnd"`
	Type      LoggerType `json:"type"`
	Keywords  *string    `json:"keywords"`
}
