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
	"google.golang.org/protobuf/types/known/timestamppb"
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

func TestStore(t *testing.T) {
	var count int64
	timestamp, _ := time.Parse("2006-01-02 15:04:00.000", "2023-01-16 03:14:17.635") // convertendo string para time.Time
	agmi := models.AgencyMonthlyInfo{
		AgencyID: "mpsp",
		Month:    12,
		Year:     2022,
		Backups: []models.Backup{
			{
				URL:  "https://dadosjusbr-public.s3.amazonaws.com/stf/backups/stf-2022-12.zip",
				Hash: "67e0928dbb026752637ad489bdbf9045",
				Size: 140939,
			},
		},
		Summary: &models.Summary{
			Count: 11,
			BaseRemuneration: models.DataSummary{
				Max:     45710.19,
				Min:     39293.32,
				Average: 43376.78272727273,
				Total:   477144.61000000004,
			},
			OtherRemunerations: models.DataSummary{
				Max:     34727.98,
				Min:     13097.77,
				Average: 17030.535454545454,
				Total:   187335.88999999998,
			},
			IncomeHistogram: map[int]int{-1: 0, 10000: 0, 20000: 0, 30000: 0, 40000: 4, 50000: 7},
		},
		CrawlerVersion:    "sha256:28763548a598f7b2754c770735453bdc94c400d2d923636fb52d64b851a2055d",
		CrawlerRepo:       "https://github.com/dadosjusbr/coletor-stf",
		CrawlingTimestamp: timestamppb.New(timestamp),
		Package: &models.Backup{
			URL:  "https://dadosjusbr-public.s3.amazonaws.com/stf/datapackage/stf-2022-12.zip",
			Hash: "3f500b7d2d99b02ff5f4a4e58a6e04b7",
			Size: 5653,
		},
		Meta: &models.Meta{
			OpenFormat:       true,
			Access:           "ACESSO_DIRETO",
			Extension:        "HTML",
			StrictlyTabular:  true,
			ConsistentFormat: true,
			HaveEnrollment:   true,
			ThereIsACapacity: true,
			HasPosition:      true,
			BaseRevenue:      "DETALHADO",
			OtherRecipes:     "SUMARIZADO",
			Expenditure:      "DETALHADO",
		},
		Score: &models.Score{
			Score:             0.95652174949646,
			CompletenessScore: 0.9166666865348816,
			EasinessScore:     1,
		},
		Duration: 114,
	}
	err := postgresDb.Store(agmi)
	m := postgresDb.db.Model(dto.AgencyMonthlyInfoDTO{}).Where("id = 'mpsp/12/2022' AND atual = true").Count(&count)
	if m.Error != nil {
		fmt.Errorf("error finding agmi: %v", err)
	}
	// Verificando se o método Store deu erro
	assert.Nil(t, err)
	// Verificando se realmente foi armazenado e se tem apenas 1 com atual = true.
	assert.Equal(t, true, count == 1)
}

func truncateAgencies() error {
	tx := postgresDb.db.Exec(`TRUNCATE TABLE "coletas", "remuneracoes_zips","orgaos"`)
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
