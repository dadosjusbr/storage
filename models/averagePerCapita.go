package models

type AveragePerCapita struct {
	ID                          string  `json:"orgao,omitempty"`
	Year                        int     `json:"ano,omitempty"`
	BaseRemunerationPerCapita   float64 `json:"remuneracao_base,omitempty"`
	OtherRemunerationsPerCapita float64 `json:"outras_remuneracoes,omitempty"`
	DiscountsPerCapita          float64 `json:"descontos,omitempty"`
	RemunerationsPerCapita      float64 `json:"remuneracoes,omitempty"`
}
