// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "forge-std/console.sol";

/**
 * @title Prague
 * @notice Simple test for all Prague EIPs in a single contract
 */
contract Prague {
    // Precompile addresses for EIP-2537 (BLS12-381)
    address constant BLS12_G1ADD = 0x000000000000000000000000000000000000000b;
    address constant BLS12_G1MSM = 0x000000000000000000000000000000000000000C;
    address constant BLS12_G2ADD = 0x000000000000000000000000000000000000000d;
    address constant BLS12_G2MSM = 0x000000000000000000000000000000000000000E;

    address constant BLS12_PAIRING = 0x000000000000000000000000000000000000000F;
    address constant BLS12_MAP_FP_TO_G1 = 0x0000000000000000000000000000000000000010;
    address constant BLS12_MAP_FP2_TO_G2 = 0x0000000000000000000000000000000000000011;

    // Contract address for EIP-2935 (historical block hashes)
    address constant HISTORICAL_BLOCKHASH_ADDRESS = 0x0000F90827F1C53a10cb7A02335B175320002935;

    // System contract addresses for EIP-6110, EIP-7002, EIP-7251
    address constant DEPOSIT_CONTRACT_ADDRESS = 0x00000000219ab540356cBB839Cbe05303d7705Fa;
    address constant WITHDRAWAL_REQUEST_ADDRESS = 0x00000961Ef480Eb55e80D19ad83579A64c007002;
    address constant CONSOLIDATION_REQUEST_ADDRESS = 0x0000BBdDc7CE488642fb579F8B00f3a590007251;

    uint256 public state;
    bytes32 public lastResult;

    /**
     * @notice Test EIP-7702: Set EOA account code
     * @dev This tests delegation designation functionality
     * @dev EOA can designate a contract to execute on its behalf using auth tuples
     * @param delegateContract Address to delegate to
     */
    function testSetEOACode(address delegateContract) public {
        // Under EIP-7702, an EOA can have code set via authorization
        // This would typically be done via transaction type 0x04 with auth tuples
        // Here we test interaction with potentially delegated EOAs

        uint256 codeSize;
        assembly {
            codeSize := extcodesize(delegateContract)
        }

        state = codeSize;
    }

    /**
     * @notice Test EIP-7691: Blob throughput increase
     * @dev Tests the new blob gas parameters (target: 6, max: 9 blobs)
     * @dev Reads blob hashes which would be populated in type 3 transactions
     */
    function testBlobThroughput() public {
        bytes32 result;

        assembly {
            // Read all 9 blob hashes and XOR them together
            let hash0 := blobhash(0)
            let hash1 := blobhash(1)
            let hash2 := blobhash(2)
            let hash3 := blobhash(3)
            let hash4 := blobhash(4)
            let hash5 := blobhash(5)
            let hash6 := blobhash(6)
            let hash7 := blobhash(7)
            let hash8 := blobhash(8)

            // XOR all hashes together
            result := xor(hash0, hash1)
            result := xor(result, hash2)
            result := xor(result, hash3)
            result := xor(result, hash4)
            result := xor(result, hash5)
            result := xor(result, hash6)
            result := xor(result, hash7)
            result := xor(result, hash8)
        }

        // Store XOR of all 9 blob hashes (max blobs in EIP-7691)
        lastResult = result;

        if (result != bytes32(0)) {
            state = state + 1;
        }
    }

    /**
     * @notice Test EIP-7623: Increase calldata cost
     * @dev Tests behavior with increased calldata costs (new floor pricing)
     * @dev This primarily affects gas accounting, observable through gas usage
     * @param data Large calldata to test cost increase
     */
    function testCalldataCost(bytes calldata data) public returns (uint256 dataLength, uint256 gasUsed) {
        uint256 startGas = gasleft();

        // Process the calldata
        dataLength = data.length;

        // Do some operations with the calldata
        bytes32 hash = keccak256(data);
        lastResult = hash;

        gasUsed = startGas - gasleft();
        state = dataLength;

        return (dataLength, gasUsed);
    }

    /**
     * @notice Test EIP-2537: BLS12-381 G1 addition precompile
     * @dev Tests the G1 point addition precompile
     */
    function testBLS12G1Add() public {
        // Two valid G1 points (128 bytes each) - using generator point G1 twice
        bytes memory input = new bytes(256);

        (bool success,) = BLS12_G1ADD.staticcall(input);

        if (success) {
            state = state + 1;
        } else {
            revert("BLS12 G1 Add failed");
        }
    }

    /**
     * @notice Test EIP-2537: BLS12-381 G1 multi-scalar multiplication precompile
     * @dev Tests the G1 MSM precompile
     */
    function testBLS12G1MSM() public {
        // MSM with 1 point: G1 point (128 bytes) + scalar (32 bytes) = 160 bytes
        bytes memory input = new bytes(160);

        (bool success,) = BLS12_G1MSM.staticcall(input);
        if (success) {
            state = state + 1;
        } else {
            revert("BLS12 G1 MSM failed");
        }
    }

    /**
     * @notice Test EIP-2537: BLS12-381 G2 addition precompile
     * @dev Tests the G2 point addition precompile
     */
    function testBLS12G2Add() public {
        // Two G2 points (256 bytes each) = 512 bytes total
        bytes memory input = new bytes(512);

        (bool success,) = BLS12_G2ADD.staticcall(input);

        if (success) {
            state = state + 1;
        } else {
            revert("BLS12 G2 Add failed");
        }
    }

    /**
     * @notice Test EIP-2537: BLS12-381 G2 multi-scalar multiplication precompile
     * @dev Tests the G2 MSM precompile
     */
    function testBLS12G2MSM() public {
        // MSM with 1 point: G2 point (256 bytes) + scalar (32 bytes) = 288 bytes
        bytes memory input = new bytes(288);

        (bool success,) = BLS12_G2MSM.staticcall(input);

        if (success) {
            state = state + 1;
        } else {
            revert("BLS12 G2 MSM failed");
        }
    }

    /**
     * @notice Test EIP-2537: BLS12-381 pairing precompile
     * @dev Tests the pairing check precompile
     */
    function testBLS12Pairing() public {
        // One pairing: G1 point (128 bytes) + G2 point (256 bytes) = 384 bytes
        bytes memory input = new bytes(384);

        (bool success,) = BLS12_PAIRING.staticcall(input);
        if (success) {
            state = state + 1;
        } else {
            revert("BLS12 Pairing failed");
        }
    }

    /**
     * @notice Test EIP-2537: BLS12-381 map field element to G1 precompile
     * @dev Tests the map Fp to G1 precompile
     */
    function testBLS12MapFpToG1() public {
        // Field element: 64 bytes
        bytes memory input = new bytes(64);

        (bool success,) = BLS12_MAP_FP_TO_G1.staticcall(input);

        if (success) {
            state = state + 1;
        } else {
            revert("BLS12 Map Fp to G1 failed");
        }
    }

    /**
     * @notice Test EIP-2537: BLS12-381 map field element to G2 precompile
     * @dev Tests the map Fp2 to G2 precompile
     */
    function testBLS12MapFp2ToG2() public {
        // Fp2 element: 128 bytes (two 64-byte field elements)
        bytes memory input = new bytes(128);

        (bool success,) = BLS12_MAP_FP2_TO_G2.staticcall(input);

        if (success) {
            state = state + 1;
        } else {
            revert("BLS12 Map Fp2 to G2 failed");
        }
    }

    /**
     * @notice Test EIP-2935: Save historical block hashes in state
     * @dev Accesses historical block hash from the new system contract
     * @param blockNumber The block number to query (must be within HISTORY_SERVE_WINDOW)
     */
    function testHistoricalBlockHash(uint256 blockNumber) public returns (bytes32 blockHash) {
        // Call the historical blockhash contract
        (bool success, bytes memory data) = HISTORICAL_BLOCKHASH_ADDRESS.staticcall(abi.encode(blockNumber));

        console.log("Called historical blockhash for block:", blockNumber);
        console.log("success", blockNumber, "data length:", data.length);

        if (success && data.length >= 32) {
            blockHash = abi.decode(data, (bytes32));
            lastResult = blockHash;
            state = state + 1;
        } else if (!success) {
            revert("Historical blockhash call failed");
        }

        return blockHash;
    }

    /**
     * @notice Test EIP-7685: General purpose execution layer requests
     * @dev Tests the general request framework (deposits, exits, consolidations)
     */
    function testExecutionLayerRequests() public {
        // EIP-7685 provides a general framework for EL requests
        // This encompasses deposits (EIP-6110), exits (EIP-7002), and consolidations (EIP-7251)

        // Check all three system contracts exist
        uint256 depositCodeSize;
        uint256 withdrawalCodeSize;
        uint256 consolidationCodeSize;
        uint256 systemContractsFound = 0;

        assembly {
            depositCodeSize := extcodesize(DEPOSIT_CONTRACT_ADDRESS)
            withdrawalCodeSize := extcodesize(WITHDRAWAL_REQUEST_ADDRESS)
            consolidationCodeSize := extcodesize(CONSOLIDATION_REQUEST_ADDRESS)
        }

        if (depositCodeSize > 0) systemContractsFound++;
        if (withdrawalCodeSize > 0) systemContractsFound++;
        if (consolidationCodeSize > 0) systemContractsFound++;

        state = state + systemContractsFound;
    }
}

//Simple contract that returns a constant value and has a state variable
contract Test {
    uint256 public value;

    function getValue() public view returns (uint256) {
        return value;
    }

    function setValue(uint256 _value) public {
        value = _value;
    }

    function getConstant() public pure returns (uint256) {
        return 12345;
    }
}
