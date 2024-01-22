package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dadosjusbr/proto/coleta"
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
	postgresDb.Disconnect()
	os.Exit(exitValue)
}

func TestGetStateAgencies(t *testing.T) {
	tests := getStateAgencies{}
	t.Run("Test TestGetStateAgencies when agencies exists", tests.testWhenAgenciesExists)
	t.Run("Test TestGetStateAgencies when UF not exists", tests.testWhenUFNotExists)
	t.Run("Test TestGetStateAgencies when UF is in lower case", tests.testWhenUFIsInLowerCase)
}

type getStateAgencies struct{}

func (g getStateAgencies) testWhenAgenciesExists(t *testing.T) {
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
	returnedAgencies, err := postgresDb.GetStateAgencies("SP")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
	truncateTables()
}

func (g getStateAgencies) testWhenUFNotExists(t *testing.T) {
	truncateTables()

	returnedAgencies, err := postgresDb.GetStateAgencies("SP")

	assert.Nil(t, err)
	assert.Empty(t, returnedAgencies)
}

func (g getStateAgencies) testWhenUFIsInLowerCase(t *testing.T) {
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

	returnedAgencies, err := postgresDb.GetStateAgencies("sp")

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
	truncateTables()
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
	truncateTables()
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
	truncateTables()
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
	truncateTables()
}

func TestGetAllAgencies(t *testing.T) {
	tests := getAllAgencies{}
	t.Run("Test GetAllAgencies when agencies exists", tests.testWhenAgenciesExists)
	t.Run("Test GetAllAgencies when agencies not exists", tests.testWhenAgenciesNotExists)
}

type getAllAgencies struct{}

func (g getAllAgencies) testWhenAgenciesExists(t *testing.T) {
	timestamp := int64(1643724131)
	agencies := []models.Agency{
		{
			ID:   "tjsp",
			Name: "Tribunal de Justiça do Estado de São Paulo",

			Type:   "Estadual",
			Entity: "Tribunal",
			UF:     "SP",
			Collecting: []models.Collecting{{
				Timestamp:   &timestamp,
				Description: []string{"Não há dados abertos disponíveis"},
				Collecting:  true,
			},
			},
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
	assert.NotNil(t, returnedAgencies[0].Collecting)
	assert.True(t, returnedAgencies[0].Collecting[0].Collecting)
}

func (g getAllAgencies) testWhenAgenciesNotExists(t *testing.T) {
	truncateTables()

	returnedAgencies, err := postgresDb.GetAllAgencies()

	assert.Nil(t, err)
	assert.Empty(t, returnedAgencies)
}

func TestGetGeneralMonthlyInfo(t *testing.T) {
	tests := getGeneralMonthlyInfo{}
	t.Run("Test GetGeneralMonthlyInfo when monthly info exists", tests.testWhenMonthlyInfoExists)
	t.Run("Test GetGeneralMonthlyInfo when monthly info not exists", tests.testWhenMonthlyInfoNotExists)
}

type getGeneralMonthlyInfo struct{}

func (g getGeneralMonthlyInfo) testWhenMonthlyInfoExists(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjsp",
		},
		{
			ID: "tjal",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID:          "tjal",
			Year:              2020,
			Month:             3,
			CrawlingTimestamp: timestamppb.Now(),
			Summary: &models.Summary{
				BaseRemuneration: models.DataSummary{
					Total: 1200,
				},
				OtherRemunerations: models.DataSummary{
					Total: 1000,
				},
			},
		},
		{
			AgencyID:          "tjsp",
			Year:              2020,
			Month:             3,
			CrawlingTimestamp: timestamppb.Now(),
			Summary: &models.Summary{
				BaseRemuneration: models.DataSummary{
					Total: 1000,
				},
				OtherRemunerations: models.DataSummary{
					Total: 1300,
				},
			},
		},
	}

	var total float64
	for _, agmi := range agmis {
		total += agmi.Summary.BaseRemuneration.Total
		total += agmi.Summary.OtherRemunerations.Total
	}

	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	value, err := postgresDb.GetGeneralMonthlyInfo()

	assert.Nil(t, err)
	assert.Equal(t, value, total)
	truncateTables()
}

func (g getGeneralMonthlyInfo) testWhenMonthlyInfoNotExists(t *testing.T) {
	truncateTables()

	value, err := postgresDb.GetGeneralMonthlyInfo()

	assert.Nil(t, err)
	assert.Equal(t, value, float64(0))
}

func TestGetLastDateWithMonthlyInfo(t *testing.T) {
	tests := getLastDateWithMonthlyInfo{}
	t.Run("Test GetLastDateWithMonthlyInfo when monthly infos exists", tests.testWhenMonthlyInfosExists)
	t.Run("Test GetLastDateWithMonthlyInfo when monthly infos is empty", tests.testWhenMonthlyInfosIsEmpty)
	t.Run("Test GetLastDateWithMonthlyInfo when monthly infos is equal", tests.testWhenMonthlyInfosIsEqual)
}

type getLastDateWithMonthlyInfo struct{}

func (g getLastDateWithMonthlyInfo) testWhenMonthlyInfosExists(t *testing.T) {
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
	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID:          "tjba",
			Year:              2022,
			Month:             3,
			CrawlingTimestamp: timestamppb.Now(),
		},
		{
			AgencyID:          "tjal",
			Year:              2022,
			Month:             2,
			CrawlingTimestamp: timestamppb.Now(),
		},
		{
			AgencyID:          "tjsp",
			Year:              2020,
			Month:             4,
			CrawlingTimestamp: timestamppb.Now(),
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	month, year, err := postgresDb.GetLastDateWithMonthlyInfo()

	assert.Nil(t, err)
	assert.Equal(t, agmis[0].Month, month)
	assert.Equal(t, agmis[0].Year, year)
	truncateTables()
}

func (g getLastDateWithMonthlyInfo) testWhenMonthlyInfosIsEmpty(t *testing.T) {
	truncateTables()

	month, year, err := postgresDb.GetLastDateWithMonthlyInfo()

	assert.NotEmpty(t, err)
	assert.Equal(t, 0, month)
	assert.Equal(t, 0, year)
}

