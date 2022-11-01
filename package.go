package storage

type Package struct {
	IdOrgao *string `json:"aid,omitempty" bson:"aid,omitempty"`
	Mes     *int    `json:"mes,omitempty" bson:"month,omitempty"`
	Ano     *int    `json:"ano,omitempty" bson:"year,omitempty"`
	Grupo   *string `json:"grupo,omitempty" bson:"group,omitempty"`
	Package Backup
}
type PackageFilterOpts struct {
	IdOrgao *string `json:"id_orgao,omitempty" bson:"aid,omitempty"`
	Mes     *int    `json:"mes,omitempty" bson:"month,omitempty"`
	Ano     *int    `json:"ano,omitempty" bson:"year,omitempty"`
	Grupo   *string `json:"grupo,omitempty" bson:"group,omitempty"`
}
