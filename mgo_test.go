package storage

import (
	"context"
	"fmt"
	"reflect"
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
	opts   []*options.ReplaceOptions
	check  bool
	err    bool
}

func (c *checkCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	c.check = true
	if c.err {
		return nil, fmt.Errorf("replace one error")
	}
	assert.Equal(c.t, c.filter, filter)
	assert.Equal(c.t, c.value, replacement)
	assert.Equal(c.t, c.opts, opts)
	return &mongo.UpdateResult{}, nil
}

func (c *checkCollection) calledReplaceOne() bool {
	return c.check
}

func TestClient_Store(t *testing.T) {
	/*
		file, err := ioutil.ReadFile("teste.json")
		assert.NoError(t, err)
		employee := []Employee{}
		err = json.Unmarshal([]byte(file), &employee)
		assert.NoError(t, err)
	*/
	//Summary := Summary{Count: 2, Wage: {Max: }}

	crawler := Crawler{CrawlerID: "123132", CrawlerVersion: "v.1"}
	cr := CrawlingResult{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler}
	col := checkCollection{
		t:      t,
		filter: bson.D{{Key: "aid", Value: "a"}, {Key: "year", Value: 2019}, {Key: "month", Value: 9}},
		value:  AgencyMonthlyInfo{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler},
		opts:   []*options.ReplaceOptions{options.Replace().SetUpsert(true)},
		err:    false,
	}
	colErr := checkCollection{
		t:      t,
		filter: bson.D{{Key: "aid", Value: "a"}, {Key: "year", Value: 2019}, {Key: "month", Value: 9}},
		value:  AgencyMonthlyInfo{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler},
		opts:   []*options.ReplaceOptions{options.Replace().SetUpsert(true)},
		err:    true,
	}

	tests := []struct {
		name           string
		col            *checkCollection
		cr             CrawlingResult
		wantErr        bool
		wantReplaceOne bool
	}{
		{name: "ok", col: &col, cr: cr, wantErr: false, wantReplaceOne: true},
		{name: "replaceOne error", col: &colErr, cr: cr, wantErr: true, wantReplaceOne: true},
		{name: "missing collection error", cr: cr, wantErr: true, wantReplaceOne: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}
			if tt.col != nil {
				c.c = tt.col
			}

			if err := c.Store(tt.cr); (err != nil) != tt.wantErr {
				t.Errorf("Client.Store() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.col != nil && (tt.wantReplaceOne != tt.col.calledReplaceOne()) {
				t.Errorf("Client.Store() error calledReplaceOne != wantReplaceOne")
			}
		})
	}
}

func Test_summary(t *testing.T) {

	tests := []struct {
		name      string
		Employees []Employee
		want      Summary
	}{
		{name: "no employee"},
		{name: "1 employee"},
		{name: "1+ employee"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := summary(tt.Employees); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("summary() = %v, want %v", got, tt.want)
			}
		})
	}
}
