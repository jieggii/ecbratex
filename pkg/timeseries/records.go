package timeseries

import (
	"errors"
	"github.com/jieggii/ecbratex/pkg/date"
	"github.com/jieggii/ecbratex/pkg/record"
)

var (
	ErrRatesRecordNotFound     = errors.New("exchange rates record was not found on the given date")
	ErrRateApproximationFailed = errors.New("approximation of one or multiple exchange rates did not succeed within the given date and range limit")
)

// DefaultRangeLim is a sane default for rangeLim argument passed to some functions of Records.
// At the moment of 2024-04-09 (when I write this comment), the largest time gap between two records is 5 days,
// meaning, that rangeLim equal to 5 will work for any date within [record.MinDate] - *now* interval.
// However, rangeLim argument does not affect performance, so, I think, that 100 is a good value in such situation.
const DefaultRangeLim = 100

// Records describes a wrapper around historical rates data.
// It provides useful functions for retrieving or approximating rates and
// functions for converting amounts from one currency to another.
type Records interface {
	// Slice returns slice of all exchange rates records in the anti-chronological order.
	Slice() []record.WithDate

	// Map returns map containing all exchange rates records indexed by their date.
	Map() map[record.Date]record.Record

	// Rates retrieves currency rates for the given date.
	// It returns the rates record for the specified date and a boolean indicating whether the record was found.
	Rates(date date.Date) (record.Record, bool)

	// Rate retrieves the rate of the given currency for the given date.
	// It returns the rate for the specified currency and a boolean indicating whether the rate was found.
	Rate(date date.Date, currency string) (float32, bool)

	// ApproximateRates calculates approximate rates for the given date within a specified days range.
	// It searches for the nearest earlier and later rate records within the given range (rangeLim).
	// If neither earlier nor later rate records are found within the range limit, it returns false.
	// If only one of either earlier or later rate records is found, it returns the rates of that record.
	// If both earlier and later rate records are found, it approximates rates by interpolating between
	// the rates of the closest earlier and later records, or using the rates.
	// Is useful when there is no rates record on the desired date.
	ApproximateRates(date date.Date, rangeLim int) (record.Record, bool)

	// ApproximateRate calculates an approximate rate for the given currency on a given date within a specified days range.
	// It searches for the closest earlier and later rate records within the given range (rangeLim).
	// If neither earlier nor later rate records are found within the range limit, it returns false.
	// If only one of either earlier or later rate records is found, it returns the rate of that record.
	// If both earlier and later rate records are found, it approximates the rate by interpolating between
	// the rates of the closest earlier and later records, or using the rates if interpolation is not possible.
	// Is useful when there is no rates record on the desired date.
	ApproximateRate(date date.Date, currency string, rangeLim int) (float32, bool)

	// Convert converts the specified amount from one currency to another on the given date.
	// It returns the converted amount as a float32 or an error if conversion fails.
	Convert(date date.Date, amount float32, from string, to string) (float32, error)

	// ConvertApproximate converts the specified amount of currency from one currency to another on the given date
	// using approximate rates within a specified days determined by rangeLim.
	// It returns the converted amount as a float32 or an error if conversion fails.
	ConvertApproximate(date date.Date, amount float32, from string, to string, rangeLim int) (float32, error)

	// ConvertMinors converts the specified amount of currency in minor units from one currency to another on the given date.
	// It returns the converted amount as an int or an error if conversion fails.
	ConvertMinors(date date.Date, amount int, from string, to string) (int, error)

	// ConvertMinorsApproximate converts the specified amount of currency in minor units from one currency to another on the given date
	// using approximate rates within a specified days determined by rangeLim.
	// It returns the converted amount as an int or an error if conversion fails.
	ConvertMinorsApproximate(date date.Date, amount int, from string, to string, rangeLim int) (int, error)
}
