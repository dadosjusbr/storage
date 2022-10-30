package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db       *gorm.DB
	newrelic *newrelic.Application
	user     string
	password string
	dbName   string
	host     string
	port     string
	dsn      string
}

func NewPostgresDB(user, password, dbName, host, port string) (*PostgresDB, error) {
	// check if parameters are not empty
	if user == "" {
		return nil, fmt.Errorf("user cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}
	if dbName == "" {
		return nil, fmt.Errorf("dbName cannot be empty")
	}
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}
	if port == "" {
		return nil, fmt.Errorf("port cannot be empty")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, password)

	return &PostgresDB{
		user:     user,
		password: password,
		dbName:   dbName,
		host:     host,
		port:     port,
		dsn:      dsn,
	}, nil
}

func (p *PostgresDB) Connect() error {
	conn, err := sql.Open("nrpostgres", p.dsn)
	if err != nil {
		panic(err)
	}
	ctx, canc := context.WithTimeout(context.Background(), 30*time.Second)
	defer canc()
	if err := conn.PingContext(ctx); err != nil {
		return fmt.Errorf("error connecting to postgres (creds:%s):%q", p.dsn, err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}))
	if err != nil {
		return fmt.Errorf("error initializing gorm: %q", err)
	}
	p.db = db
	return nil
}

func (p *PostgresDB) Disconnect() error {
	db, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("error returning sql DB: %q", err)
	}
	err = db.Close()
	if err != nil {
		return fmt.Errorf("error closing DB connection: %q", err)
	}
	return nil
}

func (p *PostgresDB) Store(agmi Coleta) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&agmi).Error; err != nil {
			return fmt.Errorf("error inserting 'coleta': %q", err)
		}

		if err := tx.Model(&Coleta{}).Update("atual", false).Where("id = ?", agmi.ID).Error; err != nil {
			return fmt.Errorf("error seting 'atual' to false: %q", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error performing transaction: %q", err)
	}

	return nil
}

func (p *PostgresDB) StorePackage(newPackage Package) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetOPE(uf string, year int) ([]Orgao, map[string][]Coleta, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetAgenciesCount() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetNumberOfMonthsCollected() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetAgencies(uf string) ([]Orgao, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetAgency(aid string) (*Orgao, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetAllAgencies() ([]Orgao, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetMonthlyInfo(agencies []Orgao, year int) (map[string][]Coleta, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetMonthlyInfoSummary(agencies []Orgao, year int) (map[string][]Coleta, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetOMA(month int, year int, agency string) (*Coleta, *Orgao, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetGeneralMonthlyInfosFromYear(year int) ([]GeneralMonthlyInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetFirstDateWithMonthlyInfo() (int, int, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetLastDateWithMonthlyInfo() (int, int, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetRemunerationSummary() (*RemmunerationSummary, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetPackage(pkgOpts PackageFilterOpts) (*Package, error) {
	//TODO implement me
	panic("implement me")
}
