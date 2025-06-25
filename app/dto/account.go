package dto

import "github.com/shopspring/decimal"

type CreateAccountRequest struct {
	AccountID      int             `json:"account_id" binding:"required,gt=0"`
	InitialBalance decimal.Decimal `json:"initial_balance" binding:"required,nonnegdecimal"`
}

type AccountResponse struct {
	AccountID int             `json:"account_id"`
	Balance   decimal.Decimal `json:"balance"`
}
