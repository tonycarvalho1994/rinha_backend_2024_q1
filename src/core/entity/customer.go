package entity

import (
	"errors"
	"fmt"
)

type Customer struct {
	ID           string        `json:"id"`
	Limit        float64       `json:"limite"`
	Transactions []Transaction `json:"transactions"`
}

func (c *Customer) CalculateBalance() float64 {
	balance := 0.0

	for _, transaction := range c.Transactions {
		if transaction.Type == Credit {
			balance += transaction.Value
		} else if transaction.Type == Debit {
			balance -= transaction.Value
		}
	}

	return balance
}

func (c *Customer) AddTransaction(transaction Transaction) error {
	err := c.validateTransaction(transaction)
	if err != nil {
		return err
	}
	c.Transactions = append(c.Transactions, transaction)
	return nil
}

func (c *Customer) validateTransaction(transaction Transaction) error {
	balance := c.CalculateBalance()
	if transaction.Type == Credit {
		balance += transaction.Value
	} else if transaction.Type == Debit {
		balance -= transaction.Value
	} else {
		return errors.New("invalid transaction type. available types: 'c' for credit, 'd' for debit")
	}

	_min := c.Limit * -1
	_max := c.Limit

	if balance > _max || balance < _min {
		msg := fmt.Sprintf("invalid to proceed transaction. limit exceeded. available: %v", c.Limit)
		return errors.New(msg)
	}

	return nil
}
