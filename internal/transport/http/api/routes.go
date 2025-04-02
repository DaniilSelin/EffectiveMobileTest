package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

type IHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Insert(w http.ResponseWriter, r *http.Request)
}

func NewRouter(h IHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("api/get", h.Get).Methods("GET")
	router.HandleFunc("api/delete", h.Delete).Methods("DELETE")
	router.HandleFunc("api/update", h.Update).Methods("PATCH")
	router.HandleFunc("api/insert", h.Insert).Methods("POST")

	return router
}