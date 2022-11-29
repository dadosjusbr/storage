package dto

import "github.com/dadosjusbr/storage/models"

type RemunerationsDTO struct {
  AgencyID          string `gorm:"column:id_orgao"`
  Month             int    `gorm:"column:mes"`
  Year              int    `gorm:"column:ano"`
  NumDiscounts     int    `gorm:"column:linhas_descontos"`
  NumBase         int    `gorm:"column:linhas_base"`
  NumOther         int    `gorm:"column:linhas_outras"`
  ZipUrl          string `gorm:"column:zip_url"`
}

func (RemunerationsDTO) TableName() string {
	return "remuneracoes_zips"
}

func (a RemunerationsDTO) ConvertToModel() *models.Remunerations {
  return &models.Remunerations{
    AgencyID:          a.AgencyID,
    Month:             a.Month,
    Year:              a.Year,
    NumDiscounts:      a.NumDiscounts,
    NumBase:           a.NumBase,
    NumOther:          a.NumOther,
    ZipUrl:           a.ZipUrl,
  }
}

func NewRemunerationsDTO(r models.Remunerations) *RemunerationsDTO {
  return &RemunerationsDTO{
    AgencyID:          r.AgencyID,
    Month:             r.Month,
    Year:              r.Year,
    NumDiscounts:      r.NumDiscounts,
    NumBase:           r.NumBase,
    NumOther:          r.NumOther,
    ZipUrl:           r.ZipUrl,
  }
}

