// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/shared/common/LibStrings.sol";
import "../tiers/ITierProvider.sol";
import "../tiers/LibTiers.sol";

/// @title SurgeTierProviderBase
/// @notice This contract is a version of Taiko's TierProviderBase modified for Nethermind's Surge
/// @dev Modification include:
/// - Removed guardian tiers
/// - Only one proving tier i.e ZK + TEE
/// - No contestation for the proof
abstract contract SurgeTierProviderBase is ITierProvider {
    uint96 public constant BOND_UNIT = 0.04 ether;

    /// @inheritdoc ITierProvider
    function getTier(uint16 _tierId) public pure virtual returns (ITierProvider.Tier memory) {
        if (_tierId == LibTiers.TIER_ZKVM_AND_TEE) {
            // No contestation, but validity bond is still required to stake on the correctness of the 
            // initial state. This is returned during block verification. 
            return _buildTier(LibStrings.B_TIER_ZKVM_AND_TEE, 3, 0, 90);
        }

        revert TIER_NOT_FOUND();
    }

    /// @dev Builds a generic tier with specified parameters.
    /// @param _verifierName The name of the verifier.
    /// @param _validityBondUnits The units of validity bonds.
    /// @param _cooldownWindow The cooldown window duration in minutes.
    /// @param _provingWindow The proving window duration in minutes.
    /// @return A Tier struct with the provided parameters.
    function _buildTier(
        bytes32 _verifierName,
        uint8 _validityBondUnits,
        uint16 _cooldownWindow,
        uint16 _provingWindow
    )
        private
        pure
        returns (ITierProvider.Tier memory)
    {
        uint96 validityBond = BOND_UNIT * _validityBondUnits;
        return ITierProvider.Tier({
            verifierName: _verifierName,
            validityBond: validityBond,
            contestBond: validityBond / 10_000 * 65_625,
            cooldownWindow: _cooldownWindow,
            provingWindow: _provingWindow,
            maxBlocksToVerifyPerProof: 0
        });
    }
}