func (g getLastDateWithMonthlyInfo) testWhenMonthlyInfosIsEqual(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjsp",
		},
		{
			ID: "tjal",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID:          "tjal",
			Year:              2020,
			Month:             3,
			CrawlingTimestamp: timestamppb.Now(),
		},
		{
			AgencyID:          "tjsp",
			Year:              2020,
			Month:             3,
			CrawlingTimestamp: timestamppb.Now(),
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	month, year, err := postgresDb.GetFirstDateWithMonthlyInfo()

	assert.Nil(t, err)
	assert.Equal(t, agmis[0].Month, month)
	assert.Equal(t, agmis[0].Year, year)
	truncateTables()
}

func TestGetFirstDateWithMonthlyInfo(t *testing.T) {
	tests := getFirstDateWithMonthlyInfo{}
	t.Run("Test GetFirstDateWithMonthlyInfo when monthly infos exists", tests.testWhenMonthlyInfosExists)
	t.Run("Test GetFirstDateWithMonthlyInfo when monthly infos is empty", tests.testWhenMonthlyInfosIsEmpty)
	t.Run("Test GetFirstDateWithMonthlyInfo when monthly infos is equal", tests.testWhenMonthlyInfosIsEqual)
}

type getFirstDateWithMonthlyInfo struct{}

func (g getFirstDateWithMonthlyInfo) testWhenMonthlyInfosExists(t *testing.T) {
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
	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID:          "tjsp",
			Year:              2020,
			Month:             4,
			CrawlingTimestamp: timestamppb.Now(),
		},
		{
			AgencyID:          "tjba",
			Year:              2022,
			Month:             3,
			CrawlingTimestamp: timestamppb.Now(),
		},
		{
			AgencyID:          "tjal",
			Year:              2022,
			Month:             2,
			CrawlingTimestamp: timestamppb.Now(),
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	month, year, err := postgresDb.GetFirstDateWithMonthlyInfo()

	assert.Nil(t, err)
	assert.Equal(t, agmis[0].Month, month)
	assert.Equal(t, agmis[0].Year, year)
	truncateTables()
}

func (g getFirstDateWithMonthlyInfo) testWhenMonthlyInfosIsEmpty(t *testing.T) {
	truncateTables()

	month, year, err := postgresDb.GetFirstDateWithMonthlyInfo()

	assert.NotEmpty(t, err)
	assert.Equal(t, 0, month)
	assert.Equal(t, 0, year)
}

func (g getFirstDateWithMonthlyInfo) testWhenMonthlyInfosIsEqual(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjsp",
		},
		{
			ID: "tjal",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID:          "tjal",
			Year:              2020,
			Month:             3,
			CrawlingTimestamp: timestamppb.Now(),
		},
		{
			AgencyID:          "tjsp",
			Year:              2020,
			Month:             3,
			CrawlingTimestamp: timestamppb.Now(),
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	month, year, err := postgresDb.GetFirstDateWithMonthlyInfo()

	assert.Nil(t, err)
	assert.Equal(t, agmis[0].Month, month)
	assert.Equal(t, agmis[0].Year, year)
	truncateTables()
}

func TestGetAgenciesByUF(t *testing.T) {
	tests := getAgenciesByUF{}
	t.Run("Test GetAgenciesByUF when agencies exists", tests.testWhenAgenciesExists)
	t.Run("Test GetAgenciesByUF when UF not exists", tests.testWhenUFNotExists)
	t.Run("Test GetAgenciesByUF when UF is in irregular case", tests.testWhenUFIsInIrregularCase)
}

type getAgenciesByUF struct{}

func (g getAgenciesByUF) testWhenAgenciesExists(t *testing.T) {
	agencies := []models.Agency{
		{
			ID:           "mpsp",
			Type:         "Ministério",
			UF:           "SP",
			OmbudsmanURL: "https://sis.mpsp.mp.br/atendimentocidadao/Ouvidoria/Manifestacao/EscolherTipoDeIdentificacao",
		},
		{
			ID:   "tjsp",
			Type: "Estadual",
			UF:   "SP",
		},
		{
			ID:   "tjmsp",
			Type: "Militar",
			UF:   "SP",
		},
		{
			ID:   "tjal",
			Type: "Estadual",
			UF:   "AL",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgencies, err := postgresDb.GetAgenciesByUF("SP")

	assert.Nil(t, err)
	assert.Equal(t, agencies[:3], returnedAgencies)
	truncateTables()
}

func (g getAgenciesByUF) testWhenUFNotExists(t *testing.T) {
	truncateTables()

	returnedAgencies, err := postgresDb.GetAgenciesByUF("SP")

	assert.Nil(t, err)
	assert.Empty(t, returnedAgencies)
}

func (g getAgenciesByUF) testWhenUFIsInIrregularCase(t *testing.T) {
	agencies := []models.Agency{
		{
			ID:           "mpsp",
			Type:         "Ministério",
			UF:           "SP",
			OmbudsmanURL: "https://sis.mpsp.mp.br/atendimentocidadao/Ouvidoria/Manifestacao/EscolherTipoDeIdentificacao",
		},
		{
			ID:   "tjsp",
			Type: "Estadual",
			UF:   "SP",
		},
		{
			ID:   "tjmsp",
			Type: "Militar",
			UF:   "SP",
		},
		{
			ID:   "tjal",
			Type: "Estadual",
			UF:   "AL",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgencies, err := postgresDb.GetAgenciesByUF("sP")

	assert.Nil(t, err)
	assert.Equal(t, agencies[:3], returnedAgencies)
	truncateTables()
}

func TestGetMonthlyInfo(t *testing.T) {
	tests := getMonthlyInfo{}
	t.Run("Test GetMonthlyInfo when monthly info exists", tests.testWhenMonthlyInfoExists)
	t.Run("Test GetMonthlyInfo when agency not exists", tests.testWhenAgencyNotExists)
	t.Run("Test GetMonthlyInfo when year not exists", tests.testWhenYearNotExists)
	t.Run("Test GetMonthlyInfo when procinfo is not null", tests.testWhenProcInfoIsNotNull)
}

type getMonthlyInfo struct{}

func (g getMonthlyInfo) testWhenMonthlyInfoExists(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjsp",
		},
		{
			ID: "tjal",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID:          "tjal",
			Year:              2020,
			Month:             1,
			CrawlingTimestamp: timestamppb.Now(),
		},
		{
			AgencyID:          "tjsp",
			Year:              2020,
			Month:             1,
			CrawlingTimestamp: timestamppb.Now(),
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agency monthly info: %q", err)
	}
	var agmiMap = make(map[string][]models.AgencyMonthlyInfo)
	for _, agmi := range agmis {
		agmiMap[agmi.AgencyID] = append(agmiMap[agmi.AgencyID], agmi)
	}

	returnedAgmis, err := postgresDb.GetMonthlyInfo(agencies, 2020)

	assert.Nil(t, err)
	assert.Equal(t, agmiMap["tjal"][0].AgencyID, returnedAgmis["tjal"][0].AgencyID)
	assert.Equal(t, agmiMap["tjal"][0].Year, returnedAgmis["tjal"][0].Year)
	assert.Equal(t, agmiMap["tjal"][0].Month, returnedAgmis["tjal"][0].Month)
	assert.Equal(t, agmiMap["tjsp"][0].AgencyID, returnedAgmis["tjsp"][0].AgencyID)
	assert.Equal(t, agmiMap["tjsp"][0].Year, returnedAgmis["tjsp"][0].Year)
	assert.Equal(t, agmiMap["tjsp"][0].Month, returnedAgmis["tjsp"][0].Month)
	truncateTables()
}

func (g getMonthlyInfo) testWhenAgencyNotExists(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjal",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID:          "tjal",
			Year:              2020,
			Month:             1,
			CrawlingTimestamp: timestamppb.Now(),
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agency monthly info: %q", err)
	}

	returnedAgmis, err := postgresDb.GetMonthlyInfo([]models.Agency{{ID: "tjsp"}}, 2020)

	assert.Nil(t, err)
	assert.Empty(t, returnedAgmis)
	truncateTables()
}

func (g getMonthlyInfo) testWhenYearNotExists(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjsp",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID:          "tjsp",
			Year:              2021,
			Month:             1,
			CrawlingTimestamp: timestamppb.Now(),
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agency monthly info: %q", err)
	}

	returnedAgmis, err := postgresDb.GetMonthlyInfo([]models.Agency{{ID: "tjsp"}}, 2020)

	assert.Nil(t, err)
	assert.Empty(t, returnedAgmis)
}

func (g getMonthlyInfo) testWhenProcInfoIsNotNull(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjsp",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID:          "tjsp",
			Year:              2020,
			Month:             1,
			CrawlingTimestamp: timestamppb.Now(),
			ProcInfo: &coleta.ProcInfo{
				Stdin:  "stdin",
				Stdout: "stdout",
				Stderr: "stderr",
				Cmd:    "cmd",
				CmdDir: "cmdDir",
				Status: 4,
				Env:    []string{"env"},
			},
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agency monthly info: %q", err)
	}

	returnedAgmis, err := postgresDb.GetMonthlyInfo([]models.Agency{{ID: "tjsp"}}, 2020)

	assert.Nil(t, err)
	assert.Empty(t, returnedAgmis)
	truncateTables()
}

