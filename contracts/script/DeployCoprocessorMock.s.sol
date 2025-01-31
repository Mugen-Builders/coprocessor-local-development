// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {CoprocessorMock} from "../src/CoprocessorMock.sol";

contract DeployCoprocessorMock is Script {
    CoprocessorMock public coprocessorMock;

    function setUp() public {}

    function run() public {
        vm.startBroadcast();

        coprocessorMock = new CoprocessorMock();

        vm.stopBroadcast();
        
        console.log("Deployed CoprocessorMock at address: {}", address(coprocessorMock));
    }
}
