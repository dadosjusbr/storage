package database

import (
	"github.com/dadosjusbr/storage/models"
)

type Interface interface {
	Connect() error
	Disconnect() error
	Store(agmi models.AgencyMonthlyInfo) error
	StoreRemunerations(remu models.Remunerations) error
	GetStateAgencies(uf string) ([]models.Agency, error)
	// OPJ: Órgãos por jurisdição.
	GetOPJ(group string) ([]models.Agency, error)
	GetNumberOfMonthsCollected() (int, error)
	GetAgenciesCount() (int, error)
	GetAgenciesByUF(uf string) ([]models.Agency, error)
	GetAgency(aid string) (*models.Agency, error)
	GetAllAgencies() ([]models.Agency, error)
	GetMonthlyInfo(agencies []models.Agency, year int) (map[string][]models.AgencyMonthlyInfo, error)
	GetAnnualSummary(agency string) ([]models.AnnualSummary, error)
	// OMA: Órgão Mês Ano
	GetOMA(month int, year int, agency string) (*models.AgencyMonthlyInfo, *models.Agency, error)
	GetGeneralMonthlyInfosFromYear(year int) ([]models.GeneralMonthlyInfo, error)
	GetFirstDateWithMonthlyInfo() (int, int, error)
	GetLastDateWithMonthlyInfo() (int, int, error)
	GetGeneralMonthlyInfo() (float64, error)
	GetAllIndices(groupName string) (map[string][]models.Detail, error)
}
