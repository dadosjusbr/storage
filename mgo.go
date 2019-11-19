package storage

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName             = "db"
	monthlyInfoColName = "mi"
	maxValue           = 3.40282346638528859811704183484516925440e+38
)

type collection interface {
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
}

type Client struct {
	mgoClient *mongo.Client
	C         collection
}

func NewClient(url string) (*Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	return &Client{mgoClient: client}, nil
}

func (c *Client) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := c.mgoClient.Connect(ctx); err != nil {
		return fmt.Errorf("error connection with mongo:%q", err)
	}
	c.C = c.mgoClient.Database(dbName).Collection(monthlyInfoColName)
	return nil
}

func (c *Client) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return c.mgoClient.Disconnect(ctx)
}

// Store processes and stores the crawling results.
func (c *Client) Store(cr CrawlingResult) error {
	summary := summary(cr.Employees)
	agmi := AgencyMonthlyInfo{AgencyID: cr.AgencyID, Month: cr.Month, Year: cr.Year, Crawler: cr.Crawler, Employee: cr.Employees, Summary: summary}
	_, err := c.C.ReplaceOne(context.TODO(), bson.D{{Key: "aid", Value: cr.AgencyID}, {Key: "year", Value: cr.Year}, {Key: "month", Value: cr.Month}}, agmi, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}
	// armazenar os empregados
	// armazenar o sumário
	// armazenar o backup
	return nil
}

func summary(Employees []Employee) Summary {
	wage := DataSummary{Max: 0.0, Min: maxValue, Total: 0.0}
	perks := DataSummary{Max: 0.0, Min: maxValue, Total: 0.0}
	others := DataSummary{Max: 0.0, Min: maxValue, Total: 0.0}
	fmt.Println(len(Employees))
	count := len(Employees)
	for _, value := range Employees {
		wage.Max = math.Max(wage.Max, *value.Income.Wage)
		perks.Max = math.Max(perks.Max, value.Income.Perks.Total)
		others.Max = math.Max(others.Max, value.Income.Other.Total)
		wage.Min = math.Min(wage.Min, *value.Income.Wage)
		perks.Min = math.Min(perks.Min, value.Income.Perks.Total)
		others.Min = math.Min(others.Min, value.Income.Other.Total)
		wage.Total += *value.Income.Wage
		perks.Total += value.Income.Perks.Total
		others.Total += value.Income.Other.Total
	}
	wage.Average = wage.Total / float64(count)
	perks.Average = perks.Total / float64(count)
	others.Average = others.Total / float64(count)
	summary := Summary{Count: count, Wage: wage, Perks: perks, Others: others}
	return summary
}
