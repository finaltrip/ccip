// Code generated by mockery v2.42.2. DO NOT EDIT.

package mock_contracts

import (
	big "math/big"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	common "github.com/ethereum/go-ethereum/common"

	event "github.com/ethereum/go-ethereum/event"

	generated "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"

	link_token_interface "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// LinkTokenInterface is an autogenerated mock type for the LinkTokenInterface type
type LinkTokenInterface struct {
	mock.Mock
}

// Address provides a mock function with given fields:
func (_m *LinkTokenInterface) Address() common.Address {
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

// Allowance provides a mock function with given fields: opts, _owner, _spender
func (_m *LinkTokenInterface) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	ret := _m.Called(opts, _owner, _spender)

	if len(ret) == 0 {
		panic("no return value specified for Allowance")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, common.Address, common.Address) (*big.Int, error)); ok {
		return rf(opts, _owner, _spender)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, common.Address, common.Address) *big.Int); ok {
		r0 = rf(opts, _owner, _spender)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts, common.Address, common.Address) error); ok {
		r1 = rf(opts, _owner, _spender)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Approve provides a mock function with given fields: opts, _spender, _value
func (_m *LinkTokenInterface) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	ret := _m.Called(opts, _spender, _value)

	if len(ret) == 0 {
		panic("no return value specified for Approve")
	}

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, *big.Int) (*types.Transaction, error)); ok {
		return rf(opts, _spender, _value)
	}
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, *big.Int) *types.Transaction); ok {
		r0 = rf(opts, _spender, _value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.TransactOpts, common.Address, *big.Int) error); ok {
		r1 = rf(opts, _spender, _value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BalanceOf provides a mock function with given fields: opts, _owner
func (_m *LinkTokenInterface) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	ret := _m.Called(opts, _owner)

	if len(ret) == 0 {
		panic("no return value specified for BalanceOf")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, common.Address) (*big.Int, error)); ok {
		return rf(opts, _owner)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, common.Address) *big.Int); ok {
		r0 = rf(opts, _owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts, common.Address) error); ok {
		r1 = rf(opts, _owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Decimals provides a mock function with given fields: opts
func (_m *LinkTokenInterface) Decimals(opts *bind.CallOpts) (uint8, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for Decimals")
	}

	var r0 uint8
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (uint8, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) uint8); ok {
		r0 = rf(opts)
	} else {
		r0 = ret.Get(0).(uint8)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecreaseApproval provides a mock function with given fields: opts, _spender, _subtractedValue
func (_m *LinkTokenInterface) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	ret := _m.Called(opts, _spender, _subtractedValue)

	if len(ret) == 0 {
		panic("no return value specified for DecreaseApproval")
	}

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, *big.Int) (*types.Transaction, error)); ok {
		return rf(opts, _spender, _subtractedValue)
	}
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, *big.Int) *types.Transaction); ok {
		r0 = rf(opts, _spender, _subtractedValue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.TransactOpts, common.Address, *big.Int) error); ok {
		r1 = rf(opts, _spender, _subtractedValue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FilterApproval provides a mock function with given fields: opts, owner, spender
func (_m *LinkTokenInterface) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*link_token_interface.LinkTokenApprovalIterator, error) {
	ret := _m.Called(opts, owner, spender)

	if len(ret) == 0 {
		panic("no return value specified for FilterApproval")
	}

	var r0 *link_token_interface.LinkTokenApprovalIterator
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.FilterOpts, []common.Address, []common.Address) (*link_token_interface.LinkTokenApprovalIterator, error)); ok {
		return rf(opts, owner, spender)
	}
	if rf, ok := ret.Get(0).(func(*bind.FilterOpts, []common.Address, []common.Address) *link_token_interface.LinkTokenApprovalIterator); ok {
		r0 = rf(opts, owner, spender)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*link_token_interface.LinkTokenApprovalIterator)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.FilterOpts, []common.Address, []common.Address) error); ok {
		r1 = rf(opts, owner, spender)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FilterTransfer provides a mock function with given fields: opts, from, to
func (_m *LinkTokenInterface) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*link_token_interface.LinkTokenTransferIterator, error) {
	ret := _m.Called(opts, from, to)

	if len(ret) == 0 {
		panic("no return value specified for FilterTransfer")
	}

	var r0 *link_token_interface.LinkTokenTransferIterator
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.FilterOpts, []common.Address, []common.Address) (*link_token_interface.LinkTokenTransferIterator, error)); ok {
		return rf(opts, from, to)
	}
	if rf, ok := ret.Get(0).(func(*bind.FilterOpts, []common.Address, []common.Address) *link_token_interface.LinkTokenTransferIterator); ok {
		r0 = rf(opts, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*link_token_interface.LinkTokenTransferIterator)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.FilterOpts, []common.Address, []common.Address) error); ok {
		r1 = rf(opts, from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IncreaseApproval provides a mock function with given fields: opts, _spender, _addedValue
func (_m *LinkTokenInterface) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	ret := _m.Called(opts, _spender, _addedValue)

	if len(ret) == 0 {
		panic("no return value specified for IncreaseApproval")
	}

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, *big.Int) (*types.Transaction, error)); ok {
		return rf(opts, _spender, _addedValue)
	}
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, *big.Int) *types.Transaction); ok {
		r0 = rf(opts, _spender, _addedValue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.TransactOpts, common.Address, *big.Int) error); ok {
		r1 = rf(opts, _spender, _addedValue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Name provides a mock function with given fields: opts
func (_m *LinkTokenInterface) Name(opts *bind.CallOpts) (string, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (string, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) string); ok {
		r0 = rf(opts)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseApproval provides a mock function with given fields: log
func (_m *LinkTokenInterface) ParseApproval(log types.Log) (*link_token_interface.LinkTokenApproval, error) {
	ret := _m.Called(log)

	if len(ret) == 0 {
		panic("no return value specified for ParseApproval")
	}

	var r0 *link_token_interface.LinkTokenApproval
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Log) (*link_token_interface.LinkTokenApproval, error)); ok {
		return rf(log)
	}
	if rf, ok := ret.Get(0).(func(types.Log) *link_token_interface.LinkTokenApproval); ok {
		r0 = rf(log)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*link_token_interface.LinkTokenApproval)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Log) error); ok {
		r1 = rf(log)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseLog provides a mock function with given fields: log
func (_m *LinkTokenInterface) ParseLog(log types.Log) (generated.AbigenLog, error) {
	ret := _m.Called(log)

	if len(ret) == 0 {
		panic("no return value specified for ParseLog")
	}

	var r0 generated.AbigenLog
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Log) (generated.AbigenLog, error)); ok {
		return rf(log)
	}
	if rf, ok := ret.Get(0).(func(types.Log) generated.AbigenLog); ok {
		r0 = rf(log)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(generated.AbigenLog)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Log) error); ok {
		r1 = rf(log)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseTransfer provides a mock function with given fields: log
func (_m *LinkTokenInterface) ParseTransfer(log types.Log) (*link_token_interface.LinkTokenTransfer, error) {
	ret := _m.Called(log)

	if len(ret) == 0 {
		panic("no return value specified for ParseTransfer")
	}

	var r0 *link_token_interface.LinkTokenTransfer
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Log) (*link_token_interface.LinkTokenTransfer, error)); ok {
		return rf(log)
	}
	if rf, ok := ret.Get(0).(func(types.Log) *link_token_interface.LinkTokenTransfer); ok {
		r0 = rf(log)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*link_token_interface.LinkTokenTransfer)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Log) error); ok {
		r1 = rf(log)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Symbol provides a mock function with given fields: opts
func (_m *LinkTokenInterface) Symbol(opts *bind.CallOpts) (string, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for Symbol")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (string, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) string); ok {
		r0 = rf(opts)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TotalSupply provides a mock function with given fields: opts
func (_m *LinkTokenInterface) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for TotalSupply")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (*big.Int, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) *big.Int); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Transfer provides a mock function with given fields: opts, _to, _value
func (_m *LinkTokenInterface) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	ret := _m.Called(opts, _to, _value)

	if len(ret) == 0 {
		panic("no return value specified for Transfer")
	}

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, *big.Int) (*types.Transaction, error)); ok {
		return rf(opts, _to, _value)
	}
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, *big.Int) *types.Transaction); ok {
		r0 = rf(opts, _to, _value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.TransactOpts, common.Address, *big.Int) error); ok {
		r1 = rf(opts, _to, _value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransferAndCall provides a mock function with given fields: opts, _to, _value, _data
func (_m *LinkTokenInterface) TransferAndCall(opts *bind.TransactOpts, _to common.Address, _value *big.Int, _data []byte) (*types.Transaction, error) {
	ret := _m.Called(opts, _to, _value, _data)

	if len(ret) == 0 {
		panic("no return value specified for TransferAndCall")
	}

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, *big.Int, []byte) (*types.Transaction, error)); ok {
		return rf(opts, _to, _value, _data)
	}
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, *big.Int, []byte) *types.Transaction); ok {
		r0 = rf(opts, _to, _value, _data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.TransactOpts, common.Address, *big.Int, []byte) error); ok {
		r1 = rf(opts, _to, _value, _data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransferFrom provides a mock function with given fields: opts, _from, _to, _value
func (_m *LinkTokenInterface) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	ret := _m.Called(opts, _from, _to, _value)

	if len(ret) == 0 {
		panic("no return value specified for TransferFrom")
	}

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, common.Address, *big.Int) (*types.Transaction, error)); ok {
		return rf(opts, _from, _to, _value)
	}
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, common.Address, common.Address, *big.Int) *types.Transaction); ok {
		r0 = rf(opts, _from, _to, _value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.TransactOpts, common.Address, common.Address, *big.Int) error); ok {
		r1 = rf(opts, _from, _to, _value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WatchApproval provides a mock function with given fields: opts, sink, owner, spender
func (_m *LinkTokenInterface) WatchApproval(opts *bind.WatchOpts, sink chan<- *link_token_interface.LinkTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
	ret := _m.Called(opts, sink, owner, spender)

	if len(ret) == 0 {
		panic("no return value specified for WatchApproval")
	}

	var r0 event.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *link_token_interface.LinkTokenApproval, []common.Address, []common.Address) (event.Subscription, error)); ok {
		return rf(opts, sink, owner, spender)
	}
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *link_token_interface.LinkTokenApproval, []common.Address, []common.Address) event.Subscription); ok {
		r0 = rf(opts, sink, owner, spender)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(event.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.WatchOpts, chan<- *link_token_interface.LinkTokenApproval, []common.Address, []common.Address) error); ok {
		r1 = rf(opts, sink, owner, spender)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WatchTransfer provides a mock function with given fields: opts, sink, from, to
func (_m *LinkTokenInterface) WatchTransfer(opts *bind.WatchOpts, sink chan<- *link_token_interface.LinkTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
	ret := _m.Called(opts, sink, from, to)

	if len(ret) == 0 {
		panic("no return value specified for WatchTransfer")
	}

	var r0 event.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *link_token_interface.LinkTokenTransfer, []common.Address, []common.Address) (event.Subscription, error)); ok {
		return rf(opts, sink, from, to)
	}
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *link_token_interface.LinkTokenTransfer, []common.Address, []common.Address) event.Subscription); ok {
		r0 = rf(opts, sink, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(event.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.WatchOpts, chan<- *link_token_interface.LinkTokenTransfer, []common.Address, []common.Address) error); ok {
		r1 = rf(opts, sink, from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewLinkTokenInterface creates a new instance of LinkTokenInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLinkTokenInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *LinkTokenInterface {
	mock := &LinkTokenInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
