/*
This simple example illustrates usage of ConvertMinors function.
It is useful, because sometimes amounts are stored as an integer number in minor units:
for example, 16.25$ are stored as 1625 cents.
*/
package main

import (
	"fmt"
	"github.com/jieggii/ecbratex"
)

func main() {
	// fetch latest rates record:
	record, _ := ecbratex.FetchLatest()

	var eurMinors int64 = 1625 // 16.25 EUR
	usdMinors, _ := record.ConvertMinors(eurMinors, "EUR", "USD")
	fmt.Printf("%f EUR = %f USD\n", float32(eurMinors)/100, float32(usdMinors)/100)
}
