# Surge Protocol Deployment Guide

This document describes the deployment sequence for the Surge protocol on both L1 and L2.

> **Prerequisite**: It is assumed that the genesis/chainspec file for the L2 network has already been generated.

---

## Deployment Overview

| Step | Script                            | Network | Description                                         |
| ---- | --------------------------------- | ------- | --------------------------------------------------- |
| 1    | `DeploySurgeL1.s.sol`             | L1      | Deploy all L1 contracts                             |
| 2    | Verifier setup scripts            | L1      | Configure prover image IDs                          |
| 3    | `AcceptOwnership.s.sol`           | L1      | Accept pending ownership transfers                  |
| 4    | `SetupSurgeL2.s.sol`              | L2      | Register L1 contracts and setup delegate controller |
| 5    | Bridge-based ownership acceptance | L1→L2   | Accept L2 ownership via delegate controller         |

---

## Step 1: Deploy L1 Contracts

**Script**: `script/layer1/surge/DeploySurgeL1.s.sol`  
**Shell wrapper**: `script/layer1/surge/deploy_surge_l1.sh`

### What it deploys

#### Rollup Contracts

- **Inbox** (proxy) - Main rollup contract for propose, prove, and finalization
- **Proof Verifier** (`SurgeVerifier` or `SurgeVerifierDummy` if `USE_DUMMY_VERIFIER=true`)

#### Shared Contracts

- **SharedResolver** - Cross-contract discovery
- **SignalService** - Cross-chain signal relay
- **Bridge** - Cross-chain messaging
- **ERC20Vault** - ERC20 token bridging
- **ERC721Vault** - ERC721 token bridging
- **ERC1155Vault** - ERC1155 token bridging
- **BridgedERC20/721/1155** - Bridged token implementations (clone pattern)

#### Preconf Contracts

- **PreconfWhitelist** - Whitelisted preconfirmation operators store

#### Internal Verifiers (optional)

- **Risc0Verifier** (if `DEPLOY_RISC0_RETH_VERIFIER=true`)
- **SP1Verifier** (if `DEPLOY_SP1_RETH_VERIFIER=true`)

### Ownership Configuration

The `CONTRACT_OWNER` environment variable specifies the intended owner of all contracts (typically a timelocked Security Council, DAO, or an EOA on devnet).

#### Contracts with immediate ownership (`owner = CONTRACT_OWNER`)

These contracts have their ownership set directly during deployment:

- SignalService
- Bridge
- ERC20Vault
- ERC721Vault
- ERC1155Vault
- PreconfWhitelist

#### Contracts with pending ownership (`pendingOwner = CONTRACT_OWNER`)

These contracts use the 2-step ownership transfer pattern and require manual acceptance:

- **Proof Verifier** (`SurgeVerifier` / `SurgeVerifierDummy`)
- **Inbox** (SurgeInbox proxy)
- **SharedResolver**
- **Risc0Verifier** (if deployed)
- **SP1Verifier** (if deployed)

> ⚠️ The pending owner must explicitly accept ownership in **Step 3**.

### Environment Variables

```bash
# Required
PRIVATE_KEY          # Deployer private key
CONTRACT_OWNER       # Address that will own all contracts
L2_CHAIN_ID          # Chain ID of the L2 network
L2_GENESIS_HASH      # Genesis block hash of L2

# Verifier Configuration
USE_DUMMY_VERIFIER=false           # Set true for devnet testing
DEPLOY_RISC0_RETH_VERIFIER=true    # Deploy RISC0 verifier
DEPLOY_SP1_RETH_VERIFIER=true      # Deploy SP1 verifier

# Inbox Configuration
BOND_TOKEN=0x0...                  # Bond token (0-address for ETH)
PROVING_WINDOW=7200                # Proving window in seconds (2 hours)
EXTENDED_PROVING_WINDOW=14400      # Extended proving window (4 hours)
MAX_FINALIZATION_COUNT=16          # Max proposals finalized per tx
FINALIZATION_GRACE_PERIOD=768      # Grace period (2 epochs)
RING_BUFFER_SIZE=100               # Proposal hash ring buffer size
BASEFEE_SHARING_PCTG=75            # Basefee sharing percentage
MIN_FORCED_INCLUSION_COUNT=1       # Min forced inclusions to process
FORCED_INCLUSION_DELAY=0           # Forced inclusion delay (seconds)
FORCED_INCLUSION_FEE_IN_GWEI=10000000  # Base fee (0.01 ETH)
FORCED_INCLUSION_FEE_DOUBLE_THRESHOLD=50  # Queue size for fee doubling
MIN_CHECKPOINT_DELAY=384           # Min checkpoint delay (1 epoch)
PERMISSIONLESS_INCLUSION_MULTIPLIER=5
COMPOSITE_KEY_VERSION=1

# Finality Gadget Configuration
OPTIMISTIC_FALLBACK_DELAY=604800   # 7 days for single proof finalization
FINALISING_PROOF_COUNT=2           # Proofs needed for immediate finalization
```

### Running the deployment

```bash
cd packages/protocol

# Simulation (dry run)
./script/layer1/surge/deploy_surge_l1.sh

# Broadcast transactions
BROADCAST=true ./script/layer1/surge/deploy_surge_l1.sh

# With contract verification
BROADCAST=true VERIFY=true ./script/layer1/surge/deploy_surge_l1.sh
```

### Output

Deployment addresses are written to `deployments/` directory as JSON files.

