package dto

import "github.com/dadosjusbr/storage/models"

type AveragePerCapita struct {
	ID                          string  `gorm:"column:orgao"`
	Year                        int     `gorm:"column:ano"`
	BaseRemunerationPerCapita   float64 `gorm:"column:salario"`
	OtherRemunerationsPerCapita float64 `gorm:"column:beneficios"`
	DiscountsPerCapita          float64 `gorm:"column:descontos"`
	RemunerationsPerCapita      float64 `gorm:"column:remuneracao"`
}

func (AveragePerCapita) TableName() string {
	return "media_por_membro"
}

func (a *AveragePerCapita) ConvertToModel() *models.AveragePerCapita {
	return &models.AveragePerCapita{
		ID:                          a.ID,
		Year:                        a.Year,
		BaseRemunerationPerCapita:   a.BaseRemunerationPerCapita,
		OtherRemunerationsPerCapita: a.OtherRemunerationsPerCapita,
		DiscountsPerCapita:          a.DiscountsPerCapita,
		RemunerationsPerCapita:      a.RemunerationsPerCapita,
	}
}
