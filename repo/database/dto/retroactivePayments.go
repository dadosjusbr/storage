package dto

import "github.com/dadosjusbr/storage/models"

type RetroactivePaymentsDTO struct {
	ID                 int     `gorm:"column:id"`
	PaycheckID         int     `gorm:"column:id_contracheque"`
	Agency             string  `gorm:"column:orgao"`
	Month              int     `gorm:"column:mes"`
	Year               int     `gorm:"column:ano"`
	Name               string  `gorm:"column:nome"`
	RegisterID         string  `gorm:"column:matricula"`
	Role               string  `gorm:"column:funcao"`
	Workplace          string  `gorm:"column:local_trabalho"`
	ProcessNumber      string  `gorm:"column:numero_processo"`
	ProcessObject      string  `gorm:"column:objeto_processo"`
	ProcessOrigin      string  `gorm:"column:origem_processo"`
	Value              float64 `gorm:"column:valor_bruto"`
	SocialContribution float64 `gorm:"column:contribuicao_previdenciaria"`
	IncomeTax          float64 `gorm:"column:imposto_de_renda"`
	SalaryCapDeduction float64 `gorm:"column:abate_teto"`
	Discounts          float64 `gorm:"column:descontos"`
	NetValue           float64 `gorm:"column:valor_liquido"`
	SanitizedName      string  `gorm:"column:nome_sanitizado"`
}

func (RetroactivePaymentsDTO) TableName() string {
	return "retroativos"
}

func (p RetroactivePaymentsDTO) ConvertToModel() *models.RetroactivePayments {
	return &models.RetroactivePayments{
		ID:                 p.ID,
		PaycheckID:         p.PaycheckID,
		Agency:             p.Agency,
		Month:              p.Month,
		Year:               p.Year,
		Name:               p.Name,
		RegisterID:         p.RegisterID,
		Role:               p.Role,
		Workplace:          p.Workplace,
		ProcessNumber:      p.ProcessNumber,
		ProcessObject:      p.ProcessObject,
		ProcessOrigin:      p.ProcessOrigin,
		Value:              p.Value,
		SocialContribution: p.SocialContribution,
		IncomeTax:          p.IncomeTax,
		SalaryCapDeduction: p.SalaryCapDeduction,
		Discounts:          p.Discounts,
		NetValue:           p.NetValue,
		SanitizedName:      p.SanitizedName,
	}
}

func NewRetroactivePaymentsDTO(p models.RetroactivePayments) *RetroactivePaymentsDTO {
	return &RetroactivePaymentsDTO{
		ID:                 p.ID,
		PaycheckID:         p.PaycheckID,
		Agency:             p.Agency,
		Month:              p.Month,
		Year:               p.Year,
		Name:               p.Name,
		RegisterID:         p.RegisterID,
		Role:               p.Role,
		Workplace:          p.Workplace,
		ProcessNumber:      p.ProcessNumber,
		ProcessObject:      p.ProcessObject,
		ProcessOrigin:      p.ProcessOrigin,
		Value:              p.Value,
		SocialContribution: p.SocialContribution,
		IncomeTax:          p.IncomeTax,
		SalaryCapDeduction: p.SalaryCapDeduction,
		Discounts:          p.Discounts,
		NetValue:           p.NetValue,
		SanitizedName:      p.SanitizedName,
	}
}
