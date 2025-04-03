package api

import (
	"net/http"
	"context"

	"github.com/gorilla/mux"
)

type IHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Insert(w http.ResponseWriter, r *http.Request)
}

func NewRouter(ctx context.Context, h IHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/persons", h.Get).Methods("GET")
	router.HandleFunc("/api/persons/{id}", h.Delete).Methods("DELETE")
	router.HandleFunc("/api/persons/{id}", h.Update).Methods("PATCH")
	router.HandleFunc("/api/persons", h.Insert).Methods("POST")

	return router
}