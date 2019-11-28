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
	maxValue           = math.MaxFloat64
)

//collection is a private interface to create a mongo's ReplaceOne method and their signatures to be used and tested.
type collection interface {
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
}

//Client is a mongodb Client instance
type Client struct {
	mgoClient *mongo.Client
	col       collection
}

//NewClient instantiates a new client, but will not connect to the specified URL. Please use Client.Connect before using the client.
func NewClient(url string) (*Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	return &Client{mgoClient: client}, nil
}

//Connect establishes a connection to MongoDB using the previously specified URL
func (c *Client) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := c.mgoClient.Connect(ctx); err != nil {
		return fmt.Errorf("error connection with mongo:%q", err)
	}
	c.col = c.mgoClient.Database(dbName).Collection(monthlyInfoColName)
	return nil
}

//Disconnect closes the connections to MongoDB. It does nothing if the connection had already been closed.
func (c *Client) Disconnect() error {
	if c.col == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	c.col = nil
	return c.mgoClient.Disconnect(ctx)
}

// Store processes and stores the crawling results.
func (c *Client) Store(cr CrawlingResult) error {
	if c.col == nil {
		return fmt.Errorf("Client is not connected")
	}
	summary := summary(cr.Employees)
	agmi := AgencyMonthlyInfo{AgencyID: cr.AgencyID, Month: cr.Month, Year: cr.Year, Crawler: cr.Crawler, Employee: cr.Employees, Summary: summary}
	_, err := c.col.ReplaceOne(context.TODO(), bson.D{{Key: "aid", Value: cr.AgencyID}, {Key: "year", Value: cr.Year}, {Key: "month", Value: cr.Month}}, agmi, options.Replace().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("error trying to update mongodb with value {%+v}: %q", agmi, err)
	}
	// armazenar o backup
	return nil
}

// summary aux func to make all necessary calculations to DataSummary Struct
func summary(Employees []Employee) Summary {
	wage := DataSummary{Min: maxValue}
	perks := DataSummary{Min: maxValue}
	others := DataSummary{Min: maxValue}
	count := len(Employees)
	if count == 0 {
		return Summary{}
	}
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
	return Summary{
		Count:  count,
		Wage:   wage,
		Perks:  perks,
		Others: others,
	}
}

func GetAgenciesOfStateByName(stateName string) []Agency {
	return []Agency{
		Agency{
			ID:     "tjpb",
			Name:   "Tribunal de Justiça da Paraíba",
			Type:   "E",
			Entity: "J",
			UF:     "PB",
		},
		Agency{
			ID:     "MPPB",
			Name:   "Ministério Público da Paraíba",
			Type:   "E",
			Entity: "M",
			UF:     "PB",
		},
	}
}
