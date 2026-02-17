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
	SystemLogger     = "system"
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
	ResetUserSensetiveData Action = "has reset user mail and IPs"
	EditRank               Action = "has edited rank"
)

type BaseLog struct {
	ID             uint64    `bun:"id,pk,autoincrement" json:"id"`
	Date           time.Time `bun:"date,notnull" json:"date"`
	AdminID        uint64    `bun:"admin_id,notnull" json:"adminId"`
	Action         Action    `bun:"action,notnull" json:"action"`
	UserID         *uint64   `bun:"user_id" json:"userId"`
	AdditionalInfo *string   `bun:"additional_information" json:"additionalInfo,omitempty"`

	Admin *User `bun:"rel:belongs-to,join:admin_id=id" json:"admin"`
	User  *User `bun:"rel:belongs-to,join:user_id=id" json:"user"`
}

type CommonLog struct {
	bun.BaseModel `bun:"table:admin_common"`
	BaseLog
	Limit int `bun:"-" json:"limit"`
}
type PunishmentLog struct {
	bun.BaseModel `bun:"table:admin_punishments"`
	BaseLog
	Limit int `bun:"-" json:"limit"`
}

type SystemLog struct {
	bun.BaseModel `bun:"table:system"`
	ID            uint64    `bun:"id,pk,autoincrement" json:"id"`
	Date          time.Time `bun:"date,notnull" json:"date"`
	Event         string    `bun:"event,notnull" json:"event"`
}

type LogPopulate struct {
	StartDate string     `json:"dateStart"`
	EndDate   string     `json:"dateEnd"`
	Type      LoggerType `json:"type"`
	Keywords  *string    `json:"keywords"`
}
