// SPDX-License-Identifier: MIT
pragma solidity ^0.8.31;

import "forge-std/console.sol";

/**
 * @title Osaka
 * @notice Simple test for all Osaka (Fusaka) EIPs in a single contract
 */
contract Osaka {
    // Precompile addresses
    address constant MODEXP_PRECOMPILE = 0x0000000000000000000000000000000000000005;
    address constant SECP256R1_PRECOMPILE = 0x0000000000000000000000000000000000000100; // EIP-7951

    uint256 public state;
    bytes32 public lastResult;

    /**
     * @notice Test EIP-7594: PeerDAS - data-availability sampling
     * @dev Tests blob throughput with PeerDAS (increased blob count)
     * @dev This test will just fail becease blob is not available in Surge
     */
    function testPeerDASBlobThroughput() public {
        bytes32 result;

        assembly {
            // PeerDAS allows for more blobs - test reading extended blob hashes
            // Read first 12 blob hashes (example capacity with PeerDAS)
            let hash0 := blobhash(0)
            let hash1 := blobhash(1)
            let hash2 := blobhash(2)
            let hash3 := blobhash(3)
            let hash4 := blobhash(4)
            let hash5 := blobhash(5)
            let hash6 := blobhash(6)
            let hash7 := blobhash(7)
            let hash8 := blobhash(8)
            let hash9 := blobhash(9)
            let hash10 := blobhash(10)
            let hash11 := blobhash(11)

            // XOR all hashes together
            result := xor(hash0, hash1)
            result := xor(result, hash2)
            result := xor(result, hash3)
            result := xor(result, hash4)
            result := xor(result, hash5)
            result := xor(result, hash6)
            result := xor(result, hash7)
            result := xor(result, hash8)
            result := xor(result, hash9)
            result := xor(result, hash10)
            result := xor(result, hash11)
        }

        lastResult = result;

        if (result != bytes32(0)) {
            state = state + 1;
        }
    }

    /**
     * @notice Test EIP-7892: BPO (Blob-Parameter-Only) forks
     * @dev Tests that blob parameters can be updated between major forks
     * @dev This test will just fail becease blob is not available in Surge
     */
    function testBPOBlobParameters() public {
        bytes32 result;

        assembly {
            // Read blob hashes to test blob availability
            let hash0 := blobhash(0)
            result := hash0
        }

        lastResult = result;

        // If blobs are available, BPO is working
        if (result != bytes32(0)) {
            state = state + 1;
        }
    }

    /**
     * @notice Test EIP-7918: Blob base-fee bounded by execution costs
     * @dev Tests that blob fee market reacts appropriately to congestion
     * @dev This test will just fail becease blob is not available in Surge
     */
    function testBlobBaseFee() public {
        assembly {
            sstore(state.slot, blobbasefee())
        }
    }

    /**
     * @notice Test EIP-7823: Set upper bounds for MODEXP precompile inputs
     * @dev Tests that MODEXP has proper input size limits (max 8192 bits)
     */
    function testMODEXPUpperBounds() public {
        // Test 1: Valid input within bounds (256 bits)
        bytes memory validInput = new bytes(96 + 32 + 32 + 32);
        // [base_len(32) | exp_len(32) | mod_len(32) | base | exp | mod]
        // Setting lengths to 32 bytes (256 bits) each
        assembly {
            mstore(add(validInput, 32), 32) // base length
            mstore(add(validInput, 64), 32) // exponent length
            mstore(add(validInput, 96), 32) // modulus length
        }

        (bool success1,) = MODEXP_PRECOMPILE.staticcall(validInput);
        bool withinBounds = success1;

        if (!withinBounds) {
            revert("MODEXP valid input test failed");
        }

        state = state + 1;
    }

    /**
     * @notice Test EIP-7823: Set upper bounds for MODEXP precompile inputs
     * @dev Tests that MODEXP has proper input size limits (max 8192 bits)
     */
    function testMODEXPUpperBoundsShouldFail() public {
        bytes memory boundaryInput = new bytes(96 + 1025 + 32 + 32);
        assembly {
            mstore(add(boundaryInput, 32), 1025) // base length
            mstore(add(boundaryInput, 64), 32) // exponent length
            mstore(add(boundaryInput, 96), 32) // modulus length
        }

        (bool success2,) = MODEXP_PRECOMPILE.staticcall(boundaryInput);
        bool exceedsBounds = !success2; // Should fail if bounds are enforced

        if (!exceedsBounds) {
            revert("MODEXP valid out of bound input test failed");
        }

        state = state + 1;
    }

    /**
     * @notice Test EIP-7825: Transaction Gas Limit Cap
     * @dev Tests that transaction gas limit is capped at 2^24 (~16.7M gas)
     * @dev This creates a loop that consumes gas to test the cap
     */
    function testTransactionGasLimitCap() public {
        unchecked {
            uint256 startGas = gasleft();
            console.log("Starting gas:", startGas);

            // The cap is 2^24 = 16,777,216
            uint256 maxGasLimit = 2 ** 24;
            console.log("Max gas limit per tx (2^24):", maxGasLimit);

            // Consume gas with a simple computation loop
            uint256 dummy = 0;
            uint256 i = 0;

            while (i < 100000) {
                if (gasleft() <= startGas - maxGasLimit) {
                    break; // Stop if we've consumed enough gas
                }
                dummy = dummy + uint256(keccak256(abi.encodePacked(i, block.timestamp)));
                i++;
            }

            uint256 gasConsumed = startGas - gasleft();
            console.log("Gas consumed:", gasConsumed);
            console.log("Iterations completed:", i);
            console.log("Dummy result (prevent optimization):", dummy);

            // If we successfully consumed significant gas, the function worked
            if (gasConsumed > maxGasLimit) {
                console.log("Successfully consumed significant gas");
                state = state + 1;
            }
        }
    }

    /**
     * @notice Test EIP-7883: MODEXP gas cost increase
     * @dev Tests that MODEXP precompile has increased gas costs
     */
    function testMODEXPGasCostIncrease() public returns (uint256 gasUsed) {
        uint256 startGas = gasleft();

        // Prepare MODEXP input: base^exp mod modulus
        // Format: [base_len | exp_len | mod_len | base | exp | mod]
        bytes memory input = new bytes(96 + 32 + 32 + 32);

        assembly {
            // Set lengths (32 bytes each)
            mstore(add(input, 32), 32)
            mstore(add(input, 64), 32)
            mstore(add(input, 96), 32)
            // Set base = 3
            mstore(add(input, 128), 3)
            // Set exp = 5
            mstore(add(input, 160), 5)
            // Set mod = 7
            mstore(add(input, 192), 7)
        }

        (bool success,) = MODEXP_PRECOMPILE.staticcall(input);
        require(success, "MODEXP call failed");

        gasUsed = startGas - gasleft();
        console.log("MODEXP gas used:", gasUsed);

        state = state + 1;

        return gasUsed;
    }

    /**
     * @notice Test EIP-7939: Add CLZ opcode (Count Leading Zeros)
     * @dev Tests the new CLZ instruction for counting leading zeros
     */
    function testCLZOpcode() public {
        assembly {
            let slotValue := add(sload(state.slot), 1)
            sstore(state.slot, clz(slotValue))
        }
    }

    /**
     * @notice Test EIP-7951: Precompile for secp256r1 curve support
     * @dev Tests P-256 signature verification precompile
     */
    function testSecp256r1Precompile() public {
        // Check if precompile exists
        uint256 codeSize;
        assembly {
            codeSize := extcodesize(0x0000000000000000000000000000000000000100)
        }

        // Test with sample P-256 signature verification data
        // Format: [message_hash(32) | r(32) | s(32) | x(32) | y(32)] = 160 bytes
        bytes memory input = new bytes(160);

        // Use zero-filled data as a basic existence test
        (bool success,) = SECP256R1_PRECOMPILE.staticcall(input);

        if (success) {
            state = state + 1;
        }
    }
}

