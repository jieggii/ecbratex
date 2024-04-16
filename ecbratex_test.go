package ecbratex

import (
	"github.com/jieggii/ecbratex/mocks"
	"github.com/jieggii/ecbratex/pkg/provider"
	"github.com/jieggii/ecbratex/pkg/record"
	"github.com/jieggii/ecbratex/tests"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

const testDataPath = "./testdata"

const (
	expectedTimeSeriesDataLen           = 61
	expectedTimeSeriesLast90DaysDataLen = 61
)

func TestTimeSeriesPeriod_DataKind(t *testing.T) {
	t.Run("whole period", func(t *testing.T) {
		var period = PeriodWhole

		kind, err := period.DataKind()
		if assert.NoError(t, err) {
			assert.Equal(t, provider.DataKindTimeSeries, kind)
		}
	})

	t.Run("last 90 days period", func(t *testing.T) {
		var period = PeriodLast90Days

		kind, err := period.DataKind()
		if assert.NoError(t, err) {
			assert.Equal(t, provider.DataKindTimeSeriesLast90Days, kind)
		}
	})

	t.Run("unexpected period", func(t *testing.T) {
		var period = Period(100)

		kind, err := period.DataKind()
		if assert.Error(t, err) {
			assert.Zero(t, kind)
		}
	})

}

func TestFetchLatest(t *testing.T) {
	var (
		server = tests.NewTestHTTPServer(testDataPath, false)
	)
	defer server.Close() // todo: t.Cleanup

	t.Run("working provider", func(t *testing.T) {
		var provider = provider.NewHTTPProvider(server.URLLatest, server.URLTimeSeries, server.URLTimeSeriesLast90Days)

		oldProvider := Provider
		SetProvider(provider)
		defer SetProvider(oldProvider)

		data, err := FetchLatest()
		if assert.NoError(t, err) {
			assert.Equal(t, record.NewDate(2024, 2, 27), data.Date)
			assert.Len(t, data.Record, 30)
		}
	})

	t.Run("provider returning invalid XML raw data", func(t *testing.T) {
		var (
			invalidProvider = mocks.NewInvalidProvider()
		)

		oldProvider := Provider
		SetProvider(invalidProvider)
		defer SetProvider(oldProvider)

		data, err := FetchLatest()
		if assert.ErrorIs(t, err, io.EOF) {
			assert.Empty(t, data)
		}
	})

	t.Run("provider returns an error", func(t *testing.T) {
		var (
			brokenProvider = mocks.NewBrokenProvider()
		)

		oldProvider := Provider
		SetProvider(brokenProvider)
		defer SetProvider(oldProvider)

		data, err := FetchLatest()
		if assert.Error(t, err) {
			assert.Empty(t, data)
		}
	})
}

func TestFetchTimeSeries(t *testing.T) {
	var (
		server   = tests.NewTestHTTPServer(testDataPath, false)
		provider = provider.NewHTTPProvider(server.URLLatest, server.URLTimeSeries, server.URLTimeSeriesLast90Days)
	)
	defer server.Close()

	oldProvider := Provider
	SetProvider(provider)
	defer SetProvider(oldProvider)

	t.Run("whole period", func(t *testing.T) {
		var period = PeriodWhole

		data, err := FetchTimeSeries(period)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Len(t, data, expectedTimeSeriesDataLen)
	})

	t.Run("last 90 days period", func(t *testing.T) {
		var period = PeriodLast90Days

		data, err := FetchTimeSeries(period)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Len(t, data, expectedTimeSeriesLast90DaysDataLen)
	})

	t.Run("unexpected period", func(t *testing.T) {
		var period = Period(10)

		data, err := FetchTimeSeries(period)
		if assert.ErrorIs(t, ErrUnexpectedPeriod, err) {
			assert.Empty(t, data)
		}
	})

	t.Run("provider returns invalid XML raw data", func(t *testing.T) {
		var (
			period          = PeriodLast90Days
			invalidProvider = mocks.NewInvalidProvider()
		)

		oldProvider := Provider
		SetProvider(invalidProvider)
		defer SetProvider(oldProvider)

		data, err := FetchTimeSeries(period)
		if assert.ErrorIs(t, err, io.EOF) {
			assert.Empty(t, data)
		}
	})

	t.Run("provider returns an error", func(t *testing.T) {
		var (
			period         = PeriodLast90Days
			brokenProvider = mocks.NewBrokenProvider()
		)

		oldProvider := Provider
		SetProvider(brokenProvider)
		defer SetProvider(oldProvider)

		data, err := FetchTimeSeries(period)
		if assert.Error(t, err) {
			assert.Empty(t, data)
		}
	})
}

func TestFetchOrderedTimeSeries(t *testing.T) {
	var (
		server   = tests.NewTestHTTPServer(testDataPath, false)
		provider = provider.NewHTTPProvider(server.URLLatest, server.URLTimeSeries, server.URLTimeSeriesLast90Days)
	)
	defer server.Close()

	oldProvider := Provider
	SetProvider(provider)
	defer SetProvider(oldProvider)

	t.Run("whole period", func(t *testing.T) {
		var period = PeriodWhole

		data, err := FetchOrderedTimeSeries(period)
		assert.NoError(t, err)
		assert.Len(t, data, expectedTimeSeriesDataLen)
	})

	t.Run("last 90 days period", func(t *testing.T) {
		var period = PeriodLast90Days

		data, err := FetchOrderedTimeSeries(period)
		if assert.NoError(t, err) {
			assert.Len(t, data, expectedTimeSeriesLast90DaysDataLen)
		}
	})

	t.Run("unexpected period", func(t *testing.T) {
		var period = Period(10)

		data, err := FetchOrderedTimeSeries(period)
		if assert.ErrorIs(t, ErrUnexpectedPeriod, err) {
			assert.Empty(t, data)
		}
	})

	t.Run("provider returns invalid XML raw data", func(t *testing.T) {
		var (
			period          = PeriodLast90Days
			invalidProvider = mocks.NewInvalidProvider()
		)

		oldProvider := Provider
		SetProvider(invalidProvider)
		defer SetProvider(oldProvider)

		data, err := FetchOrderedTimeSeries(period)
		if assert.ErrorIs(t, err, io.EOF) {
			assert.Empty(t, data)
		}
	})

	t.Run("provider returns an error", func(t *testing.T) {
		var (
			period         = PeriodLast90Days
			brokenProvider = mocks.NewBrokenProvider()
		)

		oldProvider := Provider
		SetProvider(brokenProvider)
		defer SetProvider(oldProvider)

		data, err := FetchOrderedTimeSeries(period)
		if assert.Error(t, err) {
			assert.Empty(t, data)
		}
	})
}

func TestFetchOrderedUnorderedTimeSeries(t *testing.T) {
	var (
		server   = tests.NewTestHTTPServer(testDataPath, false)
		provider = provider.NewHTTPProvider(server.URLLatest, server.URLTimeSeries, server.URLTimeSeriesLast90Days)
	)
	defer server.Close()

	oldProvider := Provider
	SetProvider(provider)
	defer SetProvider(oldProvider)

	t.Run("whole period", func(t *testing.T) {
		var period = PeriodWhole

		data, err := FetchOrderedUnorderedTimeSeries(period)
		if assert.NoError(t, err) {
			assert.Len(t, data.Map(), expectedTimeSeriesDataLen)
			assert.Len(t, data.Slice(), expectedTimeSeriesDataLen)
		}
	})

	t.Run("last 90 days period", func(t *testing.T) {
		var period = PeriodLast90Days

		data, err := FetchOrderedUnorderedTimeSeries(period)
		if assert.NoError(t, err) {
			assert.Len(t, data.Map(), expectedTimeSeriesLast90DaysDataLen)
			assert.Len(t, data.Slice(), expectedTimeSeriesLast90DaysDataLen)
		}
	})

	t.Run("unexpected period", func(t *testing.T) {
		var period = Period(10)

		data, err := FetchOrderedUnorderedTimeSeries(period)
		if assert.ErrorIs(t, ErrUnexpectedPeriod, err) {
			assert.Empty(t, data)
		}
	})

	t.Run("provider returns invalid XML raw data", func(t *testing.T) {
		var (
			period          = PeriodLast90Days
			invalidProvider = mocks.NewInvalidProvider()
		)

		oldProvider := Provider
		SetProvider(invalidProvider)
		defer SetProvider(oldProvider)

		data, err := FetchOrderedUnorderedTimeSeries(period)
		if assert.ErrorIs(t, err, io.EOF) {
			assert.Empty(t, data)
		}
	})

	t.Run("provider returns an error", func(t *testing.T) {
		var (
			period         = PeriodLast90Days
			brokenProvider = mocks.NewBrokenProvider()
		)

		oldProvider := Provider
		SetProvider(brokenProvider)
		defer SetProvider(oldProvider)

		data, err := FetchOrderedUnorderedTimeSeries(period)
		if assert.Error(t, err) {
			assert.Empty(t, data)
		}
	})
}
