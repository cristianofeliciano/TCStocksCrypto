// Code generated by MockGen. DO NOT EDIT.
// Source: ./app/crypto/crypto.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/tradersclub/TCStocksCrypto/model"
)

// MockApp is a mock of App interface.
type MockApp struct {
	ctrl     *gomock.Controller
	recorder *MockAppMockRecorder
}

// MockAppMockRecorder is the mock recorder for MockApp.
type MockAppMockRecorder struct {
	mock *MockApp
}

// NewMockApp creates a new mock instance.
func NewMockApp(ctrl *gomock.Controller) *MockApp {
	mock := &MockApp{ctrl: ctrl}
	mock.recorder = &MockAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApp) EXPECT() *MockAppMockRecorder {
	return m.recorder
}

// GetCryptoCategories mocks base method.
func (m *MockApp) GetCryptoCategories(ctx context.Context) ([]model.CryptoCategories, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCryptoCategories", ctx)
	ret0, _ := ret[0].([]model.CryptoCategories)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCryptoCategories indicates an expected call of GetCryptoCategories.
func (mr *MockAppMockRecorder) GetCryptoCategories(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCryptoCategories", reflect.TypeOf((*MockApp)(nil).GetCryptoCategories), ctx)
}

// GetCryptoMarkets mocks base method.
func (m *MockApp) GetCryptoMarkets(ctx context.Context) ([]model.Market, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCryptoMarkets", ctx)
	ret0, _ := ret[0].([]model.Market)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCryptoMarkets indicates an expected call of GetCryptoMarkets.
func (mr *MockAppMockRecorder) GetCryptoMarkets(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCryptoMarkets", reflect.TypeOf((*MockApp)(nil).GetCryptoMarkets), ctx)
}

// GetGlobalInfos mocks base method.
func (m *MockApp) GetGlobalInfos(ctx context.Context) (*model.GlobalInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGlobalInfos", ctx)
	ret0, _ := ret[0].(*model.GlobalInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGlobalInfos indicates an expected call of GetGlobalInfos.
func (mr *MockAppMockRecorder) GetGlobalInfos(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGlobalInfos", reflect.TypeOf((*MockApp)(nil).GetGlobalInfos), ctx)
}
