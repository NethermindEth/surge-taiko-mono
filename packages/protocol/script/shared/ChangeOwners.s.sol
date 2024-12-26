// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/src/Script.sol";
import "forge-std/src/console2.sol";
import "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

import "src/shared/common/AddressManager.sol";

contract ChangeOwners is Script {
    uint256 public adminPrivateKey = vm.envUint("PRIVATE_KEY");
    address public newOwner = vm.envAddress("NEW_OWNER");
    
    function run() external {
        address[] memory contracts = vm.envAddress("CONTRACTS", ",");
        vm.startBroadcast(adminPrivateKey);
        for(uint i; i < contracts.length; ++i) {
            OwnableUpgradeable(contracts[0]).transferOwnership(newOwner);
        }
        vm.stopBroadcast();
    }
}