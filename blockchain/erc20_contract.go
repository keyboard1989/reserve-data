// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Erc20ContractABI is the input ABI used to generate the binding from.
const Erc20ContractABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"saleStartTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenSaleContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emergencyERC20Drain\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"saleEndTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"name\":\"tokenTotalAmount\",\"type\":\"uint256\"},{\"name\":\"startTime\",\"type\":\"uint256\"},{\"name\":\"endTime\",\"type\":\"uint256\"},{\"name\":\"admin\",\"type\":\"address\"}],\"payable\":false,\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_burner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// Erc20Contract is an auto generated Go binding around an Ethereum contract.
type Erc20Contract struct {
	Erc20ContractCaller     // Read-only binding to the contract
	Erc20ContractTransactor // Write-only binding to the contract
}

// Erc20ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type Erc20ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Erc20ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Erc20ContractSession struct {
	Contract     *Erc20Contract    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc20ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Erc20ContractCallerSession struct {
	Contract *Erc20ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// Erc20ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Erc20ContractTransactorSession struct {
	Contract     *Erc20ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// Erc20ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type Erc20ContractRaw struct {
	Contract *Erc20Contract // Generic contract binding to access the raw methods on
}

// Erc20ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Erc20ContractCallerRaw struct {
	Contract *Erc20ContractCaller // Generic read-only contract binding to access the raw methods on
}

// Erc20ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Erc20ContractTransactorRaw struct {
	Contract *Erc20ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewErc20Contract creates a new instance of Erc20Contract, bound to a specific deployed contract.
