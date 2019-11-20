package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// checkCollection is a struct containing necessary tools to make tests and simulate mongoDb func
type checkCollection struct {
	t      *testing.T
	filter bson.D
	value  interface{}
	opts   []*options.ReplaceOptions
	check  bool
	err    bool
}

// Samples to use for case tests, employee and summary's for 1 row and 1+ rows
var (
	employeeSample1 = []byte(`[{"reg": "","name": "Abiaci De Carvalho Silva","role": "Inativo","type": "","workplace": "inativo","active": false,"income": {
	  "total": 30368.59,"wage": 7000,"perks": {"total": 600,"food": null,"transportation": null,"pre_school": null,"health": null,"birth_aid": null,"housing_aid": null,  
	  "subsistence": null, "others": null}, "other": {"total": 100,"person_benefits": 7475.71,"eventual_benefits": 0,"trust_position": 5990.88,"daily": null,	  
	  "gratific": 0, "origin_pos": 0, "others": null}}, "discounts": {"total": 8930.05,"prev_contribution": 2719.5,"ceil_retention": 0,"income_tax": 6210.55,	  
	  "sundry": {"Sundry": 0}}}]`)

	employeeSample2 = []byte(`[{"reg":"","name":"Abiaci De Carvalho Silva","role":"Inativo","type":"","workplace":"inativo","active":false,"income":
	{"total":30368.59,"wage":7000,"perks":{"total":600,"food":null,"transportation":null,"pre_school":null,"health":null,"birth_aid":null,
	"housing_aid":null,"subsistence":null,"others":null},"other":{"total":100,"person_benefits":7475.71,"eventual_benefits":0,"trust_position":5990.88,
	"daily":null,"gratific":0,"origin_pos":0,"others":null}},"discounts":{"total":8930.05,"prev_contribution":2719.5,"ceil_retention":0,"income_tax":6210.55,
	"sundry":{"Sundry":0}}},{"reg":"","name":"Abraao Falcao De Carvalho","role":"Promotor Eleitoral","type":"",
	"workplace":"10ª zona eleitoral - guarabira/pb","active":true,"income":{"total":10000,"wage":5000,"perks":{"total":200,"food":null,
	"transportation":null,"pre_school":null,"health":null,"birth_aid":null,"housing_aid":null,"subsistence":null,"others":null},
	"other":{"total":500,"person_benefits":0,"eventual_benefits":0,"trust_position":4631.61,"daily":null,"gratific":0,"origin_pos":0,"others":null}},
	"discounts":{"total":405.98,"prev_contribution":0,"ceil_retention":0,"income_tax":405.98,"sundry":{"Sundry":0}}}]
	`)

	summFor1Row = Summary{
		Count:  1,
		Wage:   DataSummary{Max: 7000.00, Min: 7000.00, Average: 7000.00, Total: 7000.00},
		Perks:  DataSummary{Max: 600.00, Min: 600.00, Average: 600.00, Total: 600.00},
		Others: DataSummary{Max: 100.00, Min: 100.00, Average: 100.00, Total: 100.00},
	}
	summFor2Row = Summary{
		Count:  2,
		Wage:   DataSummary{Max: 7000.00, Min: 5000.00, Average: 6000.00, Total: 12000.00},
		Perks:  DataSummary{Max: 600.00, Min: 200.00, Average: 400.00, Total: 800.00},
		Others: DataSummary{Max: 500.00, Min: 100.00, Average: 300.00, Total: 600.00},
	}
)

// ReplaceOne is a checkCollection func that use same signature of collection interface, which is the same as the method signature with the same name in mongoDb
// We assert if filter, value and opts are the same in mongo, and turn into true c.check if ReplaceOne is called
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

// calledReplaceOne is a checkCollection func that returns a bool checking if ReplaceOne was called or not.
func (c *checkCollection) calledReplaceOne() bool {
	return c.check
}

//TestClient_Store test Store func if is everything is ok, if replaceOne method is called and if we can connect into a mongoDb and get collection,
// if we cant get a collection, its because we cant be able to connect with mongo.
func TestClient_Store(t *testing.T) {
	emp2Row := []Employee{}
	err := json.Unmarshal(employeeSample2, &emp2Row)
	assert.NoError(t, err)
	crawler := Crawler{CrawlerID: "123132", CrawlerVersion: "v.1"}
	cr := CrawlingResult{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler, Employees: emp2Row}
	col := checkCollection{
		t:      t,
		filter: bson.D{{Key: "aid", Value: "a"}, {Key: "year", Value: 2019}, {Key: "month", Value: 9}},
		value:  AgencyMonthlyInfo{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler, Employee: emp2Row, Summary: summFor2Row},
		opts:   []*options.ReplaceOptions{options.Replace().SetUpsert(true)},
		err:    false,
	}
	colErr := checkCollection{
		t:      t,
		filter: bson.D{{Key: "aid", Value: "a"}, {Key: "year", Value: 2019}, {Key: "month", Value: 9}},
		value:  AgencyMonthlyInfo{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler, Employee: emp2Row, Summary: summFor2Row},
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
				c.col = tt.col
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

// Test_summary func when we have no rows, one or more rows in employee slice.
func Test_summary(t *testing.T) {
	emp2Row := []Employee{}
	err := json.Unmarshal(employeeSample2, &emp2Row)
	emp1Row := emp2Row[:1]
	assert.NoError(t, err)

	tests := []struct {
		name      string
		Employees []Employee
		want      Summary
	}{
		{name: "no employee"},
		{name: "1 employee", Employees: emp1Row, want: summFor1Row},
		{name: "1+ employee", Employees: emp2Row, want: summFor2Row},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := summary(tt.Employees); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("summary() = %v, want %v", got, tt.want)
			}
		})
	}
}
