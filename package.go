package storage

type Package struct {
	AgencyID *string `json:"aid,omitempty" bson:"aid,omitempty"`
	Month    *int    `json:"month,omitempty" bson:"month,omitempty"`
	Year     *int    `json:"year,omitempty" bson:"year,omitempty"`
	Group    *string `json:"group,omitempty" bson:"group,omitempty"`
	Package  Backup
}
type PackageFilterOpts struct {
	AgencyID *string `json:"aid,omitempty" bson:"aid,omitempty"`
	Month    *int    `json:"month,omitempty" bson:"month,omitempty"`
	Year     *int    `json:"year,omitempty" bson:"year,omitempty"`
	Group    *string `json:"group,omitempty" bson:"group,omitempty"`
}
