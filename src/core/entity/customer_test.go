package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomer(t *testing.T) {
	transaction1 := Transaction{
		Value:       100,
		Type:        "c",
		Description: "test",
		CarriedOut:  "",
	}
	transaction2 := Transaction{
		Value:       50,
		Type:        "c",
		Description: "test",
		CarriedOut:  "",
	}
	transaction3 := Transaction{
		Value:       10,
		Type:        "d",
		Description: "test",
		CarriedOut:  "",
	}
	customer := Customer{
		ID:           "1",
		Limit:        1000,
		Transactions: []Transaction{},
	}
	err := customer.AddTransaction(transaction1)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = customer.AddTransaction(transaction2)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = customer.AddTransaction(transaction3)
	if err != nil {
		t.Errorf(err.Error())
	}

	balance := customer.CalculateBalance()
	expectedBalance := 140.0
	assert.Equal(t, balance, expectedBalance)
}

func TestCustomerShouldRaiseErrorLimitExceeded(t *testing.T) {
	transaction1 := Transaction{
		Value:       100,
		Type:        "c",
		Description: "test",
		CarriedOut:  "",
	}
	transaction2 := Transaction{
		Value:       50,
		Type:        "c",
		Description: "test",
		CarriedOut:  "",
	}
	customer := Customer{
		ID:           "1",
		Limit:        100,
		Transactions: []Transaction{},
	}
	err := customer.AddTransaction(transaction1)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = customer.AddTransaction(transaction2)

	assert.Equal(t, err.Error(), "invalid to proceed transaction. limit exceeded. available: 100")
}
