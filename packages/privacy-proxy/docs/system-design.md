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
   `auth_tokens`, joins to `members` and `roles`, loads the role's typed
   attributes (from `user_attributes` for `user`; identity-only for
   `admin`), and inserts a `CallerCtx { eoa, attributes:
CallerAttributes::{Admin|User}(...) }` into the request's
   extensions. Missing / invalid / expired tokens fall back to
   `CallerCtx::anonymous()`.
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
| `roles`               | Mirror of the static `ROLES: &[&str]` in `src/roles.rs`. Reconciled at boot. Adding a role is a code change + migration; no runtime mutation API exists.               |
| `members`             | `eoa_address → role_id`. Identity-only; per-role attributes live in role-specific tables. Both admins and users are rows here.                                         |
| `user_attributes`     | Per `user`-role attributes: `(eoa_address PK FK → members, kyc bool, blacklisted bool)`. Admin role has no attribute table — an admin's only attribute is their EOA.   |
| `auth_tokens`         | `sha256(token) → eoa_address` with expiry. Token plaintext is never persisted.                                                                                         |
| `challenges`          | Short-lived sign-in nonces. Required to prove (a) the wallet signed _this_ server's freshly issued nonce (no replay) and (b) the wallet has the private key right now. |
| `access_rules`        | `(contract_address, function_selector) → mode`. Unique on the pair.                                                                                                    |
| `access_rule_entries` | One per role under a rule. Optional `lambda_name` ties to an in-build lambda. Unique on `(rule_id, role_id)`.                                                          |

## 5. Lambdas

The base types in [src/acl/lambdas/mod.rs](../src/acl/lambdas/mod.rs)
are generic over the caller's attribute type `C`:

```rust
pub struct LambdaCtx<'a, C> {
    pub caller_info: &'a C,
    pub selector:    [u8; 4],
    pub call_data:   &'a [u8],
}
pub struct LambdaSpec<C: 'static> {
    pub name: &'static str,
    pub description: &'static str,
    pub expected_selectors: &'static [[u8; 4]],
    pub run: fn(&LambdaCtx<C>) -> bool,
}
```

Lambdas live in **role-specific directories** under `src/acl/lambdas/<role>/`.
Each role declares its attribute struct in its `mod.rs` and exposes a
`registry()` returning `&'static HashMap<&'static str, &'static LambdaSpec<RoleAttrs>>`.
Currently:

- `src/acl/lambdas/user/` — `UserCallerInfo { eoa, kyc, blacklisted }`.
  Lambdas: `require_kyc`, `erc20_self_only`.
- Admin has no lambda registry: an admin's only attribute is their
  EOA, and write-time validation in [`admin/registry.rs`](../src/admin/registry.rs)
  rejects `lambda_name` on admin entries.

The evaluator unwraps `CallerCtx.attributes` (a `CallerAttributes`
tagged enum) and dispatches to the matching role's registry. Each
lambda receives the typed struct directly — no JSON probing, no
silently-missing fields.

