package timeseries

import (
	"github.com/jieggii/ecbratex/pkg/record"
	"github.com/jieggii/ecbratex/pkg/xml"
)

// OrderedUnorderedRecords (yes) is an implementation of the Records interface that combines
// the benefits of UnorderedRecords with the ability to store the order of rate records.
// It uses an underlying UnorderedRecords to store rates and perform calculations, while also
// keeping track of the order of all rate records for easy accessibility.
//
// Use OrderedUnorderedRecords only if memory constraints are not a concern, and you need to
// know the order of rate records. Otherwise, consider using UnorderedRecords or OrderedRecords.
type OrderedUnorderedRecords struct {
	// Dates of all rate records in chronological order.
	Dates []record.Date

	// UnorderedRecords is the underlying data structure to store rate records.
	UnorderedRecords
}

// NewOrderedUnorderedRecordsFromXML creates a new OrderedUnorderedRecords from [xml.Data].
func NewOrderedUnorderedRecordsFromXML(xmlData *xml.Data) (*OrderedUnorderedRecords, error) {
	rates := make(UnorderedRecords)
	dates := make([]record.Date, 0)

	for _, cube := range xmlData.Cubes {
		date, err := record.DateFromString(cube.Date)
		if err != nil {
			return nil, err
		}
		dates = append(dates, date)

		rates[date] = record.New()
		for _, rate := range cube.Rates {
			rates[date][rate.Currency] = rate.Rate
		}
		rates[date]["EUR"] = 1 // add EUR rate for convenience
	}

	return &OrderedUnorderedRecords{
		Dates:            dates,
		UnorderedRecords: rates,
	}, nil
}

func (r OrderedUnorderedRecords) Slice() []record.WithDate {
	records := make([]record.WithDate, 0)
	for _, date := range r.Dates {
		rec := r.UnorderedRecords[date]
		records = append(records, record.NewWithDate(rec, date))
	}
	return records
}
