package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dadosjusbr/storage/models"
	"github.com/dadosjusbr/storage/repo/database/dto"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresDB struct {
	db       *gorm.DB
	user     string
	password string
	dbName   string
	host     string
	port     string
	uri      string
}

func NewPostgresDB(user, password, dbName, host, port string) (*PostgresDB, error) {
	// Verificando se as credenciais de conexão não estão vazias
	if user == "" {
		return nil, fmt.Errorf("user cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}
	if dbName == "" {
		return nil, fmt.Errorf("dbName cannot be empty")
	}
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}
	if port == "" {
		return nil, fmt.Errorf("port cannot be empty")
	}

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, password)
	postgresDB := &PostgresDB{
		user:     user,
		password: password,
		dbName:   dbName,
		host:     host,
		port:     port,
		uri:      uri,
	}
	//Conectando ao postgres
	if err := postgresDB.Connect(); err != nil {
		return nil, fmt.Errorf("error connecting to postgres (creds:%s):%q", uri, err)
	}
	return postgresDB, nil
}

func (p *PostgresDB) Connect() error {
	if p.db != nil {
		return nil
	} else {
		conn, err := sql.Open("nrpostgres", p.uri)
		if err != nil {
			panic(err)
		}
		ctx, canc := context.WithTimeout(context.Background(), 30*time.Second)
		defer canc()
		if err := conn.PingContext(ctx); err != nil {
			return fmt.Errorf("error connecting to postgres (creds:%s):%q", p.uri, err)
		}
		db, err := gorm.Open(postgres.New(postgres.Config{
			Conn: conn,
		}))
		if err != nil {
			return fmt.Errorf("error initializing gorm: %q", err)
		}
		p.db = db
		return nil
	}
}

func (p *PostgresDB) Disconnect() error {
	db, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("error returning sql DB: %q", err)
	}
	err = db.Close()
	if err != nil {
		return fmt.Errorf("error closing DB connection: %q", err)
	}
	return nil
}

func (p *PostgresDB) GetConnection() (*gorm.DB, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database not connected!")
	}
	return p.db, nil
}

func (p *PostgresDB) SetConnection(conn *gorm.DB) {
	p.db = conn
}

func (p *PostgresDB) Store(agmi models.AgencyMonthlyInfo) error {
	/*Criando o DTO da coleta a partir de um modelo. É necessário a utilização de
	DTO's para melhor escalabilidade de bancos de dados. Caso não fosse utilizado,
	não seria possível utilizar outros frameworks/bancos além do GORM, pois ele
	afeta diretamente os tipos e campos de uma struct.*/
	coletas, err := dto.NewAgencyMonthlyInfoDTO(agmi)
	if err != nil {
		return fmt.Errorf("error converting agency monthly info to dto: %q", err)
	}

	/* Iniciando a transação. É necessário que seja uma transação porque queremos
	executar vários scripts que são dependentes um do outro. Ou seja, se um falhar
	todos falham. Isso nos dá uma maior segurança ao executar a inserção. */
	err = p.db.Transaction(func(tx *gorm.DB) error {
		// Definindo atual como false para todos os registros com o mesmo ID.
		ID := fmt.Sprintf("%s/%s/%d", agmi.AgencyID, dto.AddZeroes(agmi.Month), agmi.Year)
		if err := tx.Model(dto.AgencyMonthlyInfoDTO{}).Where("id = ?", ID).Update("atual", false).Error; err != nil {
			return fmt.Errorf("error seting 'atual' to false: %q", err)
		}

		if err := tx.Model(dto.AgencyMonthlyInfoDTO{}).Create(coletas).Error; err != nil {
			return fmt.Errorf("error inserting 'coleta': %q", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error performing transaction: %q", err)
	}

	return nil
}

func (p *PostgresDB) StorePaychecks(paychecks []models.Paycheck, remunerations []models.PaycheckItem) error {
	// Armazenando contracheques
	var payc []*dto.PaycheckDTO
	for _, pc := range paychecks {
		payc = append(payc, dto.NewPaycheckDTO(pc))
	}
	if err := p.db.Model(dto.PaycheckDTO{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "orgao"}, {Name: "mes"}, {Name: "ano"}, {Name: "id"}},
		UpdateAll: true,
	}).Create(payc).Error; err != nil {
		return fmt.Errorf("error inserting 'contracheques': %w", err)
	}

	// Armazenando o detalhamento das remunerações
	if len(remunerations) != 0 {
		var rem []*dto.PaycheckItemDTO
		for _, r := range remunerations {
			rem = append(rem, dto.NewPaycheckItemDTO(r))
		}
		if err := p.db.Model(dto.PaycheckItemDTO{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "orgao"}, {Name: "mes"}, {Name: "ano"}, {Name: "id"}, {Name: "id_contracheque"}},
			UpdateAll: true,
		}).CreateInBatches(rem, 5000).Error; err != nil {
			return fmt.Errorf("error inserting 'remuneracoes': %w", err)
		}
	}

	return nil
}