func TestGetAnnualSummary(t *testing.T) {
	tests := getAnnualSummary{}

	t.Run("Test GetAnnualSummary when monthly info exists", tests.testWhenMonthlyInfoExists)
	t.Run("Test GetAnnualSummary when agency not exists", tests.testWhenAgencyNotExists)
}

type getAnnualSummary struct{}

func (g getAnnualSummary) testWhenMonthlyInfoExists(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjal",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID:          "tjal",
			Year:              2020,
			Month:             1,
			CrawlingTimestamp: timestamppb.Now(),
			Summary: &models.Summary{
				Count: 100,
				BaseRemuneration: models.DataSummary{
					Total: 1000,
				},
				OtherRemunerations: models.DataSummary{
					Total: 500,
				},
				Discounts: models.DataSummary{
					Total: 500,
				},
				Remunerations: models.DataSummary{
					Total: 1000,
				},
				ItemSummary: models.ItemSummary{
					Others: 100,
				},
			},
		},
		{
			AgencyID:          "tjal",
			Year:              2020,
			Month:             2,
			CrawlingTimestamp: timestamppb.Now(),
			Summary: &models.Summary{
				Count: 150,
				BaseRemuneration: models.DataSummary{
					Total: 1200,
				},
				OtherRemunerations: models.DataSummary{
					Total: 500,
				},
				Discounts: models.DataSummary{
					Total: 500,
				},
				Remunerations: models.DataSummary{
					Total: 1200,
				},
			},
		},
		{
			AgencyID:          "tjal",
			Year:              2021,
			Month:             1,
			CrawlingTimestamp: timestamppb.Now(),
			Summary: &models.Summary{
				Count: 200,
				BaseRemuneration: models.DataSummary{
					Total: 1500,
				},
				OtherRemunerations: models.DataSummary{
					Total: 500,
				},
				Discounts: models.DataSummary{
					Total: 500,
				},
				Remunerations: models.DataSummary{
					Total: 1500,
				},
				ItemSummary: models.ItemSummary{
					Others:        100,
					FoodAllowance: 150,
					BonusLicense:  200,
				},
			},
		},
		{
			AgencyID:          "tjal",
			Year:              2021,
			Month:             2,
			CrawlingTimestamp: timestamppb.Now(),
			Summary: &models.Summary{
				Count: 300,
				BaseRemuneration: models.DataSummary{
					Total: 1000,
				},
				OtherRemunerations: models.DataSummary{
					Total: 500,
				},
				Discounts: models.DataSummary{
					Total: 500,
				},
				Remunerations: models.DataSummary{
					Total: 1000,
				},
			},
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agency monthly info: %q", err)
	}

	var amis []models.AnnualSummary
	//Realizando a soma das remunerações por ano
	for _, agmi := range agmis {
		for _, agmi2 := range agmis {
			exists := false
			for _, gmi := range amis {
				if gmi.Year == agmi.Year && agmi.Month != agmi2.Month {
					exists = true
				}
			}
			if !exists && agmi.Year == agmi2.Year && agmi.Month != agmi2.Month {
				if agmi.Year == agmi2.Year && agmi.Month != agmi2.Month {
					amis = append(amis, models.AnnualSummary{
						Year:               agmi.Year,
						AverageCount:       (agmi.Summary.Count + agmi2.Summary.Count) / 2,
						TotalCount:         agmi.Summary.Count + agmi2.Summary.Count,
						BaseRemuneration:   agmi.Summary.BaseRemuneration.Total + agmi2.Summary.BaseRemuneration.Total,
						OtherRemunerations: agmi.Summary.OtherRemunerations.Total + agmi2.Summary.OtherRemunerations.Total,
						Discounts:          agmi.Summary.Discounts.Total + agmi2.Summary.Discounts.Total,
						Remunerations:      agmi.Summary.Remunerations.Total + agmi2.Summary.Remunerations.Total,
						ItemSummary: models.ItemSummary{
							Others:        agmi.Summary.ItemSummary.Others + agmi2.Summary.ItemSummary.Others,
							BonusLicense:  agmi.Summary.ItemSummary.BonusLicense + agmi2.Summary.ItemSummary.BonusLicense,
							FoodAllowance: agmi.Summary.ItemSummary.FoodAllowance + agmi2.Summary.ItemSummary.FoodAllowance,
						},
					})
				}
			}
		}
	}

	returnedAmis, err := postgresDb.GetAnnualSummary("tjal")

	assert.Nil(t, err)
	assert.Equal(t, amis[0].Year, returnedAmis[0].Year)
	assert.Equal(t, amis[0].BaseRemuneration, returnedAmis[0].BaseRemuneration)
	assert.Equal(t, amis[0].OtherRemunerations, returnedAmis[0].OtherRemunerations)
	assert.Equal(t, amis[0].Discounts, returnedAmis[0].Discounts)
	assert.Equal(t, amis[0].Remunerations, returnedAmis[0].Remunerations)
	assert.Equal(t, amis[0].AverageCount, returnedAmis[0].AverageCount)
	assert.Equal(t, amis[1].AverageCount, returnedAmis[1].AverageCount)
	assert.Equal(t, amis[0].TotalCount, returnedAmis[0].TotalCount)
	assert.Equal(t, amis[1].TotalCount, returnedAmis[1].TotalCount)
	assert.Equal(t, 2, returnedAmis[0].NumMonthsWithData)
	assert.Equal(t, amis[0].ItemSummary.Others, returnedAmis[0].ItemSummary.Others)
	assert.Equal(t, amis[1].ItemSummary.BonusLicense, returnedAmis[1].ItemSummary.BonusLicense)
	assert.Equal(t, amis[1].ItemSummary.FoodAllowance, returnedAmis[1].ItemSummary.FoodAllowance)
	truncateTables()
}

