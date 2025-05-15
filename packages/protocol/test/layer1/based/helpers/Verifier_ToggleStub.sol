// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/layer1/surge/verifiers/ISurgeVerifier.sol";
import "src/layer1/based/ITaikoInbox.sol";

// Surge: change the contract to ISurgeVerifier
contract Verifier_ToggleStub is ISurgeVerifier {
    bool private shouldFail;
    ITaikoInbox.ProofType private proofType;

    constructor() {
        proofType = ITaikoInbox.ProofType.ZK_TEE;
    }

    function makeVerifierToFail() external {
        shouldFail = true;
    }

    function makeVerifierToSucceed() external {
        shouldFail = false;
    }

    function setProofType(ITaikoInbox.ProofType _proofType) external {
        proofType = _proofType;
    }

    function verifyProof(
        Context[] calldata,
        bytes calldata
    )
        external
        view
        returns (ITaikoInbox.ProofType)
    {
        require(!shouldFail, "ISurgeVerifier failure");
        return proofType;
    }
}
