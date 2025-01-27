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

// MockCoprocessorMetaData contains all meta data concerning the MockCoprocessor contract.
var MockCoprocessorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"machineHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"callback\",\"type\":\"address\"}],\"name\":\"TaskIssued\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"machineHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"callback\",\"type\":\"address\"}],\"name\":\"issueTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_machineHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_payloadHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"outputs\",\"type\":\"bytes[]\"},{\"internalType\":\"address\",\"name\":\"_coprocessorAdapter\",\"type\":\"address\"}],\"name\":\"solverCallbackOutputsOnly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506105ed8061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610034575f3560e01c80632b2020d614610038578063d1543fdd14610054575b5f5ffd5b610052600480360381019061004d9190610221565b610070565b005b61006e600480360381019061006991906102e7565b6100b3565b005b7f8c891d6fd9e74f3e3543a699806523f3157ed53345e124d00ca981f6e3d2892a848484846040516100a594939291906103e3565b60405180910390a150505050565b5f8190508073ffffffffffffffffffffffffffffffffffffffff166358f6e29f878787876040518563ffffffff1660e01b81526004016100f69493929190610579565b5f604051808303815f87803b15801561010d575f5ffd5b505af115801561011f573d5f5f3e3d5ffd5b50505050505050505050565b5f5ffd5b5f5ffd5b5f819050919050565b61014581610133565b811461014f575f5ffd5b50565b5f813590506101608161013c565b92915050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f84011261018757610186610166565b5b8235905067ffffffffffffffff8111156101a4576101a361016a565b5b6020830191508360018202830111156101c0576101bf61016e565b5b9250929050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6101f0826101c7565b9050919050565b610200816101e6565b811461020a575f5ffd5b50565b5f8135905061021b816101f7565b92915050565b5f5f5f5f606085870312156102395761023861012b565b5b5f61024687828801610152565b945050602085013567ffffffffffffffff8111156102675761026661012f565b5b61027387828801610172565b935093505060406102868782880161020d565b91505092959194509250565b5f5f83601f8401126102a7576102a6610166565b5b8235905067ffffffffffffffff8111156102c4576102c361016a565b5b6020830191508360208202830111156102e0576102df61016e565b5b9250929050565b5f5f5f5f5f60808688031215610300576102ff61012b565b5b5f61030d88828901610152565b955050602061031e88828901610152565b945050604086013567ffffffffffffffff81111561033f5761033e61012f565b5b61034b88828901610292565b9350935050606061035e8882890161020d565b9150509295509295909350565b61037481610133565b82525050565b5f82825260208201905092915050565b828183375f83830152505050565b5f601f19601f8301169050919050565b5f6103b3838561037a565b93506103c083858461038a565b6103c983610398565b840190509392505050565b6103dd816101e6565b82525050565b5f6060820190506103f65f83018761036b565b81810360208301526104098185876103a8565b905061041860408301846103d4565b95945050505050565b5f82825260208201905092915050565b5f819050919050565b5f82825260208201905092915050565b5f610455838561043a565b935061046283858461038a565b61046b83610398565b840190509392505050565b5f61048284848461044a565b90509392505050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f833560016020038436030381126104b3576104b2610493565b5b83810192508235915060208301925067ffffffffffffffff8211156104db576104da61048b565b5b6001820236038313156104f1576104f061048f565b5b509250929050565b5f602082019050919050565b5f6105108385610421565b93508360208402850161052284610431565b805f5b8781101561056757848403895261053c8284610497565b610547868284610476565b9550610552846104f9565b935060208b019a505050600181019050610525565b50829750879450505050509392505050565b5f60608201905061058c5f83018761036b565b610599602083018661036b565b81810360408301526105ac818486610505565b90509594505050505056fea26469706673582212203f82e07158c4409abded5827b411fbe8cbf90c677e7697e5b09c996616b7f36c64736f6c634300081c0033",
}

// MockCoprocessorABI is the input ABI used to generate the binding from.
// Deprecated: Use MockCoprocessorMetaData.ABI instead.
var MockCoprocessorABI = MockCoprocessorMetaData.ABI

// MockCoprocessorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockCoprocessorMetaData.Bin instead.
var MockCoprocessorBin = MockCoprocessorMetaData.Bin

// DeployMockCoprocessor deploys a new Ethereum contract, binding an instance of MockCoprocessor to it.
func DeployMockCoprocessor(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MockCoprocessor, error) {
	parsed, err := MockCoprocessorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockCoprocessorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockCoprocessor{MockCoprocessorCaller: MockCoprocessorCaller{contract: contract}, MockCoprocessorTransactor: MockCoprocessorTransactor{contract: contract}, MockCoprocessorFilterer: MockCoprocessorFilterer{contract: contract}}, nil
}

