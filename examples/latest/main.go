package main

import (
	"fmt"
	"github.com/jieggii/ecbratex"
)

func main() {
	// fetch latest rates record:
	record, _ := ecbratex.FetchLatest()

	// get latest USD exchange rate:
	usdRate, _ := record.Rate("USD")
	fmt.Printf("Latest USD rate (%s): %f\n", record.Date.String(), usdRate)

	// convert 500 USD to EUR:
	amount, _ := record.Convert(500, "USD", "EUR")
	fmt.Printf("500 USD = %f EUR\n", amount)
}
