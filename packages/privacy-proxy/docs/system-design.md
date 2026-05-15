# privacy-proxy — system design

Engineering-internal reference for anyone modifying or extending the
proxy. Read this alongside the source.

## 1. Architecture

```
                  ┌──────────────────────────────────────────────────┐
   wallet  ─────► │  privacy-proxy (HTTP, axum)                      │
   (token)        │                                                  │
                  │  ┌────────────┐   ┌──────────┐   ┌────────────┐  │
                  │  │ auth layer │ → │ rpc      │ → │ acl filter │  │
                  │  │ (token →   │   │ dispatch │   │ (registry  │  │
                  │  │  role +    │   │          │   │  + lambdas)│  │
                  │  │  info)     │   │          │   │            │  │
                  │  └────────────┘   └──────────┘   └─────┬──────┘  │
                  │                                        │         │
                  │                          debug_traceCall│         │
                  │                                        ▼         │
                  │                                  ┌─────────────┐ │
                  │                                  │ upstream    │ │
                  │                                  │ forwarder   │ │
                  │                                  └─────┬───────┘ │
                  └────────────────────────────────────────│─────────┘
                                                           ▼
                                                    Nethermind RPC
```

## 2. Request flow

For every HTTP request:

1. **Auth resolution.** The `caller_ctx_layer` middleware
   ([src/auth/middleware.rs](../src/auth/middleware.rs)) extracts the
   `Authorization: Bearer <token>` header, looks up the sha256 hash in
   `auth_tokens`, joins to `users` and `roles`, and inserts a
   `CallerCtx` into the request's extensions. Missing / invalid /
   expired tokens fall back to `CallerCtx::anonymous()`.
2. **JSON-RPC dispatch.** `rpc::dispatch`
   ([src/rpc/handlers.rs](../src/rpc/handlers.rs)) parses the JSON
   body (rejecting batches), inspects `method`, and:
   - Outside `eth_` namespace → JSON-RPC error `-32601`.
   - `eth_` with no call data (`eth_blockNumber`, `eth_getBalance`,
     `eth_chainId`, etc.) → forward unchanged.
   - **Call-bearing**: `eth_call`, `eth_estimateGas`,
     `eth_createAccessList`, `eth_sendTransaction`,
     `eth_sendRawTransaction` → ACL filter.
3. **ACL filter.**
   1. Extract `(to, data)`. For `eth_sendRawTransaction`, RLP-decode
      via [`tracer::decode_raw_tx`](../src/tracer/mod.rs) to recover
      `from`, `to`, `value`, and `input`.
   2. Top-level check via [`acl::check_call`](../src/acl/evaluator.rs).
      Cheap; aborts before tracer if obviously denied.
   3. Issue `debug_traceCall` with `callTracer` to upstream via
      [`tracer::trace_call`](../src/tracer/mod.rs).
   4. Walk every `CALL`/`STATICCALL`/`DELEGATECALL`/`CALLCODE` frame
      (`CallFrame::flatten`). For each, run `check_call`.
4. **Forward or deny.** If every frame passes, forward the original
   request unchanged. Otherwise return JSON-RPC error `-32001` with
   `data: { contract, selector, detail }`.

Admin endpoints are mounted on `/admin/*` and gated by
[`admin::middleware::admin_gate`](../src/admin/middleware.rs): `401` if
no token, `403` if the resolved role is not `admin`, otherwise the
handler runs.

## 2.5 Method-level gating (synthetic selectors)

Several `eth_` methods read account state directly without going
through a contract call. These are gated using **synthetic 4-byte
selectors** so the same `access_rules` table works without a schema
change. See [`src/rpc/gated_methods.rs`](../src/rpc/gated_methods.rs).

Reserved range: `0xff______`. Current allocations:

| RPC method                | Synthetic selector |
| ------------------------- | ------------------ |
| `eth_getBalance`          | `0xff010001`       |
| `eth_getTransactionCount` | `0xff010002`       |
| `eth_getCode`             | `0xff010003`       |
| `eth_getStorageAt`        | `0xff010004`       |
| `eth_getProof`            | `0xff010005`       |

