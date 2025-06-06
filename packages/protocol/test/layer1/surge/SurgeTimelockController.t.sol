// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.24;

import "../../shared/CommonTest.sol";
import "src/layer1/surge/common/SurgeTimelockController.sol";
import "@openzeppelin/contracts/governance/TimelockController.sol";

contract MockTaikoInbox {
    uint256 internal verificationStreakStartedAt;

    function getVerificationStreakStartedAt() external view returns (uint256) {
        return verificationStreakStartedAt;
    }

    function resetVerificationStreakStartedAt() external {
        verificationStreakStartedAt = block.timestamp;
    }
}

contract MockStore {
    uint256 internal store;

    function setStore(uint256 _store) external {
        store = _store;
    }

    function getStore() external view returns (uint256) {
        return store;
    }
}

contract SurgeTimelockControllerTestBase is CommonTest {
    MockTaikoInbox internal mockTaikoInbox;
    SurgeTimelockController internal timelockController;
    MockStore internal mockStore;

    // For the calls
    address internal target;
    uint256 internal value;
    bytes internal data;
    bytes32 internal predecessor;
    bytes32 internal salt;

    // For the timelock and stage2
    uint256 internal minDelay;
    uint256 internal minDelayPlusOne;
    uint256 internal startTimestamp;
    uint256 internal minVerificationStreak;

    function setUpOnEthereum() internal override {
        mockTaikoInbox = new MockTaikoInbox();
        mockStore = new MockStore();

        minVerificationStreak = 45 days;
        minDelay = 45 days;
        minDelayPlusOne = minDelay + 1;

        address[] memory proposers = new address[](1);
        address[] memory executors = new address[](1);
        proposers[0] = Alice;
        executors[0] = Alice;

        timelockController = new SurgeTimelockController(
            uint64(minVerificationStreak), // _minVerificationStreak
            uint64(minDelay), // _minDelay
            proposers,
            executors,
            address(0)
        );
        timelockController.init(address(mockTaikoInbox), address(0));

        target = address(mockStore);
        value = 0;
        data = abi.encodeWithSelector(MockStore.setStore.selector, 1);
        predecessor = bytes32(0);
        salt = bytes32(0);
        minDelayPlusOne = 45 days + 1;
        startTimestamp = block.timestamp;
    }
}

contract SurgeTimelockControllerTest is SurgeTimelockControllerTestBase {
    function test_execute_fails_when_min_delay_is_not_met() external transactBy(Alice) {
        // Schedule a call to set the store to 1
        timelockController.schedule(target, value, data, predecessor, salt, minDelayPlusOne);

        // Warp time to before the timelock delay is covered
        vm.warp(startTimestamp + minDelay - 1);

        // Execution reverts because the timelock delay is not covered
        vm.expectRevert("TimelockController: operation is not ready");
        timelockController.execute(target, value, data, predecessor, salt);
    }

    function test_execute_reverts_when_verification_streak_is_disrupted()
        external
        transactBy(Alice)
    {
        // Schedule a call to set the store to 1
        timelockController.schedule(target, value, data, predecessor, salt, minDelayPlusOne);

        // Warp time to somewhere between day 0 and day 45, to reset the verification streak
        vm.warp(startTimestamp + minDelayPlusOne / 2);
        mockTaikoInbox.resetVerificationStreakStartedAt();

        // Warp time to cover the timelock delay
        vm.warp(startTimestamp + minDelayPlusOne);

        // Execution reverts because the verification streak was disrupted in the
        // last 45 days
        vm.expectRevert(SurgeTimelockController.VerificationStreakDisrupted.selector);
        timelockController.execute(target, value, data, predecessor, salt);
    }

    function test_execute_passes_when_verification_streak_holds() external transactBy(Alice) {
        // Schedule a call to set the store to 1
        timelockController.schedule(target, value, data, predecessor, salt, minDelayPlusOne);

        // Warp time to cover the timelock delay
        vm.warp(startTimestamp + minDelayPlusOne);

        // Execution passes because the verification streak is held for the last 45 days
        timelockController.execute(target, value, data, predecessor, salt);

        // Assert that the store was set to 1
        assertEq(mockStore.getStore(), 1);
    }

    function test_executeBatch_fails_when_min_delay_is_not_met() external transactBy(Alice) {
        // Create arrays for batch operation
        address[] memory targets = new address[](2);
        uint256[] memory values = new uint256[](2);
        bytes[] memory payloads = new bytes[](2);

        targets[0] = target;
        targets[1] = target;
        values[0] = value;
        values[1] = value;
        payloads[0] = abi.encodeWithSelector(MockStore.setStore.selector, 1);
        payloads[1] = abi.encodeWithSelector(MockStore.setStore.selector, 2);

        // Schedule batch calls
        timelockController.scheduleBatch(
            targets, values, payloads, predecessor, salt, minDelayPlusOne
        );

        // Warp time to before the timelock delay is covered
        vm.warp(startTimestamp + minDelay - 1);

        // Execution reverts because the timelock delay is not covered
        vm.expectRevert("TimelockController: operation is not ready");
        timelockController.executeBatch(targets, values, payloads, predecessor, salt);
    }

    function test_executeBatch_reverts_when_verification_streak_is_disrupted()
        external
        transactBy(Alice)
    {
        // Create arrays for batch operation
        address[] memory targets = new address[](2);
        uint256[] memory values = new uint256[](2);
        bytes[] memory payloads = new bytes[](2);

        targets[0] = target;
        targets[1] = target;
        values[0] = value;
        values[1] = value;
        payloads[0] = abi.encodeWithSelector(MockStore.setStore.selector, 1);
        payloads[1] = abi.encodeWithSelector(MockStore.setStore.selector, 2);

        // Schedule batch calls
        timelockController.scheduleBatch(
            targets, values, payloads, predecessor, salt, minDelayPlusOne
        );

        // Warp time to somewhere between day 0 and day 45, to reset the verification streak
        vm.warp(startTimestamp + minDelayPlusOne / 2);
        mockTaikoInbox.resetVerificationStreakStartedAt();

        // Warp time to cover the timelock delay
        vm.warp(startTimestamp + minDelayPlusOne);

        // Execution reverts because the verification streak was disrupted in the
        // last 45 days
        vm.expectRevert(SurgeTimelockController.VerificationStreakDisrupted.selector);
        timelockController.executeBatch(targets, values, payloads, predecessor, salt);
    }

    function test_executeBatch_passes_when_verification_streak_holds() external transactBy(Alice) {
        // Create arrays for batch operation
        address[] memory targets = new address[](2);
        uint256[] memory values = new uint256[](2);
        bytes[] memory payloads = new bytes[](2);

        targets[0] = target;
        targets[1] = target;
        values[0] = value;
        values[1] = value;
        payloads[0] = abi.encodeWithSelector(MockStore.setStore.selector, 1);
        payloads[1] = abi.encodeWithSelector(MockStore.setStore.selector, 2);

        // Schedule batch calls
        timelockController.scheduleBatch(
            targets, values, payloads, predecessor, salt, minDelayPlusOne
        );

        // Warp time to cover the timelock delay
        vm.warp(startTimestamp + minDelayPlusOne);

        // Execution passes because the verification streak is held for the last 45 days
        timelockController.executeBatch(targets, values, payloads, predecessor, salt);

        // Assert that the store was set to 2 (last operation in batch)
        assertEq(mockStore.getStore(), 2);
    }
}
