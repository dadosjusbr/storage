package models

type RetroactivePayments struct {
	ID                 int     `json:"id,omitempty"`
	PaycheckID         int     `json:"id_contracheque,omitempty"`
	Agency             string  `json:"orgao,omitempty"`
	Month              int     `json:"mes,omitempty"`
	Year               int     `json:"ano,omitempty"`
	Name               string  `json:"nome,omitempty"`
	RegisterID         string  `json:"matricula,omitempty"`
	Role               string  `json:"funcao,omitempty"`
	Workplace          string  `json:"local_trabalho,omitempty"`
	ProcessNumber      string  `json:"numero_processo,omitempty"`
	ProcessObject      string  `json:"objeto_processo,omitempty"`
	ProcessOrigin      string  `json:"origem_processo,omitempty"`
	Value              float64 `json:"valor_bruto,omitempty"`
	SocialContribution float64 `json:"contribuicao_previdenciaria,omitempty"`
	IncomeTax          float64 `json:"imposto_de_renda,omitempty"`
	SalaryCapDeduction float64 `json:"abate_teto,omitempty"`
	Discounts          float64 `json:"descontos,omitempty"`
	NetValue           float64 `json:"valor_liquido,omitempty"`
	SanitizedName      string  `json:"nome_sanitizado,omitempty"`
}
