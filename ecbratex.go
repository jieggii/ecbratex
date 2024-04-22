// Package ecbratex provides a convenient interface and useful data structures for fetching
// and manipulating currency exchange rate records from the European Central Bank (https://ecb.europa.eu).
package ecbratex

import (
	"errors"
	"github.com/jieggii/ecbratex/pkg/provider"
	"github.com/jieggii/ecbratex/pkg/record"
	"github.com/jieggii/ecbratex/pkg/timeseries"
	"github.com/jieggii/ecbratex/pkg/xml"
)

// Provider is [provider.Provider] which will be used to fetch rates data across the ecbratex.
// provider.HTTPProvider with URLs to the ECB website is used by default.
var Provider provider.Provider = provider.NewHTTPProvider(
	"https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml",
	"https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml",
	"https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml",
)

// SetProvider sets data provider which will be used across the ecbratex to fetch exchange rates records.
func SetProvider(p provider.Provider) {
	Provider = p
}

var ErrUnexpectedPeriod = errors.New("unexpected period")

type Period uint8

const (
	PeriodWhole Period = iota
	PeriodLast90Days
)

// DataKind converts the period to provider.DataKind.
func (p Period) DataKind() (provider.DataKind, error) {
	switch p {
	case PeriodWhole:
		return provider.DataKindTimeSeries, nil
	case PeriodLast90Days:
		return provider.DataKindTimeSeriesLast90Days, nil
	default:
		return 0, ErrUnexpectedPeriod
	}
}

// FetchLatest fetches latest available exchange rates using Provider.
func FetchLatest() (*record.WithDate, error) {
	rawData, err := Provider.GetRatesData(provider.DataKindLatest)
	if err != nil {
		return nil, err
	}

	xmlData, err := xml.NewData(rawData)
	if err != nil {
		return nil, err
	}

	return record.NewWithDateFromXMLData(xmlData)
}

// FetchTimeSeries fetches rate records within the given period using Provider.
// Returns records represented as timeseries.RecordsMap.
func FetchTimeSeries(period Period) (timeseries.UnorderedRecords, error) {
	dataKind, err := period.DataKind()
	if err != nil {
		return nil, err
	}

	rawData, err := Provider.GetRatesData(dataKind)
	if err != nil {
		return nil, err
	}

	xmlData, err := xml.NewData(rawData)
	if err != nil {
		return nil, err
	}

	return timeseries.NewUnorderedRecordsFromXML(xmlData)
}

// FetchOrderedTimeSeries fetches rate records within the given period using Provider.
// Returns records represented as timeseries.RecordsSlice.
func FetchOrderedTimeSeries(period Period) (timeseries.OrderedRecords, error) {
	dataKind, err := period.DataKind()
	if err != nil {
		return nil, err
	}

	rawData, err := Provider.GetRatesData(dataKind)
	if err != nil {
		return nil, err
	}

	xmlData, err := xml.NewData(rawData)
	if err != nil {
		return nil, err
	}

	return timeseries.NewOrderedRecordsFromXML(xmlData)
}

// FetchOrderedUnorderedTimeSeries fetches rate records within the given period using Provider.
// Returns records represented as timeseries.RecordsSliceMap.
func FetchOrderedUnorderedTimeSeries(period Period) (*timeseries.OrderedUnorderedRecords, error) {
	dataKind, err := period.DataKind()
	if err != nil {
		return nil, err
	}

	rawData, err := Provider.GetRatesData(dataKind)
	if err != nil {
		return nil, err
	}

	xmlData, err := xml.NewData(rawData)
	if err != nil {
		return nil, err
	}

	return timeseries.NewOrderedUnorderedRecordsFromXML(xmlData)
}
