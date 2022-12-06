package dto

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dadosjusbr/proto/coleta"
	"github.com/dadosjusbr/storage/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"
)

// AgencyMonthlyInfo A Struct containing a snapshot of a agency in a month.
type AgencyMonthlyInfoDTO struct {
	ID             string         `gorm:"column:id"` // 'trt13/01/2020'
	AgencyID       string         `gorm:"column:id_orgao"`
	Month          int            `gorm:"column:mes"`
	Year           int            `gorm:"column:ano"`
	Actual         bool           `gorm:"column:atual"`
	Backup         datatypes.JSON `gorm:"column:backups"`
	Summary        datatypes.JSON `gorm:"column:sumario"`
	CrawlerVersion string         `gorm:"column:versao_coletor"`
	CrawlerRepo    string         `gorm:"column:repositorio_coletor"` // The github Repository of MI Crawler
	ParserRepo     string         `gorm:"column:repositorio_parser"`  // The github Repository of MI Parser
	ParserVersion  string         `gorm:"column:versao_parser"`
	Timestamp      time.Time      `gorm:"column:timestamp"`
	ProcInfo       datatypes.JSON `gorm:"column:procinfo"`
	Package        datatypes.JSON `gorm:"column:package"`
	Meta
	Score
	//TODO: Add ExectionTime
}

func (AgencyMonthlyInfoDTO) TableName() string {
	return "coletas"
}

type Meta struct {
	OpenFormat       bool   `gorm:"column:formato_aberto"`
	Access           string `gorm:"column:acesso"`
	Extension        string `gorm:"column:extensao"`
	StrictlyTabular  bool   `gorm:"column:estritamente_tabular"`
	ConsistentFormat bool   `gorm:"column:formato_consistente"`
	HaveEnrollment   bool   `gorm:"column:tem_matricula"`
	ThereIsACapacity bool   `gorm:"column:tem_lotacao"`
	HasPosition      bool   `gorm:"column:tem_cargo"`
	BaseRevenue      string `gorm:"column:detalhamento_receita_base"`
	OtherRecipes     string `gorm:"column:detalhamento_outras_receitas"`
	Expenditure      string `gorm:"column:detalhamento_descontos"`
}

type Score struct {
	Score             float64 `gorm:"column:indice_transparencia"`
	CompletenessScore float64 `gorm:"column:indice_completude"`
	EasinessScore     float64 `gorm:"column:indice_facilidade"`
}

func (a AgencyMonthlyInfoDTO) ConvertToModel() (*models.AgencyMonthlyInfo, error) {
	var backup models.Backup
	var summary models.Summary
	var procInfo coleta.ProcInfo
	var pkg models.Backup

	backupBytes, err := a.Backup.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("error while marshaling backup: %q", err)
	}
	err = json.Unmarshal(backupBytes, &backup)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling backup: %q", err)
	}

	summaryBytes, err := a.Summary.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("error while marshaling summary: %q", err)
	}
	err = json.Unmarshal(summaryBytes, &summary)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling summary: %q", err)
	}

	procInfoBytes, err := a.ProcInfo.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("error while marshaling procInfo: %q", err)
	}
	err = json.Unmarshal(procInfoBytes, &procInfo)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling procInfo: %q", err)
	}

	pkgBytes, err := a.Package.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("error while marshaling package: %q", err)
	}
	err = json.Unmarshal(pkgBytes, &pkg)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling package: %q", err)
	}

	return &models.AgencyMonthlyInfo{
		AgencyID:          a.AgencyID,
		Month:             a.Month,
		Year:              a.Year,
		CrawlerVersion:    a.CrawlerVersion,
		CrawlerRepo:       a.CrawlerRepo,
		ParserRepo:        a.ParserRepo,
		ParserVersion:     a.ParserVersion,
		CrawlingTimestamp: timestamppb.New(a.Timestamp),
		Score: &models.Score{
			Score:             a.Score.Score,
			CompletenessScore: a.Score.CompletenessScore,
			EasinessScore:     a.Score.EasinessScore,
		},
		Meta: &models.Meta{
			OpenFormat:       a.Meta.OpenFormat,
			Expenditure:      a.Meta.Expenditure,
			Access:           a.Meta.Access,
			Extension:        a.Meta.Extension,
			StrictlyTabular:  a.Meta.StrictlyTabular,
			ConsistentFormat: a.Meta.ConsistentFormat,
			HaveEnrollment:   a.Meta.HaveEnrollment,
			ThereIsACapacity: a.Meta.ThereIsACapacity,
			HasPosition:      a.Meta.HasPosition,
			BaseRevenue:      a.Meta.BaseRevenue,
			OtherRecipes:     a.Meta.OtherRecipes,
		},
		Summary:  summary,
		Backups:  []models.Backup{backup},
		ProcInfo: &procInfo,
		Package:  &pkg,
	}, nil
}

func NewAgencyMonthlyInfoDTO(agmi models.AgencyMonthlyInfo) (*AgencyMonthlyInfoDTO, error) {
	backup, err := json.Marshal(agmi.Backups[0])
	if err != nil {
		return nil, fmt.Errorf("error while marshaling backup: %q", err)
	}
	summary, err := json.Marshal(agmi.Summary)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling summary: %q", err)
	}
	procInfo, err := json.Marshal(agmi.ProcInfo)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling procInfo: %q", err)
	}
	pkg, err := json.Marshal(agmi.Package)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling package: %q", err)
	}

	return &AgencyMonthlyInfoDTO{
		ID:             fmt.Sprintf("%s/%s/%d", agmi.AgencyID, AddZeroes(agmi.Month), agmi.Year),
		Actual:         true,
		AgencyID:       agmi.AgencyID,
		Month:          agmi.Month,
		Year:           agmi.Year,
		CrawlerVersion: agmi.CrawlerVersion,
		CrawlerRepo:    agmi.CrawlerRepo,
		ParserRepo:     agmi.ParserRepo,
		ParserVersion:  agmi.ParserVersion,
		Timestamp:      time.Unix(agmi.CrawlingTimestamp.Seconds, int64(agmi.CrawlingTimestamp.Nanos)),
		Score: Score{
			Score:             agmi.Score.Score,
			CompletenessScore: agmi.Score.CompletenessScore,
			EasinessScore:     agmi.Score.EasinessScore,
		},
		Meta: Meta{
			OpenFormat:       agmi.Meta.OpenFormat,
			Expenditure:      agmi.Meta.Expenditure,
			Access:           agmi.Meta.Access,
			Extension:        agmi.Meta.Extension,
			StrictlyTabular:  agmi.Meta.StrictlyTabular,
			ConsistentFormat: agmi.Meta.ConsistentFormat,
			HaveEnrollment:   agmi.Meta.HaveEnrollment,
			ThereIsACapacity: agmi.Meta.ThereIsACapacity,
			HasPosition:      agmi.Meta.HasPosition,
			BaseRevenue:      agmi.Meta.BaseRevenue,
			OtherRecipes:     agmi.Meta.OtherRecipes,
		},
		Summary:  summary,
		Backup:   backup,
		ProcInfo: procInfo,
		Package:  pkg,
	}, nil
}

// Funcão que adiciona um zero a esquerda a um número caso ele seja menor que 10
func AddZeroes(num int) string {
	numStr := strconv.Itoa(num)
	if len(numStr) == 1 {
		numStr = "0" + numStr
	}
	return numStr
}
