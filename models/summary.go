package models

// Summaries contains all summary detailed information
type Summaries struct {
	General       Summary `json:"general,omitempty" bson:"general"`
	MemberActive  Summary `json:"memberactive,omitempty" bson:"memberactive"`
	Undefined     Summary `json:"undefined,omitempty" bson:"undefined"`
	ServantActive Summary `json:"servantactive,omitempty" bson:"servantactive"`
}

// Summary A Struct containing summarized  information about a agency/month stats
type Summary struct {
	Count              int         `json:"membros" bson:"count,omitempty"`                             // Number of employees
	BaseRemuneration   DataSummary `json:"remuneracao_base" bson:"base_remuneration,omitempty"`     //  Statistics (Max, Min, Median, Total)
	OtherRemunerations DataSummary `json:"outras_remuneracoes" bson:"other_remunerations,omitempty"` //  Statistics (Max, Min, Median, Total)
	IncomeHistogram    map[int]int `json:"histograma_renda" bson:"hist,omitempty"`
}

// DataSummary A Struct containing data summary with statistics.
type DataSummary struct {
	Max     float64 `json:"maximo" bson:"max,omitempty"`
	Min     float64 `json:"minimo" bson:"min,omitempty"`
	Average float64 `json:"media" bson:"avg,omitempty"`
	Total   float64 `json:"total" bson:"total,omitempty"`
}
