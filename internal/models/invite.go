package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Invite struct {
	bun.BaseModel `bun:"table:invites,alias:i"`

	ID        uint64    `bun:"id,pk,autoincrement" json:"id"`
	Code      string    `bun:"code" json:"code"`
	CreatedBy uint64    `bun:"created_by" json:"createdBy"`
	Used      bool      `bun:"used" json:"used"`
	UsedBy    *uint64   `bun:"used_by" json:"usedBy"`
	CreatedAt time.Time `bun:"created_at" json:"createdAt"`

	Creator *User `bun:"rel:belongs-to,join:created_by=id" json:"creator"`
	User    *User `bun:"rel:belongs-to,join:used_by=id" json:"user"`
}
