// Package entity has a data model and some constant variables.
package entity

import "time"

// Following variables represent types of meal.
const (
	Breakfast int = iota
	Lunch
	Dinner
	Others
)

// JournalID is a parimary key of the table.
type JournalID uint64

// Journal represents a record of meal.
type Journal struct {
	ID       JournalID `json:"id"`
	Name     string    `json:"name"`
	Cateogry int       `json:"category"`
	Created  time.Time `json:"created"`
}
