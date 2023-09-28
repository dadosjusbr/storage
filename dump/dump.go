package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dadosjusbr/datapackage"
	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/repo/database"
	"github.com/dadosjusbr/storage/repo/file_storage"
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
	coletas, contracheques, remuneracoes, err := postgresDb.Dump()
	if err != nil {
		log.Fatalf("error Dump(): %v", err)
	}

	var coletaCSV []datapackage.Coleta_CSV_V2
	var metadadosCSV []datapackage.Metadados_CSV_V2

	for _, c := range coletas {
		chave := fmt.Sprintf("%s/%s/%d", c.AgencyID, addZeroes(c.Month), c.Year)
		coletaCSV = append(coletaCSV, datapackage.Coleta_CSV_V2{
			ChaveColeta:        chave,
			Orgao:              c.AgencyID,
			Mes:                int32(c.Month),
			Ano:                int32(c.Year),
			TimestampColeta:    c.CrawlingTimestamp.AsTime(),
			RepositorioColetor: c.CrawlerRepo,
			VersaoColetor:      c.CrawlerVersion,
			RepositorioParser:  c.ParserRepo,
			VersaoParser:       c.ParserVersion,
		})

		metadadosCSV = append(metadadosCSV, datapackage.Metadados_CSV_V2{
			Orgao:                      c.AgencyID,
			Mes:                        int32(c.Month),
			Ano:                        int32(c.Year),
			FormatoAberto:              c.Meta.OpenFormat,
			Acesso:                     c.Meta.Access,
			Extensao:                   c.Meta.Extension,
			EstritamenteTabular:        c.Meta.StrictlyTabular,
			FormatoConsistente:         c.Meta.ConsistentFormat,
			TemMatricula:               c.Meta.HaveEnrollment,
			TemLotacao:                 c.Meta.ThereIsACapacity,
			TemCargo:                   c.Meta.HasPosition,
			DetalhamentoReceitaBase:    c.Meta.BaseRevenue,
			DetalhamentoOutrasReceitas: c.Meta.OtherRecipes,
			DetalhamentoDescontos:      c.Meta.Expenditure,
			IndiceCompletude:           float32(c.Score.CompletenessScore),
			IndiceFacilidade:           float32(c.Score.EasinessScore),
			IndiceTransparencia:        float32(c.Score.Score),
		})
	}

	var contrachequeCSV []datapackage.Contracheque_CSV_V2
	for _, cc := range contracheques {
		contrachequeCSV = append(contrachequeCSV, datapackage.Contracheque_CSV_V2{
			IdContracheque: cc.ID,
			Orgao:          cc.Agency,
			Mes:            int32(cc.Month),
			Ano:            int32(cc.Year),
			Nome:           cc.Name,
			Matricula:      cc.RegisterID,
			Funcao:         cc.Role,
			LocalTrabalho:  cc.Workplace,
			Salario:        cc.Salary,
			Beneficios:     cc.Benefits,
			Descontos:      cc.Discounts,
			Remuneracao:    cc.Remuneration,
			Situacao:       dereference(cc.Situation),
		})
	}

	var remuneracaoCSV []datapackage.Remuneracao_CSV_V2
	for _, r := range remuneracoes {
		remuneracaoCSV = append(remuneracaoCSV, datapackage.Remuneracao_CSV_V2{
			IdContracheque: r.PaycheckID,
			Orgao:          r.Agency,
			Mes:            int32(r.Month),
			Ano:            int32(r.Year),
			Tipo:           r.Type,
			Categoria:      r.Category,
			Item:           r.Item,
			Valor:          r.Value,
		})
	}

	rc := datapackage.ResultadoColeta_CSV_V2{
		Coleta:       coletaCSV,
		Remuneracoes: remuneracaoCSV,
		Folha:        contrachequeCSV,
		Metadados:    metadadosCSV,
	}

	// Criando o pacote
	year, month, _ := time.Now().Date()
	pkgName := fmt.Sprintf("dadosjusbr-%d-%d.zip", year, month)
	if err := datapackage.ZipV2(pkgName, rc, true); err != nil {
		log.Fatalf("error ZipV2(): %w", err)
	}

	// Armazenando no S3
	_, err = pgS3Client.Cloud.UploadFile(pkgName, "dumps/"+pkgName)
	if err != nil {
		log.Fatalf("error while uploading dump (%s): %v", pkgName, err)
	}
}

func addZeroes(num int) string {
	numStr := strconv.Itoa(num)
	if len(numStr) == 1 {
		numStr = "0" + numStr
	}
	return numStr
}

func dereference(p *string) string {
	if p != nil {
		return *p
	} else {
		return ""
	}
}
