package api 

import (
	"encoding/json"
	"net/http"
	"strings"
	"log"

	"EffectiveMobile/internal/models"
	"EffectiveMobile/internal/service"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

type Handler struct {
	US *service.UserService
}

func NewHandler(us *service.UserService) *Handler {
	return &Handler{
		US: us,
	}	
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
}

func (h Handler) Insert(w http.ResponseWriter, r *http.Request) {
}