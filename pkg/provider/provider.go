package provider

import "errors"

var ErrUnexpectedDataKind = errors.New("unexpected data kind")

type DataKind uint8

const (
	DataKindLatest DataKind = iota
	DataKindTimeSeries
	DataKindTimeSeriesLast90Days
)

// Provider is an interface that defines behavior for fetching currency exchange rate data.
type Provider interface {
	GetRatesData(kind DataKind) ([]byte, error)
}
