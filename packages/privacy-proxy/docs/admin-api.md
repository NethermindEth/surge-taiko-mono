# Admin API reference

Full endpoint reference for the privacy-proxy admin surface. Every endpoint
below is gated by the admin middleware (`Authorization: Bearer
<admin-token>` resolving to `role = admin`). No admin endpoint is
anonymous.

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
| `409` | Uniqueness / referential conflict.           |

Error body shape:

```json
{ "error": { "code": "string_id", "message": "human readable" } }
```

## Capability index

| Method   | Endpoint                                      | Purpose                                        |
| -------- | --------------------------------------------- | ---------------------------------------------- |
| `GET`    | `/admin/registry/rules`                       | List access rules.                             |
| `GET`    | `/admin/registry/rules/:id`                   | Get a rule + its entries.                      |
| `POST`   | `/admin/registry/rules`                       | Create a rule + entries.                       |
| `PUT`    | `/admin/registry/rules/:id`                   | Replace a rule's mode + entries.               |
| `DELETE` | `/admin/registry/rules/:id`                   | Delete a rule.                                 |
| `POST`   | `/admin/registry/rules/:id/entries`           | Add one entry to a rule.                       |
| `PUT`    | `/admin/registry/rules/:id/entries/:entry_id` | Update an entry's `lambda_id`.                 |
| `DELETE` | `/admin/registry/rules/:id/entries/:entry_id` | Remove one entry.                              |
| `GET`    | `/admin/registry/lambdas`                     | List stored lambdas grouped by role.           |
| `GET`    | `/admin/registry/lambdas/:id`                 | Get one lambda + its comparison rules.         |
| `POST`   | `/admin/registry/lambdas`                     | Create a lambda + its rules.                   |
| `DELETE` | `/admin/registry/lambdas/:id`                 | Delete a lambda (fails if referenced).         |
| `GET`    | `/admin/registry/role-attributes`             | Lambda LHS attribute names per role.           |
| `GET`    | `/admin/registry/synthetic-selectors`         | Synthetic selector ↔ JSON-RPC method mapping. |
| `GET`    | `/admin/roles`                                | Enumerate the code-declared roles.             |
| `GET`    | `/admin/members`                              | List members (admins + users).                 |
| `GET`    | `/admin/members/:eoa`                         | Get one member.                                |
| `PUT`    | `/admin/members/:eoa`                         | Upsert a member with role + typed attributes.  |
| `DELETE` | `/admin/members/:eoa`                         | Delete a member.                               |
| `DELETE` | `/admin/members/:eoa/tokens`                  | Revoke a member's active sessions.             |

## Access registry

The `function_selector` field on rule writes accepts either:

- a 4-byte hex value (e.g. `"0xa9059cbb"`); or
- a JSON-RPC method name (e.g. `"eth_getBalance"`). The server
  rewrites it to the synthetic selector. The accepted method names
  are those returned by `/admin/registry/synthetic-selectors`.

### `GET /admin/registry/rules`

Query params (all optional): `contract` (filter by contract address),
`limit` (default 100, max 1000), `offset` (default 0).

Response:

```json
[
  {
    "id": 1,
    "contract_address": "0x...",
    "function_selector": "0xa9059cbb",
    "mode": "allow",
    "entries": [
      {
        "id": 7,
        "role": "user",
        "lambda_id": 3,
        "lambda_name": "transfer_under_limit"
      }
    ]
  }
]
```

`lambda_name` is echoed for display; writes are by `lambda_id` only.

### `GET /admin/registry/rules/:id`

Returns one rule with the same shape as the array element above.

### `POST /admin/registry/rules`

Body:

```json
{
  "contract_address": "0x...",
  "function_selector": "0xa9059cbb",
  "mode": "allow",
  "entries": [{ "role": "user", "lambda_id": 3 }, { "role": "premium_user" }]
}
```

Validation:

- `mode` ∈ {`allow`, `deny`}.
- `(contract_address, function_selector)` globally unique.
- Each `role` must exist.
- Each `lambda_id`, if set, must exist and belong to the entry's role.

Response: `201` with the created `RuleView`.

### `PUT /admin/registry/rules/:id`

```json
{
  "mode": "deny",
  "entries": [{ "role": "user", "lambda_id": 3 }]
}
```

Replaces mode and entries atomically.

### `DELETE /admin/registry/rules/:id`

`204 No Content` on success. Cascades to entries.

### `POST /admin/registry/rules/:id/entries`

```json
{ "role": "user", "lambda_id": 5 }
```

Response: `201` with `{ id, role, lambda_id, lambda_name }`. Conflict
(`409`) if the role already has an entry under this rule.

### `PUT /admin/registry/rules/:id/entries/:entry_id`

```json
{ "lambda_id": 5 }
```

`lambda_id: null` clears the lambda. The entry's role cannot be
changed via this endpoint — delete and re-add to change role.

