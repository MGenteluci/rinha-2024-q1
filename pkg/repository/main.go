package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/mgenteluci/rinha-2024-q1/pkg/types"
)

type ClientsRepository struct {
	database *pgx.Conn
}

func NewClientsRepository() *ClientsRepository {
	database, err := pgx.Connect(context.Background(), "host=postgres user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		panic(err)
	}

	return &ClientsRepository{database}
}

func (c *ClientsRepository) GetClient(clientID string) (*types.Client, error) {
	query := `SELECT id, client_limit, balance FROM clients WHERE id=$1`
	var client types.Client
	err := c.database.QueryRow(context.Background(), query, clientID).Scan(&client.ID, &client.Limit, &client.Balance)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("recurso nao encontrado")
		}
		return nil, err
	}

	return &client, nil
}

func (c *ClientsRepository) SaveTransaction(clientID string, transaction *types.NewTransactionRequestPayload) (*types.NewTransactionResponse, error) {
	tx, err := c.database.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	query := `SELECT id, client_limit, balance FROM clients WHERE id=$1 FOR UPDATE`
	var client types.Client
	err = tx.
		QueryRow(context.Background(), query, clientID).
		Scan(&client.ID, &client.Limit, &client.Balance)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, fmt.Errorf("recurso nao encontrado")
		}
		return nil, err
	}

	if transaction.Type == "d" {
		if absInt(client.Balance-transaction.Value) > client.Limit {
			return nil, fmt.Errorf("operação não permitida")
		}
	}

	query = `
		INSERT INTO transactions(client_id, transaction_value, transaction_type, transaction_description, transaction_date)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(
		context.Background(),
		query,
		clientID,
		transaction.Value,
		transaction.Type,
		transaction.Description,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	newBalance := NewBalance(&client, transaction)
	query = `UPDATE clients SET balance = $1 WHERE id = $2`
	_, err = tx.Exec(context.Background(), query, newBalance, clientID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &types.NewTransactionResponse{
		Balance: newBalance,
		Limit:   client.Limit,
	}, nil
}

func (c *ClientsRepository) GetClientDetails(clientID string) (*types.GetDetailsResponse, error) {
	query := `
		SELECT clients.client_limit, clients.balance, transactions.transaction_value, transactions.transaction_type, transactions.transaction_description, transactions.transaction_date
		FROM clients
		LEFT JOIN transactions ON transactions.client_id = clients.id
		WHERE clients.id = $1
		ORDER BY transactions.id DESC
		LIMIT 10
	`
	rows, err := c.database.Query(context.Background(), query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balance := types.GetDetailsBalance{SearchDate: time.Now()}
	transactions := []types.GetDetailsTransaction{}

	if !rows.Next() {
		return nil, fmt.Errorf("recurso nao encontrado")
	} else {
		transaction, err := scanTransaction(rows, &balance)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *transaction)
	}

	for rows.Next() {
		transaction, err := scanTransaction(rows, &balance)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *transaction)
	}

	response := types.GetDetailsResponse{
		Balance:          balance,
		LastTransactions: transactions,
	}

	return &response, nil
}

func scanTransaction(rows pgx.Rows, balance *types.GetDetailsBalance) (*types.GetDetailsTransaction, error) {
	var transaction types.GetDetailsTransaction
	err := rows.Scan(
		&balance.Limit,
		&balance.Total,
		&transaction.Value,
		&transaction.Type,
		&transaction.Description,
		&transaction.TransactionDate,
	)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func NewBalance(client *types.Client, transaction *types.NewTransactionRequestPayload) int {
	if transaction.Type == "d" {
		return client.Balance - transaction.Value
	}

	return client.Balance + transaction.Value
}
