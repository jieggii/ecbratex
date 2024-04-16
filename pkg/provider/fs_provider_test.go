package provider

import (
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func TestNewFSProvider(t *testing.T) {
	const (
		pathLatest               = "latest"
		pathTimeSeries           = "time-series"
		pathTimeSeriesLast90Days = "time-series-last-90-days"
	)

	provider := NewFSProvider(pathLatest, pathTimeSeries, pathTimeSeriesLast90Days)

	assert.Equal(t, pathLatest, provider.pathLatest)
	assert.Equal(t, pathTimeSeries, provider.pathTimeSeries)
	assert.Equal(t, pathTimeSeriesLast90Days, provider.pathTimeSeriesLast90Days)
}

func TestFSProvider_GetRatesData(t *testing.T) {
	provider := NewFSProvider(
		path.Join(testDataPath, "eurofxref-daily.xml"),
		path.Join(testDataPath, "eurofxref-hist.xml"),
		path.Join(testDataPath, "eurofxref-hist-90d.xml"),
	)

	t.Run("latest kind", func(t *testing.T) {
		var kind = DataKindLatest

		data, err := provider.GetRatesData(kind)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, data)
		}
	})

	t.Run("time series kind", func(t *testing.T) {
		var kind = DataKindTimeSeries

		data, err := provider.GetRatesData(kind)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, data)
		}
	})

	t.Run("time series last 90 days kind", func(t *testing.T) {
		var kind = DataKindTimeSeriesLast90Days

		data, err := provider.GetRatesData(kind)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, data)
		}
	})

	t.Run("unexpected data kind", func(t *testing.T) {
		var kind = DataKind(100)

		data, err := provider.GetRatesData(kind)
		if assert.Error(t, err) {
			assert.Empty(t, data)
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		var (
			kind     = DataKindLatest
			provider = NewFSProvider("i-dont-exist", "and-me", "same-lol")
		)

		data, err := provider.GetRatesData(kind)
		if assert.Error(t, err) {
			assert.Empty(t, data)
		}
	})
}
