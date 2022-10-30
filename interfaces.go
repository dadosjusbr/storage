package storage

type IDatabase interface {
  Connect() error
  Disconnect() error
  Store(agmi Coleta) error
  StorePackage(newPackage Package) error
  GetOPE(uf string, year int) ([]Orgao, map[string][]Coleta, error)
  GetAgenciesCount() (int64, error)
  GetNumberOfMonthsCollected() (int64, error)
  GetAgencies(uf string) ([]Orgao, error)
  GetAgency(aid string)(*Orgao, error)
  GetAllAgencies() ([]Orgao, error)
  GetMonthlyInfo(agencies []Orgao, year int) (map[string][]Coleta, error)
  GetMonthlyInfoSummary(agencies []Orgao, year int) (map[string][]Coleta, error)
  GetOMA(month int, year int, agency string) (*Coleta, *Orgao, error)
  GetGeneralMonthlyInfosFromYear(year int) ([]GeneralMonthlyInfo, error)
  GetFirstDateWithMonthlyInfo() (int, int, error)
  GetLastDateWithMonthlyInfo() (int, int, error)
  GetRemunerationSummary() (*RemmunerationSummary, error)
  GetPackage(pkgOpts PackageFilterOpts) (*Package, error)
}

type IStorageService interface {
  UploadFile(srcPath string, dstFolder string) (*Backup, error)
  Backup(Files []string, dstFolder string) ([]Backup, error)
}