`rpc::handlers::dispatch` pre-checks each request before the
call-bearing classification. The flow for a gated method:

1. Extract `target = params[0]`. Bail `-32602` if missing/invalid.
2. Look up an `access_rules` row keyed by `(target, synthetic_selector)`.
3. If the row exists, run `acl::check_call` with synthetic
   `call_data = synthetic_selector ++ pad32(target)` (and, for
   `eth_getStorageAt`, the 32-byte slot appended). Lambdas can decode
   the target from `call_data[16..36]`.
4. If no row exists, apply the default policy:
   - Anonymous caller → deny (`AnonymousAgainstGatedCall`).
   - `state.upstream.is_contract(target)` → on `true`, allow;
     on `false`, allow iff `target == caller.eoa`, else deny
     (`DefaultEoaSelfOnly`).

`UpstreamClient::is_contract` calls upstream `eth_getCode(target,
"latest")` and caches the result in a `Mutex<HashMap<Address, (bool,
Instant)>>` with a 60s TTL — bounds the staleness window of an EOA
that later deploys to its address.

`admin::normalize_selector` accepts a JSON-RPC method name (e.g.
`"eth_getBalance"`) in place of the synthetic hex selector and
rewrites it on write, so operators authoring rules through the admin
API don't have to memorize the synthetic values. The new
`GET /admin/registry/synthetic-selectors` endpoint surfaces the full
map for the future UI.

## 3. ACL semantics

For a given `(contract, selector)`, the decision is:

| rule mode | entry for caller role? | lambda | result                               |
| --------- | ---------------------- | ------ | ------------------------------------ |
| no rule   | —                      | —      | allow                                |
| `allow`   | no                     | —      | DENY (`NotInAllowList`)              |
| `allow`   | yes                    | none   | allow                                |
| `allow`   | yes                    | true   | allow                                |
| `allow`   | yes                    | false  | DENY (`LambdaRejected`)              |
| `deny`    | no                     | —      | allow                                |
| `deny`    | yes                    | none   | DENY (`InDenyList`)                  |
| `deny`    | yes                    | true   | DENY (`InDenyList`)                  |
| `deny`    | yes                    | false  | allow                                |
| any       | unknown lambda name    | —      | DENY (`UnknownLambda`) — fail closed |
| `allow`   | anonymous caller       | —      | DENY (`AnonymousAgainstGatedCall`)   |
| `deny`    | anonymous caller       | —      | allow                                |

Empty call data (< 4 bytes) is always allowed — no selector to gate.

For **gated RPC methods** (synthetic-selector path), the dispatcher
adds a default-policy layer on top of this table:

| rule presence for `(target, synthetic_selector)` | target kind | caller               | result                             |
| ------------------------------------------------ | ----------- | -------------------- | ---------------------------------- |
| present                                          | —           | —                    | apply the table above              |
| absent                                           | EOA         | anonymous            | DENY (`AnonymousAgainstGatedCall`) |
| absent                                           | EOA         | target == caller.eoa | allow                              |
| absent                                           | EOA         | target != caller.eoa | DENY (`DefaultEoaSelfOnly`)        |
| absent                                           | contract    | any non-anon         | allow                              |
| absent                                           | contract    | anonymous            | DENY (`AnonymousAgainstGatedCall`) |

Implementation: [`acl::evaluator::check_call`](../src/acl/evaluator.rs).

## 4. DB schema

SQLite, file-backed by default. Migrations run on startup via
`sqlx::migrate!`. See [migrations/0001_init.sql](../migrations/0001_init.sql).