### `DELETE /admin/registry/rules/:id/entries/:entry_id`

`204` on success.

## Lambdas

A lambda is a named set of comparison rules, scoped to a single role.
It evaluates to true when **all** of its rules whose `selector` matches
the call's selector evaluate to true; a lambda with zero rules for
the call's selector evaluates vacuously to true. Any error during
evaluation (out-of-bounds calldata, unknown attribute, malformed
literal) returns false.

A comparison rule has:

- `selector` — the 4-byte selector this rule applies to.
- `lhs` — either `{ kind: "calldata", offset }` (a 32-byte word read
  at byte `offset` of the call's input) or `{ kind: "attribute",
name }` (the named role attribute coerced to a 32-byte word).
- `condition` — one of `eq | neq | gt | lt | gte | lte`. Comparisons are
  unsigned 32-byte numeric.
- `rhs` — `{ kind: "tx_origin" }`, `{ kind: "msg_sender" }`, or
  `{ kind: "literal", value: "0x...64 hex" }`. `tx_origin` is the original
  transaction's `from` (the EOA holding the auth token). `msg_sender` is the
  immediate caller of the current frame — for the top-level call it equals
  `tx_origin`; for internal calls it's the parent contract's address. Use
  `msg_sender` when you want to enforce that a call is reached via a specific
  intermediate contract.

### `GET /admin/registry/lambdas`

```json
[
  { "role": "admin", "lambdas": [] },
  {
    "role": "user",
    "lambdas": [
      {
        "id": 7,
        "name": "require_kyc",
        "role": "user",
        "description": "kyc must be true",
        "in_use": true,
        "rules": [
          {
            "id": 12,
            "selector": "0x70a08231",
            "lhs_kind": "attribute",
            "lhs_offset": null,
            "lhs_attribute": "kyc",
            "condition": "eq",
            "rhs_kind": "literal",
            "rhs_value": "0x0000000000000000000000000000000000000000000000000000000000000001"
          }
        ]
      }
    ]
  }
]
```

`in_use: true` means at least one `access_rule_entries` row references
this lambda — deletion is blocked while this is true.

### `GET /admin/registry/lambdas/:id`

Returns the same shape as one element of `lambdas` above. `404` if
unknown.

### `POST /admin/registry/lambdas`

```json
{
  "name": "transfer_self_only",
  "role": "user",
  "description": "Sender of transfer must equal caller's EOA.",
  "rules": [
    {
      "selector": "0xa9059cbb",
      "lhs_kind": "calldata",
      "lhs_offset": 4,
      "condition": "eq",
      "rhs_kind": "tx_origin"
    }
  ]
}
```

Validation:

- `name` non-empty; `(name, role)` unique.
- `role` must be a known role.
- At least one rule.
- For each rule:
  - `selector` is a 4-byte hex or a synthetic method name.
  - If `lhs_kind = "calldata"`: `lhs_offset` is required; `lhs_attribute`
    must be null.
  - If `lhs_kind = "attribute"`: `lhs_attribute` is required and must be
    in `/admin/registry/role-attributes` for `role`; `lhs_offset` must
    be null.
  - `condition` ∈ `{ eq, neq, gt, lt, gte, lte }`.
  - `rhs_kind` ∈ `{ tx_origin, msg_sender, literal }`.
  - If `rhs_kind = "literal"`: `rhs_value` is required, `0x` + 64 hex
    chars; otherwise must be null.

Response: `201` with the created `LambdaView`.

### `DELETE /admin/registry/lambdas/:id`

`204` on success. `409` if the lambda is still referenced by any rule
entry (detach before deleting). `404` if unknown.

### `GET /admin/registry/role-attributes`

Lambda authors choose `lhs_kind = "attribute"` LHSes from this list,
per role.

```json
[
  { "role": "admin", "attributes": ["eoa"] },
  { "role": "user", "attributes": ["eoa", "kyc", "blacklisted"] }
]
```

The values are coerced to 32-byte words at evaluation: bools become
`0x..00` / `0x..01`, addresses are left-padded to 32 bytes.

### `GET /admin/registry/synthetic-selectors`

The proxy gates a small set of address-parameterized read methods
using synthetic selectors in the reserved `0xff______` range.

```json
[
  { "method": "eth_getBalance", "selector": "0xff010001" },
  { "method": "eth_getTransactionCount", "selector": "0xff010002" },
  { "method": "eth_getCode", "selector": "0xff010003" },
  { "method": "eth_getStorageAt", "selector": "0xff010004" },
  { "method": "eth_getProof", "selector": "0xff010005" }
]
```

When evaluating a gated read against a lambda, the proxy synthesizes
calldata in the layout `selector | left-padded target address`
(`getStorageAt` additionally appends the slot at `36..68`). Use byte
offset `4` to compare the target, `36` to compare the storage slot.

## Walkthroughs

### Gating ERC-20 self-reads with a calldata lambda

The "user may only read their own balance / allowance" lambda becomes
two POSTs: create the lambda, then attach it to the rules.

```bash
# 1. Create the lambda.
curl -X POST "$PROXY/admin/registry/lambdas" \
  -H "authorization: Bearer $ADMIN_TOKEN" \
  -H "content-type: application/json" \
  -d '{
    "name": "erc20_self_only",
    "role": "user",
    "description": "balanceOf / allowance must use caller'\''s own EOA",
    "rules": [
      { "selector": "0x70a08231", "lhs_kind": "calldata", "lhs_offset": 4,
        "condition": "eq", "rhs_kind": "tx_origin" },
      { "selector": "0xdd62ed3e", "lhs_kind": "calldata", "lhs_offset": 4,
        "condition": "eq", "rhs_kind": "tx_origin" }
    ]
  }'