func (g getAnnualSummary) testWhenAgencyNotExists(t *testing.T) {
	truncateTables()
	returnedAmis, err := postgresDb.GetAnnualSummary("tjsp")

	assert.Nil(t, err)
	assert.Empty(t, returnedAmis)
}

func TestGetOMA(t *testing.T) {
	tests := getOMA{}

	t.Run("Test GetOMA when data exists", tests.testWhenDataExists)
	t.Run("Test GetOMA when data not exists", tests.testWhenDataNotExists)
	t.Run("Test GetOMA when agency is in irregular case", tests.testWhenAgencyIsInIrregularCase)
}

type getOMA struct{}

func (g getOMA) testWhenDataExists(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjsp",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	agmi := models.AgencyMonthlyInfo{
		AgencyID:          "tjsp",
		Month:             12,
		Year:              2022,
		CrawlingTimestamp: timestamppb.Now(),
	}
	if err := insertMonthlyInfos([]models.AgencyMonthlyInfo{agmi}); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgmi, agency, err := postgresDb.GetOMA(12, 2022, "tjsp")

	assert.Nil(t, err)
	assert.Equal(t, agmi.AgencyID, returnedAgmi.AgencyID)
	assert.Equal(t, agmi.Month, returnedAgmi.Month)
	assert.Equal(t, agmi.Year, returnedAgmi.Year)
	assert.Equal(t, agencies[0], *agency)
	truncateTables()
}

func (g getOMA) testWhenDataNotExists(t *testing.T) {
	truncateTables()
	expecErr := fmt.Errorf("there is no data with this parameters")
	returnedAgmi, agency, err := postgresDb.GetOMA(12, 2022, "tjba")

	assert.Equal(t, err, expecErr)
	assert.Nil(t, returnedAgmi)
	assert.Nil(t, agency)
}

func (g getOMA) testWhenAgencyIsInIrregularCase(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjsp",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	agmi := models.AgencyMonthlyInfo{
		AgencyID:          "tjsp",
		Month:             12,
		Year:              2022,
		CrawlingTimestamp: timestamppb.Now(),
	}
	if err := insertMonthlyInfos([]models.AgencyMonthlyInfo{agmi}); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgmi, agency, err := postgresDb.GetOMA(12, 2022, "tJsp")

	assert.Nil(t, err)
	assert.Equal(t, agmi.AgencyID, returnedAgmi.AgencyID)
	assert.Equal(t, agmi.Month, returnedAgmi.Month)
	assert.Equal(t, agmi.Year, returnedAgmi.Year)
	assert.Equal(t, agencies[0], *agency)
	truncateTables()
}

func TestGetGeneralMonthlyInfosFromYear(t *testing.T) {
	tests := getGeneralMonthlyInfoFromYear{}

	t.Run("Test GetGeneralMonthlyInfosFromYear when monthly infos exists", tests.testWhenDataExists)
	t.Run("Test GetGeneralMonthlyInfosFromYear when monthly infos not exists", tests.testWhenDataNotExists)
}

type getGeneralMonthlyInfoFromYear struct{}

