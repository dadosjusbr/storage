package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type checkCollection struct {
	t      *testing.T
	filter bson.D
	value  interface{}
	check  bool
}

func (c *checkCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	assert.Equal(c.t, c.filter, filter)
	assert.Equal(c.t, c.value, replacement)
	c.check = true
	return &mongo.UpdateResult{}, nil
}

func (c *checkCollection) calledReplaceOne() bool {
	return c.check
}

func TestClient_Store(t *testing.T) {
	c, err := NewClient("mongodb://localhost:666")
	assert.NoError(t, err)

	crawler := Crawler{CrawlerID: "123132", CrawlerVersion: "v.1"}
	col := checkCollection{t: t, filter: bson.D{{Key: "aid", Value: "a"}, {Key: "year", Value: 2019}, {Key: "month", Value: 9}}, value: AgencyMonthlyInfo{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler}}
	c.C = &col

	assert.NoError(t, c.Store(CrawlingResult{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler}))
	assert.True(t, col.calledReplaceOne())
}
