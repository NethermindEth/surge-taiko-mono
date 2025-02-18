// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "../based/ITaikoL1.sol";
import "@openzeppelin/contracts/governance/TimelockController.sol";

/// @title SurgeTimelockController
/// @dev Staisfies stage-2 rollup requirements by blocking executions if a block 
/// has not been verified in a while. 
/// @custom:security-contact security@nethermind.io
contract SurgeTimelockedController is TimelockController {
    /// @notice Address of taiko's inbox contract
    ITaikoL1 internal taikoL1;

    /// @notice Minimum period for which the liveness disruption must not have exceeded the 
    /// `maxLivenessDisruption` period in the L1 config
    uint64 internal minLivenessStreak;

    error ALREADY_INITIALIZED();
    error LIVENESS_DISRUPTED();
    
    constructor(uint64 _minLivenessStreak, uint256 minDelay, address[] memory proposers, address[] memory executors, address admin) TimelockController(minDelay, proposers, executors, admin) {
        minLivenessStreak = _minLivenessStreak;
    }

    function init(address _taikoL1) external {
        if(address(taikoL1) != address(0)) {
            revert ALREADY_INITIALIZED();
        }
        taikoL1 = ITaikoL1(_taikoL1);
    }

    function execute(
        address target,
        uint256 value,
        bytes calldata payload,
        bytes32 predecessor,
        bytes32 salt
    ) public payable override onlyRoleOrOpenRole(EXECUTOR_ROLE) {
        if(_isLivenessDisrupted()) { 
            revert LIVENESS_DISRUPTED();
        }

        super.execute(target, value, payload, predecessor, salt);
    }

    function executeBatch(
        address[] calldata targets,
        uint256[] calldata values,
        bytes[] calldata payloads,
        bytes32 predecessor,
        bytes32 salt
    ) public payable override onlyRoleOrOpenRole(EXECUTOR_ROLE) {
        if(_isLivenessDisrupted()) { 
            revert LIVENESS_DISRUPTED();
        }

        super.executeBatch(targets, values, payloads, predecessor, salt);
    }

    function updateMinLivenessStreak(uint64 _minLivenessStreak) external onlyRole(TIMELOCK_ADMIN_ROLE) {
        minLivenessStreak = _minLivenessStreak;
    }

    /// @dev Returns `true` if an L2 block has not been proposed + verified in a gap of greater 
    // than `Config.maxLivenessDisruptionPeriod` seconds within the last `minLivenessStreak`
    function _isLivenessDisrupted() internal view returns(bool) {
        uint256 verificationStreakStartedAt = taikoL1.getVerificationStreakStartAt();
        return (block.timestamp - verificationStreakStartedAt) < minLivenessStreak;
    }
}