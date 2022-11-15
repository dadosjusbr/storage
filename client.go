package storage

import (
	"fmt"
)

//Client is composed by mongoDbClient and Cloud5 client (used for backup).
type Client struct {
	Db    IDatabaseService
	Cloud IStorageService
}

// NewClient NewClient
func NewClient(db IDatabaseService, cloud IStorageService) (*Client, error) {
	c := Client{Db: db, Cloud: cloud}
	if err := c.Db.Connect(); err != nil {
		return nil, err
	}
	return &c, nil
}

// Close Connection with DB
func (c *Client) Close() error {
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
	if err := c.Db.Store(agmi); err != nil {
		return fmt.Errorf("Store() error: %q", err)
	}
	return nil
}

// StorePackage update an package in the database.
func (c *Client) StorePackage(newPackage Package) error {
	if err := c.Db.StorePackage(newPackage); err != nil {
		return fmt.Errorf("StorePackage() error %q", err)
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
