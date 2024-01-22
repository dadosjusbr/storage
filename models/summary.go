package models

// Summaries contains all summary detailed information
type Summaries struct {
	General       Summary `json:"general,omitempty"`
	MemberActive  Summary `json:"memberactive,omitempty"`
	Undefined     Summary `json:"undefined,omitempty"`
	ServantActive Summary `json:"servantactive,omitempty"`
}

// Summary A Struct containing summarized  information about a agency/month stats
type Summary struct {
	Count              int         `json:"membros,omitempty"`             // Number of employees
	BaseRemuneration   DataSummary `json:"remuneracao_base,omitempty"`    //  Statistics (Max, Min, Median, Total)
	OtherRemunerations DataSummary `json:"outras_remuneracoes,omitempty"` //  Statistics (Max, Min, Median, Total)
	Discounts          DataSummary `json:"descontos,omitempty"`           //  Statistics (Max, Min, Median, Total)
	Remunerations      DataSummary `json:"remuneracoes,omitempty"`        //  Statistics (Max, Min, Median, Total)
	IncomeHistogram    map[int]int `json:"histograma_renda,omitempty"`
	ItemSummary        ItemSummary `json:"resumo_rubricas,omitempty"` // resumo de valores agregadado por rubrica
}

// DataSummary A Struct containing data summary with statistics.
type DataSummary struct {
	Max     float64 `json:"maximo,omitempty"`
	Min     float64 `json:"minimo,omitempty"`
	Average float64 `json:"media,omitempty"`
	Total   float64 `json:"total,omitempty"`
}

type ItemSummary struct {
	FoodAllowance float64 `json:"auxilio_alimentacao,omitempty"`
	BonusLicense  float64 `json:"licenca_premio,omitempty"`
	Others        float64 `json:"outras,omitempty"` // valor agregado de outras rubricas não identificadas
}
