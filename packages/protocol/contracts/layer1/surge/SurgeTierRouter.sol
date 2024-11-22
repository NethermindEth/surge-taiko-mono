// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./SurgeTierProviderBase.sol";
import "../tiers/ITierRouter.sol";

/// @title SurgeTierRouter
contract SurgeTierRouter is SurgeTierProviderBase, ITierRouter {
    /// @inheritdoc ITierRouter
    function getProvider(uint256) external view returns (address) {
        return address(this);
    }

    /// @inheritdoc ITierProvider
    function getTierIds() external pure returns (uint16[] memory tiers_) {
        tiers_ = new uint16[](1);
        tiers_[0] = LibTiers.TIER_ZKVM_AND_TEE;
    }

    /// @inheritdoc ITierProvider
    function getMinTier(address, uint256) public pure override returns (uint16) {
        return LibTiers.TIER_ZKVM_AND_TEE;
    }
}
