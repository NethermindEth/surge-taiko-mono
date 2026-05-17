# privacy-proxy-ui

Enterprise admin panel for the [`privacy-proxy`](../privacy-proxy) package.
Talks to the proxy's `/admin/*` and `/auth/*` endpoints over plain HTTP after a
one-time wallet-signature sign-in.

## Local development

```bash
# 1. Terminal A — start the privacy-proxy pointed at a public Sepolia RPC.
cd ../privacy-proxy
export UPSTREAM_URL=https://ethereum-sepolia-rpc.publicnode.com
export ADMIN_EOAS=0x3e95dFbBaF6B348396E6674C7871546dCC568e56
export DATABASE_URL=sqlite://./privacy-proxy.db
export BIND_ADDR=0.0.0.0:8080
export AUTH_DOMAIN=privacy-proxy
cargo run --release

# 2. Terminal B — start the UI.
cd ../privacy-proxy-ui
cp .env.example .env
pnpm install
pnpm dev
```

Open <http://localhost:5173>, connect the seed admin EOA in your wallet, sign
the EIP-191 challenge → admin dashboard renders.

## Ports

| Service          | Port | URL                   |
| ---------------- | ---- | --------------------- |
| privacy-proxy    | 8080 | http://localhost:8080 |
| privacy-proxy-ui | 5173 | http://localhost:5173 |

The Vite dev server proxies every `/api/*` request through to the proxy at
`VITE_PROXY_URL` (default `http://localhost:8080`), stripping the `/api`
prefix. Browsers see same-origin requests, so no CORS preflight is needed —
the proxy package stays untouched.

## Auth

The proxy's tokens are 32-byte opaque hex strings, not actual JWTs, but they
behave the same way: bearer header, server-stored sha256 hash, TTL'd. The UI
treats them as session credentials, persists `{ token, expiresAt, eoa }` in
`localStorage` under `privacy-proxy-ui:auth`, and re-signs only when the token
expires.

## Features

See [the plan / spec](../../packages-privacy-proxy-i-do-not-vast-zebra.md) and
[`docs/admin-api.md`](../privacy-proxy/docs/admin-api.md) for the full
capability list. The UI exposes all 16 admin capabilities of the proxy plus
the sign-in / sign-out flow.

## Stack

React 18 · Vite 5 · TypeScript 5 · Tailwind CSS · wagmi + RainbowKit + viem ·
@tanstack/react-query · react-router-dom · react-hot-toast.
