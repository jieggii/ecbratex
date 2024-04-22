# ecbratex
Go package for fetching and manipulating exchange rates from the [European Central Bank](https://www.ecb.europa.eu/).

**ecbratex** is inspired by [openprovider/ecbrates](https://github.com/openprovider/ecbrates), but has some additional features.

## Project status
Work in progress, the module is almost done! No API changes expected.

## Simple examples
### Fetch latest rates

```go
package main

import (
    "fmt"
    "github.com/jieggii/ecbratex"
)

func main() {
    record, _ := ecbratex.FetchLatest()
  
    // get latest USD exchange rate:
    usdRate, _ := record.Rate("USD")
    fmt.Printf("Latest USD rate (%s): %f\n", record.Date.String(), usdRate)
  
    // convert 500 USD to EUR:
    amount, _ := record.Convert(500, "USD", "EUR")
    fmt.Printf("500 USD = %f EUR\n", amount)
}
```

### Fetch time series rates
```go
package main

import (
    "fmt"
    "github.com/jieggii/ecbratex"
    "github.com/jieggii/ecbratex/pkg/record"
)

func main() {
    // fetch all exchange rates records:
    records, _ := ecbratex.FetchTimeSeries(ecbratex.PeriodWhole)
  
    // get USD rate on 2000-02-14:
    date := record.NewDate(2000, 2, 14)
    //date := time.Date(2000, 2, 14, 0, 0, 0, 0, time.UTC) // you can also use time module instead
  
    usdRate, _ := records.Rate(date, "USD")
    fmt.Printf("USD rate on %s: %f\n", date.String(), usdRate)
  
    // convert 999 USD to EUR on 2000-02-14:
    result, _ := records.Convert(date, 999, "USD", "EUR")
    fmt.Printf("999 USD = %f EUR on %s\n", result, date.String())
}

```
