# Admin API reference

Full endpoint reference for the privacy-proxy admin surface. Intended
for whoever builds the operator UI. Every endpoint below is gated by
the admin middleware (`Authorization: Bearer <admin-token>` resolving
to `role = admin`). No admin endpoint is anonymous.

## Authentication

Admins sign in via the same `/auth/challenge` → `/auth/verify` flow as
users (see [wallet-integration.md](wallet-integration.md)). The
resulting token is identical in shape; what makes it an "admin token"
is that its EOA's `users.role_id` resolves to the `admin` role —
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
| 2   | `GET`    | `/admin/registry/rules/:id`                   | Get a rule + entries.                                  |
| 3   | `POST`   | `/admin/registry/rules`                       | Create rule + entries.                                 |
| 4   | `PUT`    | `/admin/registry/rules/:id`                   | Replace rule mode + entries.                           |
| 5   | `DELETE` | `/admin/registry/rules/:id`                   | Delete rule (cascades entries).                        |
| 6   | `POST`   | `/admin/registry/rules/:id/entries`           | Add one entry.                                         |
| 7   | `PUT`    | `/admin/registry/rules/:id/entries/:entry_id` | Update entry's lambda_name.                            |
| 8   | `DELETE` | `/admin/registry/rules/:id/entries/:entry_id` | Remove one entry.                                      |
| 9   | `GET`    | `/admin/registry/lambdas`                     | List in-build lambdas.                                 |
| 10  | `GET`    | `/admin/roles`                                | List roles.                                            |
| 11  | `POST`   | `/admin/roles`                                | Create role.                                           |
| 12  | `DELETE` | `/admin/roles/:id`                            | Delete role.                                           |
| 13  | `GET`    | `/admin/users`                                | List users.                                            |
| 14  | `GET`    | `/admin/users/:eoa`                           | Get one user.                                          |
| 15  | `PUT`    | `/admin/users/:eoa`                           | Upsert user (sets role + caller_info).                 |
| 16  | `PUT`    | `/admin/users/:eoa`                           | Promote/demote admin (= 15 with role: "admin").        |
| 17  | `DELETE` | `/admin/users/:eoa`                           | Delete user (cascades tokens).                         |
| 18  | `DELETE` | `/admin/users/:eoa/tokens`                    | Revoke all active sessions.                            |
| 19  | `GET`    | `/admin/registry/synthetic-selectors`         | List the synthetic selector for each gated RPC method. |

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
- Each `role` must exist in `roles`.
- Each `lambda_name`, if set, must exist in the in-build lambda registry.

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

Returns every lambda compiled into the running binary. Used to
populate UI dropdowns when authoring entries.

Response:

```json
[
  {
    "name": "require_kyc",
    "description": "Allow only callers whose stored caller_info has `kyc: true`.",
    "expected_selector": null
  },
  {
    "name": "transfer_under_limit",
    "description": "For ERC-20 transfer(address,uint256): require amount <= caller_info.max_transfer (decimal string in wei).",
    "expected_selector": "0xa9059cbb"
  },
  {
    "name": "target_in_caller_allowlist",
    "description": "For gated address-parameterized reads (eth_getBalance et al.): allow if params[0] is present in caller_info.balance_allowlist (array of lowercase hex addresses).",
    "expected_selector": null
  },
  {
    "name": "erc20_self_only",
    "description": "For ERC-20 balanceOf(address) and allowance(address,address): allow only when the queried account (balanceOf) or owner (allowance) equals the caller's EOA injected as caller_info.eoa.",
    "expected_selector": null
  }
]
```

`expected_selector` is advisory — the UI should warn (but not block)
if an operator attaches a lambda to a rule whose `function_selector`
differs.

### 19. `GET /admin/registry/synthetic-selectors`

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

**Example 2 — allow a user to read another EOA's balance based on a
per-user allowlist.**

