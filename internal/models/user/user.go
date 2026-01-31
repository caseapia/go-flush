package models

import "time"

type User struct {
	ID        uint64     `bun:"id,pk,autoincrement" json:"id"`
	Name      string     `bun:"name,unique,notnull" json:"name"`
	IsBanned  bool       `bun:"is_banned" json:"isBanned,omitempty"`
	BanReason *string    `bun:"ban_reason" json:"banReason,omitempty"`
	IsDeleted bool       `bun:"is_deleted" json:"isDeleted,omitempty"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time  `bun:"updated_at,notnull,default:current_timestamp" json:"updatedAt"`
	DeletedAt *time.Time `bun:"deleted_at,nullzero" json:"-"`
}

func (User) TableName() string {
	return "users"
}
