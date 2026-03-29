// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { UserOpsSubmitter } from "./UserOpsSubmitter.sol";

contract UserOpsSubmitterFactory {
    /// @notice Mapping from owner EOA to their submitter contract
    mapping(address owner => address submitter) public submitters;

    event SubmitterCreated(
        address indexed submitter, address indexed owner, address indexed deployer
    );

    /// @notice Create a new UserOpsSubmitter for the given owner using CREATE2.
    /// @dev Uses keccak256(owner) as the salt, so the same owner always gets the
    ///      same submitter address on every chain where this factory is deployed
    ///      at the same address.
    /// @param _owner The EOA that will own the submitter
    /// @return submitter_ The address of the created submitter
    function createSubmitter(address _owner) external returns (address submitter_) {
        if (submitters[_owner] != address(0)) revert SUBMITTER_EXISTS();

        bytes32 salt = keccak256(abi.encodePacked(_owner));
        UserOpsSubmitter submitter = new UserOpsSubmitter{salt: salt}(_owner);
        submitter_ = address(submitter);

        submitters[_owner] = submitter_;

        emit SubmitterCreated(submitter_, _owner, msg.sender);
    }

    /// @notice Get the submitter for a given owner
    /// @param _owner The EOA to look up
    /// @return submitter_ The submitter address, or address(0) if none exists
    function getSubmitter(address _owner) external view returns (address submitter_) {
        return submitters[_owner];
    }

    /// @notice Predict the CREATE2 address for a given owner (before deployment)
    /// @param _owner The EOA that will own the submitter
    /// @return predicted_ The address the submitter will be deployed to
    function predictSubmitterAddress(address _owner) external view returns (address predicted_) {
        bytes32 salt = keccak256(abi.encodePacked(_owner));
        bytes32 hash = keccak256(
            abi.encodePacked(
                bytes1(0xff),
                address(this),
                salt,
                keccak256(
                    abi.encodePacked(type(UserOpsSubmitter).creationCode, abi.encode(_owner))
                )
            )
        );
        predicted_ = address(uint160(uint256(hash)));
    }

    error SUBMITTER_EXISTS();
}
