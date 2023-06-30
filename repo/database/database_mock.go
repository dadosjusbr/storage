// Code generated by MockGen. DO NOT EDIT.
// Source: ./repo/database/interface.go

// Package mock_database is a generated GoMock package.
package database

import (
	reflect "reflect"

	models "github.com/dadosjusbr/storage/models"
	gomock "github.com/golang/mock/gomock"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// Connect mocks base method.
func (m *MockInterface) Connect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockInterfaceMockRecorder) Connect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockInterface)(nil).Connect))
}

// Disconnect mocks base method.
func (m *MockInterface) Disconnect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnect indicates an expected call of Disconnect.
func (mr *MockInterfaceMockRecorder) Disconnect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockInterface)(nil).Disconnect))
}

// GetAgenciesByUF mocks base method.
func (m *MockInterface) GetAgenciesByUF(uf string) ([]models.Agency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAgenciesByUF", uf)
	ret0, _ := ret[0].([]models.Agency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgenciesByUF indicates an expected call of GetAgenciesByUF.
func (mr *MockInterfaceMockRecorder) GetAgenciesByUF(uf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgenciesByUF", reflect.TypeOf((*MockInterface)(nil).GetAgenciesByUF), uf)
}

// GetAgenciesCount mocks base method.
func (m *MockInterface) GetAgenciesCount() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAgenciesCount")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgenciesCount indicates an expected call of GetAgenciesCount.
func (mr *MockInterfaceMockRecorder) GetAgenciesCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgenciesCount", reflect.TypeOf((*MockInterface)(nil).GetAgenciesCount))
}

// GetAgency mocks base method.
func (m *MockInterface) GetAgency(aid string) (*models.Agency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAgency", aid)
	ret0, _ := ret[0].(*models.Agency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgency indicates an expected call of GetAgency.
func (mr *MockInterfaceMockRecorder) GetAgency(aid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgency", reflect.TypeOf((*MockInterface)(nil).GetAgency), aid)
}

// GetAllAgencies mocks base method.
func (m *MockInterface) GetAllAgencies() ([]models.Agency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllAgencies")
	ret0, _ := ret[0].([]models.Agency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllAgencies indicates an expected call of GetAllAgencies.
func (mr *MockInterfaceMockRecorder) GetAllAgencies() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllAgencies", reflect.TypeOf((*MockInterface)(nil).GetAllAgencies))
}

// GetAllAgencyCollection mocks base method.
func (m *MockInterface) GetAllAgencyCollection(agency string) ([]models.AgencyMonthlyInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllAgencyCollection", agency)
	ret0, _ := ret[0].([]models.AgencyMonthlyInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllAgencyCollection indicates an expected call of GetAllAgencyCollection.
func (mr *MockInterfaceMockRecorder) GetAllAgencyCollection(agency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllAgencyCollection", reflect.TypeOf((*MockInterface)(nil).GetAllAgencyCollection), agency)
}

// GetAnnualSummary mocks base method.
func (m *MockInterface) GetAnnualSummary(agency string) ([]models.AnnualSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnualSummary", agency)
	ret0, _ := ret[0].([]models.AnnualSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnualSummary indicates an expected call of GetAnnualSummary.
func (mr *MockInterfaceMockRecorder) GetAnnualSummary(agency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnualSummary", reflect.TypeOf((*MockInterface)(nil).GetAnnualSummary), agency)
}

// GetFirstDateWithMonthlyInfo mocks base method.
func (m *MockInterface) GetFirstDateWithMonthlyInfo() (int, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFirstDateWithMonthlyInfo")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFirstDateWithMonthlyInfo indicates an expected call of GetFirstDateWithMonthlyInfo.
func (mr *MockInterfaceMockRecorder) GetFirstDateWithMonthlyInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFirstDateWithMonthlyInfo", reflect.TypeOf((*MockInterface)(nil).GetFirstDateWithMonthlyInfo))
}

// GetGeneralMonthlyInfo mocks base method.
func (m *MockInterface) GetGeneralMonthlyInfo() (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGeneralMonthlyInfo")
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGeneralMonthlyInfo indicates an expected call of GetGeneralMonthlyInfo.
func (mr *MockInterfaceMockRecorder) GetGeneralMonthlyInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGeneralMonthlyInfo", reflect.TypeOf((*MockInterface)(nil).GetGeneralMonthlyInfo))
}

// GetGeneralMonthlyInfosFromYear mocks base method.
func (m *MockInterface) GetGeneralMonthlyInfosFromYear(year int) ([]models.GeneralMonthlyInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGeneralMonthlyInfosFromYear", year)
	ret0, _ := ret[0].([]models.GeneralMonthlyInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGeneralMonthlyInfosFromYear indicates an expected call of GetGeneralMonthlyInfosFromYear.
func (mr *MockInterfaceMockRecorder) GetGeneralMonthlyInfosFromYear(year interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGeneralMonthlyInfosFromYear", reflect.TypeOf((*MockInterface)(nil).GetGeneralMonthlyInfosFromYear), year)
}

// GetIndexInformation mocks base method.
func (m *MockInterface) GetIndexInformation(name string, month, year int) (map[string][]models.IndexInformation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIndexInformation", name, month, year)
	ret0, _ := ret[0].(map[string][]models.IndexInformation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIndexInformation indicates an expected call of GetIndexInformation.
func (mr *MockInterfaceMockRecorder) GetIndexInformation(name, month, year interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIndexInformation", reflect.TypeOf((*MockInterface)(nil).GetIndexInformation), name, month, year)
}

// GetLastDateWithMonthlyInfo mocks base method.
func (m *MockInterface) GetLastDateWithMonthlyInfo() (int, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastDateWithMonthlyInfo")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetLastDateWithMonthlyInfo indicates an expected call of GetLastDateWithMonthlyInfo.
func (mr *MockInterfaceMockRecorder) GetLastDateWithMonthlyInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastDateWithMonthlyInfo", reflect.TypeOf((*MockInterface)(nil).GetLastDateWithMonthlyInfo))
}

// GetMonthlyInfo mocks base method.
func (m *MockInterface) GetMonthlyInfo(agencies []models.Agency, year int) (map[string][]models.AgencyMonthlyInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMonthlyInfo", agencies, year)
	ret0, _ := ret[0].(map[string][]models.AgencyMonthlyInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMonthlyInfo indicates an expected call of GetMonthlyInfo.
func (mr *MockInterfaceMockRecorder) GetMonthlyInfo(agencies, year interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMonthlyInfo", reflect.TypeOf((*MockInterface)(nil).GetMonthlyInfo), agencies, year)
}

// GetNumberOfMonthsCollected mocks base method.
func (m *MockInterface) GetNumberOfMonthsCollected() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNumberOfMonthsCollected")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNumberOfMonthsCollected indicates an expected call of GetNumberOfMonthsCollected.
func (mr *MockInterfaceMockRecorder) GetNumberOfMonthsCollected() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNumberOfMonthsCollected", reflect.TypeOf((*MockInterface)(nil).GetNumberOfMonthsCollected))
}

// GetOMA mocks base method.
func (m *MockInterface) GetOMA(month, year int, agency string) (*models.AgencyMonthlyInfo, *models.Agency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOMA", month, year, agency)
	ret0, _ := ret[0].(*models.AgencyMonthlyInfo)
	ret1, _ := ret[1].(*models.Agency)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetOMA indicates an expected call of GetOMA.
func (mr *MockInterfaceMockRecorder) GetOMA(month, year, agency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOMA", reflect.TypeOf((*MockInterface)(nil).GetOMA), month, year, agency)
}

// GetOPJ mocks base method.
func (m *MockInterface) GetOPJ(group string) ([]models.Agency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOPJ", group)
	ret0, _ := ret[0].([]models.Agency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOPJ indicates an expected call of GetOPJ.
func (mr *MockInterfaceMockRecorder) GetOPJ(group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOPJ", reflect.TypeOf((*MockInterface)(nil).GetOPJ), group)
}

// GetStateAgencies mocks base method.
func (m *MockInterface) GetStateAgencies(uf string) ([]models.Agency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateAgencies", uf)
	ret0, _ := ret[0].([]models.Agency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateAgencies indicates an expected call of GetStateAgencies.
func (mr *MockInterfaceMockRecorder) GetStateAgencies(uf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateAgencies", reflect.TypeOf((*MockInterface)(nil).GetStateAgencies), uf)
}

// Store mocks base method.
func (m *MockInterface) Store(agmi models.AgencyMonthlyInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", agmi)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockInterfaceMockRecorder) Store(agmi interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockInterface)(nil).Store), agmi)
}

// StorePaychecks mocks base method.
func (m *MockInterface) StorePaychecks(p []models.Paycheck, r []models.Remuneration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorePaychecks", p, r)
	ret0, _ := ret[0].(error)
	return ret0
}

// StorePaychecks indicates an expected call of StorePaychecks.
func (mr *MockInterfaceMockRecorder) StorePaychecks(p, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorePaychecks", reflect.TypeOf((*MockInterface)(nil).StorePaychecks), p, r)
}

// StoreRemunerations mocks base method.
func (m *MockInterface) StoreRemunerations(remu models.Remunerations) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreRemunerations", remu)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreRemunerations indicates an expected call of StoreRemunerations.
func (mr *MockInterfaceMockRecorder) StoreRemunerations(remu interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreRemunerations", reflect.TypeOf((*MockInterface)(nil).StoreRemunerations), remu)
}
