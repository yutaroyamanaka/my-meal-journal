package service

import (
	"context"
	"errors"
	"os"
	"testing"

	log "github.com/go-kit/log"
	"github.com/yutaroyamanaka/my-meal-journal/internal/entity"
	"github.com/yutaroyamanaka/my-meal-journal/internal/util"
)

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
			name: "logger must not be nil",
			repo: AdderJournalFunc(func(ctx context.Context, j *entity.Journal) error {
				return nil
			}),
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
			if !util.EqualError(err, tt.want) {
				t.Errorf("\nGot %#v\nwanted %#v", err, tt.want)
			}
		})
	}
}

func TestAddService_Add(t *testing.T) {
	type args struct {
		name     string
		category uint8
	}
	tests := []struct {
		name string
		repo JournalAdder
		args args
		want error
	}{
		{
			name: "normal",
			repo: AdderJournalFunc(func(ctx context.Context, j *entity.Journal) error {
				return nil
			}),
			args: args{"sunny side up", 0},
			want: nil,
		},
		{
			name: "JournalAdder returns an error",
			repo: AdderJournalFunc(func(ctx context.Context, j *entity.Journal) error {
				return errors.New("unexpected error is occurred")
			}),
			args: args{"sunny side up", 0},
			want: errors.New("failed to register your meal's information"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := NewAddService(tt.repo, log.NewLogfmtLogger(os.Stderr))
			if err != nil {
				t.Fatalf("Got an unexpected error: %#v", err)
			}
			err = svc.Add(context.Background(), tt.args.name, tt.args.category)
			if !util.EqualError(err, tt.want) {
				t.Errorf("\nGot %#v\nwanted %#v", err, tt.want)
			}
		})
	}
}