func (g getGeneralMonthlyInfoFromYear) testWhenDataExists(t *testing.T) {
	agencies := []models.Agency{
		{
			ID: "tjsp",
		},
		{
			ID: "tjba",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}
	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID:          "tjsp",
			Month:             1,
			Year:              2022,
			CrawlingTimestamp: timestamppb.Now(),
			Summary: &models.Summary{
				Count: 100,
				BaseRemuneration: models.DataSummary{
					Total: 1000,
				},
				OtherRemunerations: models.DataSummary{
					Total: 450,
				},
				Discounts: models.DataSummary{
					Total: 450,
				},
				Remunerations: models.DataSummary{
					Total: 1000,
				},
				ItemSummary: models.ItemSummary{
					FoodAllowance: 200,
				},
			},
		},
		{
			AgencyID:          "tjba",
			Month:             1,
			Year:              2022,
			CrawlingTimestamp: timestamppb.Now(),
			Summary: &models.Summary{
				Count: 300,
				BaseRemuneration: models.DataSummary{
					Total: 3000,
				},
				OtherRemunerations: models.DataSummary{
					Total: 1200,
				},
				Discounts: models.DataSummary{
					Total: 450,
				},
				Remunerations: models.DataSummary{
					Total: 3750,
				},
				ItemSummary: models.ItemSummary{
					BonusLicense: 400,
				},
			},
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	var gmis []models.GeneralMonthlyInfo
	for _, agmi := range agmis {
		for _, agmi2 := range agmis {
			exists := false
			for _, gmi := range gmis {
				if gmi.Month == agmi.Month {
					gmi.BaseRemuneration += agmi.Summary.BaseRemuneration.Total
					gmi.OtherRemunerations += agmi.Summary.OtherRemunerations.Total
					gmi.Discounts += agmi.Summary.Discounts.Total
					gmi.Remunerations += agmi.Summary.Remunerations.Total
					gmi.Count += agmi.Summary.Count
					exists = true
				}
			}
			if !exists && agmi.Month == agmi2.Month && agmi.AgencyID != agmi2.AgencyID {
				if agmi.Month == agmi2.Month && agmi.AgencyID != agmi2.AgencyID {
					gmis = append(gmis, models.GeneralMonthlyInfo{
						Month:              agmi.Month,
						Count:              agmi.Summary.Count + agmi2.Summary.Count,
						BaseRemuneration:   agmi.Summary.BaseRemuneration.Total + agmi2.Summary.BaseRemuneration.Total,
						OtherRemunerations: agmi.Summary.OtherRemunerations.Total + agmi2.Summary.OtherRemunerations.Total,
						Discounts:          agmi.Summary.Discounts.Total + agmi2.Summary.Discounts.Total,
						Remunerations:      agmi.Summary.Remunerations.Total + agmi2.Summary.Remunerations.Total,
						ItemSummary: models.ItemSummary{
							FoodAllowance: agmi.Summary.ItemSummary.FoodAllowance + agmi2.Summary.ItemSummary.FoodAllowance,
							BonusLicense:  agmi.Summary.ItemSummary.BonusLicense + agmi2.Summary.ItemSummary.BonusLicense,
						},
					})
				}
			}
		}
	}
	returnedGmis, err := postgresDb.GetGeneralMonthlyInfosFromYear(2022)

	assert.Nil(t, err)
	assert.Equal(t, gmis, returnedGmis)
	truncateTables()
}

func (g getGeneralMonthlyInfoFromYear) testWhenDataNotExists(t *testing.T) {
	truncateTables()
	returnedGmis, err := postgresDb.GetGeneralMonthlyInfosFromYear(2022)

	assert.Nil(t, err)
	assert.Empty(t, returnedGmis)
}

func TestStoreRemunerations(t *testing.T) {
	tests := storeRemunerations{}

	t.Run("Test StoreRemunerations when data is ok", tests.testWhenDataIsOk)
	t.Run("Test StoreRemunerations when ID already exists", tests.testWhenIDAlreadyExists)
}

type storeRemunerations struct{}

func (s storeRemunerations) testWhenDataIsOk(t *testing.T) {
	remunerations := models.Remunerations{
		AgencyID:     "tjsp",
		Year:         2020,
		Month:        1,
		NumBase:      100,
		NumDiscounts: 100,
		NumOther:     100,
		ZipUrl:       "https://dadosjusbr-public.s3.amazonaws.com/tjsp/remunerations/tjsp-2020-01.zip",
	}
	err := postgresDb.StoreRemunerations(remunerations)

	var count int64
	var remuDTO dto.RemunerationsDTO

	m := postgresDb.db.Model(&dto.RemunerationsDTO{}).Count(&count).Where("id_orgao = ? AND ano = ? AND mes = ?", remunerations.AgencyID, remunerations.Year, remunerations.Month).First(&remuDTO)
	if m.Error != nil {
		t.Fatalf("error getting remunerations: %q", m.Error)
	}

	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, remunerations.AgencyID, remuDTO.AgencyID)
	assert.Equal(t, remunerations.Year, remuDTO.Year)
	assert.Equal(t, remunerations.Month, remuDTO.Month)
	assert.Equal(t, remunerations.NumBase, remuDTO.NumBase)
	assert.Equal(t, remunerations.NumDiscounts, remuDTO.NumDiscounts)
	assert.Equal(t, remunerations.NumOther, remuDTO.NumOther)
	truncateTables()
}

func (s storeRemunerations) testWhenIDAlreadyExists(t *testing.T) {
	remunerations := []models.Remunerations{
		{
			AgencyID:     "tjsp",
			Year:         2020,
			Month:        1,
			NumBase:      100,
			NumDiscounts: 100,
			NumOther:     100,
			ZipUrl:       "https://dadosjusbr-public.s3.amazonaws.com/tjsp/remunerations/tjsp-2020-01.zip",
		},
	}
	if err := insertRemunerations(remunerations); err != nil {
		t.Fatalf("error inserting remunerations: %q", err)
	}
	err := postgresDb.StoreRemunerations(remunerations[0])

	var count int64
	var remuDTO dto.RemunerationsDTO

	m := postgresDb.db.Model(&dto.RemunerationsDTO{}).Count(&count).Where("id_orgao = ? AND ano = ? AND mes = ?", remunerations[0].AgencyID, remunerations[0].Year, remunerations[0].Month).First(&remuDTO)
	if m.Error != nil {
		t.Fatalf("error getting remunerations: %q", m.Error)
	}

	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, remunerations[0].AgencyID, remuDTO.AgencyID)
	assert.Equal(t, remunerations[0].Year, remuDTO.Year)
	assert.Equal(t, remunerations[0].Month, remuDTO.Month)
	assert.Equal(t, remunerations[0].NumBase, remuDTO.NumBase)
	assert.Equal(t, remunerations[0].NumDiscounts, remuDTO.NumDiscounts)
	assert.Equal(t, remunerations[0].NumOther, remuDTO.NumOther)
	truncateTables()
}

func TestStore(t *testing.T) {
	tests := store{}

	t.Run("Test Store when data is OK", tests.testWhenDataIsOK)
	t.Run("Test Store when ID already exists", tests.testWhenIDAlreadyExists)
}

type store struct{}

