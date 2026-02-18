package models

import "time"

type NotificationsType = string

const (
	Information = "information"
	Error       = "error"
	Success     = "success"
)

type Notification struct {
	ID        uint64            `bun:"id,pk,autoincrement,unique" json:"id"`
	CreatedAt time.Time         `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	Type      NotificationsType `bun:"type,notnull,default:information" json:"type"`
	Title     string            `bun:"title" json:"title"`
	SenderID  *uint64           `bun:"sender_id" json:"senderId"`
	UserID    uint64            `bun:"user_id" json:"userId"`
	Text      string            `bun:"text" json:"text"`
	IsReaded  bool              `bun:"is_readed" json:"isReaded"`

	Sender *User `bun:"rel:belongs-to,join:sender_id=id" json:"sender,omitempty"`
	User   *User `bun:"rel:belongs-to,join:user_id=id" json:"user"`
}

type NotificationsInput struct {
	ID *uint64 `json:"id"`
}

type SendNotificationInput struct {
	Type     NotificationsType `json:"type"`
	Title    string            `json:"title"`
	SenderID *uint64           `json:"senderId"`
	UserID   uint64            `json:"userId"`
	Text     string            `json:"text"`
}
