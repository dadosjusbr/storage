package dto

import (
	"github.com/dadosjusbr/storage/models"
)

type GeneralMonthlyInfoDTO struct {
	Month              int     `gorm:"column:mes"`
	Count              int     `gorm:"column:num_membros"`
	BaseRemuneration   float64 `gorm:"column:remuneracao_base"`
	OtherRemunerations float64 `gorm:"column:outras_remuneracoes"`
	Discounts          float64 `gorm:"column:descontos"`
	Remunerations      float64 `gorm:"column:remuneracoes"`
}

func NewGeneralMonthlyInfoDTO(gmi models.GeneralMonthlyInfo) *GeneralMonthlyInfoDTO {
	return &GeneralMonthlyInfoDTO{
		Month:              gmi.Month,
		Count:              gmi.Count,
		BaseRemuneration:   gmi.BaseRemuneration,
		OtherRemunerations: gmi.OtherRemunerations,
		Discounts:          gmi.Discounts,
		Remunerations:      gmi.Remunerations,
	}
}

func (gmi *GeneralMonthlyInfoDTO) ConvertToModel() *models.GeneralMonthlyInfo {
	return &models.GeneralMonthlyInfo{
		Month:              gmi.Month,
		Count:              gmi.Count,
		BaseRemuneration:   gmi.BaseRemuneration,
		OtherRemunerations: gmi.OtherRemunerations,
		Discounts:          gmi.Discounts,
		Remunerations:      gmi.Remunerations,
	}
}
