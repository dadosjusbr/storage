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
	if err := c.Db.Connect(); err != nil {
		return nil, err
	}
	return &c, nil
}

// GetDataForFirstScreen Connect to db to collect data to build first screen
func (c *Client) GetDataForFirstScreen(Uf string, Year int) ([]Agency, map[string][]AgencyMonthlyInfo, error) {
	ags, agsMR, err := c.Db.GetDataForFirstScreen(Uf, Year)
	if err != nil {
		return nil, nil, fmt.Errorf("GetDataForFirstScreen() error: %q", err)
	}
	c.Db.Disconnect()
	return ags, agsMR, err
}

// GetDataForSecondScreen Connect to db to collect data for a month including all employees
func (c *Client) GetDataForSecondScreen(month int, year int, agency string) (*AgencyMonthlyInfo, error) {
	agsMR, err := c.Db.GetDataForSecondScreen(month, year, agency)
	if err != nil {
		return nil, fmt.Errorf("GetDataForSecondScreen() error: %q", err)
	}
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
	agmi := AgencyMonthlyInfo{AgencyID: cr.AgencyID, Month: cr.Month, Year: cr.Year, Crawler: cr.Crawler, Employee: cr.Employees, Summary: summary, Backups: backup, CrawlingTimestamp: cr.Timestamp}
	_, err = c.Db.col.ReplaceOne(context.TODO(), bson.D{{Key: "aid", Value: cr.AgencyID}, {Key: "year", Value: cr.Year}, {Key: "month", Value: cr.Month}}, agmi, options.Replace().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("error trying to update mongodb with value {%v}: %q", agmi, err)
	}
	return nil
}

// summary aux func to make all necessary calculations to DataSummary Struct
func summary(Employees []Employee) Summaries {
	general := createSummary()
	memberA := createSummary()
	memberI := createSummary()
	serverA := createSummary()
	serverI := createSummary()
	for _, emp := range Employees {
		updateSummary(&general, emp)
		switch {
		case emp.Type == "membro" && emp.Active:
			updateSummary(&memberA, emp)
		case emp.Type == "membro" && !emp.Active:
			updateSummary(&memberI, emp)
		case emp.Type == "servidor" && emp.Active:
			updateSummary(&serverA, emp)
		case emp.Type == "servidor" && !emp.Active:
			updateSummary(&serverI, emp)
		}
	}

	if general.Count == 0 {
		return Summaries{}
	}
	checkCountIsO(&memberA)
	checkCountIsO(&memberI)
	checkCountIsO(&serverA)
	checkCountIsO(&serverI)

	return Summaries{
		General: general,
		MemberA: memberA,
		MemberI: memberI,
		ServerA: serverA,
		ServerI: serverI,
	}
}

//checkCOuntIsO check if the number of employees is 0 of each employee type
func checkCountIsO(s *Summary) {
	if s.Count == 0 {
		s.Wage.Min = 0
		s.Others.Min = 0
		s.Perks.Min = 0
	}
}

//updateDataSummary auxiliary function that updates the summary data at each employee value
func updateDataSummary(d *DataSummary, value float64, count int) {
	d.Max = math.Max(d.Max, value)
	d.Min = math.Min(d.Min, value)
	d.Total += value
	d.Average = d.Total / float64(count)
}

//updateSummary count number of employees
func updateSummary(s *Summary, emp Employee) {
	s.Count++
	updateDataSummary(&s.Wage, *emp.Income.Wage, s.Count)
	updateDataSummary(&s.Perks, emp.Income.Perks.Total, s.Count)
	updateDataSummary(&s.Others, emp.Income.Other.Total, s.Count)
}

//createSummary instanciates all employee group Type plus Active fields.
func createSummary() Summary {
	return Summary{
		Count:  0,
		Wage:   DataSummary{Min: math.MaxFloat64},
		Perks:  DataSummary{Min: math.MaxFloat64},
		Others: DataSummary{Min: math.MaxFloat64},
	}
}
