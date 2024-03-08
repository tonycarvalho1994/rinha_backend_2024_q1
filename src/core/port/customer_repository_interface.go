package port

import "github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/entity"

type CustomerRepositoryInterface interface {
	AddTransaction(customerId string, transaction *entity.Transaction) (int, int, error)
	FindHistoryByCustomerId(customerId string) (int, int, []entity.Transaction, error)
}
