package mocks

import "github.com/jieggii/ecbratex/pkg/provider"

type InvalidProvider struct{}

func NewInvalidProvider() *InvalidProvider {
	return &InvalidProvider{}
}

func (p *InvalidProvider) GetRatesData(kind provider.DataKind) ([]byte, error) {
	if kind != provider.DataKindLatest && kind != provider.DataKindTimeSeries && kind != provider.DataKindTimeSeriesLast90Days {
		return nil, provider.ErrUnexpectedDataKind
	}
	return []byte("some invalid data"), nil
}
