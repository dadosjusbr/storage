package dto

import (
	"github.com/dadosjusbr/storage/models"
)

type AnnualSummaryDTO struct {
	Year               int     `gorm:"column:ano"`
	Count              int     `gorm:"column:num_membros"`
	BaseRemuneration   float64 `gorm:"column:remuneracao_base"`
	OtherRemunerations float64 `gorm:"column:outras_remuneracoes"`
}

func NewAnnualSummaryDTO(ami models.AnnualSummary) *AnnualSummaryDTO {
	return &AnnualSummaryDTO{
		Year:               ami.Year,
		Count:              ami.Count,
		BaseRemuneration:   ami.BaseRemuneration,
		OtherRemunerations: ami.OtherRemunerations,
	}
}

func (ami *AnnualSummaryDTO) ConvertToModel() *models.AnnualSummary {
	return &models.AnnualSummary{
		Year:               ami.Year,
		Count:              ami.Count,
		BaseRemuneration:   ami.BaseRemuneration,
		OtherRemunerations: ami.OtherRemunerations,
	}
}
