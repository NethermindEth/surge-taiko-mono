// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "forge-std/Script.sol";
import "../src/Cancun.sol";

/**
 * @title ExecuteEIPTests
 * @notice Deploy and execute all Cancun EIP tests ON-CHAIN
 * @dev This will create real transactions and modify blockchain state
 */
contract ExecuteEIPTests is Script {
    function run() external {
        uint256 deployerPrivateKey = 0x94eb3102993b41ec55c241060f47daa0f6372e2e3ad7e91612ae36c364042e44;
        vm.startBroadcast(deployerPrivateKey);

        console.log("\n========================================");
        console.log("EXECUTING CANCUN EIP TESTS ON-CHAIN");
        console.log("========================================\n");

        Cancun cancunTests = new Cancun();

        cancunTests.testBlobBaseFee();

        cancunTests.testMCopy();

        cancunTests.testTransientStorage();

        cancunTests.getBeaconRoot();

        cancunTests.testSameTransactionDestroy();

        address destroyedContract2 = cancunTests.deployVictim(false);
        cancunTests.testDifferentTransactionDestroy(destroyedContract2);

        cancunTests.testBlobHash();

        console.log("\n========================================");
        console.log("TESTS COMPLETE");
        console.log("========================================\n");

        vm.stopBroadcast();
    }
}
