package main

import (
	"log"

	"github.com/dadosjusbr/datapackage"
	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/repo/database"
	"github.com/dadosjusbr/storage/repo/file_storage"
	dpkg "github.com/frictionlessdata/datapackage-go/datapackage"
	"github.com/joho/sqltocsv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	PostgresUser     string `envconfig:"POSTGRES_USER" required:"true"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	PostgresDBName   string `envconfig:"POSTGRES_DB" required:"true"`
	PostgresHost     string `envconfig:"POSTGRES_HOST" required:"true"`
	PostgresPort     string `envconfig:"POSTGRES_PORT" required:"true"`

	AWSRegion string `envconfig:"AWS_REGION" required:"true"`
	S3Bucket  string `envconfig:"S3_BUCKET" required:"true"`
}

func main() {
	var conf config
	if err := envconfig.Process("", &conf); err != nil {
		log.Fatal(err)
	}
	//Criando o client do Postgres
	postgresDb, err := database.NewPostgresDB(conf.PostgresUser, conf.PostgresPassword, conf.PostgresDBName, conf.PostgresHost, conf.PostgresPort)
	if err != nil {
		log.Fatalf("error creating Postgres client: %v", err.Error())
	}
	// Criando o client do S3
	s3Client, err := file_storage.NewS3Client(conf.AWSRegion, conf.S3Bucket)
	if err != nil {
		log.Fatalf("error creating S3 client: %v", err.Error())
	}
	// Criando o client do storage a partir do banco postgres e do client do s3
	pgS3Client, err := storage.NewClient(postgresDb, s3Client)
	if err != nil {
		log.Fatalf("error setting up postgres storage client: %v", err)
	}
	defer pgS3Client.Db.Disconnect()

	// Consultando os dados de todas as tabelas
	dados, err := postgresDb.Dump()
	if err != nil {
		log.Fatalf("error Dump(): %v", err)
	}

	// Criando os CSVs
	for csvName, rows := range dados {
		csv := sqltocsv.New(rows)
		csv.WriteFile(csvName)
	}

	// Criando o pacote
	pkgName := "dump-dadosjusbr.zip"
	desc, err := datapackage.DescriptorMapV2()
	if err != nil {
		log.Fatalf("error DescriptorMapV2(): %v", err)
	}
	pkg, err := dpkg.New(desc, ".")
	if err != nil {
		log.Fatalf("error create datapackage: %v", err)
	}
	if err := pkg.Zip(pkgName); err != nil {
		log.Fatalf("error zipping datapackage: %v", err)
	}

	// Armazenando no S3
	_, err = pgS3Client.Cloud.UploadFile(pkgName, pkgName)
	if err != nil {
		log.Fatalf("error while uploading dump (%s): %v", pkgName, err)
	}
}
