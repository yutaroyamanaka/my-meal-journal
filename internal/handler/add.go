// Package handler has functions which handle application logics.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/yutaroyamanaka/my-meal-journal/internal/entity"
	"github.com/yutaroyamanaka/my-meal-journal/internal/service"
)

// Error messages
var (
	ErrDiallowedMethod    = "only post method is allowed"
	ErrInvalidRequestBody = "request body is invalid"
	ErrName               = "name must be a non-empty string"
	ErrCategory           = fmt.Sprintf("category must be integer between %d (%s) and %d (%s)",
		entity.Breakfast, entity.CategoryBreakfast, entity.Others, entity.CateogoryOthersName)
	ErrUnknown = "An unknown errors is occurred"
)

// Service receives meal's information and returns the error of business logic.
type Service interface {
	Add(context.Context, string, int) error
}

var _ Service = (*service.AddService)(nil)

// AddFunc is a stub function for mocking Service interface.
type AddFunc func(context.Context, string, int) error

// Add receives meal's information and returns the error of business logic.
func (f AddFunc) Add(ctx context.Context, name string, category int) error {
	return f(ctx, name, category)
}

func writeErrorResponse(w http.ResponseWriter, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	body, _ := json.Marshal(map[string]string{"error": msg})
	w.Write(body)
}

// NewAddHandler returns http.Handler which paases meal's information to business logic
// and returns the response with the record information.
func NewAddHandler(svc Service, logger log.Logger) (http.Handler, error) {
	if svc == nil {
		return nil, errors.New("serivce must not be nil")
	}
	if logger == nil {
		return nil, errors.New("logger must not be nil")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		if r.Method != http.MethodPost {
			writeErrorResponse(w, ErrDiallowedMethod, http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			level.Error(logger).Log("msg", "failed to read a new body from request", "err", err)
			writeErrorResponse(w, ErrUnknown, http.StatusInternalServerError)
			return
		}
		defer func() {
			io.Copy(io.Discard, r.Body)
			r.Body.Close()
		}()
		var i struct {
			Name     *string `json:"name"`
			Category *int    `json:"category"`
		}
		if err := json.Unmarshal(body, &i); err != nil {
			level.Error(logger).Log("msg", "failed to unmarshal request body", "err", err)
			writeErrorResponse(w, ErrInvalidRequestBody, http.StatusBadRequest)
			return
		}
		if i.Name == nil || *i.Name == "" {
			writeErrorResponse(w, ErrName, http.StatusBadRequest)
			return
		}
		if i.Category == nil || *i.Category < entity.Breakfast || *i.Category > entity.Others {
			writeErrorResponse(w, ErrCategory, http.StatusBadRequest)
			return
		}

		err = svc.Add(r.Context(), *i.Name, *i.Category)
		if err != nil {
			level.Error(logger).Log("msg", "failed to create a new journal", "err", err)
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}), nil
}
