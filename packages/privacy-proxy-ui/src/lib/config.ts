import { getDefaultConfig } from "@rainbow-me/rainbowkit";
import { sepolia } from "wagmi/chains";
import { http } from "wagmi";

const projectId =
  (import.meta.env.VITE_WALLETCONNECT_PROJECT_ID as string | undefined) ||
  "privacy-proxy-local";

export const APP_NAME =
  (import.meta.env.VITE_APP_NAME as string | undefined) ||
  "Surge Privacy Proxy";

// The proxy's challenge is a plain EIP-191 personal_sign string — wallets can
// be on any chain when signing. We still register Sepolia as the default so
// the user has something to switch to in the wallet UI and to display chain
// status in the connect button. The chain doesn't affect sign-in.
export const wagmiConfig = getDefaultConfig({
  appName: APP_NAME,
  projectId,
  chains: [sepolia],
  transports: {
    [sepolia.id]: http(),
  },
});