func (p *PostgresDB) GetStateAgencies(uf string) ([]models.Agency, error) {
	uf = strings.ToUpper(uf)
	var dtoOrgaos []dto.AgencyDTO
	if err := p.db.Model(&dto.AgencyDTO{}).Where("jurisdicao = 'Estadual' AND uf = ?", uf).Find(&dtoOrgaos).Error; err != nil {
		return nil, fmt.Errorf("error getting agencies: %q", err)
	}

	var orgaos []models.Agency
	for _, dtoOrgao := range dtoOrgaos {
		orgao, err := dtoOrgao.ConvertToModel()
		if err != nil {
			return nil, fmt.Errorf("error converting agency dto to model: %q", err)
		}
		orgaos = append(orgaos, *orgao)
	}
	return orgaos, nil
}

func (p *PostgresDB) GetOPJ(group string) ([]models.Agency, error) {
	var dtoOrgaos []dto.AgencyDTO
	group = strings.ToLower(group)
	if err := p.db.Model(&dto.AgencyDTO{}).Where("LOWER(jurisdicao) = ?", group).Find(&dtoOrgaos).Error; err != nil {
		return nil, fmt.Errorf("error getting agencies by type: %q", err)
	}

	var orgaos []models.Agency
	for _, dtoOrgao := range dtoOrgaos {
		orgao, err := dtoOrgao.ConvertToModel()
		if err != nil {
			return nil, fmt.Errorf("error converting agency dto to model: %q", err)
		}
		orgaos = append(orgaos, *orgao)
	}
	return orgaos, nil
}

func (p *PostgresDB) StoreRemunerations(remu models.Remunerations) error {
	remuneracoes := dto.NewRemunerationsDTO(remu)
	if err := p.db.Model(dto.RemunerationsDTO{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id_orgao"}, {Name: "mes"}, {Name: "ano"}},
		UpdateAll: true,
	}).Create(remuneracoes).Error; err != nil {
		return fmt.Errorf("error inserting 'remuneracoes_zips': %q", err)
	}
	return nil
}

func (p *PostgresDB) GetAgenciesCount() (int, error) {
	var count int64
	if err := p.db.Model(&dto.AgencyDTO{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error getting agencies count: %q", err)
	}
	return int(count), nil
}

func (p *PostgresDB) GetNumberOfMonthsCollected() (int, error) {
	var count int64
	if err := p.db.Model(&dto.AgencyMonthlyInfoDTO{}).Where("atual = true").Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error getting agencies count: %q", err)
	}
	return int(count), nil
}

func (p *PostgresDB) GetAgenciesByUF(uf string) ([]models.Agency, error) {
	var dtoOrgaos []dto.AgencyDTO
	uf = strings.ToUpper(uf)
	if err := p.db.Model(&dto.AgencyDTO{}).Where("uf = ?", uf).Find(&dtoOrgaos).Error; err != nil {
		return nil, fmt.Errorf("error getting agencies: %q", err)
	}
	var orgaos []models.Agency
	for _, dtoOrgao := range dtoOrgaos {
		orgao, err := dtoOrgao.ConvertToModel()
		if err != nil {
			return nil, fmt.Errorf("error converting agency dto to model: %q", err)
		}
		orgaos = append(orgaos, *orgao)
	}
	return orgaos, nil
}

func (p *PostgresDB) GetAgency(aid string) (*models.Agency, error) {
	var dtoOrgao dto.AgencyDTO
	aid = strings.ToLower(aid)
	if err := p.db.Model(&dto.AgencyDTO{}).Where("id = ?", aid).First(&dtoOrgao).Error; err != nil {
		return nil, fmt.Errorf("error getting agency '%s': %q", aid, err)
	}
	orgao, err := dtoOrgao.ConvertToModel()
	if err != nil {
		return nil, fmt.Errorf("error converting agency dto to model: %q", err)
	}
	return orgao, nil
}

func (p *PostgresDB) GetAllAgencies() ([]models.Agency, error) {
	var dtoOrgaos []dto.AgencyDTO
	if err := p.db.Model(&dto.AgencyDTO{}).Find(&dtoOrgaos).Error; err != nil {
		return nil, fmt.Errorf("error getting agencies: %q", err)
	}
	var orgaos []models.Agency
	for _, dtoOrgao := range dtoOrgaos {
		orgao, err := dtoOrgao.ConvertToModel()
		if err != nil {
			return nil, fmt.Errorf("error converting agency dto to model: %q", err)
		}
		orgaos = append(orgaos, *orgao)
	}
	return orgaos, nil
}

