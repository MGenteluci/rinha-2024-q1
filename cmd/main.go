package main

import (
	"context"
	"net/http"

	"github.com/mgenteluci/rinha-2024-q1/pkg/handlers"
	"github.com/mgenteluci/rinha-2024-q1/pkg/repository"
	"github.com/mgenteluci/rinha-2024-q1/pkg/services"

	"github.com/go-playground/validator/v10"
)

var (
	validate          *validator.Validate
	clientsRepository *repository.ClientsRepository
	clientsService    *services.ClientsService
	clientsHandler    *handlers.ClientsHandler
)

func init() {
	ctx := context.Background()
	validate = validator.New(validator.WithRequiredStructEnabled())

	clientsRepository = repository.NewClientsRepository(ctx)
	clientsService = services.NewClientsService(validate, clientsRepository)
	clientsHandler = handlers.NewClientsHandler(clientsService)
}

func main() {
	defer clientsRepository.Close()

	http.HandleFunc("GET /clientes/{id}/extrato", clientsHandler.GetClientDetails)
	http.HandleFunc("POST /clientes/{id}/transacoes", clientsHandler.CreateTransaction)

	http.ListenAndServe(":8080", nil)
}