# Suppose the response says { "id": 11 }.

# 2. Gate the token.
curl -X POST "$PROXY/admin/registry/rules" \
  -H "authorization: Bearer $ADMIN_TOKEN" \
  -H "content-type: application/json" \
  -d '{
    "contract_address": "0xT...",
    "function_selector": "0x70a08231",
    "mode": "allow",
    "entries": [ { "role": "user", "lambda_id": 11 } ]
  }'

curl -X POST "$PROXY/admin/registry/rules" \
  -H "authorization: Bearer $ADMIN_TOKEN" \
  -H "content-type: application/json" \
  -d '{
    "contract_address": "0xT...",
    "function_selector": "0xdd62ed3e",
    "mode": "allow",
    "entries": [ { "role": "user", "lambda_id": 11 } ]
  }'
```

After this, a user-role caller `A` can `eth_call`:

- `T.balanceOf(A)` → forwarded.
- `T.balanceOf(B)` → `-32001` with `data.detail = "LambdaRejected"`.
- `T.allowance(A, anySpender)` → forwarded.
- `T.allowance(B, anySpender)` → `-32001`.

### Requiring KYC via an attribute lambda

```bash
curl -X POST "$PROXY/admin/registry/lambdas" \
  -H "authorization: Bearer $ADMIN_TOKEN" \
  -H "content-type: application/json" \
  -d '{
    "name": "require_kyc",
    "role": "user",
    "rules": [{
      "selector": "0xa9059cbb",
      "lhs_kind": "attribute", "lhs_attribute": "kyc",
      "condition": "eq",
      "rhs_kind": "literal",
      "rhs_value": "0x0000000000000000000000000000000000000000000000000000000000000001"
    }]
  }'
```

A user-role caller whose stored `user_attributes.kyc = 0` will be
rejected with `LambdaRejected` on `transfer(...)`.

## Roles

### `GET /admin/roles`

```json
[
  { "id": 1, "name": "admin" },
  { "id": 2, "name": "user" }
]
```

The set is `ROLES` in [src/roles.rs](../src/roles.rs); the API only
reads it.

## Members

A _member_ is any authenticated account. Both admins and users are
members; `role` distinguishes them.

### `GET /admin/members`

Query params: `role`, `limit` (default 100, max 1000), `offset`.

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

### `GET /admin/members/:eoa`

Returns the single member record or `404`.

### `PUT /admin/members/:eoa`

User member:

```json
{ "role": "user", "attributes": { "kyc": true, "blacklisted": false } }
```

Admin member:

```json
{ "role": "admin" }
```

Submitting attributes alongside `role: "admin"` returns `400`.
Switching an EOA from `user` to `admin` drops its `user_attributes`
row.

### `DELETE /admin/members/:eoa`

`204` on success.

### `DELETE /admin/members/:eoa/tokens`

```json
{ "revoked": 3 }
```

## Validation summary

| Field                  | Format                                                                      |
| ---------------------- | --------------------------------------------------------------------------- |
| EOA / contract address | `0x` + 40 lowercase hex (server normalizes from any case)                   |
| Function selector      | `0x` + 8 lowercase hex                                                      |
| Mode                   | `"allow"` or `"deny"`                                                       |
| Role name              | one of the names returned by `/admin/roles`                                 |
| Lambda id              | must exist and belong to the entry's role                                   |
| Lambda rule `lhs_kind` | `calldata` (with `lhs_offset`) or `attribute` (with `lhs_attribute`)        |
| Lambda rule `rhs_kind` | `tx_origin`, `msg_sender`, or `literal` (with `rhs_value` as `0x` + 64 hex) |
| `attributes` (user)    | object with bool fields `kyc`, `blacklisted` (each optional on partial PUT) |
| `attributes` (admin)   | must be absent or `null`                                                    |

## What the API doesn't do

- No bulk import/export.
- No diff or dry-run endpoints.
- No pagination cursors; offset-based pagination only.
- No filtering by role on the rules list.
- No PUT for lambdas — delete and recreate to edit.
- No audit log endpoint.
