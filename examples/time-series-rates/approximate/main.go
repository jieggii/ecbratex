/*
	This code snipped illustrates how to fetch currency rates records and approximate
	rates in case if the record with with the desired date is not present in the dataset.
*/

package main

import (
	"fmt"
	"github.com/jieggii/ecbratex"
	"github.com/jieggii/ecbratex/pkg/record"
	"github.com/jieggii/ecbratex/pkg/timeseries"
)

func main() {
	// fetch all exchange rates records:
	records, _ := ecbratex.FetchTimeSeries(ecbratex.PeriodWhole)

	date := record.NewDate(2005, 1, 1)
	usdRate, _ := records.ApproximateRate(date, "USD", timeseries.DefaultRangeLim)
	fmt.Printf("approximated USD rate on %s: %f\n", date.String(), usdRate)

	// convert 500 USD to EUR on 2005-01-01:
	result, _ := records.ConvertApproximate(date, 500, "USD", "EUR", timeseries.DefaultRangeLim)
	fmt.Printf("500 USD was approximately equal to %f EUR on %s\n", result, date.String())
}
