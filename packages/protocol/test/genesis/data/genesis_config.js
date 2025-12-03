"use strict";
const ADDRESS_LENGTH = 40;

// Surge: make owner configurable
const ownerAddress =
  process.env.CONTRACT_OWNER || "0xDf08F82De32B8d460adbE8D72043E3a7e25A3B39";

module.exports = {
  // Owner address of the pre-deployed L2 contracts.
  contractOwner: ownerAddress,
  // Chain ID of the Surge L2 network.
  // Surge: make chainId configurable
  chainId: parseInt(process.env.L2_CHAINID) || 167,
  l1ChainId: parseInt(process.env.L1_CHAINID) || 31337,
  // Account address and pre-mint ETH amount as key-value pairs.
  seedAccounts: [{ [ownerAddress]: 1000 }],
  // Owner Chain ID, Security Council, and Timelock Controller
  ownerSecurityCouncil: ownerAddress,
  ownerTimelockController: ownerAddress,
  get contractAddresses() {
    return {
      // ============ Implementations ============
      // Shared Contracts
      BridgeImpl: getConstantAddress(`0${this.chainId}`, 1),
      ERC20VaultImpl: getConstantAddress(`0${this.chainId}`, 2),
      ERC721VaultImpl: getConstantAddress(`0${this.chainId}`, 3),
      ERC1155VaultImpl: getConstantAddress(`0${this.chainId}`, 4),
      SignalServiceImpl: getConstantAddress(`0${this.chainId}`, 5),
      SharedResolverImpl: getConstantAddress(`0${this.chainId}`, 6),
      BridgedERC20Impl: getConstantAddress(`0${this.chainId}`, 10096),
      BridgedERC721Impl: getConstantAddress(`0${this.chainId}`, 10097),
      BridgedERC1155Impl: getConstantAddress(`0${this.chainId}`, 10098),
      RegularERC20: getConstantAddress(`0${this.chainId}`, 10099),
      // Rollup Contracts
      TaikoAnchorImpl: getConstantAddress(`0${this.chainId}`, 10001),
      RollupResolverImpl: getConstantAddress(`0${this.chainId}`, 10002),
      BondManagerImpl: getConstantAddress(`0${this.chainId}`, 10003),
      AnchorForkRouterImpl: getConstantAddress(`0${this.chainId}`, 10004),
      // ============ Proxies ============
      // Shared Contracts
      Bridge: getConstantAddress(this.chainId, 1),
      ERC20Vault: getConstantAddress(this.chainId, 2),
      ERC721Vault: getConstantAddress(this.chainId, 3),
      ERC1155Vault: getConstantAddress(this.chainId, 4),
      SignalService: getConstantAddress(this.chainId, 5),
      SharedResolver: getConstantAddress(this.chainId, 6),
      // Rollup Contracts
      TaikoAnchor: getConstantAddress(this.chainId, 10001),
      RollupResolver: getConstantAddress(this.chainId, 10002),
      BondManager: getConstantAddress(this.chainId, 10003),
    };
  },
  // L2 EIP-1559 baseFee calculation related fields.
  param1559: {
    gasExcess: 1,
  },
  // Option to pre-deploy an ERC-20 token.
  predeployERC20: process.env.PREDEPLOY_ERC20 === "true" || true,
  // Bond-related configurations
  livenessBond: process.env.LIVENESS_BOND || "128000000000000000000",
  provabilityBond: process.env.PROVABILITY_BOND || "128000000000000000000",
  withdrawalDelay: parseInt(process.env.WITHDRAWAL_DELAY) || 3600,
  minBond: parseInt(process.env.MIN_BOND) || 0,
  bondToken:
    process.env.BOND_TOKEN || "0x0000000000000000000000000000000000000000",
  remoteSignalService:
    process.env.REMOTE_SIGNAL_SERVICE ||
    "0x0000000000000000000000000000000000000000",
  pacayaTaikoAnchor:
    process.env.PACAYA_TAIKO_ANCHOR ||
    "0x0000000000000000000000000000000000000000",
};

function getConstantAddress(prefix, suffix) {
  return `0x${prefix}${"0".repeat(
    ADDRESS_LENGTH - String(prefix).length - String(suffix).length,
  )}${suffix}`;
}
