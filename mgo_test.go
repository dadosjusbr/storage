package storage

import (
	"context"
	"encoding/json"
	"io/ioutil"
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
	file, err := ioutil.ReadFile("teste.json")
	assert.NoError(t, err)
	employee := []Employee{}
	err = json.Unmarshal([]byte(file), &employee)
	assert.NoError(t, err)

	//Summary := Summary{Count: 2, Wage: {Max: }}
	crawler := Crawler{CrawlerID: "123132", CrawlerVersion: "v.1"}
	col := checkCollection{t: t, filter: bson.D{{Key: "aid", Value: "a"}, {Key: "year", Value: 2019}, {Key: "month", Value: 9}}, value: AgencyMonthlyInfo{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler, Employee: employee}}
	c.C = &col

	assert.NoError(t, c.Store(CrawlingResult{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler, Employees: employee}))
	assert.True(t, col.calledReplaceOne())
}
