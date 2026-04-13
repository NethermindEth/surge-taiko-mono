// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { FlashLoanToken } from "../../../../contracts/layer1/surge/flash-loan/FlashLoanToken.sol";
import {
    FlashLoanExecutorL2
} from "../../../../contracts/layer2/surge/flash-loan/FlashLoanExecutorL2.sol";
import { FlashLoanPool } from "../../../../contracts/layer2/surge/flash-loan/FlashLoanPool.sol";
import { Script } from "forge-std/src/Script.sol";
import { console2 } from "forge-std/src/console2.sol";

/// @title DeployFlashLoanL2
/// @notice Deploys the L2-side of the L2→L1→L2 synchronous flash loan demo:
///         a freely-mintable ERC20 (decimals 6), the flash loan pool, and the
///         L2 executor. Mints an initial supply into the pool.
contract DeployFlashLoanL2 is Script {
    address internal immutable bridge = vm.envAddress("L2_BRIDGE");
    uint64 internal immutable l1ChainId = uint64(vm.envUint("L1_CHAIN_ID"));
    uint8 internal immutable tokenDecimals = uint8(vm.envOr("TOKEN_DECIMALS", uint256(6)));
    uint256 internal immutable initialPoolLiquidity =
        vm.envOr("INITIAL_POOL_LIQUIDITY", uint256(1_000_000 * 10 ** 6));

    modifier broadcast() {
        vm.startBroadcast();
        _;
        vm.stopBroadcast();
    }

    function run() external broadcast returns (address token_, address pool_, address executor_) {
        address deployer = msg.sender;

        console2.log("=====================================");
        console2.log("Deploying Flash Loan L2");
        console2.log("=====================================");
        console2.log("Deployer:         ", deployer);
        console2.log("Bridge:           ", bridge);
        console2.log("L1 Chain ID:      ", l1ChainId);
        console2.log("Token decimals:   ", uint256(tokenDecimals));
        console2.log("Pool liquidity:   ", initialPoolLiquidity);
        console2.log("");

        FlashLoanToken token = new FlashLoanToken("FlashUSD", "fUSD", tokenDecimals);
        token_ = address(token);
        console2.log("FlashLoanToken (L2):    ", token_);

        FlashLoanPool pool = new FlashLoanPool(token_);
        pool_ = address(pool);
        console2.log("FlashLoanPool:          ", pool_);

        FlashLoanExecutorL2 executor =
            new FlashLoanExecutorL2(bridge, l1ChainId, token_, pool_, deployer);
        executor_ = address(executor);
        console2.log("FlashLoanExecutorL2:    ", executor_);

        // Seed the pool with initial liquidity (freely minted since the token is open).
        token.mint(pool_, initialPoolLiquidity);
        console2.log("Seeded pool with        ", initialPoolLiquidity);

        writeJson("FlashLoanToken", token_);
        writeJson("FlashLoanPool", pool_);
        writeJson("FlashLoanExecutorL2", executor_);

        console2.log("");
        console2.log("Next:");
        console2.log(
            "1. Link: setL2Executor(addr) on L1 callback and setL1Callback(addr) on L2 executor"
        );
    }

    function writeJson(string memory name, address addr) internal {
        vm.writeJson(
            vm.serializeAddress("deployment", name, addr),
            string.concat(vm.projectRoot(), "/deployments/flash-loan-l2.json")
        );
    }
}
