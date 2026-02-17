package models

import "github.com/uptrace/bun"

type RankStructure struct {
	bun.BaseModel `bun:"table:ranks"`
	ID            int      `bun:"column:id,pk,autoincrement" json:"id"`
	Name          string   `bun:"column:name" json:"name"`
	Color         string   `bun:"column:color" json:"color"`
	Flags         []string `bun:"column:flags" json:"flags"`

	Users      []User `bun:"rel:has-many,join:id=staff_rank" json:"users,omitempty"`
	Developers []User `bun:"rel:has-many,join:id=developer_rank" json:"developers,omitempty"`
}

type CreateRankBody struct {
	Name  string   `json:"name"`
	Color string   `json:"color"`
	Flags []string `json:"flags"`
}

func (r *RankStructure) HasFlag(flag string) bool {
	for _, f := range r.Flags {
		if f == flag {
			return true
		}
	}
	return false
}
