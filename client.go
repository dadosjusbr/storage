package storage

import (
	"context"
	"fmt"
	"time"

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
	// TODO: avaliar a necessidade de utilizar transações.

	// Armazenando sempre duas cópias no novo item. Tomamos a decisão de
	// armazenar uma cópia para evitar a complexidade e a perda de desempenho de
	// gerenciar a manutenção de apenas uma cópia entre as coleções (tirar de uma
	// coleção e colocar em outra).

	// ## Armazenando versão corrente
	c.Db.Collection(c.Db.monthlyInfoCol)
	key := bson.D{{Key: "aid", Value: agmi.AgencyID}, {Key: "year", Value: agmi.Year}, {Key: "month", Value: agmi.Month}}
	opts := options.Replace().SetUpsert(true)
	if _, err := c.Db.col.ReplaceOne(context.TODO(), key, agmi, opts); err != nil {
		return fmt.Errorf("error trying to update current monthly info with value {%v}: %q", agmi, err)
	}
	// ## Armazenando revisão.
	c.Db.Collection(c.Db.revCol)
	rev := MonthlyInfoVersion{
		AgencyID:  agmi.AgencyID,
		Month:     agmi.Month,
		Year:      agmi.Year,
		VersionID: time.Now().Unix(),
		Version:   agmi,
	}
	if _, err := c.Db.col.InsertOne(context.TODO(), rev); err != nil {
		return fmt.Errorf("error trying to insert monthly info revision with value {%v}: %q", agmi, err)
	}
	return nil
}

// StorePackage update an package in the database.
func (c *Client) StorePackage(newPackage Package) error {
	c.Db.Collection(c.Db.packageCol)
	filter := bson.M{
		"aid": bson.M{
			"$eq": newPackage.AgencyID,
		},
		"month": bson.M{
			"$eq": newPackage.Month,
		},
		"year": bson.M{
			"$eq": newPackage.Year,
		},
	}
	update := bson.M{
		"$set": bson.M{
			"aid":     newPackage.AgencyID,
			"group":   newPackage.Group,
			"month":   newPackage.Month,
			"year":    newPackage.Year,
			"package": newPackage.Package,
		},
	}
	_, err := c.Db.col.ReplaceOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error while updating a agreggation: %q", err)
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
