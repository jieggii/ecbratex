package timeseries

import (
	"fmt"
	"github.com/jieggii/ecbratex/pkg/date"
	"github.com/jieggii/ecbratex/pkg/record"
	"github.com/jieggii/ecbratex/pkg/xml"
)

// OrderedRecords is an implementation of the Records interface.
// It stores records in a sorted slice, so the order of records is easily accessible.
type OrderedRecords []record.WithDate

// NewOrderedRecords creates new empty OrderedRecords.
func NewOrderedRecords() OrderedRecords {
	return make(OrderedRecords, 0)
}

// NewOrderedRecordsFromXML creates new OrderedRecords from [xml.Data].
func NewOrderedRecordsFromXML(xmlData *xml.Data) (OrderedRecords, error) {
	records := NewOrderedRecords()

	for _, cube := range xmlData.Cubes {
		rec := record.New()
		for _, rate := range cube.Rates {
			rec[rate.Currency] = rate.Rate
		}
		rec["EUR"] = 1 // add EUR rate for convenience

		recDate, err := record.DateFromString(cube.Date)
		if err != nil {
			return nil, err
		}
		records = append(records, record.NewWithDate(rec, recDate))
	}
	return records, nil
}

// Slice returns the underlying slice containing all records in anti-chronological order.
// Operates on O(1) time complexity.
func (r OrderedRecords) Slice() []record.WithDate {
	return r
}

// Map creates and returns map of all records indexed by date.
// Operates on O(n) time complexity.
func (r OrderedRecords) Map() map[record.Date]record.Record {
	result := make(map[record.Date]record.Record)
	for _, rec := range r {
		result[rec.Date] = rec.Record
	}
	return result
}

// Rates returns rates on the given date.
// Operates on O(n) time complexity.
func (r OrderedRecords) Rates(date date.Date) (record.Record, bool) {
	recDate := record.DateFromDate(date)
	for _, rec := range r {
		if rec.Date == recDate {
			return rec.Record, true
		}
	}
	return nil, false
}

// Rate returns rate of the given currency on the given date.
// Operates on O(n) time complexity.
func (r OrderedRecords) Rate(date date.Date, string string) (float32, bool) {
	rec, found := r.Rates(date)
	if !found {
		return 0, false
	}

	rate, found := rec[string]
	if !found {
		return 0, false
	}

	return rate, true
}

// ApproximateRates approximates and returns approximated rates on the given date.
// Operates on O(rangeLim) time complexity.
func (r OrderedRecords) ApproximateRates(date date.Date, rangeLim int) (record.Record, bool) {
	recDate := record.DateFromDate(date)

	// find the nearest earlier and later rate records:
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
	rec := record.New()
	for currency, earlierRate := range earlierRec {
		laterRate, found := laterRec.Rate(currency)
		if !found {
			rec[currency] = earlierRate
			continue
		}
		rec[currency] = (earlierRate + laterRate) / 2
	}

	// fill ratesRecord with possible missing rates from later records:
	for currency, laterRate := range laterRec {
		_, found := rec[currency]
		if !found {
			rec[currency] = laterRate
		}
	}

	return rec, true
}

// ApproximateRate approximates and returns approximated rate of the given currency on the given date.
// Operates on O(rangeLim) time complexity.
func (r OrderedRecords) ApproximateRate(date date.Date, string string, rangeLim int) (float32, bool) {
	recDate := record.DateFromDate(date)

	// find the nearest earlier and later rate records:
	earlierRecord, earlierRecordFound := r.nearestEarlierRecord(recDate, rangeLim)
	laterRecord, laterRecordFound := r.nearestLaterRecord(recDate, rangeLim)

	if !earlierRecordFound && !laterRecordFound {
		return 0, false
	}

	if !earlierRecordFound {
		return laterRecord.Rate(string)
	}

	if !laterRecordFound {
		return earlierRecord.Rate(string)
	}

	earlierRate, foundEarlierRate := earlierRecord.Rate(string)
	laterRate, foundLaterRate := laterRecord.Rate(string)

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
// Operates on O(n) time complexity.
func (r OrderedRecords) Convert(date date.Date, amount float32, from string, to string) (float32, error) {
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

// ConvertApproximate approximates rates on the given date
// and uses them to convert amount from one currency to another on the given date.
// Operates on O(n) time complexity.
func (r OrderedRecords) ConvertApproximate(date date.Date, amount float32, from string, to string, rangeLim int) (float32, error) {
	rates, found := r.ApproximateRates(date, rangeLim)
	if !found {
		return 0, fmt.Errorf("approximate rates on %s within range of %d days: %w", date.String(), rangeLim, ErrRateApproximationFailed)
	}

	result, err := rates.Convert(amount, from, to)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// nearestEarlierRecord finds the closest earlier rate record to the given date within the specified range.
// Operates on O(rangeLim) time complexity.
func (r OrderedRecords) nearestEarlierRecord(recDate record.Date, rangeLim int) (record.Record, bool) {
	earlierRecIndex := -1
	for i, rec := range r {
		if rec.Date.Before(recDate) {
			earlierRecIndex = i
			break
		}
	}
	if earlierRecIndex == -1 {
		return nil, false
	}

	rec := r[earlierRecIndex]
	if recDate.SubDays(rec.Date) > rangeLim {
		return nil, false
	}
	return rec.Record, true
}

// nearestLaterRecord finds the closest later rate record to the given date within the specified range.
// Operates on O(rangeLim) time complexity.
func (r OrderedRecords) nearestLaterRecord(recDate record.Date, rangeLim int) (record.Record, bool) {
	laterRecIndex := -1
	for i, rec := range r {
		if rec.Date.After(recDate) {
			laterRecIndex = i
			continue
		}
		break
	}

	if laterRecIndex == -1 {
		return nil, false
	}

	rec := r[laterRecIndex]
	if rec.Date.SubDays(recDate) > rangeLim {
		return nil, false
	}

	return rec.Record, true
}
