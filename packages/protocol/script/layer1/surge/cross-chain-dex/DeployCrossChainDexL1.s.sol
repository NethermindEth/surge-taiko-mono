// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {
    CrossChainSwapVaultL1
} from "../../../../contracts/layer1/surge/cross-chain-dex/CrossChainSwapVaultL1.sol";
import { SimpleDEXL1 } from "../../../../contracts/layer1/surge/cross-chain-dex/SimpleDEXL1.sol";
import { SwapToken } from "../../../../contracts/layer1/surge/cross-chain-dex/SwapToken.sol";
import { WETH9Stub } from "../../../../contracts/layer1/surge/cross-chain-dex/WETH9Stub.sol";
import { Script } from "forge-std/src/Script.sol";
import { console2 } from "forge-std/src/console2.sol";

/// @title DeployCrossChainDexL1
/// @notice Deploys the L1 side of the cross-chain DEX (vault + optional test L1 DEX).
/// @dev Two L1 DEX modes:
///        - Test mode (default): if `L1_DEX_ROUTER` is unset, deploys `WETH9Stub` +
///          `SimpleDEXL1` and uses them as the L1 DEX. Optionally seeds the test DEX
///          with `L1_DEX_SEED_ETH` + `L1_DEX_SEED_TOKEN`.
///        - Live mode: if `L1_DEX_ROUTER` is set, uses it as-is. `L1_DEX_WETH` must
///          also be set in this case.
///      Always seeds the L1 vault inventory (env: `L1_VAULT_SEED_TOKEN`,
///      `L1_VAULT_SEED_ETH`) so that the L2→L1→L2 token→ETH direction works on day one.
contract DeployCrossChainDexL1 is Script {
    address internal immutable bridge = vm.envAddress("L1_BRIDGE");
    uint64 internal immutable l2ChainId = uint64(vm.envUint("L2_CHAIN_ID"));
    uint256 internal immutable initialTokenSupply = vm.envUint("INITIAL_TOKEN_SUPPLY");
    uint8 internal immutable tokenDecimals = uint8(vm.envOr("TOKEN_DECIMALS", uint256(18)));
    address internal immutable existingToken = vm.envOr("SWAP_TOKEN", address(0));

    // L1 DEX configuration
    address internal immutable existingL1Router = vm.envOr("L1_DEX_ROUTER", address(0));
    address internal immutable existingL1Weth = vm.envOr("L1_DEX_WETH", address(0));
    uint256 internal immutable l1DexSeedEth = vm.envOr("L1_DEX_SEED_ETH", uint256(0));
    uint256 internal immutable l1DexSeedToken = vm.envOr("L1_DEX_SEED_TOKEN", uint256(0));

    // L1 vault inventory seeding
    uint256 internal immutable l1VaultSeedToken = vm.envOr("L1_VAULT_SEED_TOKEN", uint256(0));
    uint256 internal immutable l1VaultSeedEth = vm.envOr("L1_VAULT_SEED_ETH", uint256(0));

    modifier broadcast() {
        vm.startBroadcast();
        _;
        vm.stopBroadcast();
    }

    function run()
        external
        broadcast
        returns (address swapToken_, address l1Vault_, address l1Router_, address weth_)
    {
        address deployer = msg.sender;

        console2.log("=====================================");
        console2.log("Deploying Cross-Chain DEX L1 (Vault)");
        console2.log("=====================================");
        console2.log("Deployer:", deployer);
        console2.log("Bridge:", bridge);
        console2.log("L2 Chain ID:", l2ChainId);
        console2.log("");

        // ---- Step 1: SwapToken ----
        if (existingToken != address(0)) {
            swapToken_ = existingToken;
            console2.log("Using existing SwapToken at:", swapToken_);
        } else {
            SwapToken swapToken = new SwapToken("USD Coin", "USDC", deployer, 0, tokenDecimals);
            swapToken_ = address(swapToken);
            console2.log("SwapToken deployed at:", swapToken_);
            console2.log("Token decimals:", uint256(tokenDecimals));

            swapToken.mint(deployer, initialTokenSupply);
            console2.log("Minted", initialTokenSupply, "tokens to deployer");
        }

        // ---- Step 2: L1 DEX (test mode auto-deploy, live mode use env) ----
        if (existingL1Router != address(0)) {
            require(existingL1Weth != address(0), "L1_DEX_WETH must be set with L1_DEX_ROUTER");
            l1Router_ = existingL1Router;
            weth_ = existingL1Weth;
            console2.log("Using existing L1 router at:", l1Router_);
            console2.log("Using existing WETH at:", weth_);
        } else {
            WETH9Stub weth9 = new WETH9Stub();
            weth_ = address(weth9);
            console2.log("WETH9Stub deployed at:", weth_);

            SimpleDEXL1 dexL1 = new SimpleDEXL1(swapToken_, weth_, deployer);
            l1Router_ = address(dexL1);
            console2.log("SimpleDEXL1 deployed at:", l1Router_);

            // Seed the test DEX if amounts are configured
            if (l1DexSeedEth > 0 && l1DexSeedToken > 0) {
                SwapToken(swapToken_).mint(deployer, l1DexSeedToken);
                SwapToken(swapToken_).approve(l1Router_, l1DexSeedToken);
                dexL1.seedLiquidity{ value: l1DexSeedEth }(l1DexSeedToken);
                console2.log(
                    "Seeded SimpleDEXL1 with ETH:", l1DexSeedEth, "tokens:", l1DexSeedToken
                );
            }
        }

        // ---- Step 3: L1 Vault ----
        CrossChainSwapVaultL1 l1Vault =
            new CrossChainSwapVaultL1(bridge, l2ChainId, swapToken_, l1Router_, weth_, deployer);
        l1Vault_ = address(l1Vault);
        console2.log("CrossChainSwapVaultL1 deployed at:", l1Vault_);

        // ---- Step 4: Seed L1 Vault inventory (USDC + ETH) ----
        // Needed so the L2→L1→L2 token→ETH path can draw on inventory before any
        // opposite-direction trade has funded it.
        if (l1VaultSeedToken > 0) {
            SwapToken(swapToken_).mint(deployer, l1VaultSeedToken);
            SwapToken(swapToken_).transfer(l1Vault_, l1VaultSeedToken);
            console2.log("Seeded L1 Vault with token amount:", l1VaultSeedToken);
        }
        if (l1VaultSeedEth > 0) {
            (bool ok,) = l1Vault_.call{ value: l1VaultSeedEth }("");
            require(ok, "ETH seed transfer failed");
            console2.log("Seeded L1 Vault with ETH:", l1VaultSeedEth);
        }

        // ---- Step 5: Write deployment artifacts ----
        writeJson("SwapToken", swapToken_);
        writeJson("CrossChainSwapVaultL1", l1Vault_);
        writeJson("L1Router", l1Router_);
        writeJson("WETH", weth_);

        console2.log("");
        console2.log("=====================================");
        console2.log("Deployment Complete");
        console2.log("=====================================");
        console2.log("");
        console2.log("Next steps:");
        console2.log("1. Deploy L2 contracts using deploy_cross_chain_dex.sh");
        console2.log("2. Set L2Vault on L1Vault (orchestration script)");
        console2.log("3. Set L1Vault on L2Vault (orchestration script)");
        console2.log("4. Add L2 DEX liquidity from L1 via addLiquidityToL2()");
    }

    function writeJson(string memory name, address addr) internal {
        vm.writeJson(
            vm.serializeAddress("deployment", name, addr),
            string.concat(vm.projectRoot(), "/deployments/cross-chain-dex-l1.json")
        );
    }
}