**Constraint**: lambdas ship with the binary. Adding one for a role =
(a) add a sibling module under that role's directory with the
`fn(&LambdaCtx<RoleAttrs>) -> bool` implementation, (b) add a
`LambdaSpec` entry in the role's `registry()`. Admins attach lambdas
to rule entries by name via the admin API; an unknown name on write
returns `400` (validated against the _target role's_ registry), and an
unknown name discovered at evaluation time fails closed
(`UnknownLambda`).

## 6. Roles

The set of roles is declared in code, not in the database. The static
slice `ROLES: &[&str]` in [src/roles.rs](../src/roles.rs) is the source
of truth; on boot, `reconcile_roles` inserts each name into the `roles`
table so foreign keys from `members.role_id` and
`access_rule_entries.role_id` resolve. The admin API exposes a single
read endpoint (`GET /admin/roles`) for enumeration.

For each role, two more pieces may exist alongside the name:

- A **typed attribute struct** plus a backing table, if the role
  carries per-member state beyond identity. The `user` role declares
  `UserCallerInfo { eoa, kyc, blacklisted }` and persists it in
  `user_attributes`. The `admin` role is identity-only and has no
  attribute table.
- A **lambda registry** under `src/acl/lambdas/<role>/`, holding the
  `LambdaSpec<RoleAttrs>` entries the role accepts on its rule entries.
  Roles without role-specific gating (e.g. `admin`) have no registry
  and reject any `lambda_name` on their entries at write time.

### Adding a role

To introduce a new role `foo` with attributes `bar: bool`:

1. **Declare the role name.** In [src/roles.rs](../src/roles.rs), add
   `pub const ROLE_FOO: &str = "foo";` and append `ROLE_FOO` to the
   `ROLES` slice. `reconcile_roles` will insert it into the `roles`
   table on the next boot.

2. **Add the attribute table.** Create a new migration in
   `migrations/` (next numbered prefix):

   ```sql
   CREATE TABLE foo_attributes (
       eoa_address TEXT PRIMARY KEY REFERENCES members(eoa_address) ON DELETE CASCADE,
       bar         INTEGER NOT NULL DEFAULT 0 CHECK (bar IN (0, 1))
   );
   ```

   Skip this step if the role is identity-only (like `admin`).

3. **Declare the typed attribute struct.** Create
   `src/acl/lambdas/foo/mod.rs` with `pub struct FooCallerInfo { eoa,
bar }` plus a `registry()` returning
   `&'static HashMap<&str, &LambdaSpec<FooCallerInfo>>` (start empty if
   no lambdas yet). Wire it under `pub mod foo;` in
   [src/acl/lambdas/mod.rs](../src/acl/lambdas/mod.rs).

4. **Extend `CallerAttributes`.** In [src/auth/mod.rs](../src/auth/mod.rs),
   add a `Foo(FooCallerInfo)` variant to `CallerAttributes` and a
   matching arm in `role_name()` / `eoa()`.

5. **Load the attributes on token resolution.** In
   [src/auth/middleware.rs](../src/auth/middleware.rs), add an arm to
   the role match in `resolve_token` that joins `foo_attributes` for
   this role and builds `CallerAttributes::Foo(...)`.

6. **Dispatch in the evaluator.** In
   [src/acl/evaluator.rs](../src/acl/evaluator.rs)'s lambda dispatch
   match, add a `(Some(CallerAttributes::Foo(info)), Some(name))` arm
   that looks up `lambdas::foo::lookup(name)`.

7. **Allow admin writes for the new role.** In
   [src/admin/registry.rs](../src/admin/registry.rs), add the role to
   `ensure_lambda_known`. In [src/admin/members.rs](../src/admin/members.rs),
   extend `load_member` and `upsert_member` to read/write the new
   attribute table.

8. **Document it.** Update the schema table in §4, the lambdas list in
   §5, and the role-attribute summary in this section. Add typed
   attribute fields to the validation cheatsheet in
   [admin-api.md](admin-api.md).

A test that round-trips a `PUT /admin/members/:eoa { "role": "foo",
"attributes": {...} }` and asserts the resolved `CallerAttributes`
variant exercises every layer touched above.

## 7. Admin management & root trust

`ADMIN_EOAS` (comma-separated env var) is the trust anchor. On every
boot, [`admin::reconcile_seed_admins`](../src/admin/mod.rs) upserts each
EOA into `members` with `role = admin` and drops any leftover
`user_attributes` row for the EOA. Properties:

- Operator owns root via deploy config.
- All admin auth still goes through the wallet-signature flow → no
  long-lived admin shared secret.
- **Break-glass**: even if every other admin is demoted in DB, the
  next restart re-promotes everyone in `ADMIN_EOAS`. The seed cannot be
  locked out by DB state alone.
- **Key rotation**: edit `ADMIN_EOAS` and restart. To revoke active
  sessions for a compromised EOA without rotation:
  `DELETE /admin/members/:eoa/tokens`.

## 8. Authentication & tokens

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
   the consumed nonce, upserts the `members` row (default role = `user`)
   and a `user_attributes` row with defaults `kyc=false, blacklisted=false`
   in a single transaction, then mints a fresh 32-byte random token,
   stores its sha256, and returns the token plaintext (only time the
   plaintext exists outside the wallet). `kyc` and `blacklisted` are
   admin-managed only; this endpoint never modifies them after the
   initial row creation.

Token format: 64 hex chars. Stored hashed: `sha256(token)`. Tokens are
opaque; rotating them is "revoke + re-sign-in".

## 9. Filtering via `debug_traceCall`

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

## 10. Module layout

| Path                                                    | Responsibility                                                              |
| ------------------------------------------------------- | --------------------------------------------------------------------------- |
| [src/main.rs](../src/main.rs)                           | bin entry — calls `lib::run()`                                              |
| [src/lib.rs](../src/lib.rs)                             | module aggregator + bootstrap                                               |
| [src/config.rs](../src/config.rs)                       | env-based config                                                            |
| [src/db.rs](../src/db.rs)                               | sqlx pool, migration, `now_unix`                                            |
| [src/state.rs](../src/state.rs)                         | `AppState { config, pool, upstream }`                                       |
| [src/upstream.rs](../src/upstream.rs)                   | reqwest client to Nethermind                                                |
| [src/error.rs](../src/error.rs)                         | `ApiError` + JSON error response shape                                      |
| [src/server.rs](../src/server.rs)                       | axum router assembly                                                        |
| [src/auth/](../src/auth/)                               | challenge, verify, middleware, `CallerCtx`                                  |
| [src/rpc/](../src/rpc/)                                 | JSON-RPC dispatcher, method classification                                  |
| [src/rpc/gated_methods.rs](../src/rpc/gated_methods.rs) | synthetic selectors for address-parameterized read methods + helpers        |
| [src/acl/evaluator.rs](../src/acl/evaluator.rs)         | allow/deny semantics                                                        |
| [src/acl/registry.rs](../src/acl/registry.rs)           | DB reads of rules + entries                                                 |
| [src/acl/lambdas/](../src/acl/lambdas/)                 | named lambda registry + examples                                            |
| [src/tracer/mod.rs](../src/tracer/mod.rs)               | `debug_traceCall` + frame walk + raw tx decode                              |
| [src/admin/](../src/admin/)                             | `/admin/*` routes (roles, members, registry, lambdas) + seed reconciliation |

## 11. Out of scope (deferred) — engineering notes

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

## 12. Verification matrix

| Layer                                     | Test type   | Where                                                                                                                                                                 |
| ----------------------------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ACL truth table                           | unit        | [src/acl/evaluator.rs](../src/acl/evaluator.rs) — 12 cases covering allow/deny × entry × lambda × anonymous + unknown lambda.                                         |
| `require_kyc` lambda                      | unit        | [src/acl/lambdas/user/require_kyc.rs](../src/acl/lambdas/user/require_kyc.rs) — typed kyc bool flips.                                                                 |
| `erc20_self_only` lambda                  | unit        | [src/acl/lambdas/user/erc20_self_only.rs](../src/acl/lambdas/user/erc20_self_only.rs) — balanceOf and allowance × self/other × unknown selector × malformed calldata. |
| Auth-layer typed attribute resolution     | unit        | [src/auth/middleware.rs](../src/auth/middleware.rs) — admin / user / drift-default tests.                                                                             |
| Admin gate                                | integration | [tests/admin.rs](../tests/admin.rs) — `401` / `403` / `200` matrix on `/admin/roles`.                                                                                 |
| Lambda listing                            | integration | [tests/admin.rs](../tests/admin.rs) — `GET /admin/registry/lambdas` grouped by role.                                                                                  |
| Removed role create/delete endpoints      | integration | [tests/admin.rs](../tests/admin.rs).                                                                                                                                  |
| Typed user upsert + admin lambda reject   | integration | [tests/admin.rs](../tests/admin.rs).                                                                                                                                  |
| Restart reconciliation                    | integration | [tests/admin.rs](../tests/admin.rs) — pre-seed a user with attrs → reconcile → assert promoted + attrs dropped.                                                       |
| Non-eth namespace rejection               | integration | [tests/admin.rs](../tests/admin.rs).                                                                                                                                  |
| Synthetic selector encoding               | unit        | [src/rpc/gated_methods.rs](../src/rpc/gated_methods.rs) — round-trip lookup + ABI layout for target/slot.                                                             |
| Default gated-method policy               | integration | [tests/gated_methods.rs](../tests/gated_methods.rs) — self-allowed, other-EOA denied, contract free, anonymous denied.                                                |
| Admin override on gated method            | integration | [tests/gated_methods.rs](../tests/gated_methods.rs) — deny rule on contract.                                                                                          |
| Method-name selector normalization        | integration | [tests/gated_methods.rs](../tests/gated_methods.rs) — `POST /admin/registry/rules` with `function_selector: "eth_getBalance"` stores `0xff010001`.                    |
| `GET /admin/registry/synthetic-selectors` | integration | [tests/gated_methods.rs](../tests/gated_methods.rs).                                                                                                                  |
| End-to-end RPC + tracer                   | manual      | run against real Nethermind; see operator-guide.md.                                                                                                                   |