| Table                 | Purpose                                                                                                                                                                |
| --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `roles`               | Dynamic role definitions. Two seeded: `admin`, `user`.                                                                                                                 |
| `users`               | `eoa_address → role_id + caller_info_json`. JSON is opaque to the DB so new roles need no migrations.                                                                  |
| `auth_tokens`         | `sha256(token) → eoa_address` with expiry. Token plaintext is never persisted.                                                                                         |
| `challenges`          | Short-lived sign-in nonces. Required to prove (a) the wallet signed _this_ server's freshly issued nonce (no replay) and (b) the wallet has the private key right now. |
| `access_rules`        | `(contract_address, function_selector) → mode`. Unique on the pair.                                                                                                    |
| `access_rule_entries` | One per role under a rule. Optional `lambda_name` ties to an in-build lambda. Unique on `(rule_id, role_id)`.                                                          |

## 5. Lambdas

`LambdaSpec` and `LambdaFn` live in
[src/acl/lambdas/mod.rs](../src/acl/lambdas/mod.rs):

```rust
pub struct LambdaCtx<'a> {
    pub caller_info: &'a serde_json::Value,
    pub selector:    [u8; 4],
    pub call_data:   &'a [u8],
}
pub type LambdaFn = fn(&LambdaCtx) -> bool;
pub struct LambdaSpec {
    pub name: &'static str,
    pub description: &'static str,
    pub expected_selector: Option<[u8; 4]>,
    pub run: LambdaFn,
}
```

`registry()` returns a `&'static HashMap<&'static str, &'static LambdaSpec>`
built once via `LazyLock` from a static table. Description and selector
are static literals next to the `fn` pointer — code is the single
source of truth.

**Caller identity in `caller_info`.** Anything a lambda needs about the
caller flows through the single `caller_info` JSON. The auth middleware
([src/auth/middleware.rs](../src/auth/middleware.rs)) is the trusted
writer of identity fields — after deserializing the DB-stored
`caller_info_json`, it injects `caller_info.eoa = "0x{lowercase hex}"`
for every authenticated request (overwriting any admin-stored `eoa`).
Lambdas read `caller_info.eoa` like any other field; if it's missing or
unparseable, the lambda returns `false` and the rule's allow/deny
semantics decide. New auth-derived fields (role, issued_at, …) get a
single `map.insert(...)` call here and a corresponding read in
whichever lambda needs them — no per-lambda declarations, no merge
passes, no `LambdaCtx` changes.

**Constraint**: lambdas ship with the binary. Adding one = (a) implement
the fn in [src/acl/lambdas/examples.rs](../src/acl/lambdas/examples.rs)
(or a sibling), (b) add a `LambdaSpec` entry in `registry()`. Admins
attach lambdas to rule entries by name via the admin API; an unknown
name on write returns `400`, and an unknown name discovered at
evaluation time fails closed (`UnknownLambda`).

## 6. Admin management & root trust

`ADMIN_EOAS` (comma-separated env var) is the trust anchor. On every
boot, [`admin::reconcile_seed_admins`](../src/admin/mod.rs) upserts each
EOA into `users` with `role = admin`. Properties:

- Operator owns root via deploy config.
- All admin auth still goes through the wallet-signature flow → no
  long-lived admin shared secret.
- **Break-glass**: even if every other admin is demoted in DB, the
  next restart re-promotes everyone in `ADMIN_EOAS`. The seed cannot be
  locked out by DB state alone.
- **Key rotation**: edit `ADMIN_EOAS` and restart. To revoke active
  sessions for a compromised EOA without rotation:
  `DELETE /admin/users/:eoa/tokens`.

## 7. Authentication & tokens

Challenge-response sequence:

1. `GET /auth/challenge?address=0x…` →
   [`auth::challenge::handler`](../src/auth/challenge.rs) generates 16
   random bytes (hex-encoded), upserts into `challenges` keyed by EOA
   with a ~5 minute TTL, and returns the EIP-191 message:
   ```
   {domain} sign-in
   Address: 0x{lowercase eoa}
   Nonce: {nonce hex}
   ```
