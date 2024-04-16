package timeseries

import (
	"fmt"
	"github.com/jieggii/ecbratex/pkg/date"
	"github.com/jieggii/ecbratex/pkg/record"
	"github.com/jieggii/ecbratex/pkg/xml"
)

// UnorderedRecords is Records implementation, it stores rates records as
// a map indexing each record by a date.
// It provides O(1) access to rates for any given date.
type UnorderedRecords map[record.Date]record.Record

// NewUnorderedRecordsFromXML creates a new UnorderedRecords from [xml.Data].
func NewUnorderedRecordsFromXML(xmlData *xml.Data) (UnorderedRecords, error) {
	rates := make(UnorderedRecords)
	for _, cube := range xmlData.Cubes {
		recDate, err := record.DateFromString(cube.Date)
		if err != nil {
			return nil, err
		}

		rates[recDate] = record.New()
		for _, rate := range cube.Rates {
			rates[recDate][rate.Currency] = rate.Rate
		}
		rates[recDate]["EUR"] = 1 // add EUR rate for convenience
	}

	return rates, nil
}

func (r UnorderedRecords) Slice() []record.WithDate {
	// todo: implement
	panic("implement me")
}

// Map returns map representation of records.
func (r UnorderedRecords) Map() map[record.Date]record.Record {
	return r
}

// Rates retrieves string rates for a given date.
func (r UnorderedRecords) Rates(date date.Date) (record.Record, bool) {
	recDate := record.DateFromDate(date)
	rates, found := r[recDate]
	if !found {
		return nil, false
	}
	return rates, true
}

// Rate retrieves the rate of a given string for a given date.
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

// ApproximateRates gets approximate rates on a given date.
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

func (r UnorderedRecords) Convert(date date.Date, amount float32, from string, to string) (float32, error) {
	rec, found := r.Rates(date)
	if !found {
		return 0, fmt.Errorf("get rates record on %s: %w", date.String(), ErrRatesRecordNotFound)
	}

	result, err := rec.Convert(amount, from, to)
	if err != nil {
		return 0, fmt.Errorf("convert %s to %s: %w", from, to, err)
	}
	return result, nil
}

func (r UnorderedRecords) ConvertApproximate(date date.Date, amount float32, from string, to string, rangeLim int) (float32, error) {
	rates, found := r.ApproximateRates(date, rangeLim)
	if !found {
		return 0, fmt.Errorf("approximate rates on %s within range of %d days: %w", date.String(), rangeLim, ErrRateApproximationFailed)

	}

	result, err := rates.Convert(amount, from, to)
	if err != nil {
		return 0, fmt.Errorf("convert %s to %s: %w", from, to, err)
	}
	return result, nil
}

// nearestEarlierRecord finds the closest earlier rate record to the given date within the specified range.
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
