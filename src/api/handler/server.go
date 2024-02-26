package handler

import (
	"encoding/json"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/entity"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/service"
	"io/ioutil"
	"net/http"
	"strings"
)

type Server struct {
	service *service.CustomerService
	Mux     *http.ServeMux
}

func NewServer(service *service.CustomerService) *Server {
	return &Server{service: service}
}

func (s *Server) SetupRouter() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /clientes/{id}/transacoes", s.HandleAddTransaction)
	mux.HandleFunc("GET /clientes/{id}/extrato", s.HandleGetTransactionHistory)

	s.Mux = mux
}

type AddTransactionInput struct {
	Value       float64                `json:"valor"`
	Type        entity.TransactionType `json:"tipo"`
	Description string                 `json:"descricao"`
}

func (s *Server) HandleAddTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	var data AddTransactionInput
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	customerId := r.PathValue("id")
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
		} else if strings.HasPrefix(err.Error(), "invalid transaction type") {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	jsonResponse, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleGetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	customerId := r.PathValue("id")
	output, err := s.service.GetTransactionHistory(customerId)
	if err != nil {
		if err.Error() == "customer not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	jsonResponse, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
