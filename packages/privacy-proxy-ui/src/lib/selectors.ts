/**
 * Local table of well-known 4-byte selectors → function signatures.
 * Used by:
 *   - SelectorPicker  : dropdown of common picks in rule creation.
 *   - LambdasPage     : pretty-print a lambda's `expected_selectors`
 *                       chip as a signature instead of raw hex.
 *
 * Keep this list short and high-signal. The proxy does not consult it —
 * selectors on rules are always stored as 4-byte hex.
 */
export interface CommonSelector {
  signature: string;
  selector: string;
  tag?: string;
}

export const COMMON_SELECTORS: readonly CommonSelector[] = [
  { signature: "transfer(address,uint256)", selector: "0xa9059cbb", tag: "ERC-20" },
  { signature: "transferFrom(address,address,uint256)", selector: "0x23b872dd", tag: "ERC-20/721" },
  { signature: "approve(address,uint256)", selector: "0x095ea7b3", tag: "ERC-20/721" },
  { signature: "balanceOf(address)", selector: "0x70a08231", tag: "ERC-20/721" },
  { signature: "allowance(address,address)", selector: "0xdd62ed3e", tag: "ERC-20" },
  { signature: "totalSupply()", selector: "0x18160ddd", tag: "ERC-20" },
  { signature: "permit(address,address,uint256,uint256,uint8,bytes32,bytes32)", selector: "0xd505accf", tag: "EIP-2612" },
  { signature: "deposit()", selector: "0xd0e30db0", tag: "WETH" },
  { signature: "withdraw(uint256)", selector: "0x2e1a7d4d", tag: "WETH" },
  { signature: "safeTransferFrom(address,address,uint256)", selector: "0x42842e0e", tag: "ERC-721" },
];

const BY_HEX = new Map(
  COMMON_SELECTORS.map((s) => [s.selector.toLowerCase(), s]),
);

/** Look up a common selector entry by 4-byte hex. Case-insensitive. */
export function findCommonSelector(hex: string | undefined): CommonSelector | undefined {
  if (!hex) return undefined;
  return BY_HEX.get(hex.toLowerCase());
}
