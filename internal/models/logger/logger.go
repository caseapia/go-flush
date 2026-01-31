package loggermodule

import "time"

type LoggerAction string

const (
	Ban              LoggerAction = "has banned"
	Unban            LoggerAction = "has unbanned"
	Create           LoggerAction = "has created"
	SoftDelete       LoggerAction = "has soft-deleted"
	RestoreUser      LoggerAction = "has restored"
	SearchByUsername LoggerAction = "searched by username"
	SearchByUserID   LoggerAction = "searched by user ID"
	SearchByAllUsers LoggerAction = "searched all users"
	SearchLogs       LoggerAction = "searched logs"
)

type ActionLog struct {
	ID             uint64       `bun:"id,pk,autoincrement" json:"id"`
	AdminID        uint64       `bun:"admin_id,notnull" json:"adminId"`
	UserID         uint64       `bun:"user_id,notnull" json:"userId"`
	Action         LoggerAction `bun:"action,notnull" json:"action"`
	AdditionalInfo *string      `bun:"additional_info,nullzero" json:"additionalInfo,omitempty"`
	CreatedAt      time.Time    `bun:"created_at" json:"createdAt"`
}

func (ActionLog) TableName() string {
	return "action_logs"
}
