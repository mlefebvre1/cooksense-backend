package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Person struct {
	ID        int64  `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
}

func main() {
	fmt.Println("test")

	db, err := sqlx.Connect("postgres", "user=cooksense password=cooksense dbname=cooksense sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	server := http.NewServeMux()

	server.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello World4!")
	})

	server.HandleFunc("/add-user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var p Person

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
			return
		}

		err := db.QueryRowx(
			`INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING id`,
			p.FirstName, p.LastName, p.Email,
		).Scan(&p.ID)

		if err != nil {
			http.Error(w, "insert failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	})

	server.HandleFunc("/get-user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "missing email query param", http.StatusBadRequest)
			return
		}

		var p Person
		err := db.Get(&p, `SELECT id, first_name, last_name, email FROM person WHERE email = $1`, email)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "query failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	})

	if err := http.ListenAndServe(":8081", server); err != nil {
		log.Fatal(err)
	}

}
