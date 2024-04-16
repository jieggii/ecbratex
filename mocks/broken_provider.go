package mocks

import (
	"errors"
	"github.com/jieggii/ecbratex/pkg/provider"
)

type BrokenProvider struct{}

func NewBrokenProvider() *BrokenProvider {
	return &BrokenProvider{}
}

func (p *BrokenProvider) GetRatesData(kind provider.DataKind) ([]byte, error) {
	if kind != provider.DataKindLatest && kind != provider.DataKindTimeSeries && kind != provider.DataKindTimeSeriesLast90Days {
		return nil, provider.ErrUnexpectedDataKind
	}
	return nil, errors.New("I failed again")
}
