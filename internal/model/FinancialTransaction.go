package model

import (
	"errors"
	"time"
)

// FinancialTransaction represents the 'FinancialTransaction' table in your database.
type FinancialTransaction struct {
	TransactionID string    `json:"transactionID"`
	UserID        string    `json:"userID"`
	UnitID        string    `json:"unitID"`
	PaymentMethod string    `json:"paymentMethod"`
	Amount        int       `json:"amount"`
	CreateTime    time.Time `json:"createTime"`
}

func (f *FinancialTransaction) Validate() error {
	if f.TransactionID == "" {
		return errors.New("TransactionID is required")
	}
	if f.UserID == "" {
		return errors.New("UserID is required")
	}
	if f.UnitID == "" {
		return errors.New("UnitID is required")
	}
	if f.PaymentMethod == "" {
		return errors.New("PaymentMethod is required")
	}
	if f.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}
