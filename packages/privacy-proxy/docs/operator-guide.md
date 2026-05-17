# Operator guide

For the person deploying and running privacy-proxy. Covers configuration,
bootstrap, key rotation, and restart semantics.

## Configuration (env vars)

| Var                  | Required             | Default                     | Purpose                                                                                                                                            |
| -------------------- | -------------------- | --------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| `UPSTREAM_URL`       | yes                  | —                           | HTTP URL of the upstream execution client (e.g. Nethermind). Must support `debug_traceCall` with `callTracer`.                                     |
| `ADMIN_EOAS`         | strongly recommended | empty                       | Comma-separated EOAs that become admins on every boot. Empty → no seed admins → no one can use admin endpoints until the DB is edited out of band. |
| `BIND_ADDR`          | no                   | `0.0.0.0:8080`              | HTTP listener.                                                                                                                                     |
| `DATABASE_URL`       | no                   | `sqlite://privacy-proxy.db` | sqlx URL. Use `sqlite::memory:` for ephemeral test runs.                                                                                           |
| `CHALLENGE_TTL_SECS` | no                   | `300`                       | Sign-in challenge nonce TTL.                                                                                                                       |
| `TOKEN_TTL_SECS`     | no                   | `604800` (7 days)           | Auth token TTL.                                                                                                                                    |
| `AUTH_DOMAIN`        | no                   | `privacy-proxy`             | Domain string embedded in the sign-in message. Pick something stable per deployment (e.g. `surge.wtf`).                                            |
| `RUST_LOG`           | no                   | `info,privacy_proxy=debug`  | Standard `tracing_subscriber` filter.                                                                                                              |

## First-time bootstrap

1. Pick the EOA(s) that will be the initial admin. The corresponding
   private key must be held by a human or hardware wallet.
2. Set `ADMIN_EOAS` in your deploy config to the comma-separated list
   of these EOAs.
3. Start the binary. It will:
   - Create the SQLite DB file (or use the existing one).
   - Run migrations.
   - Upsert every EOA in `ADMIN_EOAS` to `role = admin`.
   - Start serving on `BIND_ADDR`.
4. From the admin's wallet, hit `/auth/challenge` and `/auth/verify` to
   obtain an admin bearer token (same flow as regular users — see
   [wallet-integration.md](wallet-integration.md)).
5. Use the admin token to populate the registry. See
   [admin-api.md](admin-api.md) for the full surface.

## Key rotation

**Scenario: a seed admin EOA is compromised.**

1. Edit `ADMIN_EOAS` to remove the compromised EOA and add its
   replacement.
2. Restart the binary. On boot:
   - The new EOA is upserted to `role = admin`.
   - The compromised EOA is **not** automatically demoted in DB — that
     row still has whatever role was last assigned to it.
3. From an unaffected admin's session, call
   `PUT /admin/members/:compromisedEoa { "role": "user" }` to demote.
4. Call `DELETE /admin/members/:compromisedEoa/tokens` to revoke any
   tokens already issued.

Until step 3, the compromised EOA's existing tokens (if any) still
work. If you cannot get an unaffected admin online quickly, you can
edit the DB directly (e.g. `sqlite3 privacy-proxy.db 'DELETE FROM
auth_tokens WHERE eoa_address = ...'`).

**Scenario: a non-seed admin EOA is compromised.**

Step 1 and 2 don't apply. Just demote (step 3) and revoke (step 4) via
the API.

## Restart semantics

Every restart re-runs `reconcile_seed_admins`:

- Each EOA in `ADMIN_EOAS` is upserted into `members` with `role =
admin`. If the row exists with a different role, it is **promoted**
  back to admin.
- EOAs that used to be in `ADMIN_EOAS` but are no longer listed are
  **not** demoted automatically. Their last-assigned role persists.
- This gives a "break-glass" property: even if every other admin is
  demoted in the DB by a compromised admin, the next restart restores
  the operator-controlled seed admins. The trust anchor is the deploy
  config, not the DB.

If `ADMIN_EOAS` is empty, the proxy logs a warning at startup and
boots normally — but no one will be able to use any `/admin/*`
endpoint until you either set the env var and restart, or insert an
admin row in the DB by hand.

## Monitoring

The proxy emits structured logs via `tracing`. Useful filters:

- `privacy_proxy::rpc=debug` — every JSON-RPC dispatch decision.
- `privacy_proxy::acl=debug` — ACL evaluations.
- `privacy_proxy::tracer=warn` — only the failures of `debug_traceCall`.
- `privacy_proxy::admin=info` — admin endpoint hits and seed
  reconciliation logs.

There is no Prometheus exporter in the POC. If you need metrics, wire
`metrics`/`metrics-exporter-prometheus` into `axum` via a custom
middleware.

## DB management

SQLite, default file path `./privacy-proxy.db`. To inspect:

```bash
sqlite3 privacy-proxy.db
.tables
SELECT * FROM members;
SELECT * FROM access_rules;
```

To back up: copy the file while the proxy is stopped, or use `sqlite3
... .backup`. To migrate to a new schema, write the migration into
`migrations/` (numbered prefix) and restart; `sqlx::migrate!` picks it
up.

## Operational hazards

- The proxy's `debug_traceCall` simulation is a snapshot of upstream
  state. A `eth_sendRawTransaction` that passes the proxy filter may
  still revert on-chain due to state change between simulate and
  submit; or, conversely, an internal call that _would_ be denied at
  execution time may slip through if state changes invalidate the
  simulated trace path. For high-stakes deployments, mitigate by
  funneling all `eth_sendRawTransaction` through a private mempool /
  builder that runs the same proxy logic post-inclusion.
- The proxy holds no chain state; if the upstream is wrong, the proxy
  is wrong. Rate-limit and authenticate the upstream connection
  separately (firewall, mTLS).
- `ADMIN_EOAS` keys are root. Treat them as you would a Kubernetes
  cluster-admin kubeconfig or an AWS root account: rotate routinely,
  store in a hardware wallet, never paste into a hot environment.
- The DB stores token _hashes_, not plaintext — but a DB compromise
  still leaks every user's EOA, role, and `caller_info` (which may
  include KYC fields). Encrypt the volume at rest if you're storing
  anything sensitive in `caller_info`.
- `eth_getBalance` and the other gated address-parameterized reads
  are filtered (default: EOA self-only), but the EVM `BALANCE` opcode
  used inside an `eth_call` is **not**. A motivated user with
  permission to call `eth_call` against any unrestricted contract — or
  to send custom bytecode in a `to`-less `eth_call` — can still read
  any account's balance through arbitrary contract code. Tighten the
  contract-level rules around `eth_call` if balance privacy matters
  for your deployment. The principled fix is to add `prestateTracer`
  support; out of scope for the POC.

## Upgrading

The proxy is a single binary with embedded migrations. Standard
deploy:

1. Stop the process.
2. Replace the binary.
3. Start the process. Migrations run automatically before the listener
   binds.

There is no online migration story for the POC. Take a short outage.
