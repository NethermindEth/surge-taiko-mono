// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/shared/common/LibStrings.sol";
import "src/layer1/tiers/ITierProvider.sol";
import "src/layer1/tiers/LibTiers.sol";
import "src/layer1/tiers/ITierRouter.sol";

/// @title TestTierProvider
/// @dev Labeled in AddressResolver as "tier_router"
/// @custom:security-contact security@taiko.xyz
contract TestTierProvider is ITierProvider, ITierRouter {
    uint256[50] private __gap;

    /// @inheritdoc ITierRouter
    function getProvider(uint256) external view returns (address) {
        return address(this);
    }
    /// @inheritdoc ITierProvider

    function getTier(uint16 _tierId) public pure override returns (ITierProvider.Tier memory) {
        if (_tierId == LibTiers.TIER_OPTIMISTIC) {
            return ITierProvider.Tier({
                verifierName: "",
                validityBond: 0.14 ether,
                contestBond: 0.28 ether,
                cooldownWindow: 1440, //24 hours
                provingWindow: 30, // 0.5 hours
                maxBlocksToVerifyPerProof: 0
            });
        }

        if (_tierId == LibTiers.TIER_SGX) {
            return ITierProvider.Tier({
                verifierName: LibStrings.B_TIER_SGX,
                validityBond: 0.14 ether,
                contestBond: 0.91875 ether, // =0.14 * 6.5625
                cooldownWindow: 1440, //24 hours
                provingWindow: 60, // 1 hours
                maxBlocksToVerifyPerProof: 0
            });
        }

        if (_tierId == LibTiers.TIER_GUARDIAN) {
            return ITierProvider.Tier({
                verifierName: LibStrings.B_TIER_GUARDIAN,
                validityBond: 0, // must be 0 for top tier
                contestBond: 0, // must be 0 for top tier
                cooldownWindow: 60, //1 hours
                provingWindow: 2880, // 48 hours
                maxBlocksToVerifyPerProof: 0
            });
        }

        revert TIER_NOT_FOUND();
    }

    /// @inheritdoc ITierProvider
    function getTierIds() public pure override returns (uint16[] memory tiers_) {
        tiers_ = new uint16[](3);
        tiers_[0] = LibTiers.TIER_OPTIMISTIC;
        tiers_[1] = LibTiers.TIER_SGX;
        tiers_[2] = LibTiers.TIER_GUARDIAN;
    }

    /// @inheritdoc ITierProvider
    function getMinTier(address, uint256 _rand) public pure override returns (uint16) {
        // 10% will be selected to require SGX proofs.
        if (_rand % 10 == 0) return LibTiers.TIER_SGX;
        // Other blocks are optimistic, without validity proofs.
        return LibTiers.TIER_OPTIMISTIC;
    }
}
