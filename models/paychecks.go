package models

type Paycheck struct {
	ID           int     `json:"id,omitempty"`
	Agency       string  `json:"orgao,omitempty"`
	Month        int     `json:"mes,omitempty"`
	Year         int     `json:"ano,omitempty"`
	CollectKey   string  `json:"chave_coleta,omitempty"`
	Name         string  `json:"nome,omitempty"`
	RegisterID   string  `json:"matricula,omitempty"`
	Role         string  `json:"funcao,omitempty"`
	Workplace    string  `json:"local_trabalho,omitempty"`
	Salary       float64 `json:"salario,omitempty"`
	Benefits     float64 `json:"beneficios,omitempty"`
	Discounts    float64 `json:"descontos,omitempty"`
	Remuneration float64 `json:"remuneracao,omitempty"`
}

type PaycheckItem struct {
	ID           int     `json:"id,omitempty"`
	PaycheckID   int     `json:"id_contracheque,omitempty"`
	Agency       string  `json:"orgao,omitempty"`
	Month        int     `json:"mes,omitempty"`
	Year         int     `json:"ano,omitempty"`
	Type         string  `json:"tipo,omitempty"`
	Category     string  `json:"categoria,omitempty"`
	Item         string  `json:"item,omitempty"`
	Value        float64 `json:"valor,omitempty"`
	Inconsistent bool    `json:"inconsistente,omitempty"`
}
