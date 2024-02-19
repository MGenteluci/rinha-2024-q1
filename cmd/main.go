package main

import (
	"context"

	"github.com/mgenteluci/rinha-2024-q1/pkg/handlers"
	"github.com/mgenteluci/rinha-2024-q1/pkg/repository"
	"github.com/mgenteluci/rinha-2024-q1/pkg/services"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
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

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Get("/clientes/:id/extrato", clientsHandler.GetClientDetails)
	app.Post("/clientes/:id/transacoes", clientsHandler.CreateTransaction)

	app.Listen(":8080")
}
