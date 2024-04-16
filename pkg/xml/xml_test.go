package xml

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const testXMLFilePath = "./../../testdata/eurofxref-daily.xml"

func TestNewData(t *testing.T) {
	t.Run("valid XML bytes", func(t *testing.T) {
		data, err := os.ReadFile(testXMLFilePath)
		if err != nil {
			panic(err)
		}

		xmlData, err := NewData(data)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, xmlData)
		}
	})

	t.Run("invalid XML bytes", func(t *testing.T) {
		var data = []byte("invalid xml data")

		xmlData, err := NewData(data)
		if assert.Error(t, err) {
			assert.Empty(t, xmlData)
		}
	})
}