2. Wallet signs the message (personal_sign).
3. `POST /auth/verify { address, signature }` →
   [`auth::verify::handler`](../src/auth/verify.rs) loads the pending
   nonce, reconstructs the exact message, recovers the signer via
   `alloy::primitives::Signature::recover_address_from_msg`, deletes
   the consumed nonce, upserts the `users` row (default role = `user`),
   mints a fresh 32-byte random token, stores its sha256, returns the
   token plaintext (only time the plaintext exists outside the wallet).

Token format: 64 hex chars. Stored hashed: `sha256(token)`. Tokens are
opaque; rotating them is "revoke + re-sign-in".

## 8. Filtering via `debug_traceCall`

`callTracer` returns the call tree as a nested JSON of frames with
`{type, from, to, input, value, calls?}`. `CallFrame::flatten`
recursively collects every CALL-family frame into a flat
`Vec<CallSite>`, which the dispatcher iterates against the registry.

**Known caveat — `eth_sendRawTransaction`:** the trace is a
_simulation_ against the latest state at the time of the proxy's call.
State may shift before the tx is mined; the ACL outcome could differ
from on-chain reality. This is acceptable for the POC since the proxy's
job is to gate **what the user can attempt**, not what eventually lands.
A future hardening would gate at execution time (rejecting receipts
post-hoc) or freeze the simulation block tag to `pending` and
re-validate aggressively.

## 9. Module layout

| Path                                                    | Responsibility                                                            |
| ------------------------------------------------------- | ------------------------------------------------------------------------- |
| [src/main.rs](../src/main.rs)                           | bin entry — calls `lib::run()`                                            |
| [src/lib.rs](../src/lib.rs)                             | module aggregator + bootstrap                                             |
| [src/config.rs](../src/config.rs)                       | env-based config                                                          |
| [src/db.rs](../src/db.rs)                               | sqlx pool, migration, `now_unix`                                          |
| [src/state.rs](../src/state.rs)                         | `AppState { config, pool, upstream }`                                     |
| [src/upstream.rs](../src/upstream.rs)                   | reqwest client to Nethermind                                              |
| [src/error.rs](../src/error.rs)                         | `ApiError` + JSON error response shape                                    |
| [src/server.rs](../src/server.rs)                       | axum router assembly                                                      |
| [src/auth/](../src/auth/)                               | challenge, verify, middleware, `CallerCtx`                                |
| [src/rpc/](../src/rpc/)                                 | JSON-RPC dispatcher, method classification                                |
| [src/rpc/gated_methods.rs](../src/rpc/gated_methods.rs) | synthetic selectors for address-parameterized read methods + helpers      |
| [src/acl/evaluator.rs](../src/acl/evaluator.rs)         | allow/deny semantics                                                      |
| [src/acl/registry.rs](../src/acl/registry.rs)           | DB reads of rules + entries                                               |
| [src/acl/lambdas/](../src/acl/lambdas/)                 | named lambda registry + examples                                          |
| [src/tracer/mod.rs](../src/tracer/mod.rs)               | `debug_traceCall` + frame walk + raw tx decode                            |
| [src/admin/](../src/admin/)                             | `/admin/*` routes (roles, users, registry, lambdas) + seed reconciliation |

## 10. Out of scope (deferred) — engineering notes

