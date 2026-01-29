package loggermodule

import "time"

type LoggerAction string

const (
	Ban        LoggerAction = "has banned"
	Unban      LoggerAction = "has unbanned"
	Create     LoggerAction = "has created"
	SoftDelete LoggerAction = "has soft-deleted"
)

type ActionLog struct {
	ID             uint `gorm:"primaryKey;autoIncrement"`
	AdminID        uint
	UserID         uint
	Action         LoggerAction
	AdditionalInfo *string
	CreatedAt      time.Time
}
