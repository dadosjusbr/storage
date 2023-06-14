package dto

import (
	"time"

	"github.com/dadosjusbr/storage/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnnualSummaryDTO struct {
	Year               int       `gorm:"column:ano"`
	AverageCount       int       `gorm:"column:media_num_membros"`
	TotalCount         int       `gorm:"column:total_num_membros"`
	BaseRemuneration   float64   `gorm:"column:remuneracao_base"`
	OtherRemunerations float64   `gorm:"column:outras_remuneracoes"`
	Discounts          float64   `gorm:"column:descontos"`
	NumMonthsWithData  int       `gorm:"column:meses_com_dados"`
	Timestamp          time.Time `gorm:"column:timestamp"`
}

func NewAnnualSummaryDTO(ami models.AnnualSummary) *AnnualSummaryDTO {
	var timestamp time.Time
	if ami.Timestamp != nil {
		timestamp = time.Unix(ami.Timestamp.Seconds, int64(ami.Timestamp.Nanos))
	} else {
		timestamp = time.Now()
	}
	return &AnnualSummaryDTO{
		Year:               ami.Year,
		AverageCount:       ami.AverageCount,
		TotalCount:         ami.TotalCount,
		BaseRemuneration:   ami.BaseRemuneration,
		OtherRemunerations: ami.OtherRemunerations,
		Discounts:          ami.Discounts,
		NumMonthsWithData:  ami.NumMonthsWithData,
		Timestamp:          timestamp,
	}
}

func (ami *AnnualSummaryDTO) ConvertToModel() *models.AnnualSummary {
	return &models.AnnualSummary{
		Year:               ami.Year,
		AverageCount:       ami.AverageCount,
		TotalCount:         ami.TotalCount,
		BaseRemuneration:   ami.BaseRemuneration,
		OtherRemunerations: ami.OtherRemunerations,
		Discounts:          ami.Discounts,
		NumMonthsWithData:  ami.NumMonthsWithData,
		Timestamp:          timestamppb.New(ami.Timestamp),
	}
}
