package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//collection is a private interface to create a mongo's ReplaceOne method and their signatures to be used and tested.
type collection interface {
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

//DBClient is a mongodb Client instance
type DBClient struct {
	mgoClient      *mongo.Client
	dbName         string
	monthlyInfoCol string
	agencyCol      string
	col            collection
}

//NewDBClient instantiates a mongo new client, but will not connect to the specified URL. Please use Client.Connect before using the client.
func NewDBClient(url, dbName, monthlyInfoCol, agencyCol string) (*DBClient, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	return &DBClient{mgoClient: client, dbName: dbName, monthlyInfoCol: monthlyInfoCol, agencyCol: agencyCol}, nil
}

//Connect establishes a connection to MongoDB using the previously specified URL
func (c *DBClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := c.mgoClient.Connect(ctx); err != nil {
		return fmt.Errorf("error connection with mongo:%q", err)
	}
	return nil
}

//Disconnect closes the connections to MongoDB. It does nothing if the connection had already been closed.
func (c *DBClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	c.col = nil
	return c.mgoClient.Disconnect(ctx)
}

// GetDataForFirstScreen GetDataForFirstScreen
func (c *DBClient) GetDataForFirstScreen(Uf string, Year int) ([]Agency, map[string][]AgencyMonthlyInfo, error) {
	var result = make(map[string][]AgencyMonthlyInfo)
	c.collection(c.agencyCol)
	resultAgencies, _ := c.col.Find(context.TODO(), bson.D{{}}, nil)

	var allAgencies []Agency
	resultAgencies.All(context.TODO(), &allAgencies)
	if err := resultAgencies.Err(); err != nil {
		return nil, nil, fmt.Errorf("Error in result %v", err)
	}

	c.collection(c.monthlyInfoCol)
	findOptions := options.Find()
	for _, agency := range allAgencies {
		resultMonthly, _ := c.col.Find(context.TODO(), bson.D{{Key: "aid", Value: agency.ID}, {Key: "year", Value: Year}},
			findOptions.SetProjection(bson.D{{Key: "aid", Value: ""}, {Key: "year", Value: ""}, {Key: "month", Value: ""}, {Key: "summary", Value: ""}}))
		var mr []AgencyMonthlyInfo
		resultMonthly.All(context.TODO(), &mr)
		result[agency.ID] = mr
	}
	return allAgencies, result, nil
}

func (c *DBClient) collection(collectionName string) {
	c.col = c.mgoClient.Database(c.dbName).Collection(collectionName)
}
