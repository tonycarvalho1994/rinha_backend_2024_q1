package database

const (
	SelectBalanceQuery     = "SELECT limite, saldo FROM customers WHERE id=$1 FOR UPDATE"
	InsertTransactionQuery = "INSERT INTO transactions (customer_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, $5)"
	UpdateBalanceQuery     = "UPDATE customers SET saldo=$1 WHERE id=$2"
	GetHistoryLimit10Query = "SELECT valor, tipo, descricao, realizada_em FROM transactions WHERE customer_id = $1 ORDER BY realizada_em DESC LIMIT 10"
)
