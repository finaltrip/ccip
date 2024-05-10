// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	discoverer "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/discoverer"
	mock "github.com/stretchr/testify/mock"

	models "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

// Factory is an autogenerated mock type for the Factory type
type Factory struct {
	mock.Mock
}

// NewDiscoverer provides a mock function with given fields: selector, rebalancerAddress
func (_m *Factory) NewDiscoverer(selector models.NetworkSelector, rebalancerAddress models.Address) (discoverer.Discoverer, error) {
	ret := _m.Called(selector, rebalancerAddress)

	if len(ret) == 0 {
		panic("no return value specified for NewDiscoverer")
	}

	var r0 discoverer.Discoverer
	var r1 error
	if rf, ok := ret.Get(0).(func(models.NetworkSelector, models.Address) (discoverer.Discoverer, error)); ok {
		return rf(selector, rebalancerAddress)
	}
	if rf, ok := ret.Get(0).(func(models.NetworkSelector, models.Address) discoverer.Discoverer); ok {
		r0 = rf(selector, rebalancerAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(discoverer.Discoverer)
		}
	}

	if rf, ok := ret.Get(1).(func(models.NetworkSelector, models.Address) error); ok {
		r1 = rf(selector, rebalancerAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFactory creates a new instance of Factory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *Factory {
	mock := &Factory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
