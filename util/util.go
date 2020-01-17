package util

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// NewID generates a random 16 byte UUID
func NewID() string {
	var buf [16]byte
	rand.Read(buf[:])
	return hex.EncodeToString(buf[:])
}

// DateRange specifies a time range
type DateRange struct {
	Start time.Time
	End   time.Time
}

// Contains returns a bool value based on if the time specified is contained within the DateRange
func (dt DateRange) Contains(t time.Time) bool {
	return !t.Before(dt.Start) && !t.After(dt.End)
}
