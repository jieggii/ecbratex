package record

import (
	"github.com/jieggii/ecbratex/pkg/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWithDate(t *testing.T) {
	rec := New()
	date := NewDate(2024, 12, 1)

	recWithDate := NewWithDate(rec, date)
	assert.Equal(t, rec, recWithDate.Record)
	assert.Equal(t, date, recWithDate.Date)
}

func TestNewWithDateFromXMLData(t *testing.T) {
	const (
		USDRate float32 = 123.1
		EURRate float32 = 1
	)

	t.Run("valid XML data", func(t *testing.T) {
		xmlData := &xml.Data{
			Cubes: []xml.DataCube{{
				Date: "2024-12-01",
				Rates: []xml.DataCubeRate{
					{Currency: "USD", Rate: USDRate},
				},
			}},
		}
		recWithDate, err := NewWithDateFromXMLData(xmlData)
		if assert.NoError(t, err) {
			assert.Equal(t, NewDate(2024, 12, 1), recWithDate.Date)
			assert.Equal(t, Record(map[string]float32{"EUR": EURRate, "USD": USDRate}), recWithDate.Record)
		}
	})

	t.Run("invalid date in XML data", func(t *testing.T) {
		xmlData := &xml.Data{
			Cubes: []xml.DataCube{{
				Date: "invalid date!",
				Rates: []xml.DataCubeRate{
					{Currency: "USD", Rate: USDRate},
				},
			}},
		}
		recWithDate, err := NewWithDateFromXMLData(xmlData)
		if assert.Error(t, err) {
			assert.Empty(t, recWithDate)
		}
	})

	t.Run("no rate records in XML data", func(t *testing.T) {
		xmlData := &xml.Data{
			Cubes: []xml.DataCube{},
		}
		recWithDate, err := NewWithDateFromXMLData(xmlData)
		if assert.Error(t, err) {
			assert.Empty(t, recWithDate)
		}

	})
}
