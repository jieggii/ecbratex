/*
	This is a simple example illustrating how to retrieve and handle historical time series rates records.
*/

package main

import (
	"fmt"
	"github.com/jieggii/ecbratex"
	"github.com/jieggii/ecbratex/pkg/record"
)

func main() {
	// fetch all exchange rates records for the whole period:
	records, _ := ecbratex.FetchTimeSeries(ecbratex.PeriodWhole)

	date := record.NewDate(2000, 2, 14)
	//date := time.Date(2000, 2, 14, 0, 0, 0, 0, time.UTC) // you can also use time module instead

	// get USD rate on 2000-02-14:
	usdRate, _ := records.Rate(date, "USD")
	fmt.Printf("USD rate on %s: %f\n", date.String(), usdRate)

	// convert 999 USD to EUR on 2000-02-14:
	result, _ := records.Convert(date, 999, "USD", "EUR")
	fmt.Printf("999 USD = %f EUR on %s\n", result, date.String())
}
