package storage

import (
	"fmt"

	"github.com/dadosjusbr/storage/models"
	"github.com/dadosjusbr/storage/repositories/interfaces"
)

//Client is composed by mongoDbClient and Cloud5 client (used for backup).
type Client struct {
	Db    interfaces.IDatabaseRepository
	Cloud interfaces.IStorageRepository
}

// NewClient NewClient
func NewClient(db interfaces.IDatabaseRepository, cloud interfaces.IStorageRepository) (*Client, error) {
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
func (c *Client) GetOPE(Group string, Uf string, Year int) ([]models.Agency, error) {
	ags, err := c.Db.GetOPE(Group, Uf, Year)
	if err != nil {
		return nil, fmt.Errorf("GetOPE() error: %q", err)
	}
	return ags, err
}

// GetOPT Connect to db to collect data to build 'Órgao por grupo' screen
func (c *Client) GetOPT(Group string, Year int) ([]models.Agency, error) {
	ags, err := c.Db.GetOPT(Group, Year)
	if err != nil {
		return nil, fmt.Errorf("GetOPT() error: %q", err)
	}
	return ags, err
}

// GetOMA Connect to db to collect data for a month including all employees
func (c *Client) GetOMA(month int, year int, agency string) (*models.AgencyMonthlyInfo, *models.Agency, error) {
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
func (c *Client) Store(agmi models.AgencyMonthlyInfo) error {
	if err := c.Db.Store(agmi); err != nil {
		return fmt.Errorf("Store() error: %q", err)
	}
	return nil
}

// StorePackage update an package in the database.
func (c *Client) StorePackage(newPackage models.Package) error {
	if err := c.Db.StorePackage(newPackage); err != nil {
		return fmt.Errorf("StorePackage() error %q", err)
	}
	return nil
}

func (c *Client) StoreRemunerations(remu models.Remunerations) error {
	if err := c.Db.StoreRemunerations(remu); err != nil {
		return fmt.Errorf("StoreRemunerations() error: %q", err)
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
