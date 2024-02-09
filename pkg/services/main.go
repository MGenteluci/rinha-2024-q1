package services

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/mgenteluci/rinha-2024-q1/pkg/repository"
	"github.com/mgenteluci/rinha-2024-q1/pkg/types"
)

type ClientsService struct {
	validate          *validator.Validate
	clientsRepository *repository.ClientsRepository
}

func NewClientsService(validate *validator.Validate, clientsRepository *repository.ClientsRepository) *ClientsService {
	return &ClientsService{validate, clientsRepository}
}

func (c *ClientsService) SaveTransaction(clientID string, transaction *types.NewTransactionRequestPayload) (*types.NewTransactionResponse, error) {
	err := c.validate.Struct(transaction)
	if err != nil {
		return nil, err
	}

	client, err := c.clientsRepository.GetClient(clientID)
	if err != nil {
		return nil, err
	}

	err = c.ValidateTransaction(client, transaction)
	if err != nil {
		return nil, err
	}

	newBalance := getNewBalance(client, transaction)
	c.clientsRepository.SaveTransaction(clientID, newBalance, transaction)

	response := types.NewTransactionResponse{
		Limit:   client.Limit,
		Balance: newBalance,
	}
	return &response, nil
}

func (c *ClientsService) GetClientDetails(clientID string) (*types.GetDetailsResponse, error) {
	return c.clientsRepository.GetClientDetails(clientID)
}

func (c *ClientsService) ValidateTransaction(client *types.Client, transaction *types.NewTransactionRequestPayload) error {
	if transaction.Type == "d" {
		if absInt(client.Balance-transaction.Value) > client.Limit {
			return errors.New("operação não permitida")
		}
	}

	return nil
}

func getNewBalance(client *types.Client, transaction *types.NewTransactionRequestPayload) int {
	if transaction.Type == "d" {
		return client.Balance - transaction.Value
	}

	return client.Balance + transaction.Value
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
