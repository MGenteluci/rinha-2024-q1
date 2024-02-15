package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mgenteluci/rinha-2024-q1/pkg/services"
	"github.com/mgenteluci/rinha-2024-q1/pkg/types"
)

type ClientsHandler struct {
	clientsService *services.ClientsService
}

func NewClientsHandler(clientsService *services.ClientsService) *ClientsHandler {
	return &ClientsHandler{clientsService}
}

func (c *ClientsHandler) GetClientDetails(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	details, err := c.clientsService.GetClientDetails(id)
	if err != nil {
		if err.Error() == "recurso nao encontrado" {
			WriteErrorResponse(w, 404, err.Error())
			return
		}

		WriteErrorResponse(w, 500, err.Error())
		return
	}

	WriteHttpResponse(w, 200, details)
}

func (c *ClientsHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction types.NewTransactionRequestPayload
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		WriteErrorResponse(w, 422, "transacao invalida")
		return
	}

	id := r.PathValue("id")

	response, err := c.clientsService.SaveTransaction(id, &transaction)
	if err != nil {
		if err.Error() == "recurso nao encontrado" {
			WriteErrorResponse(w, 404, "recurso nao encontrado")
			return
		}

		WriteErrorResponse(w, 422, err.Error())
		return
	}

	WriteHttpResponse(w, 200, response)
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	WriteHttpResponse(w, statusCode, `{"message":"`+message+`"}`)
}

func WriteHttpResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
