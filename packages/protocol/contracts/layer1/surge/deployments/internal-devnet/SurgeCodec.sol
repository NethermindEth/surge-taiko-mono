// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { Codec } from "../../../core/impl/Codec.sol";
import { FinalityGadgetCodec } from "../../features/codec/FinalityGadgetCodec.sol";

/// @title SurgeCodec
/// @notice Surge codec implementation for internal-devnet deployment
/// @dev Inherits from FinalityGadgetCodec providing Surge-specific encoding/decoding
/// @custom:security-contact security@nethermind.io
contract SurgeCodec is Codec, FinalityGadgetCodec { }
