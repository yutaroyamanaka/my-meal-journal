// Package handler has functions which handle application logics.
package handler

import (
	"context"
	"encoding/json"
	"net/http"

	log "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/yutaroyamanaka/my-httpserver-monitoring/internal/entity"
)

// Service receives meal's information, and returns the record and error.
type Service interface {
	Add(string, int) (*entity.Journal, error)
}

// AddFunc is a stub function for mocking Service interface.
type AddFunc func(string, int) (*entity.Journal, error)

// Add receives meal's information, and returns the record and error.
func (f AddFunc) Add(name string, category int) (*entity.Journal, error) {
	return f(name, category)
}

// AddHandler returns http.Handler which paases meal's information to business logic
// and returns the response with the record information.
func AddHandler(ctx context.Context, svc Service, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		ret, err := svc.Add("sunny side up", entity.Breakfast)
		if err != nil {
			level.Error(logger).Log("msg", "failed to create a new journal", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body, err := json.Marshal(ret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(body)
		if err != nil {
			level.Error(logger).Log("msg", "failed to write body to response writer", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
