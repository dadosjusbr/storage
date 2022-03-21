package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client is composed by mongoDbClient and Cloud5 client (used for backup).
type Client struct {
	Db    *DBClient
	Cloud *CloudClient
}

// NewClient NewClient
func NewClient(db *DBClient, cloud *CloudClient) (*Client, error) {
	c := Client{Db: db, Cloud: cloud}
	if err := c.Db.Connect(); err != nil {
		return nil, err
	}
	return &c, nil
}

// Close Connection with DB
func (c *Client) Close(db *DBClient, cloud *CloudClient) error {
	return c.Db.Disconnect()
}

// GetOPE Connect to db to collect data to build 'Órgao por estado' screen
func (c *Client) GetOPE(Uf string, Year int) ([]Agency, map[string][]AgencyMonthlyInfo, error) {
	ags, agsMR, err := c.Db.GetOPE(Uf, Year)
	if err != nil {
		return nil, nil, fmt.Errorf("GetOPE() error: %q", err)
	}
	return ags, agsMR, err
}

// GetOMA Connect to db to collect data for a month including all employees
func (c *Client) GetOMA(month int, year int, agency string) (*AgencyMonthlyInfo, *Agency, error) {
	agsMR, agencyObj, err := c.Db.GetOMA(month, year, agency)
	if err == nil {
		return agsMR, agencyObj, err
	}
	// It is important to let API users know when there no record/doc has been found.
	if err == ErrNothingFound {
		return nil, nil, err
	}
	return nil, nil, fmt.Errorf("error in GetOMA: %q", err)
}

// Store stores the Agency Monthly Info stats.
func (c *Client) Store(agmi AgencyMonthlyInfo) error {
	if c.Db.col == nil {
		return fmt.Errorf("missing collection")
	}
	var err error
	_, err = c.Db.col.ReplaceOne(context.TODO(), bson.D{{Key: "aid", Value: agmi.AgencyID}, {Key: "year", Value: agmi.Year}, {Key: "month", Value: agmi.Month}}, agmi, options.Replace().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("error trying to update mongodb with value {%v}: %q", agmi, err)
	}
	return nil
}

// Store stores an package in the database.
func (c *Client) StorePackage(newPackage Package) error {
	c.Db.Collection(c.Db.packageCol)
	_, err := c.Db.col.InsertOne(context.TODO(),
		bson.D{
			{Key: "aid", Value: newPackage.AgencyID},
			{Key: "group", Value: newPackage.Group},
			{Key: "month", Value: newPackage.Month},
			{Key: "year", Value: newPackage.Year},
			{Key: "package", Value: newPackage.Package}})
	if err != nil {
		return fmt.Errorf("error while storing a new agreggation %q", err)
	}
	return nil
}

// GetAgenciesCount Return the Agencies amount
func (c *Client) GetAgenciesCount() (int64, error) {
	count, err := c.Db.GetAgenciesCount()
	if err != nil {
		return count, fmt.Errorf("GetAgenciesCount() error: %q", err)
	}
	return count, nil
}

// GetNumberOfMonthsCollected Return the Agencies amount
func (c *Client) GetNumberOfMonthsCollected() (int64, error) {
	count, err := c.Db.GetNumberOfMonthsCollected()
	if err != nil {
		return count, fmt.Errorf("GetNumberOfMonthsCollected() error: %q", err)
	}
	return count, nil
}

//GetLastDateWithMonthlyInfo return the latest year and month with collected data
func (c *Client) GetLastDateWithMonthlyInfo() (int, int, error) {
	month, year, err := c.Db.GetLastDateWithMonthlyInfo()
	if err != nil {
		return 0, 0, fmt.Errorf("GetLastDateWithMonthlyInfo() error: %q", err)
	}
	return month, year, nil
}

//GetFirstDateWithMonthlyInfo return the initial year and month with collected data
func (c *Client) GetFirstDateWithMonthlyInfo() (int, int, error) {
	month, year, err := c.Db.GetFirstDateWithMonthlyInfo()
	if err != nil {
		return 0, 0, fmt.Errorf("GetFirstDateWithMonthlyInfo() error: %q", err)
	}
	return month, year, nil
}
