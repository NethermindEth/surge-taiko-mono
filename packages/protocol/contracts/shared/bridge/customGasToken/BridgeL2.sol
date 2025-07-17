// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/shared/bridge/Bridge.sol";

/// @title BridgeL2
/// @notice See the documentation for {IBridge}.
/// @dev Labeled in address resolver as "bridge_l2".
/// @custom:security-contact security@nethermind.io
contract BridgeL2 is Bridge {
    // Surge: This contract is used *only* when the rollup uses a custom gas token.
    // This requires a complementary BridgeL1 contract on the L1 side.
    // Besides acting as the cross chain messaging contract, this contract holds the native gas
    // token on L2 side.

    constructor(
        address _resolver,
        address _signalService,
        address _quotaManager
    )
        Bridge(_resolver, _signalService, _quotaManager)
    { }
}
