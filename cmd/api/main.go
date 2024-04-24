package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Users should be able to fetch a cake by ID.
	mux.HandleFunc("GET /v1/cake/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "got cake with id: %v\n", id)
	})

	// Users should be able to list all cakes.
	mux.HandleFunc("GET /v1/cakes/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "list cakes\n")
	})

	// Users should be able to search for a cake by yumFactor and/or name.
	mux.HandleFunc("POST /v1/cakes/search/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "search cakes\n")
	})

	// Users should be able to add another cake.
	mux.HandleFunc("PUT /v1/cake/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "create cake\n")
	})

	// Users should be able to delete an existing cake.
	mux.HandleFunc("DELETE /v1/cake/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "delete cake: %v\n", id)
	})

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Panic(err)
	}
}
