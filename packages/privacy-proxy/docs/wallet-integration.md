# Wallet integration guide

For a wallet developer integrating with a `privacy-proxy` endpoint.
This is the **only doc you need** to make a wallet talk to the proxy —
everything else (admin API, system internals) is unrelated to wallet
work.

## TL;DR

1. Replace the chain's normal RPC URL with the proxy URL.
2. On first use (or whenever the token is missing/expired), run the
   sign-in flow:
   - `GET /auth/challenge?address={eoa}` → get a message.
   - Have the wallet sign it (personal_sign).
   - `POST /auth/verify { address, signature }` → receive a token.
3. Attach `Authorization: Bearer <token>` to **every** JSON-RPC request
   sent to the proxy.
4. On HTTP `401`, re-run the sign-in flow.
5. On JSON-RPC error `code: -32001`, surface the denial to the user.

The wire shape of all other JSON-RPC requests is identical to a normal
Ethereum node. No params change.

## Endpoints

### `GET /auth/challenge?address={eoa}`

Returns the EIP-191 message the wallet must sign to prove it controls
the EOA.

```json
HTTP 200
{
  "message": "your-domain sign-in\nAddress: 0x...\nNonce: 8f3a...",
  "expires_at": 1715600000
}
```

Errors: `400` for an invalid address.

The `message` is a single UTF-8 string. The wallet MUST sign it as-is
with `personal_sign` (EIP-191). Do not modify, trim, or wrap it.

### `POST /auth/verify`

Submit the signature; receive a token.

Request:

```json
{
  "address": "0x{eoa}",
  "signature": "0x{65-byte hex}"
}
```

Response:

```json
HTTP 200
{
  "token":      "abc123...",      // 64 hex chars
  "expires_at": 1716200000        // unix seconds
}
```

Errors:

- `400` invalid address / signature encoding / expired challenge.
- `401` recovered signer does not match the supplied address.
- `404` no pending challenge for this address (call `/auth/challenge`
  first).

The wallet should:

1. Store the token in secure local state (keychain, IndexedDB with
   encryption, etc.).
2. Track `expires_at` to know when to re-run the flow proactively.

### POST `/` (JSON-RPC)

Standard Ethereum JSON-RPC. The only addition vs a vanilla node is the
required header:

```
Authorization: Bearer <token>
Content-Type: application/json
```

Example:

```json
POST /
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "eth_call",
  "params": [{ "to": "0x...", "data": "0x..." }, "latest"]
}
```

Most methods that don't carry call data (`eth_blockNumber`,
`eth_getTransactionReceipt`, `eth_chainId`, …) are always passed
through. **Exception**: a small set of address-parameterized reads —
`eth_getBalance`, `eth_getProof`, `eth_getTransactionCount`,
`eth_getCode`, `eth_getStorageAt` — are also gated. By default a user
can only query their **own** EOA; querying any other EOA returns
`-32001` with `data.detail = "DefaultEoaSelfOnly"`. Contract targets
are free unless the operator has installed a restriction. Methods that
carry call data may be filtered too; see below.

#### Access-denied response

When a call (or any of its internal calls) is forbidden, the response
is a standard JSON-RPC error with code `-32001`:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32001,
    "message": "access denied",
    "data": {
      "contract": "0x...",
      "selector": "0xa9059cbb",
      "detail": "NotInAllowList"
    }
  }
}
```

`contract` and `selector` identify the call frame (or, for gated
read methods, the target address and the synthetic 4-byte selector
in the reserved `0xff______` range) that triggered the denial. The
denial may originate from an internal call, not the top-level
transaction. `detail` is a stable identifier suitable for switching
on in UX code; common values:

| `detail`                    | Meaning                                                                                                                                |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| `NotInAllowList`            | Rule is `mode = allow` and the caller's role has no entry.                                                                             |
| `InDenyList`                | Rule is `mode = deny` and the caller's role has an entry.                                                                              |
| `LambdaRejected`            | The attached lambda returned `false`.                                                                                                  |
| `AnonymousAgainstGatedCall` | No auth token, request needed one.                                                                                                     |
| `DefaultEoaSelfOnly`        | Gated read method targeting another EOA, with no admin override. The user should retry with their own EOA, or ask an admin for access. |
| `UnknownLambda`             | Server-side configuration error (fail-closed). Report to operator.                                                                     |

Recommended UX: render the denial inline as a soft error with an
"unauthorized" badge; offer a link to whatever onboarding flow (KYC,
allowlist application, etc.) the operator uses to grant access.

#### Other relevant errors

| HTTP / RPC   | Meaning                                                       |
| ------------ | ------------------------------------------------------------- |
| HTTP `401`   | Token missing, invalid, or expired. Re-run sign-in.           |
| RPC `-32600` | Invalid request (e.g. batch JSON-RPC — not supported).        |
| RPC `-32601` | Method not in `eth_` namespace (everything else is rejected). |
| RPC `-32602` | Bad params (e.g. invalid raw transaction hex).                |
| RPC `-32700` | JSON parse error.                                             |
| RPC `-32000` | Internal proxy error (tracer or upstream failure).            |

## End-to-end (TypeScript pseudocode)

```ts
class PrivacyProxyClient {
  constructor(
    private url: string,
    private wallet: Wallet,
  ) {}

  private token?: string;
  private expiresAt?: number;

  async ensureToken() {
    const now = Math.floor(Date.now() / 1000);
    if (this.token && this.expiresAt! > now + 60) return;

    const eoa = await this.wallet.getAddress();
    const chal = await fetch(`${this.url}/auth/challenge?address=${eoa}`).then(
      (r) => r.json(),
    );
    const sig = await this.wallet.signMessage(chal.message);
    const verify = await fetch(`${this.url}/auth/verify`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ address: eoa, signature: sig }),
    }).then((r) => r.json());
    this.token = verify.token;
    this.expiresAt = verify.expires_at;
  }

  async rpc(method: string, params: unknown[]) {
    await this.ensureToken();
    const res = await fetch(`${this.url}/`, {
      method: "POST",
      headers: {
        "content-type": "application/json",
        authorization: `Bearer ${this.token}`,
      },
      body: JSON.stringify({ jsonrpc: "2.0", id: 1, method, params }),
    });
    if (res.status === 401) {
      this.token = undefined;
      return this.rpc(method, params); // one retry
    }
    const body = await res.json();
    if (body.error?.code === -32001) {
      throw new AccessDeniedError(body.error.data);
    }
    if (body.error) throw new Error(body.error.message);
    return body.result;
  }
}
```

## Things to avoid

- **Do not** batch JSON-RPC requests — the proxy rejects arrays.
- **Do not** subscribe via WebSocket — only HTTP is supported.
- **Do not** call non-`eth_` namespace methods (`net_version`,
  `web3_clientVersion`, `debug_*`, `trace_*`). Hard-code their answers
  client-side if you need them, or expose them via a separate
  unauthenticated read endpoint.
- **Do not** persist the `message` from `/auth/challenge` — nonces are
  one-shot and short-lived. Always fetch a fresh challenge each sign-in.
- **Do not** assume the token is bound to a session cookie — it's a
  bearer token. Treat it like an API key; never log it.
