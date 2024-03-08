package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/api/handler"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/service"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/infra/database"
)

var pool *pgxpool.Pool

func main() {
	gin.SetMode(gin.ReleaseMode)

	pool, _ = database.InitDB()
	defer pool.Close()
	repository := database.CustomerRepositoryPSQL{Pool: pool}

	customerService := service.CustomerService{
		Repository: &repository,
	}

	_handler := handler.NewServer(&customerService)

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET(
		"/clientes/:id/extrato",
		_handler.HandleGetTransactionHistory,
	)
	r.POST(
		"/clientes/:id/transacoes",
		_handler.HandleAddTransaction,
	)
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		return
	}
}
