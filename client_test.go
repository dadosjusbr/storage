package storage_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/models"
	"github.com/dadosjusbr/storage/repo/database"
	"github.com/dadosjusbr/storage/repo/file_storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetStateAgencies(t *testing.T) {
	tests := getStateAgencies{}
	t.Run("Test GetStateAgencies when repository return agencies", tests.testWhenRepositoryReturnAgencies)
	t.Run("Test GetStateAgencies when database connection fails", tests.testWhenRepositoryReturnError)
	t.Run("Test GetStateAgencies when repository return empty array", tests.testWhenRepositoryReturnEmptyArray)
}

type getStateAgencies struct{}

func (getStateAgencies) testWhenRepositoryReturnAgencies(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	tjsp := models.Agency{
		ID:     "tjsp",
		Name:   "Tribunal de Justiça do Estado de São Paulo",
		Type:   "Estadual",
		Entity: "Tribunal",
		UF:     "SP",
	}
	mpsp := models.Agency{
		ID:     "mpsp",
		Name:   "Ministério Público do Estado de São Paulo",
		Type:   "Estadual",
		Entity: "Ministério",
		UF:     "SP",
	}
	agencies := []models.Agency{tjsp, mpsp}
	uf := "SP"

	dbMock.EXPECT().GetStateAgencies(uf).Return(agencies, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)

	returnedAgencies, err := client.GetStateAgencies(uf)

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
}

func (getStateAgencies) testWhenRepositoryReturnError(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	repoErr := errors.New("error getting agencies")
	dbMock.EXPECT().GetStateAgencies("SP").Return(nil, repoErr)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	returnedAgencies, err := client.GetStateAgencies("SP")
	expectedErr := errors.New(fmt.Sprintf("GetStateAgencies() error: \"%s\"", repoErr.Error()))

	assert.Equal(t, expectedErr, err)
	assert.Nil(t, returnedAgencies)
}

