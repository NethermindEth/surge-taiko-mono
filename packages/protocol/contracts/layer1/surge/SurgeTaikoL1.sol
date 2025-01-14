// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "../based/TaikoL1.sol";

/// @title SurgeTaikoL1
/// @dev Labeled in AddressResolver as "taiko"
/// @custom:security-contact security@nethermind.io
contract SurgeTaikoL1 is TaikoL1 {
    uint64 private immutable chainId;

    constructor(uint64 _chainId) {
        chainId = _chainId;
    }

    /// @inheritdoc ITaikoL1
    function getConfig() public view override returns (TaikoData.Config memory) {
        return TaikoData.Config({
            chainId: chainId,
            blockMaxProposals: 324_000,
            blockRingBufferSize: 360_000,
            maxBlocksToVerify: 16,
            blockMaxGasLimit: 240_000_000,
            livenessBond: 0.07 ether,
            stateRootSyncInternal: 16,
            maxAnchorHeightOffset: 64,
            baseFeeConfig: LibSharedData.BaseFeeConfig({
                adjustmentQuotient: 8,
                sharingPctg: 75,
                gasIssuancePerSecond: 5_000_000,
                minGasExcess: 1_340_000_000,
                maxGasIssuancePerBlock: 600_000_000
            }),
            ontakeForkHeight: 1
        });
    }
}