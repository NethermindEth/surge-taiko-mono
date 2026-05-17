# Admin API reference

Full endpoint reference for the privacy-proxy admin surface. Intended
for whoever builds the operator UI. Every endpoint below is gated by
the admin middleware (`Authorization: Bearer <admin-token>` resolving
to `role = admin`). No admin endpoint is anonymous.

## Authentication

Admins sign in via the same `/auth/challenge` → `/auth/verify` flow as
regular users (see [wallet-integration.md](wallet-integration.md)). The
resulting token is identical in shape; what makes it an "admin token"
is that its EOA's `members.role_id` resolves to the `admin` role —
seeded via the `ADMIN_EOAS` env var or promoted by an existing admin.

| HTTP  | Meaning                                      |
| ----- | -------------------------------------------- |
| `401` | Missing / invalid / expired token.           |
| `403` | Token resolves to a role other than `admin`. |
| `400` | Bad input (see per-endpoint validation).     |
| `404` | Resource not found.                          |
| `409` | Uniqueness conflict (e.g. duplicate rule).   |

Error body shape (every code):

```json
{ "error": { "code": "string_id", "message": "human readable" } }
```

## Capability index

| #   | Method   | Endpoint                                      | Purpose                                                |
| --- | -------- | --------------------------------------------- | ------------------------------------------------------ |
| 1   | `GET`    | `/admin/registry/rules`                       | List access rules.                                     |
| 2   | `GET`    | `/admin/registry/rules/:id`                   | Get a rule + its entries.                              |
| 3   | `POST`   | `/admin/registry/rules`                       | Create a rule + entries.                               |
| 4   | `PUT`    | `/admin/registry/rules/:id`                   | Replace a rule's mode + entries.                       |
| 5   | `DELETE` | `/admin/registry/rules/:id`                   | Delete a rule (cascades to entries).                   |
| 6   | `POST`   | `/admin/registry/rules/:id/entries`           | Add one entry to a rule.                               |
| 7   | `PUT`    | `/admin/registry/rules/:id/entries/:entry_id` | Update an entry's `lambda_name`.                       |
| 8   | `DELETE` | `/admin/registry/rules/:id/entries/:entry_id` | Remove one entry.                                      |
| 9   | `GET`    | `/admin/registry/lambdas`                     | List in-build lambdas grouped by role.                 |
| 10  | `GET`    | `/admin/registry/synthetic-selectors`         | List the synthetic selector for each gated RPC method. |
| 11  | `GET`    | `/admin/roles`                                | Enumerate the roles declared by this build.            |
| 12  | `GET`    | `/admin/members`                              | List members (admins + users).                         |
| 13  | `GET`    | `/admin/members/:eoa`                         | Get one member.                                        |
| 14  | `PUT`    | `/admin/members/:eoa`                         | Upsert a member with their role and typed attributes.  |
| 15  | `DELETE` | `/admin/members/:eoa`                         | Delete a member (cascades to tokens).                  |
| 16  | `DELETE` | `/admin/members/:eoa/tokens`                  | Revoke a member's active sessions.                     |

## Access registry

The `function_selector` field on rule writes accepts either:

- a 4-byte hex value (e.g. `"0xa9059cbb"`) — for ordinary contract function rules, as before; or
- a JSON-RPC method name (e.g. `"eth_getBalance"`) — for the gated address-parameterized reads. The server rewrites it to the synthetic selector (see capability 19). Subsequent reads return the synthetic hex value.

The set of accepted method names is the response of capability 19.

### 1. `GET /admin/registry/rules`

Query params (all optional):

- `contract` — filter by contract address (0x-prefixed 20-byte hex).
- `limit` — default 100, max 1000.
- `offset` — default 0.

Response:

```json
[
  {
    "id": 1,
    "contract_address": "0x...",
    "function_selector": "0xa9059cbb",
    "mode": "allow",
    "entries": [
      { "id": 7, "role": "user", "lambda_name": "transfer_under_limit" }
    ]
  }
]
```

### 2. `GET /admin/registry/rules/:id`

Returns one rule with its entries (same shape as the array element above).

### 3. `POST /admin/registry/rules`

Body:

```json
{
  "contract_address": "0x...",
  "function_selector": "0xa9059cbb",
  "mode": "allow",
  "entries": [
    { "role": "user", "lambda_name": "transfer_under_limit" },
    { "role": "premium_user" }
  ]
}
```

Validation:

- `mode` ∈ {`allow`, `deny`}.
- `contract_address` normalized to lowercase 20-byte hex.
- `function_selector` normalized to lowercase 4-byte hex (`0x` + 8 chars).
- `(contract_address, function_selector)` must be globally unique.
- Each `role` must exist in `roles` (and therefore in [src/roles.rs](../src/roles.rs)).
- Each `lambda_name`, if set, must exist in **the target role's** lambda
  registry. Note: the `admin` role has no lambda registry — attaching a
  `lambda_name` to an admin entry returns `400`.

