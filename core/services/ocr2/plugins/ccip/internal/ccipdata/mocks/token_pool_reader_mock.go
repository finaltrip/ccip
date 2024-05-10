// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	common "github.com/ethereum/go-ethereum/common"
	mock "github.com/stretchr/testify/mock"
)

// TokenPoolReader is an autogenerated mock type for the TokenPoolReader type
type TokenPoolReader struct {
	mock.Mock
}

// Address provides a mock function with given fields:
func (_m *TokenPoolReader) Address() common.Address {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Address")
	}

	var r0 common.Address
	if rf, ok := ret.Get(0).(func() common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	return r0
}

// Type provides a mock function with given fields:
func (_m *TokenPoolReader) Type() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Type")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewTokenPoolReader creates a new instance of TokenPoolReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenPoolReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenPoolReader {
	mock := &TokenPoolReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
