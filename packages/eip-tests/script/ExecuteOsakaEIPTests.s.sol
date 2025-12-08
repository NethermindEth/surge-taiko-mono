// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "forge-std/Script.sol";
import "../src/Osaka.sol";

/**
 * @title ExecuteOsakaEIPTests
 * @notice Deploy and execute all Osaka (Fusaka) EIP tests ON-CHAIN
 * @dev This will create real transactions and modify blockchain state
 */
contract ExecuteOsakaEIPTests is Script {
    function run() external {
        uint256 deployerPrivateKey = 0x94eb3102993b41ec55c241060f47daa0f6372e2e3ad7e91612ae36c364042e44;
        vm.startBroadcast(deployerPrivateKey);

        console.log("\n========================================");
        console.log("EXECUTING OSAKA (FUSAKA) EIP TESTS ON-CHAIN");
        console.log("========================================\n");

        Osaka osakaTests = new Osaka();
        console.log("Osaka contract deployed at:", address(osakaTests));

        // EIP-7594: PeerDAS - blob throughput
        console.log("\n--- Testing EIP-7594: PeerDAS blob throughput ---");
        osakaTests.testPeerDASBlobThroughput();
        console.log("EIP-7594 test completed, state:", osakaTests.state());

        // EIP-7892: BPO (Blob-Parameter-Only) forks
        console.log("\n--- Testing EIP-7892: BPO blob parameters ---");
        osakaTests.testBPOBlobParameters();
        console.log("EIP-7892 test completed, state:", osakaTests.state());

        // EIP-7918: Blob base-fee bounded
        console.log("\n--- Testing EIP-7918: Blob base-fee ---");
        osakaTests.testBlobBaseFee();
        console.log("EIP-7918 test completed");

        // EIP-7823: MODEXP upper bounds
        console.log("\n--- Testing EIP-7823: MODEXP upper bounds ---");
        osakaTests.testMODEXPUpperBounds();
        console.log("EIP-7823 test completed");

        // EIP-7823: MODEXP upper bounds should fail
        console.log("\n--- Testing EIP-7823: MODEXP upper bounds, should fail ---");
        osakaTests.testMODEXPUpperBoundsShouldFail();
        console.log("EIP-7823 test completed");

        // // EIP-7883: MODEXP gas cost increase
        console.log("\n--- Testing EIP-7883: MODEXP gas cost increase ---");
        osakaTests.testMODEXPGasCostIncrease();
        console.log("EIP-7883 test completed");

        // EIP-7939: CLZ opcode
        console.log("\n--- Testing EIP-7939: CLZ opcode ---");
        osakaTests.testCLZOpcode();
        console.log("EIP-7939 test completed");

        // EIP-7951: secp256r1 precompile
        console.log("\n--- Testing EIP-7951: secp256r1 precompile ---");
        osakaTests.testSecp256r1Precompile();
        console.log("EIP-7951 test completed");

        // EIP-7825: Transaction gas limit cap
        console.log("\n--- Testing EIP-7825: Transaction gas limit cap ---");
        osakaTests.testTransactionGasLimitCap();
        console.log("EIP-7825 test completed");

        console.log("\n========================================");
        console.log("TESTS COMPLETE");
        console.log("Final state:", osakaTests.state());
        console.log("========================================\n");

        vm.stopBroadcast();
    }
}