// MockCoprocessor is an auto generated Go binding around an Ethereum contract.
type MockCoprocessor struct {
	MockCoprocessorCaller     // Read-only binding to the contract
	MockCoprocessorTransactor // Write-only binding to the contract
	MockCoprocessorFilterer   // Log filterer for contract events
}

// MockCoprocessorCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockCoprocessorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockCoprocessorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockCoprocessorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockCoprocessorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockCoprocessorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockCoprocessorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockCoprocessorSession struct {
	Contract     *MockCoprocessor  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MockCoprocessorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockCoprocessorCallerSession struct {
	Contract *MockCoprocessorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// MockCoprocessorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockCoprocessorTransactorSession struct {
	Contract     *MockCoprocessorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// MockCoprocessorRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockCoprocessorRaw struct {
	Contract *MockCoprocessor // Generic contract binding to access the raw methods on
}

// MockCoprocessorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockCoprocessorCallerRaw struct {
	Contract *MockCoprocessorCaller // Generic read-only contract binding to access the raw methods on
}

// MockCoprocessorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockCoprocessorTransactorRaw struct {
	Contract *MockCoprocessorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockCoprocessor creates a new instance of MockCoprocessor, bound to a specific deployed contract.
func NewMockCoprocessor(address common.Address, backend bind.ContractBackend) (*MockCoprocessor, error) {
	contract, err := bindMockCoprocessor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockCoprocessor{MockCoprocessorCaller: MockCoprocessorCaller{contract: contract}, MockCoprocessorTransactor: MockCoprocessorTransactor{contract: contract}, MockCoprocessorFilterer: MockCoprocessorFilterer{contract: contract}}, nil
}

// NewMockCoprocessorCaller creates a new read-only instance of MockCoprocessor, bound to a specific deployed contract.
func NewMockCoprocessorCaller(address common.Address, caller bind.ContractCaller) (*MockCoprocessorCaller, error) {
	contract, err := bindMockCoprocessor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockCoprocessorCaller{contract: contract}, nil
}

// NewMockCoprocessorTransactor creates a new write-only instance of MockCoprocessor, bound to a specific deployed contract.
func NewMockCoprocessorTransactor(address common.Address, transactor bind.ContractTransactor) (*MockCoprocessorTransactor, error) {
	contract, err := bindMockCoprocessor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockCoprocessorTransactor{contract: contract}, nil
}

// NewMockCoprocessorFilterer creates a new log filterer instance of MockCoprocessor, bound to a specific deployed contract.
func NewMockCoprocessorFilterer(address common.Address, filterer bind.ContractFilterer) (*MockCoprocessorFilterer, error) {
	contract, err := bindMockCoprocessor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockCoprocessorFilterer{contract: contract}, nil
}

// bindMockCoprocessor binds a generic wrapper to an already deployed contract.
func bindMockCoprocessor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockCoprocessorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockCoprocessor *MockCoprocessorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockCoprocessor.Contract.MockCoprocessorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockCoprocessor *MockCoprocessorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockCoprocessor.Contract.MockCoprocessorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockCoprocessor *MockCoprocessorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockCoprocessor.Contract.MockCoprocessorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockCoprocessor *MockCoprocessorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockCoprocessor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockCoprocessor *MockCoprocessorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockCoprocessor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockCoprocessor *MockCoprocessorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockCoprocessor.Contract.contract.Transact(opts, method, params...)
}

// IssueTask is a paid mutator transaction binding the contract method 0x2b2020d6.
//
// Solidity: function issueTask(bytes32 machineHash, bytes input, address callback) returns()
func (_MockCoprocessor *MockCoprocessorTransactor) IssueTask(opts *bind.TransactOpts, machineHash [32]byte, input []byte, callback common.Address) (*types.Transaction, error) {
	return _MockCoprocessor.contract.Transact(opts, "issueTask", machineHash, input, callback)
}

// IssueTask is a paid mutator transaction binding the contract method 0x2b2020d6.
//
// Solidity: function issueTask(bytes32 machineHash, bytes input, address callback) returns()
func (_MockCoprocessor *MockCoprocessorSession) IssueTask(machineHash [32]byte, input []byte, callback common.Address) (*types.Transaction, error) {
	return _MockCoprocessor.Contract.IssueTask(&_MockCoprocessor.TransactOpts, machineHash, input, callback)
}

// IssueTask is a paid mutator transaction binding the contract method 0x2b2020d6.
//
// Solidity: function issueTask(bytes32 machineHash, bytes input, address callback) returns()
func (_MockCoprocessor *MockCoprocessorTransactorSession) IssueTask(machineHash [32]byte, input []byte, callback common.Address) (*types.Transaction, error) {
	return _MockCoprocessor.Contract.IssueTask(&_MockCoprocessor.TransactOpts, machineHash, input, callback)
}

