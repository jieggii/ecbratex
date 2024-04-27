package timeseries

import (
	"github.com/jieggii/ecbratex/pkg/record"
	"github.com/jieggii/ecbratex/pkg/xml"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNewUnorderedRecordsFromXML(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		const USDRate float32 = 0.9
		var date = record.NewDate(2024, 12, 2)

		data := &xml.Data{
			Cubes: []xml.DataCube{
				{Date: date.String(), Rates: []xml.DataCubeRate{
					{"USD", USDRate},
				}},
			},
		}
		records, err := NewUnorderedRecordsFromXML(data)
		if assert.NoError(t, err) {
			assert.Equal(
				t,
				UnorderedRecords{
					date: record.Record{"EUR": 1, "USD": USDRate},
				},
				records,
			)
		}
	})

	t.Run("data with invalid date", func(t *testing.T) {
		data := &xml.Data{
			Cubes: []xml.DataCube{
				{Date: "invalid date", Rates: []xml.DataCubeRate{
					{"USD", 0.9},
				}},
			},
		}
		records, err := NewUnorderedRecordsFromXML(data)
		if assert.Error(t, err) {
			assert.Empty(t, records)
		}
	})

	t.Run("nil data", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = NewUnorderedRecordsFromXML(nil)
		})
	})

	t.Run("zero data", func(t *testing.T) {
		data := &xml.Data{}
		records, err := NewUnorderedRecordsFromXML(data)
		if assert.NoError(t, err) {
			assert.Empty(t, records)
		}
	})
}

func TestUnorderedRecords_Slice(t *testing.T) {
	var (
		date1 = record.NewDate(2000, 1, 3)
		rec1  = record.Record{"USD": 0.7}

		date2 = record.NewDate(2000, 1, 2)
		rec2  = record.Record{"USD": 0.8}

		date3 = record.NewDate(2000, 1, 1)
		rec3  = record.Record{"USD": 0.9}
	)

	records := UnorderedRecords{
		date1: rec1,
		date2: rec2,
		date3: rec3,
	}

	assert.Equal(
		t,
		[]record.WithDate{
			record.NewWithDate(rec1, date1),
			record.NewWithDate(rec2, date2),
			record.NewWithDate(rec3, date3),
		},
		records.Slice(),
	)
}

func TestUnorderedRecords_Map(t *testing.T) {
	var (
		recDate = record.NewDate(2000, 1, 1)
		rec     = record.Record{"USD": 0.9, "RUB": 0.01}
	)
	records := UnorderedRecords{recDate: rec}
	assert.Equal(t, map[record.Date]record.Record{recDate: rec}, records.Map())
}

