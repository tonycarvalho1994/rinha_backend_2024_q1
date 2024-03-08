package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/entity"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	pass     = os.Getenv("DB_PASS")
	dbName   = os.Getenv("DB_NAME")
	maxConns = os.Getenv("DB_MAX_CONNECTIONS")
)

func InitDB() (*pgxpool.Pool, error) {
	stringFormat := "user=%s dbname=%s password=%s host=%s port=%s sslmode=disable"
	connectionString := fmt.Sprintf(stringFormat, user, dbName, pass, host, port)
	config, err := pgxpool.ParseConfig(connectionString)
	_max, err := strconv.Atoi(maxConns)
	config.MaxConns = int32(_max)

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	return pool, err
}

type CustomerRepositoryPSQL struct {
	Pool *pgxpool.Pool
}

func (c *CustomerRepositoryPSQL) CreateCustomer(id string, limit int) error {
	return nil
}

func (c *CustomerRepositoryPSQL) FindById(id string) (*entity.Customer, error) {
	return nil, nil
}

func (c *CustomerRepositoryPSQL) Update(customer *entity.Customer) error {
	return nil
}

func (c *CustomerRepositoryPSQL) AddTransaction(customerId string, transaction *entity.Transaction) (int, int, error) {
	newBalance := 0
	tx, err := c.Pool.Begin(context.Background())
	if err != nil {
		return newBalance, 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	var limit int
	var currentBalance int
	err = tx.QueryRow(context.Background(), "SELECT limite, saldo FROM customers WHERE id=$1 FOR UPDATE", customerId).Scan(&limit, &currentBalance)
	if err != nil {
		return newBalance, 0, errors.New("customer not found")
	}

	if transaction.Type == entity.Debit {
		newBalance = currentBalance - transaction.Value
		if newBalance < limit*-1 {
			return newBalance, 0, errors.New("invalid to proceed transaction. limit exceeded")
		}
	} else if transaction.Type == entity.Credit {
		newBalance = currentBalance + transaction.Value
		if newBalance > limit {
			return newBalance, 0, errors.New("invalid to proceed transaction. limit exceeded")
		}
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO transactions (customer_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, $5)",
		customerId, transaction.Value, transaction.Type, transaction.Description, time.Now())
	if err != nil {
		return newBalance, 0, err
	}

	_, err = tx.Exec(context.Background(), "UPDATE customers SET saldo=$1 WHERE id=$2", newBalance, customerId)
	if err != nil {
		return newBalance, 0, err
	}

	return newBalance, limit, tx.Commit(context.Background())
}

func (c *CustomerRepositoryPSQL) GetBalance(customerId string) (int, error) {
	var currentBalance int

	err := c.Pool.QueryRow(context.Background(), "SELECT COALESCE(SUM(CASE WHEN tipo='c' THEN valor ELSE -valor END), 0) FROM transactions WHERE customer_id=$1", customerId).Scan(&currentBalance)
	if err != nil {
		return currentBalance, errors.New("error while getting customer balance: " + err.Error())
	}

	return currentBalance, nil
}

func (c *CustomerRepositoryPSQL) FindHistoryByCustomerId(customerId string) (int, int, []entity.Transaction, error) {
	tx, err := c.Pool.Begin(context.Background())
	if err != nil {
		return 0, 0, nil, err
	}
	defer tx.Rollback(context.Background())

	var limit int
	var currentBalance int
	err = tx.QueryRow(context.Background(), "SELECT limite, saldo FROM customers WHERE id=$1 FOR UPDATE", customerId).Scan(&limit, &currentBalance)
	if err != nil {
		return 0, 0, nil, errors.New("customer not found")
	}

	var transactions []entity.Transaction
	rows, err := tx.Query(context.Background(), "SELECT valor, tipo, descricao, realizada_em FROM transactions WHERE customer_id = $1 ORDER BY realizada_em DESC LIMIT 10", customerId)
	if err != nil {
		return 0, 0, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var transaction entity.Transaction
		var carriedOut time.Time
		err := rows.Scan(&transaction.Value, &transaction.Type, &transaction.Description, &carriedOut)
		if err != nil {
			return 0, 0, nil, err
		}
		transaction.CarriedOut = carriedOut.Format(time.RFC3339)
		transaction.Description = strings.TrimSpace(transaction.Description)
		transactions = append(transactions, transaction)
	}
	if err := tx.Commit(context.Background()); err != nil {
		return 0, 0, nil, err
	}
	return limit, currentBalance, transactions, nil
}
