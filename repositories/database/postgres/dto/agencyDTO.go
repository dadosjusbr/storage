package dto

import (
	"github.com/dadosjusbr/storage/models"
)

// Agency A Struct containing the main descriptions of each Agency.
type AgencyDTO struct {
	ID     string `gorm:"column:id"`
	Name   string `gorm:"column:nome"`
	Type   string `gorm:"column:jurisdição"`
	Entity string `gorm:"column:entidade"`
	UF     string `gorm:"column:uf"`
	//TODO: Add Collecting
}

func (AgencyDTO) TableName() string {
	return "orgaos"
}

func (a AgencyDTO) ConvertToModel() (*models.Agency, error) {
	return &models.Agency{
		ID:     a.ID,
		Name:   a.Name,
		Type:   a.Type,
		Entity: a.Entity,
		UF:     a.UF,
	}, nil
}

func NewAgencyDTO(agency models.Agency) (*AgencyDTO, error) {
	return &AgencyDTO{
		ID:     agency.ID,
		Name:   agency.Name,
		Type:   agency.Type,
		Entity: agency.Entity,
		UF:     agency.UF,
	}, nil
}
