// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./TaikoL1TestGroupBase.sol";

contract TaikoL1TestGroup7 is TaikoL1TestGroupBase {
    // Test summary:
    // 1. Alice proposes a block,
    // 2. Alice proves the block within the proving window, using the correct parent hash.
    // 3. After the cooldown window, Taylor contests Alice's proof, and fails.
    function test_taikoL1_group_7_case_1() external {
        vm.warp(1_000_000);

        giveEthAndDepositBond(Alice, 1000 ether, 1000 ether);
        giveEthAndDepositBond(Taylor, 1000 ether, 1000 ether);
        ITierProvider.Tier memory tierOp = TestTierProvider(cp).getTier(LibTiers.TIER_OPTIMISTIC);

        console2.log("====== Alice propose a block");
        TaikoData.BlockMetadata memory meta = proposeBlock(Alice, "");

        console2.log("====== Alice proves the block as the assigned prover");
        bytes32 parentHash = GENESIS_BLOCK_HASH;
        bytes32 blockHash = bytes32(uint256(10));
        bytes32 stateRoot = bytes32(uint256(11));

        mineAndWrap(10 seconds);
        proveBlock(Alice, meta, parentHash, blockHash, stateRoot, meta.minTier, "");

        mineAndWrap(tierOp.cooldownWindow * 60);
        bytes32 blockHash2 = bytes32(uint256(20));
        bytes32 stateRoot2 = bytes32(uint256(21));
        proveBlock(
            Taylor,
            meta,
            parentHash,
            blockHash2,
            stateRoot2,
            meta.minTier,
            LibProving.L1_CANNOT_CONTEST.selector
        );
        printBlockAndTrans(meta.id);
    }

    // Test summary:
    // 1. Alice proposes a block,
    // 2. Alice proves the block within the proving window, using the correct parent hash.
    // 3. Taylor contests Alice's proof.
    // 4. William attempts but fails to contest Alice again.
    function test_taikoL1_group_7_case_2() external {
        vm.warp(1_000_000);

        giveEthAndDepositBond(Alice, 1000 ether, 1000 ether);
        giveEthAndDepositBond(Taylor, 1000 ether, 1000 ether);
        giveEthAndDepositBond(William, 1000 ether, 1000 ether);
        ITierProvider.Tier memory tierOp = TestTierProvider(cp).getTier(LibTiers.TIER_OPTIMISTIC);

        console2.log("====== Alice propose a block");
        TaikoData.BlockMetadata memory meta = proposeBlock(Alice, "");

        console2.log("====== Alice proves the block as the assigned prover");
        bytes32 parentHash = GENESIS_BLOCK_HASH;
        bytes32 blockHash = bytes32(uint256(10));
        bytes32 stateRoot = bytes32(uint256(11));

        mineAndWrap(10 seconds);
        proveBlock(Alice, meta, parentHash, blockHash, stateRoot, meta.minTier, "");

        mineAndWrap(tierOp.cooldownWindow * 60 - 1);
        bytes32 blockHash2 = bytes32(uint256(20));
        bytes32 stateRoot2 = bytes32(uint256(21));
        proveBlock(Taylor, meta, parentHash, blockHash2, stateRoot2, meta.minTier, "");

        bytes32 blockHash3 = bytes32(uint256(30));
        bytes32 stateRoot3 = bytes32(uint256(31));
        proveBlock(
            William,
            meta,
            parentHash,
            blockHash3,
            stateRoot3,
            meta.minTier,
            LibProving.L1_ALREADY_CONTESTED.selector
        );

        printBlockAndTrans(meta.id);
    }
}
