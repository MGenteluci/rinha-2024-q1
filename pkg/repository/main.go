package repository

import (
	"fmt"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mgenteluci/rinha-2024-q1/pkg/types"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

type ClientsRepository struct {
	databaseURL string
	database    *sql.DB
}

func NewClientsRepository(databaseURL string) *ClientsRepository {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	database, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	return &ClientsRepository{databaseURL, database}
}

func (c *ClientsRepository) GetClient(clientID string) (*types.Client, error) {
	query := `SELECT id, client_limit, balance FROM clients WHERE id=$1`
	var client types.Client
	err := c.database.QueryRow(query, clientID).Scan(&client.ID, &client.Limit, &client.Balance)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("recurso nao encontrado")
		}
		return nil, err
	}

	return &client, nil
}

func (c *ClientsRepository) SaveTransaction(clientID string, clientBalance int, transaction *types.NewTransactionRequestPayload) (*types.NewTransactionResponse, error) {
	query := `
		INSERT INTO transactions(client_id, transaction_value, transaction_type, transaction_description, transaction_date)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := c.database.Exec(
		query,
		clientID,
		transaction.Value,
		transaction.Type,
		transaction.Description,
		time.Now(),
	)
	if err != nil {
		panic(err)
	}

	query = `UPDATE clients SET balance = $1 WHERE id = $2`
	_, err = c.database.Exec(query, clientBalance, clientID)
	if err != nil {
		panic(err)
	}

	return nil, nil
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
	rows, err := c.database.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balance := types.GetDetailsBalance{SearchDate: time.Now()}
	transactions := []types.GetDetailsTransaction{}

	if !rows.Next() {
		return nil, fmt.Errorf("recurso nao encontrado")
	} else {
		transactions = append(transactions, *scanTransaction(rows, &balance))
	}

	for rows.Next() {
		transactions = append(transactions, *scanTransaction(rows, &balance))
	}

	response := types.GetDetailsResponse{
		Balance:          balance,
		LastTransactions: transactions,
	}

	return &response, nil
}

func scanTransaction(rows *sql.Rows, balance *types.GetDetailsBalance) *types.GetDetailsTransaction {
	var transaction types.GetDetailsTransaction
	rows.Scan(
		&balance.Limit,
		&balance.Total,
		&transaction.Value,
		&transaction.Type,
		&transaction.Description,
		&transaction.TransactionDate,
	)
	return &transaction
}
