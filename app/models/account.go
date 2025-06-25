package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID        uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID int             `gorm:"uniqueIndex;not null" json:"account_id"`
	Balance   decimal.Decimal `gorm:"type:numeric(20,6)" json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
