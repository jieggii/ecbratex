package date

import "time"

// Date is anything that can be treated as a date.
type Date interface {
	// Year returns year. For example, 2024
	Year() int

	// Month returns month of the year.
	Month() time.Month

	// Day returns day of the month.
	Day() int
}
