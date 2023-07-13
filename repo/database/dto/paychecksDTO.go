package dto

import "github.com/dadosjusbr/storage/models"

type PaycheckDTO struct {
	ID           int     `gorm:"column:id"`
	Agency       string  `gorm:"column:orgao"`
	Month        int     `gorm:"column:mes"`
	Year         int     `gorm:"column:ano"`
	CollectKey   string  `gorm:"column:chave_coleta"`
	Name         string  `gorm:"column:nome"`
	RegisterID   string  `gorm:"column:matricula"`
	Role         string  `gorm:"column:funcao"`
	Workplace    string  `gorm:"column:local_trabalho"`
	Salary       float64 `gorm:"column:salario"`
	Benefits     float64 `gorm:"column:beneficios"`
	Discounts    float64 `gorm:"column:descontos"`
	Remuneration float64 `gorm:"column:remuneracao"`
}

func (PaycheckDTO) TableName() string {
	return "contracheques"
}

func (p PaycheckDTO) ConvertToModel() *models.Paycheck {
	return &models.Paycheck{
		ID:           p.ID,
		Agency:       p.Agency,
		Month:        p.Month,
		Year:         p.Year,
		CollectKey:   p.CollectKey,
		Name:         p.Name,
		RegisterID:   p.RegisterID,
		Role:         p.Role,
		Workplace:    p.Workplace,
		Salary:       p.Salary,
		Benefits:     p.Benefits,
		Discounts:    p.Discounts,
		Remuneration: p.Remuneration,
	}
}

func NewPaycheckDTO(p models.Paycheck) *PaycheckDTO {
	return &PaycheckDTO{
		ID:           p.ID,
		Agency:       p.Agency,
		Month:        p.Month,
		Year:         p.Year,
		CollectKey:   p.CollectKey,
		Name:         p.Name,
		RegisterID:   p.RegisterID,
		Role:         p.Role,
		Workplace:    p.Workplace,
		Salary:       p.Salary,
		Benefits:     p.Benefits,
		Discounts:    p.Discounts,
		Remuneration: p.Remuneration,
	}
}
