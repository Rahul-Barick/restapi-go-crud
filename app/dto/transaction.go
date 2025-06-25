package dto

import "github.com/shopspring/decimal"

type CreateTransactionRequest struct {
	SourceAccountID      int             `json:"source_account_id" binding:"required"`
	DestinationAccountID int             `json:"destination_account_id" binding:"required"`
	Amount               decimal.Decimal `json:"amount" binding:"required"`
}
