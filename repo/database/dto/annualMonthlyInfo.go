package dto

import (
	"github.com/dadosjusbr/storage/models"
)

type AnnualMonthlyInfoDTO struct {
	Year               int     `gorm:"column:ano"`
	Count              int     `gorm:"column:num_membros"`
	BaseRemuneration   float64 `gorm:"column:remuneracao_base"`
	OtherRemunerations float64 `gorm:"column:outras_remuneracoes"`
}

func NewAnnualMonthlyInfoDTO(ami models.AnnualMonthlyInfo) *AnnualMonthlyInfoDTO {
	return &AnnualMonthlyInfoDTO{
		Year:               ami.Year,
		Count:              ami.Count,
		BaseRemuneration:   ami.BaseRemuneration,
		OtherRemunerations: ami.OtherRemunerations,
	}
}

func (ami *AnnualMonthlyInfoDTO) ConvertToModel() *models.AnnualMonthlyInfo {
	return &models.AnnualMonthlyInfo{
		Year:               ami.Year,
		Count:              ami.Count,
		BaseRemuneration:   ami.BaseRemuneration,
		OtherRemunerations: ami.OtherRemunerations,
	}
}
