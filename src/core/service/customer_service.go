package service

import (
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/entity"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/port"
	"time"
)

type CustomerService struct {
	Repository port.CustomerRepositoryInterface
}

type AddTransactionOutput struct {
	Balance int `json:"saldo"`
	Limit   int `json:"limite"`
}

type BalanceDTO struct {
	Total       int    `json:"total"`
	DateHistory string `json:"data_extrato"`
	Limit       int    `json:"limite"`
}

type TransactionHistory struct {
	Balance          BalanceDTO           `json:"saldo"`
	LastTransactions []entity.Transaction `json:"ultimas_transacoes"`
}

func (c *CustomerService) AddTransaction(customerId string, transaction entity.Transaction) (*AddTransactionOutput, error) {
	newBalance, limit, err := c.Repository.AddTransaction(customerId, &transaction)
	if err != nil {
		return nil, err
	}

	output := AddTransactionOutput{
		Balance: newBalance,
		Limit:   limit,
	}

	return &output, nil
}

func (c *CustomerService) GetTransactionHistory(customerId string) (*TransactionHistory, error) {
	limit, currentBalance, transactions, err := c.Repository.FindHistoryByCustomerId(customerId)
	if err != nil {
		return nil, err
	}
	balance := BalanceDTO{
		Total:       currentBalance,
		DateHistory: time.Now().Format(time.RFC3339),
		Limit:       limit,
	}
	output := TransactionHistory{
		Balance:          balance,
		LastTransactions: transactions,
	}

	return &output, nil
}
