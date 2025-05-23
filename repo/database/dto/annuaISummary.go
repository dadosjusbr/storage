package dto

import (
	"github.com/dadosjusbr/storage/models"
)

type AnnualSummaryDTO struct {
	Year                        int                `gorm:"column:ano"`
	AverageCount                int                `gorm:"column:media_num_membros"`
	TotalCount                  int                `gorm:"column:total_num_membros"`
	BaseRemuneration            float64            `gorm:"column:remuneracao_base"`
	OtherRemunerations          float64            `gorm:"column:outras_remuneracoes"`
	Discounts                   float64            `gorm:"column:descontos"`
	Remunerations               float64            `gorm:"column:remuneracoes"`
	BaseRemunerationPerCapita   float64            `gorm:"column:remuneracao_base_membro"`
	OtherRemunerationsPerCapita float64            `gorm:"column:outras_remuneracoes_membro"`
	DiscountsPerCapita          float64            `gorm:"column:descontos_membro"`
	RemunerationsPerCapita      float64            `gorm:"column:remuneracoes_membro"`
	NumMonthsWithData           int                `gorm:"column:meses_com_dados"`
	ItemSummary                 map[string]float64 `gorm:"-" json:"item_summary"`
	Inconsistent                bool               `gorm:"column:inconsistente;<-:false"`
}

func NewAnnualSummaryDTO(ami models.AnnualSummary) *AnnualSummaryDTO {
	return &AnnualSummaryDTO{
		Year:                        ami.Year,
		AverageCount:                ami.AverageCount,
		TotalCount:                  ami.TotalCount,
		BaseRemuneration:            ami.BaseRemuneration,
		OtherRemunerations:          ami.OtherRemunerations,
		BaseRemunerationPerCapita:   ami.BaseRemunerationPerCapita,
		OtherRemunerationsPerCapita: ami.OtherRemunerationsPerCapita,
		DiscountsPerCapita:          ami.DiscountsPerCapita,
		RemunerationsPerCapita:      ami.RemunerationsPerCapita,
		Discounts:                   ami.Discounts,
		Remunerations:               ami.Remunerations,
		NumMonthsWithData:           ami.NumMonthsWithData,
		ItemSummary:                 ami.ItemSummary,
		Inconsistent:                ami.Inconsistent,
	}
}

func (ami *AnnualSummaryDTO) ConvertToModel() *models.AnnualSummary {
	return &models.AnnualSummary{
		Year:                        ami.Year,
		AverageCount:                ami.AverageCount,
		TotalCount:                  ami.TotalCount,
		BaseRemuneration:            ami.BaseRemuneration,
		OtherRemunerations:          ami.OtherRemunerations,
		BaseRemunerationPerCapita:   ami.BaseRemunerationPerCapita,
		OtherRemunerationsPerCapita: ami.OtherRemunerationsPerCapita,
		DiscountsPerCapita:          ami.DiscountsPerCapita,
		RemunerationsPerCapita:      ami.RemunerationsPerCapita,
		Discounts:                   ami.Discounts,
		Remunerations:               ami.Remunerations,
		NumMonthsWithData:           ami.NumMonthsWithData,
		ItemSummary:                 ami.ItemSummary,
		Inconsistent:                ami.Inconsistent,
	}
}
