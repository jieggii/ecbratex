package record

import (
	"fmt"
	"github.com/jieggii/ecbratex/pkg/date"
	"time"
)

var (
	// ZeroDate represents a record date with all components set to zero.
	ZeroDate = NewDate(0, 0, 0)

	// MinDate represents the minimum allowed record date, set to January 4, 1999.
	MinDate = NewDate(1999, 1, 4)
)

type (
	Year  uint16
	Month uint8
	Day   uint8
)

// Date is an implementation of the [date.Date] interface containing useful
// methods for manipulating record dates.
// Stores year, month and day in minimalistic data types to save memory.
type Date struct {
	year  Year
	month Month
	day   Day
}

// NewDate creates a new Date with the given year, month, and day components.
func NewDate(year Year, month Month, day Day) Date {
	return Date{
		year:  year,
		month: month,
		day:   day,
	}
}

// DateFromDate creates a new Date from anything that implements data.Date interface.
func DateFromDate(d date.Date) Date {
	return NewDate(
		Year(d.Year()),
		Month(d.Month()),
		Day(d.Day()),
	)
}

// DateFromTime creates a new Date from time.Time.
func DateFromTime(t time.Time) Date {
	return NewDate(
		Year(t.Year()),
		Month(t.Month()),
		Day(t.Day()),
	)
}

// DateFromString creates a new Date from a string representing a date in "YYYY-MM-DD" format.
func DateFromString(date string) (Date, error) {
	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return ZeroDate, err
	}
	return DateFromDate(t), nil
}

// Year returns year.
func (d Date) Year() int {
	return int(d.year)
}

// Month returns month of the year.
func (d Date) Month() time.Month {
	return time.Month(d.month)
}

// Day returns day of the month.
func (d Date) Day() int {
	return int(d.day)
}

// String returns string representation of the Date in "YYYY-MM-DD" format.
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}

// Time returns a time.Time representation of the Date.
func (d Date) Time() time.Time {
	return time.Date(int(d.year), time.Month(d.month), int(d.day), 0, 0, 0, 0, time.UTC)
}

// After returns true if the date is later than another date.
func (d Date) After(other Date) bool {
	t := d.Time()
	o := other.Time()
	return t.After(o)
}

// Before returns true if the date is earlier than another date.
func (d Date) Before(other Date) bool {
	t := d.Time()
	o := other.Time()
	return t.Before(o)
}

// AddDays adds the given number days to the date and returns a new date.
func (d Date) AddDays(days int) Date {
	t := d.Time()
	result := t.AddDate(0, 0, days)
	return DateFromTime(result)
}

// SubDays returns the number of days between the current date and another date.
func (d Date) SubDays(other Date) int {
	t := d.Time()
	o := other.Time()
	return int(t.Sub(o).Hours() / 24)
}
