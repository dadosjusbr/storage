package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	dto "github.com/dadosjusbr/storage/repositories/database/postgres/dto"

	"github.com/dadosjusbr/storage/models"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db       *gorm.DB
	user     string
	password string
	dbName   string
	host     string
	port     string
	uri      string
}

func NewPostgresDB(user, password, dbName, host, port string) (*PostgresDB, error) {
	// Verificando se as credenciais de conexão não estão vazias
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

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, password)
	postgresDB := &PostgresDB{
		user:     user,
		password: password,
		dbName:   dbName,
		host:     host,
		port:     port,
		uri:      uri,
	}
	//Conectando ao postgres
	if err := postgresDB.Connect(); err != nil {
		return nil, fmt.Errorf("error connecting to postgres (creds:%s):%q", uri, err)
	}
	return postgresDB, nil
}

func (p *PostgresDB) Connect() error {
	if p.db != nil {
		return nil
	} else {
		conn, err := sql.Open("nrpostgres", p.uri)
		if err != nil {
			panic(err)
		}
		ctx, canc := context.WithTimeout(context.Background(), 30*time.Second)
		defer canc()
		if err := conn.PingContext(ctx); err != nil {
			return fmt.Errorf("error connecting to postgres (creds:%s):%q", p.uri, err)
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

func (p *PostgresDB) Store(agmi models.AgencyMonthlyInfo) error {
	/*Criando o DTO da coleta a partir de um modelo. É necessário a utilização de 
	DTO's para melhor escalabilidade de bancos de dados. Caso não fosse utilizado,
	não seria possível utilizar outros frameworks/bancos além do GORM, pois ele 
	afeta diretamente os tipos e campos de uma struct.*/
	coletas, err := dto.NewAgencyMonthlyInfoDTO(agmi)
	if err != nil {
		return fmt.Errorf("error converting agency monthly info to dto: %q", err)
	}

	/* Iniciando a transação. É necessário que seja uma transação porque queremos
	executar vários scripts que são dependentes um do outro. Ou seja, se um falhar
	todos falham. Isso nos dá uma maior segurança ao executar a inserção. */
	err = p.db.Transaction(func(tx *gorm.DB) error {
		// Definindo atual como false para todos os registros com o mesmo ID.
		ID := fmt.Sprintf("%s/%d/%d", agmi.AgencyID, agmi.Month, agmi.Year)
		if err := tx.Model(dto.AgencyMonthlyInfoDTO{}).Where("id = ?", ID).Update("atual", false).Error; err != nil {
			return fmt.Errorf("error seting 'atual' to false: %q", err)
		}

		if err := tx.Model(dto.AgencyMonthlyInfoDTO{}).Create(coletas).Error; err != nil {
			return fmt.Errorf("error inserting 'coleta': %q", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error performing transaction: %q", err)
	}

	return nil
}

func (p *PostgresDB) StorePackage(newPackage models.Package) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetOPE(uf string, year int) ([]models.Agency, map[string][]models.AgencyMonthlyInfo, error) {
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

func (p *PostgresDB) GetAgencies(uf string) ([]models.Agency, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetAgency(aid string) (*models.Agency, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetAllAgencies() ([]models.Agency, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetMonthlyInfo(agencies []models.Agency, year int) (map[string][]models.AgencyMonthlyInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetMonthlyInfoSummary(agencies []models.Agency, year int) (map[string][]models.AgencyMonthlyInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetOMA(month int, year int, agency string) (*models.AgencyMonthlyInfo, *models.Agency, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetGeneralMonthlyInfosFromYear(year int) ([]models.GeneralMonthlyInfo, error) {
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

func (p *PostgresDB) GetRemunerationSummary() (*models.RemmunerationSummary, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresDB) GetPackage(pkgOpts models.PackageFilterOpts) (*models.Package, error) {
	//TODO implement me
	panic("implement me")
}