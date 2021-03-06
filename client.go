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
	if err != nil {
		return nil, nil, fmt.Errorf("GetOMA() error: %q", err)
	}
	return agsMR, agencyObj, err
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
