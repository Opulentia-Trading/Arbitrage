package models

import (
	"time"
	"github.com/uptrace/bun"
)

type Pairing struct {
	bun.BaseModel `bun:"table:pairings,alias:pairing"`

	ID        		int64     `bun:"id,pk,autoincrement"`
	Platform1_ID 	string    `bun:"platform_asset1,notnull,pk"`
	Platform2_ID 	string    `bun:"platform_asset2,notnull,pk"`
	Asset1 			string    `bun:"asset1,notnull,pk"`
	Asset2 			string    `bun:"asset2,notnull,pk"`
	CreatedAt     	time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     	time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
