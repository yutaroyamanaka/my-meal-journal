package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	log "github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
	"github.com/yutaroyamanaka/my-httpserver-monitoring/internal/entity"
)

func TestAddHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		status int
		f      AddFunc
		want   string
	}{
		{
			name:   "normal",
			method: http.MethodPost,
			status: http.StatusOK,
			f: func(string, int) (*entity.Journal, error) {
				return &entity.Journal{ID: 0, Name: "sunny side up", Cateogry: 0, Created: time.Date(2023, 2, 5, 16, 27, 56, 0, time.UTC)}, nil
			},
			want: `{"id":0,"name":"sunny side up","category":0,"created":"2023-02-05T16:27:56Z"}`,
		},
		{
			name:   "method not allowed",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			f:      func(string, int) (*entity.Journal, error) { return nil, nil },
			want:   "",
		},
		{
			name:   "service returns an error",
			method: http.MethodPost,
			status: http.StatusInternalServerError,
			f:      func(string, int) (*entity.Journal, error) { return nil, errors.New("unexpected error is occurred") },
			want:   "failed to register your meal's information\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := AddHandler(context.Background(), tt.f, log.NewLogfmtLogger(os.Stderr))
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, "http://localhost/add", nil)
			h.ServeHTTP(w, r)

			resp := w.Result()
			defer func() {
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			}()
			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Got an unexpected error: %#v", err)
			}
			if diff := cmp.Diff(resp.StatusCode, tt.status); diff != "" {
				t.Errorf("Got an unexcpeted status code: %s", diff)
			}
			if diff := cmp.Diff(resp.Header.Get("content-type"), "application/json"); diff != "" {
				t.Errorf("Got an unexcpeted content-type: %s", diff)
			}
			if diff := cmp.Diff(string(got), tt.want); diff != "" {
				t.Errorf("Got an unexcpeted state: %s", diff)
			}
		})
	}
}
