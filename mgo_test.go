package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/dadosjusbr/coletores"
	"github.com/ncw/swift"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type checkCollection struct {
	t      *testing.T
	filter bson.D
	value  interface{}
	opts   []*options.ReplaceOptions
	check  bool
	err    bool
}

type checkStorage struct {
	t           *testing.T
	check       bool
	err         bool
	container   string
	objectName  string
	contents    io.Reader
	checkHash   bool
	Hash        string
	contentType string
	h           swift.Headers
}

func makePointer(x float64) *float64 {
	return &x
}

func newString(s string) *string {
	return &s
}

var (
	emp4Row = []coletores.Employee{
		{
			Reg:       "",
			Name:      "Abiaci De Carvalho Silva",
			Role:      "Inativo",
			Type:      newString("servidor"),
			Workplace: "inativo",
			Active:    false,
			Income: &coletores.IncomeDetails{
				Total: 30368.59,
				Wage:  makePointer(7000),
				Perks: &coletores.Perks{
					Total: 600,
				},
				Other: &coletores.Funds{
					Total:            100,
					PersonalBenefits: makePointer(7475.71),
					EventualBenefits: makePointer(0),
					PositionOfTrust:  makePointer(5990.88),
					Gratification:    makePointer(0),
					OriginPosition:   makePointer(0),
				},
			},
			Discounts: &coletores.Discount{
				Total:            8930.05,
				PrevContribution: makePointer(2719.5),
				CeilRetention:    makePointer(0),
				IncomeTax:        makePointer(6210.55),
				Others: map[string]float64{
					"Sundry": 0,
				},
			},
		},
		{
			Reg:       "",
			Name:      "Abraao Falcao De Carvalho",
			Role:      "Promotor Eleitoral",
			Type:      newString("servidor"),
			Workplace: "10ª zona eleitoral - guarabira/pb",
			Active:    true,
			Income: &coletores.IncomeDetails{
				Total: 10000,
				Wage:  makePointer(5000),
				Perks: &coletores.Perks{
					Total: 200,
				},
				Other: &coletores.Funds{
					Total:            500,
					PersonalBenefits: makePointer(0),
					EventualBenefits: makePointer(0),
					PositionOfTrust:  makePointer(4631.61),
					Gratification:    makePointer(0),
					OriginPosition:   makePointer(0),
				},
			},
			Discounts: &coletores.Discount{
				Total:            405.98,
				PrevContribution: makePointer(0),
				CeilRetention:    makePointer(0),
				IncomeTax:        makePointer(405.98),
				Others: map[string]float64{
					"Sundry": 0,
				},
			},
		},
		{
			Reg:       "",
			Name:      "Abraao Galcao",
			Role:      "Promotor Eleitoral",
			Type:      newString("membro"),
			Workplace: "10ª zona eleitoral - guarabira/pb",
			Active:    false,
			Income: &coletores.IncomeDetails{
				Total: 10000,
				Wage:  makePointer(5000),
				Perks: &coletores.Perks{
					Total: 200,
				},
				Other: &coletores.Funds{
					Total:            500,
					PersonalBenefits: makePointer(0),
					EventualBenefits: makePointer(0),
					PositionOfTrust:  makePointer(4631.61),
					Gratification:    makePointer(0),
					OriginPosition:   makePointer(0),
				},
			},
			Discounts: &coletores.Discount{
				Total:            405.98,
				PrevContribution: makePointer(0),
				CeilRetention:    makePointer(0),
				IncomeTax:        makePointer(405.98),
				Others: map[string]float64{
					"Sundry": 0,
				},
			},
		},
		{
			Reg:       "",
			Name:      "Abraao Halcao",
			Role:      "Promotor Eleitoral",
			Type:      newString("membro"),
			Workplace: "10ª zona eleitoral - guarabira/pb",
			Active:    true,
			Income: &coletores.IncomeDetails{
				Total: 10000,
				Wage:  makePointer(5000),
				Perks: &coletores.Perks{
					Total: 200,
				},
				Other: &coletores.Funds{
					Total:            500,
					PersonalBenefits: makePointer(0),
					EventualBenefits: makePointer(0),
					PositionOfTrust:  makePointer(4631.61),
					Gratification:    makePointer(0),
					OriginPosition:   makePointer(0),
				},
			},
			Discounts: &coletores.Discount{
				Total:            405.98,
				PrevContribution: makePointer(0),
				CeilRetention:    makePointer(0),
				IncomeTax:        makePointer(405.98),
				Others: map[string]float64{
					"Sundry": 0,
				},
			},
		},
	}

	summFor1RowGeneral = Summary{
		Count:  1,
		Wage:   DataSummary{Max: 7000.00, Min: 7000.00, Average: 7000.00, Total: 7000.00},
		Perks:  DataSummary{Max: 600.00, Min: 600.00, Average: 600.00, Total: 600.00},
		Others: DataSummary{Max: 100.00, Min: 100.00, Average: 100.00, Total: 100.00},
	}
	summFor1RowNull = Summary{
		Count:  1,
		Wage:   DataSummary{Max: 0, Min: 0, Average: 0, Total: 0},
		Perks:  DataSummary{Max: 0, Min: 0, Average: 0, Total: 0},
		Others: DataSummary{Max: 0, Min: 0, Average: 0, Total: 0},
	}

	summFor4RowGeneral = Summary{
		Count:  4,
		Wage:   DataSummary{Max: 7000.00, Min: 5000.00, Average: 5500.00, Total: 22000.00},
		Perks:  DataSummary{Max: 600.00, Min: 200.00, Average: 300.00, Total: 1200.00},
		Others: DataSummary{Max: 500.00, Min: 100.00, Average: 400.00, Total: 1600.00},
	}
	summFor4RowSI = Summary{
		Count:  1,
		Wage:   DataSummary{Max: 7000.00, Min: 7000.00, Average: 7000.00, Total: 7000.00},
		Perks:  DataSummary{Max: 600.00, Min: 600.00, Average: 600.00, Total: 600.00},
		Others: DataSummary{Max: 100.00, Min: 100.00, Average: 100.00, Total: 100.00},
	}
	summFor4RowMI = Summary{
		Count:  1,
		Wage:   DataSummary{Max: 5000.00, Min: 5000.00, Average: 5000.00, Total: 5000.00},
		Perks:  DataSummary{Max: 200.00, Min: 200.00, Average: 200.00, Total: 200.00},
		Others: DataSummary{Max: 500.00, Min: 500.00, Average: 500.00, Total: 500.00},
	}
	summFor4RowSA = Summary{
		Count:  1,
		Wage:   DataSummary{Max: 5000.00, Min: 5000.00, Average: 5000.00, Total: 5000.00},
		Perks:  DataSummary{Max: 200.00, Min: 200.00, Average: 200.00, Total: 200.00},
		Others: DataSummary{Max: 500.00, Min: 500.00, Average: 500.00, Total: 500.00},
	}
	summFor4RowMA = Summary{
		Count:  1,
		Wage:   DataSummary{Max: 5000.00, Min: 5000.00, Average: 5000.00, Total: 5000.00},
		Perks:  DataSummary{Max: 200.00, Min: 200.00, Average: 200.00, Total: 200.00},
		Others: DataSummary{Max: 500.00, Min: 500.00, Average: 500.00, Total: 500.00},
	}

	summFor4Row = Summaries{
		General:       summFor4RowGeneral,
		MemberActive:  summFor4RowMA,
		ServantActive: summFor4RowSA,
	}
	summFor1Row = Summaries{
		General: summFor1RowGeneral,
	}

	crawler = coletores.Crawler{CrawlerID: "123132", CrawlerVersion: "v.1"}
	cr      = coletores.CrawlingResult{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler, Employees: emp4Row, Files: []string{"teste.txt", "outroTeste.txt"}}
	agmi    AgencyMonthlyInfo

	backup1 = Backup{URL: "/dadosjusbr/teste.txt", Hash: "0e30309b400c02246b6ac4f461c0fa96"}
	backup2 = Backup{URL: "/dadosjusbr/outroTeste.txt", Hash: "0e30309b400c02246b6ac4f461c0fa96"}
)

