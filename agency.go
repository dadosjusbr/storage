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

package storage

import (
	"github.com/dadosjusbr/proto/coleta"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Agency A Struct containing the main descriptions of each Agency.
type Agency struct {
	ID         string       `json:"aid" bson:"aid,omitempty" gorm:"column:id"`             // 'trt13'
	Name       string       `json:"name" bson:"name,omitempty" gorm:"column:nome"`         // 'Tribunal Regional do Trabalho 13° Região'
	Type       string       `json:"type" bson:"type,omitempty" gorm:"column:jurisdicao"`   // "R" for Regional, "M" for Municipal, "F" for Federal, "E" for State.
	Entity     string       `json:"entity" bson:"entity,omitempty" gorm:"column:entidade"` // "J" For Judiciário, "M" for Ministério Público, "P" for Procuradorias and "D" for Defensorias.
	UF         string       `json:"uf" bson:"uf,omitempty" gorm:"column:uf"`               // Short code for federative unity.
	FlagURL    string       `json:"url" bson:"url,omitempty"`                              // Link for state url
	Collecting []Collecting `json:"collecting" bson:"collecting,omitempty"`
}

// Collecting A Struct containing the day we checked the status of the data and the reasons why we didn't collected it.
type Collecting struct {
	Timestamp   *int64   `json:"timestamp" bson:"timestamp,omitempty"`     // Day(unix) we checked the status of the data
	Description []string `json:"description" bson:"description,omitempty"` // Reasons why we didn't collect the data
}

// AgencyMonthlyInfo A Struct containing a snapshot of a agency in a month.

type AgencyMonthlyInfo struct {
	ID                string                 `json:"id" bson:"-" gorm:"column:id"` // 'trt13/01/2020'
	AgencyID          string                 `json:"aid,omitempty" bson:"aid,omitempty" gorm:"column:id_orgao"`
	Month             int                    `json:"month,omitempty" bson:"month,omitempty" gorm:"column:mes"`
	Year              int                    `json:"year,omitempty" bson:"year,omitempty" gorm:"column:ano"`
	Actual            bool                   `json:"actual,omitempty" bson:"-" gorm:"column:atual"`
	Backups           []Backup               `json:"backups,omitempty" bson:"backups,omitempty" gorm:"-"`
	Summary           Summary                `json:"summary,omitempty" bson:"summary,omitempty" gorm:"-"`
	CrawlerID         string                 `json:"crawler_id,omitempty" bson:"crawler_id,omitempty"`
	CrawlerVersion    string                 `json:"crawler_version,omitempty" bson:"crawler_version,omitempty" gorm:"column:versao_coletor"`
	CrawlerRepo       string                 `json:"crawler_repo,omitempty" bson:"crawler_repo,omitempty" gorm:"column:repositorio_coletor"` // The github Repository of MI Crawler
	ParserRepo        string                 `json:"parser_repo,omitempty" bson:"parser_repo,omitempty" gorm:"column:repositorio_parser"`    // The github Repository of MI Parser
	ParserVersion     string                 `json:"parser_version,omitempty" bson:"parser_version,omitempty" gorm:"column:versao_parser"`
	CrawlingTimestamp *timestamppb.Timestamp `json:"crawling_ts,omitempty" bson:"crawling_ts,omitempty"`    // Crawling moment (always UTC)
	ProcInfo          *coleta.ProcInfo       `json:"procinfo,omitempty" bson:"procinfo,omitempty" gorm:"-"` // Making this a pointer because it should be an optional field due to backwards compatibility.
	Package           *Backup                `json:"package,omitempty" bson:"package,omitempty" gorm:"-"`   // Making this a pointer because it should be an optional field due to backwards compatibility.
	Meta              `json:"meta,omitempty" bson:"meta,omitempty"`
	Score             `json:"score,omitempty" bson:"score,omitempty"`
	ExectionTime      float64 `json:"exection_time,omitempty" bson:"exection_time,omitempty" gorm:"-"`
}

// Backup contains the URL to download a file and a hash to track if in the future will be changes in the file.
type Backup struct {
	URL  string `json:"url" bson:"url,omitempty"`
	Hash string `json:"hash" bson:"hash,omitempty"`
	Size int64  `json:"size" bson:"size,omitempty"`
}

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

type Meta struct {
	OpenFormat       bool   `json:"open_format,omitempty" bson:"open_format,omitempty" gorm:"column:formato_aberto"`
	Access           string `json:"access,omitempty" bson:"access,omitempty" gorm:"column:acesso"`
	Extension        string `json:"extension,omitempty" bson:"extension,omitempty" gorm:"column:extensao"`
	StrictlyTabular  bool   `json:"strictly_tabular,omitempty" bson:"strictly_tabular,omitempty" gorm:"column:estritamente_tabular"`
	ConsistentFormat bool   `json:"consistent_format,omitempty" bson:"consistent_format,omitempty" gorm:"column:formato_consistente"`
	HaveEnrollment   bool   `json:"have_enrollment,omitempty" bson:"have_enrollment,omitempty" gorm:"column:tem_matricula"`
	ThereIsACapacity bool   `json:"there_is_a_capacity,omitempty" bson:"there_is_a_capacity,omitempty" gorm:"column:tem_lotacao"`
	HasPosition      bool   `json:"has_position,omitempty" bson:"has_position,omitempty" gorm:"column:tem_cargo"`
	BaseRevenue      string `json:"base_revenue,omitempty" bson:"base_revenue,omitempty" gorm:"column:detalhamento_receita_base"`
	OtherRecipes     string `json:"other_recipes,omitempty" bson:"other_recipes,omitempty" gorm:"column:detalhamento_outras_receitas"`
	Expenditure      string `json:"expenditure,omitempty" bson:"expenditure,omitempty" gorm:"column:detalhamento_descontos"`
}

type Score struct {
	Score             float64 `json:"score,omitempty" bson:"score,omitempty" gorm:"column:indice_transparencia"`
	CompletenessScore float64 `json:"completeness_score,omitempty" bson:"completeness_score,omitempty" gorm:"column:indice_completude"`
	EasinessScore     float64 `json:"easiness_score,omitempty" bson:"easiness_score,omitempty" gorm:"column:indice_facilidade"`
}

// MonthlyInfoVersion é um item do histórico de coletas armazenado no banco de dados.
type MonthlyInfoVersion struct {
	AgencyID  string            `json:"aid,omitempty" bson:"aid,omitempty"`
	Month     int               `json:"month,omitempty" bson:"month,omitempty"`
	Year      int               `json:"year,omitempty" bson:"year,omitempty"`
	VersionID int64             `json:"version_id,omitempty" bson:"version_id,omitempty"` // revisão/versão do irem. O tipo é int64 pois podemos querer usar epoch para ficar mais simples.
	Version   AgencyMonthlyInfo `json:"version,omitempty" bson:"version,omitempty"`
}
