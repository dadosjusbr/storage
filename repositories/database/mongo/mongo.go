package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//collection is a private interface to create a mongo's ReplaceOne method and their signatures to be used and tested.
type collection interface {
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

//DBClient is a mongodb Client instance
type DBClient struct {
	mgoClient      *mongo.Client
	dbName         string
	monthlyInfoCol string
	agencyCol      string
	packageCol     string
	revCol         string
	col            collection
}

//NewDBClient instantiates a mongo new client, but will not connect to the specified URL. Please use Client.Connect before using the client.
func NewMongoDB(url, dbName, monthlyInfoCol, agencyCol string, packageCol string, revCol string) (*DBClient, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	return &DBClient{mgoClient: client,
		dbName:         dbName,
		monthlyInfoCol: monthlyInfoCol,
		agencyCol:      agencyCol,
		packageCol:     packageCol,
		revCol:         revCol}, nil
}

//Connect establishes a connection to MongoDB using the previously specified URL
func (c *DBClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := c.mgoClient.Connect(ctx); err != nil {
		return fmt.Errorf("error connection with mongo:%q", err)
	}
	return nil
}

//Disconnect closes the connections to MongoDB. It does nothing if the connection had already been closed.
func (c *DBClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	c.col = nil
	return c.mgoClient.Disconnect(ctx)
}

// GetOPE return agmi info to build first screen
func (c *DBClient) GetOPE(uf string, year int) ([]models.Agency, error) {
	allAgencies, err := c.GetAgencies(uf)
	if err != nil {
		return nil, fmt.Errorf("GetOPE() error: %q", err)
	}
	return allAgencies, nil
}

// GetAgenciesCount Return the Agencies amount
func (c *DBClient) GetAgenciesCount() (int64, error) {
	c.Collection(c.agencyCol)
	itemCount, err := c.col.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return itemCount, fmt.Errorf("Error in result %v", err)
	}
	return itemCount, nil
}

// GetNumberOfMonthsCollected Return the number of months collected
func (c *DBClient) GetNumberOfMonthsCollected() (int64, error) {
	c.Collection(c.monthlyInfoCol)
	itemCount, err := c.col.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return itemCount, fmt.Errorf("Error in result %v", err)
	}
	return itemCount, nil
}

//GetAgencies Return UF Agencies
func (c *DBClient) GetAgencies(uf string) ([]models.Agency, error) {
	c.Collection(c.agencyCol)
	resultAgencies, err := c.col.Find(context.TODO(), bson.M{"$and": []bson.M{{"uf": uf}}}, nil)
	if err != nil {
		return nil, fmt.Errorf("error in getAgencies %v", err)
	}
	var allAgencies []models.Agency
	resultAgencies.All(context.TODO(), &allAgencies)
	if err := resultAgencies.Err(); err != nil {
		return nil, fmt.Errorf("error in getAgencies %v", err)
	}
	return allAgencies, nil
}

//GetAgency Return Agency that match ID.
func (c *DBClient) GetAgency(aid string) (*models.Agency, error) {
	c.Collection(c.agencyCol)
	var Ag models.Agency
	if err := c.col.FindOne(context.TODO(), bson.D{{Key: "aid", Value: aid}}).Decode(&Ag); err != nil {
		return nil, fmt.Errorf("Error searching for agency id \"%s\":%q", aid, err)
	}
	return &Ag, nil
}

// GetAllAgencies returns all agencies from AG collection
func (c *DBClient) GetAllAgencies() ([]models.Agency, error) {
	c.Collection(c.agencyCol)
	var agencies []models.Agency
	agCursor, err := c.col.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("Error while indexing Agencies: %q", err)
	}
	if err := agCursor.All(context.TODO(), &agencies); err != nil {
		return nil, fmt.Errorf("Error while indexing Agencies: %q", err)
	}
	for i := range agencies {
		agencies[i].FlagURL = "v1/orgao/" + agencies[i].ID
	}
	return agencies, nil
}

//GetMonthlyInfo return summarized monthlyInfo for each agency in agencies in a specific year
func (c *DBClient) GetMonthlyInfo(agencies []models.Agency, year int) (map[string][]models.AgencyMonthlyInfo, error) {
	var result = make(map[string][]models.AgencyMonthlyInfo)
	c.Collection(c.monthlyInfoCol)
	opts := &options.FindOptions{}
	opts.SetSort(bson.M{"month": 1})
	for _, agency := range agencies {
		resultMonthly, err := c.col.Find(
			context.TODO(),
			bson.D{{Key: "aid", Value: agency.ID}, {Key: "year", Value: year}},
			opts,
		)
		if err != nil {
			return nil, fmt.Errorf("Error in GetMonthlyInfo %v", err)
		}
		var mr []models.AgencyMonthlyInfo
		resultMonthly.All(context.TODO(), &mr)
		result[agency.ID] = mr
	}
	return result, nil
}

