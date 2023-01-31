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
	truncateTables()
	exitValue := m.Run()
	truncateTables()
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
			ID:     "tjsp",
			Name:   "Tribunal de Justiça do Estado de São Paulo",
			Type:   "Estadual",
			Entity: "Tribunal",
			UF:     "SP",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgencies, err := postgresDb.GetOPE("SP")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
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
			ID:     "tjsp",
			Name:   "Tribunal de Justiça do Estado de São Paulo",
			Type:   "Estadual",
			Entity: "Tribunal",
			UF:     "SP",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgencies, err := postgresDb.GetOPE("sp")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
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

	returnedAgencies, err := postgresDb.GetOPJ("Estadual")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
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
}

func TestGetAgenciesCount(t *testing.T) {
	tests := getAgenciesCount{}
	t.Run("Test GetAgenciesCount when agencies exists", tests.testWhenAgenciesExists)
	t.Run("Test GetAgenciesCount when agencies not exists", tests.testWhenAgenciesNotExists)
}

type getAgenciesCount struct{}

func (g getAgenciesCount) testWhenAgenciesExists(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjsp",
		},
		{
			ID: "tjal",
		},
		{
			ID: "tjba",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	count, err := postgresDb.GetAgenciesCount()

	assert.Nil(t, err)
	assert.Equal(t, len(agencies), count)
}

func (g getAgenciesCount) testWhenAgenciesNotExists(t *testing.T) {
	truncateTables()

	count, err := postgresDb.GetAgenciesCount()

	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}

func TestGetNumberOfMonthsCollected(t *testing.T) {
	tests := getNumberOfMonthsCollected{}
	t.Run("Test GetNumberOfMonthsCollected when monthly infos exists", tests.testWhenMonthlyInfosExists)
	t.Run("Test GetNumberOfMonthsCollected when monthly infos not exists", tests.testWhenMonthlyInfosNotExists)
}

type getNumberOfMonthsCollected struct{}

func (g getNumberOfMonthsCollected) testWhenMonthlyInfosExists(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjsp",
		},
		{
			ID: "tjal",
		},
		{
			ID: "tjba",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	monthlyInfos := []models.AgencyMonthlyInfo{
		{
			AgencyID: "tjsp",
			Year:     2020,
			Month:    1,
		},
		{
			AgencyID: "tjal",
			Year:     2020,
			Month:    2,
		},
		{
			AgencyID: "tjba",
			Year:     2020,
			Month:    3,
		},
	}
	if err := insertMonthlyInfos(monthlyInfos); err != nil {
		t.Fatalf("error inserting monthly infos: %q", err)
	}

	count, err := postgresDb.GetNumberOfMonthsCollected()

	assert.Nil(t, err)
	assert.Equal(t, len(monthlyInfos), count)
}

func (g getNumberOfMonthsCollected) testWhenMonthlyInfosNotExists(t *testing.T) {
	truncateTables()

	count, err := postgresDb.GetNumberOfMonthsCollected()

	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}

func TestGetAgency(t *testing.T) {
	tests := getAgency{}
	t.Run("Test GetAgency when agency exists", tests.testWhenAgencyExists)
	t.Run("Test GetAgency when agency not exists", tests.testWhenAgencyNotExists)
	t.Run("Test GetAgency when agency is in irregular case", tests.testWhenAgencyIsInIrregularCase)
}

type getAgency struct{}

