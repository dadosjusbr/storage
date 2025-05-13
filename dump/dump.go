package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dadosjusbr/datapackage"
	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/repo/database"
	"github.com/dadosjusbr/storage/repo/file_storage"
	dpkg "github.com/frictionlessdata/datapackage-go/datapackage"
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

	// Criando o pacote
	year, _, _ := time.Now().Date()
	pkgName := fmt.Sprintf("dadosjusbr-%d-%d.zip", year, 4)
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

	// Testando o pacote
	if _, err = datapackage.LoadV2(pkgName); err != nil {
		log.Fatalf("error loading datapackage: %v", err)
	}

	// Armazenando no S3
	_, err = pgS3Client.Cloud.UploadFile(pkgName, "dumps/"+pkgName)
	if err != nil {
		log.Fatalf("error while uploading dump (%s): %v", pkgName, err)
	}
}
