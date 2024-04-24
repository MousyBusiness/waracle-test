package routes

import (
	"encoding/json"
	"fmt"
	"github.com/mousybusiness/waracle-test/internal/dao"
	"github.com/mousybusiness/waracle-test/internal/db"
	"github.com/mousybusiness/waracle-test/internal/handler/middleware"
	"io"
	"log"
	"net/http"
	"strconv"
)

func BuildRoutes(database db.Database) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "healthcheck\n")
	})

	// Users should be able to fetch a cake by ID.
	router.HandleFunc("GET /cake/{id}/", func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("GET /cakes/", func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("POST /cakes/search/", func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("PUT /cake/", func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("DELETE /cake/{id}/", func(w http.ResponseWriter, r *http.Request) {
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

	// Add route versioning
	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	// Apply middlewares
	middlewares := middleware.Create(
		middleware.Logging,
		middleware.IsAuthenticated,
	)

	return middlewares(v1)
}