func (p *PostgresDB) GetMonthlyInfo(agencies []models.Agency, year int) (map[string][]models.AgencyMonthlyInfo, error) {
	var results = make(map[string][]models.AgencyMonthlyInfo)
	//Mapeando os órgãos
	for _, agency := range agencies {
		var dtoAgmis []dto.AgencyMonthlyInfoDTO
		//Pegando as coletas do postgres, filtrando por órgão, ano e a coleta atual.
		m := p.db.Model(&dto.AgencyMonthlyInfoDTO{})
		m = m.Where("id_orgao = ? AND ano = ? AND atual = TRUE AND (procinfo::text = 'null' OR procinfo IS NULL) ", agency.ID, year)
		m = m.Order("mes ASC")
		if err := m.Find(&dtoAgmis).Error; err != nil {
			return nil, fmt.Errorf("error getting monthly info: %q", err)
		}
		//Convertendo os DTO's para modelos
		for _, dtoAgmi := range dtoAgmis {
			agmi, err := dtoAgmi.ConvertToModel()
			if err != nil {
				return nil, fmt.Errorf("error converting dto to model: %q", err)
			}
			results[agency.ID] = append(results[agency.ID], *agmi)
		}
	}
	return results, nil
}

func (p *PostgresDB) GetAnnualSummary(agency string) ([]models.AnnualSummary, error) {
	var dtoAgmi dto.AgencyMonthlyInfoDTO
	var dtoAmis []dto.AnnualSummaryDTO
	agency = strings.ToLower(agency)
	query := `
		ano,
		id_orgao,
		TRUNC(AVG((sumario -> 'membros')::text::int)) AS media_num_membros,
		SUM((sumario -> 'membros')::text::int) AS total_num_membros,
		SUM(CAST(sumario -> 'remuneracao_base' ->> 'total' AS DECIMAL)) AS remuneracao_base,
		SUM(CAST(sumario -> 'outras_remuneracoes' ->> 'total' AS DECIMAL)) AS outras_remuneracoes,
		SUM(CAST(sumario -> 'descontos' ->> 'total' AS DECIMAL)) AS descontos,
		SUM(CAST(sumario -> 'remuneracoes' ->> 'total' AS DECIMAL)) AS remuneracoes,
		SUM(CAST(sumario -> 'resumo_rubricas' ->> 'auxilio_alimentacao' AS DECIMAL)) AS auxilio_alimentacao,
		SUM(CAST(sumario -> 'resumo_rubricas' ->> 'outras' AS DECIMAL)) AS outras,
		COUNT(*) AS meses_com_dados`
	m := p.db.Model(&dtoAgmi).Select(query)
	m = m.Where("id_orgao = ? AND atual = TRUE AND (procinfo::text = 'null' OR procinfo IS NULL) ", agency)
	m = m.Group("ano, id_orgao").Order("ano ASC")
	if err := m.Scan(&dtoAmis).Error; err != nil {
		return nil, fmt.Errorf("error getting annual monthly info: %q", err)
	}
	var amis []models.AnnualSummary
	for _, dtoAmi := range dtoAmis {
		amis = append(amis, *dtoAmi.ConvertToModel())
	}
	return amis, nil
}

func (p *PostgresDB) GetOMA(month int, year int, agency string) (*models.AgencyMonthlyInfo, *models.Agency, error) {
	var dtoAgmi dto.AgencyMonthlyInfoDTO
	id := fmt.Sprintf("%s/%s/%d", strings.ToLower(agency), dto.AddZeroes(month), year)
	m := p.db.Model(dto.AgencyMonthlyInfoDTO{}).Where("id = ? AND atual = true", id).First(&dtoAgmi)
	if err := m.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, fmt.Errorf("there is no data with this parameters")
		}
		return nil, nil, fmt.Errorf("error getting 'coletas' with id (%s): %q", id, err)
	}
	agmi, err := dtoAgmi.ConvertToModel()
	if err != nil {
		return nil, nil, fmt.Errorf("error converting agmi dto to model: %q", err)
	}
	agencyObject, err := p.GetAgency(agency)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting 'orgaos' with id (%s): %q", agency, err)
	}
	return agmi, agencyObject, nil
}