// SolverCallbackOutputsOnly is a paid mutator transaction binding the contract method 0xd1543fdd.
//
// Solidity: function solverCallbackOutputsOnly(bytes32 _machineHash, bytes32 _payloadHash, bytes[] outputs, address _coprocessorAdapter) returns()
func (_MockCoprocessor *MockCoprocessorTransactor) SolverCallbackOutputsOnly(opts *bind.TransactOpts, _machineHash [32]byte, _payloadHash [32]byte, outputs [][]byte, _coprocessorAdapter common.Address) (*types.Transaction, error) {
	return _MockCoprocessor.contract.Transact(opts, "solverCallbackOutputsOnly", _machineHash, _payloadHash, outputs, _coprocessorAdapter)
}

// SolverCallbackOutputsOnly is a paid mutator transaction binding the contract method 0xd1543fdd.
//
// Solidity: function solverCallbackOutputsOnly(bytes32 _machineHash, bytes32 _payloadHash, bytes[] outputs, address _coprocessorAdapter) returns()
func (_MockCoprocessor *MockCoprocessorSession) SolverCallbackOutputsOnly(_machineHash [32]byte, _payloadHash [32]byte, outputs [][]byte, _coprocessorAdapter common.Address) (*types.Transaction, error) {
	return _MockCoprocessor.Contract.SolverCallbackOutputsOnly(&_MockCoprocessor.TransactOpts, _machineHash, _payloadHash, outputs, _coprocessorAdapter)
}

// SolverCallbackOutputsOnly is a paid mutator transaction binding the contract method 0xd1543fdd.
//
// Solidity: function solverCallbackOutputsOnly(bytes32 _machineHash, bytes32 _payloadHash, bytes[] outputs, address _coprocessorAdapter) returns()
func (_MockCoprocessor *MockCoprocessorTransactorSession) SolverCallbackOutputsOnly(_machineHash [32]byte, _payloadHash [32]byte, outputs [][]byte, _coprocessorAdapter common.Address) (*types.Transaction, error) {
	return _MockCoprocessor.Contract.SolverCallbackOutputsOnly(&_MockCoprocessor.TransactOpts, _machineHash, _payloadHash, outputs, _coprocessorAdapter)
}

// MockCoprocessorTaskIssuedIterator is returned from FilterTaskIssued and is used to iterate over the raw logs and unpacked data for TaskIssued events raised by the MockCoprocessor contract.
type MockCoprocessorTaskIssuedIterator struct {
	Event *MockCoprocessorTaskIssued // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockCoprocessorTaskIssuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockCoprocessorTaskIssued)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockCoprocessorTaskIssued)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockCoprocessorTaskIssuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockCoprocessorTaskIssuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockCoprocessorTaskIssued represents a TaskIssued event raised by the MockCoprocessor contract.
type MockCoprocessorTaskIssued struct {
	MachineHash [32]byte
	Input       []byte
	Callback    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTaskIssued is a free log retrieval operation binding the contract event 0x8c891d6fd9e74f3e3543a699806523f3157ed53345e124d00ca981f6e3d2892a.
//
// Solidity: event TaskIssued(bytes32 machineHash, bytes input, address callback)
func (_MockCoprocessor *MockCoprocessorFilterer) FilterTaskIssued(opts *bind.FilterOpts) (*MockCoprocessorTaskIssuedIterator, error) {

	logs, sub, err := _MockCoprocessor.contract.FilterLogs(opts, "TaskIssued")
	if err != nil {
		return nil, err
	}
	return &MockCoprocessorTaskIssuedIterator{contract: _MockCoprocessor.contract, event: "TaskIssued", logs: logs, sub: sub}, nil
}

// WatchTaskIssued is a free log subscription operation binding the contract event 0x8c891d6fd9e74f3e3543a699806523f3157ed53345e124d00ca981f6e3d2892a.
//
// Solidity: event TaskIssued(bytes32 machineHash, bytes input, address callback)
func (_MockCoprocessor *MockCoprocessorFilterer) WatchTaskIssued(opts *bind.WatchOpts, sink chan<- *MockCoprocessorTaskIssued) (event.Subscription, error) {

	logs, sub, err := _MockCoprocessor.contract.WatchLogs(opts, "TaskIssued")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockCoprocessorTaskIssued)
				if err := _MockCoprocessor.contract.UnpackLog(event, "TaskIssued", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskIssued is a log parse operation binding the contract event 0x8c891d6fd9e74f3e3543a699806523f3157ed53345e124d00ca981f6e3d2892a.
//
// Solidity: event TaskIssued(bytes32 machineHash, bytes input, address callback)
func (_MockCoprocessor *MockCoprocessorFilterer) ParseTaskIssued(log types.Log) (*MockCoprocessorTaskIssued, error) {
	event := new(MockCoprocessorTaskIssued)
	if err := _MockCoprocessor.contract.UnpackLog(event, "TaskIssued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
