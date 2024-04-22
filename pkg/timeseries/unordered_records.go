package timeseries

import (
	"github.com/jieggii/ecbratex/pkg/date"
	"github.com/jieggii/ecbratex/pkg/record"
	"github.com/jieggii/ecbratex/pkg/xml"
	"sort"
)

// UnorderedRecords is an implementation of the Records interface.
// It stores records as a map indexing each record by a date so each record is accessible for O(1).
type UnorderedRecords map[record.Date]record.Record

// NewUnorderedRecordsFromXML creates a new UnorderedRecords from [xml.Data].
func NewUnorderedRecordsFromXML(xmlData *xml.Data) (UnorderedRecords, error) {
	records := make(UnorderedRecords)
	for _, cube := range xmlData.Cubes {
		recDate, err := record.DateFromString(cube.Date)
		if err != nil {
			return nil, err
		}

		records[recDate] = record.New()
		for _, rate := range cube.Rates {
			records[recDate][rate.Currency] = rate.Rate
		}
		records[recDate]["EUR"] = 1 // add EUR rate for convenience
	}

	return records, nil
}

// Slice returns the underlying slice containing all records in anti-chronological order.
// Operates on O(n log n) time complexity.
func (r UnorderedRecords) Slice() []record.WithDate {
	recordsCount := len(r)

	// create and fill dates slice:
	dates := make([]record.Date, 0, recordsCount)
	for d := range r {
		dates = append(dates, d)
	}

	// sort dates slice:
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].After(dates[j])
	})

	// create and fill records:
	records := make([]record.WithDate, 0, recordsCount)
	for _, d := range dates {
		records = append(records, record.NewWithDate(r[d], d))
	}

	return records
}

// Map creates and returns map of all records indexed by date.
// Operates on O(1) time complexity.
func (r UnorderedRecords) Map() map[record.Date]record.Record {
	return r
}

// Rates returns rates on the given date.
// Operates on O(1) time complexity.
func (r UnorderedRecords) Rates(date date.Date) (record.Record, bool) {
	recDate := record.DateFromDate(date)
	rates, found := r[recDate]
	if !found {
		return nil, false
	}
	return rates, true
}

// Rate returns rate of the given currency on the given date.
// Operates on O(1) time complexity.
func (r UnorderedRecords) Rate(date date.Date, currency string) (float32, bool) {
	rates, found := r.Rates(date)
	if !found {
		return 0, false
	}

	rate, found := rates[currency]
	if !found {
		return 0, false
	}
	return rate, true
}

// ApproximateRates approximates and returns approximated rates on the given date.
// Operates on O(rangeLim) time complexity.
func (r UnorderedRecords) ApproximateRates(date date.Date, rangeLim int) (record.Record, bool) {
	recDate := record.DateFromDate(date)

	// find the closest earlier rate records:
	earlierRec, earlierRecFound := r.nearestEarlierRecord(recDate, rangeLim)
	laterRec, laterRecFound := r.nearestLaterRecord(recDate, rangeLim)

	// if neither earlier nor later rate records were found:
	if !earlierRecFound && !laterRecFound {
		return nil, false
	}

	// return later record if earlier record was not found:
	if !earlierRecFound {
		return laterRec, true
	}

	// return earlier record if later record was not found:
	if !laterRecFound {
		return earlierRec, true
	}

	// approximate the closest earlier and later records if both were found
	// (calculating average where possible or using later or earlier rates):
	ratesRecord := record.New()
	for currency, earlierRate := range earlierRec {
		laterRate, found := laterRec.Rate(currency)
		if !found {
			ratesRecord[currency] = earlierRate
			continue
		}
		ratesRecord[currency] = (earlierRate + laterRate) / 2
	}

	// fill ratesRecord with possible missing rates from later records:
	for currency, laterRate := range laterRec {
		_, found := ratesRecord[currency]
		if !found {
			ratesRecord[currency] = laterRate
		}
	}

	return ratesRecord, true
}

// ApproximateRate approximates and returns approximated rate of the given currency on the given date.
// Operates on O(rangeLim) time complexity.
func (r UnorderedRecords) ApproximateRate(date date.Date, currency string, rangeLim int) (float32, bool) {
	recDate := record.DateFromDate(date)
	earlierRec, earlierRecFound := r.nearestEarlierRecord(recDate, rangeLim)
	laterRec, laterRecFound := r.nearestLaterRecord(recDate, rangeLim)

	if !earlierRecFound && !laterRecFound {
		return 0, false
	}

	if !earlierRecFound {
		return laterRec.Rate(currency)
	}

	if !laterRecFound {
		return earlierRec.Rate(currency)
	}

	earlierRate, foundEarlierRate := earlierRec.Rate(currency)
	laterRate, foundLaterRate := laterRec.Rate(currency)

	if !foundEarlierRate && !foundLaterRate {
		return 0, false
	}

	if !foundEarlierRate {
		return laterRate, true
	}

	if !foundLaterRate {
		return earlierRate, true
	}

	return (earlierRate + laterRate) / 2, true
}

// Convert converts amount from one currency to another on the given date.
// Operates on O(1) time complexity.
func (r UnorderedRecords) Convert(date date.Date, amount float32, from string, to string) (float32, error) {
	rec, found := r.Rates(date)
	if !found {
		return 0, ErrRatesRecordNotFound
	}

	result, err := rec.Convert(amount, from, to)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// ConvertApproximate converts amount from one currency to another on the given date,
// using approximated rates within rangeLim days.
// Operates on O(rangeLim) time complexity.
func (r UnorderedRecords) ConvertApproximate(date date.Date, amount float32, from string, to string, rangeLim int) (float32, error) {
	rates, found := r.ApproximateRates(date, rangeLim)
	if !found {
		return 0, ErrRateApproximationFailed

	}

	result, err := rates.Convert(amount, from, to)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// ConvertMinors converts amount in minor units from one currency to another on the given date.
// Operates on O(1) time complexity.
func (r UnorderedRecords) ConvertMinors(date date.Date, amount int, from string, to string) (int, error) {
	rec, found := r.Rates(date)
	if !found {
		return 0, ErrRatesRecordNotFound
	}

	result, err := rec.ConvertMinors(amount, from, to)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// ConvertMinorsApproximate converts amount in minor units from one currency to another on the given date,
// using approximated rates within rangeLim days.
// Operates on O(rangeLim) time complexity.
func (r UnorderedRecords) ConvertMinorsApproximate(date date.Date, amount int, from string, to string, rangeLim int) (int, error) {
	rates, found := r.ApproximateRates(date, rangeLim)
	if !found {
		return 0, ErrRateApproximationFailed

	}

	result, err := rates.ConvertMinors(amount, from, to)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// nearestEarlierRecord finds the closest earlier rate record to the given date within the specified range.
// Operates on O(rangeLim) time complexity.
func (r UnorderedRecords) nearestEarlierRecord(recDate record.Date, rangeLim int) (record.Record, bool) {
	for range rangeLim {
		recDate = recDate.AddDays(-1)
		rec, found := r.Rates(recDate)
		if found {
			return rec, true
		}
	}

	return nil, false
}

// nearestLaterRecord finds the closest later rate record to the given date within the specified range.
// Operates on O(rangeLim) time complexity.
func (r UnorderedRecords) nearestLaterRecord(recDate record.Date, rangeLim int) (record.Record, bool) {
	for range rangeLim {
		recDate = recDate.AddDays(1)
		rec, found := r.Rates(recDate)
		if found {
			return rec, true
		}
	}
	return nil, false
}
