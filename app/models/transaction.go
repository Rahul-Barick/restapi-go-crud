package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Transaction table structs
type Transaction struct {
	ID                   uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	SourceAccountID      int             `gorm:"not null" json:"source_account_id"`
	DestinationAccountID int             `gorm:"not null" json:"destination_account_id"`
	Amount               decimal.Decimal `gorm:"type:numeric(20,6);not null" json:"amount"`
	ReferenceID          uuid.UUID       `gorm:"type:uuid;uniqueIndex;not null" json:"reference_id"`
	CreatedAt            time.Time       `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	SourceAccount      Account `gorm:"foreignKey:SourceAccountID;references:AccountID"`
	DestinationAccount Account `gorm:"foreignKey:DestinationAccountID;references:AccountID"`
}

type LedgerType string

const (
	CREDIT LedgerType = "CREDIT"
	DEBIT  LedgerType = "DEBIT"
)

// Ledger Entry table stuct
type LedgerEntries struct {
	ID            int             `gorm:"primaryKey;autoIncrement"`
	AccountID     int             `gorm:"not null;index" json:"account_id"`
	TransactionID uuid.UUID       `gorm:"type:uuid;not null;index" json:"transaction_id"`
	Amount        decimal.Decimal `gorm:"type:numeric(20,6);not null" json:"amount"`
	Type          LedgerType      `gorm:"type:ledger_entry_type;not null" json:"type"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`

	// Foreign keys
	Account     Account     `gorm:"foreignKey:AccountID;references:AccountID"`
	Transaction Transaction `gorm:"foreignKey:TransactionID;references:ID"`
}
