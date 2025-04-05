package api

import (
	"encoding/json"
	"fmt"
	"errors"
	"net/http"
	"strconv"
	"context"

	log "EffectiveMobile/internal/logger"
	"EffectiveMobile/internal/models"
	"EffectiveMobile/internal/service"
	// исключительно для ErrNotFound
	"EffectiveMobile/internal/repository"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Исключительно для swagger-а
type ErrorResponse struct {
    Message string `json:"message"`
    Code    int    `json:"code"`
}

type IPersonService interface {
	Get(context.Context, models.PersonFilters, int, int) (*[]models.Person, error)
	Insert(context.Context, models.Person) (*models.Person, error)
	Delete(context.Context, string) error
	Update(context.Context, models.Person, string) error
}

type Handler struct {
	PS IPersonService
	ctx context.Context
}

func NewHandler(ctx context.Context, ps *service.PersonService) *Handler {
	return &Handler{
		PS: ps,
		ctx: ctx,
	}	
}

func (h *Handler) GenerateRequestID(ctx context.Context) context.Context {
    ctx = context.WithValue(ctx, log.RequestID, uuid.New().String())
    return ctx
}

// Get godoc
// @Summary Получить список людей
// @Description Получение списка людей с фильтрами и пагинацией
// @Tags persons
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param filters body models.PersonFilters false "Фильтры"
// @Success 200 {array} models.Person "Список людей"
// @Failure 400 {object} ErrorResponse "Ошибка валидации данных"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /persons [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := h.GenerateRequestID(h.ctx)
	logger := log.GetLoggerFromCtx(ctx)

	logger.Info(
		ctx, "Received new request",
		zap.String("url", r.URL.String()), zap.String("method", r.Method),
	)

	var filters models.PersonFilters
	if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
		logger.Info(
			ctx, "Failed to decode filter data",
			zap.Error(err),
		)
		http.Error(w, "Invalid filter data", http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	logger.Debug(
		ctx, "Query parameters",
		zap.String("limit", limit), zap.String("offset", offset),
	)

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		logger.Info(
			ctx, "Invalid query parameters", 
			zap.String("limit", limit), zap.String("offset", offset),
		)
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		logger.Info(
			ctx, "Invalid query parameters", 
			zap.String("limit", limit), zap.String("offset", offset),
		)		
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}

	persons, err := h.PS.Get(r.Context(), filters, limitInt, offsetInt)
	if err != nil {
		logger.Info(
			ctx, "Failed to retrieve persons", 
			zap.Error(err),
		)
		http.Error(w, fmt.Sprintf("Failed to get persons: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(persons); err != nil {
		logger.Error(
			ctx, "Failed to encode persons", 
			zap.Error(err),
		)
		http.Error(w, "Failed to encode persons", http.StatusInternalServerError)
	}

	logger.Info(
		ctx, "Successfully retrieved persons",
	)
}

// Delete godoc
// @Summary Удалить человека по ID
// @Description Удаляет человека из базы данных по указанному ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path string true "ID человека"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse "Ошибка: ID обязателен"
// @Failure 404 {object} ErrorResponse "Человек не найден"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /persons/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := h.GenerateRequestID(h.ctx)
	logger := log.GetLoggerFromCtx(ctx)
	
	logger.Info(
		ctx, "Received DELETE request",
		zap.String("url", r.URL.String()), zap.String("method", r.Method),
	)

	vars := mux.Vars(r)
	id := vars["id"]

	logger.Debug(
		ctx, "Deleting person with ID",
		zap.String("id", id),
	)

	if id == "" {
		logger.Info(
			ctx, "Id is empty", 
		)
		http.Error(w, "Id is required", http.StatusBadRequest)
        return
	}

	err := h.PS.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Info(
				ctx, "Person not found", 
				zap.Error(err),
			)
	        http.Error(w, "Person not found", http.StatusNotFound)
	        return
	    }

	    logger.Info(
			ctx, "Failed to delete person", 
			zap.Error(err),
		)
	    http.Error(w, fmt.Sprintf("Failed to delete person: %v", err), http.StatusInternalServerError)
	    return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info(
		ctx, "Successfully deleted person",
	)
}

// Update godoc
// @Summary Обновить информацию о человеке
// @Description Обновляет данные человека по указанному ID. Если поле не указано, оно остается без изменений.
// @Tags persons
// @Accept json
// @Produce json
// @Param id path string true "ID человека"
// @Param person body models.Person true "Данные для обновления"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse "Ошибка: ID обязателен или некорректные данные"
// @Failure 404 {object} ErrorResponse "Человек не найден"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /persons/{id} [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := h.GenerateRequestID(h.ctx)
	logger := log.GetLoggerFromCtx(ctx)
	
	logger.Info(
		ctx, "Received new request",
		zap.String("url", r.URL.String()), zap.String("method", r.Method),
	)

	vars := mux.Vars(r)
	id := vars["id"]

	logger.Debug(
		ctx, "Updating person with ID",
		zap.String("id", id),
	)

	if id == "" {
		logger.Info(
			ctx, "Id is empty", 
		)
		http.Error(w, "Id is required", http.StatusBadRequest)
        return
	}

	var person models.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Info(
				ctx, "Person not found",
			)
	        http.Error(w, "Person not found", http.StatusNotFound)
	        return
	    }

	    logger.Info(
			ctx, "Failed to decode input data",
			zap.Error(err),
		)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := h.PS.Update(r.Context(), person, id)
	if err != nil {
		logger.Info(
			ctx, "Failed to update person", 
			zap.Error(err),
		)
		http.Error(w, fmt.Sprintf("Failed to update person: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info(
		ctx, "Successfully updated person",
	)
}

// Insert godoc
// @Summary Создать нового человека
// @Description Создание нового человека в базе данных
// @Tags persons
// @Accept json
// @Produce json
// @Param person body models.Person true "Новый человек"
// @Success 201 {object} models.Person "Человек успешно создан"
// @Failure 400 {object} ErrorResponse "Ошибка валидации данных"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /persons [post]
func (h *Handler) Insert(w http.ResponseWriter, r *http.Request) {
	ctx := h.GenerateRequestID(h.ctx)
	logger := log.GetLoggerFromCtx(ctx)

	logger.Info(
		ctx, "Received new request", 
		zap.String("url", r.URL.String()), zap.String("method", r.Method),
	)

	var person models.Person

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		logger.Info(
			ctx, "Failed to decode input data",
			zap.Error(err),
		)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	createdPerson, err := h.PS.Insert(r.Context(), person)
	if err != nil {
		logger.Info(
			ctx, "Failed to insert person", 
			zap.Error(err),
		)
		http.Error(w, fmt.Sprintf("Failed to insert person: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdPerson); err != nil {
		logger.Error(
			ctx, "Failed to encode response", 
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

	logger.Info(
		ctx, "Successfully created person", 
	)
}
