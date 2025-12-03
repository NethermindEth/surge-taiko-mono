// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "forge-std/Script.sol";
import "../src/Prague.sol";

/**
 * @title ExecutePragueEIPTests
 * @notice Deploy and execute all Prague EIP tests ON-CHAIN
 * @dev This will create real transactions and modify blockchain state
 */
contract ExecuteEIPTests is Script {
    function run() external {
        uint256 deployerPrivateKey = 0x94eb3102993b41ec55c241060f47daa0f6372e2e3ad7e91612ae36c364042e44;
        vm.startBroadcast(deployerPrivateKey);

        console.log("\n========================================");
        console.log("EXECUTING PRAGUE EIP TESTS ON-CHAIN");
        console.log("========================================\n");

        Prague pragueTests = new Prague();
        console.log("Prague contract deployed at:", address(pragueTests));

        // EIP-7702: Set EOA account code
        console.log("\n--- Testing EIP-7702: Set EOA account code ---");
        pragueTests.testSetEOACode(address(pragueTests));
        console.log("EIP-7702 test completed, state:", pragueTests.state());

        // EIP-7691: Blob throughput increase
        console.log("\n--- Testing EIP-7691: Blob throughput increase (9 blobs) ---");
        pragueTests.testBlobThroughput();
        console.log("EIP-7691 test completed");

        // EIP-7623: Increase calldata cost
        console.log("\n--- Testing EIP-7623: Increase calldata cost ---");
        bytes memory largeCalldata = new bytes(10000);
        for (uint256 i = 0; i < 10000; i++) {
            largeCalldata[i] = bytes1(uint8(i % 256));
        }
        (uint256 dataLength, uint256 gasUsed) = pragueTests.testCalldataCost(largeCalldata);
        console.log("EIP-7623 test completed, data length:", dataLength, "gas used:", gasUsed);

        // Individual BLS tests for detailed verification
        console.log("\n--- Testing individual BLS precompiles ---");
        pragueTests.testBLS12G1Add();
        pragueTests.testBLS12G1MSM();
        pragueTests.testBLS12G2Add();
        pragueTests.testBLS12G2MSM();
        pragueTests.testBLS12Pairing();
        pragueTests.testBLS12MapFpToG1();
        pragueTests.testBLS12MapFp2ToG2();

        // EIP-2935: Historical block hashes
        console.log("\n--- Testing EIP-2935: Historical block hashes ---");
        pragueTests.testHistoricalBlockHash(block.number > 0 ? block.number - 1 : 0);

        // EIP-7685: General purpose execution layer requests
        console.log("\n--- Testing EIP-7685: Execution layer requests ---");
        pragueTests.testExecutionLayerRequests();

        console.log("\n========================================");
        console.log("TESTS COMPLETE");
        console.log("Final state:", pragueTests.state());
        console.log("========================================\n");

        vm.stopBroadcast();
    }
}
