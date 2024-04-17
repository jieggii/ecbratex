package timeseries

import (
	"github.com/jieggii/ecbratex/pkg/record"
	"github.com/jieggii/ecbratex/pkg/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOrderedRecords(t *testing.T) {
	recs := NewOrderedRecords()
	assert.Len(t, recs, 0)
}

func TestNewOrderedRecordsFromXML(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		data := &xml.Data{
			Cubes: []xml.DataCube{
				{Date: "2024-12-02", Rates: []xml.DataCubeRate{{"USD", 0.9}}},
			},
		}

		records, err := NewOrderedRecordsFromXML(data)
		if assert.NoError(t, err) {
			assert.Len(t, records, 1)
		}
	})

	t.Run("data with invalid date", func(t *testing.T) {
		data := &xml.Data{
			Cubes: []xml.DataCube{
				{Date: "invalid date", Rates: []xml.DataCubeRate{{"USD", 0.9}}},
			},
		}
		records, err := NewOrderedRecordsFromXML(data)
		if assert.Error(t, err) {
			assert.Empty(t, records)
		}
	})

	t.Run("nil data", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = NewOrderedRecordsFromXML(nil)
		})
	})

	t.Run("zero data", func(t *testing.T) {
		data := &xml.Data{}
		records, err := NewOrderedRecordsFromXML(data)
		if assert.NoError(t, err) {
			assert.Empty(t, records)
		}
	})
}

func TestOrderedRecords_Slice(t *testing.T) {
	var rec = record.NewWithDate(
		record.Record{"USD": 0.9, "RUB": 0.01},
		record.NewDate(2000, 1, 1),
	)
	records := OrderedRecords{rec}
	assert.Equal(t, []record.WithDate{rec}, records.Slice())
}

func TestOrderedRecords_Map(t *testing.T) {
	var (
		rec  = record.Record{"USD": 0.9, "RUB": 0.01}
		date = record.NewDate(2000, 10, 3)
	)
	recWithDate := record.NewWithDate(rec, date)
	records := OrderedRecords{recWithDate}
	assert.Equal(t, map[record.Date]record.Record{date: rec}, records.Map())
}

func TestOrderedRecords_Rates(t *testing.T) {
	const (
		USDRate float32 = 0.9
		RUBRate float32 = 0.01
	)

	records := OrderedRecords{
		record.NewWithDate(
			record.Record{"USD": USDRate, "RUB": RUBRate},
			record.NewDate(2024, 12, 1),
		),
	}

	t.Run("get rates on existing date", func(t *testing.T) {
		rates, found := records.Rates(record.NewDate(2024, 12, 1))
		if assert.True(t, found) {
			assert.Equal(t, record.Record{"USD": USDRate, "RUB": RUBRate}, rates)
		}
	})

	t.Run("get rates on non-existent date", func(t *testing.T) {
		rates, found := records.Rates(record.NewDate(1090, 2, 2))
		if assert.False(t, found) {
			assert.Empty(t, rates)
		}
	})
}

func TestOrderedRecords_Rate(t *testing.T) {
	const (
		USDRate float32 = 0.9
		RUBRate float32 = 0.01
	)

	t.Run("get existing rate on existing date", func(t *testing.T) {
		var (
			date    = record.NewDate(2024, 12, 1)
			records = OrderedRecords{
				record.NewWithDate(
					record.Record{"USD": USDRate, "RUB": RUBRate},
					date,
				),
			}
		)

		rate, found := records.Rate(date, "USD")
		if assert.True(t, found) {
			assert.Equal(t, USDRate, rate)
		}
	})

	t.Run("get non-existent rate on existing date", func(t *testing.T) {
		var (
			date    = record.NewDate(2024, 12, 1)
			records = OrderedRecords{
				record.NewWithDate(
					record.Record{"USD": USDRate, "RUB": RUBRate},
					date,
				),
			}
		)

		rate, found := records.Rate(date, "XXX")
		if assert.False(t, found) {
			assert.Zero(t, rate)
		}
	})

	t.Run("get rate on non-existing date", func(t *testing.T) {
		var records = OrderedRecords{
			record.NewWithDate(
				record.Record{"USD": USDRate, "RUB": RUBRate},
				record.NewDate(2024, 12, 1),
			),
		}

		rate, found := records.Rate(record.NewDate(2000, 1, 5), "USD")
		if assert.False(t, found) {
			assert.Zero(t, rate)
		}
	})
}

