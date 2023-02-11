package service

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	log "github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
	"github.com/yutaroyamanaka/my-meal-journal/internal/entity"
)

func newAdderJournalFunc(hasError bool) AdderJournalFunc {
	if hasError {
		return AdderJournalFunc(func(ctx context.Context, name string, category int) (*entity.Journal, error) {
			return nil, errors.New("unexpected error is occurred")
		})
	}
	return AdderJournalFunc(func(ctx context.Context, name string, category int) (*entity.Journal, error) {
		return &entity.Journal{ID: 0, Name: "sunny side up", Cateogry: 0, Created: time.Date(2023, 2, 5, 16, 27, 56, 0, time.UTC)}, nil
	})
}

func TestNewAddService_error(t *testing.T) {
	tests := []struct {
		name   string
		repo   JournalAdder
		logger log.Logger
		want   error
	}{
		{
			name:   "JournalAdder must not be nil",
			repo:   nil,
			logger: log.NewNopLogger(),
			want:   errors.New("repo must not be nil"),
		},
		{
			name:   "logger must not be nil",
			repo:   newAdderJournalFunc(false),
			logger: nil,
			want:   errors.New("logger must not be nil"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := NewAddService(tt.repo, tt.logger)
			if svc != nil {
				t.Fatal("Got an non-nil AddService")
			}
			if !equalError(err, tt.want) {
				t.Errorf("\nGot %#v\nwanted %#v", err, tt.want)
			}
		})
	}
}

func TestAddService_Add(t *testing.T) {
	type args struct {
		name     string
		category int
	}
	tests := []struct {
		name    string
		repo    JournalAdder
		args    args
		want    *entity.Journal
		wantErr error
	}{
		{
			name:    "normal",
			repo:    newAdderJournalFunc(false),
			args:    args{"sunny side up", 0},
			want:    &entity.Journal{ID: 0, Name: "sunny side up", Cateogry: 0, Created: time.Date(2023, 2, 5, 16, 27, 56, 0, time.UTC)},
			wantErr: nil,
		},
		{
			name:    "JournalAdder returns an error",
			repo:    newAdderJournalFunc(true),
			args:    args{"sunny side up", 0},
			want:    nil,
			wantErr: errors.New("failed to register your meal's information"),
		},
		{
			name:    "empty name",
			repo:    newAdderJournalFunc(false),
			args:    args{"", 0},
			want:    nil,
			wantErr: errors.New("name must not be empty"),
		},
		{
			name:    "invalid category (lower than 0)",
			repo:    newAdderJournalFunc(false),
			args:    args{"sunny side up", -1},
			want:    nil,
			wantErr: errors.New(`category must be between 0("Breakfast") and 3("Others")`),
		},
		{
			name:    "invalid category (greater than 3)",
			repo:    newAdderJournalFunc(false),
			args:    args{"sunny side up", 4},
			want:    nil,
			wantErr: errors.New(`category must be between 0("Breakfast") and 3("Others")`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := NewAddService(tt.repo, log.NewLogfmtLogger(os.Stderr))
			if err != nil {
				t.Fatalf("Got an unexpected error: %#v", err)
			}
			got, err := svc.Add(context.Background(), tt.args.name, tt.args.category)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Got an unexcpeted state: %s", diff)
			}
			if !equalError(err, tt.wantErr) {
				t.Errorf("\nGot %#v\nwanted %#v", err, tt.wantErr)
			}
		})
	}
}

func equalError(e1, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if (e1 != nil && e2 != nil) && e1.Error() == e2.Error() {
		return true
	}
	return false
}
