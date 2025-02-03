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

    error ALREADY_INITIALIZED();
    error LIVENESS_DISRUPTED();
    
    constructor(uint256 minDelay, address[] memory proposers, address[] memory executors, address admin) TimelockController(minDelay, proposers, executors, admin) {}

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

    /// @dev Returns `true` if an L2 block has not been verified in the last 1 day
    function _isLivenessDisrupted() internal view returns(bool) {
        uint256 lastVerificationTimestamp = taikoL1.getLastVerificationTimestamp();
        return (block.timestamp - lastVerificationTimestamp) > 1 days;
    }
}