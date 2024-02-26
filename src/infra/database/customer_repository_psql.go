package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/entity"
	"os"
	"time"
)

var (
	host   = os.Getenv("DB_HOST")
	port   = os.Getenv("DB_PORT")
	user   = os.Getenv("DB_USER")
	pass   = os.Getenv("DB_PASS")
	dbName = os.Getenv("DB_NAME")
)

func CreateDatabasePSQL() (*sql.DB, error) {
	stringFormat := "user=%s dbname=%s password=%s host=%s port=%s sslmode=disable"
	//stringFormat := "user=postgres dbname=rinha2024q1 password=postgres host=localhost port=5432 sslmode=disable"
	connectionString := fmt.Sprintf(stringFormat, user, dbName, pass, host, port)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type CustomerRepositoryPSQL struct {
	Db *sql.DB
}

func (c *CustomerRepositoryPSQL) CreateCustomer(id string, limit float64) error {
	_, err := c.Db.Exec("INSERT INTO customers(id, limite) VALUES($1, $2)", id, limit)
	return err
}

func (c *CustomerRepositoryPSQL) FindById(id string) (*entity.Customer, error) {
	customer := entity.Customer{Transactions: make([]entity.Transaction, 0)}
	err := c.Db.QueryRow("SELECT id, limite FROM customers WHERE id = $1 FOR UPDATE", id).Scan(&customer.ID, &customer.Limit)
	if err != nil {
		return nil, err
	}

	var transactions []entity.Transaction
	rows, err := c.Db.Query("SELECT valor, tipo, descricao, realizada_em FROM transactions WHERE customer_id = $1 ORDER BY realizada_em DESC", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var transaction entity.Transaction
		var carriedOut time.Time
		err := rows.Scan(&transaction.Value, &transaction.Type, &transaction.Description, &carriedOut)
		if err != nil {
			return nil, err
		}
		transaction.CarriedOut = carriedOut.Format(time.RFC3339)
		transactions = append(transactions, transaction)
	}
	customer.Transactions = transactions

	return &customer, nil
}

func (c *CustomerRepositoryPSQL) Update(customer *entity.Customer) error {
	return nil
}

func (c *CustomerRepositoryPSQL) AddTransaction(customerId string, transaction *entity.Transaction) error {
	_, err := c.Db.Exec("INSERT INTO transactions(customer_id, valor, tipo, descricao, realizada_em) VALUES($1, $2, $3, $4, $5)",
		customerId, transaction.Value, transaction.Type, transaction.Description, transaction.CarriedOut)
	return err
}
