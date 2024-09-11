package models

type PerCapitaData struct {
	AgencyID           string  `json:"orgao,omitempty"`
	Year               int     `json:"ano,omitempty"`
	BaseRemuneration   float64 `json:"remuneracao_base,omitempty"`
	OtherRemunerations float64 `json:"outras_remuneracoes,omitempty"`
	Discounts          float64 `json:"descontos,omitempty"`
	Remunerations      float64 `json:"remuneracoes,omitempty"`
}
