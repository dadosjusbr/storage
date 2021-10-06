package main

import (
	"fmt"

	"github.com/dadosjusbr/storage"
)

type config struct {
	Port   int    `envconfig:"PORT"`
	DBUrl  string `envconfig:"MONGODB_URI"`
	DBName string `envconfig:"MONGODB_NAME"`

	// StorageDB config
	MongoURI    string `envconfig:"MONGODB_URI"`
	MongoDBName string `envconfig:"MONGODB_NAME"`
	MongoMICol  string `envconfig:"MONGODB_MICOL" required:"true"`
	MongoAgCol  string `envconfig:"MONGODB_AGCOL" required:"true"`
	MongoPkgCol string `envconfig:"MONGODB_AGRECOL" required:"true"`

	// Omited fields
	EnvOmittedFields []string `envconfig:"ENV_OMITTED_FIELDS"`
}

func newClient(c config) (*storage.Client, error) {
	if c.MongoMICol == "" || c.MongoAgCol == "" {
		return nil, fmt.Errorf("error creating storage client: db collections must not be empty. MI:\"%s\", AG:\"%s\"", c.MongoMICol, c.MongoAgCol)
	}
	db, err := storage.NewDBClient(c.MongoURI, c.MongoDBName, c.MongoMICol, c.MongoAgCol, c.MongoPkgCol)
	if err != nil {
		return nil, fmt.Errorf("error creating DB client: %q", err)
	}
	db.Collection(c.MongoMICol)
	client, err := storage.NewClient(db, &storage.CloudClient{})
	if err != nil {
		return nil, fmt.Errorf("error creating storage.client: %q", err)
	}
	return client, nil
}
func main() {
	conf := config{
		MongoURI:   "mongodb+srv://dadosjusbr:dadosjus123@cluster-798622w0.dryls.mongodb.net/heroku_798622w0",
		DBName:     "heroku_798622w0",
		MongoMICol: "miProto", MongoAgCol: "ag",
		Port:        8081,
		DBUrl:       "mongodb+srv://dadosjusbr:dadosjus123@cluster-798622w0.dryls.mongodb.net/heroku_798622w0",
		MongoDBName: "heroku_798622w0",
		MongoPkgCol: "pkg",
	}
	client, err := newClient(conf)
	if err != nil {
		fmt.Println("deu ruim")
	}
	// fmt.Println(client.GetLastDateWithMonthlyInfo())
	// fmt.Println(client.GetFirstDateWithMonthlyInfo())
	// fmt.Println(client.Db.GetRemunerationSummary())
	// mg, _ := client.Db.GetGeneralMonthlyInfosFromYear(2021)
	// fmt.Println(mg[0].Wage)
	// aID := "mppb"
	// year := 2020
	// pkg, err := client.Db.GetPackage(storage.PackageFilterOpts{AgencyID: &aID, Month: nil, Year: &year, Group: nil})
	// if err != nil {
	// 	fmt.Println("%q", err)
	// }
	// fmt.Println(pkg)
	agency, err := client.Db.GetAgency("mppb")
	agencies := []storage.Agency{*agency}
	monthlyInfos, err := client.Db.GetMonthlyInfo(agencies, 2019)
	for _, m := range monthlyInfos {
		fmt.Println(m)
	}
	agencies, err := client.Db.GetAgencies("AL")

}
