package xml

import (
	"encoding/xml"
)

// Data is a struct representing an XML document containing currency exchange rates records.
type Data struct {
	Cubes []DataCube `xml:"Cube>Cube"`
}

type DataCube struct {
	Date  string         `xml:"time,attr"`
	Rates []DataCubeRate `xml:"Cube"`
}

type DataCubeRate struct {
	Currency string  `xml:"currency,attr"`
	Rate     float32 `xml:"rate,attr"`
}

// NewData decodes XML bytes to a new Data.
func NewData(data []byte) (*Data, error) {
	xmlData := &Data{}
	if err := xml.Unmarshal(data, xmlData); err != nil {
		return nil, err
	}
	return xmlData, nil
}