---

## Step 2: Configure Verifier Image IDs

After deploying the internal verifiers (Risc0Verifier, SP1Verifier), you must configure them with the correct prover image IDs.

> **Note**: The specific scripts for this step depend on your prover implementation. Consult the prover documentation for the image ID configuration process.

Each internal verifier needs its respective image ID set before proofs can be verified.

---

## Step 3: Accept L1 Ownership

**Script**: `script/layer1/surge/AcceptOwnership.s.sol`  
**Shell wrapper**: `script/layer1/surge/accept_ownership.sh`

### Purpose

Accept pending ownership for contracts that use the 2-step ownership transfer pattern (Ownable2Step).

### Contracts requiring ownership acceptance

From Step 1, the following contracts have `CONTRACT_OWNER` as their `pendingOwner`:

- Proof Verifier address
- Inbox proxy address
- SharedResolver address
- Risc0Verifier address (if deployed)
- SP1Verifier address (if deployed)

### Environment Variables

```bash
PRIVATE_KEY          # Private key of CONTRACT_OWNER from Step 1
CONTRACT_ADDRESSES   # Comma-separated list of contract addresses
FORK_URL             # L1 RPC URL
BROADCAST            # Set to "true" to execute transactions
```

### Running the script

```bash
cd packages/protocol

# Get addresses from deployment output and set them
export CONTRACT_ADDRESSES="0x...,0x...,0x..."  # Comma-separated

# Simulation
./script/layer1/surge/accept_ownership.sh

# Broadcast
BROADCAST=true ./script/layer1/surge/accept_ownership.sh
```

---

## Step 4: Setup L2 Contracts

**Script**: `script/layer2/surge/SetupSurgeL2.s.sol`  
**Shell wrapper**: `script/layer2/surge/setup_surge_l2.sh`

### What it does

1. **Verifies L2 registrations** - Confirms all L2 contracts are properly registered in the L2 SharedResolver
2. **Registers L1 contracts** - Adds L1 contract addresses to the L2 SharedResolver:
   - L1 Bridge
   - L1 SignalService
   - L1 ERC20Vault
   - L1 ERC721Vault
   - L1 ERC1155Vault
3. **Deploys DelegateController** - Creates a new DelegateController that will be the owner of L2 contracts
4. **Initiates ownership transfer** - Initiates ownership transfers of L2 contracts to the DelegateController:
   - Bridge
   - ERC20Vault
   - ERC721Vault
   - ERC1155Vault
   - SignalService
   - TaikoAnchor
   - SharedResolver

> ⚠️ These ownership transfers are **initiated only**. The DelegateController must accept ownership via a bridge message from L1.

### Environment Variables

```bash
# Script Configuration
PRIVATE_KEY          # Private key of current L2 contract owner

# L1 Configuration (from Step 1 deployment output)
L1_CHAINID           # L1 chain ID
L1_BRIDGE            # L1 Bridge address
L1_SIGNAL_SERVICE    # L1 SignalService address
L1_ERC20_VAULT       # L1 ERC20Vault address
L1_ERC721_VAULT      # L1 ERC721Vault address
L1_ERC1155_VAULT     # L1 ERC1155Vault address

# L1 Owner Configuration
L1_OWNER             # L1 DAO/Security Council/EOA that controls the DelegateController
```

### Running the script

```bash
cd packages/protocol

# Set L1 addresses from Step 1 deployment
export L1_BRIDGE="0x..."
export L1_SIGNAL_SERVICE="0x..."
export L1_ERC20_VAULT="0x..."
export L1_ERC721_VAULT="0x..."
export L1_ERC1155_VAULT="0x..."
export L1_OWNER="0x..."  # Same as CONTRACT_OWNER from Step 1

# Simulation
FOUNDRY_PROFILE=layer2 ./script/layer2/surge/setup_surge_l2.sh

# Broadcast
FOUNDRY_PROFILE=layer2 BROADCAST=true ./script/layer2/surge/setup_surge_l2.sh
```

### Output

The DelegateController address is written to `deployments/setup_l2.json`.

---

## Step 5: Accept L2 Ownership via Bridge

### Purpose

The DelegateController deployed in Step 4 needs to accept ownership of the L2 contracts. Since the DelegateController is controlled by the `L1_OWNER` (DAO/Security Council/EOA), this requires initiating the acceptance via a cross-chain message through the Bridge.

### How it works

1. The `L1_OWNER` sends a message through the L1 Bridge
2. The message targets the DelegateController on L2
3. The DelegateController executes `acceptOwnership()` on each L2 contract

### Contracts requiring ownership acceptance on L2

- Bridge
- ERC20Vault
- ERC721Vault
- ERC1155Vault
- SignalService
- TaikoAnchor
- SharedResolver

> ⚠️ **Note**: This step involves cross-chain messaging and is more complex. It is **not required for devnet deployments** where simpler ownership patterns may be used.

---

## Summary Checklist

- [ ] Genesis/chainspec file generated
- [ ] **Step 1**: L1 contracts deployed (`DeploySurgeL1.s.sol`)
- [ ] **Step 2**: Verifier image IDs configured
- [ ] **Step 3**: L1 ownership accepted (`AcceptOwnership.s.sol`)
- [ ] **Step 4**: L2 contracts configured (`SetupSurgeL2.s.sol`)
- [ ] **Step 5**: L2 ownership accepted via bridge (production only)

---
