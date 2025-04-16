// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "../based/TaikoL1.sol";

/// @title SurgeTaikoL1
/// @dev Labeled in AddressResolver as "taiko"
/// @custom:security-contact security@nethermind.io
contract SurgeDevnetTaikoL1 is TaikoL1 {
    uint64 private immutable chainId;
    uint64 private immutable maxLivenessDisruptionPeriod;

    constructor(uint64 _chainId, uint64 _maxLivenessDisruptionPeriod) {
        chainId = _chainId;
        maxLivenessDisruptionPeriod = _maxLivenessDisruptionPeriod;
    }

    /// @inheritdoc ITaikoL1
    function getConfig() public view override returns (TaikoData.Config memory) {
        return TaikoData.Config({
            chainId: chainId,
            blockMaxProposals: 324_000,
            blockRingBufferSize: 360_000,
            maxBlocksToVerify: 4,
            blockMaxGasLimit: 600_000_000,
            livenessBond: 0.07 ether,
            stateRootSyncInternal: 4,
            maxAnchorHeightOffset: 64,
            baseFeeConfig: LibSharedData.BaseFeeConfig({
                adjustmentQuotient: 8,
                sharingPctg: 0,
                gasIssuancePerSecond: 100_000_000, // The network targets 1 block / 3 seconds
                minGasExcess: 31_136_000_000, // Resolves to ~0.09992 Gwei
                maxGasIssuancePerBlock: 6_000_000_000
            }),
            ontakeForkHeight: 1,
            maxLivenessDisruptionPeriod: maxLivenessDisruptionPeriod
        });
    }
}