package handler

import (
	"encoding/json"
	"errors"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/entity"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/service"
	"net/http"
	"strings"
)

type Handler struct {
	service *service.CustomerService
}

func NewServer(service *service.CustomerService) *Handler {
	return &Handler{service: service}
}

type AddTransactionInput struct {
	Value       int                    `json:"valor"`
	Type        entity.TransactionType `json:"tipo"`
	Description string                 `json:"descricao"`
}

func (s *Handler) HandleAddTransaction(w http.ResponseWriter, r *http.Request) {
	var data AddTransactionInput
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customerId := r.PathValue("id")
	var emptyDescription string
	if data.Description == emptyDescription {
		http.Error(w, errors.New("invalid transaction description").Error(), http.StatusUnprocessableEntity)
		return
	}
	transaction, err := entity.NewTransaction(data.Value, data.Type, data.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	output, err := s.service.AddTransaction(customerId, *transaction)
	if err != nil {
		if err.Error() == "customer not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if strings.HasPrefix(err.Error(), "invalid to proceed transaction. limit exceeded") {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		} else if strings.HasPrefix(err.Error(), "invalid transaction") {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		} else {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Handler) HandleGetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("id")
	output, err := s.service.GetTransactionHistory(customerId)
	if err != nil {
		if err.Error() == "customer not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{}`))
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{}`))
			return
		}
	}
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
