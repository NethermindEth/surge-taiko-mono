// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {
    CrossChainSwapVaultL2
} from "../../../../contracts/layer2/surge/cross-chain-dex/CrossChainSwapVaultL2.sol";
import { SimpleDEX } from "../../../../contracts/layer2/surge/cross-chain-dex/SimpleDEX.sol";
import { SwapTokenL2 } from "../../../../contracts/layer2/surge/cross-chain-dex/SwapTokenL2.sol";
import { Script } from "forge-std/src/Script.sol";
import { console2 } from "forge-std/src/console2.sol";

/// @title DeployCrossChainDexL2
/// @notice Script to deploy the Cross-Chain DEX L2 contracts (vault-based, no mock minting)
contract DeployCrossChainDexL2 is Script {
    address internal immutable bridge = vm.envAddress("L2_BRIDGE");
    uint64 internal immutable l1ChainId = uint64(vm.envUint("L1_CHAIN_ID"));
    uint8 internal immutable tokenDecimals = uint8(vm.envOr("TOKEN_DECIMALS", uint256(6)));

    // Devnet-only EOAs pre-funded with bUSDC so end-users can swap immediately after
    // deployment. Matches the L1 deploy script's funding list. Only applied when
    // `block.chainid == DEVNET_L2_CHAIN_ID`.
    uint256 internal constant DEVNET_L2_CHAIN_ID = 763_374;
    address internal constant DEVNET_DEPLOYER = 0x3e95dFbBaF6B348396E6674C7871546dCC568e56;
    address internal constant DEVNET_OPERATOR = 0x5918b2e647464d4743601a865753e64C8059Dc4F;
    address internal constant DEVNET_SUBMITTER = 0x589A698b7b7dA0Bec545177D3963A2741105C7C9;

    modifier broadcast() {
        vm.startBroadcast();
        _;
        vm.stopBroadcast();
    }

    function run()
        external
        broadcast
        returns (address swapTokenL2_, address dex_, address l2Vault_)
    {
        address deployer = msg.sender;

        console2.log("=====================================");
        console2.log("Deploying Cross-Chain DEX L2 (Vault)");
        console2.log("=====================================");
        console2.log("Deployer:", deployer);
        console2.log("Bridge:", bridge);
        console2.log("L1 Chain ID:", l1ChainId);
        console2.log("");

        // Deploy SwapTokenL2 (bridged token) with deployer as initial minter
        // Decimals must match L1 token (configurable via TOKEN_DECIMALS env var)
        SwapTokenL2 swapTokenL2 =
            new SwapTokenL2("Bridged USDC", "bUSDC", deployer, 0, tokenDecimals);
        swapTokenL2_ = address(swapTokenL2);
        console2.log("SwapTokenL2 deployed at:", swapTokenL2_);

        // Deploy SimpleDEX
        SimpleDEX dex = new SimpleDEX(swapTokenL2_, deployer);
        dex_ = address(dex);
        console2.log("SimpleDEX deployed at:", dex_);

        // Deploy L2 Vault
        CrossChainSwapVaultL2 l2Vault =
            new CrossChainSwapVaultL2(bridge, l1ChainId, dex_, swapTokenL2_, deployer);
        l2Vault_ = address(l2Vault);
        console2.log("CrossChainSwapVaultL2 deployed at:", l2Vault_);

        // Pre-fund devnet EOAs with 100 bUSDC each (devnet L2 only) — must happen
        // before setMinter transfers authority away from the deployer.
        if (block.chainid == DEVNET_L2_CHAIN_ID) {
            uint256 devnetFundAmount = 100 * 10 ** uint256(tokenDecimals);
            address[3] memory devnetEoas = [DEVNET_DEPLOYER, DEVNET_OPERATOR, DEVNET_SUBMITTER];
            for (uint256 i = 0; i < devnetEoas.length; i++) {
                swapTokenL2.mint(devnetEoas[i], devnetFundAmount);
                console2.log("Funded devnet EOA with bUSDC:", devnetEoas[i]);
            }
        }

        // Transfer minting authority to L2 vault
        swapTokenL2.setMinter(l2Vault_);
        console2.log("Transferred minting authority to L2Vault");

        // Set L2 vault as authorized liquidity provider on DEX
        dex.setLiquidityProvider(l2Vault_);
        console2.log("Set L2Vault as DEX liquidity provider");

        // NO mock minting! NO initial liquidity!
        // Liquidity will be added from L1 via addLiquidityToL2()

        // Write deployment addresses
        writeJson("SwapTokenL2", swapTokenL2_);
        writeJson("SimpleDEX", dex_);
        writeJson("CrossChainSwapVaultL2", l2Vault_);

        console2.log("");
        console2.log("=====================================");
        console2.log("Deployment Complete");
        console2.log("=====================================");
        console2.log("");
        console2.log("Next steps:");
        console2.log("1. Set L2Vault on L1Vault (setup script on L1)");
        console2.log("2. Set L1Vault on L2Vault (setup script on L2)");
        console2.log("3. Add liquidity from L1 via addLiquidityToL2()");
    }

    function writeJson(string memory name, address addr) internal {
        vm.writeJson(
            vm.serializeAddress("deployment", name, addr),
            string.concat(vm.projectRoot(), "/deployments/cross-chain-dex-l2.json")
        );
    }
}
