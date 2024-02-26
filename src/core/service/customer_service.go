package service

import (
	"errors"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/entity"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/port"
	"time"
)

type CustomerService struct {
	Repository port.CustomerRepositoryInterface
}

type AddTransactionOutput struct {
	Limit   float64 `json:"limite"`
	Balance float64 `json:"saldo"`
}

type BalanceDTO struct {
	Total       float64 `json:"total"`
	DateHistory string  `json:"data_extrato"`
	Limit       float64 `json:"limite"`
}

type TransactionHistory struct {
	Balance          BalanceDTO           `json:"saldo"`
	LastTransactions []entity.Transaction `json:"ultimas_transacoes"`
}

func (c *CustomerService) AddTransaction(customerId string, transaction entity.Transaction) (*AddTransactionOutput, error) {
	customer, err := c.Repository.FindById(customerId)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	err = customer.AddTransaction(transaction)
	if err != nil {
		return nil, err
	}

	err = c.Repository.Update(customer)
	if err != nil {
		return nil, err
	}
	err = c.Repository.AddTransaction(customer.ID, &transaction)
	if err != nil {
		return nil, err
	}

	output := AddTransactionOutput{
		Limit:   customer.Limit,
		Balance: customer.CalculateBalance(),
	}

	return &output, nil
}

func (c *CustomerService) sortByDate(transactions []entity.Transaction) []entity.Transaction {
	//sort.Slice(transactions, func(i, j int) bool {
	//	timeI, _ := time.Parse(time.RFC3339, transactions[i].CarriedOut)
	//	timeJ, _ := time.Parse(time.RFC3339, transactions[j].CarriedOut)
	//	return timeI.After(timeJ)
	//})

	if len(transactions) > 10 {
		return transactions[:10]
	}
	return transactions
}

func (c *CustomerService) GetTransactionHistory(customerId string) (*TransactionHistory, error) {
	customer, err := c.Repository.FindById(customerId)
	if err != nil {
		return nil, errors.New("customer not found")
	}
	balance := BalanceDTO{
		Total:       customer.CalculateBalance(),
		DateHistory: time.Now().Format(time.RFC3339),
		Limit:       customer.Limit,
	}
	lastTransactions := c.sortByDate(customer.Transactions)
	output := TransactionHistory{
		Balance:          balance,
		LastTransactions: lastTransactions,
	}

	return &output, nil
}