func TestUnorderedRecords_Rates(t *testing.T) {
	const (
		USDRate float32 = 0.9
		RUBRate float32 = 0.01
	)
	var date = record.NewDate(2024, 12, 1)

	records := UnorderedRecords{
		date: {
			"USD": USDRate,
			"RUB": RUBRate,
		},
	}

	t.Run("get rates on existing date", func(t *testing.T) {
		rates, found := records.Rates(date)
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

func TestUnorderedRecords_Rate(t *testing.T) {
	const (
		USDRate float32 = 0.9
		RUBRate float32 = 0.01
	)
	var date = record.NewDate(2024, 12, 1)

	t.Run("get existing rate on existing date", func(t *testing.T) {
		var records = UnorderedRecords{
			date: {
				"USD": USDRate,
				"RUB": RUBRate,
			},
		}

		rate, found := records.Rate(date, "USD")
		if assert.True(t, found) {
			assert.Equal(t, USDRate, rate)
		}
	})

	t.Run("get non-existent rate on existing date", func(t *testing.T) {
		var records = UnorderedRecords{
			date: {
				"USD": USDRate,
				"RUB": RUBRate,
			},
		}

		rate, found := records.Rate(date, "XXX")
		if assert.False(t, found) {
			assert.Zero(t, rate)
		}
	})

	t.Run("get rate on non-existing date", func(t *testing.T) {
		var records = UnorderedRecords{
			date: {
				"USD": USDRate,
				"RUB": RUBRate,
			},
		}

		rate, found := records.Rate(record.NewDate(1500, 1, 5), "USD")
		if assert.False(t, found) {
			assert.Zero(t, rate)
		}
	})
}

func TestUnorderedRecords_ApproximateRates(t *testing.T) {
	t.Run("rangeLim covering both earlier and later records", func(t *testing.T) {
		const rangeLim = 100

		var (
			date1 = record.NewDate(2000, 1, 1)
			rec1  = record.Record{"USD": 0.1, "RUB": 0.5}

			date2 = record.NewDate(2000, 1, 30)
			rec2  = record.Record{"USD": 1.0, "RUB": 0.1}
		)

		records := UnorderedRecords{
			date1: rec1,
			date2: rec2,
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 16), rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, (rec1["USD"]+rec2["USD"])/2, result["USD"])
			assert.Equal(t, (rec1["RUB"]+rec2["RUB"])/2, result["RUB"])
		}
	})

	t.Run("rangeLim covering only earlier record", func(t *testing.T) {
		const rangeLim = 1

		var (
			date1 = record.NewDate(2000, 1, 1)
			rec1  = record.Record{"USD": 0.1, "RUB": 0.5}

			date2 = record.NewDate(2000, 1, 30)
			rec2  = record.Record{"USD": 1.0, "RUB": 0.1}
		)

		records := UnorderedRecords{
			date1: rec1,
			date2: rec2,
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 2), rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, rec1["USD"], result["USD"])
			assert.Equal(t, rec1["RUB"], result["RUB"])
		}
	})

	t.Run("rangeLim covering only later record", func(t *testing.T) {
		const rangeLim = 1

		var (
			date1 = record.NewDate(2000, 1, 1)
			rec1  = record.Record{"USD": 0.1, "RUB": 0.5}

			date2 = record.NewDate(2000, 1, 30)
			rec2  = record.Record{"USD": 1.0, "RUB": 0.1}
		)

		records := UnorderedRecords{
			date1: rec1,
			date2: rec2,
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 29), rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, rec2["USD"], result["USD"])
			assert.Equal(t, rec2["RUB"], result["RUB"])
		}
	})

	t.Run("rangeLim covering neither earlier nor later record", func(t *testing.T) {
		const rangeLim = 5

		var (
			date1 = record.NewDate(2000, 1, 1)
			rec1  = record.Record{"USD": 0.1, "RUB": 0.5}

			date2 = record.NewDate(2000, 1, 30)
			rec2  = record.Record{"USD": 1.0, "RUB": 0.1}
		)

		records := UnorderedRecords{
			date1: rec1,
			date2: rec2,
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 15), rangeLim)
		if assert.False(t, found) {
			assert.Empty(t, result)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing RUB rate in the earlier record", func(t *testing.T) {
		const rangeLim = 100

		var (
			date1 = record.NewDate(2000, 1, 1)
			rec1  = record.Record{"USD": 0.1}

			date2 = record.NewDate(2000, 1, 30)
			rec2  = record.Record{"USD": 1.0, "RUB": 0.1}
		)

		records := UnorderedRecords{
			date1: rec1,
			date2: rec2,
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 16), rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, (rec1["USD"]+rec2["USD"])/2, result["USD"])
			assert.Equal(t, rec2["RUB"], result["RUB"])
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing RUB rate in the later record", func(t *testing.T) {
		const rangeLim = 100

		var (
			date1 = record.NewDate(2000, 1, 1)
			rec1  = record.Record{"USD": 0.1, "RUB": 0.5}

			date2 = record.NewDate(2000, 1, 30)
			rec2  = record.Record{"USD": 1.0}
		)

		records := UnorderedRecords{
			date1: rec1,
			date2: rec2,
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 16), rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, (rec1["USD"]+rec2["USD"])/2, result["USD"])
			assert.Equal(t, rec1["RUB"], result["RUB"])
		}
	})

	t.Run("zero rangeLim", func(t *testing.T) {
		const rangeLim = 0

		var (
			date1 = record.NewDate(2000, 1, 30)
			rec1  = record.Record{}

			date2 = record.NewDate(2000, 1, 1)
			rec2  = record.Record{}
		)

		records := OrderedRecords{
			record.NewWithDate(rec1, date1),
			record.NewWithDate(rec2, date2),
		}

		result, found := records.ApproximateRates(record.NewDate(2000, 1, 16), rangeLim)
		if assert.False(t, found) {
			assert.Zero(t, result)
		}
	})
}

