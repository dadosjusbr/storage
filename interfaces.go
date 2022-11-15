package storage

type IDatabaseService interface {
	Connect() error
	Disconnect() error
	Store(agmi AgencyMonthlyInfo) error
	StorePackage(newPackage Package) error
	GetOPE(uf string, year int) ([]Agency, map[string][]AgencyMonthlyInfo, error)
	GetAgenciesCount() (int64, error)
	GetNumberOfMonthsCollected() (int64, error)
	GetAgencies(uf string) ([]Agency, error)
	GetAgency(aid string) (*Agency, error)
	GetAllAgencies() ([]Agency, error)
	GetMonthlyInfo(agencies []Agency, year int) (map[string][]AgencyMonthlyInfo, error)
	GetMonthlyInfoSummary(agencies []Agency, year int) (map[string][]AgencyMonthlyInfo, error)
	GetOMA(month int, year int, agency string) (*AgencyMonthlyInfo, *Agency, error)
	GetGeneralMonthlyInfosFromYear(year int) ([]GeneralMonthlyInfo, error)
	GetFirstDateWithMonthlyInfo() (int, int, error)
	GetLastDateWithMonthlyInfo() (int, int, error)
	GetRemunerationSummary() (*RemmunerationSummary, error)
	GetPackage(pkgOpts PackageFilterOpts) (*Package, error)
}

type IStorageService interface {
	UploadFile(srcPath string, dstFolder string) (*Backup, error)
}
