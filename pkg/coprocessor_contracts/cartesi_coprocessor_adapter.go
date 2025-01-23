// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package coprocessor_contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// CoprocessorContractsMetaData contains all meta data concerning the CoprocessorContracts contract.
var CoprocessorContractsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"payloadHash\",\"type\":\"bytes32\"}],\"name\":\"ComputationNotFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"InvalidOutputLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"internalType\":\"bytes4\",\"name\":\"expected\",\"type\":\"bytes4\"}],\"name\":\"InvalidOutputSelector\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"current\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"}],\"name\":\"MachineHashMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"UnauthorizedCaller\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"callCoprocessor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"computationSent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"coprocessor\",\"outputs\":[{\"internalType\":\"contractICoprocessor\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_machineHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_payloadHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"outputs\",\"type\":\"bytes[]\"}],\"name\":\"coprocessorCallbackOutputsOnly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"machineHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// CoprocessorContractsABI is the input ABI used to generate the binding from.
// Deprecated: Use CoprocessorContractsMetaData.ABI instead.
var CoprocessorContractsABI = CoprocessorContractsMetaData.ABI

// CoprocessorContracts is an auto generated Go binding around an Ethereum contract.
type CoprocessorContracts struct {
	CoprocessorContractsCaller     // Read-only binding to the contract
	CoprocessorContractsTransactor // Write-only binding to the contract
	CoprocessorContractsFilterer   // Log filterer for contract events
}

// CoprocessorContractsCaller is an auto generated read-only Go binding around an Ethereum contract.
type CoprocessorContractsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoprocessorContractsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CoprocessorContractsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoprocessorContractsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CoprocessorContractsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoprocessorContractsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CoprocessorContractsSession struct {
	Contract     *CoprocessorContracts // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// CoprocessorContractsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CoprocessorContractsCallerSession struct {
	Contract *CoprocessorContractsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// CoprocessorContractsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CoprocessorContractsTransactorSession struct {
	Contract     *CoprocessorContractsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// CoprocessorContractsRaw is an auto generated low-level Go binding around an Ethereum contract.
type CoprocessorContractsRaw struct {
	Contract *CoprocessorContracts // Generic contract binding to access the raw methods on
}

// CoprocessorContractsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CoprocessorContractsCallerRaw struct {
	Contract *CoprocessorContractsCaller // Generic read-only contract binding to access the raw methods on
}

// CoprocessorContractsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CoprocessorContractsTransactorRaw struct {
	Contract *CoprocessorContractsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCoprocessorContracts creates a new instance of CoprocessorContracts, bound to a specific deployed contract.
