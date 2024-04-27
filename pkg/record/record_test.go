package record

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	record := New()
	assert.Empty(t, record)
}

func TestRecord_Rate(t *testing.T) {
	const USDRate = float32(0.9)

	var record = New()
	record["USD"] = USDRate

	t.Run("existing rate", func(t *testing.T) {
		rate, found := record.Rate("USD")
		if assert.True(t, found) {
			assert.Equal(t, USDRate, rate)
		}
	})

	t.Run("non-existent rate", func(t *testing.T) {
		rate, found := record.Rate("XXX")
		if assert.False(t, found) {
			assert.Zero(t, rate)
		}
	})
}

func TestRecord_Convert(t *testing.T) {
	const (
		USDRate = float32(0.9)
		RUBRate = float32(0.01)
	)

	t.Run("convert existing to existing", func(t *testing.T) {
		var record = New()
		record["USD"] = USDRate
		record["RUB"] = RUBRate

		USDtoRUB := map[float32]float32{
			0:   0,
			5:   450,
			-10: -900,
		}
		RUBtoUSD := map[float32]float32{
			0:    0,
			450:  5,
			-900: -10,
		}

		for amount, expectedAmount := range USDtoRUB {
			actualAmount, err := record.Convert(amount, "USD", "RUB")
			if assert.NoError(t, err) {
				assert.Equalf(t, expectedAmount, actualAmount, "amount = %f USD", amount)
			}
		}

		for amount, expectedAmount := range RUBtoUSD {
			actualAmount, err := record.Convert(amount, "RUB", "USD")
			if assert.NoError(t, err) {
				assert.Equalf(t, expectedAmount, actualAmount, "amount = %f RUB", amount)
			}
		}
	})

	t.Run("non-existent to existing", func(t *testing.T) {
		var record = New()
		record["USD"] = 0.9

		result, err := record.Convert(999, "USD", "XXX")
		if assert.ErrorIs(t, err, ErrRateNotFound) {
			assert.Zero(t, result)
		}
	})

	t.Run("existing to non-existent", func(t *testing.T) {
		var record = New()
		record["USD"] = 0.9

		result, err := record.Convert(999, "XXX", "USD")
		if assert.ErrorIs(t, err, ErrRateNotFound) {
			assert.Zero(t, result)
		}
	})
}

func TestRecord_ConvertMinors(t *testing.T) {
	const (
		USDRate = float32(0.9)
		RUBRate = float32(0.01)
	)

	t.Run("convert existing to existing", func(t *testing.T) {
		var record = New()
		record["USD"] = USDRate
		record["RUB"] = RUBRate

		USDtoRUB := map[int64]int64{
			0:     0,
			500:   45000,
			-1000: -90000,
			1:     90,
		}
		RUBtoUSD := map[int64]int64{
			0:      0,
			45000:  500,
			-90000: -1000,
			90:     1,
			1:      0,
		}

		for amount, expectedAmount := range USDtoRUB {
			actualAmount, err := record.ConvertMinors(amount, "USD", "RUB")
			if assert.NoError(t, err) {
				assert.Equalf(t, expectedAmount, actualAmount, "amount = %d USD", amount)
			}
		}

		for amount, expectedAmount := range RUBtoUSD {
			actualAmount, err := record.ConvertMinors(amount, "RUB", "USD")
			if assert.NoError(t, err) {
				assert.Equalf(t, expectedAmount, actualAmount, "amount = %d RUB", amount)
			}
		}
	})

	t.Run("non-existent to existing", func(t *testing.T) {
		var record = New()
		record["USD"] = 0.9

		result, err := record.ConvertMinors(999, "USD", "XXX")
		if assert.ErrorIs(t, err, ErrRateNotFound) {
			assert.Zero(t, result)
		}
	})

	t.Run("existing to non-existent", func(t *testing.T) {
		var record = New()
		record["USD"] = 0.9

		result, err := record.ConvertMinors(999, "XXX", "USD")
		if assert.ErrorIs(t, err, ErrRateNotFound) {
			assert.Zero(t, result)
		}
	})
}
