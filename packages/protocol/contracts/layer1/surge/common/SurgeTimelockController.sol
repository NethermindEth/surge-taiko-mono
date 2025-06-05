// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "../verifiers/ISurgeVerifier.sol";
import "../verifiers/LibProofType.sol";
import "src/layer1/based/ITaikoInbox.sol";
import "@openzeppelin/contracts/governance/TimelockController.sol";

/// @title SurgeTimelockController
/// @dev Satisfies stage-2 rollup requirements by blocking executions if a block
/// has not been verified in a while.
/// @custom:security-contact security@nethermind.io
contract SurgeTimelockController is TimelockController {
    address public taikoInbox; // Slot 1
    address public verifier; // Slot 2

    /// @notice Minimum period for which the verification streak must not have been disrupted
    uint256 public minVerificationStreak; // Slot 3

    uint256[45] private _gap;

    error AlreadyInitialized();
    error VerificationStreakDisrupted();

    constructor(
        uint64 _minVerificationStreak,
        uint256 _minDelay,
        address[] memory _proposers,
        address[] memory _executors,
        address _admin
    )
        TimelockController(_minDelay, _proposers, _executors, _admin)
    {
        minVerificationStreak = _minVerificationStreak;
    }

    function init(address _taikoInbox, address _verifier) external {
        if (taikoInbox != address(0)) {
            revert AlreadyInitialized();
        }
        taikoInbox = _taikoInbox;
        verifier = _verifier;
    }

    function execute(
        address _target,
        uint256 _value,
        bytes calldata _payload,
        bytes32 _predecessor,
        bytes32 _salt
    )
        public
        payable
        override
        onlyRoleOrOpenRole(EXECUTOR_ROLE)
    {
        if (_isVerificationStreakDisrupted()) {
            revert VerificationStreakDisrupted();
        }

        super.execute(_target, _value, _payload, _predecessor, _salt);
    }

    function executeBatch(
        address[] calldata _targets,
        uint256[] calldata _values,
        bytes[] calldata _payloads,
        bytes32 _predecessor,
        bytes32 _salt
    )
        public
        payable
        override
        onlyRoleOrOpenRole(EXECUTOR_ROLE)
    {
        if (_isVerificationStreakDisrupted()) {
            revert VerificationStreakDisrupted();
        }

        super.executeBatch(_targets, _values, _payloads, _predecessor, _salt);
    }

    // Timelock bypass functions
    // --------------------------------------------------------------------------------------------

    /// @dev Can be used to bypass the timelock when the verifier needs an instant upgrade.
    /// @dev Only the proposer role can call this function, which in the case of Surge is the
    /// primary owner multisig.
    function executeVerifierUpgrade(
        LibProofType.ProofType _proofType,
        address _newVerifier
    )
        external
        onlyRole(PROPOSER_ROLE)
    {
        ISurgeVerifier(verifier).upgradeVerifier(_proofType, _newVerifier);
    }

    // Timelocked functions
    // --------------------------------------------------------------------------------------------

    function updateMinVerificationStreak(uint64 _minVerificationStreak)
        external
        onlyRole(TIMELOCK_ADMIN_ROLE)
    {
        minVerificationStreak = _minVerificationStreak;
    }

    function updateVerifierAddress(address _newVerifier) external onlyRole(TIMELOCK_ADMIN_ROLE) {
        verifier = _newVerifier;
    }

    function updateTaikoInboxAddress(address _newTaikoInbox)
        external
        onlyRole(TIMELOCK_ADMIN_ROLE)
    {
        taikoInbox = _newTaikoInbox;
    }

    // Internal functions
    // --------------------------------------------------------------------------------------------

    /// @dev Returns `true` if an L2 block has not been proposed & verified in a gap of greater
    ///      than `Config.maxVerificationDelay` seconds within the last `minVerificationStreak`
    function _isVerificationStreakDisrupted() internal view returns (bool) {
        uint256 verificationStreakStartedAt =
            ITaikoInbox(taikoInbox).getVerificationStreakStartedAt();
        return (block.timestamp - verificationStreakStartedAt) < minVerificationStreak;
    }
}
