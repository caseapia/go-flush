package models

import (
	"time"

	"github.com/uptrace/bun"
)

type BanModel struct {
	bun.BaseModel  `bun:"table:bans"`
	ID             uint64    `bun:"id,pk,autoincrement" json:"id"`
	IssuedBy       uint64    `bun:"issued_by,notnull" json:"issuedBy"`
	IssuedTo       uint64    `bun:"issued_to,notnull" json:"issuedTo"`
	Date           time.Time `bun:"date,notnull,default:current_timestamp" json:"date"`
	ExpirationDate time.Time `bun:"expiration_date,notnull" json:"expirationDate"`
	Reason         string    `bun:"reason,notnull" json:"reason"`
}

type BanModelDTO struct {
	BanModel   `bun:",extend"`
	AdminName  string `bun:"admin_name" json:"adminName"`
	TargetName string `bun:"target_name" json:"targetName"`
}
