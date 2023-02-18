// Package store has logics for storing persistent data.
package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yutaroyamanaka/my-meal-journal/internal/clock"
	"github.com/yutaroyamanaka/my-meal-journal/internal/entity"
)

// Repository has logics that deal with queries.
type Repository struct {
	clocker clock.Clocker
	*sql.DB
}

// NewRepository returns *Repository when both clock.Clocker and *sql.DB are not nil.
func NewRepository(c clock.Clocker, db *sql.DB) (*Repository, error) {
	if c == nil {
		return nil, errors.New("clocker must not be nil")
	}
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &Repository{c, db}, nil
}

// AddJournal inserts a journal to the database and returns an error when it failed to insert.
func (r *Repository) AddJournal(ctx context.Context, j *entity.Journal) error {
	stmt, err := r.Prepare("INSERT INTO journal (name, category, created, updated) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	t := r.clocker.Now()
	result, err := stmt.ExecContext(ctx, j.Name, j.Cateogry, t, t)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	j.ID = entity.JournalID(id)
	return nil
}
