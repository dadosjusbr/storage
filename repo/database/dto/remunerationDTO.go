package dto

import "github.com/dadosjusbr/storage/models"

type RemunerationDTO struct {
	ID           int     `gorm:"column:id"`
	PaycheckID   int     `gorm:"column:id_contracheque"`
	Agency       string  `gorm:"column:orgao"`
	Month        int     `gorm:"column:mes"`
	Year         int     `gorm:"column:ano"`
	Nature       string  `gorm:"column:natureza"`
	IncomeType   string  `gorm:"column:tipo_receita"`
	Category     string  `gorm:"column:categoria"`
	Item         string  `gorm:"column:item"`
	Value        float64 `gorm:"column:valor"`
	Inconsistent bool    `gorm:"column:inconsistente"`
}

func (RemunerationDTO) TableName() string {
	return "remuneracoes"
}

func (r RemunerationDTO) ConvertToModel() *models.Remuneration {
	return &models.Remuneration{
		ID:           r.ID,
		PaycheckID:   r.PaycheckID,
		Agency:       r.Agency,
		Month:        r.Month,
		Year:         r.Year,
		Nature:       r.Nature,
		IncomeType:   r.IncomeType,
		Category:     r.Category,
		Item:         r.Item,
		Value:        r.Value,
		Inconsistent: r.Inconsistent,
	}
}

func NewRemunerationDTO(r models.Remuneration) *RemunerationDTO {
	return &RemunerationDTO{
		ID:           r.ID,
		PaycheckID:   r.PaycheckID,
		Agency:       r.Agency,
		Month:        r.Month,
		Year:         r.Year,
		Nature:       r.Nature,
		IncomeType:   r.IncomeType,
		Category:     r.Category,
		Item:         r.Item,
		Value:        r.Value,
		Inconsistent: r.Inconsistent,
	}
}
