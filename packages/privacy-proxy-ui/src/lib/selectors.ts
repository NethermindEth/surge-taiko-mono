export interface CommonSelectorParam {
  name: string;
  type: string;
}

export interface CommonSelector {
  signature: string;
  selector: string;
  tag?: string;
  params: CommonSelectorParam[];
}

export const COMMON_SELECTORS: readonly CommonSelector[] = [
  {
    signature: "transfer(address,uint256)",
    selector: "0xa9059cbb",
    tag: "ERC-20",
    params: [
      { name: "to", type: "address" },
      { name: "amount", type: "uint256" },
    ],
  },
  {
    signature: "transferFrom(address,address,uint256)",
    selector: "0x23b872dd",
    tag: "ERC-20/721",
    params: [
      { name: "from", type: "address" },
      { name: "to", type: "address" },
      { name: "amount", type: "uint256" },
    ],
  },
  {
    signature: "approve(address,uint256)",
    selector: "0x095ea7b3",
    tag: "ERC-20/721",
    params: [
      { name: "spender", type: "address" },
      { name: "amount", type: "uint256" },
    ],
  },
  {
    signature: "balanceOf(address)",
    selector: "0x70a08231",
    tag: "ERC-20/721",
    params: [{ name: "account", type: "address" }],
  },
  {
    signature: "allowance(address,address)",
    selector: "0xdd62ed3e",
    tag: "ERC-20",
    params: [
      { name: "owner", type: "address" },
      { name: "spender", type: "address" },
    ],
  },
  {
    signature: "totalSupply()",
    selector: "0x18160ddd",
    tag: "ERC-20",
    params: [],
  },
  {
    signature:
      "permit(address,address,uint256,uint256,uint8,bytes32,bytes32)",
    selector: "0xd505accf",
    tag: "EIP-2612",
    params: [
      { name: "owner", type: "address" },
      { name: "spender", type: "address" },
      { name: "value", type: "uint256" },
      { name: "deadline", type: "uint256" },
      { name: "v", type: "uint8" },
      { name: "r", type: "bytes32" },
      { name: "s", type: "bytes32" },
    ],
  },
  {
    signature: "deposit()",
    selector: "0xd0e30db0",
    tag: "WETH",
    params: [],
  },
  {
    signature: "withdraw(uint256)",
    selector: "0x2e1a7d4d",
    tag: "WETH",
    params: [{ name: "amount", type: "uint256" }],
  },
  {
    signature: "safeTransferFrom(address,address,uint256)",
    selector: "0x42842e0e",
    tag: "ERC-721",
    params: [
      { name: "from", type: "address" },
      { name: "to", type: "address" },
      { name: "tokenId", type: "uint256" },
    ],
  },
];

const BY_HEX = new Map(
  COMMON_SELECTORS.map((s) => [s.selector.toLowerCase(), s]),
);

export function findCommonSelector(
  hex: string | undefined,
): CommonSelector | undefined {
  if (!hex) return undefined;
  return BY_HEX.get(hex.toLowerCase());
}
