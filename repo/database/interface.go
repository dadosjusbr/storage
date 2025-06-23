package database

import (
	"github.com/dadosjusbr/storage/models"
)

type Interface interface {
	Connect() error
	Disconnect() error
	Store(agmi models.AgencyMonthlyInfo) error
	// StorePaychecks: armazena dados nas tabelas 'contracheques' e 'remuneracoes'
	StorePaychecks(p []models.Paycheck, r []models.PaycheckItem) error
	// StoreRemunerations: armazena dados dos zips de remunerações que estão no S3.
	StoreRemunerations(remu models.Remunerations) error
	GetStateAgencies(uf string) ([]models.Agency, error)
	// OPJ: Órgãos por jurisdição.
	GetOPJ(group string) ([]models.Agency, error)
	GetNumberOfMonthsCollected() (int, error)
	GetNumberOfPaychecksCollected() (int, error)
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
	GetIndexInformation(name string, month, year int) (map[string][]models.IndexInformation, error)
	GetAllAgencyCollection(agency string) ([]models.AgencyMonthlyInfo, error)
	GetPaychecks(agency models.Agency, year int) ([]models.Paycheck, error)
	GetPaycheckItems(agency models.Agency, year int) ([]models.PaycheckItem, error)
	GetAveragePerCapita(agency string, year int) (*models.PerCapitaData, error)
	GetNotices(agency string, year int, month int) ([]*string, error)
	GetAveragePerAgency(year int) ([]models.PerCapitaData, error)
}
