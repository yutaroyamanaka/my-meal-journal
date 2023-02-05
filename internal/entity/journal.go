package entity

import "time"

type JournalID int64

const (
	Breakfast int = iota
	Lunch
	Dinner
	Others
)

type Journal struct {
	ID       JournalID `json:"id"`
	Name     string    `json:"name"`
	Cateogry int       `json:"category"`
	Created  time.Time `json:"created"`
}
