package provider

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPProvider is a Provider which uses HTTP GET request to read currency exchange rate records.
type HTTPProvider struct {
	urlLatest               string
	urlTimeSeries           string
	urlTimeSeriesLast90days string
}

// NewHTTPProvider creates a new HTTPProvider.
func NewHTTPProvider(urlLatest string, urlTimeSeries string, urlTimeSeriesLast90Days string) *HTTPProvider {
	return &HTTPProvider{
		urlLatest:               urlLatest,
		urlTimeSeries:           urlTimeSeries,
		urlTimeSeriesLast90days: urlTimeSeriesLast90Days,
	}
}

// GetRatesData file corresponding to the given data kind by its URL.
func (f *HTTPProvider) GetRatesData(kind DataKind) ([]byte, error) {
	// choose URL:
	var url string
	switch kind {
	case DataKindLatest:
		url = f.urlLatest
	case DataKindTimeSeries:
		url = f.urlTimeSeries
	case DataKindTimeSeriesLast90Days:
		url = f.urlTimeSeriesLast90days
	default:
		return nil, ErrUnexpectedDataKind
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: unexpected HTTP status %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
