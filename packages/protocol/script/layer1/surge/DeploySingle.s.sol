// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;
import "forge-std/src/Script.sol";
import "forge-std/src/console2.sol";

import "src/layer1/surge/common/SurgeTierRouter.sol";

/// @title DeploySingle
/// @dev Helps with single contract deployment during upgrades
contract DeploySingle is Script {
    uint256 public adminPrivateKey = vm.envUint("PRIVATE_KEY");
    
    function run() external {
        vm.startBroadcast(adminPrivateKey);
        address surgeTierRouter = address(new SurgeTierRouter());
        console2.log("SurgeTierRouter deployed at:", surgeTierRouter);
        vm.stopBroadcast();
    }
}