func TestUnorderedRecords_ApproximateRate(t *testing.T) {
	t.Run("rangeLim covering both earlier and later records", func(t *testing.T) {
		const rangeLim = 100

		var (
			date1 = record.NewDate(2000, 1, 1)
			rec1  = record.Record{"USD": 0.1, "RUB": 0.5}

			date2 = record.NewDate(2000, 1, 30)
			rec2  = record.Record{"USD": 1.0, "RUB": 0.1}
		)

		records := UnorderedRecords{
			date1: rec1,
			date2: rec2,
		}

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 16), "USD", rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, (rec1["USD"]+rec2["USD"])/2, rate)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing rate in the earlier record", func(t *testing.T) {
		const rangeLim = 100

		var (
			date1 = record.NewDate(2000, 1, 1)
			rec1  = record.Record{"USD": 0.1}

			date2 = record.NewDate(2000, 1, 30)
			rec2  = record.Record{"USD": 1.0, "RUB": 0.22}
		)

		records := UnorderedRecords{
			date1: rec1,
			date2: rec2,
		}

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 16), "RUB", rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, rec2["RUB"], rate)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing rate in the later record", func(t *testing.T) {
		const rangeLim = 100

		var (
			date1 = record.NewDate(2000, 1, 1)
			rec1  = record.Record{"USD": 0.1, "RUB": 0.22}

			date2 = record.NewDate(2000, 1, 30)
			rec2  = record.Record{"USD": 1.0}
		)

		records := UnorderedRecords{
			date1: rec1,
			date2: rec2,
		}

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 16), "RUB", rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, rec1["RUB"], rate)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing rate in the both records", func(t *testing.T) {
		const rangeLim = 100

		var (
			date1 = record.NewDate(2000, 1, 1)
			rec1  = record.Record{"USD": 0.1, "RUB": 0.22}

			date2 = record.NewDate(2000, 1, 30)
			rec2  = record.Record{"USD": 1.0, "RUB": 0.123}
		)

		records := UnorderedRecords{
			date1: rec1,
			date2: rec2,
		}

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 16), "XXX", rangeLim)
		if assert.False(t, found) {
			assert.Zero(t, rate)
		}
	})

	t.Run("rangeLim covering only earlier record", func(t *testing.T) {
		const rangeLim = 1

		var (
			record1 = record.Record{"USD": 0.1, "RUB": 0.5}
			record2 = record.Record{"USD": 1.0, "RUB": 0.1}
		)

		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 2), "USD", rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, record1["USD"], rate)
		}
	})

	t.Run("rangeLim covering only later record", func(t *testing.T) {
		const rangeLim = 1

		var (
			record1 = record.Record{"USD": 0.1, "RUB": 0.5}
			record2 = record.Record{"USD": 1.0, "RUB": 0.1}
		)

		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 29), "USD", rangeLim)
		if assert.True(t, found) {
			assert.Equal(t, record2["USD"], rate)
		}
	})

	t.Run("rangeLim covering neither earlier nor later record", func(t *testing.T) {
		const rangeLim = 5

		var (
			record1 = record.Record{"USD": 0.1, "RUB": 0.5}
			record2 = record.Record{"USD": 1.0, "RUB": 0.1}
		)

		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

		rate, found := records.ApproximateRate(record.NewDate(2000, 1, 15), "USD", rangeLim)
		if assert.False(t, found) {
			assert.Zero(t, rate)
		}
	})
}

