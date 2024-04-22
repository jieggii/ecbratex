package record

import (
	"errors"
	"fmt"
	"github.com/jieggii/ecbratex/pkg/xml"
)

// WithDate represents exchange rates on a specific date.
type WithDate struct {
	// Record contains exchange rates for currencies
	Record

	// Date represents the specific date for which exchange rates are recorded.
	Date Date
}

// NewWithDate creates a new WithDate and returns it.
func NewWithDate(record Record, date Date) WithDate {
	return WithDate{Record: record, Date: date}
}

// NewWithDateFromXMLData creates a new WithDate from xml.Data.
func NewWithDateFromXMLData(xmlData *xml.Data) (*WithDate, error) {
	if len(xmlData.Cubes) == 0 {
		return nil, errors.New("document does not contain any rate records")
	}
	cube := xmlData.Cubes[0]

	// represent the document as [WithDate]:
	date, err := DateFromString(cube.Date)
	if err != nil {
		return nil, fmt.Errorf("parse record date: %w", err)
	}

	record := New()
	for _, rate := range cube.Rates {
		record[rate.Currency] = rate.Rate
	}
	record["EUR"] = 1 // add EUR rate for convenience

	return &WithDate{
		Date:   date,
		Record: record,
	}, nil
}