func (p *PostgresDB) GetGeneralMonthlyInfosFromYear(year int) ([]models.GeneralMonthlyInfo, error) {
	var dtoAgmi dto.AgencyMonthlyInfoDTO
	var dtoGmi []dto.GeneralMonthlyInfoDTO
	query := `
		mes,
		SUM((sumario -> 'membros')::text::int) AS num_membros,
		SUM(CAST(sumario -> 'remuneracao_base' ->> 'total' AS DECIMAL)) AS remuneracao_base,
		SUM(CAST(sumario -> 'outras_remuneracoes' ->> 'total' AS DECIMAL)) AS outras_remuneracoes,
		SUM(CAST(sumario -> 'descontos' ->> 'total' AS DECIMAL)) AS descontos,
		SUM(CAST(sumario -> 'remuneracoes' ->> 'total' AS DECIMAL)) AS remuneracoes,
		SUM(CAST(sumario -> 'resumo_rubricas' ->> 'auxilio_alimentacao' AS DECIMAL)) AS auxilio_alimentacao,
		SUM(CAST(sumario -> 'resumo_rubricas' ->> 'outras' AS DECIMAL)) AS outras`
	m := p.db.Model(&dtoAgmi).Select(query)
	m = m.Where("ano = ? AND atual=true AND (procinfo IS NULL OR procinfo::text = 'null')", year)
	m = m.Group("mes").Order("mes ASC")
	if err := m.Scan(&dtoGmi).Error; err != nil {
		return nil, fmt.Errorf("error getting general remuneration value: %q", err)
	}
	var gmis []models.GeneralMonthlyInfo
	for _, gmi := range dtoGmi {
		gmis = append(gmis, *gmi.ConvertToModel())
	}
	return gmis, nil
}

func (p *PostgresDB) GetFirstDateWithMonthlyInfo() (int, int, error) {
	var dtoAgmi dto.AgencyMonthlyInfoDTO
	var year, month int
	m := p.db.Model(&dtoAgmi).Select("MIN(ano), MIN(mes)")
	m = m.Where("atual=true AND (procinfo IS NULL OR procinfo::text = 'null')")
	m = m.Where("ano = (SELECT min(ano) FROM coletas)")
	if err := m.Row().Scan(&year, &month); err != nil {
		return 0, 0, fmt.Errorf("error getting first date with monthly info: %q", err)
	}
	return month, year, nil
}

func (p *PostgresDB) GetLastDateWithMonthlyInfo() (int, int, error) {
	var dtoAgmi dto.AgencyMonthlyInfoDTO
	var year, month int
	m := p.db.Model(&dtoAgmi).Select("MAX(ano),MAX(mes)")
	m = m.Where("atual=true AND (procinfo IS NULL OR procinfo::text='null')")
	m = m.Where("ano = (SELECT MAX(ano) FROM coletas)")
	if err := m.Row().Scan(&year, &month); err != nil {
		return 0, 0, fmt.Errorf("error getting last date with monthly info: %q", err)
	}
	return month, year, nil
}

func (p *PostgresDB) GetGeneralMonthlyInfo() (float64, error) {
	var dtoAgmi dto.AgencyMonthlyInfoDTO
	var value float64
	query := `
		COALESCE(
			SUM(
				CAST(sumario -> 'remuneracao_base' ->> 'total' AS DECIMAL) +
				CAST(sumario -> 'outras_remuneracoes' ->> 'total' AS DECIMAL)
			), 0
		)
		`
	m := p.db.Model(&dtoAgmi).Select(query)
	m = m.Where("atual=true AND (procinfo IS NULL OR procinfo::text = 'null')")
	if err := m.Scan(&value).Error; err != nil {
		return 0, fmt.Errorf("error getting general remuneration value: %q", err)
	}
	return value, nil
}