func TestUnorderedRecords_Convert(t *testing.T) {
	var date = record.NewDate(2020, 1, 1)

	const (
		USDRate float32 = 0.9
		RUBRate float32 = 0.01
	)

	records := UnorderedRecords{
		date: record.Record{"USD": USDRate, "RUB": RUBRate},
	}

	t.Run("existing date, existing from and existing to", func(t *testing.T) {
		var amount float32 = 12

		result, err := records.Convert(date, amount, "USD", "RUB")
		if assert.NoError(t, err) {
			assert.Equal(t, amount*(USDRate/RUBRate), result)
		}
	})

	t.Run("existing date, non-existing from and existing to", func(t *testing.T) {
		result, err := records.Convert(date, 99, "XXX", "RUB")
		if assert.ErrorIs(t, err, record.ErrRateNotFound) {
			assert.Zero(t, result)
		}
	})

	t.Run("existing date, existing from and non-existing to", func(t *testing.T) {
		result, err := records.Convert(date, 99, "USD", "XXX")
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

func TestUnorderedRecords_ConvertApproximate(t *testing.T) {
	t.Run("rangeLim covering both earlier and later records", func(t *testing.T) {
		const (
			amount   float32 = 15
			rangeLim         = 100
		)

		var (
			record1 = record.Record{"USD": 0.1, "RUB": 0.5}
			record2 = record.Record{"USD": 1.0, "RUB": 0.1}
		)
		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

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
			record1 = record.Record{"USD": 0.1}
			record2 = record.Record{"USD": 1.0, "RUB": 0.22}
		)

		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 16), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = (record1["USD"] + record2["USD"]) / 2
				rubRate = record2["RUB"]
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
			record1 = record.Record{"USD": 0.1, "RUB": 0.22}
			record2 = record.Record{"USD": 1.0}
		)

		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 16), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = (record1["USD"] + record2["USD"]) / 2
				rubRate = record1["RUB"]
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
			record1 = record.Record{"USD": 0.1, "RUB": 0.22}
			record2 = record.Record{"USD": 1.0, "RUB": 0.123}
		)
		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

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
			record1 = record.Record{"USD": 0.1, "RUB": 0.5}
			record2 = record.Record{"USD": 1.0, "RUB": 0.1}
		)
		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 2), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = record1["USD"]
				rubRate = record1["RUB"]
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
			record1 = record.Record{"USD": 0.1, "RUB": 0.5}
			record2 = record.Record{"USD": 1.0, "RUB": 0.1}
		)
		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 29), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = record2["USD"]
				rubRate = record2["RUB"]
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
			record1 = record.Record{"USD": 0.1, "RUB": 0.5}
			record2 = record.Record{"USD": 1.0, "RUB": 0.1}
		)
		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

		result, err := records.ConvertApproximate(record.NewDate(2000, 1, 15), amount, "USD", "RUB", rangeLim)
		if assert.ErrorIs(t, err, ErrRateApproximationFailed) {
			assert.Zero(t, result)
		}
	})
}

func TestUnorderedRecords_ConvertMinors(t *testing.T) {
	var date = record.NewDate(2020, 1, 1)

	const (
		USDRate float32 = 0.9
		RUBRate float32 = 0.01
	)

	records := UnorderedRecords{
		date: record.Record{"USD": USDRate, "RUB": RUBRate},
	}

	t.Run("existing date, existing from and existing to", func(t *testing.T) {
		var amount int64 = 1234 // 12.34 USD

		result, err := records.ConvertMinors(date, amount, "USD", "RUB")
		if assert.NoError(t, err) {
			expectedResult := int64(math.Round(
				float64(
					float32(amount) * (USDRate / RUBRate),
				),
			))
			assert.Equal(t, expectedResult, result)
		}
	})

	t.Run("existing date, non-existing from and existing to", func(t *testing.T) {
		result, err := records.ConvertMinors(date, 1234, "XXX", "RUB")
		if assert.ErrorIs(t, err, record.ErrRateNotFound) {
			assert.Zero(t, result)
		}
	})

	t.Run("existing date, existing from and non-existing to", func(t *testing.T) {
		result, err := records.ConvertMinors(date, 1234, "USD", "XXX")
		if assert.ErrorIs(t, err, record.ErrRateNotFound) {
			assert.Zero(t, result)
		}
	})

	t.Run("non-existing date", func(t *testing.T) {
		result, err := records.ConvertMinors(record.NewDate(1500, 1, 1), 1234, "USD", "XXX")
		if assert.ErrorIs(t, err, ErrRatesRecordNotFound) {
			assert.Zero(t, result)
		}
	})
}

