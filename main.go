package main

import (
	"log"
	"net/http"

	handlers "onlineshop/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/json", handlers.JSONHandler).Methods("POST", "GET")

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
