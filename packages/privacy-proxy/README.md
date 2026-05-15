# privacy-proxy

A drop-in replacement for an Ethereum JSON-RPC endpoint that adds
wallet-signature-based authentication and role-based access control over
contract calls — **including internal calls inside any transaction or
`eth_call`**. Wallets talk to it exactly as they would talk to a normal
Ethereum RPC node; the only addition is one HTTP header carrying an auth
token. Admins configure who can call what; users sign in with their
wallet and receive a token for subsequent requests.

This package is a POC scoped to the `eth_` namespace and HTTP transport.

---

## How users interact with it

1. **Sign in once.** The wallet asks for a challenge, signs it with the
   user's EOA, and receives a bearer token (TTL ~7 days). See
   [docs/wallet-integration.md](docs/wallet-integration.md) for the
   wire-level flow.
2. **Send normal JSON-RPC.** Every subsequent request goes to the proxy
   URL with `Authorization: Bearer <token>`. Methods that don't carry a
   call (`eth_blockNumber`, `eth_getBalance`, `eth_chainId`, …) are
   always forwarded.
3. **Call-bearing methods are gated.** `eth_call`, `eth_estimateGas`,
   `eth_sendRawTransaction`, `eth_createAccessList`, and
   `eth_sendTransaction` are checked against the access registry. The
   proxy validates both the top-level call **and every internal call**
   the transaction would make; any forbidden frame rejects the whole
   request. A small set of address-parameterized reads
   (`eth_getBalance`, `eth_getProof`, `eth_getTransactionCount`,
   `eth_getCode`, `eth_getStorageAt`) is also gated — by default an
   EOA can only query its own state; admins can restrict contract
   targets explicitly.
4. **Denials are explicit.** A blocked call returns JSON-RPC error code
   `-32001` with the offending `(contract, selector)` in `data`. The
   wallet should surface this to the user.
5. **Token expiry.** When the token expires, repeat the sign-in flow.
   There is no refresh endpoint by design.

A user-side flow (curl):

```bash
# 1. ask for a challenge
curl "$PROXY_URL/auth/challenge?address=0xYourEoa"
# → { "message": "...", "expires_at": 1715600000 }

# 2. wallet signs the message; POST to /auth/verify
curl -X POST "$PROXY_URL/auth/verify" -H 'content-type: application/json' \
  -d '{"address":"0xYourEoa","signature":"0x...65bytes..."}'
# → { "token": "abcdef...", "expires_at": 1716200000 }

# 3. use the token on any JSON-RPC request
curl -X POST "$PROXY_URL/" -H "authorization: Bearer $TOKEN" \
  -H 'content-type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"eth_blockNumber","params":[]}'
```

---

## How admins interact with it

The operator seeds initial admins via the `ADMIN_EOAS` env var at deploy
time; those EOAs are promoted on every restart. After that, admins sign
in with their wallet exactly like users — their token resolves to the
`admin` role. There is no separate admin password.

Admins use 18 endpoints under `/admin/*` (see
[docs/admin-api.md](docs/admin-api.md)) to:

- **Manage the access registry**: list/create/update/delete rules and
  per-role entries; pick from the in-build list of named lambdas via
  `GET /admin/registry/lambdas`.
- **Manage roles**: list/create/delete.
- **Manage users**: list/get/upsert (role + caller info)/delete; revoke
  all active sessions for an EOA.

All admin endpoints require an admin-role token; a user-role token gets
`403`, no token gets `401`.

A future operator UI will consume these endpoints directly; the
response shapes are already designed for it (paginated lists, single
gets, full upserts).

---

## The access model in plain English

- The proxy holds an **access registry** keyed by
  `(contract address, function selector)`.
- A registry entry is either an **allow** list of roles or a **deny**
  list — never both.
- Each role in the entry may have an optional named **lambda**: a
  code-defined predicate over the caller's stored info **plus the
  function's arguments**. Admins pick a lambda from a fixed in-build
  list.
- Contracts or functions not in the registry are **freely callable** —
  the registry is opt-in.
- Filtering applies to internal calls too. If any contract call inside
  the execution of your transaction would be blocked for your role, the
  whole request is rejected **before it reaches the chain**.
- Address-parameterized reads (e.g. `eth_getBalance`) are gated
  separately by the same registry, keyed by `(target_address,
synthetic_selector)`. EOAs default to self-only; contracts are free
  unless the admin installs a rule.

---

## What's NOT included (deferred)

- WebSocket transport and JSON-RPC subscriptions (`eth_subscribe`).
- Any non-`eth_` namespace (`net_`, `web3_`, `debug_`, `trace_`, …) —
  rejected, not passed through.
- Token refresh endpoint — expired tokens require re-signing.
- Batch JSON-RPC requests — a single request per HTTP call.
- Rate limiting and abuse protection.
- Per-admin permission scopes — all admins are equal; no read-only or
  registry-only admin.
- Audit log of admin actions.
- Step-up authentication for sensitive admin endpoints (fresh-signature
  requirement).
- Runtime lambda upload or hot-reload — lambdas ship with the binary.
- Multiple roles per EOA — one user, one role.
- Multi-chain support — a deployment serves exactly one upstream chain.
- Reading balances via the `BALANCE` EVM opcode inside an `eth_call`
  (`a.balance` in attacker-supplied bytecode). The `callTracer` emits
  no frame for opcode-level state access; closing this would require
  switching to `prestateTracer`.

---

## Deeper docs

- [docs/wallet-integration.md](docs/wallet-integration.md) — for wallet
  developers integrating the sign-in flow.
- [docs/admin-api.md](docs/admin-api.md) — the full admin endpoint
  reference for the future admin UI.
- [docs/system-design.md](docs/system-design.md) — engineering-internal
  reference (architecture, request flow, schema, modules).
- [docs/operator-guide.md](docs/operator-guide.md) — deploying,
  `ADMIN_EOAS`, key rotation, restart semantics.

## Quickstart for local dev

```bash
export UPSTREAM_URL=http://localhost:8545                  # your Nethermind
export ADMIN_EOAS=0xYourSeedAdminEoa                       # seed admins
export DATABASE_URL=sqlite://./privacy-proxy.db            # default
export BIND_ADDR=0.0.0.0:8080                              # default

cargo run --release
```

Migrations run automatically on boot.
