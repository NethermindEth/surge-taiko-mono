import { getDefaultConfig } from "@rainbow-me/rainbowkit";
import { http } from "wagmi";
import { defineChain } from "viem";

const projectId =
  (import.meta.env.VITE_WALLETCONNECT_PROJECT_ID as string | undefined) ||
  "privacy-proxy-local";

export const APP_NAME =
  (import.meta.env.VITE_APP_NAME as string | undefined) ||
  "Surge Privacy Proxy";

/**
 * Public URL of the privacy-proxy as seen by the wallet (NOT the Vite dev
 * proxy target). Wallets like Ambire will route every JSON-RPC call to this
 * URL and detect it as a Surge proxy via `GET /info`. In dev it defaults
 * to the local proxy; in deployed envs set `VITE_PUBLIC_PROXY_URL`.
 */
const PROXY_RPC_URL =
  (import.meta.env.VITE_PUBLIC_PROXY_URL as string | undefined) ||
  (import.meta.env.VITE_PROXY_URL as string | undefined) ||
  "http://localhost:8080";

/**
 * Chain id the proxy advertises via `/info`. Must match the upstream the
 * proxy is configured to forward to. Defaults to Surge devnet L2
 * (chainId 763374). Override via `VITE_PROXY_CHAIN_ID` if the proxy is
 * configured to forward to a different upstream.
 */
const PROXY_CHAIN_ID = Number(
  (import.meta.env.VITE_PROXY_CHAIN_ID as string | undefined) || "763374",
);

const PROXY_NATIVE_SYMBOL =
  (import.meta.env.VITE_PROXY_NATIVE_SYMBOL as string | undefined) || "ETH";
const PROXY_NATIVE_NAME =
  (import.meta.env.VITE_PROXY_NATIVE_NAME as string | undefined) || "Ether";

const PROXY_NETWORK_NAME =
  (import.meta.env.VITE_PROXY_NETWORK_NAME as string | undefined) ||
  "Surge Private Network";

/**
 * Custom chain pointing at the privacy proxy. The native currency mirrors
 * the upstream — change via env if the proxy is wired to a non-ETH-native
 * upstream. wagmi's connect button may issue auto-RPCs (`eth_chainId`,
 * `eth_blockNumber`) against this URL; those are in the proxy's public
 * whitelist and pass through, so connection works. The `eth_getBalance`
 * widget in RainbowKit will silently fail (`-32001`) since wagmi doesn't
 * carry the proxy bearer — `showBalance={false}` on `ConnectButton`
 * (already set in `LoginGate.tsx`) hides the empty widget.
 */
export const surgePrivateNetwork = defineChain({
  id: PROXY_CHAIN_ID,
  name: PROXY_NETWORK_NAME,
  nativeCurrency: {
    name: PROXY_NATIVE_NAME,
    symbol: PROXY_NATIVE_SYMBOL,
    decimals: 18,
  },
  rpcUrls: {
    default: { http: [PROXY_RPC_URL] },
  },
});

// The proxy's challenge is a plain EIP-191 personal_sign string — wallets
// can be on any chain when signing. The chain we register here is purely
// what the wallet UI shows after connect; the sign-in flow itself does
// not depend on it.
export const wagmiConfig = getDefaultConfig({
  appName: APP_NAME,
  projectId,
  chains: [surgePrivateNetwork],
  transports: {
    [surgePrivateNetwork.id]: http(PROXY_RPC_URL),
  },
});
