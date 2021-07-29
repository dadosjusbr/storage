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
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
}

//the GeneralMonthlyInfo is used to struct the agregation used to get the remuneration info from all angencies in a given month
type GeneralMonthlyInfo struct {
	Month  int     `json:"_id,omitempty" bson:"_id,omitempty"`
	Wage   float64 `json:"wage,omitempty" bson:"wage,omitempty"`
	Perks  float64 `json:"perks,omitempty" bson:"perks,omitempty"`
	Others float64 `json:"others,omitempty" bson:"others,omitempty"`
	Count  int     `json:"count,omitempty" bson:"count,omitempty"`
}

// Errors raised by package storage.
var (
	ErrNothingFound = fmt.Errorf("There is no document with this parameters")
)

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

// GetOPE return agmi info to build first screen
func (c *DBClient) GetOPE(uf string, year int) ([]Agency, map[string][]AgencyMonthlyInfo, error) {
	allAgencies, err := c.GetAgencies(uf)
	if err != nil {
		return nil, nil, fmt.Errorf("GetOPE() error: %q", err)
	}
	result, err := c.GetMonthlyInfo(allAgencies, year)
	if err != nil {
		return nil, nil, fmt.Errorf("GetOPE() error: %q", err)
	}
	return allAgencies, result, nil
}

// GetAgenciesCount Return the Agencies amount
func (c *DBClient) GetAgenciesCount() (int64, error) {
	c.Collection(c.agencyCol)
	itemCount, err := c.col.CountDocuments(context.TODO(), bson.D{}, nil)
	if err != nil {
		return itemCount, fmt.Errorf("Error in result %v", err)
	}
	return itemCount, nil
}

// GetNumberOfMonthsCollected Return the number of months collected
func (c *DBClient) GetNumberOfMonthsCollected() (int64, error) {
	c.Collection(c.monthlyInfoCol)
	itemCount, err := c.col.CountDocuments(context.TODO(), bson.D{}, nil)
	if err != nil {
		return itemCount, fmt.Errorf("Error in result %v", err)
	}
	return itemCount, nil
}

//GetAgencies Return UF Agencies
func (c *DBClient) GetAgencies(uf string) ([]Agency, error) {
	c.Collection(c.agencyCol)
	resultAgencies, err := c.col.Find(context.TODO(), bson.D{{Key: "uf", Value: uf}}, nil)
	if err != nil {
		return nil, fmt.Errorf("Find error in getAgencies %v", err)
	}
	var allAgencies []Agency
	resultAgencies.All(context.TODO(), &allAgencies)
	if err := resultAgencies.Err(); err != nil {
		return nil, fmt.Errorf("Error in result %v", err)
	}
	return allAgencies, nil
}

//GetAgency Return Agency that match ID.
func (c *DBClient) GetAgency(aid string) (*Agency, error) {
	c.Collection(c.agencyCol)
	var Ag Agency
	if err := c.col.FindOne(context.TODO(), bson.D{{Key: "aid", Value: aid}}).Decode(&Ag); err != nil {
		return nil, fmt.Errorf("Error searching for agency id \"%s\":%q", aid, err)
	}
	return &Ag, nil
}

//GetMonthlyInfo return summarized monthlyInfo for each agency in agencies in a specific year
func (c *DBClient) GetMonthlyInfo(agencies []Agency, year int) (map[string][]AgencyMonthlyInfo, error) {
	var result = make(map[string][]AgencyMonthlyInfo)
	c.Collection(c.monthlyInfoCol)
	for _, agency := range agencies {
		resultMonthly, err := c.col.Find(
			context.TODO(), bson.D{{Key: "aid", Value: agency.ID}, {Key: "year", Value: year}},
			options.Find().SetProjection(bson.D{{"aid", 1}, {"year", 1}, {"month", 1}, {"summary", 1}}))
		if err != nil {
			return nil, fmt.Errorf("Error in GetMonthlyInfo %v", err)
		}
		var mr []AgencyMonthlyInfo
		resultMonthly.All(context.TODO(), &mr)
		result[agency.ID] = mr
	}

	return result, nil
}

//GetOMA Search if DB has a match for filters
func (c *DBClient) GetOMA(month int, year int, agency string) (*AgencyMonthlyInfo, *Agency, error) {
	c.Collection(c.monthlyInfoCol)
	var resultMonthly AgencyMonthlyInfo
	err := c.col.FindOne(context.TODO(), bson.D{{Key: "aid", Value: agency}, {Key: "year", Value: year}, {Key: "month", Value: month}}).Decode(&resultMonthly)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		return nil, nil, ErrNothingFound
	}
	agencyObject, err := c.GetAgency(agency)
	if err != nil {
		return nil, nil, fmt.Errorf("Error in GetAgency %v", err)
	}
	return &resultMonthly, agencyObject, nil
}

