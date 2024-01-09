// About used pointers.
// All pointers are important to know if in the field has information and this information is 0 or if we do not have information about that field.
// This is justified because of the use of omitempty. If a collected float64 is 0, it will not appear in the json fields, cause that's it's zero value.
// Any application consuming this data might not know if the field is really 0 or data is unavailable.
// For a example, a Funds Daily field with null will represent that we do not have that information, but a Dialy field with 0, represents that we have that information and the employee received 0 Reais in Daily Funds
// On the other hand, if we dont put pointer in those fields, Funds daily will be setted 0 as a float64 primitive number, and we will not be able to
// diferenciate if we have the 0 information or if we dont know about it.
// The point here is just to guarantee that what appears in the system are real collected data.
// As disavantage we add some complexity to code knowing that the final value will not be changed anyway.
// Use Case:
// Pointers                                 No Pointers
// daily: nil                              daily: 0
// perks: nil							   perks: 0
// total: 0								   total: 0
package models

import (
	"github.com/dadosjusbr/proto/coleta"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AgencyMonthlyInfo A Struct containing a snapshot of a agency in a month.
type AgencyMonthlyInfo struct {
	AgencyID          string                 `json:"aid,omitempty"`
	Month             int                    `json:"month,omitempty"`
	Year              int                    `json:"year,omitempty"`
	Backups           []Backup               `json:"backups,omitempty"`
	Summary           *Summary               `json:"summary,omitempty"`
	CrawlerVersion    string                 `json:"crawler_version,omitempty"`
	CrawlerRepo       string                 `json:"crawler_repo,omitempty"` // The github Repository of MI Crawler
	ParserRepo        string                 `json:"parser_repo,omitempty"`  // The github Repository of MI Parser
	ParserVersion     string                 `json:"parser_version,omitempty"`
	CrawlingTimestamp *timestamppb.Timestamp `json:"crawling_ts,omitempty"` // Crawling moment (always UTC)
	ProcInfo          *coleta.ProcInfo       `json:"procinfo,omitempty"`    // Making this a pointer because it should be an optional field due to backwards compatibility.
	Package           *Backup                `json:"package,omitempty"`     // Making this a pointer because it should be an optional field due to backwards compatibility.
	Meta              *Meta                  `json:"meta,omitempty"`
	Score             *Score                 `json:"score,omitempty"`
	Duration          float64                `json:"duration,omitempty"` // Crawling duration (seconds)
}

type Meta struct {
	OpenFormat       bool   `json:"open_format,omitempty"`
	Access           string `json:"access,omitempty"`
	Extension        string `json:"extension,omitempty"`
	StrictlyTabular  bool   `json:"strictly_tabular,omitempty"`
	ConsistentFormat bool   `json:"consistent_format,omitempty"`
	HaveEnrollment   bool   `json:"have_enrollment,omitempty"`
	ThereIsACapacity bool   `json:"there_is_a_capacity,omitempty"`
	HasPosition      bool   `json:"has_position,omitempty"`
	BaseRevenue      string `json:"base_revenue,omitempty"`
	OtherRecipes     string `json:"other_recipes,omitempty"`
	Expenditure      string `json:"expenditure,omitempty"`
}

type Score struct {
	Score             float64 `json:"score,omitempty"`
	CompletenessScore float64 `json:"completeness_score,omitempty"`
	EasinessScore     float64 `json:"easiness_score,omitempty"`
}

// MonthlyInfoVersion é um item do histórico de coletas armazenado no banco de dados.
type MonthlyInfoVersion struct {
	AgencyID  string            `json:"aid,omitempty"`
	Month     int               `json:"month,omitempty"`
	Year      int               `json:"year,omitempty"`
	VersionID int64             `json:"version_id,omitempty"` // revisão/versão do irem. O tipo é int64 pois podemos querer usar epoch para ficar mais simples.
	Version   AgencyMonthlyInfo `json:"version,omitempty"`
}

// the GeneralMonthlyInfo is used to struct the agregation used to get the remuneration info from all angencies in a given month
type GeneralMonthlyInfo struct {
	Month              int         `json:"_id,omitempty"`
	Count              int         `json:"count,omitempty"`               // Number of employees
	BaseRemuneration   float64     `json:"base_remuneration,omitempty"`   //  Statistics (Max, Min, Median, Total)
	OtherRemunerations float64     `json:"other_remunerations,omitempty"` //  Statistics (Max, Min, Median, Total)
	Discounts          float64     `json:"discounts,omitempty"`           //  Statistics (Max, Min, Median, Total)
	Remunerations      float64     `json:"remunerations,omitempty"`       //  Statistics (Max, Min, Median, Total)
	ItemSummary        ItemSummary `json:"item_summary,omitempty"`
}

type AnnualSummary struct {
	Year               int         `json:"year,omitempty"`                // Year of the data
	AverageCount       int         `json:"average_count,omitempty"`       // Average number of employees
	TotalCount         int         `json:"total_count,omitempty"`         // Total number of employees
	BaseRemuneration   float64     `json:"base_remuneration,omitempty"`   //  Statistics (Max, Min, Median, Total)
	OtherRemunerations float64     `json:"other_remunerations,omitempty"` //  Statistics (Max, Min, Median, Total)
	Discounts          float64     `json:"discounts,omitempty"`           //  Statistics (Max, Min, Median, Total)
	Remunerations      float64     `json:"remunerations,omitempty"`       //  Statistics (Max, Min, Median, Total)
	NumMonthsWithData  int         `json:"months_with_data,omitempty"`
	Package            *Backup     `json:"package,omitempty"`
	ItemSummary        ItemSummary `json:"item_summary,omitempty"`
}

type RemmunerationSummary struct {
	Count int
	Value float64
}

type Remunerations struct {
	AgencyID     string `json:"aid,omitempty"`
	Year         int    `json:"year,omitempty"`
	Month        int    `json:"month,omitempty"`
	NumBase      int    `json:"num_base,omitempty"`
	NumDiscounts int    `json:"num_descontos,omitempty"`
	NumOther     int    `json:"num_outras,omitempty"`
	ZipUrl       string `json:"zip_url,omitempty"`
}
