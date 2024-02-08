package main

import (
	"fmt"
	"net/http"

	"database/sql"

	"github.com/mgenteluci/rinha-2024-q1/pkg/types"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_ = validator.New(validator.WithRequiredStructEnabled())

	r := chi.NewRouter()

	r.Get("/clientes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		fmt.Printf("id: %s", id)
		var c types.Client

		err := db.QueryRow("SELECT client_limit, balance FROM clients WHERE id = $1", id).Scan(&c.ID, &c.Limit, &c.Balance)
		if err != nil {
			panic(err)
		}

		fmt.Printf("result: %v", *c)
	})

	r.Get("/clientes/{id}/extrato", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		fmt.Printf("id: %s", id)
	})

	r.Post("/clientes{id}/transacoes", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		fmt.Printf("id: %s", id)

		w.Write([]byte("welcome post"))
	})

	http.ListenAndServe(":8080", r)
}
