package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mgenteluci/rinha-2024-q1/pkg/types"
)

type ClientsRepository struct {
	database *pgxpool.Pool
	ctx      context.Context
}

func NewClientsRepository(ctx context.Context) *ClientsRepository {
	database, err := pgxpool.New(ctx, "host=postgres user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		panic(err)
	}

	return &ClientsRepository{database, ctx}
}

func (c *ClientsRepository) SaveTransaction(clientID string, transaction *types.NewTransactionRequestPayload) (*types.NewTransactionResponse, error) {
	tx, err := c.database.Begin(c.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c.ctx)

	var limit, balance, newBalance int

	query := `SELECT client_limit, balance FROM clients WHERE id=$1 FOR UPDATE`
	err = tx.
		QueryRow(c.ctx, query, clientID).
		Scan(&limit, &balance)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, fmt.Errorf("recurso nao encontrado")
		}
		return nil, err
	}

	if transaction.Type == "d" {
		newBalance = balance - transaction.Value
	} else {
		newBalance = balance + transaction.Value
	}

	if (limit + newBalance) < 0 {
		return nil, fmt.Errorf("operaçao nao permitida")
	}

	batch := pgx.Batch{}
	batch.Queue(
		`
		INSERT INTO transactions(client_id, transaction_value, transaction_type, transaction_description, transaction_date)
		VALUES ($1, $2, $3, $4, $5)
		`,
		clientID,
		transaction.Value,
		transaction.Type,
		transaction.Description,
		time.Now(),
	)
	batch.Queue(
		`UPDATE clients SET balance = $1 WHERE id = $2`,
		newBalance,
		clientID,
	)

	batchResults := tx.SendBatch(c.ctx, &batch)
	_, err = batchResults.Exec()
	if err != nil {
		return nil, err
	}
	err = batchResults.Close()
	if err != nil {
		return nil, err
	}

	err = tx.Commit(c.ctx)
	if err != nil {
		return nil, err
	}

	return &types.NewTransactionResponse{
		Balance: newBalance,
		Limit:   limit,
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
	rows, err := c.database.Query(c.ctx, query, clientID)
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

func (c *ClientsRepository) Close() {
	c.database.Close()
}
