// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./TierProviderBase.sol";
import "./ITierRouter.sol";

/// @title TierProviderV2
/// @custom:security-contact security@taiko.xyz
contract TierProviderV2 is TierProviderBase, ITierRouter {
    /// @inheritdoc ITierRouter
    function getProvider(uint256) external view returns (address) {
        return address(this);
    }

    /// @inheritdoc ITierProvider
    function getTierIds() public pure override returns (uint16[] memory tiers_) {
        tiers_ = new uint16[](1);
        tiers_[0] = LibTiers.TIER_TWO_OF_THREE;
    }

    /// @inheritdoc ITierProvider
    function getMinTier(address, uint256) public pure override returns (uint16) {
        return LibTiers.TIER_TWO_OF_THREE;
    }
}
