# Surge Realtime Privacy Stack

How transaction privacy is wired across the Surge realtime fork. This is the durable reference engineers should read months/years from now to understand the design — independent of any one PR.

The privacy feature lives at the **blob payload level**: every L2 transaction list posted to L1 as an EIP-4844 blob is encrypted before broadcast and decrypted off-chain by trusted system components (driver + prover). The L1 protocol contracts contain no encryption logic.

## 1. Overview

Components touched:

| Component                     | Role                                                                               | Path                                                                                                          |
| ----------------------------- | ---------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| **L1 protocol contracts**     | Inbox + forced-inclusion queue. Encryption-agnostic.                               | [`packages/protocol/contracts/layer1/core/`](packages/protocol/contracts/layer1/core/)                        |
| **Catalyst (Rust)**           | Proposer/sequencer. **Encrypts** every blob it posts.                              | [`Catalyst/realtime/src/`](https://github.com/NethermindEth/Catalyst/tree/feat/realtime-privacy/realtime/src) |
| **Driver (Go, taiko-client)** | Follows L1 events, fetches blobs, **decrypts** for the L2 EL.                      | [`packages/taiko-client/`](packages/taiko-client/)                                                            |
| **raiko (Rust)**              | ZK prover. Host fetches encrypted blobs; **guest decrypts** under hash-bound keys. | [`raiko/`](https://github.com/NethermindEth/raiko/tree/feat/realtime-privacy)                                 |

The privacy boundary:

- **Private**: L2 transaction calldata posted on L1 blobs.
- **NOT private**: L2 P2P mempool, L1 blob hashes, L1 propose tx metadata (block number, block hash, state root checkpoint), the proposer's EOA, FI submitter EOAs and fees.

Privacy mode is toggled by a single shared env var, `SURGE_PRIVACY_MODE`, that all four components must agree on.

## 2. Threat model

| Adversary                                                   | Mitigated?                                                                                                                            |
| ----------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| Passive L1 observer (sees blobs, derives nothing about txs) | Yes (under classical crypto)                                                                                                          |
| Active L1 frontrunner of `propose` calls                    | Out of scope (existing concern; no privacy regression)                                                                                |
| Compromise of the symmetric key `K_sym`                     | All historical and future blobs become decryptable. Operator-managed rotation procedure required.                                     |
| Compromise of the FI private key `SK_sys`                   | All historical and future FI blobs become decryptable. Same rotation procedure.                                                       |
| Future cryptographically-relevant quantum computer (CRQC)   | Scheme 0x01 is PQ-safe (AES-256). Scheme 0x02 (ECIES on secp256k1) is **broken** by Shor — harvest-now-decrypt-later applies. See §9. |

## 3. Cipher schemes

Every privacy blob payload begins with a 1-byte **scheme id** so decoders can dispatch without a separate config flag per source:

| Scheme              | Algorithm                                          | Used for                                            | Key                               |
| ------------------- | -------------------------------------------------- | --------------------------------------------------- | --------------------------------- |
| `0x00`              | Plaintext (no encryption)                          | Out of privacy mode; FI blobs in non-privacy chains | n/a                               |
| `0x01`              | AES-256-GCM                                        | Catalyst's normal proposal blobs in privacy mode    | shared `K_sym`                    |
| `0x02`              | ECIES = secp256k1 ECDH ⊕ HKDF-SHA256 ⊕ AES-256-GCM | Forced-inclusion blobs in privacy mode              | system keypair `(SK_sys, PK_sys)` |
| `0x03` _(reserved)_ | ML-KEM-768 (FIPS 203)                              | Future post-quantum FI replacement                  | —                                 |
| `0x04` _(reserved)_ | Hybrid ML-KEM ⊕ X25519                             | Future PQ + classical hybrid FI                     | —                                 |

Adding a new scheme = one new impl + one new scheme id + one match arm in each component's dispatcher.

## 4. Blob layouts

A **blob** as broadcast on L1 (EIP-4844 / EIP-7594) is a 131,072-byte sidecar. Catalyst's existing path (`SidecarBuilder::from_slice`) packs a variable-length **payload buffer** into the blob's field elements with the standard padding scheme. Privacy mode does NOT change the EIP-4844 packing.

What it changes is the inner **payload buffer** layout:

```
payload_buffer = [version (1B)] [size (3B BE)] [scheme (1B)] [scheme_body (size-1 bytes)]
```

The outer `version`/`size` framing is unchanged; the scheme byte is the first byte of the inner content. Define `M = zlib(RLP(DerivationSourceManifest))` — the compressed manifest, exactly the same buffer that today is the inner content. This is the plaintext of every encryption operation.

### Scheme 0x00 — plaintext

```
[0x00 (1B)] [M (var)]
```

### Scheme 0x01 — AES-256-GCM (Catalyst's normal proposals)

```
[0x01 (1B)] [nonce (12B)] [C (len(M) bytes)] [tag (16B)]
```

- `nonce` is **96 bits, freshly drawn from a CSPRNG** per blob (`OsRng` in Rust, `crypto/rand` in Go). Public; in-band.
- `C = AES-256-GCM-encrypt(K_sym, nonce, M, aad=∅)` — same byte length as `M`.
- `tag` is the 16-byte GCM authentication tag.

### Scheme 0x02 — ECIES (forced-inclusion blobs)

```
[0x02 (1B)] [pk_eph (33B compressed secp256k1)] [C (len(M) bytes)] [tag (16B)]
```

Submitter side:

1. Draw ephemeral keypair `(sk_eph, pk_eph)` on secp256k1 from a CSPRNG.
2. `s = ECDH(sk_eph, PK_sys)` — 32 bytes (raw x-coordinate).
3. `(K_eph || nonce_eph) = HKDF-SHA256(salt=∅, ikm=s, info="surge-fi-v1", L=44)` — first 32B = AES key, last 12B = nonce.
4. `C || tag = AES-256-GCM-encrypt(K_eph, nonce_eph, M, aad=∅)`.
5. Emit `[0x02 || pk_eph || C || tag]`; **discard `sk_eph`**.

System side reverses with `s = ECDH(SK_sys, pk_eph)` (same shared secret by ECDH symmetry) and re-runs the HKDF.

The AES-GCM nonce is **not on the wire** for scheme 0x02 — both sides re-derive it. This is safe because `K_eph` is unique per submission (fresh `pk_eph` → fresh shared secret → fresh key); a constant nonce would also be safe but HKDF-deriving keeps both sides deterministic with no extra plumbing.

There is **no on-chain submitter pubkey registry**. `pk_eph` is generated fresh per submission, embedded in the blob, and discarded by the submitter. The submitter's only on-chain identity is the EOA paying the FI fee.

## 5. Key management

### Keys that exist

| Key      | Length                          | Where it lives                                                | Where its hash lives                    |
| -------- | ------------------------------- | ------------------------------------------------------------- | --------------------------------------- |
| `K_sym`  | 32 bytes (AES-256)              | Catalyst, driver, raiko host (env var); raiko guest (witness) | raiko guest binary (compile-time const) |
| `SK_sys` | 32 bytes (secp256k1 scalar)     | driver, raiko host (env var); raiko guest (witness)           | raiko guest binary (compile-time const) |
| `PK_sys` | 33 bytes (compressed secp256k1) | published in Surge docs / chain spec                          | n/a (public)                            |

Catalyst does NOT need `SK_sys` — it only references on-chain FI blob hashes via `numForcedInclusions`; it never decrypts.

### Hash-bound, witness-passed (for the raiko guest)

The raiko guest receives the secret keys via the witness (untrusted from the guest's POV) but verifies them against `keccak256` hashes baked into its binary at compile time:

- `SURGE_PRIVACY_SYMMETRIC_KEY_HASH` — env var read at compile time via `option_env!`.
- `SURGE_PRIVACY_FI_PRIVKEY_HASH` — same.

If either env var is set to a non-zero hash and the witness key's keccak256 doesn't match, the guest panics (the proof becomes unverifiable). When the env var is unset (defaults to all-zero hash), the check is bypassed — useful for non-privacy builds and CI.

**The vkey deployed on L1 thus commits to the keys** without ever putting the secret bytes into the public input. Rotation = recompile guest + redeploy verifier vkey. This is an explicit ops event, not silent.

### Generating keys

`packages/protocol/script/keygen/surge-privacy-keygen.sh` emits a complete env-var bundle:

```sh
$ bash packages/protocol/script/keygen/surge-privacy-keygen.sh
# === Surge realtime privacy key bundle ===
SURGE_PRIVACY_MODE=true
SURGE_PRIVACY_SYMMETRIC_KEY=0x...           # runtime: Catalyst, driver, raiko host
SURGE_PRIVACY_FI_PRIVKEY=0x...              # runtime: driver, raiko host
SURGE_PRIVACY_SYMMETRIC_KEY_HASH=0x...      # build-time: raiko guests
SURGE_PRIVACY_FI_PRIVKEY_HASH=0x...         # build-time: raiko guests
SURGE_PRIVACY_FI_PUBKEY=0x...               # public — share with FI submitters
```

Requires `openssl`, `cast` (foundry), `jq`. Generation is a single out-of-band step — components never auto-generate.

## 6. Forced-inclusion lifecycle

End-to-end flow (mirrors the legacy Pacaya design, ported to RealTimeInbox):

1. **Submission** (off-system, by an external user): builds an L2 tx list as a `DerivationSourceManifest`, RLP-encodes + zlib-compresses, ECIES-encrypts to `PK_sys` if privacy mode is in effect (or prepends `0x00` for plaintext), wraps in an EIP-4844 blob tx calling `RealTimeInbox.saveForcedInclusion(BlobReference)` with the current FI fee in ETH.
2. **On-chain queueing**: [`RealTimeInbox.saveForcedInclusion`](packages/protocol/contracts/layer1/core/impl/RealTimeInbox.sol) validates the blob via `LibBlobs.validateBlobReference` (the `blobhash` opcode resolves the hash of the same tx's blob), enqueues `ForcedInclusion { feeInGwei, blobSlice }` at `tail++`, refunds excess ETH.
3. **Catalyst consumption** (per proposal): reads `getForcedInclusionState() -> (head, tail)`, computes `numForcedInclusions = min(tail - head, fi_max_per_proposal)` and sets it on `ProposeInput`.
4. **On-chain dequeue** (inside `RealTimeInbox._consumeForcedInclusions`): pops `numForcedInclusions` from the queue, prepends them to the proposal's `sources[]` array (proposer's own blob last), forwards accumulated fees to `msg.sender`. If unconsumed FIs remain past `forcedInclusionDelay`, reverts with `UnprocessedForcedInclusionIsDue`.
5. **Driver derivation** (per source, in order): fetches blob bytes from beacon node, dispatches on the scheme byte, decrypts if needed, decompresses + RLP-decodes the manifest, applies blocks. On any failure for FI sources → produces a single empty L2 block with the anchor tx only.
6. **Prover (raiko guest) handling**: identical iteration. Each source's blob bytes are KZG-hash-bound to the on-chain blob hash; the guest dispatches on scheme byte, decrypts under hash-bound keys, replays L2 transactions. Same FI fallback to empty-block on decrypt/decode failure.
7. **Bootstrap & recovery**: empty queue → `numForcedInclusions = 0`. Stale Catalyst → permissionless fallback after `permissionlessInclusionMultiplier × forcedInclusionDelay`. Mode mismatch (FI submitted with wrong scheme byte) → empty-block fallback at the driver/guest level.

## 7. Privacy mode toggle

Single shared env var across all four components: `SURGE_PRIVACY_MODE=true|false`.

| Component             | What `SURGE_PRIVACY_MODE=true` requires                             | Behavior on missing key                                                           |
| --------------------- | ------------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| Catalyst              | `SURGE_PRIVACY_SYMMETRIC_KEY`                                       | Refuses to start.                                                                 |
| Driver (taiko-client) | `--privacy.symmetricKey` (`SURGE_PRIVACY_SYMMETRIC_KEY`)            | Refuses to start.                                                                 |
| raiko host            | `SURGE_PRIVACY_SYMMETRIC_KEY`, `SURGE_PRIVACY_FI_PRIVKEY`           | Forwards `None` to the guest; guest will fail decrypt for any non-plaintext blob. |
| raiko guest (build)   | `SURGE_PRIVACY_SYMMETRIC_KEY_HASH`, `SURGE_PRIVACY_FI_PRIVKEY_HASH` | Default zero-hash bypasses the binding check (non-privacy build).                 |

**Mismatch detection**: each blob is self-describing via the scheme byte. A non-FI source with a scheme the component cannot handle is a hard error (loud, not silent). FI sources fall back to empty-block-with-anchor. This means partial privacy-mode deployment cannot silently corrupt the chain — the worst case is a chain halt at the driver if the proposer encrypts but the driver lacks the key.

Operators verify lockstep by inspecting startup logs for the privacy-mode banner each component emits.

## 8. Code map

| Repo / area | File                                                                                                                               | Role                                                                                                                  |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------- |
| Protocol    | [`IRealTimeInbox.sol`](packages/protocol/contracts/layer1/core/iface/IRealTimeInbox.sol)                                           | Extends `IForcedInclusionStore`; `Config` carries FI delay/fee fields; `ProposeInput*` carries `numForcedInclusions`. |
| Protocol    | [`RealTimeInbox.sol`](packages/protocol/contracts/layer1/core/impl/RealTimeInbox.sol)                                              | FI queue (`LibForcedInclusion.Storage`), `saveForcedInclusion`, `_consumeForcedInclusions`. Encryption-agnostic.      |
| Protocol    | [`LibForcedInclusion.sol`](packages/protocol/contracts/layer1/core/libs/LibForcedInclusion.sol)                                    | Reused as-is from legacy Inbox.                                                                                       |
| Driver (Go) | [`pkg/privacy/`](packages/taiko-client/pkg/privacy/)                                                                               | `Cipher` interface, `Dispatch`, `Aes`, `Ecies`.                                                                       |
| Driver (Go) | [`shasta.go::manifestFromBlobBytesRealTime`](packages/taiko-client/driver/chain_syncer/event/manifest/shasta.go)                   | Per-source decrypt; FI fallback to default payload.                                                                   |
| Driver (Go) | [`driver/config.go`](packages/taiko-client/driver/config.go), [`cmd/flags/driver.go`](packages/taiko-client/cmd/flags/driver.go)   | CLI flags, env var parsing.                                                                                           |
| Catalyst    | `realtime/src/privacy/` (in [Catalyst](https://github.com/NethermindEth/Catalyst/tree/feat/realtime-privacy/realtime/src/privacy)) | `ProposalCipher::wrap` (encrypt-only).                                                                                |
| Catalyst    | `realtime/src/l1/proposal_tx_builder.rs` (Catalyst)                                                                                | Encryption hook between `encode_and_compress` and `SidecarBuilder`.                                                   |
| Catalyst    | `realtime/src/node/proposal_manager/async_submitter.rs` (Catalyst)                                                                 | Same encryption, applied to the Raiko proof request blob.                                                             |
| Catalyst    | `realtime/src/l1/execution_layer.rs` (Catalyst)                                                                                    | Reads FI queue state and sets `numForcedInclusions` on the propose input.                                             |
| raiko       | `lib/src/privacy/` (in [raiko](https://github.com/NethermindEth/raiko/tree/feat/realtime-privacy/lib/src/privacy))                 | `Cipher` traits, `dispatch_decrypt`, `Aes`, `Ecies`. no_std-safe.                                                     |
| raiko       | `lib/src/utils/realtime.rs` (raiko)                                                                                                | Per-source decrypt + compile-time hash binding via `option_env!`.                                                     |
| raiko       | `lib/src/input.rs` (raiko)                                                                                                         | `TaikoGuestBatchInput::privacy_*` witness fields.                                                                     |
| raiko       | `core/src/preflight/util.rs::prepare_taiko_chain_batch_input_realtime` (raiko)                                                     | Reads `SURGE_PRIVACY_*` env vars and populates the witness.                                                           |
| raiko       | `host/src/lib.rs::Opts` (raiko)                                                                                                    | CLI flags.                                                                                                            |
| Tooling     | [`packages/protocol/script/keygen/surge-privacy-keygen.sh`](packages/protocol/script/keygen/surge-privacy-keygen.sh)               | One-shot bundle generator.                                                                                            |

## 9. Post-quantum analysis

| Scheme | Primitive                | PQ confidentiality   | PQ authenticity          | Harvest-now-decrypt-later?                                                 |
| ------ | ------------------------ | -------------------- | ------------------------ | -------------------------------------------------------------------------- |
| 0x01   | AES-256-GCM              | **128-bit** (Grover) | ~64-bit forgery (online) | **No** — AES-256 is PQ-safe                                                |
| 0x02   | secp256k1 ECDH + AES-GCM | **Broken** by Shor   | Broken                   | **Yes** — encrypted FI archived today is decryptable in a post-CRQC future |

The PQ-vulnerable surface is bounded — FI is a low-volume censorship-resistance channel. Migration plan when needed: implement `0x03` = ML-KEM-768 or `0x04` = ML-KEM ⊕ X25519, deploy alongside (the dispatcher keeps the legacy 0x02 arm so old blobs still decrypt). Zero on-chain protocol change.

## 10. Operational runbook

### Bootstrap

1. Run `surge-privacy-keygen.sh`. Save the output securely.
2. Distribute `SURGE_PRIVACY_FI_PUBKEY` publicly (docs, chain spec).
3. Deploy raiko guests with `SURGE_PRIVACY_*_KEY_HASH` baked in. Note the resulting vkey.
4. Deploy `SurgeVerifier` on L1 with the new vkey.
5. Set `SURGE_PRIVACY_MODE=true` and the runtime keys on Catalyst + driver + raiko host.
6. Restart all four components. Verify each logs the "privacy mode: enabled" banner.

### Rotation

1. Generate a new bundle with `surge-privacy-keygen.sh`.
2. Recompile raiko guests with the new `*_HASH` env vars.
3. Deploy the new vkey to L1 (`SurgeVerifier`).
4. Drain in-flight proposals; restart Catalyst + driver + raiko host with the new runtime keys.
5. Old blobs proven under the old vkey remain verifiable forever (they're signed by the old vkey on L1). New blobs use the new keys.

### Mismatch recovery

If logs show "privacy dispatch failed" repeatedly, one component has the wrong keys. Stop the chain by stopping Catalyst, audit env vars on each component against the bundle file, restart in lockstep.

## 11. Open issues / future work

- **HKDF-per-proposal symmetric key derivation**: gives forward secrecy on master-key rotation. v2.
- **Threshold key management**: t-of-n distribution of `K_sym` and `SK_sys` so no single party holds the secret. v2.
- **Padding blobs to fixed size**: eliminates the side-channel where blob size leaks the tx-batch size. v2.
- **Encrypted P2P mempool**: privacy from validators/peers, not just L1. Out of scope for this stack.
- **Standalone FI submitter CLI**: a tool that takes a tx list + system pubkey and emits a ready-to-sign blob tx. Currently submitters do this manually.
