package api

import (
	"net/http"
	"context"

	"EffectiveMobile/internal/logger"
	
	"github.com/google/uuid"
)

type MiddleWareHandler struct {
	ctx context.Context
	IHandler
}

func NewMiddleWareHandler(handler IHandler) *MiddleWareHandler {
	return &MiddleWareHandler{
		IHandler: handler,
	}
}

func (mwlh *MiddleWareHandler) MiddleWareFunc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), logger.RequestID, uuid.New().String())
		r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}