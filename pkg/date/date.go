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

	// String returns the string representation of the date in "YYYY-MM-DD" (ISO-8601) format.
	// For example, "2024-04-24".
	String() string
}
