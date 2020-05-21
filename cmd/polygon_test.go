package cmd

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func MockPolygonTickers(file string) (*PolygonTickers, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	p := &PolygonTickers{}
	jsonFile, _ := ioutil.ReadFile(file)
	err := json.Unmarshal(jsonFile, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func Test_PolygonStocks(t *testing.T) {
	dgraphClient := NewDgraphClient(":9080")
	p, err := MockPolygonTickers("PolygonStock.json")
	assert.NoError(t, err)

	err = p.Save(dgraphClient)
	assert.NoError(t, err)
}
