// Package service has business logic of the application.
package service

import (
	"context"
	"errors"
	"fmt"

	log "github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/yutaroyamanaka/my-meal-journal/internal/entity"
)

// JournalAdder deals with business logic about new journal registration.
type JournalAdder interface {
	AddJournal(context.Context, string, int) (*entity.Journal, error)
}

// AdderJournalFunc is a stub function for mocking JournalAdder interface.
type AdderJournalFunc func(context.Context, string, int) (*entity.Journal, error)

// AddJournal deals with business logic about new journal registration.
func (f AdderJournalFunc) AddJournal(ctx context.Context, name string, category int) (*entity.Journal, error) {
	return f(ctx, name, category)
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
func (s *AddService) Add(ctx context.Context, name string, category int) (*entity.Journal, error) {
	if name == "" {
		return nil, errors.New("name must not be empty")
	}
	if category < 0 || category > entity.Others {
		return nil, fmt.Errorf("category must be between %d(%q) and %d(%q)",
			entity.Breakfast, entity.CategoryBreakfast, entity.Others, entity.CateogoryOthersName)
	}
	journal, err := s.repo.AddJournal(ctx, name, category)
	if err != nil {
		level.Error(s.logger).Log("msg", "failed to add a journal in JournalAdder", "err", err)
		return nil, errors.New("failed to register your meal's information")
	}
	return journal, err
}
