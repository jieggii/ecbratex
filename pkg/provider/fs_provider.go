package provider

import (
	"os"
)

// FSProvider is the Provider interface implementation which reads
// currency exchange rate records from files stored in the file system.
type FSProvider struct {
	pathLatest               string
	pathTimeSeries           string
	pathTimeSeriesLast90Days string
}

// NewFSProvider creates a new FSProvider.
func NewFSProvider(pathLatest string, pathTimeSeries string, pathTimeSeriesLast90Days string) *FSProvider {
	return &FSProvider{
		pathLatest:               pathLatest,
		pathTimeSeries:           pathTimeSeries,
		pathTimeSeriesLast90Days: pathTimeSeriesLast90Days,
	}
}

// GetRatesData reads a file corresponding to the given data kind and returns its content.
func (f FSProvider) GetRatesData(kind DataKind) ([]byte, error) {
	var path string
	switch kind {
	case DataKindLatest:
		path = f.pathLatest
	case DataKindTimeSeries:
		path = f.pathTimeSeries
	case DataKindTimeSeriesLast90Days:
		path = f.pathTimeSeriesLast90Days
	default:
		return nil, ErrUnexpectedDataKind
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}