func (s store) testWhenDataIsOK(t *testing.T) {
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
			Discounts: models.DataSummary{
				Max:     243308.90999999997,
				Min:     35974.35,
				Average: 96290.11472809668,
				Total:   63744055.95,
			},
			Remunerations: models.DataSummary{
				Max:     243308.90999999997,
				Min:     35974.35,
				Average: 96290.11472809668,
				Total:   63744055.95,
			},
			IncomeHistogram: map[int]int{-1: 0, 10000: 0, 20000: 0, 30000: 116, 40000: 546, 50000: 0},
			ItemSummary: models.ItemSummary{
				FoodAllowance: 100,
				BonusLicense:  150,
				Others:        200,
			},
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
	assert.Equal(t, agmi.Summary.OtherRemunerations, result.Summary.OtherRemunerations)
	assert.Equal(t, agmi.Summary.Remunerations, result.Summary.Remunerations)
	assert.Equal(t, agmi.Summary.Discounts, result.Summary.Discounts)
	assert.Equal(t, agmi.Meta.Extension, result.Meta.Extension)
	assert.Equal(t, agmi.Score.Score, result.Score.Score)
	assert.Equal(t, agmi.Duration, result.Duration)
	assert.Equal(t, agmi.Summary.ItemSummary.FoodAllowance, result.Summary.ItemSummary.FoodAllowance)
	assert.Equal(t, agmi.Summary.ItemSummary.BonusLicense, result.Summary.ItemSummary.BonusLicense)
	truncateTables()
}

func (s store) testWhenIDAlreadyExists(t *testing.T) {
	agency := models.Agency{
		ID: "tjba",
	}
	if err := insertAgencies([]models.Agency{agency}); err != nil {
		t.Errorf("error inserting agency: %v", err)
	}
	agmi := models.AgencyMonthlyInfo{
		AgencyID:          "tjba",
		Month:             12,
		Year:              2022,
		CrawlingTimestamp: timestamppb.New(time.Now()),
	}
	if err := insertMonthlyInfos([]models.AgencyMonthlyInfo{agmi}); err != nil {
		t.Errorf("error inserting agmi: %v", err)
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

	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, agmi.AgencyID, result.AgencyID)
	assert.Equal(t, agmi.Year, result.Year)
	assert.Equal(t, agmi.Month, result.Month)
	truncateTables()
}

func TestGetIndexInformation(t *testing.T) {
	tests := indexInformation{}

	t.Run("Test GetIndexInformation() by group", tests.testGetIndexInformationByGroup)
	t.Run("Test GetIndexInformation() by group and year", tests.testGetIndexInformationByYear)
	t.Run("Test GetIndexInformation() by group, month and year", tests.testGetIndexInformationByMonthAndYear)
	t.Run("Test GetIndexInformation() without parameters (all agencies)", tests.testGetAllIndexInformation)
	t.Run("Test GetAllIndexInformation() by year (all agencies)", tests.testGetAllIndexInformationByYear)
	t.Run("Test GetAllIndexInformation() by month and year (all agencies)", tests.testGetAllIndexInformationByMonthAndYear)
}

type indexInformation struct{}