// ReplaceOne is a checkCollection func that use same signature of collection interface, which is the same as the method signature with the same name in mongoDb
func (c *checkCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	c.check = true
	if c.err {
		return nil, fmt.Errorf("replace one error")
	}

	assert.Equal(c.t, c.filter, filter)
	assert.Equal(c.t, c.value, replacement)
	assert.Equal(c.t, c.opts, opts)
	return &mongo.UpdateResult{}, nil
}

func (c *checkCollection) calledReplaceOne() bool {
	return c.check
}

func (c *checkCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return nil, nil
}

func (c *checkCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return nil
}

func (cs *checkStorage) ObjectPut(container string, objectName string, contents io.Reader, checkHash bool, Hash string, contentType string, h swift.Headers) (headers swift.Headers, err error) {
	if cs.check {
		assert.Equal(cs.t, cs.container, container)
		assert.Equal(cs.t, cs.objectName, objectName)
		assert.Equal(cs.t, cs.checkHash, checkHash)
		assert.Equal(cs.t, cs.Hash, Hash)
		assert.Equal(cs.t, cs.contentType, contentType)
		assert.Equal(cs.t, cs.h, h)
	}

	if cs.err {
		return nil, fmt.Errorf("Object Put Error")
	}

	return swift.Headers{"Etag": "0e30309b400c02246b6ac4f461c0fa96"}, nil
}

