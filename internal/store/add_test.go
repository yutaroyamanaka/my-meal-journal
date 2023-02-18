package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/caarlos0/env"
	"github.com/google/go-cmp/cmp"
	"github.com/yutaroyamanaka/my-meal-journal/internal/clock"
	"github.com/yutaroyamanaka/my-meal-journal/internal/entity"
	"github.com/yutaroyamanaka/my-meal-journal/internal/util"
)

func TestNewRepository(t *testing.T) {
	tests := []struct {
		name    string
		clocker clock.Clocker
		db      *sql.DB
		want    error
	}{
		{
			name:    "Clocker must not be nil",
			clocker: nil,
			db:      &sql.DB{},
			want:    errors.New("clocker must not be nil"),
		},
		{
			name:    "db must not be nil",
			clocker: &clock.RealTimeClocker{},
			db:      nil,
			want:    errors.New("db must not be nil"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRepository(tt.clocker, tt.db)
			if r != nil {
				t.Fatal("Got a non-nil Repository")
			}
			if !util.EqualError(err, tt.want) {
				t.Errorf("\nGot %#v\nwanted %#v", err, tt.want)
			}
		})
	}
}

func TestRepository_AddNewJournal(t *testing.T) {
	// set some envs
	t.Setenv("DB_USER", "test")
	t.Setenv("DB_PASSWORD", "test")
	t.Setenv("DB_NAME", "test")
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "3306")

	c := Config{}
	if err := env.Parse(&c); err != nil {
		t.Fatalf("Got an unexpected error: %#v", err)
	}
	db, cleanup, err := Open(&c)
	if err != nil {
		t.Fatalf("Got an unexpected error: %#v", err)
	}
	defer cleanup()

	r, err := NewRepository(&clock.StubTimeClocker{}, db)
	if err != nil {
		t.Fatalf("Got an unexpected error: %#v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Got an unexpected error: %#v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })

	// delete all records before running tests
	_, err = db.ExecContext(ctx, "DELETE FROM journal")
	if err != nil {
		t.Fatalf("Got an unexpected error: %#v", err)
	}

	// insert a new journal
	input := &entity.Journal{Name: "sunny side up", Cateogry: 0}
	err = r.AddJournal(ctx, input)
	if err != nil {
		t.Fatalf("Got an unexpected error: %#v", err)
	}

	// get the first record
	var (
		id       uint64
		name     string
		category uint8
	)
	err = db.QueryRowContext(ctx, "SELECT id, name, category FROM journal limit 1").Scan(&id, &name, &category)
	if err != nil {
		t.Fatalf("Got an unexpected error: %#v", err)
	}
	if diff := cmp.Diff(entity.JournalID(id), input.ID); diff != "" {
		t.Errorf("Got an unexcpeted id: %s", diff)
	}
	if diff := cmp.Diff(name, input.Name); diff != "" {
		t.Errorf("Got an unexcpeted name: %s", diff)
	}
	if diff := cmp.Diff(category, input.Cateogry); diff != "" {
		t.Errorf("Got an unexcpeted category: %s", diff)
	}

	// cancel context intentionally here
	cancel()
	// try to insert another journal
	err = r.AddJournal(ctx, &entity.Journal{Name: "pizza", Cateogry: 1})
	if err == nil {
		t.Error("Got nil, but want error")
	}
}