func TestOrderedRecords_ApproximateRates(t *testing.T) {
	t.Run("rangeLim covering both earlier and later records", func(t *testing.T) {
		const rangeLim = 100

		var (
			record1Date = record.NewDate(2000, 1, 30)
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}

			record2Date = record.NewDate(2000, 1, 1)
			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
		)

		records := OrderedRecords{
			record.NewWithDate(record1, record1Date),
			record.NewWithDate(record2, record2Date),
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 16), rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, (record2["USD"]+record1["USD"])/2, result["USD"])
			assert.Equal(t, (record2["RUB"]+record1["RUB"])/2, result["RUB"])
		}
	})

	t.Run("rangeLim covering only earlier record", func(t *testing.T) {
		const rangeLim = 1

		var (
			record1Date = record.NewDate(2000, 1, 30)
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}

			record2Date = record.NewDate(2000, 1, 1)
			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
		)

		records := OrderedRecords{
			record.NewWithDate(record1, record1Date),
			record.NewWithDate(record2, record2Date),
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 2), rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, record2["USD"], result["USD"])
			assert.Equal(t, record2["RUB"], result["RUB"])
		}
	})

	t.Run("rangeLim covering only later record", func(t *testing.T) {
		const rangeLim = 1

		var (
			record1Date = record.NewDate(2000, 1, 30)
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}

			record2Date = record.NewDate(2000, 1, 1)
			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
		)

		records := OrderedRecords{
			record.NewWithDate(record1, record1Date),
			record.NewWithDate(record2, record2Date),
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 29), rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, record1["USD"], result["USD"])
			assert.Equal(t, record1["RUB"], result["RUB"])
		}
	})

	t.Run("rangeLim covering neither earlier nor later record", func(t *testing.T) {
		const rangeLim = 5

		var (
			record1Date = record.NewDate(2000, 1, 30)
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}

			record2Date = record.NewDate(2000, 1, 1)
			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
		)

		records := OrderedRecords{
			record.NewWithDate(record1, record1Date),
			record.NewWithDate(record2, record2Date),
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 15), rangeLim)
		if assert.False(t, found) {
			assert.Empty(t, result)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing RUB rate in the earlier record", func(t *testing.T) {
		const rangeLim = 100

		var (
			record1Date = record.NewDate(2000, 1, 30)
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}

			record2Date = record.NewDate(2000, 1, 1)
			record2     = record.Record{"USD": 0.1}
		)

		records := OrderedRecords{
			record.NewWithDate(record1, record1Date),
			record.NewWithDate(record2, record2Date),
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 16), rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, (record1["USD"]+record2["USD"])/2, result["USD"])
			assert.Equal(t, record1["RUB"], result["RUB"])
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing RUB rate in the later record", func(t *testing.T) {
		const rangeLim = 100

		var (
			record1Date = record.NewDate(2000, 1, 30)
			record1     = record.Record{"USD": 1.0}

			record2Date = record.NewDate(2000, 1, 1)
			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
		)

		records := OrderedRecords{
			record.NewWithDate(record1, record1Date),
			record.NewWithDate(record2, record2Date),
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 16), rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, (record2["USD"]+record1["USD"])/2, result["USD"])
			assert.Equal(t, record2["RUB"], result["RUB"])
		}
	})

	t.Run("empty records", func(t *testing.T) {
		const rangeLim = 100

		records := OrderedRecords{}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 16), rangeLim)
		if assert.False(t, found) {
			assert.Zero(t, result)
		}
	})
}

func TestOrderedRecords_ApproximateRate(t *testing.T) {
	t.Run("rangeLim covering both earlier and later records", func(t *testing.T) {
		const rangeLim = 100

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 16), "USD", rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, (record1["USD"]+record2["USD"])/2, rate)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing rate in the earlier record", func(t *testing.T) {
		const rangeLim = 100

		var (
			record1     = record.Record{"USD": 1.0}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 16), "RUB", rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, record2["RUB"], rate)
		}
	})

	// todo
	t.Run("rangeLim covering both earlier and later records, missing rate in the later record", func(t *testing.T) {
		const rangeLim = 100

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.1}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 16), "RUB", rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, record1["RUB"], rate)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing rate in the both records", func(t *testing.T) {
		const rangeLim = 100

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 16), "XXX", rangeLim)
		if assert.False(t, found) {
			assert.Zero(t, rate)
		}
	})

	t.Run("rangeLim covering only earlier record", func(t *testing.T) {
		const rangeLim = 1

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 2), "USD", rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, record2["USD"], rate)
		}
	})

	t.Run("rangeLim covering only later record", func(t *testing.T) {
		const rangeLim = 1

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 29), "USD", rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, record1["USD"], rate)
		}
	})

	t.Run("rangeLim covering neither earlier nor later record", func(t *testing.T) {
		const rangeLim = 5

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 15), "USD", rangeLim)
		if assert.False(t, found) {
			assert.Zero(t, rate)
		}
	})

}

