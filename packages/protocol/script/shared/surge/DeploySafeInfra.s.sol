// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "forge-std/src/Script.sol";
import "forge-std/src/console2.sol";

import { SafeL2 } from "@safe/SafeL2.sol";
import { SafeProxyFactory } from "@safe/proxies/SafeProxyFactory.sol";
import { MultiSend } from "@safe/libraries/MultiSend.sol";
import { MultiSendCallOnly } from "@safe/libraries/MultiSendCallOnly.sol";
import { CompatibilityFallbackHandler } from "@safe/handler/CompatibilityFallbackHandler.sol";

/// @title DeploySafeInfra
/// @notice Deploys the core Safe infrastructure contracts (singleton, factory, and utility
///         libraries) needed to instantiate Gnosis Safe wallets on any EVM chain. Deploy on both
///         L1 and L2 with the same nonce so that addresses match across chains.
contract DeploySafeInfra is Script {
    /// @notice Deploys SafeL2 singleton, SafeProxyFactory, MultiSend, MultiSendCallOnly, and
    ///         CompatibilityFallbackHandler and logs their addresses to the console.
    function run() external {
        vm.startBroadcast();

        SafeL2 singleton = new SafeL2();
        console2.log("SafeL2 singleton:", address(singleton));

        SafeProxyFactory factory = new SafeProxyFactory();
        console2.log("SafeProxyFactory:", address(factory));

        MultiSend multiSend = new MultiSend();
        console2.log("MultiSend:", address(multiSend));

        MultiSendCallOnly multiSendCallOnly = new MultiSendCallOnly();
        console2.log("MultiSendCallOnly:", address(multiSendCallOnly));

        CompatibilityFallbackHandler fallbackHandler = new CompatibilityFallbackHandler();
        console2.log("FallbackHandler:", address(fallbackHandler));

        vm.stopBroadcast();
    }
}
