package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName             = "db"
	monthlyInfoColName = "mi"
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

	summary := Summary(cr.Employees)
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

func Summary(Employees []Employee) (int, int) {
    max := 0
	min := 3.40282346638528859811704183484516925440e+38
	count := len(array)
	total = 0
    for _, value := range array {
        if max < value {
            max = value
        }
        if min > value {
            min = value
		}
		total = total + value
    }
    return min, max
}
