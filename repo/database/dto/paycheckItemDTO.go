package dto

import "github.com/dadosjusbr/storage/models"

type PaycheckItemDTO struct {
	ID            int     `gorm:"column:id"`
	PaycheckID    int     `gorm:"column:id_contracheque"`
	Agency        string  `gorm:"column:orgao"`
	Month         int     `gorm:"column:mes"`
	Year          int     `gorm:"column:ano"`
	Type          string  `gorm:"column:tipo"`
	Category      string  `gorm:"column:categoria"`
	Item          string  `gorm:"column:item"`
	Value         float64 `gorm:"column:valor"`
	Inconsistent  bool    `gorm:"column:inconsistente"`
	SanitizedItem *string `gorm:"column:item_sanitizado"`
}

func (PaycheckItemDTO) TableName() string {
	return "remuneracoes"
}

func (r PaycheckItemDTO) ConvertToModel() *models.PaycheckItem {
	return &models.PaycheckItem{
		ID:            r.ID,
		PaycheckID:    r.PaycheckID,
		Agency:        r.Agency,
		Month:         r.Month,
		Year:          r.Year,
		Type:          r.Type,
		Category:      r.Category,
		Item:          r.Item,
		Value:         r.Value,
		Inconsistent:  r.Inconsistent,
		SanitizedItem: r.SanitizedItem,
	}
}

func NewPaycheckItemDTO(r models.PaycheckItem) *PaycheckItemDTO {
	return &PaycheckItemDTO{
		ID:            r.ID,
		PaycheckID:    r.PaycheckID,
		Agency:        r.Agency,
		Month:         r.Month,
		Year:          r.Year,
		Type:          r.Type,
		Category:      r.Category,
		Item:          r.Item,
		Value:         r.Value,
		Inconsistent:  r.Inconsistent,
		SanitizedItem: r.SanitizedItem,
	}
}