func (cs *checkStorage) ObjectDelete(container string, objectName string) error {
	return nil
}

func TestClient_Store(t *testing.T) {
	err := createFiles(cr.Files)
	assert.NoError(t, err)
	bc := &CloudClient{conn: &checkStorage{check: false, container: "dadosjusbr"}}
	col := checkCollection{
		t:      t,
		filter: bson.D{{Key: "aid", Value: "a"}, {Key: "year", Value: 2019}, {Key: "month", Value: 9}},
		value:  AgencyMonthlyInfo{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler, Employee: emp4Row, Summary: summFor4Row, Backups: []Backup{backup1, backup2}},
		opts:   []*options.ReplaceOptions{options.Replace().SetUpsert(true)},
		err:    false,
	}
	colErr := checkCollection{
		t:      t,
		filter: bson.D{{Key: "aid", Value: "a"}, {Key: "year", Value: 2019}, {Key: "month", Value: 9}},
		value:  AgencyMonthlyInfo{AgencyID: "a", Year: 2019, Month: 9, Crawler: crawler, Employee: emp4Row, Summary: summFor4Row, Backups: []Backup{backup1, backup2}},
		opts:   []*options.ReplaceOptions{options.Replace().SetUpsert(true)},
		err:    true,
	}
	tests := []struct {
		name           string
		col            *checkCollection
		agmi           AgencyMonthlyInfo
		wantErr        bool
		wantReplaceOne bool
	}{
		//Test if everything is OK!
		{name: "ok", col: &col, agmi: agmi, wantErr: false, wantReplaceOne: true},
		// Test if the replaceOne error reflects in store error!
		{name: "replaceOne error", col: &colErr, agmi: agmi, wantErr: true, wantReplaceOne: true},
		// Check if has some connection if mongoDb, if does the collection wont be nil.
		{name: "missing collection error", agmi: agmi, wantErr: true, wantReplaceOne: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{Cloud: bc, Db: &DBClient{}}
			if tt.col != nil {
				c.Db.col = tt.col
			}
			if err := c.Store(tt.agmi); (err != nil) != tt.wantErr {
				t.Errorf("Client.Store() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.col != nil && (tt.wantReplaceOne != tt.col.calledReplaceOne()) {
				t.Errorf("Client.Store() error calledReplaceOne != wantReplaceOne")
			}
		})
	}
	err = deleteFiles(cr.Files)
	assert.NoError(t, err)
}

func Test_Backup(t *testing.T) {
	cs1 := checkStorage{
		t:           t,
		container:   "dadosjusbr",
		objectName:  "teste.txt",
		checkHash:   true,
		Hash:        "",
		contentType: "",
		h:           nil,
		err:         false,
		check:       true,
	}

	cs2Err := checkStorage{
		t:           t,
		container:   "dadosjusbr",
		objectName:  "teste.txt",
		contents:    strings.NewReader("teste.txt"),
		checkHash:   true,
		Hash:        "",
		contentType: "",
		h:           nil,
		err:         true,
		check:       true,
	}

	tests := []struct {
		name    string
		Files   []string
		cs      *checkStorage
		want    []Backup
		wantErr bool
		errMsg  string
	}{
		{name: "OK", Files: []string{"teste.txt"}, want: []Backup{backup1}, cs: &cs1},
		{name: "No Files", Files: []string{}, want: []Backup{}, wantErr: false, errMsg: "is no file"},
		{name: "Object Put Error", Files: []string{"teste.txt"}, want: []Backup{}, wantErr: true, errMsg: "Object Put Error", cs: &cs2Err},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createFiles(tt.Files)
			assert.NoError(t, err)
			bc := &CloudClient{conn: tt.cs}
			if tt.cs != nil {
				bc.container = tt.cs.container
			}
			got, err := bc.Backup(tt.Files, "dest")
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("error = %v, errExpected = %v", err, tt.errMsg)
			}
			if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
				t.Errorf("backup() = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
			err = deleteFiles(tt.Files)
			assert.NoError(t, err)
		})
	}
}

func createFiles(files []string) error {
	for _, f := range files {
		fileNew, err := os.Create(f)
		if err != nil {
			return fmt.Errorf("Error trying to create a file %v", err)
		}
		_, err = fileNew.Write([]byte("Lorem ipsum dolor sit amet consectetuer"))
		if err != nil {
			return fmt.Errorf("Error trying to write a file %v", err)
		}
	}
	return nil
}

func deleteFiles(files []string) error {
	//To test, uncomment line below and insert auth parameters.
	for _, rem := range files {
		err := os.Remove("./" + rem)
		if err != nil {
			return fmt.Errorf("Error trying to delete a file from local %v", err)
		}
	}

	return nil
}
