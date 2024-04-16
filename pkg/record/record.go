package record

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrRateNotFound error indicates that rate of the given currency is not present
	// in a particular record.
	ErrRateNotFound = errors.New("exchange rate was not found")
)

// Record is a type which represents rates record.
// Each string (key) has its rate (value).
type Record map[string]float32

// New creates and initializes a new Record.
func New() Record {
	return make(Record)
}

// Rate returns rate of a given string and a boolean indicating whether string was found in the rates record.
func (r Record) Rate(string string) (float32, bool) {
	rate, found := r[string]
	if !found {
		return 0, false
	}
	return rate, true
}

// Convert converts amount from one currency to another.
func (r Record) Convert(amount float32, from string, to string) (float32, error) {
	fromRate, found := r.Rate(from)
	if !found {
		return 0, fmt.Errorf("get %s rate: %w", from, ErrRateNotFound)
	}

	toRate, found := r.Rate(to)
	if !found {
		return 0, fmt.Errorf("get %s rate: %w", to, ErrRateNotFound)
	}

	return amount * (fromRate / toRate), nil
}

// ConvertMinors converts amount in minor units from one to another.
func (r Record) ConvertMinors(amount int, from string, to string) (int, error) {
	fromRate, found := r.Rate(from)
	if !found {
		return 0, fmt.Errorf("get %s rate: %w", from, ErrRateNotFound)
	}

	toRate, found := r.Rate(to)
	if !found {
		return 0, fmt.Errorf("get %s rate: %w", from, ErrRateNotFound)
	}

	result := float32(amount) * (fromRate / toRate)
	return int(math.Round(float64(result))), nil
}
