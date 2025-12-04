// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { FinalityGadgetCodecSimple } from "../../features/codec/FinalityGadgetCodecSimple.sol";

/// @title SurgeCodecSimple
/// @notice Surge codec implementation for internal-devnet deployment
/// @dev Inherits from FinalityGadgetCodecSimple providing Surge-specific encoding/decoding
/// @custom:security-contact security@nethermind.io
contract SurgeCodecSimple is FinalityGadgetCodecSimple { }

