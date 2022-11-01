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

	"github.com/dadosjusbr/proto/coleta"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Estrutura que contém as principais descrições de cada órgão.
type Orgao struct { 
	ID         string      `json:"id" bson:"id,omitempty" gorm:"column:id"`                 // 'trt13'
	Nome       string      `json:"nome" bson:"nome,omitempty" gorm:"column:nome"`             // 'Tribunal Regional do Trabalho 13° Região'
	Jurisdicao string      `json:"jurisdicao" bson:"jurisdicao,omitempty" gorm:"column:jurisdicao"` // Estadual, Trabalho, Federal
	Entidade   string      `json:"entidade" bson:"entidade,omitempty" gorm:"column:entidade"`     // Ministério, Tribunal
	UF         string      `json:"uf" bson:"uf,omitempty" gorm:"column:uf"`                 // Unidade Federativa
	FlagURL    string      `json:"url" bson:"url,omitempty"`               // Link for state url
	Coletando  []Coletando `json:"coletando" bson:"coletando,omitempty"`
}

// Estrutura que contém o dia em que verificamos o status dos dados e as razões pelas quais não os coletamos.
type Coletando struct {
	Timestamp *timestamppb.Timestamp `json:"timestamp" bson:"timestamp,omitempty"` // Day(unix) we checked the status of the data
	Descricao []string               `json:"descricao" bson:"descricao,omitempty"` // Reasons why we didn't collect the data
}

// Estrutura que contém as informações de uma coleta.
type Coleta struct {
	ID                 string           `json:"id" bson:"-" gorm:"primaryKey;column:id"` // 'trt13/01/2020'
	Timestamp          time.Time        `json:"timestamp,omitempty" bson:"-" gorm:"primaryKey;column:timestamp"`
	IdOrgao            string           `json:"id_orgao,omitempty" bson:"aid,omitempty" gorm:"column:id_orgao"`
	Mes                int              `json:"mes,omitempty" bson:"month,omitempty" gorm:"column:mes"`
	Ano                int              `json:"ano,omitempty" bson:"year,omitempty" gorm:"column:ano"`
	Atual              bool             `json:"atual,omitempty" bson:"-" gorm:"column:atual"`
	Backup             []Backup         `json:"backups,omitempty" bson:"backups,omitempty" gorm:"-"`
	Sumario            Sumario         	`json:"sumario,omitempty" bson:"summary,omitempty" gorm:"-"` 
	RepositorioColetor string           `json:"repositorio_coletor,omitempty" bson:"crawler_repo,omitempty" gorm:"column:repositorio_coletor"`
	VersaoColetor      string           `json:"versao_coletor,omitempty" bson:"crawler_version,omitempty" gorm:"column:versao_coletor"`
	RepositorioParser  string           `json:"repositorio_parser,omitempty" bson:"parser_repository,omitempty" gorm:"column:repositorio_parser"`
	VersaoParser       string           `json:"versao_parser,omitempty" bson:"parser_version,omitempty" gorm:"column:versao_parser"`
	ProcInfo           *coleta.ProcInfo `json:"procinfo,omitempty" bson:"procinfo,omitempty" gorm:"-"`
	Package            *Backup          `json:"package,omitempty" bson:"package,omitempty" gorm:"-"`
	Meta       													`json:"meta,omitempty" bson:"meta,omitempty"`
	Indice          										`json:"indices,omitempty" bson:"score,omitempty"`
	CrawlerID         string            `json:"crawler_id,omitempty" bson:"crawler_id,omitempty" gorm:"-"`
	CrawlingTimestamp *timestamppb.Timestamp `json:"crawling_ts,omitempty" bson:"crawling_ts,omitempty" gorm:"-"`   // Crawling moment (always UTC)
}

// Backup contains the URL to download a file and a hash to track if in the future will be changes in the file.
type Backup struct {
	URL  string `json:"url" bson:"url,omitempty"`
	Hash string `json:"hash" bson:"hash,omitempty"`
	Size int64  `json:"size" bson:"size,omitempty"`
}

