package dto

import "github.com/dadosjusbr/storage/models"

type PerCapitaData struct {
	AgencyID           string  `gorm:"column:orgao"`
	Year               int     `gorm:"column:ano"`
	BaseRemuneration   float64 `gorm:"column:salario"`
	OtherRemunerations float64 `gorm:"column:beneficios"`
	Discounts          float64 `gorm:"column:descontos"`
	Remunerations      float64 `gorm:"column:remuneracao"`
}

func (PerCapitaData) TableName() string {
	return "media_por_membro"
}

func (a *PerCapitaData) ConvertToModel() *models.PerCapitaData {
	return &models.PerCapitaData{
		AgencyID:           a.AgencyID,
		Year:               a.Year,
		BaseRemuneration:   a.BaseRemuneration,
		OtherRemunerations: a.OtherRemunerations,
		Discounts:          a.Discounts,
		Remunerations:      a.Remunerations,
	}
}
