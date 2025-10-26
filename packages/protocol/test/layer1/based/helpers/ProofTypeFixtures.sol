// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "src/layer1/surge/verifiers/LibProofType.sol";

contract ProofTypeFixtures {
    using LibProofType for LibProofType.ProofType;

    // Fuzz testing input fixtures
    // ---------------------------

    LibProofType.ProofType[] internal zkTeeProofTypes = [
        LibProofType.tdxNethermind().combine(LibProofType.sp1Reth()),
        LibProofType.tdxNethermind().combine(LibProofType.risc0Reth()),
        LibProofType.azureTdxNethermind().combine(LibProofType.sp1Reth()),
        LibProofType.azureTdxNethermind().combine(LibProofType.risc0Reth()),
        LibProofType.sgxReth().combine(LibProofType.sp1Reth()),
        LibProofType.sgxReth().combine(LibProofType.risc0Reth()),
        LibProofType.sgxGeth().combine(LibProofType.sp1Reth()),
        LibProofType.sgxGeth().combine(LibProofType.risc0Reth())
    ];

    LibProofType.ProofType[] internal zkProofTypes =
        [LibProofType.sp1Reth(), LibProofType.risc0Reth()];

    LibProofType.ProofType[] internal teeProofTypes = [
        LibProofType.sgxReth(),
        LibProofType.sgxGeth(),
        LibProofType.tdxNethermind(),
        LibProofType.azureTdxNethermind()
    ];

    LibProofType.ProofType[] internal allProofTypes = [
        LibProofType.sgxReth(),
        LibProofType.sgxGeth(),
        LibProofType.tdxNethermind(),
        LibProofType.azureTdxNethermind(),
        LibProofType.sp1Reth(),
        LibProofType.risc0Reth()
    ];
}
