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
	"time"

	"github.com/dadosjusbr/coletores"
)

// Agency A Struct containing the main descriptions of each Agency.
type Agency struct {
	ID      string `json:"aid" bson:"aid,omitempty"`       // 'trt13'
	Name    string `json:"name" bson:"name,omitempty"`     // 'Tribunal Regional do Trabalho 13° Região'
	Type    string `json:"type" bson:"type,omitempty"`     // "R" for Regional, "M" for Municipal, "F" for Federal, "E" for State.
	Entity  string `json:"entity" bson:"entity,omitempty"` // "J" For Judiciário, "M" for Ministério Público, "P" for Procuradorias and "D" for Defensorias.
	UF      string `json:"uf" bson:"uf,omitempty"`         // Short code for federative unity.
	FlagURL string `json:"url" bson:"url,omitempty"`       //Link for state url
}

// AgencyMonthlyInfo A Struct containing a snapshot of a agency in a month.
type AgencyMonthlyInfo struct {
	AgencyID          string               `json:"aid,omitempty" bson:"aid,omitempty"`
	Month             int                  `json:"month,omitempty" bson:"month,omitempty"`
	Year              int                  `json:"year,omitempty" bson:"year,omitempty"`
	Backups           []Backup             `json:"backups,omitempty" bson:"backups,omitempty"`
	Summary           Summaries            `json:"summary,omitempty" bson:"summary,omitempty"`
	Employee          []coletores.Employee `json:"employee,omitempty" bson:"employee,omitempty"`
	Crawler           coletores.Crawler    `json:"crawler,omitempty" bson:"crawler,omitempty"`
	CrawlingTimestamp time.Time            `json:"ts,omitempty" bson:"ts,omitempty"`             // Crawling moment (always UTC)
	ProcInfo          *coletores.ProcInfo  `json:"procinfo,omitempty" bson:"procinfo,omitempty"` // Making this a pointer because it should be an optional field due to backwards compatibility.
	Package           *Backup              `json:"package,omitempty" bson:"package,omitempty"`   // Making this a pointer because it should be an optional field due to backwards compatibility.
}

// Backup contains the URL to download a file and a hash to track if in the future will be changes in the file.
type Backup struct {
	URL  string `json:"url" bson:"url,omitempty"`
	Hash string `json:"hash" bson:"hash,omitempty"`
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
	Count           int         `json:"count" bson:"count,omitempty"`   // Number of employees
	Wage            DataSummary `json:"wage" bson:"wage,omitempty"`     //  Statistics (Max, Min, Median, Total)
	Perks           DataSummary `json:"perks" bson:"perks,omitempty"`   //  Statistics (Max, Min, Median, Total)
	Others          DataSummary `json:"others" bson:"others,omitempty"` //  Statistics (Max, Min, Median, Total)
	IncomeHistogram map[int]int `json:"hist" bson:"hist,omitempty"`
	Benefits        DataSummary `json:"benefits" bson:"benefits,omitempty"` // Agreggated summary of Perks+Others.
}

// DataSummary A Struct containing data summary with statistics.
type DataSummary struct {
	Max     float64 `json:"max" bson:"max,omitempty"`
	Min     float64 `json:"min" bson:"min,omitempty"`
	Average float64 `json:"avg" bson:"avg,omitempty"`
	Total   float64 `json:"total" bson:"total,omitempty"`
}