func (p *PostgresDB) GetIndexInformation(name string, month, year int) (map[string][]models.IndexInformation, error) {
	// name: ID do órgão (e.g. "trt12") ou jurisdição.
	groupMap := map[string]struct{}{"eleitoral": {}, "ministério": {}, "estadual": {}, "trabalho": {}, "federal": {}, "militar": {}, "superior": {}, "conselho": {}}
	params := []interface{}{}

	// somente considerar os dados da coleta mais recente de cada órgão.
	// lembrar que a gente guarda um histórico de coletas (revisões)
	query := "INNER JOIN orgaos ON coletas.id_orgao = orgaos.id AND coletas.atual = true"

	// verificar se devemos considerar o ano como parâmetro e adicionar a query.
	if year != 0 {
		query += " AND coletas.ano = ?"
		params = append(params, year)

		// verificar se devemos considerar o mês. como parâmetro e adicionar a query.
		// só verificamos o mês se o ano for passado.
		if month != 0 {
			query += " AND coletas.mes = ?"
			params = append(params, month)
		}
	}
	var dtoIndex []dto.IndexInformation
	var d *gorm.DB
	_, porJurisdicao := groupMap[strings.ToLower(name)]
	if porJurisdicao {
		// Consultando e mapeando os índices e metadados por jurisdição.
		// Para tal, precisamos realizar um join com a tabela de órgãos.
		query += " AND orgaos.jurisdicao = ?"
		params = append(params, name)
	} else {
		if name != "" {
			// Consultando e mapeando os índices e metadados por id do órgão
			query += " AND coletas.id_orgao = ?"
			params = append(params, name)
		}
	}
	d = p.db.Model(&dtoIndex).Select("coletas.*, orgaos.jurisdicao as jurisdicao").Joins(query, params...)
	if err := d.Scan(&dtoIndex).Error; err != nil {
		return nil, fmt.Errorf("error getting all indexes: %w", err)
	}
	// Agrupando os índices por órgão
	indexes := make(map[string][]models.IndexInformation)
	for _, d := range dtoIndex {
		d.Score.EasinessScore = calcEasinessScore(d.ID, d.Score.EasinessScore)
		indexes[d.ID] = append(indexes[d.ID], *d.ConvertToModel())
	}
	return indexes, nil
}

func (p *PostgresDB) GetAllAgencyCollection(agency string) ([]models.AgencyMonthlyInfo, error) {
	var dtoAgmis []dto.AgencyMonthlyInfoDTO
	//Pegando todas as coletas atuais de um determinado órgão.
	m := p.db.Model(&dto.AgencyMonthlyInfoDTO{})
	m = m.Where("id_orgao = ? AND atual = TRUE", agency)
	m = m.Order("(ano, mes) ASC")
	if err := m.Find(&dtoAgmis).Error; err != nil {
		return nil, fmt.Errorf("error getting all agency collections: %q", err)
	}

	var collections []models.AgencyMonthlyInfo
	for _, dtoAgmi := range dtoAgmis {
		agmi, err := dtoAgmi.ConvertToModel()
		if err != nil {
			return nil, fmt.Errorf("error converting dto to model: %q", err)
		}
		agmi.Score.EasinessScore = calcEasinessScore(agency, agmi.Score.EasinessScore)
		collections = append(collections, *agmi)
	}
	return collections, nil
}

// Verificamos se o órgão pertence ao painel do CNJ (ou se é um ministério público)
// O índice de facilidade para os órgãos do CNJ é padronizado, mesmo quando não há dados para o mês.
// obs.: o "STF" é o único tribunal que monitoramos e que não pertence ao CNJ
func calcEasinessScore(agency string, easinessScore float64) float64 {
	if !strings.Contains(strings.ToLower(agency), "mp") && agency != "stf" {
		return 0.5
	} else {
		return easinessScore
	}
}

func (p *PostgresDB) GetPaychecks(agency models.Agency, year int) ([]models.Paycheck, error) {
	var results []models.Paycheck
	var dtoPaychecks []dto.PaycheckDTO
	//Pegando os contracheques do postgres, filtrando por órgão e ano
	m := p.db.Model(&dto.PaycheckDTO{})
	m = m.Where("orgao = ? AND ano = ? ", agency.ID, year)
	m = m.Order("mes, id ASC")
	if err := m.Find(&dtoPaychecks).Error; err != nil {
		return nil, fmt.Errorf("error getting paychecks: %q", err)
	}
	//Convertendo os DTO's para modelos
	for _, dtoPaycheck := range dtoPaychecks {
		p := dtoPaycheck.ConvertToModel()
		results = append(results, *p)
	}
	return results, nil
}

func (p *PostgresDB) GetPaycheckItems(agency models.Agency, year int) ([]models.PaycheckItem, error) {
	var results []models.PaycheckItem
	var dtoPaycheckItems []dto.PaycheckItemDTO
	//Pegando as remuneracoes do postgres, filtrando por órgão e ano
	m := p.db.Model(&dto.PaycheckItemDTO{})
	m = m.Where("orgao = ? AND ano = ?", agency.ID, year)
	m = m.Order("mes, id_contracheque, id ASC")
	if err := m.Find(&dtoPaycheckItems).Error; err != nil {
		return nil, fmt.Errorf("error getting paycheck items: %q", err)
	}
	//Convertendo os DTO's para modelos
	for _, dtoPaycheckItem := range dtoPaycheckItems {
		p := dtoPaycheckItem.ConvertToModel()
		results = append(results, *p)
	}
	return results, nil
}
