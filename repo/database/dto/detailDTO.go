package dto

import "github.com/dadosjusbr/storage/models"

// Detail A struct contains a summary of the agency's indices and their metadata
type Detail struct {
	ID    string `gorm:"column:id_orgao"`
	Month int    `gorm:"column:mes"`
	Year  int    `gorm:"column:ano"`
	Score
	Meta
}

func (Detail) TableName() string {
	return "coletas"
}

func (d *Detail) ConvertToModel() *models.Detail {
	return &models.Detail{
		Month: d.Month,
		Year:  d.Year,
		Score: &models.Score{
			Score:             d.Score.Score,
			CompletenessScore: d.Score.CompletenessScore,
			EasinessScore:     d.Score.EasinessScore,
		},
		Meta: &models.Meta{
			OpenFormat:       d.Meta.OpenFormat,
			Expenditure:      d.Meta.Expenditure,
			Access:           d.Meta.Access,
			Extension:        d.Meta.Extension,
			StrictlyTabular:  d.Meta.StrictlyTabular,
			ConsistentFormat: d.Meta.ConsistentFormat,
			HaveEnrollment:   d.Meta.HaveEnrollment,
			ThereIsACapacity: d.Meta.ThereIsACapacity,
			HasPosition:      d.Meta.HasPosition,
			BaseRevenue:      d.Meta.BaseRevenue,
			OtherRecipes:     d.Meta.OtherRecipes,
		},
	}
}
