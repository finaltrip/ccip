// Code generated by mockery v2.42.2. DO NOT EDIT.

package tokendata

import (
	context "context"

	ccip "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	mock "github.com/stretchr/testify/mock"
)

// MockReader is an autogenerated mock type for the Reader type
type MockReader struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MockReader) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadTokenData provides a mock function with given fields: ctx, msg, tokenIndex
func (_m *MockReader) ReadTokenData(ctx context.Context, msg ccip.EVM2EVMOnRampCCIPSendRequestedWithMeta, tokenIndex int) ([]byte, error) {
	ret := _m.Called(ctx, msg, tokenIndex)

	if len(ret) == 0 {
		panic("no return value specified for ReadTokenData")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ccip.EVM2EVMOnRampCCIPSendRequestedWithMeta, int) ([]byte, error)); ok {
		return rf(ctx, msg, tokenIndex)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ccip.EVM2EVMOnRampCCIPSendRequestedWithMeta, int) []byte); ok {
		r0 = rf(ctx, msg, tokenIndex)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ccip.EVM2EVMOnRampCCIPSendRequestedWithMeta, int) error); ok {
		r1 = rf(ctx, msg, tokenIndex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockReader creates a new instance of MockReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockReader {
	mock := &MockReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
