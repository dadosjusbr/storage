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
	"gorm.io/driver/postgres"
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
	agencies := []models.Agency{
		{
			ID:   "tjsp",
			Type: "Estadual",
			UF:   "SP",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	returnedAgencies, err := postgresDb.GetOPE("SP")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
	truncateTables()
}

func (g getOPE) testWhenUFNotExists(t *testing.T) {
	truncateTables()

	returnedAgencies, err := postgresDb.GetOPE("SP")

	assert.Nil(t, err)
	assert.Empty(t, returnedAgencies)
}

func (g getOPE) testWhenUFIsInLowerCase(t *testing.T) {
	agencies := []models.Agency{
		{
			ID:   "tjsp",
			Type: "Estadual",
			UF:   "SP",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgencies, err := postgresDb.GetOPE("sp")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
	truncateTables()
}

func TestGetOPJ(t *testing.T) {
	tests := getOPJ{}
	t.Run("Test GetOPJ when agencies exists", tests.testWhenAgenciesExists)
	t.Run("Test GetOPJ when group not exists", tests.testWhenGroupNotExists)
	t.Run("Test GetOPJ when Group is in irregular case", tests.testWhenGroupIsInIrregularCase)
}

type getOPJ struct{}

func (g getOPJ) testWhenAgenciesExists(t *testing.T) {
	agencies := []models.Agency{
		{
			ID:   "tjsp",
			Type: "Estadual",
		},
		{
			ID:   "tjal",
			Type: "Estadual",
		},
		{
			ID:   "tjba",
			Type: "Estadual",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgencies, err := postgresDb.GetOPJ("Estadual")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
	truncateTables()
}

func (g getOPJ) testWhenGroupNotExists(t *testing.T) {
	truncateTables()

	returnedAgencies, err := postgresDb.GetOPJ("Estadual")

	assert.Nil(t, err)
	assert.Empty(t, returnedAgencies)
}

func (g getOPJ) testWhenGroupIsInIrregularCase(t *testing.T) {
	agencies := []models.Agency{
		{
			ID:     "tjsp",
			Name:   "Tribunal de Justiça do Estado de São Paulo",
			Type:   "Estadual",
			Entity: "Tribunal",
			UF:     "SP",
		},
		{
			ID:     "tjal",
			Name:   "Tribunal de Justiça do Estado de Alagoas",
			Type:   "Estadual",
			Entity: "Tribunal",
			UF:     "AL",
		},
		{
			ID:     "tjba",
			Name:   "Tribunal de Justiça do Estado da Bahia",
			Type:   "Estadual",
			Entity: "Tribunal",
			UF:     "BA",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgencies, err := postgresDb.GetOPJ("eStAdUaL")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
	truncateTables()
}

func insertAgencies(agencies []models.Agency) error {
	for _, agency := range agencies {
		agencyDto, err := dto.NewAgencyDTO(agency)
		if err != nil {
			return fmt.Errorf("error creating agency dto %s: %q", agency.ID, err)
		}
		tx := postgresDb.db.Model(dto.AgencyDTO{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(agencyDto)
		if tx.Error != nil {
			return fmt.Errorf("error inserting agency %s: %q", agency.ID, tx.Error)
		}
	}
	return nil
}

func truncateTables() error {
	tx := postgresDb.db.Exec(`TRUNCATE TABLE coletas, remuneracoes_zips,orgaos CASCADE`)
	if tx.Error != nil {
		return fmt.Errorf("error truncating agencies: %q", tx.Error)
	}
	return nil
}

func getDbTestConnection() error {
	godotenv.Load(".env.test")
	/*Url do banco de dados que será utilizada nos testes. É importante
	que os valores das credenciais dessa Url sejam iguais as que estão no Dockerfile.
	Formato da URL: postgres://{usuario}:{senha}@{host}:{porta}/{banco_de_dados}?sslmode=disable*/
	url := os.Getenv("POSTGRES_CONNECTION_URL")

	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	ctx, canc := context.WithTimeout(context.Background(), 30*time.Second)
	defer canc()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("error connecting to postgres (creds:%s):%q", url, err)
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	if err != nil {
		return fmt.Errorf("error initializing gorm (creds: %s): %q", url, err)
	}
	postgresDb = &PostgresDB{}
	postgresDb.SetConnection(gormDb)
	return nil
}
