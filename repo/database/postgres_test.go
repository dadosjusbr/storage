package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dadosjusbr/storage/models"
	"github.com/dadosjusbr/storage/repo/database/dto"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var postgresDb *PostgresDB

func TestMain(m *testing.M) {
	if err := getDbTestConnection(); err != nil {
		panic(err)
	}
	exitValue := m.Run()
	postgresDb.Disconnect()
	os.Exit(exitValue)
}

func TestGetOPE(t *testing.T) {
	tests := getOPE{}
	t.Run("Test GetOPE when agencies exists", tests.testWhenAgenciesExists)
	t.Run("Test GetOPE when UF not exists", tests.testWhenUFNotExists)
	t.Run("Test GetOPE when UF is in lower case", tests.testWhenUFIsInLowerCase)
}

type getOPE struct{}

func (g getOPE) testWhenAgenciesExists(t *testing.T) {
	agencies, err := g.insertAgencies()
	if err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgencies, err := postgresDb.GetOPE("SP")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
	truncateAgencies()
}

func (g getOPE) testWhenUFNotExists(t *testing.T) {
	truncateAgencies()

	returnedAgencies, err := postgresDb.GetOPE("SP")

	assert.Nil(t, err)
	assert.Empty(t, returnedAgencies)
}

func (g getOPE) testWhenUFIsInLowerCase(t *testing.T) {
	agencies, err := g.insertAgencies()
	if err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgencies, err := postgresDb.GetOPE("sp")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
}

func (getOPE) insertAgencies() ([]models.Agency, error) {
	agencies := []models.Agency{
		{
			ID:     "tjsp",
			Name:   "Tribunal de Justiça do Estado de São Paulo",
			Type:   "Estadual",
			Entity: "Tribunal",
			UF:     "SP",
		},
		{
			ID:     "mpsp",
			Name:   "Ministério Público do Estado de São Paulo",
			Type:   "Estadual",
			Entity: "Ministério",
			UF:     "SP",
		},
	}
	for _, agency := range agencies {
		agencyDto, err := dto.NewAgencyDTO(agency)
		if err != nil {
			return nil, fmt.Errorf("error creating agency dto %s: %q", agency.ID, err)
		}
		tx := postgresDb.db.Model(dto.AgencyDTO{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(agencyDto)
		if tx.Error != nil {
			return nil, fmt.Errorf("error inserting agency %s: %q", agency.ID, tx.Error)
		}
	}
	return agencies, nil
}

func truncateAgencies() error {
	tx := postgresDb.db.Exec(`TRUNCATE TABLE "coletas", "remuneracoes_zips","orgaos"`)
	if tx.Error != nil {
		return fmt.Errorf("error truncating agencies: %q", tx.Error)
	}
	return nil
}

func getDbTestConnection() error {
	/*Credenciais do banco de dados que serão utilizadas nos testes. Esse é o
	formato que o GoORM utiliza para se conectar ao banco de dados. É importante
	que os valores dessas credenciais sejam iguais as que estão no Dockerfile*/
	credentials := "port=5432 user=dadosjusbr_test dbname=dadosjusbr_test password=dadosjusbr_test sslmode=disable"
	db, err := sql.Open("postgres", credentials)
	if err != nil {
		panic(err)
	}
	ctx, canc := context.WithTimeout(context.Background(), 30*time.Second)
	defer canc()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("error connecting to postgres (creds:%s):%q", credentials, err)
	}
	gormDb, err := gorm.Open(pgdriver.New(pgdriver.Config{
		Conn: db,
	}))
	if err != nil {
		return fmt.Errorf("error initializing gorm: %q", err)
	}
	conn := gormDb
	postgresDb = &PostgresDB{}
	postgresDb.SetConnection(conn)
	return nil
}
