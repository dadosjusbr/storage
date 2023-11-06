package dto

import "github.com/dadosjusbr/storage/models"

// IndexInformation A struct contains a summary of the agency's indexes and their metadata
type IndexInformation struct {
	ID    string `gorm:"column:id_orgao"`
	Month int    `gorm:"column:mes"`
	Year  int    `gorm:"column:ano"`
	Score
	Meta
	Type string `gorm:"column:jurisdicao"`
}

func (IndexInformation) TableName() string {
	return "coletas"
}

func (d *IndexInformation) ConvertToModel() *models.IndexInformation {
	return &models.IndexInformation{
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
		Type: d.Type,
	}
}
