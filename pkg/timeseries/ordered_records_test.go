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
	var (
		rec = record.NewWithDate(record.Record{"USD": 0.9, "RUB": 0.01}, record.NewDate(2000, 1, 1))
	)

	records := OrderedRecords{rec}
	assert.Equal(t, []record.WithDate{rec}, records.Slice())
}

func TestOrderedRecords_Rates(t *testing.T) {

}

func TestOrderedRecords_Rate(t *testing.T) {

}

func TestOrderedRecords_ApproximateRates(t *testing.T) {

}

func TestOrderedRecords_ApproximateRate(t *testing.T) {

}

func TestOrderedRecords_Convert(t *testing.T) {

}

func TestOrderedRecords_ConvertApproximate(t *testing.T) {

}
