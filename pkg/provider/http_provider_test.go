package provider

import (
	"github.com/jieggii/ecbratex/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testDataPath = "./../../testdata"

func TestNewHTTPProvider(t *testing.T) {
	const (
		urlLatest               = "url-latest"
		urlTimeSeries           = "url-time-series"
		urlTimeSeriesLast90Days = "time-series-last-90-days"
	)

	provider := NewHTTPProvider(urlLatest, urlTimeSeries, urlTimeSeriesLast90Days)
	assert.Equal(t, urlLatest, provider.urlLatest)
	assert.Equal(t, urlTimeSeries, provider.urlTimeSeries)
	assert.Equal(t, urlTimeSeriesLast90Days, provider.urlTimeSeriesLast90days)
}

func TestHTTPProvider_GetRatesData(t *testing.T) {
	var (
		server   = tests.NewTestHTTPServer(testDataPath, false)
		provider = NewHTTPProvider(server.URLLatest, server.URLTimeSeries, server.URLTimeSeriesLast90Days)
	)
	defer server.Close()

	t.Run("latest", func(t *testing.T) {
		data, err := provider.GetRatesData(DataKindLatest)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, data)
		}
	})

	t.Run("time series", func(t *testing.T) {
		data, err := provider.GetRatesData(DataKindTimeSeries)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, data)
		}
	})

	t.Run("time series latest", func(t *testing.T) {
		data, err := provider.GetRatesData(DataKindTimeSeriesLast90Days)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, data)
		}
	})

	t.Run("unexpected kind", func(t *testing.T) {
		data, err := provider.GetRatesData(DataKind(10))
		if assert.Error(t, err) {
			assert.Empty(t, data)
		}
	})

	t.Run("with broken server", func(t *testing.T) {
		var (
			brokenServer = tests.NewTestHTTPServer(testDataPath, true)
			provider     = NewHTTPProvider(brokenServer.URLLatest, brokenServer.URLTimeSeries, brokenServer.URLTimeSeriesLast90Days)
		)
		defer brokenServer.Close()

		data, err := provider.GetRatesData(DataKindLatest)
		if assert.Error(t, err) {
			assert.Empty(t, data)
		}
	})

	t.Run("without server", func(t *testing.T) {
		var provider = NewHTTPProvider("latest", "time-series", "time-series-90-days")

		data, err := provider.GetRatesData(DataKindLatest)
		if assert.Error(t, err) {
			assert.Empty(t, data)
		}
	})

}
