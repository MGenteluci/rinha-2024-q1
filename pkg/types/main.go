package types

import "time"

type Client struct {
	ID      int `json:"id"`
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}

type NewTransactionRequestPayload struct {
	Value       int    `json:"valor" validate:"required,gt=0"`
	Type        string `json:"tipo" validate:"required,oneof=c d"`
	Description string `json:"descricao" validate:"required,gte=1,lte=10"`
}

type NewTransactionResponse struct {
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}

type GetDetailsResponse struct {
	Balance          GetDetailsBalance       `json:"saldo"`
	LastTransactions []GetDetailsTransaction `json:"ultimas_transacoes"`
}

type GetDetailsBalance struct {
	Total      *int      `json:"total"`
	SearchDate time.Time `json:"data_extrato"`
	Limit      *int      `json:"limite"`
}

type GetDetailsTransaction struct {
	Value           *int      `json:"valor"`
	Type            *string   `json:"tipo"`
	Description     *string   `json:"descricao"`
	TransactionDate time.Time `json:"realizada_em"`
}

type HttpMessage struct {
	Message string `json:"mensagem"`
}
