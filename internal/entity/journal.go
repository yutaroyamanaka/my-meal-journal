// Package entity has a data model and some constant variables.
package entity

import "time"

// Following variables represent types of meal.
const (
	Breakfast uint8 = iota
	Lunch
	Dinner
	Others
)

// JournalID is a parimary key of the table.
type JournalID uint64

// CategoryName is a string type for meal category.
type CategoryName string

// Following variables represent string expressions of meal category.
const (
	CategoryBreakfast   CategoryName = "Breakfast"
	CategoryLunchName   CategoryName = "Lunch"
	CategoryDinnerName  CategoryName = "Dinner"
	CateogoryOthersName CategoryName = "Others"
)

// Journal represents a record of meal.
type Journal struct {
	ID       JournalID `json:"id"`
	Name     string    `json:"name"`
	Cateogry uint8     `json:"category"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}
