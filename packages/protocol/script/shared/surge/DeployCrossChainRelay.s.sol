// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { Script } from "forge-std/src/Script.sol";
import { console2 } from "forge-std/src/console2.sol";

import { CrossChainRelay } from "../../../contracts/shared/bridge/CrossChainRelay.sol";

contract DeployCrossChainRelay is Script {
    function run() external {
        vm.startBroadcast();
        CrossChainRelay relay = new CrossChainRelay();
        console2.log("CrossChainRelay:", address(relay));
        vm.stopBroadcast();
    }
}
