package dto

import (
	"github.com/dadosjusbr/storage/models"
)

type GeneralMonthlyInfoDTO struct {
	Month              int         `gorm:"column:mes"`
	Count              int         `gorm:"column:num_membros"`
	BaseRemuneration   float64     `gorm:"column:remuneracao_base"`
	OtherRemunerations float64     `gorm:"column:outras_remuneracoes"`
	Discounts          float64     `gorm:"column:descontos"`
	Remunerations      float64     `gorm:"column:remuneracoes"`
	ItemSummary        ItemSummary `gorm:"embedded"`
}

type ItemSummary struct {
	FoodAllowance        float64 `gorm:"column:auxilio_alimentacao"`
	BonusLicense         float64 `gorm:"column:licenca_premio"`
	VacationCompensation float64 `gorm:"column:indenizacao_de_ferias"`
	Others               float64 `gorm:"column:outras"`
}

func NewGeneralMonthlyInfoDTO(gmi models.GeneralMonthlyInfo) *GeneralMonthlyInfoDTO {
	return &GeneralMonthlyInfoDTO{
		Month:              gmi.Month,
		Count:              gmi.Count,
		BaseRemuneration:   gmi.BaseRemuneration,
		OtherRemunerations: gmi.OtherRemunerations,
		Discounts:          gmi.Discounts,
		Remunerations:      gmi.Remunerations,
		ItemSummary: ItemSummary{
			FoodAllowance:        gmi.ItemSummary.FoodAllowance,
			BonusLicense:         gmi.ItemSummary.BonusLicense,
			VacationCompensation: gmi.ItemSummary.VacationCompensation,
			Others:               gmi.ItemSummary.Others,
		},
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
		ItemSummary: models.ItemSummary{
			FoodAllowance:        gmi.ItemSummary.FoodAllowance,
			BonusLicense:         gmi.ItemSummary.BonusLicense,
			VacationCompensation: gmi.ItemSummary.VacationCompensation,
			Others:               gmi.ItemSummary.Others,
		},
	}
}