func NewErc20Contract(address common.Address, backend bind.ContractBackend) (*Erc20Contract, error) {
	contract, err := bindErc20Contract(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc20Contract{Erc20ContractCaller: Erc20ContractCaller{contract: contract}, Erc20ContractTransactor: Erc20ContractTransactor{contract: contract}}, nil
}

// NewErc20ContractCaller creates a new read-only instance of Erc20Contract, bound to a specific deployed contract.
func NewErc20ContractCaller(address common.Address, caller bind.ContractCaller) (*Erc20ContractCaller, error) {
	contract, err := bindErc20Contract(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20ContractCaller{contract: contract}, nil
}

// NewErc20ContractTransactor creates a new write-only instance of Erc20Contract, bound to a specific deployed contract.
func NewErc20ContractTransactor(address common.Address, transactor bind.ContractTransactor) (*Erc20ContractTransactor, error) {
	contract, err := bindErc20Contract(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &Erc20ContractTransactor{contract: contract}, nil
}

// bindErc20Contract binds a generic wrapper to an already deployed contract.
func bindErc20Contract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc20ContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20Contract *Erc20ContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Erc20Contract.Contract.Erc20ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20Contract *Erc20ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20Contract.Contract.Erc20ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20Contract *Erc20ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20Contract.Contract.Erc20ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20Contract *Erc20ContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Erc20Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20Contract *Erc20ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20Contract *Erc20ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20Contract.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_Erc20Contract *Erc20ContractCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Erc20Contract.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_Erc20Contract *Erc20ContractSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Erc20Contract.Contract.Allowance(&_Erc20Contract.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_Erc20Contract *Erc20ContractCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Erc20Contract.Contract.Allowance(&_Erc20Contract.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_Erc20Contract *Erc20ContractCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Erc20Contract.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_Erc20Contract *Erc20ContractSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Erc20Contract.Contract.BalanceOf(&_Erc20Contract.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_Erc20Contract *Erc20ContractCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Erc20Contract.Contract.BalanceOf(&_Erc20Contract.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint256)
func (_Erc20Contract *Erc20ContractCaller) Decimals(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Erc20Contract.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint256)
func (_Erc20Contract *Erc20ContractSession) Decimals() (*big.Int, error) {
	return _Erc20Contract.Contract.Decimals(&_Erc20Contract.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint256)
func (_Erc20Contract *Erc20ContractCallerSession) Decimals() (*big.Int, error) {
	return _Erc20Contract.Contract.Decimals(&_Erc20Contract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Erc20Contract *Erc20ContractCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Erc20Contract.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Erc20Contract *Erc20ContractSession) Name() (string, error) {
	return _Erc20Contract.Contract.Name(&_Erc20Contract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Erc20Contract *Erc20ContractCallerSession) Name() (string, error) {
	return _Erc20Contract.Contract.Name(&_Erc20Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Erc20Contract *Erc20ContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Erc20Contract.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Erc20Contract *Erc20ContractSession) Owner() (common.Address, error) {
	return _Erc20Contract.Contract.Owner(&_Erc20Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Erc20Contract *Erc20ContractCallerSession) Owner() (common.Address, error) {
	return _Erc20Contract.Contract.Owner(&_Erc20Contract.CallOpts)
}

// SaleEndTime is a free data retrieval call binding the contract method 0xed338ff1.
//
// Solidity: function saleEndTime() constant returns(uint256)
func (_Erc20Contract *Erc20ContractCaller) SaleEndTime(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Erc20Contract.contract.Call(opts, out, "saleEndTime")
	return *ret0, err
}

// SaleEndTime is a free data retrieval call binding the contract method 0xed338ff1.
//
// Solidity: function saleEndTime() constant returns(uint256)
func (_Erc20Contract *Erc20ContractSession) SaleEndTime() (*big.Int, error) {
	return _Erc20Contract.Contract.SaleEndTime(&_Erc20Contract.CallOpts)
}

// SaleEndTime is a free data retrieval call binding the contract method 0xed338ff1.
//
// Solidity: function saleEndTime() constant returns(uint256)
func (_Erc20Contract *Erc20ContractCallerSession) SaleEndTime() (*big.Int, error) {
	return _Erc20Contract.Contract.SaleEndTime(&_Erc20Contract.CallOpts)
}

// SaleStartTime is a free data retrieval call binding the contract method 0x1cbaee2d.
//
// Solidity: function saleStartTime() constant returns(uint256)
func (_Erc20Contract *Erc20ContractCaller) SaleStartTime(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Erc20Contract.contract.Call(opts, out, "saleStartTime")
	return *ret0, err
}

// SaleStartTime is a free data retrieval call binding the contract method 0x1cbaee2d.
//
// Solidity: function saleStartTime() constant returns(uint256)
func (_Erc20Contract *Erc20ContractSession) SaleStartTime() (*big.Int, error) {
	return _Erc20Contract.Contract.SaleStartTime(&_Erc20Contract.CallOpts)
}

// SaleStartTime is a free data retrieval call binding the contract method 0x1cbaee2d.
//
// Solidity: function saleStartTime() constant returns(uint256)
func (_Erc20Contract *Erc20ContractCallerSession) SaleStartTime() (*big.Int, error) {
	return _Erc20Contract.Contract.SaleStartTime(&_Erc20Contract.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Erc20Contract *Erc20ContractCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Erc20Contract.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Erc20Contract *Erc20ContractSession) Symbol() (string, error) {
	return _Erc20Contract.Contract.Symbol(&_Erc20Contract.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Erc20Contract *Erc20ContractCallerSession) Symbol() (string, error) {
	return _Erc20Contract.Contract.Symbol(&_Erc20Contract.CallOpts)
}

// TokenSaleContract is a free data retrieval call binding the contract method 0x5d5aa277.
//
// Solidity: function tokenSaleContract() constant returns(address)
func (_Erc20Contract *Erc20ContractCaller) TokenSaleContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Erc20Contract.contract.Call(opts, out, "tokenSaleContract")
	return *ret0, err
}

// TokenSaleContract is a free data retrieval call binding the contract method 0x5d5aa277.
//
// Solidity: function tokenSaleContract() constant returns(address)
func (_Erc20Contract *Erc20ContractSession) TokenSaleContract() (common.Address, error) {
	return _Erc20Contract.Contract.TokenSaleContract(&_Erc20Contract.CallOpts)
}

// TokenSaleContract is a free data retrieval call binding the contract method 0x5d5aa277.
//
// Solidity: function tokenSaleContract() constant returns(address)
func (_Erc20Contract *Erc20ContractCallerSession) TokenSaleContract() (common.Address, error) {
	return _Erc20Contract.Contract.TokenSaleContract(&_Erc20Contract.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Erc20Contract *Erc20ContractCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Erc20Contract.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Erc20Contract *Erc20ContractSession) TotalSupply() (*big.Int, error) {
	return _Erc20Contract.Contract.TotalSupply(&_Erc20Contract.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Erc20Contract *Erc20ContractCallerSession) TotalSupply() (*big.Int, error) {
	return _Erc20Contract.Contract.TotalSupply(&_Erc20Contract.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.Approve(&_Erc20Contract.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.Approve(&_Erc20Contract.TransactOpts, _spender, _value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(_value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractTransactor) Burn(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.contract.Transact(opts, "burn", _value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(_value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractSession) Burn(_value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.Burn(&_Erc20Contract.TransactOpts, _value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(_value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractTransactorSession) Burn(_value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.Burn(&_Erc20Contract.TransactOpts, _value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(_from address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractTransactor) BurnFrom(opts *bind.TransactOpts, _from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.contract.Transact(opts, "burnFrom", _from, _value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(_from address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractSession) BurnFrom(_from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.BurnFrom(&_Erc20Contract.TransactOpts, _from, _value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(_from address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractTransactorSession) BurnFrom(_from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.BurnFrom(&_Erc20Contract.TransactOpts, _from, _value)
}

// EmergencyERC20Drain is a paid mutator transaction binding the contract method 0xdb0e16f1.
//
// Solidity: function emergencyERC20Drain(token address, amount uint256) returns()
func (_Erc20Contract *Erc20ContractTransactor) EmergencyERC20Drain(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.contract.Transact(opts, "emergencyERC20Drain", token, amount)
}

// EmergencyERC20Drain is a paid mutator transaction binding the contract method 0xdb0e16f1.
//
// Solidity: function emergencyERC20Drain(token address, amount uint256) returns()
func (_Erc20Contract *Erc20ContractSession) EmergencyERC20Drain(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.EmergencyERC20Drain(&_Erc20Contract.TransactOpts, token, amount)
}

// EmergencyERC20Drain is a paid mutator transaction binding the contract method 0xdb0e16f1.
//
// Solidity: function emergencyERC20Drain(token address, amount uint256) returns()
func (_Erc20Contract *Erc20ContractTransactorSession) EmergencyERC20Drain(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.EmergencyERC20Drain(&_Erc20Contract.TransactOpts, token, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.Transfer(&_Erc20Contract.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.Transfer(&_Erc20Contract.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.TransferFrom(&_Erc20Contract.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_Erc20Contract *Erc20ContractTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Erc20Contract.Contract.TransferFrom(&_Erc20Contract.TransactOpts, _from, _to, _value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Erc20Contract *Erc20ContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Erc20Contract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Erc20Contract *Erc20ContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Erc20Contract.Contract.TransferOwnership(&_Erc20Contract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Erc20Contract *Erc20ContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Erc20Contract.Contract.TransferOwnership(&_Erc20Contract.TransactOpts, newOwner)
}
