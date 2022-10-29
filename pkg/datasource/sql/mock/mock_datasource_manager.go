// Code generated by MockGen. DO NOT EDIT.
// Source: ../datasource/datasource_manager.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	datasource "github.com/seata/seata-go/pkg/datasource/sql/datasource"
	types "github.com/seata/seata-go/pkg/datasource/sql/types"
	branch "github.com/seata/seata-go/pkg/protocol/branch"
	message "github.com/seata/seata-go/pkg/protocol/message"
	rm "github.com/seata/seata-go/pkg/rm"
)

// MockDataSourceManager is a mock of DataSourceManager interface.
type MockDataSourceManager struct {
	ctrl     *gomock.Controller
	recorder *MockDataSourceManagerMockRecorder
}

// MockDataSourceManagerMockRecorder is the mock recorder for MockDataSourceManager.
type MockDataSourceManagerMockRecorder struct {
	mock *MockDataSourceManager
}

// NewMockDataSourceManager creates a new mock instance.
func NewMockDataSourceManager(ctrl *gomock.Controller) *MockDataSourceManager {
	mock := &MockDataSourceManager{ctrl: ctrl}
	mock.recorder = &MockDataSourceManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataSourceManager) EXPECT() *MockDataSourceManagerMockRecorder {
	return m.recorder
}

// BranchCommit mocks base method.
func (m *MockDataSourceManager) BranchCommit(ctx context.Context, req message.BranchCommitRequest) (branch.BranchStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BranchCommit", ctx, req)
	ret0, _ := ret[0].(branch.BranchStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BranchCommit indicates an expected call of BranchCommit.
func (mr *MockDataSourceManagerMockRecorder) BranchCommit(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BranchCommit", reflect.TypeOf((*MockDataSourceManager)(nil).BranchCommit), ctx, req)
}

// BranchRegister mocks base method.
func (m *MockDataSourceManager) BranchRegister(ctx context.Context, clientId string, req message.BranchRegisterRequest) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BranchRegister", ctx, clientId, req)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BranchRegister indicates an expected call of BranchRegister.
func (mr *MockDataSourceManagerMockRecorder) BranchRegister(ctx, clientId, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BranchRegister", reflect.TypeOf((*MockDataSourceManager)(nil).BranchRegister), ctx, clientId, req)
}

// BranchReport mocks base method.
func (m *MockDataSourceManager) BranchReport(ctx context.Context, req message.BranchReportRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BranchReport", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// BranchReport indicates an expected call of BranchReport.
func (mr *MockDataSourceManagerMockRecorder) BranchReport(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BranchReport", reflect.TypeOf((*MockDataSourceManager)(nil).BranchReport), ctx, req)
}

// BranchRollback mocks base method.
func (m *MockDataSourceManager) BranchRollback(ctx context.Context, req message.BranchRollbackRequest) (branch.BranchStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BranchRollback", ctx, req)
	ret0, _ := ret[0].(branch.BranchStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BranchRollback indicates an expected call of BranchRollback.
func (mr *MockDataSourceManagerMockRecorder) BranchRollback(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BranchRollback", reflect.TypeOf((*MockDataSourceManager)(nil).BranchRollback), ctx, req)
}

// CreateTableMetaCache mocks base method.
func (m *MockDataSourceManager) CreateTableMetaCache(ctx context.Context, resID string, dbType types.DBType, db *sql.DB) (datasource.TableMetaCache, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTableMetaCache", ctx, resID, dbType, db)
	ret0, _ := ret[0].(datasource.TableMetaCache)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTableMetaCache indicates an expected call of CreateTableMetaCache.
func (mr *MockDataSourceManagerMockRecorder) CreateTableMetaCache(ctx, resID, dbType, db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTableMetaCache", reflect.TypeOf((*MockDataSourceManager)(nil).CreateTableMetaCache), ctx, resID, dbType, db)
}

// GetManagedResources mocks base method.
func (m *MockDataSourceManager) GetManagedResources() map[string]rm.Resource {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagedResources")
	ret0, _ := ret[0].(map[string]rm.Resource)
	return ret0
}

// GetManagedResources indicates an expected call of GetManagedResources.
func (mr *MockDataSourceManagerMockRecorder) GetManagedResources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagedResources", reflect.TypeOf((*MockDataSourceManager)(nil).GetManagedResources))
}

// LockQuery mocks base method.
func (m *MockDataSourceManager) LockQuery(ctx context.Context, req message.GlobalLockQueryRequest) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LockQuery", ctx, req)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LockQuery indicates an expected call of LockQuery.
func (mr *MockDataSourceManagerMockRecorder) LockQuery(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockQuery", reflect.TypeOf((*MockDataSourceManager)(nil).LockQuery), ctx, req)
}

// RegisterResource mocks base method.
func (m *MockDataSourceManager) RegisterResource(resource rm.Resource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterResource", resource)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterResource indicates an expected call of RegisterResource.
func (mr *MockDataSourceManagerMockRecorder) RegisterResource(resource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterResource", reflect.TypeOf((*MockDataSourceManager)(nil).RegisterResource), resource)
}

// UnregisterResource mocks base method.
func (m *MockDataSourceManager) UnregisterResource(resource rm.Resource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterResource", resource)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterResource indicates an expected call of UnregisterResource.
func (mr *MockDataSourceManagerMockRecorder) UnregisterResource(resource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterResource", reflect.TypeOf((*MockDataSourceManager)(nil).UnregisterResource), resource)
}

// MockTableMetaCache is a mock of TableMetaCache interface.
type MockTableMetaCache struct {
	ctrl     *gomock.Controller
	recorder *MockTableMetaCacheMockRecorder
}

// MockTableMetaCacheMockRecorder is the mock recorder for MockTableMetaCache.
type MockTableMetaCacheMockRecorder struct {
	mock *MockTableMetaCache
}

// NewMockTableMetaCache creates a new mock instance.
func NewMockTableMetaCache(ctrl *gomock.Controller) *MockTableMetaCache {
	mock := &MockTableMetaCache{ctrl: ctrl}
	mock.recorder = &MockTableMetaCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTableMetaCache) EXPECT() *MockTableMetaCacheMockRecorder {
	return m.recorder
}

// Destroy mocks base method.
func (m *MockTableMetaCache) Destroy() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy")
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy.
func (mr *MockTableMetaCacheMockRecorder) Destroy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockTableMetaCache)(nil).Destroy))
}

// GetTableMeta mocks base method.
func (m *MockTableMetaCache) GetTableMeta(table string) (types.TableMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTableMeta", table)
	ret0, _ := ret[0].(types.TableMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTableMeta indicates an expected call of GetTableMeta.
func (mr *MockTableMetaCacheMockRecorder) GetTableMeta(table interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTableMeta", reflect.TypeOf((*MockTableMetaCache)(nil).GetTableMeta), table)
}

// Init mocks base method.
func (m *MockTableMetaCache) Init(ctx context.Context, conn *sql.DB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", ctx, conn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockTableMetaCacheMockRecorder) Init(ctx, conn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockTableMetaCache)(nil).Init), ctx, conn)
}