func NewCoprocessorContracts(address common.Address, backend bind.ContractBackend) (*CoprocessorContracts, error) {
	contract, err := bindCoprocessorContracts(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CoprocessorContracts{CoprocessorContractsCaller: CoprocessorContractsCaller{contract: contract}, CoprocessorContractsTransactor: CoprocessorContractsTransactor{contract: contract}, CoprocessorContractsFilterer: CoprocessorContractsFilterer{contract: contract}}, nil
}

// NewCoprocessorContractsCaller creates a new read-only instance of CoprocessorContracts, bound to a specific deployed contract.
func NewCoprocessorContractsCaller(address common.Address, caller bind.ContractCaller) (*CoprocessorContractsCaller, error) {
	contract, err := bindCoprocessorContracts(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CoprocessorContractsCaller{contract: contract}, nil
}

// NewCoprocessorContractsTransactor creates a new write-only instance of CoprocessorContracts, bound to a specific deployed contract.
func NewCoprocessorContractsTransactor(address common.Address, transactor bind.ContractTransactor) (*CoprocessorContractsTransactor, error) {
	contract, err := bindCoprocessorContracts(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CoprocessorContractsTransactor{contract: contract}, nil
}

// NewCoprocessorContractsFilterer creates a new log filterer instance of CoprocessorContracts, bound to a specific deployed contract.
func NewCoprocessorContractsFilterer(address common.Address, filterer bind.ContractFilterer) (*CoprocessorContractsFilterer, error) {
	contract, err := bindCoprocessorContracts(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CoprocessorContractsFilterer{contract: contract}, nil
}

// bindCoprocessorContracts binds a generic wrapper to an already deployed contract.
func bindCoprocessorContracts(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CoprocessorContractsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CoprocessorContracts *CoprocessorContractsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CoprocessorContracts.Contract.CoprocessorContractsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CoprocessorContracts *CoprocessorContractsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CoprocessorContracts.Contract.CoprocessorContractsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CoprocessorContracts *CoprocessorContractsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CoprocessorContracts.Contract.CoprocessorContractsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CoprocessorContracts *CoprocessorContractsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CoprocessorContracts.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CoprocessorContracts *CoprocessorContractsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CoprocessorContracts.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CoprocessorContracts *CoprocessorContractsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CoprocessorContracts.Contract.contract.Transact(opts, method, params...)
}

// ComputationSent is a free data retrieval call binding the contract method 0xdca588c1.
//
// Solidity: function computationSent(bytes32 ) view returns(bool)
func (_CoprocessorContracts *CoprocessorContractsCaller) ComputationSent(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _CoprocessorContracts.contract.Call(opts, &out, "computationSent", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ComputationSent is a free data retrieval call binding the contract method 0xdca588c1.
//
// Solidity: function computationSent(bytes32 ) view returns(bool)
func (_CoprocessorContracts *CoprocessorContractsSession) ComputationSent(arg0 [32]byte) (bool, error) {
	return _CoprocessorContracts.Contract.ComputationSent(&_CoprocessorContracts.CallOpts, arg0)
}

// ComputationSent is a free data retrieval call binding the contract method 0xdca588c1.
//
// Solidity: function computationSent(bytes32 ) view returns(bool)
func (_CoprocessorContracts *CoprocessorContractsCallerSession) ComputationSent(arg0 [32]byte) (bool, error) {
	return _CoprocessorContracts.Contract.ComputationSent(&_CoprocessorContracts.CallOpts, arg0)
}

// Coprocessor is a free data retrieval call binding the contract method 0x7382084a.
//
// Solidity: function coprocessor() view returns(address)
func (_CoprocessorContracts *CoprocessorContractsCaller) Coprocessor(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CoprocessorContracts.contract.Call(opts, &out, "coprocessor")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Coprocessor is a free data retrieval call binding the contract method 0x7382084a.
//
// Solidity: function coprocessor() view returns(address)
func (_CoprocessorContracts *CoprocessorContractsSession) Coprocessor() (common.Address, error) {
	return _CoprocessorContracts.Contract.Coprocessor(&_CoprocessorContracts.CallOpts)
}

// Coprocessor is a free data retrieval call binding the contract method 0x7382084a.
//
// Solidity: function coprocessor() view returns(address)
func (_CoprocessorContracts *CoprocessorContractsCallerSession) Coprocessor() (common.Address, error) {
	return _CoprocessorContracts.Contract.Coprocessor(&_CoprocessorContracts.CallOpts)
}

// MachineHash is a free data retrieval call binding the contract method 0x25daa706.
//
// Solidity: function machineHash() view returns(bytes32)
func (_CoprocessorContracts *CoprocessorContractsCaller) MachineHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _CoprocessorContracts.contract.Call(opts, &out, "machineHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MachineHash is a free data retrieval call binding the contract method 0x25daa706.
//
// Solidity: function machineHash() view returns(bytes32)
func (_CoprocessorContracts *CoprocessorContractsSession) MachineHash() ([32]byte, error) {
	return _CoprocessorContracts.Contract.MachineHash(&_CoprocessorContracts.CallOpts)
}

// MachineHash is a free data retrieval call binding the contract method 0x25daa706.
//
// Solidity: function machineHash() view returns(bytes32)
func (_CoprocessorContracts *CoprocessorContractsCallerSession) MachineHash() ([32]byte, error) {
	return _CoprocessorContracts.Contract.MachineHash(&_CoprocessorContracts.CallOpts)
}

// CallCoprocessor is a paid mutator transaction binding the contract method 0x127f92a3.
//
// Solidity: function callCoprocessor(bytes input) returns()
func (_CoprocessorContracts *CoprocessorContractsTransactor) CallCoprocessor(opts *bind.TransactOpts, input []byte) (*types.Transaction, error) {
	return _CoprocessorContracts.contract.Transact(opts, "callCoprocessor", input)
}

// CallCoprocessor is a paid mutator transaction binding the contract method 0x127f92a3.
//
// Solidity: function callCoprocessor(bytes input) returns()
func (_CoprocessorContracts *CoprocessorContractsSession) CallCoprocessor(input []byte) (*types.Transaction, error) {
	return _CoprocessorContracts.Contract.CallCoprocessor(&_CoprocessorContracts.TransactOpts, input)
}

// CallCoprocessor is a paid mutator transaction binding the contract method 0x127f92a3.
//
// Solidity: function callCoprocessor(bytes input) returns()
func (_CoprocessorContracts *CoprocessorContractsTransactorSession) CallCoprocessor(input []byte) (*types.Transaction, error) {
	return _CoprocessorContracts.Contract.CallCoprocessor(&_CoprocessorContracts.TransactOpts, input)
}

// CoprocessorCallbackOutputsOnly is a paid mutator transaction binding the contract method 0x58f6e29f.
//
// Solidity: function coprocessorCallbackOutputsOnly(bytes32 _machineHash, bytes32 _payloadHash, bytes[] outputs) returns()
func (_CoprocessorContracts *CoprocessorContractsTransactor) CoprocessorCallbackOutputsOnly(opts *bind.TransactOpts, _machineHash [32]byte, _payloadHash [32]byte, outputs [][]byte) (*types.Transaction, error) {
	return _CoprocessorContracts.contract.Transact(opts, "coprocessorCallbackOutputsOnly", _machineHash, _payloadHash, outputs)
}

// CoprocessorCallbackOutputsOnly is a paid mutator transaction binding the contract method 0x58f6e29f.
//
// Solidity: function coprocessorCallbackOutputsOnly(bytes32 _machineHash, bytes32 _payloadHash, bytes[] outputs) returns()
func (_CoprocessorContracts *CoprocessorContractsSession) CoprocessorCallbackOutputsOnly(_machineHash [32]byte, _payloadHash [32]byte, outputs [][]byte) (*types.Transaction, error) {
	return _CoprocessorContracts.Contract.CoprocessorCallbackOutputsOnly(&_CoprocessorContracts.TransactOpts, _machineHash, _payloadHash, outputs)
}

// CoprocessorCallbackOutputsOnly is a paid mutator transaction binding the contract method 0x58f6e29f.
//
// Solidity: function coprocessorCallbackOutputsOnly(bytes32 _machineHash, bytes32 _payloadHash, bytes[] outputs) returns()
func (_CoprocessorContracts *CoprocessorContractsTransactorSession) CoprocessorCallbackOutputsOnly(_machineHash [32]byte, _payloadHash [32]byte, outputs [][]byte) (*types.Transaction, error) {
	return _CoprocessorContracts.Contract.CoprocessorCallbackOutputsOnly(&_CoprocessorContracts.TransactOpts, _machineHash, _payloadHash, outputs)
}
