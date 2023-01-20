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

func TestGetOPE(t *testing.T) {
	tests := getOPE{}
	t.Run("Test GetOPE when repository return agencies", tests.testWhenRepositoryReturnAgencies)
	t.Run("Test GetOPE when database connection fails", tests.testWhenRepositoryReturnError)
	t.Run("Test GetOPE when repository return empty array", tests.testWhenRepositoryReturnEmptyArray)
}

type getOPE struct{}

func (getOPE) testWhenRepositoryReturnAgencies(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockIDatabaseRepository(mockCrl)
	fsMock := file_storage.NewMockIStorageRepository(mockCrl)

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

	dbMock.EXPECT().GetOPE(uf).Return(agencies, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)

	returnedAgencies, err := client.GetOPE(uf)

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
}

func (getOPE) testWhenRepositoryReturnError(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockIDatabaseRepository(mockCrl)
	fsMock := file_storage.NewMockIStorageRepository(mockCrl)

	repoErr := errors.New("error getting agencies")
	dbMock.EXPECT().GetOPE("SP").Return(nil, repoErr)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	returnedAgencies, err := client.GetOPE("SP")
	expectedErr := errors.New(fmt.Sprintf("GetOPE() error: \"%s\"", repoErr.Error()))

	assert.Equal(t, expectedErr, err)
	assert.Nil(t, returnedAgencies)
}

func (getOPE) testWhenRepositoryReturnEmptyArray(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := database.NewMockIDatabaseRepository(mockCrl)
	fsMock := file_storage.NewMockIStorageRepository(mockCrl)

	agencies := []models.Agency{}
	dbMock.EXPECT().GetOPE("SP").Return(agencies, nil)
	dbMock.EXPECT().Connect().Return(nil)

	client, err := storage.NewClient(dbMock, fsMock)
	returnedAgencies, err := client.GetOPE("SP")

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
}
