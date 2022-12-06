package models

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

// Backup contains the URL to download a file and a hash to track if in the future will be changes in the file.
type Backup struct {
	URL  string `json:"url" bson:"url,omitempty"`
	Hash string `json:"hash" bson:"hash,omitempty"`
	Size int64  `json:"size" bson:"size,omitempty"`
}
