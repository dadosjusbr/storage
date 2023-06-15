package dto

import (
	"github.com/dadosjusbr/storage/models"
)

type AnnualSummaryDTO struct {
	Year               int     `gorm:"column:ano"`
	AverageCount       int     `gorm:"column:media_num_membros"`
	TotalCount         int     `gorm:"column:total_num_membros"`
	BaseRemuneration   float64 `gorm:"column:remuneracao_base"`
	OtherRemunerations float64 `gorm:"column:outras_remuneracoes"`
	Discounts          float64 `gorm:"column:descontos"`
	NumMonthsWithData  int     `gorm:"column:meses_com_dados"`
}

func NewAnnualSummaryDTO(ami models.AnnualSummary) *AnnualSummaryDTO {
	return &AnnualSummaryDTO{
		Year:               ami.Year,
		AverageCount:       ami.AverageCount,
		TotalCount:         ami.TotalCount,
		BaseRemuneration:   ami.BaseRemuneration,
		OtherRemunerations: ami.OtherRemunerations,
		Discounts:          ami.Discounts,
		NumMonthsWithData:  ami.NumMonthsWithData,
	}
}

func (ami *AnnualSummaryDTO) ConvertToModel() *models.AnnualSummary {
	return &models.AnnualSummary{
		Year:               ami.Year,
		AverageCount:       ami.AverageCount,
		TotalCount:         ami.TotalCount,
		BaseRemuneration:   ami.BaseRemuneration,
		OtherRemunerations: ami.OtherRemunerations,
		Discounts:          ami.Discounts,
		NumMonthsWithData:  ami.NumMonthsWithData,
	}
}
