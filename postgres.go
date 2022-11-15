package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/dadosjusbr/proto/coleta"
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

func (p *PostgresDB) Store(agmi AgencyMonthlyInfo) error {
	var bkpJson string
	var pkgJson string
	var procInfoJson string
	var summaryJson string

	if !reflect.DeepEqual(&agmi.Backups[0], &Backup{}) {
		bkp, err := json.Marshal(agmi.Backups[0])
		if err != nil {
			return fmt.Errorf("error marshaling backup: %q", err)
		}
		bkpJson = string(bkp)
	}

	if !reflect.DeepEqual(&agmi.Summary, &Summary{}) {
		summary, err := json.Marshal(agmi.Summary)
		if err != nil {
			return fmt.Errorf("error marshaling summary: %q", err)
		}
		summaryJson = string(summary)
	}

	if !reflect.DeepEqual(&agmi.Package, &Package{}) {
		pkg, err := json.Marshal(agmi.Package)
		if err != nil {
			return fmt.Errorf("error marshaling package: %q", err)
		}
		pkgJson = string(pkg)
	}

	if !reflect.DeepEqual(&agmi.ProcInfo, coleta.ProcInfo{}) {
		procInfo, err := json.Marshal(agmi.ProcInfo)
		if err != nil {
			return fmt.Errorf("error marshaling procInfo: %q", err)
		}
		procInfoJson = string(procInfo)
	}

	err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("coletas").Where("id = ?", agmi.ID).Update("atual", false).Error; err != nil {
			return fmt.Errorf("error seting 'atual' to false: %q", err)
		}

		if err := tx.Table("coletas").Create(map[string]interface{}{
			"id":                           agmi.ID,
			"ano":                          agmi.Year,
			"mes":                          agmi.Month,
			"id_orgao":                     agmi.AgencyID,
			"timestamp":                    time.Unix(agmi.CrawlingTimestamp.Seconds, int64(agmi.CrawlingTimestamp.Nanos)),
			"atual":                        true,
			"repositorio_coletor":          agmi.CrawlerRepo,
			"versao_coletor":               agmi.CrawlerVersion,
			"repositorio_parser":           agmi.ParserRepo,
			"versao_parser":                agmi.ParserVersion,
			"estritamente_tabular":         agmi.StrictlyTabular,
			"formato_consistente":          agmi.ConsistentFormat,
			"tem_matricula":                agmi.HaveEnrollment,
			"tem_lotacao":                  agmi.ThereIsACapacity,
			"tem_cargo":                    agmi.HasPosition,
			"acesso":                       agmi.Access,
			"extensao":                     agmi.Extension,
			"formato_aberto":               agmi.OpenFormat,
			"detalhamento_receita_base":    agmi.BaseRevenue,
			"detalhamento_outras_receitas": agmi.OtherRecipes,
			"detalhamento_descontos":       agmi.Expenditure,
			"indice_completude":            agmi.CompletenessScore,
			"indice_facilidade":            agmi.EasinessScore,
			"indice_transparencia":         agmi.Score.Score,
			"backups":                      bkpJson,
			"procinfo":                     procInfoJson,
			"package":                      pkgJson,
			"sumario":                      summaryJson,
		}).Error; err != nil {
			return fmt.Errorf("error inserting 'coleta': %q", err)
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

func (p *PostgresDB) GetOPE(uf string, year int) ([]Agency, map[string][]AgencyMonthlyInfo, error) {
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

func (p *PostgresDB) GetAgencies(uf string) ([]Agency, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetAgency(aid string) (*Agency, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetAllAgencies() ([]Agency, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetMonthlyInfo(agencies []Agency, year int) (map[string][]AgencyMonthlyInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetMonthlyInfoSummary(agencies []Agency, year int) (map[string][]AgencyMonthlyInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetOMA(month int, year int, agency string) (*AgencyMonthlyInfo, *Agency, error) {
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