| Feature                              | Notes for whoever picks this up                                                                                                                                                                                                                                                                                                                                                      |
| ------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| WebSocket transport                  | Add a `/ws` route + `axum::extract::ws`; per-message auth via initial frame; subscription state on the proxy mapping client subs → upstream subs.                                                                                                                                                                                                                                    |
| Non-`eth_` namespaces                | Add classifier branches in `rpc::handlers::dispatch`; decide per-namespace whether to passthrough or gate.                                                                                                                                                                                                                                                                           |
| Token refresh                        | New endpoint `POST /auth/refresh` taking the current token + signature over a fresh nonce. Or shorten TTL and add a sliding window.                                                                                                                                                                                                                                                  |
| Batch JSON-RPC                       | Parse top-level array; map each element through the same dispatcher; assemble response array. Watch for partial-failure semantics.                                                                                                                                                                                                                                                   |
| Rate limiting                        | tower-http `RequestBodyLimitLayer` plus a per-token counter in DB or in-memory (`DashMap`).                                                                                                                                                                                                                                                                                          |
| Per-admin scopes                     | New `admin_scopes` table; admin gate becomes a per-route capability check instead of a single boolean.                                                                                                                                                                                                                                                                               |
| Audit log                            | Append-only `admin_audit` table written by an axum middleware on `/admin/*` success.                                                                                                                                                                                                                                                                                                 |
| Step-up auth                         | Tag specific routes; require a fresh signature timestamp within N minutes from `auth_tokens`.                                                                                                                                                                                                                                                                                        |
| Multi-chain                          | Per-`chain_id` config and DB partitioning; one binary per chain remains simplest.                                                                                                                                                                                                                                                                                                    |
| `BALANCE` opcode leak via `eth_call` | The callTracer used today emits no frame for opcode-level state reads (`BALANCE`, `EXTCODESIZE`, `EXTCODEHASH`). Plug by switching to `prestateTracer` (returns every account touched during simulation) and adding a new check pass that compares the touched-set against the per-role ACL. Higher latency, more false positives — only worth doing if the threat model demands it. |

## 11. Verification matrix

| Layer                                     | Test type   | Where                                                                                                                                                                   |
| ----------------------------------------- | ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ACL truth table                           | unit        | [src/acl/evaluator.rs](../src/acl/evaluator.rs) — 12 cases covering allow/deny × entry × lambda × anonymous + unknown lambda.                                           |
| Lambda decoding                           | unit        | [src/acl/lambdas/examples.rs](../src/acl/lambdas/examples.rs) — `require_kyc`, `transfer_under_limit`.                                                                  |
| Admin gate                                | integration | [tests/admin.rs](../tests/admin.rs) — `401` / `403` / `200` matrix on `/admin/roles`.                                                                                   |
| Lambda listing                            | integration | [tests/admin.rs](../tests/admin.rs) — `GET /admin/registry/lambdas`.                                                                                                    |
| Restart reconciliation                    | integration | [tests/admin.rs](../tests/admin.rs) — pre-seed a user → reconcile → assert promoted to admin.                                                                           |
| Non-eth namespace rejection               | integration | [tests/admin.rs](../tests/admin.rs).                                                                                                                                    |
| Synthetic selector encoding               | unit        | [src/rpc/gated_methods.rs](../src/rpc/gated_methods.rs) — round-trip lookup + ABI layout for target/slot.                                                               |
| `target_in_caller_allowlist` lambda       | unit        | covered indirectly via the allowlist integration test below.                                                                                                            |
| `erc20_self_only` lambda                  | unit        | [src/acl/lambdas/examples.rs](../src/acl/lambdas/examples.rs) — balanceOf and allowance × self/other × missing/unparseable eoa × unknown selector × malformed calldata. |
| Auth-layer injects `caller_info.eoa`      | unit        | [src/auth/middleware.rs](../src/auth/middleware.rs) — preserves admin-set fields, overwrites stale `eoa`, synthesizes object when caller_info is empty.                 |
| Default gated-method policy               | integration | [tests/gated_methods.rs](../tests/gated_methods.rs) — self-allowed, other-EOA denied, contract free, anonymous denied.                                                  |
| Admin override on gated method            | integration | [tests/gated_methods.rs](../tests/gated_methods.rs) — deny rule on contract; allow rule with `target_in_caller_allowlist`.                                              |
| Method-name selector normalization        | integration | [tests/gated_methods.rs](../tests/gated_methods.rs) — `POST /admin/registry/rules` with `function_selector: "eth_getBalance"` stores `0xff010001`.                      |
| `GET /admin/registry/synthetic-selectors` | integration | [tests/gated_methods.rs](../tests/gated_methods.rs).                                                                                                                    |
| End-to-end RPC + tracer                   | manual      | run against real Nethermind; see operator-guide.md.                                                                                                                     |
