package entity

import (
	"errors"
	"strings"
	"time"
)

type TransactionType string

const (
	Credit TransactionType = "c"
	Debit  TransactionType = "d"
)

type Transaction struct {
	Value       int             `json:"valor"`
	Type        TransactionType `json:"tipo"`
	Description string          `json:"descricao"`
	CarriedOut  string          `json:"realizada_em"`
}

func validateString(transactionType TransactionType) error {
	switch transactionType {
	case "c", "d":
		return nil
	default:
		return errors.New("invalid transaction type")
	}
}

func NewTransaction(value int, transactionType TransactionType, description string) (*Transaction, error) {
	if description == "" {
		return nil, errors.New("invalid transaction, description length must be max 10")
	}
	if len(description) > 10 {
		return nil, errors.New("invalid transaction, description length must be max 10")
	}
	err := validateString(transactionType)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Value:       value,
		Type:        transactionType,
		Description: strings.TrimSpace(description),
		CarriedOut:  time.Now().Format(time.RFC3339),
	}, nil
}
