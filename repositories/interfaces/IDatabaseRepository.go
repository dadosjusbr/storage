package interfaces

import (
	"github.com/dadosjusbr/storage/models"
)

type IDatabaseRepository interface {
	Connect() error
	Disconnect() error
	Store(agmi models.AgencyMonthlyInfo) error
	StorePackage(newPackage models.Package) error
	// OPE : Órgãos Por Estado
	GetOPE(uf string, year int) ([]models.Agency, error)
	GetAgenciesCount() (int64, error)
	GetNumberOfMonthsCollected() (int64, error)
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
}