```bash
# 1. Put the readable EOAs into the user's caller_info.
curl -X PUT "$PROXY/admin/users/0xUserEoa..." \
  -H "authorization: Bearer $ADMIN_TOKEN" \
  -H "content-type: application/json" \
  -d '{
    "role": "user",
    "caller_info": { "balance_allowlist": ["0xOtherEoa..."] }
  }'

# 2. Install an allow rule on the other EOA, with the allowlist lambda.
curl -X POST "$PROXY/admin/registry/rules" \
  -H "authorization: Bearer $ADMIN_TOKEN" \
  -H "content-type: application/json" \
  -d '{
    "contract_address": "0xOtherEoa...",
    "function_selector": "eth_getBalance",
    "mode": "allow",
    "entries": [
      { "role": "user", "lambda_name": "target_in_caller_allowlist" }
    ]
  }'
```

Now any user-role caller whose `caller_info.balance_allowlist`
contains `0xOtherEoa...` can read its balance through the proxy.

## Gating ERC-20 reads — walkthrough

`balanceOf(address)` and `allowance(address,address)` are ordinary
contract reads, gated through the standard `(contract_address,
function_selector)` rules. Pair both selectors with the in-build
`erc20_self_only` lambda to enforce _"a user can only inspect their
own balance / allowance on this token"_ — without writing a custom
predicate.

The lambda reads `caller_info.eoa`, which the proxy's auth layer
auto-injects on every authenticated request, so admins don't have to
populate it manually in `caller_info_json`.

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

### 10. `GET /admin/roles`

```json
[
  { "id": 1, "name": "admin" },
  { "id": 2, "name": "user" }
]
```

### 11. `POST /admin/roles`

```json
{ "name": "premium_user" }
```

`409` if a role with that name exists.

### 12. `DELETE /admin/roles/:id`

`204` on success. `409` if any user or rule entry references the role.

## Users

### 13. `GET /admin/users`

Query params:

- `role` — filter by role name.
- `limit` — default 100, max 1000.
- `offset` — default 0.

```json
[
  {
    "eoa_address": "0x...",
    "role": "user",
    "caller_info": { "kyc": true, "max_transfer": "1000000000000000000" },
    "created_at": 1715000000
  }
]
```

### 14. `GET /admin/users/:eoa`

Returns the single user record or `404`.

### 15 & 16. `PUT /admin/users/:eoa`

Body:

```json
{
  "role": "user",
  "caller_info": { "kyc": true }
}
```

Upserts the row. To promote an EOA to admin (capability 16), use
`"role": "admin"`. To demote, set any non-admin role.

Note: if the EOA is listed in `ADMIN_EOAS`, demoting it persists until
the next restart, at which point env-driven reconciliation re-promotes
it. To remove a seed admin permanently, edit the env var and restart.

### 17. `DELETE /admin/users/:eoa`

`204` on success. Cascades to `auth_tokens` (the EOA's active sessions
are revoked).

### 18. `DELETE /admin/users/:eoa/tokens`

Revoke every active token for the EOA without deleting the user. Use
when you want the user to re-authenticate without losing their role
assignment.

Response:

```json
{ "revoked": 3 }
```

## Validation summary (cheatsheet)

| Field                  | Format                                                    |
| ---------------------- | --------------------------------------------------------- |
| EOA / contract address | `0x` + 40 lowercase hex (server normalizes from any case) |
| Function selector      | `0x` + 8 lowercase hex                                    |
| Mode                   | `"allow"` or `"deny"`                                     |
| Role name              | non-empty string, unique                                  |
| Lambda name            | must appear in the response of capability 9               |
| `caller_info`          | any JSON value; opaque to the proxy                       |

## What the API doesn't do

- No bulk import/export.
- No diff or dry-run endpoints — writes are immediate.
- No pagination cursors; offset-based pagination only.
- No filtering by role on the rules list — fetch all and filter
  client-side, or filter by `contract`.
- No write endpoint for lambdas (in-build only).
- No audit log endpoint (admin actions are not currently recorded).