//Collection Changes active collection
func (c *DBClient) Collection(collectionName string) {
	c.col = c.mgoClient.Database(c.dbName).Collection(collectionName)
}

//GetGeneralMonthlyInfosFromYear return the sum from all remuneration info from all months of a given year
func (c *DBClient) GetGeneralMonthlyInfosFromYear(year int) ([]GeneralMonthlyInfo, error) {
	c.Collection(c.monthlyInfoCol)
	resultMonthly, err := c.col.Aggregate(context.TODO(),
		mongo.Pipeline{
			bson.D{{"$match",
				bson.D{{"year", year}}}},
			bson.D{{"$group",
				bson.D{
					{"_id", "$month"},
					{"wage", bson.D{{"$sum", "$summary.memberactive.wage.total"}}},
					{"perks", bson.D{{"$sum", "$summary.memberactive.perks.total"}}},
					{"others", bson.D{{"$sum", "$summary.memberactive.others.total"}}},
					{"count", bson.D{{"$sum", "$summary.memberactive.count"}}}}}},
			bson.D{{"$sort",
				bson.D{
					{"month", 1}}}}})
	if err != nil {
		return nil, fmt.Errorf("Error in GetMonthlyInfo %v", err)
	}
	var mr []GeneralMonthlyInfo
	resultMonthly.All(context.TODO(), &mr)
	return mr, nil
}

//GetFirstDateWithMonthlyInfo return the initial year and month with collected data
func (c *DBClient) GetFirstDateWithMonthlyInfo() (int, int, error) {
	var resultMonthly AgencyMonthlyInfo
	firstDateQueryOptions := options.FindOne().SetSort(bson.D{{Key: "year", Value: +1}, {Key: "month", Value: +1}})
	err := c.col.FindOne(
		context.TODO(),
		bson.D{}, firstDateQueryOptions).Decode(&resultMonthly)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		return 0, 0, fmt.Errorf("Error in result %v", err)
	}
	return resultMonthly.Month, resultMonthly.Year, nil
}

//GetLastDateWithMonthlyInfo return the latest year and month with collected data
func (c *DBClient) GetLastDateWithMonthlyInfo() (int, int, error) {
	var resultMonthly AgencyMonthlyInfo
	lastDateQueryOptions := options.FindOne().SetSort(bson.D{{Key: "year", Value: -1}, {Key: "month", Value: -1}})
	err := c.col.FindOne(
		context.TODO(),
		bson.D{}, lastDateQueryOptions).Decode(&resultMonthly)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		return 0, 0, fmt.Errorf("Error in result %v", err)
	}
	return resultMonthly.Month, resultMonthly.Year, nil
}

//GetAmountOfRemunerationRecords return the amount of remuneration records from all agencies
func (c *DBClient) GetAmountOfRemunerationRecords() (int, error) {
	c.Collection(c.monthlyInfoCol)
	amountCursor, err := c.col.Aggregate(context.TODO(),
		mongo.Pipeline{bson.D{{"$group", bson.D{{"_id", ""},
			{"amount",
				bson.D{{"$sum", "$summary.memberactive.count"}}}}}}})
	if err != nil {
		return 0, fmt.Errorf("Error in GetAmountOfRemunerationRecords %v", err)
	}
	var result struct {
		Amount int `bson:"amount,omitempty"`
	}
	if amountCursor.Next(context.TODO()) {
		amountCursor.Decode(&result)
	}
	return result.Amount, nil
}

func (c *DBClient) GetGeneralRemunerationValue() (float64, error) {
	c.Collection(c.monthlyInfoCol)
	generalRemuneration, err := c.col.Aggregate(context.TODO(),
		mongo.Pipeline{bson.D{{"$group",
			bson.D{
				{"_id", ""},
				{"wage", bson.D{{"$sum", "$summary.memberactive.wage.total"}}},
				{"perks", bson.D{{"$sum", "$summary.memberactive.perks.total"}}},
				{"others", bson.D{{"$sum", "$summary.memberactive.others.total"}}}}}}})
	if err != nil {
		return 0, fmt.Errorf("Error in GetGeneralRemunerationValue %v", err)
	}
	var result struct {
		Wage   float64 `bson:"wage,omitempty"`
		Perks  float64 `bson:"perks,omitempty"`
		Others float64 `json:"others,omitempty" bson:"others,omitempty"`
	}
	if generalRemuneration.Next(context.TODO()) {
		generalRemuneration.Decode(&result)
	}
	return result.Others + result.Perks + result.Wage, nil
}
