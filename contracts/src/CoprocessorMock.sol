// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "./ICoprocessorCallback.sol";

contract CoprocessorMock {
    event TaskIssued(bytes32 machineHash, bytes input, address callback);

    function issueTask(bytes32 machineHash, bytes calldata input, address callback) public {
        emit TaskIssued(machineHash, input, callback);
    }

    function solverCallbackOutputsOnly(
        bytes32 _machineHash, 
        bytes32 _payloadHash, 
        bytes[] calldata _outputs,
        address _callback
    ) public {
        ICoprocessorCallback callbackContract = ICoprocessorCallback(_callback);
        callbackContract.coprocessorCallbackOutputsOnly(_machineHash, _payloadHash, _outputs);
    }
}