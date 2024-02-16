package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mgenteluci/rinha-2024-q1/pkg/services"
	"github.com/mgenteluci/rinha-2024-q1/pkg/types"
)

type ClientsHandler struct {
	clientsService *services.ClientsService
}

func NewClientsHandler(clientsService *services.ClientsService) *ClientsHandler {
	return &ClientsHandler{clientsService}
}

func (h *ClientsHandler) GetClientDetails(c *fiber.Ctx) error {
	id := c.Params("id")
	details, err := h.clientsService.GetClientDetails(id)
	if err != nil {
		if err.Error() == "recurso nao encontrado" {
			return c.SendStatus(404)
		}

		return c.SendStatus(500)
	}

	return c.Status(200).JSON(details)
}

func (h *ClientsHandler) CreateTransaction(c *fiber.Ctx) error {
	var transaction types.NewTransactionRequestPayload
	if err := c.BodyParser(&transaction); err != nil {
		return c.SendStatus(422)
	}

	id := c.Params("id")

	response, err := h.clientsService.SaveTransaction(id, &transaction)
	if err != nil {
		if err.Error() == "recurso nao encontrado" {
			return c.SendStatus(404)
		}

		return c.SendStatus(422)
	}

	return c.Status(200).JSON(response)
}
