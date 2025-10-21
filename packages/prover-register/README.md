# Register CLI

A CLI tool for registering and trusting prover instances collaterals in their verifier contracts based on their
Guest Data API.

## Usage

### Basic Command Structure

```bash
./prover-register \
  --verifier <VERIFIER_ADDRESS> \
  --type <VERIFIER_TYPE> \
  --prover <PROVER_URL> \
  --rpc <RPC_URL> \
  --private-key <PRIVATE_KEY> \
  [--trust] \
  [--register] \
  [--dry] \
  [--dry-as-owner] \
  [--env <ENV_FILE>]
```

### Parameters

- `--verifier`: Verifier contract address (required)
- `--type`: Verifier type - one of: `sgx`, `sp1`, `risc0`, `tdx`, or `azure-tdx` (required)
- `--prover`: Prover address URL (e.g., `http://192.168.1.100:9000`) (required)
- `--rpc`: Blockchain RPC URL (required)
- `--private-key`: Private key for transaction signing (can be set via PRIVATE_KEY env var)
- `--trust`: Trust the collateral (optional flag)
- `--register`: Register the instance (optional flag)
- `--dry`: Dry run mode - simulate transactions without sending (optional flag)
- `--dry-as-owner`: Dry run mode as the contract owner (optional flag)
- `--env`: Path to environment file (default: `.env`)

**Note**: At least one of `--trust` or `--register` must be specified.

#### Dry Run Mode

The `--dry` flag allows you to simulate transactions without actually sending them to the blockchain. This is useful for:

- Testing your configuration
- Verifying that transactions would succeed
- Estimating gas costs
- Debugging issues without spending gas

In dry run mode:

- No private key is required
- Transactions are simulated using `eth_call` and `eth_estimateGas`
- All operations are logged with `[DRY RUN]` prefix
- The tool validates that the transaction would succeed

### Environment Variables

You can use a `.env` file to store sensitive information:

```env
PRIVATE_KEY=your_private_key_here
```

### Examples

#### Register an SGX Instance

```bash
./prover-register \
  --verifier 0x1234567890123456789012345678901234567890 \
  --type sgx \
  --prover http://192.168.1.100:9000 \
  --rpc https://ethereum-rpc.example.com \
  --register \
  --env .env
```

#### Trust SP1 Collateral

```bash
./prover-register \
  --verifier 0x1234567890123456789012345678901234567890 \
  --type sp1 \
  --prover http://192.168.1.100:9000 \
  --rpc https://ethereum-rpc.example.com \
  --trust \
  --private-key 0xYourPrivateKey
```

#### Trust and Register TDX Instance (Raw TDX)

```bash
./prover-register \
  --verifier 0x1234567890123456789012345678901234567890 \
  --type tdx \
  --prover http://192.168.1.100:9000 \
  --rpc https://ethereum-rpc.example.com \
  --trust \
  --register \
  --env .env
```

#### Trust and Register Azure TDX Instance

```bash
./prover-register \
  --verifier 0x1234567890123456789012345678901234567890 \
  --type azure-tdx \
  --prover http://192.168.1.100:9000 \
  --rpc https://ethereum-rpc.example.com \
  --trust \
  --register \
  --env .env
```

#### Dry Run Example

```bash
./prover-register \
  --verifier 0x1234567890123456789012345678901234567890 \
  --type tdx \
  --prover http://192.168.1.100:9000 \
  --rpc https://ethereum-rpc.example.com \
  --trust \
  --register \
  --dry
```

## Verifier Type Details

### SGX (Software Guard Extensions)

- Fetches: `mr_enclave`, `mr_signer`, and `quote` from the prover
- Operations:
  - Register: Registers the instance using the attestation quote
  - Trust: No separate trust operation (trust is established through attestation)

### SP1

- Fetches: `aggregation_program_hash` and `block_program_hash`
- Operations:
  - Trust: Sets the program verification keys as trusted
  - Register: No instance registration required

### RISC0

- Fetches: `aggregation_program_hash` and `block_program_hash`
- Operations:
  - Trust: Sets the image IDs as trusted
  - Register: No instance registration required

### TDX (Trust Domain Extensions) - Raw

- Fetches: Raw TDX quote, nonce, and metadata
- Operations:
  - Trust: Sets trusted parameters for raw TDX instances (MR_TD, RT_MR, XFAM, etc.)
  - Register: Registers instance with raw attestation data

### Azure TDX

- Fetches: Complex attestation data with TPM quotes, PCR values, and runtime data
- Operations:
  - Trust: Sets trusted parameters for Azure TDX instances (TEE_TCB_SVN, MR_SEAM, MR_TD, PCR)
  - Register: Registers instance with Azure attestation data

## Guest Data API

The prover must expose a `/guest_data` endpoint that returns JSON data specific to the verifier type.

### Expected Response Format

#### SGX

```json
{
  "mr_enclave": "0x...",
  "mr_signer": "0x...",
  "quote": "0x..."
}
```

#### SP1

```json
{
  "sp1": {
    "aggregation_program_hash": "0x...",
    "block_program_hash": "0x..."
  }
}
```

#### RISC0

```json
{
  "risc0": {
    "aggregation_program_hash": "0x...",
    "block_program_hash": "0x..."
  }
}
```

#### TDX (Raw)

```json
{
  "issuer_type": "...",
  "public_key": "0x...",
  "quote": "0x...raw_tdx_quote...",
  "nonce": "0x...",
  "metadata": {}
}
```

#### Azure TDX

```json
{
  "issuer_type": "...",
  "public_key": "0x...",
  "quote": "0x...azure_tdx_quote...",
  "nonce": "0x...",
  "metadata": {}
}
```

### Building from Source

```bash
go mod tidy
go build -o prover-register
```