// GetMonthlyInfoSummary returns summarized monthlyInfo for each agency in agencies in a specific year with packages
func (c *DBClient) GetMonthlyInfoSummary(agencies []models.Agency, year int) (map[string][]models.AgencyMonthlyInfo, error) {
	var result = make(map[string][]models.AgencyMonthlyInfo)
	c.Collection(c.monthlyInfoCol)
	for _, agency := range agencies {
		resultMonthly, err := c.col.Find(
			context.TODO(), bson.D{{Key: "aid", Value: agency.ID}, {Key: "year", Value: year}},
			options.Find().SetProjection(bson.D{{"aid", 1}, {"year", 1}, {"month", 1}, {"summary", 1}, {"package", 1}}))
		if err != nil {
			return nil, fmt.Errorf("Error in GetMonthlyInfo %v", err)
		}
		var mr []models.AgencyMonthlyInfo
		resultMonthly.All(context.TODO(), &mr)
		result[agency.ID] = mr
	}
	return result, nil
}

//GetOMA Search if DB has a match for filters
func (c *DBClient) GetOMA(month int, year int, agency string) (*models.AgencyMonthlyInfo, *models.Agency, error) {
	c.Collection(c.monthlyInfoCol)
	var resultMonthly models.AgencyMonthlyInfo
	err := c.col.FindOne(context.TODO(), bson.D{{Key: "aid", Value: agency}, {Key: "year", Value: year}, {Key: "month", Value: month}}).Decode(&resultMonthly)
	// ErrNoDocuments means that the filter did not match any documents in the collection
	if err == mongo.ErrNoDocuments {
		return nil, nil, storage.ErrNothingFound
	}
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch AGMI information from bd: %w", err)
	}
	agencyObject, err := c.GetAgency(agency)
	if err != nil {
		return nil, nil, fmt.Errorf("error in GetAgency: %w", err)
	}
	return &resultMonthly, agencyObject, nil
}

func (c *DBClient) Store(agmi models.AgencyMonthlyInfo) error {
	// Armazenando sempre duas cópias no novo item. Tomamos a decisão de
	// armazenar uma cópia para evitar a complexidade e a perda de desempenho de
	// gerenciar a manutenção de apenas uma cópia entre as coleções (tirar de uma
	// coleção e colocar em outra).

	// ## Armazenando versão corrente
	c.Collection(c.monthlyInfoCol)
	key := bson.D{{Key: "aid", Value: agmi.AgencyID}, {Key: "year", Value: agmi.Year}, {Key: "month", Value: agmi.Month}}
	opts := options.Replace().SetUpsert(true)
	if _, err := c.col.ReplaceOne(context.TODO(), key, agmi, opts); err != nil {
		return fmt.Errorf("error trying to update current monthly info with value {%v}: %q", agmi, err)
	}
	// ## Armazenando revisão.
	c.Collection(c.revCol)
	rev := models.MonthlyInfoVersion{
		AgencyID:  agmi.AgencyID,
		Month:     agmi.Month,
		Year:      agmi.Year,
		VersionID: time.Now().Unix(),
		Version:   agmi,
	}
	if _, err := c.col.InsertOne(context.TODO(), rev); err != nil {
		return fmt.Errorf("error trying to insert monthly info revision with value {%v}: %q", agmi, err)
	}
	return nil
}

func (c *DBClient) StorePackage(newPackage models.Package) error {
	c.Collection(c.packageCol)
	filter := bson.M{
		"aid":   newPackage.AgencyID,
		"month": newPackage.Month,
		"year":  newPackage.Year,
	}
	update := bson.M{
		"aid":     newPackage.AgencyID,
		"group":   newPackage.Group,
		"month":   newPackage.Month,
		"year":    newPackage.Year,
		"package": newPackage.Package,
	}
	opts := options.Replace().SetUpsert(true)
	_, err := c.col.ReplaceOne(context.TODO(), filter, update, opts)
	if err != nil {
		return fmt.Errorf("error while updating a agreggation: %q", err)
	}
	return nil
}

func (p *DBClient) StoreRemunerations(remu models.Remunerations) error {
	//TODO implement me
	panic("implement me")
}

//Collection Changes active collection
func (c *DBClient) Collection(collectionName string) {
	c.col = c.mgoClient.Database(c.dbName).Collection(collectionName)
}