func TestOrderedRecords_Convert(t *testing.T) {
	const (
		USDRate float32 = 1
		RUBRate float32 = 0.5
	)

	var (
		recDate = record.NewDate(2020, 1, 1)
		records = OrderedRecords{
			record.NewWithDate(record.Record{"USD": 1, "RUB": 0.5}, recDate),
		}
	)

	t.Run("existing date, existing from and existing to", func(t *testing.T) {
		var amount float32 = 12

		result, err := records.Convert(recDate, amount, "USD", "RUB")
		if assert.NoError(t, err) {
			assert.Equal(t, amount*(USDRate/RUBRate), result)
		}
	})

	t.Run("existing date, non-existing from and existing to", func(t *testing.T) {
		result, err := records.Convert(recDate, 1, "XXX", "RUB")
		if assert.ErrorIs(t, err, record.ErrRateNotFound) {
			assert.Zero(t, result)
		}
	})

	t.Run("existing date, existing from and non-existing to", func(t *testing.T) {
		result, err := records.Convert(recDate, 1, "USD", "XXX")
		if assert.ErrorIs(t, err, record.ErrRateNotFound) {
			assert.Zero(t, result)
		}
	})

	t.Run("non-existing date", func(t *testing.T) {
		result, err := records.Convert(record.NewDate(1500, 1, 1), 1, "USD", "XXX")
		if assert.ErrorIs(t, err, ErrRatesRecordNotFound) {
			assert.Zero(t, result)
		}
	})

}

func TestOrderedRecords_ConvertApproximate(t *testing.T) {
	t.Run("rangeLim covering both earlier and later records", func(t *testing.T) {
		const (
			amount   float32 = 15
			rangeLim         = 100
		)

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.1, "RUB": 0.5}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 16), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = (record1["USD"] + record2["USD"]) / 2
				rubRate = (record1["RUB"] + record2["RUB"]) / 2
			)
			assert.Equal(t, amount*(usdRate/rubRate), result)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing TO rate in the earlier record", func(t *testing.T) {
		const (
			amount   float32 = 12
			rangeLim         = 100
		)

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.1}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.1}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 16), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = (record1["USD"] + record2["USD"]) / 2
				rubRate = record1["RUB"]
			)
			assert.Equal(t, amount*(usdRate/rubRate), result)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing TO rate in the later record", func(t *testing.T) {
		const (
			amount   float32 = 11
			rangeLim         = 100
		)

		var (
			record1     = record.Record{"USD": 1.0}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.9, "RUB": 0.01}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 16), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = (record1["USD"] + record2["USD"]) / 2
				rubRate = record2["RUB"]
			)
			assert.Equal(t, amount*(usdRate/rubRate), result)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing TO rate in the both records", func(t *testing.T) {
		const (
			amount   float32 = 9999
			rangeLim         = 100
		)

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.01}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.9, "RUB": 0.01}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 16), amount, "USD", "XXX", rangeLim)
		if assert.Error(t, err) {
			assert.Zero(t, result)
		}
	})

	t.Run("rangeLim covering only earlier record", func(t *testing.T) {
		const (
			amount   float32 = 2
			rangeLim         = 1
		)

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.02}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.9, "RUB": 0.01}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 2), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = record2["USD"]
				rubRate = record2["RUB"]
			)
			assert.Equal(t, amount*(usdRate/rubRate), result)
		}
	})

	t.Run("rangeLim covering only later record", func(t *testing.T) {
		const (
			amount   float32 = 123
			rangeLim         = 1
		)

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.02}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.9, "RUB": 0.01}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 29), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = record1["USD"]
				rubRate = record1["RUB"]
			)
			assert.Equal(t, amount*(usdRate/rubRate), result)
		}
	})

	t.Run("rangeLim covering neither earlier nor later record", func(t *testing.T) {
		const (
			amount   float32 = 123123
			rangeLim         = 5
		)

		var (
			record1     = record.Record{"USD": 1.0, "RUB": 0.02}
			record1Date = record.NewDate(2000, 1, 30)

			record2     = record.Record{"USD": 0.9, "RUB": 0.01}
			record2Date = record.NewDate(2000, 1, 1)

			records = OrderedRecords{
				record.NewWithDate(record1, record1Date),
				record.NewWithDate(record2, record2Date),
			}
		)

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 15), amount, "USD", "RUB", rangeLim)
		if assert.ErrorIs(t, err, ErrRateApproximationFailed) {
			assert.Zero(t, result)
		}
	})
}
