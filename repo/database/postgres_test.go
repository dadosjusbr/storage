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

func TestWhenOmbudsmanURLExists(t *testing.T) {
	agency := []models.Agency{
		{
			ID:           "mpmg",
			Name:         "Estadual",
			UF:           "MG",
			OmbudsmanURL: "https://aplicacao.mpmg.mp.br/ouvidoria/",
		},
	}

	if err := insertAgencies(agency); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgency, err := postgresDb.GetAgency("mpmg")

	assert.Nil(t, err)
	assert.Equal(t, agency[0].OmbudsmanURL, returnedAgency.OmbudsmanURL)
}

func TestWhenOmbudsmanURLNotExists(t *testing.T) {
	truncateTables()

	agency := []models.Agency{
		{
			ID:   "mpmg",
			Name: "Estadual",
			UF:   "MG",
		},
	}

	if err := insertAgencies(agency); err != nil {
		t.Fatalf("error inserting agencies: %q", err)
	}

	returnedAgency, err := postgresDb.GetAgency("mpmg")

	assert.Nil(t, err)
	assert.Equal(t, agency[0].OmbudsmanURL, returnedAgency.OmbudsmanURL)
}

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
	agencies := []models.Agency{
		{
			ID:   "tjsp",
			Name: "Tribunal de Justiça do Estado de São Paulo",

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
					Total: 700,
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
					Total: 2000,
				},
				OtherRemunerations: models.DataSummary{
					Total: 1000,
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
			ID:   "mpsp",
			Type: "Ministério",
			UF:   "SP",
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
			ID:   "mpsp",
			Type: "Ministério",
			UF:   "SP",
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
