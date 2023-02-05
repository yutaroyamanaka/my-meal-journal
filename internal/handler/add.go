package handler

import (
	"context"
	"encoding/json"
	"net/http"

	log "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/yutaroyamanaka/my-httpserver-monitoring/internal/entity"
)

type Service interface {
	Add(string, int) (*entity.Journal, error)
}

func AddHandler(ctx context.Context, svc Service, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			_, err := w.Write([]byte("Only post method is allowed"))
			if err != nil {
				level.Error(logger).Log("msg", "failed to write body to response writer", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
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
