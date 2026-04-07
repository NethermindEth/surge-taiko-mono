// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { ProofVerifierDummy } from "../common/ProofVerifierDummy.sol";
import {
    ControlID,
    RiscZeroGroth16Verifier
} from "@risc0/contracts/groth16/RiscZeroGroth16Verifier.sol";
import { SP1Verifier as SuccinctVerifier } from "@sp1-contracts/src/v5.0.0/SP1VerifierPlonk.sol";
import { console2 } from "forge-std/src/console2.sol";
import { SurgeVerifier } from "src/layer1/surge/SurgeVerifier.sol";
import { LibProofBitmap } from "src/layer1/surge/libs/LibProofBitmap.sol";
import { Risc0Verifier } from "src/layer1/verifiers/Risc0Verifier.sol";
import { SP1Verifier } from "src/layer1/verifiers/SP1Verifier.sol";
import { DeployCapability } from "test/shared/DeployCapability.sol";

/// @title DeployVerifiers
/// @notice Step 2: Deploys internal verifiers and configures SurgeVerifier.
/// @dev Requires SURGE_VERIFIER and EFFECTIVE_OWNER from Step 1.
///      Must use the same PRIVATE_KEY as Step 1 (SurgeVerifier owner = Step 1 deployer).
/// @custom:security-contact security@nethermind.io
contract DeployVerifiers is DeployCapability {
    // Signer configuration
    uint256 internal immutable privateKey = vm.envUint("PRIVATE_KEY");

    // L2 configuration
    uint64 internal immutable l2ChainId = uint64(vm.envUint("L2_CHAIN_ID"));

    // Verifier configuration
    bool internal immutable useDummyVerifier = vm.envBool("USE_DUMMY_VERIFIER");
    address internal immutable dummyVerifierSigner = vm.envAddress("DUMMY_VERIFIER_SIGNER");
    bool internal immutable deployRisc0RethVerifier = vm.envBool("DEPLOY_RISC0_RETH_VERIFIER");
    bool internal immutable deploySp1RethVerifier = vm.envBool("DEPLOY_SP1_RETH_VERIFIER");

    // Addresses from Step 1
    address internal immutable surgeVerifier = vm.envAddress("SURGE_VERIFIER");
    address internal immutable effectiveOwner = vm.envAddress("EFFECTIVE_OWNER");

    struct VerifierContracts {
        address risc0RethVerifier;
        address sp1RethVerifier;
    }

    modifier broadcast() {
        require(privateKey != 0, "invalid private key");
        vm.startBroadcast(privateKey);
        _;
        vm.stopBroadcast();
    }

    function run() external broadcast {
        require(surgeVerifier != address(0), "config: SURGE_VERIFIER");
        require(effectiveOwner != address(0), "config: EFFECTIVE_OWNER");

        console2.log("=====================================");
        console2.log("Step 2: Deploy Verifiers");
        console2.log("=====================================");
        console2.log("** SurgeVerifier:", surgeVerifier);
        console2.log("** Effective owner:", effectiveOwner);

        // Deploy internal verifiers
        VerifierContracts memory verifiers = _deployInternalVerifiers();

        // Configure SurgeVerifier with internal verifiers
        _setupProofVerifier(verifiers);

        console2.log("=====================================");
        console2.log("Step 2 Complete");
        console2.log("=====================================");
    }

    /// @dev Deploys internal verifiers. The deployer is the initial owner of real verifiers.
    ///      Ownership is transferred to the effective owner and must be accepted.
    function _deployInternalVerifiers() internal returns (VerifierContracts memory verifiers) {
        if (useDummyVerifier) {
            require(dummyVerifierSigner != address(0), "config: DUMMY_VERIFIER_SIGNER");

            ProofVerifierDummy dummyVerifier = new ProofVerifierDummy(dummyVerifierSigner);
            address dummyAddr = address(dummyVerifier);
            writeJson("proof_verifier_dummy", dummyAddr);
            console2.log("** Deployed ProofVerifierDummy:", dummyAddr);
            console2.log("** ProofVerifierDummy signer:", dummyVerifierSigner);

            if (deployRisc0RethVerifier) {
                verifiers.risc0RethVerifier = dummyAddr;
                console2.log("** Using ProofVerifierDummy for RISC0");
            }
            if (deploySp1RethVerifier) {
                verifiers.sp1RethVerifier = dummyAddr;
                console2.log("** Using ProofVerifierDummy for SP1");
            }
        } else {
            if (deployRisc0RethVerifier) {
                RiscZeroGroth16Verifier verifier =
                    new RiscZeroGroth16Verifier(ControlID.CONTROL_ROOT, ControlID.BN254_CONTROL_ID);
                writeJson("risc0_groth16_verifier", address(verifier));
                console2.log("** Deployed Risc0 groth16 verifier:", address(verifier));

                Risc0Verifier risc0Verifier =
                    new Risc0Verifier(l2ChainId, address(verifier), msg.sender);
                verifiers.risc0RethVerifier = address(risc0Verifier);
                writeJson("risc0_verifier", address(risc0Verifier));
                console2.log("** Deployed Risc0 verifier:", address(risc0Verifier));

                risc0Verifier.transferOwnership(effectiveOwner);
                console2.log("** Risc0 verifier ownership transfer initiated to:", effectiveOwner);
            }

            if (deploySp1RethVerifier) {
                SuccinctVerifier succinctVerifier = new SuccinctVerifier();
                writeJson("succinct_verifier", address(succinctVerifier));
                console2.log("** Deployed Succinct verifier:", address(succinctVerifier));

                SP1Verifier sp1Verifier =
                    new SP1Verifier(l2ChainId, address(succinctVerifier), msg.sender);
                verifiers.sp1RethVerifier = address(sp1Verifier);
                writeJson("sp1_verifier", address(sp1Verifier));
                console2.log("** Deployed SP1 verifier:", address(sp1Verifier));

                sp1Verifier.transferOwnership(effectiveOwner);
                console2.log("** SP1 verifier ownership transfer initiated to:", effectiveOwner);
            }
        }
    }

    function _setupProofVerifier(VerifierContracts memory _verifiers) internal {
        SurgeVerifier proofVerifier = SurgeVerifier(surgeVerifier);

        if (_verifiers.risc0RethVerifier != address(0)) {
            proofVerifier.setVerifier(
                LibProofBitmap.ProofBitmap.wrap(proofVerifier.RISC0_RETH()),
                _verifiers.risc0RethVerifier
            );
            console2.log("** Set RISC0 verifier:", _verifiers.risc0RethVerifier);
        }

        if (_verifiers.sp1RethVerifier != address(0)) {
            proofVerifier.setVerifier(
                LibProofBitmap.ProofBitmap.wrap(proofVerifier.SP1_RETH()),
                _verifiers.sp1RethVerifier
            );
            console2.log("** Set SP1 verifier:", _verifiers.sp1RethVerifier);
        }

        // Transfer SurgeVerifier ownership (requires acceptance)
        proofVerifier.transferOwnership(effectiveOwner);
        console2.log("** SurgeVerifier ownership transfer initiated to:", effectiveOwner);
    }

    function writeJson(string memory name, address addr) internal {
        vm.writeJson(
            vm.serializeAddress("deploy_verifiers", name, addr),
            string.concat(vm.projectRoot(), "/deployments/deploy_verifiers.json")
        );
    }
}
