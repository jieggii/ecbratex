package timeseries

import (
	"fmt"
	"github.com/jieggii/ecbratex/pkg/date"
	"github.com/jieggii/ecbratex/pkg/record"
	"github.com/jieggii/ecbratex/pkg/xml"
)

// OrderedRecords is an implementation of the Records interface.
// It provides O(n) access to rates for any given date. The only benefit of
// using OrderedRecords is that it maintains the order of rates and is
// sorted in anti-chronological order. This means that rates are stored in a slice,
// ensuring they are easily accessible and in the desired order.
//
// Use OrderedRecords when memory constraints are not a concern, and you need
// rates to be ordered for specific operations. If memory is not an issue or the order
// of rates is not important, consider using UnorderedRecords or OrderedUnorderedRecords instead.
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

func (r OrderedRecords) Slice() []record.WithDate {
	return r
}

func (r OrderedRecords) Map() map[record.Date]record.Record {
	result := make(map[record.Date]record.Record)
	for _, rec := range r {
		result[rec.Date] = rec.Record
	}
	return result
}

func (r OrderedRecords) Rates(date date.Date) (record.Record, bool) {
	recDate := record.DateFromDate(date)
	for _, rec := range r {
		if rec.Date == recDate {
			return rec.Record, true
		}
	}
	return nil, false
}

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
