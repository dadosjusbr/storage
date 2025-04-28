package dto

import (
	"github.com/dadosjusbr/storage/models"
)

type GeneralMonthlyInfoDTO struct {
	Month              int                `gorm:"column:mes"`
	Count              int                `gorm:"column:num_membros"`
	BaseRemuneration   float64            `gorm:"column:remuneracao_base"`
	OtherRemunerations float64            `gorm:"column:outras_remuneracoes"`
	Discounts          float64            `gorm:"column:descontos"`
	Remunerations      float64            `gorm:"column:remuneracoes"`
	ItemSummary        map[string]float64 `gorm:"-" json:"item_summary"`
}

type ItemSummary struct {
	FoodAllowance        float64 `gorm:"column:auxilio_alimentacao"`
	BonusLicense         float64 `gorm:"column:licenca_premio"`
	VacationCompensation float64 `gorm:"column:indenizacao_de_ferias"`
	Vacation             float64 `gorm:"column:ferias"`
	ChristmasBonus       float64 `gorm:"column:gratificacao_natalina"`
	CompensatoryLicense  float64 `gorm:"column:licenca_compensatoria"`
	HealthAllowance      float64 `gorm:"column:auxilio_saude"`
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
		ItemSummary:        gmi.ItemSummary,
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
		ItemSummary:        gmi.ItemSummary,
	}
}
