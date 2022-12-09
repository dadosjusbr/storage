package dto

import "github.com/dadosjusbr/storage/models"

type GeneralMonthlyInfoDTO struct {
	Month              int     `gorm:"column:mes"`
	Count              int     `gorm:"column:count"`
	BaseRemuneration   float64 `gorm:"column:remuneracao_base"`
	OtherRemunerations float64 `gorm:"column:outras_remuneracoes"`
}

func (g GeneralMonthlyInfoDTO) ConvertToModel() *models.GeneralMonthlyInfo {
	return &models.GeneralMonthlyInfo{
		Month:              g.Month,
		Count:              g.Count,
		BaseRemuneration:   g.BaseRemuneration,
		OtherRemunerations: g.OtherRemunerations,
	}
}

func NewGeneralMonthlyInfoDTO(g models.GeneralMonthlyInfo) *GeneralMonthlyInfoDTO {
	return &GeneralMonthlyInfoDTO{
		Month:              g.Month,
		Count:              g.Count,
		BaseRemuneration:   g.BaseRemuneration,
		OtherRemunerations: g.OtherRemunerations,
	}
}