func TestUnorderedRecords_ConvertMinorsApproximate(t *testing.T) {
	t.Run("rangeLim covering both earlier and later records", func(t *testing.T) {
		const (
			amount   = 1599 // 15.99 USD
			rangeLim = 100
		)

		var (
			rec1 = record.Record{"USD": 0.1, "RUB": 0.5}
			rec2 = record.Record{"USD": 1.0, "RUB": 0.1}
		)
		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  rec1,
			record.NewDate(2000, 1, 30): rec2,
		}

		result, err := records.ConvertMinorsApproximate(record.NewDate(2000, 1, 16), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = (rec1["USD"] + rec2["USD"]) / 2
				rubRate = (rec1["RUB"] + rec2["RUB"]) / 2
			)
			expectedResult := int64(math.Round(amount * float64(usdRate/rubRate)))
			assert.Equal(t, expectedResult, result)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing TO rate in the earlier record", func(t *testing.T) {
		const (
			amount   = 1299 // 12.99 USD
			rangeLim = 100
		)

		var (
			rec1 = record.Record{"USD": 0.1}
			rec2 = record.Record{"USD": 1.0, "RUB": 0.22}
		)

		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  rec1,
			record.NewDate(2000, 1, 30): rec2,
		}

		result, err := records.ConvertMinorsApproximate(record.NewDate(2000, 1, 16), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = (rec1["USD"] + rec2["USD"]) / 2
				rubRate = rec2["RUB"]
			)
			expectedResult := int64(math.Round(amount * float64(usdRate/rubRate)))
			assert.Equal(t, expectedResult, result)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing TO rate in the later record", func(t *testing.T) {
		const (
			amount   = 1199 // 11.99 USD
			rangeLim = 100
		)

		var (
			rec1 = record.Record{"USD": 0.1, "RUB": 0.22}
			rec2 = record.Record{"USD": 1.0}
		)

		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  rec1,
			record.NewDate(2000, 1, 30): rec2,
		}

		result, err := records.ConvertMinorsApproximate(record.NewDate(2000, 1, 16), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = (rec1["USD"] + rec2["USD"]) / 2
				rubRate = rec1["RUB"]
			)
			expectedResult := int64(math.Round(amount * float64(usdRate/rubRate)))
			assert.Equal(t, expectedResult, result)
		}
	})

	t.Run("rangeLim covering both earlier and later records, missing TO rate in the both records", func(t *testing.T) {
		const (
			amount   = 9999
			rangeLim = 100
		)

		var (
			record1 = record.Record{"USD": 0.1, "RUB": 0.22}
			record2 = record.Record{"USD": 1.0, "RUB": 0.123}
		)
		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

		result, err := records.ConvertMinorsApproximate(record.NewDate(2000, 1, 16), amount, "USD", "XXX", rangeLim)
		if assert.Error(t, err) {
			assert.Zero(t, result)
		}
	})

	t.Run("rangeLim covering only earlier record", func(t *testing.T) {
		const (
			amount   = 235 // 2.35
			rangeLim = 1
		)

		var (
			rec1 = record.Record{"USD": 0.1, "RUB": 0.5}
			rec2 = record.Record{"USD": 1.0, "RUB": 0.1}
		)
		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  rec1,
			record.NewDate(2000, 1, 30): rec2,
		}

		result, err := records.ConvertMinorsApproximate(record.NewDate(2000, 1, 2), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = rec1["USD"]
				rubRate = rec1["RUB"]
			)
			expectedResult := int64(math.Round(amount * float64(usdRate/rubRate)))
			assert.Equal(t, expectedResult, result)
		}
	})

	t.Run("rangeLim covering only later record", func(t *testing.T) {
		const (
			amount   = 123 // 1.23
			rangeLim = 1
		)

		var (
			rec1 = record.Record{"USD": 0.1, "RUB": 0.5}
			rec2 = record.Record{"USD": 1.0, "RUB": 0.1}
		)
		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  rec1,
			record.NewDate(2000, 1, 30): rec2,
		}

		result, err := records.ConvertMinorsApproximate(record.NewDate(2000, 1, 29), amount, "USD", "RUB", rangeLim)
		if assert.NoError(t, err) {
			var (
				usdRate = rec2["USD"]
				rubRate = rec2["RUB"]
			)
			expectedResult := int64(math.Round(amount * float64(usdRate/rubRate)))
			assert.Equal(t, expectedResult, result)
		}
	})

	t.Run("rangeLim covering neither earlier nor later record", func(t *testing.T) {
		const (
			amount   = 123123
			rangeLim = 5
		)

		var (
			record1 = record.Record{"USD": 0.1, "RUB": 0.5}
			record2 = record.Record{"USD": 1.0, "RUB": 0.1}
		)
		records := UnorderedRecords{
			record.NewDate(2000, 1, 1):  record1,
			record.NewDate(2000, 1, 30): record2,
		}

		result, err := records.ConvertMinorsApproximate(record.NewDate(2000, 1, 15), amount, "USD", "RUB", rangeLim)
		if assert.ErrorIs(t, err, ErrRateApproximationFailed) {
			assert.Zero(t, result)
		}
	})
}
