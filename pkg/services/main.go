package services

import (
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

	return c.clientsRepository.SaveTransaction(clientID, transaction)
}

func (c *ClientsService) GetClientDetails(clientID string) (*types.GetDetailsResponse, error) {
	return c.clientsRepository.GetClientDetails(clientID)
}

func getNewBalance(client *types.Client, transaction *types.NewTransactionRequestPayload) int {
	if transaction.Type == "d" {
		return client.Balance - transaction.Value
	}

	return client.Balance + transaction.Value
}
