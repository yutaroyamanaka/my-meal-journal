// Package service has business logic of the application.
package service

import (
	"context"
	"errors"

	log "github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/yutaroyamanaka/my-meal-journal/internal/entity"
)

// JournalAdder deals with business logic about new journal registration.
type JournalAdder interface {
	AddJournal(context.Context, *entity.Journal) error
}

// AdderJournalFunc is a stub function for mocking JournalAdder interface.
type AdderJournalFunc func(context.Context, *entity.Journal) error

// AddJournal deals with business logic about new journal registration.
func (f AdderJournalFunc) AddJournal(ctx context.Context, j *entity.Journal) error {
	return f(ctx, j)
}

// AddService has JournalAdder interface and log.Logger as fields.
type AddService struct {
	repo   JournalAdder
	logger log.Logger
}

// NewAddService returns AddService and error.
// This function doesn't accept "nil" as JournalAdder and log.Logger.
func NewAddService(repo JournalAdder, logger log.Logger) (*AddService, error) {
	if repo == nil {
		return nil, errors.New("repo must not be nil")
	}
	if logger == nil {
		return nil, errors.New("logger must not be nil")
	}
	return &AddService{repo, logger}, nil
}

// Add validates name and category, and calls AddJournal method of JournalAdder interface.
func (s *AddService) Add(ctx context.Context, name string, category uint8) error {
	j := &entity.Journal{
		Name:     name,
		Cateogry: category,
	}
	err := s.repo.AddJournal(ctx, j)
	if err != nil {
		level.Error(s.logger).Log("msg", "failed to add a journal in JournalAdder", "err", err)
		return errors.New("failed to register your meal's information")
	}
	return nil
}
