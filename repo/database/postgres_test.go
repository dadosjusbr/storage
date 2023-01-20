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
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var conn *gorm.DB
var postgresDb *PostgresDB

func TestGetOPE(t *testing.T) {
	if err := getDbTestConnection(); err != nil {
		t.Fatal(err)
	}
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
		tx := conn.Model(dto.AgencyDTO{}).Clauses(clause.OnConflict{
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
	tx := conn.Exec(`TRUNCATE TABLE "coletas", "remuneracoes_zips","orgaos"`)
	if tx.Error != nil {
		return fmt.Errorf("error truncating agencies: %q", tx.Error)
	}
	return nil
}

func getDbTestConnection() error {
	godotenv.Load("../../.env")
	uri := os.Getenv("POSTGRES_TEST_URL")
	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	ctx, canc := context.WithTimeout(context.Background(), 30*time.Second)
	defer canc()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("error connecting to postgres (creds:%s):%q", uri, err)
	}
	gormDb, err := gorm.Open(pgdriver.New(pgdriver.Config{
		Conn: db,
	}))
	if err != nil {
		return fmt.Errorf("error initializing gorm: %q", err)
	}
	conn = gormDb
	postgresDb = &PostgresDB{}
	postgresDb.SetConnection(conn)
	return nil
}
