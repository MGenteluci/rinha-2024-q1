package main

import (
	"net/http"

	"github.com/mgenteluci/rinha-2024-q1/pkg/controllers"
	"github.com/mgenteluci/rinha-2024-q1/pkg/repository"
	"github.com/mgenteluci/rinha-2024-q1/pkg/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
)

var (
	validate          *validator.Validate
	clientsRepository *repository.ClientsRepository
	clientsService    *services.ClientsService
	clientsController *controllers.ClientsController
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	clientsRepository = repository.NewClientsRepository("")
	clientsService = services.NewClientsService(validate, clientsRepository)
	clientsController = controllers.NewClientsController(clientsService)
}

func main() {
	r := chi.NewRouter()

	r.Get("/clientes/{id}/extrato", clientsController.GetClientDetails)
	r.Post("/clientes/{id}/transacoes", clientsController.CreateTransaction)

	http.ListenAndServe(":8080", r)
}
