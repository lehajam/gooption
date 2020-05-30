package polygon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPolygonStockFeed_Fetch(t *testing.T) {
	feed := NewPolygonStockFeed(":9080")
	err := feed.Fetch()
	assert.NoError(t, err)
}
