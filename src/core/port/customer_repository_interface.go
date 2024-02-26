package port

import "github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/entity"

type CustomerRepositoryInterface interface {
	CreateCustomer(id string, limit float64) error
	FindById(id string) (*entity.Customer, error)
	Update(customer *entity.Customer) error
	AddTransaction(customerId string, transaction *entity.Transaction) error
}