Response: `201` with the created `RuleView`.

### 4. `PUT /admin/registry/rules/:id`

Body:

```json
{
  "mode": "deny",
  "entries": [{ "role": "user" }]
}
```

Replaces mode and entries atomically. Use this when the operator UI's
"save" button submits the full edited rule.

### 5. `DELETE /admin/registry/rules/:id`

`204 No Content` on success. Cascades to `access_rule_entries`.

### 6. `POST /admin/registry/rules/:id/entries`

Add a single entry without rewriting the rule.

Body:

```json
{ "role": "user", "lambda_name": "require_kyc" }
```

Response: `201` with `{ id, role, lambda_name }`.

Conflict (`409`) if the role already has an entry under this rule (the
`(rule_id, role_id)` unique constraint).

### 7. `PUT /admin/registry/rules/:id/entries/:entry_id`

Body:

```json
{ "lambda_name": "transfer_under_limit" }
```

`lambda_name: null` is allowed and clears the lambda. The entry's role
cannot be changed via this endpoint — delete and re-add to change role.

### 8. `DELETE /admin/registry/rules/:id/entries/:entry_id`

`204` on success.

### 9. `GET /admin/registry/lambdas`

Returns every lambda compiled into the running binary, grouped by the
role its attribute struct targets. Used to populate UI dropdowns when
authoring entries — the available lambdas for an entry depend on the
entry's role.

Response:

```json
[
  { "role": "admin", "lambdas": [] },
  {
    "role": "user",
    "lambdas": [
      {
        "name": "require_kyc",
        "description": "Allow only callers whose stored attributes have kyc=true.",
        "expected_selector": null
      },
      {
        "name": "erc20_self_only",
        "description": "For ERC-20 balanceOf(address) and allowance(address,address): allow only when the queried account (balanceOf) or owner (allowance) equals the caller's EOA.",
        "expected_selector": null
      }
    ]
  }
]
```

`expected_selector` is advisory — the UI should warn (but not block)
if an operator attaches a lambda to a rule whose `function_selector`
differs. Empty groups (e.g. `admin`) are returned for shape stability.

### 10. `GET /admin/registry/synthetic-selectors`

The proxy gates a small set of address-parameterized read methods
(`eth_getBalance`, `eth_getTransactionCount`, `eth_getCode`,
`eth_getStorageAt`, `eth_getProof`) using synthetic 4-byte selectors
in the reserved `0xff______` range. This endpoint exposes the
mapping so the UI can render method-name dropdowns when authoring
rules and translate between names and the on-disk selectors.

Response:

```json
[
  { "method": "eth_getBalance", "selector": "0xff010001" },
  { "method": "eth_getTransactionCount", "selector": "0xff010002" },
  { "method": "eth_getCode", "selector": "0xff010003" },
  { "method": "eth_getStorageAt", "selector": "0xff010004" },
  { "method": "eth_getProof", "selector": "0xff010005" }
]
```

## Gating balance-style reads — walkthrough

These methods follow the same `(contract_address, function_selector)`
schema as contract calls. The `contract_address` is the **target
address** being queried, and the `function_selector` is the synthetic
value (or the method name as a synonym).

**Default policy without any rule installed:**

- An EOA caller can read its own state. Querying any other EOA returns
  `-32001`. Contract targets are free. Anonymous callers are always
  denied.

**Example 1 — restrict who can read a treasury contract's balance.**

Operator wants only the `admin` role (not regular `user`) to be able
to call `eth_getBalance` on `0xTreasury`.

```bash
curl -X POST "$PROXY/admin/registry/rules" \
  -H "authorization: Bearer $ADMIN_TOKEN" \
  -H "content-type: application/json" \
  -d '{
    "contract_address": "0xTreasury...",
    "function_selector": "eth_getBalance",
    "mode": "deny",
    "entries": [ { "role": "user" } ]
  }'
```

Result: user-role tokens get `-32001` on `eth_getBalance(0xTreasury)`;
admin tokens are unaffected by this rule and still allowed (admins
have no `entry` under `mode = deny`, so the deny doesn't apply to
them).

## Gating ERC-20 reads — walkthrough

`balanceOf(address)` and `allowance(address,address)` are ordinary
contract reads, gated through the standard `(contract_address,
function_selector)` rules. Pair both selectors with the in-build
`erc20_self_only` lambda to enforce _"a user can only inspect their
own balance / allowance on this token"_ — without writing a custom
predicate.