//GetGeneralMonthlyInfosFromYear return the sum from all remuneration info from all months of a given year
func (c *DBClient) GetGeneralMonthlyInfosFromYear(year int) ([]models.GeneralMonthlyInfo, error) {
	c.Collection(c.monthlyInfoCol)
	resultMonthly, err := c.col.Aggregate(context.TODO(),
		mongo.Pipeline{
			bson.D{{"$match",
				bson.D{{"year", year}}}},
			bson.D{{"$group",
				bson.D{
					{"_id", "$month"},
					{"base_remuneration", bson.D{{"$sum", "$summary.base_remuneration.total"}}},
					{"other_remunerations", bson.D{{"$sum", "$summary.other_remunerations.total"}}},
					{"count", bson.D{{"$sum", "$summary.count"}}}}}},
			bson.D{{"$sort",
				bson.D{
					{"_id", 1}}}}})
	if err != nil {
		return nil, fmt.Errorf("Error in GetMonthlyInfo %v", err)
	}
	var mr []models.GeneralMonthlyInfo
	resultMonthly.All(context.TODO(), &mr)
	return mr, nil
}

//GetFirstDateWithMonthlyInfo return the initial year and month with collected data
func (c *DBClient) GetFirstDateWithMonthlyInfo() (int, int, error) {
	var resultMonthly models.AgencyMonthlyInfo
	c.Collection(c.monthlyInfoCol)
	firstDateQueryOptions := options.FindOne().SetSort(bson.D{{Key: "year", Value: +1}, {Key: "month", Value: +1}})
	err := c.col.FindOne(
		context.TODO(),
		bson.D{{"month", bson.M{"$exists": true}}}, firstDateQueryOptions).Decode(&resultMonthly)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		return 0, 0, fmt.Errorf("Error in result %v", err)
	}
	return resultMonthly.Month, resultMonthly.Year, nil
}

//GetLastDateWithMonthlyInfo return the latest year and month with collected data
func (c *DBClient) GetLastDateWithMonthlyInfo() (int, int, error) {
	var resultMonthly models.AgencyMonthlyInfo
	c.Collection(c.monthlyInfoCol)
	lastDateQueryOptions := options.FindOne().SetSort(bson.D{{Key: "year", Value: -1}, {Key: "month", Value: -1}})
	err := c.col.FindOne(
		context.TODO(),
		bson.D{{"month", bson.M{"$exists": true}}}, lastDateQueryOptions).Decode(&resultMonthly)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		return 0, 0, fmt.Errorf("Error in result %v", err)
	}
	return resultMonthly.Month, resultMonthly.Year, nil
}

//GetRemunerationSummary return the amount  of remuneration records from all agencies and the final remuneration value
func (c *DBClient) GetRemunerationSummary() (*models.RemmunerationSummary, error) {
	c.Collection(c.monthlyInfoCol)
	// NOTA: Não estamos usando a função de agregação do mongo pois a camada gratuita do Atlas não
	// permite a utilização de filtros enquanto estamos agregando.
	var amis []models.AgencyMonthlyInfo
	resultMonthly, err := c.col.Find(
		context.TODO(), bson.D{},
		options.Find().SetProjection(bson.D{{Key: "summary", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("error querying data: %q", err)
	}
	if err := resultMonthly.All(context.TODO(), &amis); err != nil {
		log.Printf("Error querying data: %q", err)
		return nil, fmt.Errorf("error querying data: %q", err)
	}
	var result struct {
		BaseRemuneration   float64 `json:"base_remuneration" bson:"base_remuneration,omitempty"`     //  Statistics (Max, Min, Median, Total)
		OtherRemunerations float64 `json:"other_remunerations" bson:"other_remunerations,omitempty"` //  Statistics (Max, Min, Median, Total)
		Count              int     `bson:"count,omitempty"`
	}
	for _, ami := range amis {
		result.Count++
		result.BaseRemuneration += ami.Summary.BaseRemuneration.Total
		result.OtherRemunerations += ami.Summary.OtherRemunerations.Total
	}
	return &models.RemmunerationSummary{Count: result.Count, Value: result.BaseRemuneration + result.OtherRemunerations}, nil
}

//GetAggregation return an aggregation who attends the given params
func (c *DBClient) GetPackage(pkgOpts models.PackageFilterOpts) (*models.Package, error) {
	c.Collection(c.packageCol)
	var pkg models.Package
	err := c.col.FindOne(context.TODO(), bson.D{{
		Key: "aid", Value: pkgOpts.AgencyID},
		{Key: "year", Value: pkgOpts.Year},
		{Key: "month", Value: pkgOpts.Month},
		{Key: "group", Value: pkgOpts.Group}}).Decode(&pkg)
	if err != nil {
		return nil, fmt.Errorf("Error searching for datapackage: %q", err)
	}
	return &pkg, nil
}

func (c *DBClient) GetGeneralMonthlyInfo() (float64, error) {
	//TODO implement me
	panic("implement me")
}
