package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/api/handler"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/service"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/infra/database"
	"log"
	"net/http"
)

var pool *pgxpool.Pool

func main() {
	pool, _ = database.InitDB()
	defer pool.Close()
	repository := database.CustomerRepositoryPSQL{Pool: pool}

	customerService := service.CustomerService{
		Repository: &repository,
	}

	_handler := handler.NewServer(&customerService)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /clientes/{id}/extrato", _handler.HandleGetTransactionHistory)
	mux.HandleFunc("POST /clientes/{id}/transacoes", _handler.HandleAddTransaction)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
