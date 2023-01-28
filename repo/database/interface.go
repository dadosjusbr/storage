package database

import (
	"github.com/dadosjusbr/storage/models"
)

type Interface interface {
	Connect() error
	Disconnect() error
	Store(agmi models.AgencyMonthlyInfo) error
	StorePackage(newPackage models.Package) error
	StoreRemunerations(remu models.Remunerations) error
	// OPE : Órgãos Por Estado
	GetOPE(uf string) ([]models.Agency, error)
	// OPJ: Órgãos por jurisdição.
	GetOPJ(group string) ([]models.Agency, error)
	GetAgenciesCount() (int64, error)
	GetNumberOfMonthsCollected() (int, error)
	GetAgencies(uf string) ([]models.Agency, error)
	GetAgency(aid string) (*models.Agency, error)
	GetAllAgencies() ([]models.Agency, error)
	GetMonthlyInfo(agencies []models.Agency, year int) (map[string][]models.AgencyMonthlyInfo, error)
	GetMonthlyInfoSummary(agencies []models.Agency, year int) (map[string][]models.AgencyMonthlyInfo, error)
	// OMA: Órgão Mês Ano
	GetOMA(month int, year int, agency string) (*models.AgencyMonthlyInfo, *models.Agency, error)
	GetGeneralMonthlyInfosFromYear(year int) ([]models.GeneralMonthlyInfo, error)
	GetFirstDateWithMonthlyInfo() (int, int, error)
	GetLastDateWithMonthlyInfo() (int, int, error)
	GetRemunerationSummary() (*models.RemmunerationSummary, error)
	GetPackage(pkgOpts models.PackageFilterOpts) (*models.Package, error)
	GetGeneralMonthlyInfo() (float64, error)
}
