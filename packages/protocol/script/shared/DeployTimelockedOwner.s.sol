// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/governance/TimelockController.sol";
import "test/shared/DeployCapability.sol";

contract DeployTimelockedOwner is DeployCapability {
    uint256 public deployerPrivKey = vm.envUint("PRIVATE_KEY");
    function run() external {
        require(deployerPrivKey != 0, "invalid deployer priv key");
        vm.startBroadcast(deployerPrivKey);
        address[] memory executors = vm.envAddress("OWNER_MULTISIG_SIGNERS", ",");
        
        address ownerMultisig = vm.envAddress("OWNER_MULTISIG");
        addressNotNull(ownerMultisig, "ownerMultisig");
        
        address[] memory proposers = new address[](1);
        proposers[0] = ownerMultisig; 
        
        // Setup timelock controller with 45 day (86400 seconds * 45) delay
        address timelockController = address(
            new TimelockController(86400 * 45, proposers, executors, address(0))
        );
        console2.log("Timelocked owner: ", timelockController);
        vm.stopBroadcast();
    }
    
    function addressNotNull(address addr, string memory err) private pure {
        require(addr != address(0), err);
    }
}