func (g getAgency) testWhenAgencyExists(t *testing.T) {
	agencies := []models.Agency{
		{ID: "tjsp",
			Name:   "Tribunal de Justiça do Estado de São Paulo",
			Type:   "Estadual",
			Entity: "Tribunal",
			UF:     "SP",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agency: %q", err)
	}

	returnedAgency, err := postgresDb.GetAgency("tjsp")

	assert.Nil(t, err)
	assert.Equal(t, agencies[0], *returnedAgency)
}

func (g getAgency) testWhenAgencyNotExists(t *testing.T) {
	truncateTables()

	returnedAgency, err := postgresDb.GetAgency("tjsp")

	expectedErr := fmt.Errorf("error getting agency 'tjsp': %q", gorm.ErrRecordNotFound)
	assert.Nil(t, returnedAgency)
	assert.Equal(t, expectedErr, err)
}

func (g getAgency) testWhenAgencyIsInIrregularCase(t *testing.T) {
	agencies := []models.Agency{
		{ID: "tjsp",
			Name:   "Tribunal de Justiça do Estado de São Paulo",
			Type:   "Estadual",
			Entity: "Tribunal",
			UF:     "SP",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agency: %q", err)
	}

	returnedAgency, err := postgresDb.GetAgency("tJsp")

	assert.Nil(t, err)
	assert.Equal(t, agencies[0], *returnedAgency)
}

func TestGetAllAgencies(t *testing.T) {
	tests := getAllAgencies{}
	t.Run("Test GetAllAgencies when agencies exists", tests.testWhenAgenciesExists)
	t.Run("Test GetAllAgencies when agencies not exists", tests.testWhenAgenciesNotExists)
}

type getAllAgencies struct{}

func (g getAllAgencies) testWhenAgenciesExists(t *testing.T) {
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

	returnedAgencies, err := postgresDb.GetAllAgencies()

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
}

func (g getAllAgencies) testWhenAgenciesNotExists(t *testing.T) {
	truncateTables()

	returnedAgencies, err := postgresDb.GetAllAgencies()

	assert.Nil(t, err)
	assert.Empty(t, returnedAgencies)
}

func TestStore(t *testing.T) {
	if err := insertAgencies([]models.Agency{{ID: "tjba"}}); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	timestamp, _ := time.Parse("2006-01-02 15:04:05.999", "2023-01-16 04:55:11.930") // convertendo string para time.Time
	agmi := models.AgencyMonthlyInfo{
		AgencyID: "tjba",
		Month:    12,
		Year:     2022,
		Backups: []models.Backup{
			{
				URL:  "https://dadosjusbr-public.s3.amazonaws.com/tjba/backups/tjba-2022-12.zip",
				Hash: "2cc54da4571ca9ff2d416a198cd09669",
				Size: 173253,
			},
		},
		Summary: &models.Summary{
			Count: 662,
			BaseRemuneration: models.DataSummary{
				Max:     35462.22,
				Min:     27098.07,
				Average: 31930.475453172014,
				Total:   21137974.749999873,
			},
			OtherRemunerations: models.DataSummary{
				Max:     243308.90999999997,
				Min:     35974.35,
				Average: 96290.11472809668,
				Total:   63744055.95,
			},
			IncomeHistogram: map[int]int{-1: 0, 10000: 0, 20000: 0, 30000: 116, 40000: 546, 50000: 0},
		},
		CrawlerVersion:    "b9ec52df612cda045544543a3b0387842475764d",
		CrawlerRepo:       "https://github.com/dadosjusbr/coletor-cnj",
		ParserVersion:     "sha256:e0b5858e2d11a2e4183a32c490517ec440020ad8ca549ae86544dbc7683dcfbb",
		ParserRepo:        "https://github.com/dadosjusbr/parser-cnj",
		CrawlingTimestamp: timestamppb.New(timestamp),
		Package: &models.Backup{
			URL:  "https://dadosjusbr-public.s3.amazonaws.com/tjba/datapackage/tjba-2022-12.zip",
			Hash: "ec2651e8e9068a1c2f7e1bfec10ce718",
			Size: 94219,
		},
		Meta: &models.Meta{
			OpenFormat:       false,
			Access:           "NECESSITA_SIMULACAO_USUARIO",
			Extension:        "XLS",
			StrictlyTabular:  true,
			ConsistentFormat: true,
			HaveEnrollment:   false,
			ThereIsACapacity: false,
			HasPosition:      false,
			BaseRevenue:      "DETALHADO",
			OtherRecipes:     "DETALHADO",
			Expenditure:      "DETALHADO",
		},
		Score: &models.Score{
			Score:             0.5,
			CompletenessScore: 0.5,
			EasinessScore:     0.5,
		},
		Duration: 305,
	}

	err := postgresDb.Store(agmi)

	var count int64
	var dtoAgmi dto.AgencyMonthlyInfoDTO

	m := postgresDb.db.Model(dto.AgencyMonthlyInfoDTO{}).Where("id = 'tjba/12/2022' AND atual = true").Count(&count).Find(&dtoAgmi)
	if m.Error != nil {
		fmt.Errorf("error finding agmi: %v", err)
	}

	result, err := dtoAgmi.ConvertToModel()
	if err != nil {
		fmt.Errorf("error converting agmi dto to model: %q", err)
	}

	// Verificando se o método Store deu erro,
	// se tem apenas 1 com atual == true e se todos os campos foram armazenados.
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, agmi.AgencyID, result.AgencyID)
	assert.Equal(t, agmi.Backups, result.Backups)
	assert.Equal(t, agmi.Package.Hash, result.Package.Hash)
	assert.Equal(t, agmi.Summary.BaseRemuneration, result.Summary.BaseRemuneration)
	assert.Equal(t, agmi.Meta.Extension, result.Meta.Extension)
	assert.Equal(t, agmi.Score.Score, result.Score.Score)
	assert.Equal(t, agmi.Duration, result.Duration)
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

func insertMonthlyInfos(monthlyInfos []models.AgencyMonthlyInfo) error {
	for _, monthlyInfo := range monthlyInfos {
		monthlyInfoDto, err := dto.NewAgencyMonthlyInfoDTO(monthlyInfo)
		if err != nil {
			return fmt.Errorf("error creating monthly info dto: %q", err)
		}
		tx := postgresDb.db.Model(dto.AgencyMonthlyInfoDTO{}).Create(monthlyInfoDto)
		if tx.Error != nil {
			return fmt.Errorf("error inserting monthly info: %q", tx.Error)
		}
	}
	return nil
}

func truncateTables() error {
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
