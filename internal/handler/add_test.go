package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	log "github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
	"github.com/yutaroyamanaka/my-meal-journal/internal/util"
)

func TestNewAddNewHandler_error(t *testing.T) {
	tests := []struct {
		name   string
		svc    Service
		logger log.Logger
		want   error
	}{
		{
			name:   "serivce must not be nil",
			svc:    nil,
			logger: log.NewNopLogger(),
			want:   errors.New("serivce must not be nil"),
		},
		{
			name: "logger must not be nil",
			svc: AddFunc(func(ctx context.Context, name string, category int) error {
				return nil
			}),
			logger: nil,
			want:   errors.New("logger must not be nil"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := NewAddHandler(tt.svc, tt.logger)
			if h != nil {
				t.Fatal("Got an non-nil handler")
			}
			if !util.EqualError(err, tt.want) {
				t.Errorf("\nGot %#v\nwanted %#v", err, tt.want)
			}
		})
	}
}

func TestNewAddHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name   string
		method string
		body   string
		status int
		f      AddFunc
		want   string
	}{
		{
			name:   "normal",
			method: http.MethodPost,
			body:   `{"name":"sunny side up","category":0}`,
			status: http.StatusCreated,
			f: func(context.Context, string, int) error {
				return nil
			},
			want: "",
		},
		{
			name:   "method not allowed",
			method: http.MethodGet,
			body:   "",
			status: http.StatusMethodNotAllowed,
			f: func(context.Context, string, int) error {
				return nil
			},
			want: fmt.Sprintf(`{"error":%q}`, ErrDiallowedMethod),
		},
		{
			name:   "not json",
			method: http.MethodPost,
			body:   "not json",
			status: http.StatusBadRequest,
			f: func(context.Context, string, int) error {
				return nil
			},
			want: fmt.Sprintf(`{"error":%q}`, ErrInvalidRequestBody),
		},
		{
			name:   "name is missing",
			method: http.MethodPost,
			body:   `{"category":0}`,
			status: http.StatusBadRequest,
			f: func(context.Context, string, int) error {
				return nil
			},
			want: fmt.Sprintf(`{"error":%q}`, ErrName),
		},
		{
			name:   "name is empty",
			method: http.MethodPost,
			body:   `{"name":"","category":0}`,
			status: http.StatusBadRequest,
			f: func(context.Context, string, int) error {
				return nil
			},
			want: fmt.Sprintf(`{"error":%q}`, ErrName),
		},
		{
			name:   "category is missing",
			method: http.MethodPost,
			body:   `{"name":"sunny side up"}`,
			status: http.StatusBadRequest,
			f: func(context.Context, string, int) error {
				return nil
			},
			want: fmt.Sprintf(`{"error":%q}`, ErrCategory),
		},
		{
			name:   "category must not be negative",
			method: http.MethodPost,
			body:   `{"name":"sunny side up","category":-1}`,
			status: http.StatusBadRequest,
			f: func(context.Context, string, int) error {
				return nil
			},
			want: fmt.Sprintf(`{"error":%q}`, ErrCategory),
		},
		{
			name:   "category must not greater than 3",
			method: http.MethodPost,
			body:   `{"name":"sunny side up","category":4}`,
			status: http.StatusBadRequest,
			f: func(context.Context, string, int) error {
				return nil
			},
			want: fmt.Sprintf(`{"error":%q}`, ErrCategory),
		},
		{
			name:   "service returns an error",
			method: http.MethodPost,
			body:   `{"name":"sunny side up","category":0}`,
			status: http.StatusInternalServerError,
			f: func(context.Context, string, int) error {
				return errors.New("service returns an error")
			},
			want: `{"error":"service returns an error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := NewAddHandler(tt.f, log.NewLogfmtLogger(os.Stderr))
			if err != nil {
				t.Fatalf("Got an unexpected error: %#v", err)
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, "http://localhost/add", bytes.NewBuffer([]byte(tt.body)))
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
