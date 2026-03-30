// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

/**
 * @title Cancun
 * @notice Simple test for all Cancun EIPs in a single contract
 */
contract Cancun {
    address constant BEACON_ROOTS_ADDRESS = 0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;

    uint256 public state;

    /**
     * @notice Test the EIP4844 BLOBBASEFEE opcode
     */
    function testBlobBaseFee() public {
        assembly {
            sstore(state.slot, blobbasefee())
        }
    }

    function testSomething() public returns (uint256) {
        return state;
    }

    /**
     * @notice Test the EIP5656 MCOPY opcode
     */
    function testMCopy() public {
        bytes memory data = new bytes(32);

        // Fill source with test data
        assembly {
            mstore(add(data, 0x20), 0x123456789ABCDEF)
        }

        // Use MCOPY to copy
        assembly {
            let sourcePtr := add(data, 0x20)
            mcopy(state.slot, sourcePtr, 32)
        }

        state = state + 1;
    }

    /**
     * @notice Test the EIP1153 transient storage
     */
    function testTransientStorage() public {
        uint256 slot = 0;
        uint256 value = 42;

        // Store value in transient storage
        assembly {
            tstore(slot, value)
            sstore(state.slot, tload(slot))
        }
    }

    /**
     * @notice Try to the EIP4788 beacon root access
     */
    function getBeaconRoot() public {
        uint256 timestamp = block.timestamp;
        (bool success, bytes memory data) = BEACON_ROOTS_ADDRESS.staticcall(abi.encode(timestamp));

        if (success && data.length >= 32) {
            state = abi.decode(data, (uint256));
            revert("Beacon root retrieved");
        }
    }

    /**
     * @notice Test EIP-4844 BLOBHASH opcode
     * @dev This function reads blob hashes from a blob-carrying transaction
     * @dev Must be called from a type 3 (blob) transaction to return non-zero hashes
     */
    function testBlobHash() public {
        assembly {
            sstore(state.slot, blobhash(0))
        }
    }

    /**
     * @notice Process data with blob hashes
     * @dev This is a simple function that can be called from a blob transaction
     * @param data Some calldata to process
     */
    function processBlobData(bytes calldata data) public returns (bytes32 firstBlobHash, uint256 dataLength) {
        // Read the first blob hash
        assembly {
            firstBlobHash := blobhash(0)
        }

        dataLength = data.length;

        // Store blob hash info
        if (firstBlobHash != bytes32(0)) {
            assembly {
                sstore(state.slot, firstBlobHash)
            }
            state = state + 1;
        }

        return (firstBlobHash, dataLength);
    }

    /**
     * @notice Test the EIP6780 SELFDESTRUCT behavior
     * @dev Creates and destroys a contract, then returns its address for later verification
     */
    function testSameTransactionDestroy() public {
        EIP6780_Victim victim = EIP6780_Victim(payable(deployVictim(true)));

        // Note: EXTCODESIZE will still return non-zero within the same transaction
        // The account is only deleted during state finalization at transaction end
        state = state + 1;

        return;
    }

    function deployVictim(bool is_destroyed) public returns (address) {
        EIP6780_Victim victim = new EIP6780_Victim(is_destroyed);
        return address(victim);
    }

    function testDifferentTransactionDestroy(address victimAddress) public {
        EIP6780_Victim victim = EIP6780_Victim(payable(victimAddress));

        // Destroy in same transaction (should delete code under EIP-6780)
        victim.destroy(payable(msg.sender));

        // Note: EXTCODESIZE will still return non-zero within the same transaction
        // The account is only deleted during state finalization at transaction end
        state = state + 1;

        return;
    }
}

/**
 * @title EIP6780_SelfDestruct
 * @notice Simple test for EIP-6780: SELFDESTRUCT only in same transaction
 * @dev Under EIP-6780, SELFDESTRUCT only deletes code if called in the same transaction as creation
 */
contract EIP6780_Victim {
    uint256 public value = 42;

    constructor(bool is_destroyed) {
        if (is_destroyed) {
            selfdestruct(payable(msg.sender));
        }
    }

    function destroy(address payable beneficiary) public {
        selfdestruct(beneficiary);
    }

    receive() external payable {}
}
