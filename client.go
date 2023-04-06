package storage

import (
	"fmt"

	"github.com/dadosjusbr/storage/models"
	"github.com/dadosjusbr/storage/repo/database"
	"github.com/dadosjusbr/storage/repo/file_storage"
)

// Client is composed by mongoDbClient and Cloud5 client (used for backup).
type Client struct {
	Db    database.Interface
	Cloud file_storage.Interface
}

// NewClient NewClient
func NewClient(db database.Interface, cloud file_storage.Interface) (*Client, error) {
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

// GetStateAgencies Connect to db to collect state agencies by UF
func (c *Client) GetStateAgencies(uf string) ([]models.Agency, error) {
	ags, err := c.Db.GetStateAgencies(uf)
	if err != nil {
		return nil, fmt.Errorf("GetStateAgencies() error: %q", err)
	}
	return ags, err
}

// GetOPJ Connect to db to collect data to build 'Órgao por jurisdição' screen
func (c *Client) GetOPJ(group string) ([]models.Agency, error) {
	ags, err := c.Db.GetOPJ(group)
	if err != nil {
		return nil, fmt.Errorf("GetOPJ() error: %q", err)
	}
	return ags, err
}

// GetOMA Connect to db to collect data for a month including all employees
func (c *Client) GetOMA(month int, year int, agency string) (*models.AgencyMonthlyInfo, *models.Agency, error) {
	agsMR, agencyObj, err := c.Db.GetOMA(month, year, agency)
	if err != nil {
		return nil, nil, fmt.Errorf("GetOMA() error: %q", err)
	}
	return agsMR, agencyObj, nil
}

// Store stores the Agency Monthly Info stats.
func (c *Client) Store(agmi models.AgencyMonthlyInfo) error {
	if err := c.Db.Store(agmi); err != nil {
		return fmt.Errorf("Store() error: %q", err)
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
func (c *Client) GetAgenciesCount() (int, error) {
	count, err := c.Db.GetAgenciesCount()
	if err != nil {
		return count, fmt.Errorf("GetAgenciesCount() error: %q", err)
	}
	return count, nil
}

// GetNumberOfMonthsCollected Return the Agencies amount
func (c *Client) GetNumberOfMonthsCollected() (int, error) {
	count, err := c.Db.GetNumberOfMonthsCollected()
	if err != nil {
		return count, fmt.Errorf("GetNumberOfMonthsCollected() error: %q", err)
	}
	return count, nil
}

// GetLastDateWithMonthlyInfo return the latest year and month with collected data
func (c *Client) GetLastDateWithMonthlyInfo() (int, int, error) {
	month, year, err := c.Db.GetLastDateWithMonthlyInfo()
	if err != nil {
		return 0, 0, fmt.Errorf("GetLastDateWithMonthlyInfo() error: %q", err)
	}
	return month, year, nil
}

// GetFirstDateWithMonthlyInfo return the initial year and month with collected data
func (c *Client) GetFirstDateWithMonthlyInfo() (int, int, error) {
	month, year, err := c.Db.GetFirstDateWithMonthlyInfo()
	if err != nil {
		return 0, 0, fmt.Errorf("GetFirstDateWithMonthlyInfo() error: %q", err)
	}
	return month, year, nil
}

func (c *Client) GetAnnualSummary(agency string) ([]models.AnnualSummary, error) {
	summary, err := c.Db.GetAnnualSummary(agency)
	if err != nil {
		return nil, fmt.Errorf("Error getting annual data from database: %q", err)
	}
	for i := range summary {
		dstKey := fmt.Sprintf("%s/datapackage/%s-%d.zip", agency, agency, summary[i].Year)
		pkg, err := c.Cloud.GetFile(dstKey)
		if err != nil {
			return nil, fmt.Errorf("Error getting annual data from file storage: %q", err)
		}
		summary[i].Package = pkg
	}
	return summary, nil
}

// Get index information by agency's ID or group (name)
func (c *Client) GetIndexInformation(name string, month, year int) (map[string][]models.IndexInformation, error) {
	agg, err := c.Db.GetIndexInformation(name, month, year)
	if err != nil {
		return nil, fmt.Errorf("GetIndexInformation() error: %w", err)
	}
	return agg, nil
}

func (c *Client) GetAllIndexInformation(month, year int) (map[string][]models.IndexInformation, error) {
	agg, err := c.Db.GetAllIndexInformation(month, year)
	if err != nil {
		return nil, fmt.Errorf("GetIndexInformation() error: %w", err)
	}
	return agg, nil
}
