package timeseries

import (
	"github.com/jieggii/ecbratex/pkg/record"
	"github.com/jieggii/ecbratex/pkg/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOrderedUnorderedRecordsFromXML(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		data := &xml.Data{
			Cubes: []xml.DataCube{
				{Date: "2024-12-02", Rates: []xml.DataCubeRate{{"USD", 0.9}}},
			},
		}

		records, err := NewOrderedUnorderedRecordsFromXML(data)
		if assert.NoError(t, err) {
			assert.Len(t, records.UnorderedRecords, 1)
		}
	})

	t.Run("data with invalid date", func(t *testing.T) {
		data := &xml.Data{
			Cubes: []xml.DataCube{
				{Date: "invalid date", Rates: []xml.DataCubeRate{{"USD", 0.9}}},
			},
		}
		records, err := NewOrderedUnorderedRecordsFromXML(data)
		if assert.Error(t, err) {
			assert.Empty(t, records)
		}
	})

	t.Run("nil data", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = NewOrderedUnorderedRecordsFromXML(nil)
		})
	})

	t.Run("empty data", func(t *testing.T) {
		data := &xml.Data{}
		records, err := NewOrderedUnorderedRecordsFromXML(data)
		if assert.NoError(t, err) {
			assert.Empty(t, records.Dates)
			assert.Empty(t, records.UnorderedRecords)
		}
	})
}

func TestOrderedUnorderedRecords_Slice(t *testing.T) {
	var (
		date1 = record.NewDate(2000, 1, 3)
		date2 = record.NewDate(2000, 1, 2)
		date3 = record.NewDate(2000, 1, 1)

		rec1 = record.Record{"USD": 0.7}
		rec2 = record.Record{"USD": 0.8}
		rec3 = record.Record{"USD": 0.9}
	)

	records := OrderedUnorderedRecords{
		Dates: []record.Date{date1, date2, date3},
		UnorderedRecords: UnorderedRecords{
			date1: rec1,
			date2: rec2,
			date3: rec3,
		},
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
