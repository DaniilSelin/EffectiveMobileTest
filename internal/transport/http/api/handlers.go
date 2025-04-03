package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"context"

	"EffectiveMobile/internal/models"
	"EffectiveMobile/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	PS *service.PersonService
	ctx context.Context
}

func NewHandler(ctx context.Context, ps *service.PersonService) *Handler {
	return &Handler{
		PS: ps,
		ctx: ctx,
	}	
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var filters models.PersonFilters
	if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
		http.Error(w, "Invalid filter data", http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}

	persons, err := h.PS.Get(r.Context(), filters, limitInt, offsetInt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get persons: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(persons); err != nil {
		http.Error(w, "Failed to encode persons", http.StatusInternalServerError)
	}
}


func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.PS.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete person: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var person models.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := h.PS.Update(r.Context(), person, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update person: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Insert(w http.ResponseWriter, r *http.Request) {
	var person models.Person

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	createdPerson, err := h.PS.Insert(r.Context(), person)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert person: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdPerson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
