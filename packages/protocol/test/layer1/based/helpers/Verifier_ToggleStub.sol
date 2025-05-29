// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/layer1/surge/verifiers/ISurgeVerifier.sol";
import "src/layer1/surge/verifiers/LibProofType.sol";

// Surge: change the contract to ISurgeVerifier
contract Verifier_ToggleStub is ISurgeVerifier {
    LibProofType.ProofType public proofType;
    LibProofType.ProofType public proofTypeToUpgrade;

    constructor() {
        proofType = LibProofType.ProofType.SGX_SP1;
        proofTypeToUpgrade = LibProofType.ProofType.INVALID;
    }

    function setProofType(LibProofType.ProofType _proofType) external {
        proofType = _proofType;
    }

    function verifyProof(
        Context[] calldata,
        bytes calldata
    )
        external
        view
        returns (LibProofType.ProofType)
    {
        return proofType;
    }

    function markUpgradeable(LibProofType.ProofType _proofType) external {
        proofTypeToUpgrade = _proofType;
    }

    function upgradeVerifier(LibProofType.ProofType _proofType, address _newVerifier) external { }
}
