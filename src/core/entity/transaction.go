package entity

import (
	"errors"
	"time"
)

type TransactionType string

const (
	Credit TransactionType = "c"
	Debit  TransactionType = "d"
)

type Transaction struct {
	Value       float64         `json:"valor"`
	Type        TransactionType `json:"tipo"`
	Description string          `json:"descricao"`
	CarriedOut  string          `json:"realizada_em"`
}

func NewTransaction(value float64, transactionType TransactionType, description string) (*Transaction, error) {
	if len(description) > 10 {
		return nil, errors.New("description must have a maximum of 10 characters")
	}

	return &Transaction{
		Value:       value,
		Type:        transactionType,
		Description: description,
		CarriedOut:  time.Now().Format(time.RFC3339),
	}, nil
}
