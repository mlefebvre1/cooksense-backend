package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("test")
	server := http.NewServeMux()

	server.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello World4!")
	})

	if err := http.ListenAndServe(":8081", server); err != nil {
		log.Fatal(err)
	}

}