The lambda reads `UserCallerInfo.eoa`, which the proxy's auth layer
sets from the resolved token (it is the row's primary key in `members`)
— admins never set the EOA manually.

**Example — gate one ERC-20 token `T`.**

```bash
# balanceOf rule
curl -X POST "$PROXY/admin/registry/rules" \
  -H "authorization: Bearer $ADMIN_TOKEN" \
  -H "content-type: application/json" \
  -d '{
    "contract_address": "0xT...",
    "function_selector": "0x70a08231",
    "mode": "allow",
    "entries": [
      { "role": "user", "lambda_name": "erc20_self_only" }
    ]
  }'

# allowance rule
curl -X POST "$PROXY/admin/registry/rules" \
  -H "authorization: Bearer $ADMIN_TOKEN" \
  -H "content-type: application/json" \
  -d '{
    "contract_address": "0xT...",
    "function_selector": "0xdd62ed3e",
    "mode": "allow",
    "entries": [
      { "role": "user", "lambda_name": "erc20_self_only" }
    ]
  }'
```

After this, a user-role caller with EOA `A` can `eth_call`:

- `T.balanceOf(A)` → forwarded.
- `T.balanceOf(B)` → `-32001` with `data.detail = "LambdaRejected"`.
- `T.allowance(A, anySpender)` → forwarded.
- `T.allowance(B, anySpender)` → `-32001`.

Repeat the two POSTs for every token you want to gate. The same
lambda name is reused across all of them.

## Roles

### 11. `GET /admin/roles`

Enumerates the role names recognized by this build. The set is the
static `ROLES` slice in [src/roles.rs](../src/roles.rs); the API only
reads it. Adding or removing a role is a code change — see
[system-design.md §6](system-design.md) for the procedure.

```json
[
  { "id": 1, "name": "admin" },
  { "id": 2, "name": "user" }
]
```

## Members

A _member_ is any authenticated account. Both admins and users are
members; the `role` field distinguishes them. The `members` table holds
identity (`eoa_address`, `role`, `created_at`) and role-specific
attributes live in role-named tables (e.g. `user_attributes`).

### 12. `GET /admin/members`

Query params:

- `role` — filter by role name.
- `limit` — default 100, max 1000.
- `offset` — default 0.

Response items are typed by role. Members whose role carries
attributes return them in the `attributes` object; members of a role
without attributes (e.g. `admin`) return `attributes: null`.

```json
[
  {
    "eoa_address": "0x...",
    "role": "user",
    "attributes": { "kyc": true, "blacklisted": false },
    "created_at": 1715000000
  },
  {
    "eoa_address": "0x...",
    "role": "admin",
    "attributes": null,
    "created_at": 1715000000
  }
]
```

### 13. `GET /admin/members/:eoa`

Returns the single member record or `404`. Same shape as one element above.

### 14. `PUT /admin/members/:eoa`

Upserts a member's role and (when applicable) attributes.

User member:

```json
{
  "role": "user",
  "attributes": { "kyc": true, "blacklisted": false }
}
```

- The `attributes` object is optional, and individual fields within it
  are optional. Omitted fields preserve the row's current value; on
  first insert, omitted fields default to `false`.

Admin member:

```json
{ "role": "admin" }
```

The admin role has no attributes — `attributes` must be absent or
`null`. Submitting attributes alongside `role: "admin"` returns `400
admin role does not accept attributes`. Switching an EOA from `user` to
`admin` removes its `user_attributes` row.

Note: if the EOA is listed in `ADMIN_EOAS`, demoting it via this
endpoint persists until the next restart, at which point env-driven
reconciliation re-promotes it. To take a seed admin off the list
permanently, edit the env var and restart.

### 15. `DELETE /admin/members/:eoa`

`204` on success. Cascades to `auth_tokens` (the EOA's active sessions
are revoked) and to any role-specific attribute row.

### 16. `DELETE /admin/members/:eoa/tokens`

Revoke every active token for the EOA without deleting the member. Use
when you want the member to re-authenticate without losing their role
assignment.

Response:

```json
{ "revoked": 3 }
```

## Validation summary (cheatsheet)

| Field                  | Format                                                                      |
| ---------------------- | --------------------------------------------------------------------------- |
| EOA / contract address | `0x` + 40 lowercase hex (server normalizes from any case)                   |
| Function selector      | `0x` + 8 lowercase hex                                                      |
| Mode                   | `"allow"` or `"deny"`                                                       |
| Role name              | one of the names returned by capability 11                                  |
| Lambda name            | must appear under the target role in the response of capability 9           |
| `attributes` (user)    | object with bool fields `kyc`, `blacklisted` (each optional on partial PUT) |
| `attributes` (admin)   | must be absent or `null`                                                    |

## What the API doesn't do

- No bulk import/export.
- No diff or dry-run endpoints — writes are immediate.
- No pagination cursors; offset-based pagination only.
- No filtering by role on the rules list — fetch all and filter
  client-side, or filter by `contract`.
- No write endpoint for lambdas (in-build only).
- No audit log endpoint (admin actions are not currently recorded).
