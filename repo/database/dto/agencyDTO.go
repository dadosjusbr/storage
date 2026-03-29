package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dadosjusbr/storage/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"
)

// Agency A Struct containing the main descriptions of each Agency.
type AgencyDTO struct {
	ID             string         `gorm:"column:id"`
	Name           string         `gorm:"column:nome"`
	Type           string         `gorm:"column:jurisdicao"`
	Entity         string         `gorm:"column:entidade"`
	UF             string         `gorm:"column:uf"`
	Collecting     datatypes.JSON `gorm:"column:coletando"`
	TwitterHandle  string         `gorm:"column:twitter_handle"`
	OmbudsmanURL   string         `gorm:"column:ouvidoria"`
	LastCollection time.Time      `gorm:"column:ultima_coleta"`
}

func (AgencyDTO) TableName() string {
	return "orgaos"
}

func (a AgencyDTO) ConvertToModel() (*models.Agency, error) {
	var collecting []models.Collecting
	collectingBytes, err := a.Collecting.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("error while marshaling collecting: %q", err)
	}
	err = json.Unmarshal(collectingBytes, &collecting)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling collecting: %q", err)
	}
	return &models.Agency{
		ID:             a.ID,
		Name:           a.Name,
		Type:           a.Type,
		Entity:         a.Entity,
		UF:             a.UF,
		Collecting:     collecting,
		TwitterHandle:  a.TwitterHandle,
		OmbudsmanURL:   a.OmbudsmanURL,
		LastCollection: timestamppb.New(a.LastCollection),
	}, nil
}

func NewAgencyDTO(agency models.Agency) (*AgencyDTO, error) {
	collecting, err := json.Marshal(agency.Collecting)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling collecting: %q", err)
	}
	var lastCollection time.Time
	if agency.LastCollection != nil {
		lastCollection = time.Unix(agency.LastCollection.Seconds, int64(agency.LastCollection.Nanos))
	} else {
		lastCollection = time.Now()
	}
	return &AgencyDTO{
		ID:             agency.ID,
		Name:           agency.Name,
		Type:           agency.Type,
		Entity:         agency.Entity,
		UF:             agency.UF,
		Collecting:     collecting,
		TwitterHandle:  agency.TwitterHandle,
		OmbudsmanURL:   agency.OmbudsmanURL,
		LastCollection: lastCollection,
	}, nil
}
