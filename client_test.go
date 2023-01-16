package storage_test

import (
	"testing"

	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/mocks"
	"github.com/dadosjusbr/storage/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func GetOPE(t *testing.T) {
	tests := getOPE{}
	t.Run("Test GetOPE when all agencies are returned", tests.testWhenAllAgenciesAreReturned)
}

type getOPE struct{}

func (getOPE) testWhenAllAgenciesAreReturned(t *testing.T) {
	mockCrl := gomock.NewController(t)
	dbMock := mocks.NewMockIDatabaseRepository(mockCrl)
	fsMock := mocks.NewMockIStorageRepository(mockCrl)

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
	tjpb := models.Agency{
		ID:     "tjpb",
		Name:   "Tribunal de Justiça do Estado de Pernambuco",
		Type:   "Estadual",
		Entity: "Tribunal",
		UF:     "PB",
	}
	agencies := []models.Agency{tjsp, mpsp, tjpb}

	dbMock.EXPECT().GetOPE("SP", 2018).Return(agencies, nil)

	client, _ := storage.NewClient(dbMock, fsMock)

	returnedAgencies, err := client.GetOPE("SP", 2018)

	assert.Nil(t, err)
	assert.Equal(t, agencies, returnedAgencies)
}
