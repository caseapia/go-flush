package logger

import (
	"time"

	"github.com/uptrace/bun"
)

type LoggerType string

const (
	CommonLogger     = "common"
	PunishmentLogger = "punish"
)

type BaseLog struct {
	ID             uint64    `bun:"id,pk,autoincrement" json:"id"`
	Date           time.Time `bun:"date,notnull" json:"date"`
	AdminName      string    `bun:"admin_name,notnull" json:"adminName"`
	AdminID        uint64    `bun:"admin_id,notnull" json:"adminId"`
	UserName       *string   `bun:"user_name" json:"userName"`
	UserID         *uint64   `bun:"user_id" json:"userId"`
	AdditionalInfo *string   `bun:"additional_information" json:"additionalInfo,omitempty"`
}

type CommonLog struct {
	bun.BaseModel `bun:"table:admin_common"`
	BaseLog
	Action CommonAction `bun:"action,notnull" json:"action"`
}

type PunishmentLog struct {
	bun.BaseModel `bun:"table:admin_punishments"`
	BaseLog
	Action UserPunishment `bun:"action,notnull" json:"action"`
}