// Summaries contains all summary detailed information
type Summaries struct {
	General       Sumario `json:"general,omitempty" bson:"general"`
	MemberActive  Sumario `json:"memberactive,omitempty" bson:"memberactive"`
	Undefined     Sumario `json:"undefined,omitempty" bson:"undefined"`
	ServantActive Sumario `json:"servantactive,omitempty" bson:"servantactive"`
}

// O sumário contém estatisitcas sobre a coleta de um mês de um órgão.
type Sumario struct {
	Membros            int         `json:"membros,omitempty" bson:"count,omitempty"`
	RemuneracaoBase    DataSummary `json:"remuneracao_base,omitempty" bson:"base_remuneration,omitempty"`
	OutrasRemuneracoes DataSummary `json:"outras_remuneracoes,omitempty" bson:"other_remunerations,omitempty"`
	HistogramaRenda    map[int]int `json:"histograma_renda,omitempty" bson:"hist,omitempty"`
}

// DataSummary A Struct containing data summary with statistics.
type DataSummary struct {
	Max     float64 `json:"max" bson:"max,omitempty"` 
	Min     float64 `json:"min" bson:"min,omitempty"`
	Average float64 `json:"avg" bson:"avg,omitempty"`
	Total   float64 `json:"total" bson:"total,omitempty"`
}

// Estrutura que contém os metadados de uma coleta.
type Meta struct {
	NaoRequerLogin             bool   `json:"no_login_required,omitempty" bson:"no_login_required,omitempty" gorm:"column:nao_requer_login"`
	NaoRequerCaptcha           bool   `json:"no_captcha_required,omitempty" bson:"no_captcha_required,omitempty" gorm:"column:nao_requer_captcha"`
	Acesso                     string `json:"access,omitempty" bson:"access,omitempty" gorm:"column:acesso"`
	Extensao                   string `json:"extension,omitempty" bson:"extension,omitempty" gorm:"column:extensao"`
	EstritamenteTabular        bool   `json:"strictly_tabular,omitempty" bson:"strictly_tabular,omitempty" gorm:"column:estritamente_tabular"`
	FormatoConsistente         bool   `json:"consistent_format,omitempty" bson:"consistent_format,omitempty" gorm:"column:formato_consistente"`
	TemMatricula               bool   `json:"have_enrollment,omitempty" bson:"have_enrollment,omitempty" gorm:"column:tem_matricula"`
	TemLotacao                 bool   `json:"there_is_a_capacity,omitempty" bson:"there_is_a_capacity,omitempty" gorm:"column:tem_lotacao"`
	TemCargo                   bool   `json:"has_position,omitempty" bson:"has_position,omitempty" gorm:"column:tem_cargo"`
	DetalhamentoReceitaBase    string `json:"detalhamento_receita_base,omitempty" bson:"base_revenue,omitempty" gorm:"column:detalhamento_receita_base"`
	DetalhamentoOutrasReceitas string `json:"detalhamento_outras_receitas,omitempty" bson:"other_recipes,omitempty" gorm:"column:detalhamento_outras_receitas"`
	DetalhamentoDescontos      string `json:"detalhamento_descontos,omitempty" bson:"expenditure,omitempty" gorm:"column:detalhamento_descontos"`
}

// Indices de uma coleta.
type Indice struct {
	IndiceTransparencia float64 `json:"indice_transparencia,omitempty" bson:"score,omitempty" gorm:"column:indice_transparencia"`
	IndiceCompletude    float64 `json:"indice_completude,omitempty" bson:"completeness_score,omitempty" gorm:"column:indice_completude"`
	IndiceFacilidade    float64 `json:"indice_facilidade,omitempty" bson:"easiness_score,omitempty" gorm:"column:indice_facilidade"`
}

// MonthlyInfoVersion é um item do histórico de coletas armazenado no banco de dados.
type MonthlyInfoVersion struct {
	AgencyID  string            `json:"aid,omitempty" bson:"aid,omitempty"`
	Month     int               `json:"month,omitempty" bson:"month,omitempty"`
	Year      int               `json:"year,omitempty" bson:"year,omitempty"`
	VersionID int64             `json:"version_id,omitempty" bson:"version_id,omitempty"` // revisão/versão do irem. O tipo é int64 pois podemos querer usar epoch para ficar mais simples.
	Version   Coleta `json:"version,omitempty" bson:"version,omitempty"`
}
