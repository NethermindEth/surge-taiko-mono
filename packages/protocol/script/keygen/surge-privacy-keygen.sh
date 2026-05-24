#!/usr/bin/env bash
#
# surge-privacy-keygen — generate a complete privacy-mode key bundle for the
# Surge realtime fork.
#
# Outputs env-var lines for:
#   - SURGE_PRIVACY_SYMMETRIC_KEY      (32-byte AES-256-GCM key for scheme 0x01)
#   - SURGE_PRIVACY_SYMMETRIC_KEY_HASH (keccak256 of the above; baked into raiko guest)
#   - SURGE_PRIVACY_FI_PRIVKEY         (32-byte secp256k1 scalar for scheme 0x02 ECIES)
#   - SURGE_PRIVACY_FI_PRIVKEY_HASH    (keccak256 of the above; baked into raiko guest)
#   - SURGE_PRIVACY_FI_PUBKEY          (33-byte compressed secp256k1 pubkey, share publicly)
#
# Distribution:
#   - K_sym + FI_PRIVKEY  → set on Catalyst, driver, and raiko host (env vars at runtime).
#   - *_HASH              → set when compiling the raiko guest binaries (env vars at build).
#   - FI_PUBKEY           → publish in Surge docs / chain spec so external FI submitters
#                           can encrypt their blobs to it.
#
# Requires: openssl (1.1+), cast (foundry), jq.

set -euo pipefail

if ! command -v openssl >/dev/null; then
  echo "error: openssl not found" >&2
  exit 1
fi
if ! command -v cast >/dev/null; then
  echo "error: cast (foundry) not found — install via https://book.getfoundry.sh/getting-started/installation" >&2
  exit 1
fi
if ! command -v jq >/dev/null; then
  echo "error: jq not found" >&2
  exit 1
fi

# ---- Symmetric key (AES-256-GCM, scheme 0x01) ----
SYM_HEX=$(openssl rand -hex 32)
SYM_HASH=$(cast keccak "0x${SYM_HEX}")

# ---- Asymmetric keypair (secp256k1, scheme 0x02 / ECIES) ----
WALLET_JSON=$(cast wallet new --json)
FI_PRIVKEY=$(echo "${WALLET_JSON}" | jq -r '.[0].private_key')
FI_PRIVKEY_HASH=$(cast keccak "${FI_PRIVKEY}")

# Cast emits the uncompressed pubkey as 64 bytes: X (32B) || Y (32B), no leading 0x04.
# Compose the SEC1-compressed form: 0x02 if Y is even, 0x03 if Y is odd, followed by X.
PUB_UNCOMPRESSED=$(cast wallet public-key --raw-private-key "${FI_PRIVKEY}")
PUB_UNCOMPRESSED=${PUB_UNCOMPRESSED#0x}
PUB_X=${PUB_UNCOMPRESSED:0:64}
PUB_Y=${PUB_UNCOMPRESSED:64:64}
LAST_HEX_DIGIT=${PUB_Y: -1}
LAST_NIBBLE=$((16#${LAST_HEX_DIGIT}))
if (( LAST_NIBBLE % 2 == 0 )); then
  COMPRESS_PREFIX="02"
else
  COMPRESS_PREFIX="03"
fi
FI_PUBKEY="0x${COMPRESS_PREFIX}${PUB_X}"

cat <<EOF
# === Surge realtime privacy key bundle ===
# Generated $(date -u +%FT%TZ)

# Runtime env vars (Catalyst / driver / raiko host):
SURGE_PRIVACY_MODE=true
SURGE_PRIVACY_SYMMETRIC_KEY=0x${SYM_HEX}
SURGE_PRIVACY_FI_PRIVKEY=${FI_PRIVKEY}

# Build-time env vars (raiko guest compile):
SURGE_PRIVACY_SYMMETRIC_KEY_HASH=${SYM_HASH}
SURGE_PRIVACY_FI_PRIVKEY_HASH=${FI_PRIVKEY_HASH}

# Public — share with external FI submitters:
SURGE_PRIVACY_FI_PUBKEY=${FI_PUBKEY}
EOF
