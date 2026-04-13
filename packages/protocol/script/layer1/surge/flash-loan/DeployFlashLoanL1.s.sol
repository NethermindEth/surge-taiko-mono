// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {
    FlashLoanCallbackL1
} from "../../../../contracts/layer1/surge/flash-loan/FlashLoanCallbackL1.sol";
import { FlashLoanToken } from "../../../../contracts/layer1/surge/flash-loan/FlashLoanToken.sol";
import { Script } from "forge-std/src/Script.sol";
import { console2 } from "forge-std/src/console2.sol";

/// @title DeployFlashLoanL1
/// @notice Deploys the L1-side of the L2→L1→L2 synchronous flash loan demo:
///         a freely-mintable ERC20 (decimals 6) and the L1 callback.
contract DeployFlashLoanL1 is Script {
    address internal immutable bridge = vm.envAddress("L1_BRIDGE");
    uint64 internal immutable l2ChainId = uint64(vm.envUint("L2_CHAIN_ID"));
    uint8 internal immutable tokenDecimals = uint8(vm.envOr("TOKEN_DECIMALS", uint256(6)));

    modifier broadcast() {
        vm.startBroadcast();
        _;
        vm.stopBroadcast();
    }

    function run() external broadcast returns (address token_, address callback_) {
        address deployer = msg.sender;

        console2.log("=====================================");
        console2.log("Deploying Flash Loan L1");
        console2.log("=====================================");
        console2.log("Deployer:     ", deployer);
        console2.log("Bridge:       ", bridge);
        console2.log("L2 Chain ID:  ", l2ChainId);
        console2.log("Token decimals:", uint256(tokenDecimals));
        console2.log("");

        FlashLoanToken token = new FlashLoanToken("FlashUSD", "fUSD", tokenDecimals);
        token_ = address(token);
        console2.log("FlashLoanToken (L1):     ", token_);

        FlashLoanCallbackL1 cb = new FlashLoanCallbackL1(bridge, l2ChainId, token_, deployer);
        callback_ = address(cb);
        console2.log("FlashLoanCallbackL1:     ", callback_);

        writeJson("FlashLoanToken", token_);
        writeJson("FlashLoanCallbackL1", callback_);

        console2.log("");
        console2.log("Next:");
        console2.log("1. Deploy L2 side via script/layer2/surge/flash-loan/deploy_flash_loan_l2.sh");
        console2.log(
            "2. Link: setL2Executor(addr) on L1 callback and setL1Callback(addr) on L2 executor"
        );
    }

    function writeJson(string memory name, address addr) internal {
        vm.writeJson(
            vm.serializeAddress("deployment", name, addr),
            string.concat(vm.projectRoot(), "/deployments/flash-loan-l1.json")
        );
    }
}