func (indexInformation) testGetIndexInformationByGroup(t *testing.T) {
	agencies := []models.Agency{
		{
			ID:     "tjsp",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
		{
			ID:     "tjba",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID: "tjsp",
			Month:    1,
			Year:     2022,
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
		},
		{
			AgencyID: "tjsp",
			Month:    2,
			Year:     2022,
			Score: &models.Score{
				Score:             0,
				CompletenessScore: 0,
				EasinessScore:     0,
			},
		},
		{
			AgencyID: "tjba",
			Month:    1,
			Year:     2022,
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
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agg, err := postgresDb.GetIndexInformation("Estadual", 0, 0)
	if err != nil {
		t.Fatalf("error GetIndexInformation(): %q", err)
	}

	assert.Equal(t, len(agg), 2)
	assert.Equal(t, len(agg["tjsp"]), 2)
	assert.Equal(t, len(agg["tjba"]), 1)
	assert.Equal(t, agg["tjsp"][0].Score, agmis[0].Score)
	assert.Equal(t, agg["tjsp"][1].Score.EasinessScore, 0.5)
	assert.Equal(t, agg["tjsp"][1].Score.CompletenessScore, 0.0)
	assert.Equal(t, agg["tjba"][0].Score, agmis[2].Score)
	truncateTables()
}

func (indexInformation) testGetIndexInformationByYear(t *testing.T) {
	agencies := []models.Agency{
		{
			ID:     "tjsp",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
		{
			ID:     "tjba",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID: "tjsp",
			Month:    1,
			Year:     2022,
			Score: &models.Score{
				Score:             0.5,
				CompletenessScore: 0.5,
				EasinessScore:     0.5,
			},
		},
		{
			AgencyID: "tjsp",
			Month:    2,
			Year:     2022,
			Score: &models.Score{
				Score:             0,
				CompletenessScore: 0,
				EasinessScore:     0,
			},
		},
		{
			AgencyID: "tjba",
			Month:    1,
			Year:     2022,
			Score: &models.Score{
				Score:             0.5,
				CompletenessScore: 0.5,
				EasinessScore:     0.5,
			},
		},
		{
			AgencyID: "tjba",
			Month:    1,
			Year:     2020,
			Score: &models.Score{
				Score:             0.5,
				CompletenessScore: 0.5,
				EasinessScore:     0.5,
			},
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agg, err := postgresDb.GetIndexInformation("Estadual", 0, 2022)
	if err != nil {
		t.Fatalf("error GetIndexInformation(): %q", err)
	}

	assert.Equal(t, len(agg), 2)
	assert.Equal(t, len(agg["tjsp"]), 2)
	assert.Equal(t, len(agg["tjba"]), 1)
	assert.Equal(t, agg["tjsp"][0].Year, 2022)
	assert.Equal(t, agg["tjsp"][1].Year, 2022)
	assert.Equal(t, agg["tjba"][0].Year, 2022)
	truncateTables()
}

func (indexInformation) testGetIndexInformationByMonthAndYear(t *testing.T) {
	agencies := []models.Agency{
		{
			ID:     "tjsp",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
		{
			ID:     "tjba",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID: "tjsp",
			Month:    1,
			Year:     2022,
			Score: &models.Score{
				Score:             0.5,
				CompletenessScore: 0.5,
				EasinessScore:     0.5,
			},
		},
		{
			AgencyID: "tjsp",
			Month:    2,
			Year:     2022,
			Score: &models.Score{
				Score:             0,
				CompletenessScore: 0,
				EasinessScore:     0,
			},
		},
		{
			AgencyID: "tjba",
			Month:    1,
			Year:     2022,
			Score: &models.Score{
				Score:             0.5,
				CompletenessScore: 0.5,
				EasinessScore:     0.5,
			},
		},
		{
			AgencyID: "tjba",
			Month:    2,
			Year:     2022,
			Score: &models.Score{
				Score:             0.5,
				CompletenessScore: 0.5,
				EasinessScore:     0.5,
			},
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agg, err := postgresDb.GetIndexInformation("Estadual", 1, 2022)
	if err != nil {
		t.Fatalf("error GetIndexInformation(): %q", err)
	}

	assert.Equal(t, len(agg), 2)
	assert.Equal(t, len(agg["tjsp"]), 1)
	assert.Equal(t, len(agg["tjba"]), 1)
	assert.Equal(t, agg["tjsp"][0].Year, 2022)
	assert.Equal(t, agg["tjba"][0].Year, 2022)
	assert.Equal(t, agg["tjsp"][0].Month, 1)
	assert.Equal(t, agg["tjba"][0].Month, 1)
	truncateTables()
}

func (indexInformation) testGetAllIndexInformation(t *testing.T) {
	agencies := []models.Agency{
		{
			ID:     "tjsp",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
		{
			ID:     "tjba",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID: "tjsp",
			Month:    1,
			Year:     2022,
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
		},
		{
			AgencyID: "tjsp",
			Month:    2,
			Year:     2022,
			Score: &models.Score{
				Score:             0,
				CompletenessScore: 0,
				EasinessScore:     0,
			},
		},
		{
			AgencyID: "tjba",
			Month:    1,
			Year:     2022,
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
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agg, err := postgresDb.GetIndexInformation("", 0, 0)
	if err != nil {
		t.Fatalf("error GetIndexInformation(): %q", err)
	}

	assert.Equal(t, len(agg), 2)
	assert.Equal(t, len(agg["tjsp"]), 2)
	assert.Equal(t, len(agg["tjba"]), 1)
	assert.Equal(t, agg["tjsp"][0].Type, "Estadual")
	truncateTables()
}

func (indexInformation) testGetAllIndexInformationByYear(t *testing.T) {
	agencies := []models.Agency{
		{
			ID:     "tjsp",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
		{
			ID:     "tjba",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID: "tjsp",
			Month:    1,
			Year:     2022,
			Score: &models.Score{
				Score:             0.5,
				CompletenessScore: 0.5,
				EasinessScore:     0.5,
			},
		},
		{
			AgencyID: "tjsp",
			Month:    2,
			Year:     2021,
			Score: &models.Score{
				Score:             0,
				CompletenessScore: 0,
				EasinessScore:     0,
			},
		},
		{
			AgencyID: "tjba",
			Month:    1,
			Year:     2022,
			Score: &models.Score{
				Score:             0.5,
				CompletenessScore: 0.5,
				EasinessScore:     0.5,
			},
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agg, err := postgresDb.GetIndexInformation("", 0, 2022)
	if err != nil {
		t.Fatalf("error GetIndexInformation(): %q", err)
	}

	assert.Equal(t, len(agg), 2)
	assert.Equal(t, len(agg["tjsp"]), 1)
	assert.Equal(t, len(agg["tjba"]), 1)
	assert.Equal(t, agg["tjsp"][0].Year, 2022)
	assert.Equal(t, agg["tjba"][0].Year, 2022)
	truncateTables()
}

func (indexInformation) testGetAllIndexInformationByMonthAndYear(t *testing.T) {
	agencies := []models.Agency{
		{
			ID:     "tjsp",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
		{
			ID:     "tjba",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
	}
	if err := insertAgencies(agencies); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID: "tjsp",
			Month:    1,
			Year:     2021,
			Score: &models.Score{
				Score:             0.5,
				CompletenessScore: 0.5,
				EasinessScore:     0.5,
			},
		},
		{
			AgencyID: "tjsp",
			Month:    1,
			Year:     2022,
			Score: &models.Score{
				Score:             0,
				CompletenessScore: 0,
				EasinessScore:     0,
			},
		},
		{
			AgencyID: "tjba",
			Month:    1,
			Year:     2021,
			Score: &models.Score{
				Score:             0.5,
				CompletenessScore: 0.5,
				EasinessScore:     0.5,
			},
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agg, err := postgresDb.GetIndexInformation("", 1, 2021)
	if err != nil {
		t.Fatalf("error GetIndexInformation(): %q", err)
	}

	assert.Equal(t, len(agg), 2)
	assert.Equal(t, len(agg["tjsp"]), 1)
	assert.Equal(t, len(agg["tjba"]), 1)
	assert.Equal(t, agg["tjsp"][0].Year, 2021)
	assert.Equal(t, agg["tjba"][0].Year, 2021)
	assert.Equal(t, agg["tjsp"][0].Month, 1)
	assert.Equal(t, agg["tjba"][0].Month, 1)
	truncateTables()
}

func TestGetAllAgencyCollection(t *testing.T) {
	agency := []models.Agency{
		{
			ID:     "tjsp",
			Entity: "Tribunal",
			Type:   "Estadual",
		},
	}
	if err := insertAgencies(agency); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	agmis := []models.AgencyMonthlyInfo{
		{
			AgencyID: "tjsp",
			Month:    1,
			Year:     2021,
			Summary: &models.Summary{
				Count: 3407,
				BaseRemuneration: models.DataSummary{
					Max:     47052.46,
					Min:     11735.02,
					Average: 33298.721643673955,
					Total:   113448744.63999715,
				},
				OtherRemunerations: models.DataSummary{
					Max:     82942.56,
					Total:   58807454.34999981,
					Average: 17260.773216906313,
				},
				Discounts: models.DataSummary{
					Max:     82942.56,
					Min:     1756.22,
					Total:   58807454.34999981,
					Average: 17260.773216906313,
				},
				Remunerations: models.DataSummary{
					Max:     47052.46,
					Min:     11735.02,
					Average: 33298.721643673955,
					Total:   113448744.63999715,
				},
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
		},
		{
			AgencyID: "tjsp",
			Month:    2,
			Year:     2022,
			Score: &models.Score{
				Score:             0,
				CompletenessScore: 0,
				EasinessScore:     0,
			},
		},
	}
	if err := insertMonthlyInfos(agmis); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	collections, err := postgresDb.GetAllAgencyCollection("tjsp")
	if err != nil {
		t.Fatalf("error GetAllAgencyCollection(): %q", err)
	}

	assert.Equal(t, len(collections), 2)
	assert.Equal(t, agmis[0].Summary, collections[0].Summary)
	assert.Equal(t, agmis[0].Meta, collections[0].Meta)
	assert.Equal(t, agmis[0].Score, collections[0].Score)
	assert.Equal(t, collections[1].Score.CompletenessScore, 0.0)
	assert.Equal(t, collections[1].Score.EasinessScore, 0.5)
	truncateTables()
}

type paycheck struct{}

func TestPaychecks(t *testing.T) {
	tests := paycheck{}

	t.Run("Test StorePaychecks", tests.testStorePaychecks)
	t.Run("Test StorePaychecks when paycheck already exists", tests.testWhenPaycheckAlreadyExists)
	t.Run("Test GetPaycheck()", tests.testGetPaychecks)
	t.Run("Test GetPaycheckItems()", tests.testGetPaycheckItems)
	t.Run("Test StorePaychecks when paycheck items not exist", tests.testWhenPaycheckItemsNotExist)
}

func (paycheck) testStorePaychecks(t *testing.T) {
	truncateTables()
	p, pi := paychecks()

	err := postgresDb.StorePaychecks(p, pi)
	var dtoPaychecks []dto.PaycheckDTO
	var dtoPaycheckItems []dto.PaycheckItemDTO

	m := postgresDb.db.Model(dto.PaycheckDTO{}).Scan(&dtoPaychecks)
	if m.Error != nil {
		t.Fatalf("error finding payckecks: %v", err)
	}

	n := postgresDb.db.Model(dto.PaycheckItemDTO{}).Scan(&dtoPaycheckItems)
	if n.Error != nil {
		t.Fatalf("error finding payckeck items: %v", err)
	}

	itemSanitizado := "subsidio"
	assert.Nil(t, err)
	assert.Equal(t, len(dtoPaychecks), 1)
	assert.Equal(t, len(dtoPaycheckItems), 3)
	assert.Equal(t, dtoPaychecks[0].Name, "nome")
	assert.Equal(t, dtoPaychecks[0].Remuneration, 2000.0)
	assert.Equal(t, dtoPaycheckItems[0].Type, "R/B")
	assert.Equal(t, dtoPaycheckItems[1].Inconsistent, true)
	assert.Equal(t, dtoPaycheckItems[2].Value, 200.0)
	assert.Equal(t, dtoPaycheckItems[0].SanitizedItem, &itemSanitizado)
}

func (paycheck) testWhenPaycheckAlreadyExists(t *testing.T) {
	p, pi := paychecks()
	err := postgresDb.StorePaychecks(p, pi)

	assert.Nil(t, err)
}

func (paycheck) testWhenPaycheckItemsNotExist(t *testing.T) {
	truncateTables()
	p, _ := paychecks()
	err := postgresDb.StorePaychecks(p, *new([]models.PaycheckItem))

	assert.Nil(t, err)
}

func (paycheck) testGetPaychecks(t *testing.T) {
	p, _ := paychecks()
	ps, err := postgresDb.GetPaychecks(models.Agency{ID: "tjal"}, 2023)
	if err != nil {
		t.Fatalf("error GetPaychecks(): %v", err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ps))
	assert.Equal(t, 2023, ps[0].Year)
	assert.Equal(t, p[0], ps[0])
}

func (paycheck) testGetPaycheckItems(t *testing.T) {
	_, pi := paychecks()
	pis, err := postgresDb.GetPaycheckItems(models.Agency{ID: "tjal"}, 2023)
	if err != nil {
		t.Fatalf("error GetPaycheckItems(): %v", err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 3, len(pis))
	assert.Equal(t, 2023, pis[0].Year)
	assert.Equal(t, pi[0], pis[0])
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

func insertRemunerations(remunerations []models.Remunerations) error {
	for _, remuneration := range remunerations {
		remunerationDto := dto.NewRemunerationsDTO(remuneration)
		tx := postgresDb.db.Model(dto.RemunerationsDTO{}).Create(remunerationDto)
		if tx.Error != nil {
			return fmt.Errorf("error inserting remuneration: %q", tx.Error)
		}
	}
	return nil
}

func truncateTables() error {
	tx := postgresDb.db.Exec(`TRUNCATE TABLE coletas, remuneracoes_zips, orgaos, contracheques, remuneracoes CASCADE`)
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

func paychecks() ([]models.Paycheck, []models.PaycheckItem) {
	situation := "A"
	p := []models.Paycheck{
		{
			ID:           1,
			Agency:       "tjal",
			Month:        5,
			Year:         2023,
			CollectKey:   "tjal/05/2023",
			Name:         "nome",
			RegisterID:   "123",
			Role:         "funcao",
			Workplace:    "local de trabalho",
			Salary:       1000,
			Benefits:     1200,
			Discounts:    200,
			Remuneration: 2000,
			Situation:    &situation,
		},
	}
	itemSanitizado := []string{"subsidio", "descontos diversos"}
	pi := []models.PaycheckItem{
		{
			ID:            1,
			PaycheckID:    1,
			Agency:        "tjal",
			Month:         5,
			Year:          2023,
			Type:          "R/B",
			Category:      "contracheque",
			Item:          "subsídio",
			Value:         1000,
			Inconsistent:  false,
			SanitizedItem: &itemSanitizado[0],
		},
		{
			ID:           2,
			PaycheckID:   1,
			Agency:       "tjal",
			Month:        5,
			Year:         2023,
			Type:         "R/O",
			Category:     "indenizações",
			Item:         "0",
			Value:        1200,
			Inconsistent: true,
		},
		{
			ID:            3,
			PaycheckID:    1,
			Agency:        "tjal",
			Month:         5,
			Year:          2023,
			Type:          "D",
			Category:      "contracheque",
			Item:          "descontos diversos",
			Value:         200,
			Inconsistent:  false,
			SanitizedItem: &itemSanitizado[1],
		},
	}

	return p, pi
}
