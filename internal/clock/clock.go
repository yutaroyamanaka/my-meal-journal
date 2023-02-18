// Package clock has RealTimeClocker and StubTimeClocker.
package clock

import "time"

// Clocker returns time.
type Clocker interface {
	Now() time.Time
}

// RealTimeClocker returns the current time.
type RealTimeClocker struct{}

// Now calls just time.Now().
func (c *RealTimeClocker) Now() time.Time {
	return time.Now()
}

// StubTimeClocker returns the certain time.
type StubTimeClocker struct{}

// Now calls the specific date.
func (c *StubTimeClocker) Now() time.Time {
	return time.Date(1997, 8, 22, 12, 10, 20, 0, time.UTC)
}
