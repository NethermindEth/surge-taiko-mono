// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/src/Script.sol";
import "forge-std/src/console2.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

import "src/shared/common/AddressManager.sol";

/// @title ChangeOwners
/// @notice This script is used to set the timelocked owner as the new owner of surge contracts
contract ChangeOwners is Script {
    uint256 public adminPrivateKey = vm.envUint("PRIVATE_KEY");
    address public newOwner = vm.envAddress("NEW_OWNER");
    
    function run() external {
        address[] memory contracts = vm.envAddress("CONTRACTS", ",");
        vm.startBroadcast(adminPrivateKey);
        for(uint i; i < contracts.length; ++i) {
            OwnableUpgradeable(contracts[i]).transferOwnership(newOwner);
        }
        vm.stopBroadcast();
    }
}