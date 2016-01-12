package database

import (
	"strconv"
	"time"
)

// Measurement holds the information of a single occurence
type Measurement struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

// SetValueFromString converts a string to a float64 and set the field value of the structure
func (m *Measurement) SetValueFromString(str string) {
	fValue, _ := strconv.ParseFloat(str, 64)
	m.Value = fValue
}

// GetTimeAsString returns the time in a string format.
// Currently RFC3339
func (m *Measurement) GetTimeAsString() string {
	return m.Timestamp.Format(time.RFC3339)
}

// Database is uses to have a general interface to access different databases
type Database interface {
	WriteElement(element Measurement) error
}
