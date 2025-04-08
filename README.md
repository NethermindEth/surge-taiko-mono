<p align="center">
  <img src="./.github/SurgeLogoOnly.svg" width="80" alt="Logo for Surge" />
</p>

<h1 align="center">
  Surge Rollup
</h1>

<p align="center">
  A based rollup template of the Taiko Alethia stack. 
</p>

## Overview

This repository is a collection of Surge-related tools and components designed to simplify the setup and operation of layer-2 (L2) rollups. It includes everything needed to build, deploy, and monitor L2 solutions. From protocol specifications and smart contracts to monitoring tools and user interfaces, this repository provides a comprehensive suite of resources to run Surge Rollup.

## Documentation

Get started with Surge:

- [Surge docs](https://docs.surge.wtf) — High-level concepts, guides, resources, and reference pages for getting started.
- [Protocol specs](./packages/protocol/docs/README.md) — In-depth specifications of the Taiko protocol for deeper understanding.
- [Smart contracts](./packages/protocol/contracts/) — Taiko protocol smart contracts, fully documented with NatSpec.

## Repository Structure

- **`packages/balance-monitor`**: Tracks and monitors balances across different accounts and contracts to ensure proper fund management.
- **`packages/blobstorage`**: Manages blob data storage and retrieval, including APIs for interacting with blob hashes and metadata.
- **`packages/bridge-ui`**: Provides a user interface for bridging assets between layer-1 and layer-2 networks.
- **`packages/docs-site`**: Hosts the documentation site for the Surge Rollup, including guides, FAQs, and protocol details.
- **`packages/eventindexer`**: Indexes blockchain events for querying and analysis, providing APIs for accessing event data.
- **`packages/fork-diff`**: A tool for comparing and visualizing differences between forks of the Taiko protocol.
- **`packages/guardian-prover-health-check`**: Monitors the health and uptime of guardian provers, ensuring their availability and reliability.
- **`packages/guardian-prover-health-check-ui`**: A Svelte-based UI for visualizing the health and status of guardian provers.
- **`packages/monitors`**: Provides monitoring tools for tracking the health and performance of the rollup system.
- **`packages/nfts`**: Manages NFT-related contracts, scripts, and metadata for the Taiko ecosystem.
- **`packages/protocol`**: Contains the core protocol specifications and smart contracts. This package defines the rules and mechanisms for the rollup, including consensus, fraud proofs, and other protocol-level details.
- **`packages/relayer`**: Facilitates cross-chain communication by relaying messages and transactions between layer-1 and layer-2.
- **`packages/snaefell-ui`**: A UI library for building user interfaces related to the Surge Rollup ecosystem.
- **`packages/supplementary-contracts`**: Contains additional smart contracts that extend the functionality of the core protocol.
- **`packages/taiko-client`**: Implements the client-side logic for interacting with the Taiko protocol, including proving, proposing, and syncing blocks.
- **`packages/taikoon-ui`**: Implements the frontend for Taikoon NFTs, including collection rendering and interaction.
- **`packages/ui-lib`**: A library of reusable UI components for building decentralized applications.

## Issues

If you find a bug or have a feature request, please [open an issue](https://github.com/NethermindEth/surge-taiko-mono/issues/new/choose).

## Contributing

Check out [CONTRIBUTING.md](./CONTRIBUTING.md) for details on how to contribute.
