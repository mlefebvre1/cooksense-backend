package main

import (
	_ "database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS person (
    first_name text,
    last_name text,
    email text
);`

func main() {
	fmt.Println("test")

	db, err := sqlx.Connect("postgres", "user=cooksense password=cooksense dbname=cooksense sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	fmt.Println("Bob")

	server := http.NewServeMux()

	server.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello World4!")
	})

	if err := http.ListenAndServe(":8081", server); err != nil {
		log.Fatal(err)
	}

}
