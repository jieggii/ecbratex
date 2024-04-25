package record

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDate(t *testing.T) {
	const (
		year  Year  = 2020
		month Month = 99
		day   Day   = 31
	)

	date := NewDate(year, month, day)
	assert.Equal(t, year, date.year)
	assert.Equal(t, month, date.month)
	assert.Equal(t, day, date.day)
}

func TestDate_Year(t *testing.T) {
	const year Year = 1976
	date := NewDate(year, 11, 1)
	assert.Equal(t, int(year), date.Year())
}

func TestDate_Month(t *testing.T) {
	const month Month = 4
	date := NewDate(1999, month, 1)
	assert.Equal(t, time.Month(month), date.Month())
}

func TestDate_Day(t *testing.T) {
	const day Day = 16
	date := NewDate(1999, 2, day)
	assert.Equal(t, int(day), date.Day())
}

func TestDateFromTime(t *testing.T) {
	const (
		year  Year  = 2050
		month Month = 12
		day   Day   = 31
	)

	timeDate := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
	date := DateFromTime(timeDate)
	assert.Equal(t, NewDate(year, month, day), date)
}

func TestDateFromString(t *testing.T) {
	type output struct {
		date  Date
		isErr bool
	}
	type testCase struct {
		input  string
		output output
	}

	testCases := []testCase{
		{"2005-05-05", output{Date{2005, 5, 5}, false}},
		{"", output{ZeroDate, true}},
		{"invalid string", output{ZeroDate, true}},
		{"2004-2-15", output{ZeroDate, true}},
	}

	for i, c := range testCases {
		date, err := DateFromString(c.input)
		msg := fmt.Sprintf("case %d, input=%s", i, c.input)
		assert.Equal(t, c.output.date, date, msg)
		assert.Equal(t, c.output.isErr, err != nil, err)
	}
}

func TestDate_Time(t *testing.T) {
	const (
		year  Year  = 2050
		month Month = 12
		day   Day   = 31
	)
	date := NewDate(year, month, day)
	assert.Equal(t, time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC), date.Time())
}

func TestDate_After(t *testing.T) {
	date1 := NewDate(2000, 12, 31)
	date2 := NewDate(2020, 12, 1)
	assert.True(t, date2.After(date1))
}

func TestDate_Before(t *testing.T) {
	date1 := NewDate(2000, 12, 31)
	date2 := NewDate(2020, 12, 1)
	assert.True(t, date1.Before(date2))
}

func TestDate_AddDays(t *testing.T) {
	type in struct {
		date Date
		days int
	}
	type _testCase struct {
		in  in
		out Date
	}

	testCases := []_testCase{
		{in{date: NewDate(2000, 1, 1), days: 1}, NewDate(2000, 1, 2)},
		{in{date: NewDate(2000, 1, 31), days: 1}, NewDate(2000, 2, 1)},
		{in{date: NewDate(2000, 12, 31), days: 1}, NewDate(2001, 1, 1)},
		{in{date: NewDate(2000, 4, 20), days: 9}, NewDate(2000, 4, 29)},

		{in{date: NewDate(2000, 1, 1), days: -1}, NewDate(1999, 12, 31)},
		{in{date: NewDate(2000, 2, 2), days: -1}, NewDate(2000, 2, 1)},
		{in{date: NewDate(2010, 6, 5), days: -5}, NewDate(2010, 5, 31)},
	}

	for _, testCase := range testCases {
		result := testCase.in.date.AddDays(testCase.in.days)
		assert.Equalf(t, testCase.out, result, "input=%s", testCase.in)
	}
}

func TestDate_SubDays(t *testing.T) {
	type input struct {
		date  Date
		other Date
	}
	cases := map[input]int{
		{NewDate(2000, 1, 2), NewDate(2000, 1, 2)}:  0,
		{NewDate(2000, 1, 1), NewDate(2000, 1, 10)}: -9,
		{NewDate(2000, 1, 10), NewDate(2000, 1, 1)}: 9,
	}

	for in, out := range cases {
		assert.Equalf(t, out, in.date.SubDays(in.other), "input=%s", in)
	}
}

func TestDate_String(t *testing.T) {
	cases := map[Date]string{
		NewDate(1974, 12, 7): "1974-12-07",
		NewDate(2005, 1, 1):  "2005-01-01",
		NewDate(0, 0, 0):     "0000-00-00",
	}
	for input, output := range cases {
		dateString := input.String()
		assert.Equalf(t, output, dateString, "input=%s", input)
	}
}
