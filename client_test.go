package storage_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/mocks"
	"github.com/dadosjusbr/storage/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	TestGetOPE(t)
}

func TestGetOPE(t *testing.T) {
	tests := getOPE{}
	t.Run("Test GetOPE when repository return agencies", tests.testWhenRepositoryReturnAgencies)
	t.Run("Test GetOPE when database connection fails", tests.testWhenRepositoryReturnError)
	t.Run("Test GetOPE when repository return empty array", tests.testWhenRepositoryReturnEmptyArray)
}

type getOPE struct{}

func (getOPE) testWhenRepositoryReturnAgencies(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := mocks.NewMockIDatabaseRepository(mockCrl)
	fsMock := mocks.NewMockIStorageRepository(mockCrl)

	tjsp := models.Agency{
		ID:      "tjsp",
		Name:    "Tribunal de Justiça do Estado de São Paulo",
		Type:    "Estadual",
		Entity:  "Tribunal",
		UF:      "SP",
		FlagURL: "v1/orgao/tjsp",
	}
	mpsp := models.Agency{
		ID:      "mpsp",
		Name:    "Ministério Público do Estado de São Paulo",
		Type:    "Estadual",
		Entity:  "Ministério",
		UF:      "SP",
		FlagURL: "v1/orgao/mpsp",
	}
	agencies := []models.Agency{tjsp, mpsp}
	uf := "SP"
	year := 2018

	dbMock.EXPECT().GetOPE(uf, year).Return(agencies, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)

	returnedAgencies, err := client.GetOPE(uf, year)

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
}

func (getOPE) testWhenRepositoryReturnError(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := mocks.NewMockIDatabaseRepository(mockCrl)
	fsMock := mocks.NewMockIStorageRepository(mockCrl)

	repoErr := errors.New("error getting agencies")
	dbMock.EXPECT().GetOPE("SP", 2018).Return(nil, repoErr)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	returnedAgencies, err := client.GetOPE("SP", 2018)
	expectedErr := errors.New(fmt.Sprintf("GetOPE() error: \"%s\"", repoErr.Error()))

	assert.Equal(t, expectedErr, err)
	assert.Nil(t, returnedAgencies)
}

func (getOPE) testWhenRepositoryReturnEmptyArray(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := mocks.NewMockIDatabaseRepository(mockCrl)
	fsMock := mocks.NewMockIStorageRepository(mockCrl)

	agencies := []models.Agency{}
	dbMock.EXPECT().GetOPE("SP", 2018).Return(agencies, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	returnedAgencies, err := client.GetOPE("SP", 2018)

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
}