func (getStateAgencies) testWhenRepositoryReturnEmptyArray(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	agencies := []models.Agency{}
	dbMock.EXPECT().GetStateAgencies("SP").Return(agencies, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	returnedAgencies, err := client.GetStateAgencies("SP")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
}

func TestGetOPJ(t *testing.T) {
	tests := getOPJ{}
	t.Run("Test GetOPJ when repository return agencies", tests.testWhenRepositoryReturnAgencies)
	t.Run("Test GetOPJ when database connection fails", tests.testWhenRepositoryReturnError)
	t.Run("Test GetOPJ when repository return empty array", tests.testWhenRepositoryReturnEmptyArray)
}

type getOPJ struct{}

func (getOPJ) testWhenRepositoryReturnAgencies(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	tjsp := models.Agency{
		ID:     "tjsp",
		Name:   "Tribunal de Justiça do Estado de São Paulo",
		Type:   "Estadual",
		Entity: "Tribunal",
		UF:     "SP",
	}
	tjal := models.Agency{
		ID:     "tjal",
		Name:   "Tribunal de Justiça do Estado de Alagoas",
		Type:   "Estadual",
		Entity: "Tribunal",
		UF:     "AL",
	}
	agencies := []models.Agency{tjsp, tjal}
	group := "Estadual"

	dbMock.EXPECT().GetOPJ(group).Return(agencies, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)

	returnedAgencies, err := client.GetOPJ(group)

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
}

func (getOPJ) testWhenRepositoryReturnError(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	repoErr := errors.New("error getting agencies")
	dbMock.EXPECT().GetOPJ("Estadual").Return(nil, repoErr)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	returnedAgencies, err := client.GetOPJ("Estadual")
	expectedErr := errors.New(fmt.Sprintf("GetOPJ() error: \"%s\"", repoErr.Error()))

	assert.Equal(t, expectedErr, err)
	assert.Nil(t, returnedAgencies)
}

func (getOPJ) testWhenRepositoryReturnEmptyArray(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	agencies := []models.Agency{}
	dbMock.EXPECT().GetOPJ("Estadual").Return(agencies, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	returnedAgencies, err := client.GetOPJ("Estadual")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
}

func TestGetFirstDateWithMonthlyInfo(t *testing.T) {
	tests := getFirstDateWithMonthlyInfo{}
	t.Run("Test GetFirstDateWithMonthlyInfo when repository return date", tests.testWhenRepositoryReturnDate)
	t.Run("Test GetFirstDateWithMonthlyInfo when repository return error", tests.testWhenRepositoryReturnError)
}

type getFirstDateWithMonthlyInfo struct{}

func TestGetLastDateWithMonthlyInfo(t *testing.T) {
	tests := getLastDateWithMonthlyInfo{}
	t.Run("Test GetLastDateWithMonthlyInfo when repository return date", tests.testWhenRepositoryReturnDate)
	t.Run("Test GetLastDateWithMonthlyInfo when repository return error", tests.testWhenRepositoryReturnError)
}

type getLastDateWithMonthlyInfo struct{}

func (getLastDateWithMonthlyInfo) testWhenRepositoryReturnDate(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	expecMonth := 12
	expecYear := 2022
	dbMock.EXPECT().GetLastDateWithMonthlyInfo().Return(expecMonth, expecYear, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	month, year, err := client.GetLastDateWithMonthlyInfo()

	assert.Nil(t, err)
	assert.Equal(t, month, expecMonth)
	assert.Equal(t, year, expecYear)
}

func (getFirstDateWithMonthlyInfo) testWhenRepositoryReturnDate(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	expecMonth := 1
	expecYear := 2018
	dbMock.EXPECT().GetFirstDateWithMonthlyInfo().Return(expecMonth, expecYear, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	month, year, err := client.GetFirstDateWithMonthlyInfo()

	assert.Nil(t, err)
	assert.Equal(t, month, expecMonth)
	assert.Equal(t, year, expecYear)
}

func (getLastDateWithMonthlyInfo) testWhenRepositoryReturnError(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	repoErr := errors.New("error getting last date")
	dbMock.EXPECT().GetLastDateWithMonthlyInfo().Return(0, 0, repoErr)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	month, year, err := client.GetLastDateWithMonthlyInfo()
	expectedErr := errors.New(fmt.Sprintf("GetLastDateWithMonthlyInfo() error: \"%s\"", repoErr.Error()))

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 0, month)
	assert.Equal(t, 0, year)
}

func (getFirstDateWithMonthlyInfo) testWhenRepositoryReturnError(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	repoErr := errors.New("error getting first date")
	dbMock.EXPECT().GetFirstDateWithMonthlyInfo().Return(0, 0, repoErr)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	month, year, err := client.GetFirstDateWithMonthlyInfo()
	expectedErr := errors.New(fmt.Sprintf("GetFirstDateWithMonthlyInfo() error: \"%s\"", repoErr.Error()))

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 0, month)
	assert.Equal(t, 0, year)
}

func TestGetNumberOfMonthsCollected(t *testing.T) {
	tests := getNumberOfMonthsCollected{}
	t.Run("Test GetNumberOfMonthsCollected when repository return number of months", tests.testWhenRepositoryReturnNumberOfMonths)
	t.Run("Test GetNumberOfMonthsCollected when database connection fails", tests.testWhenRepositoryReturnError)
}

type getNumberOfMonthsCollected struct{}

func (getNumberOfMonthsCollected) testWhenRepositoryReturnNumberOfMonths(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	count := 200
	dbMock.EXPECT().GetNumberOfMonthsCollected().Return(count, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)

	returnedCount, err := client.GetNumberOfMonthsCollected()

	assert.Nil(t, err)
	assert.Equal(t, count, returnedCount)
}

func (getNumberOfMonthsCollected) testWhenRepositoryReturnError(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	repoErr := errors.New("error getting number of months")
	dbMock.EXPECT().GetNumberOfMonthsCollected().Return(0, repoErr)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	returnedMonths, err := client.GetNumberOfMonthsCollected()
	expectedErr := errors.New(fmt.Sprintf("GetNumberOfMonthsCollected() error: \"%s\"", repoErr.Error()))

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 0, returnedMonths)
}

func TestGetAgenciesCount(t *testing.T) {
	tests := getAgenciesCount{}
	t.Run("Test GetAgenciesCount when repository return agencies count", tests.testWhenRepositoryReturnAgenciesCount)
	t.Run("Test GetAgenciesCount when database connection fails", tests.testWhenRepositoryReturnError)
}

type getAgenciesCount struct{}

func (getAgenciesCount) testWhenRepositoryReturnAgenciesCount(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	agenciesCount := 3
	dbMock.EXPECT().GetAgenciesCount().Return(agenciesCount, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)

	returnedAgenciesCount, err := client.GetAgenciesCount()

	assert.Nil(t, err)
	assert.Equal(t, agenciesCount, returnedAgenciesCount)
}

func (getAgenciesCount) testWhenRepositoryReturnError(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockInterface(mockCrl)
	fsMock := file_storage.NewMockInterface(mockCrl)

	repoErr := errors.New("error getting agencies count")
	dbMock.EXPECT().GetAgenciesCount().Return(0, repoErr)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	returnedAgenciesCount, err := client.GetAgenciesCount()
	expectedErr := errors.New(fmt.Sprintf("GetAgenciesCount() error: \"%s\"", repoErr.Error()))

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 0, returnedAgenciesCount)
}
