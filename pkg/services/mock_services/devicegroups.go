// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/services/devicegroups.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/redhatinsights/edge-api/pkg/models"
	gorm "gorm.io/gorm"
	reflect "reflect"
)

// MockDeviceGroupsServiceInterface is a mock of DeviceGroupsServiceInterface interface
type MockDeviceGroupsServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceGroupsServiceInterfaceMockRecorder
}

// MockDeviceGroupsServiceInterfaceMockRecorder is the mock recorder for MockDeviceGroupsServiceInterface
type MockDeviceGroupsServiceInterfaceMockRecorder struct {
	mock *MockDeviceGroupsServiceInterface
}

// NewMockDeviceGroupsServiceInterface creates a new mock instance
func NewMockDeviceGroupsServiceInterface(ctrl *gomock.Controller) *MockDeviceGroupsServiceInterface {
	mock := &MockDeviceGroupsServiceInterface{ctrl: ctrl}
	mock.recorder = &MockDeviceGroupsServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDeviceGroupsServiceInterface) EXPECT() *MockDeviceGroupsServiceInterfaceMockRecorder {
	return m.recorder
}

// CreateDeviceGroup mocks base method
func (m *MockDeviceGroupsServiceInterface) CreateDeviceGroup(deviceGroup *models.DeviceGroup) (*models.DeviceGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDeviceGroup", deviceGroup)
	ret0, _ := ret[0].(*models.DeviceGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDeviceGroup indicates an expected call of CreateDeviceGroup
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) CreateDeviceGroup(deviceGroup interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDeviceGroup", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).CreateDeviceGroup), deviceGroup)
}

// GetDeviceGroups mocks base method
func (m *MockDeviceGroupsServiceInterface) GetDeviceGroups(account string, limit, offset int, tx *gorm.DB) (*[]models.DeviceGroupListDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceGroups", account, limit, offset, tx)
	ret0, _ := ret[0].(*[]models.DeviceGroupListDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceGroups indicates an expected call of GetDeviceGroups
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) GetDeviceGroups(account, limit, offset, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceGroups", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).GetDeviceGroups), account, limit, offset, tx)
}

// GetDeviceGroupsCount mocks base method
func (m *MockDeviceGroupsServiceInterface) GetDeviceGroupsCount(account string, tx *gorm.DB) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceGroupsCount", account, tx)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceGroupsCount indicates an expected call of GetDeviceGroupsCount
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) GetDeviceGroupsCount(account, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceGroupsCount", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).GetDeviceGroupsCount), account, tx)
}

// GetDeviceGroupByID mocks base method
func (m *MockDeviceGroupsServiceInterface) GetDeviceGroupByID(ID string) (*models.DeviceGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceGroupByID", ID)
	ret0, _ := ret[0].(*models.DeviceGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceGroupByID indicates an expected call of GetDeviceGroupByID
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) GetDeviceGroupByID(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceGroupByID", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).GetDeviceGroupByID), ID)
}

// GetDeviceGroupDetailsByID mocks base method
func (m *MockDeviceGroupsServiceInterface) GetDeviceGroupDetailsByID(ID string) (*models.DeviceGroupDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceGroupDetailsByID", ID)
	ret0, _ := ret[0].(*models.DeviceGroupDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceGroupDetailsByID indicates an expected call of GetDeviceGroupDetailsByID
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) GetDeviceGroupDetailsByID(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceGroupDetailsByID", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).GetDeviceGroupDetailsByID), ID)
}

// DeleteDeviceGroupByID mocks base method
func (m *MockDeviceGroupsServiceInterface) DeleteDeviceGroupByID(ID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDeviceGroupByID", ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDeviceGroupByID indicates an expected call of DeleteDeviceGroupByID
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) DeleteDeviceGroupByID(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDeviceGroupByID", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).DeleteDeviceGroupByID), ID)
}

// UpdateDeviceGroup mocks base method
func (m *MockDeviceGroupsServiceInterface) UpdateDeviceGroup(deviceGroup *models.DeviceGroup, account, ID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDeviceGroup", deviceGroup, account, ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDeviceGroup indicates an expected call of UpdateDeviceGroup
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) UpdateDeviceGroup(deviceGroup, account, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDeviceGroup", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).UpdateDeviceGroup), deviceGroup, account, ID)
}

// GetDeviceGroupDeviceByID mocks base method
func (m *MockDeviceGroupsServiceInterface) GetDeviceGroupDeviceByID(account string, deviceGroupID, deviceID uint) (*models.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceGroupDeviceByID", account, deviceGroupID, deviceID)
	ret0, _ := ret[0].(*models.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceGroupDeviceByID indicates an expected call of GetDeviceGroupDeviceByID
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) GetDeviceGroupDeviceByID(account, deviceGroupID, deviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceGroupDeviceByID", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).GetDeviceGroupDeviceByID), account, deviceGroupID, deviceID)
}

// AddDeviceGroupDevices mocks base method
func (m *MockDeviceGroupsServiceInterface) AddDeviceGroupDevices(account string, deviceGroupID uint, devices []models.Device) (*[]models.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDeviceGroupDevices", account, deviceGroupID, devices)
	ret0, _ := ret[0].(*[]models.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDeviceGroupDevices indicates an expected call of AddDeviceGroupDevices
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) AddDeviceGroupDevices(account, deviceGroupID, devices interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDeviceGroupDevices", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).AddDeviceGroupDevices), account, deviceGroupID, devices)
}

// DeleteDeviceGroupDevices mocks base method
func (m *MockDeviceGroupsServiceInterface) DeleteDeviceGroupDevices(account string, deviceGroupID uint, devices []models.Device) (*[]models.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDeviceGroupDevices", account, deviceGroupID, devices)
	ret0, _ := ret[0].(*[]models.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteDeviceGroupDevices indicates an expected call of DeleteDeviceGroupDevices
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) DeleteDeviceGroupDevices(account, deviceGroupID, devices interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDeviceGroupDevices", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).DeleteDeviceGroupDevices), account, deviceGroupID, devices)
}

// GetDeviceImageInfo mocks base method
func (m *MockDeviceGroupsServiceInterface) GetDeviceImageInfo(setOfImages map[int]models.DeviceImageInfo, account string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceImageInfo", setOfImages, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetDeviceImageInfo indicates an expected call of GetDeviceImageInfo
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) GetDeviceImageInfo(setOfImages, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceImageInfo", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).GetDeviceImageInfo), setOfImages, account)
}

// DeviceGroupNameExists mocks base method
func (m *MockDeviceGroupsServiceInterface) DeviceGroupNameExists(account, name string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeviceGroupNameExists", account, name)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeviceGroupNameExists indicates an expected call of DeviceGroupNameExists
func (mr *MockDeviceGroupsServiceInterfaceMockRecorder) DeviceGroupNameExists(account, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeviceGroupNameExists", reflect.TypeOf((*MockDeviceGroupsServiceInterface)(nil).DeviceGroupNameExists), account, name)
}
