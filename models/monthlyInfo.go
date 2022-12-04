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
	AgencyID          string                 `json:"aid,omitempty" bson:"aid,omitempty"`
	Month             int                    `json:"month,omitempty" bson:"month,omitempty"`
	Year              int                    `json:"year,omitempty" bson:"year,omitempty"`
	Backups           []Backup               `json:"backups,omitempty" bson:"backups,omitempty"`
	Summary           Summary                `json:"summary,omitempty" bson:"summary,omitempty"`
	CrawlerVersion    string                 `json:"crawler_version,omitempty" bson:"crawler_version,omitempty"`
	CrawlerRepo       string                 `json:"crawler_repo,omitempty" bson:"crawler_repo,omitempty"` // The github Repository of MI Crawler
	ParserRepo        string                 `json:"parser_repo,omitempty" bson:"parser_repo,omitempty"`   // The github Repository of MI Parser
	ParserVersion     string                 `json:"parser_version,omitempty" bson:"parser_version,omitempty"`
	CrawlingTimestamp *timestamppb.Timestamp `json:"crawling_ts,omitempty" bson:"crawling_ts,omitempty"` // Crawling moment (always UTC)
	ProcInfo          *coleta.ProcInfo       `json:"procinfo,omitempty" bson:"procinfo,omitempty"`       // Making this a pointer because it should be an optional field due to backwards compatibility.
	Package           *Backup                `json:"package,omitempty" bson:"package,omitempty"`         // Making this a pointer because it should be an optional field due to backwards compatibility.
	Meta              *Meta                  `json:"meta,omitempty" bson:"meta,omitempty"`
	Score             *Score                 `json:"score,omitempty" bson:"score,omitempty"`
	ExectionTime      float64                `json:"exection_time,omitempty" bson:"exection_time,omitempty"`
}

type Meta struct {
	OpenFormat       bool   `json:"open_format,omitempty" bson:"open_format,omitempty"`
	Access           string `json:"access,omitempty" bson:"access,omitempty"`
	Extension        string `json:"extension,omitempty" bson:"extension,omitempty"`
	StrictlyTabular  bool   `json:"strictly_tabular,omitempty" bson:"strictly_tabular,omitempty"`
	ConsistentFormat bool   `json:"consistent_format,omitempty" bson:"consistent_format,omitempty"`
	HaveEnrollment   bool   `json:"have_enrollment,omitempty" bson:"have_enrollment,omitempty"`
	ThereIsACapacity bool   `json:"there_is_a_capacity,omitempty" bson:"there_is_a_capacity,omitempty"`
	HasPosition      bool   `json:"has_position,omitempty" bson:"has_position,omitempty"`
	BaseRevenue      string `json:"base_revenue,omitempty" bson:"base_revenue,omitempty"`
	OtherRecipes     string `json:"other_recipes,omitempty" bson:"other_recipes,omitempty"`
	Expenditure      string `json:"expenditure,omitempty" bson:"expenditure,omitempty"`
}

type Score struct {
	Score             float64 `json:"score,omitempty" bson:"score,omitempty"`
	CompletenessScore float64 `json:"completeness_score,omitempty" bson:"completeness_score,omitempty"`
	EasinessScore     float64 `json:"easiness_score,omitempty" bson:"easiness_score,omitempty"`
}

// MonthlyInfoVersion é um item do histórico de coletas armazenado no banco de dados.
type MonthlyInfoVersion struct {
	AgencyID  string            `json:"aid,omitempty" bson:"aid,omitempty"`
	Month     int               `json:"month,omitempty" bson:"month,omitempty"`
	Year      int               `json:"year,omitempty" bson:"year,omitempty"`
	VersionID int64             `json:"version_id,omitempty" bson:"version_id,omitempty"` // revisão/versão do irem. O tipo é int64 pois podemos querer usar epoch para ficar mais simples.
	Version   AgencyMonthlyInfo `json:"version,omitempty" bson:"version,omitempty"`
}

//the GeneralMonthlyInfo is used to struct the agregation used to get the remuneration info from all angencies in a given month
type GeneralMonthlyInfo struct {
	Month              int     `json:"_id,omitempty" bson:"_id,omitempty"`
	Count              int     `json:"count" bson:"count,omitempty"`                             // Number of employees
	BaseRemuneration   float64 `json:"base_remuneration" bson:"base_remuneration,omitempty"`     //  Statistics (Max, Min, Median, Total)
	OtherRemunerations float64 `json:"other_remunerations" bson:"other_remunerations,omitempty"` //  Statistics (Max, Min, Median, Total)
}

type RemmunerationSummary struct {
	Count int
	Value float64
}

type Remunerations struct {
	AgencyID     string `json:"aid,omitempty" bson:"-"`
	Year         int    `json:"year,omitempty" bson:"-"`
	Month        int    `json:"month,omitempty" bson:"-"`
	NumBase      int    `json:"num_base,omitempty" bson:"-"`
	NumDiscounts int    `json:"num_descontos,omitempty" bson:"-"`
	NumOther     int    `json:"num_outras,omitempty" bson:"-"`
	ZipUrl       string `json:"zip_url,omitempty" bson:"-"`
}
