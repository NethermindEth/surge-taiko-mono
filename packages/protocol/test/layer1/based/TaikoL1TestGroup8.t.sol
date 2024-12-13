// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./TaikoL1TestGroupBase.sol";

contract TaikoL1TestGroup8 is TaikoL1TestGroupBase {
    // Test summary:
    // 1. Gets a block that doesn't exist
    // 2. Gets a transition by ID & hash that doesn't exist.
    function test_taikoL1_group_8_case_3() external {
        vm.expectRevert(LibUtils.L1_INVALID_BLOCK_ID.selector);
        L1.getBlock(2);

        vm.expectRevert(LibUtils.L1_TRANSITION_NOT_FOUND.selector);
        L1.getTransition(0, 2);

        vm.expectRevert(LibUtils.L1_TRANSITION_NOT_FOUND.selector);
        L1.getTransition(0, randBytes32());

        vm.expectRevert(LibUtils.L1_INVALID_BLOCK_ID.selector);
        L1.getTransition(3, randBytes32());
    }
}
