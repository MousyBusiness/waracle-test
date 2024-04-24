package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mousybusiness/waracle-test/internal/dao"
	"github.com/mousybusiness/waracle-test/internal/db"
	"github.com/mousybusiness/waracle-test/internal/handler/middleware"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {

	database, err := db.ConnectToDatastore(context.Background())
	if err != nil {
		log.Panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "healthcheck\n")
	})

	// Users should be able to fetch a cake by ID.
	mux.HandleFunc("GET /cake/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		atoi, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cake, err := dao.GetCake(r.Context(), database, atoi)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		b, err := json.Marshal(cake)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(b); err != nil {
			log.Println(err)
			return
		}
	})

	// Users should be able to list all cakes.
	mux.HandleFunc("GET /cakes/", func(w http.ResponseWriter, r *http.Request) {
		cakes, err := dao.ListCakes(r.Context(), database)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println(cakes)

		if len(cakes) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		b, err := json.Marshal(cakes)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(b); err != nil {
			log.Println(err)
			return
		}
	})

	// Users should be able to search for a cake by yumFactor and/or name.
	mux.HandleFunc("POST /cakes/search/", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var request dao.SearchRequest
		if err := json.Unmarshal(b, &request); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if request.Name == nil && request.YumFactor == nil {
			log.Println("Require name or yum_factor")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cakes, err := dao.SearchCakes(r.Context(), database, request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNoContent) // TODO check error types
			return
		}

		fmt.Println(cakes)

		if len(cakes) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		b, err = json.Marshal(cakes)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(b); err != nil {
			log.Println(err)
			return
		}
	})

	// Users should be able to add another cake.
	mux.HandleFunc("PUT /cake/", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var cake dao.Cake
		if err := json.Unmarshal(b, &cake); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := dao.CreateCake(r.Context(), database, &cake); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// Users should be able to delete an existing cake.
	mux.HandleFunc("DELETE /cake/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		atoi, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := dao.DeleteCake(r.Context(), database, atoi); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", mux))

	middlewares := middleware.Create(middleware.Logging)

	server := http.Server{
		Addr:              ":8080",
		Handler:           middlewares(v1),
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       90 * time.Second,
	}
	fmt.Println("Starting listening on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
