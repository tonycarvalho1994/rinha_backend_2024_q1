package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
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

func (s *Handler) HandleAddTransaction(c *gin.Context) {
	var data AddTransactionInput
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	customerId := c.Param("id")
	var emptyDescription string
	if data.Description == emptyDescription {
		c.JSON(http.StatusUnprocessableEntity, errors.New("invalid transaction description"))
	}
	transaction, err := entity.NewTransaction(data.Value, data.Type, data.Description)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
	}

	output, err := s.service.AddTransaction(customerId, *transaction)
	if err != nil {
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, err)
		} else if strings.HasPrefix(err.Error(), "invalid to proceed transaction. limit exceeded") {
			c.JSON(http.StatusUnprocessableEntity, err)
		} else if strings.HasPrefix(err.Error(), "invalid transaction") {
			c.JSON(http.StatusUnprocessableEntity, err)
		} else {
			c.JSON(http.StatusUnprocessableEntity, err)
		}
	}

	c.JSON(http.StatusOK, output)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func (s *Handler) HandleGetTransactionHistory(c *gin.Context) {
	customerId := c.Param("id")
	output, err := s.service.GetTransactionHistory(customerId)
	if err != nil {
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, err)
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
	}

	c.JSON(http.StatusOK, output)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
}
