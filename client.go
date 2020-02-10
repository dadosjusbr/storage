package storage

import (
	"context"
	"fmt"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client is composed by mongoDbClient and Cloud5 client (used for backup).
type Client struct {
	Db *DBClient
	Bc *BackupClient
}

// NewClient NewClient
func NewClient(db *DBClient, bc *BackupClient) (*Client, error) {
	c := Client{Db: db, Bc: bc}
	err := c.Db.Connect()
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetDataForFirstScreen Connect to db to collect data to build first screen
func (c *Client) GetDataForFirstScreen(Uf string, Year int) ([]Agency, map[string][]AgencyMonthlyInfo, error) {
	err := c.Db.Connect()
	//if err != nil {
	//	return nil, nil, fmt.Errorf("GetDataForFirstScreen() error: Unable to connect to DB")
	//}
	ags, agsMR, err := c.Db.GetDataForFirstScreen(Uf, Year)
	c.Db.Disconnect()
	return ags, agsMR, err
}

// GetDataForSecondScreen Connect to db to collect data for a month including all employees
func (c *Client) GetDataForSecondScreen(month int, year int, agency string) (AgencyMonthlyInfo, error) {
	err := c.Db.Connect()
	//if err != nil {
	//	return nil, fmt.Errorf("GetDataForSecondScreen() error: Unable to connect to DB")
	//}
	agsMR, err := c.Db.GetDataForSecondScreen(month, year, agency)
	c.Db.Disconnect()
	return agsMR, err
}

// Store processes and stores the crawling results.
func (c *Client) Store(cr CrawlingResult) error {
	if c.Db.col == nil {
		return fmt.Errorf("missing collection")
	}
	summary := summary(cr.Employees)
	backup, err := c.Bc.backup(cr.Files)
	if err != nil {
		return fmt.Errorf("error trying to get Backup files: %v, error: %q", cr.Files, err)
	}
	agmi := AgencyMonthlyInfo{AgencyID: cr.AgencyID, Month: cr.Month, Year: cr.Year, Crawler: cr.Crawler, Employee: cr.Employees, Summary: summary, Backups: backup}
	_, err = c.Db.col.ReplaceOne(context.TODO(), bson.D{{Key: "aid", Value: cr.AgencyID}, {Key: "year", Value: cr.Year}, {Key: "month", Value: cr.Month}}, agmi, options.Replace().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("error trying to update mongodb with value {%v}: %q", agmi, err)
	}
	return nil
}

// summary aux func to make all necessary calculations to DataSummary Struct
func summary(Employees []Employee) Summary {
	wage := DataSummary{Min: math.MaxFloat64}
	perks := DataSummary{Min: math.MaxFloat64}
	others := DataSummary{Min: math.MaxFloat64}
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
