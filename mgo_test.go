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

func makePointer(x float64) *float64 {
	return &x
}

var (
	emp2Row = []Employee{
		Employee{
			Reg:       "",
			Name:      "Abiaci De Carvalho Silva",
			Role:      "Inativo",
			Type:      "",
			Workplace: "inativo",
			Active:    false,
			Income: &IncomeDetails{
				Total: 30368.59,
				Wage:  makePointer(7000),
				Perks: &Perks{
					Total: 600,
				},
				Other: &Funds{
					Total:            100,
					PersonalBenefits: makePointer(7475.71),
					EventualBenefits: makePointer(0),
					PositionOfTrust:  makePointer(5990.88),
					Gratification:    makePointer(0),
					OriginPosition:   makePointer(0),
				},
			},
			Discounts: &Discount{
				Total:            8930.05,
				PrevContribution: makePointer(2719.5),
				CeilRetention:    makePointer(0),
				IncomeTax:        makePointer(6210.55),
				Others: map[string]float64{
					"Sundry": 0,
				},
			},
		},
		Employee{
			Reg:       "",
			Name:      "Abraao Falcao De Carvalho",
			Role:      "Promotor Eleitoral",
			Type:      "",
			Workplace: "10ª zona eleitoral - guarabira/pb",
			Active:    true,
			Income: &IncomeDetails{
				Total: 10000,
				Wage:  makePointer(5000),
				Perks: &Perks{
					Total: 200,
				},
				Other: &Funds{
					Total:            500,
					PersonalBenefits: makePointer(0),
					EventualBenefits: makePointer(0),
					PositionOfTrust:  makePointer(4631.61),
					Gratification:    makePointer(0),
					OriginPosition:   makePointer(0),
				},
			},
			Discounts: &Discount{
				Total:            405.98,
				PrevContribution: makePointer(0),
				CeilRetention:    makePointer(0),
				IncomeTax:        makePointer(405.98),
				Others: map[string]float64{
					"Sundry": 0,
				},
			},
		},
	}

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
		//Test if everything is OK!
		{name: "ok", col: &col, cr: cr, wantErr: false, wantReplaceOne: true},
		// Test if the replaceOne function works!
		{name: "replaceOne error", col: &colErr, cr: cr, wantErr: true, wantReplaceOne: true},
		// Check if has some connection if mongoDb, if does the collection wont be nil.
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

func Test_summary(t *testing.T) {
	emp1Row := emp2Row[:1]
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