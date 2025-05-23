// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./TaikoData.sol";

/// @title ITaikoL1
/// @custom:security-contact security@taiko.xyz
interface ITaikoL1 {
    /// @notice Proposes a Taiko L2 block.
    /// @param _params Block parameters, currently an encoded BlockParams object.
    /// @param _txList txList data if calldata is used for DA.
    /// @return meta_ The metadata of the proposed L2 block.
    /// @return deposits_ The Ether deposits processed.
    function proposeBlock(
        bytes calldata _params,
        bytes calldata _txList
    )
        external
        payable
        returns (TaikoData.BlockMetadata memory meta_, TaikoData.EthDeposit[] memory deposits_);

    /// @notice Proposes a Taiko L2 block (version 2)
    /// @param _params Block parameters, an encoded BlockParamsV2 object.
    /// @param _txList txList data if calldata is used for DA.
    /// @return meta_ The metadata of the proposed L2 block.
    function proposeBlockV2(
        bytes calldata _params,
        bytes calldata _txList
    )
        external
        returns (TaikoData.BlockMetadataV2 memory meta_);

    /// @notice Proposes multiple Taiko L2 blocks (version 2)
    /// @param _paramsArr A list of encoded BlockParamsV2 objects.
    /// @param _txListArr A list of txList.
    /// @return metaArr_ The metadata objects of the proposed L2 blocks.
    function proposeBlocksV2(
        bytes[] calldata _paramsArr,
        bytes[] calldata _txListArr
    )
        external
        returns (TaikoData.BlockMetadataV2[] memory metaArr_);

    /// @notice Proves or contests a block transition.
    /// @param _blockId The index of the block to prove. This is also used to
    /// select the right implementation version.
    /// @param _input An abi-encoded (TaikoData.BlockMetadata, TaikoData.Transition,
    /// TaikoData.TierProof) tuple.
    function proveBlock(uint64 _blockId, bytes calldata _input) external;

    /// @notice Proves or contests multiple block transitions (version 2)
    /// @param _blockIds The indices of the blocks to prove.
    /// @param _inputs An list of abi-encoded (TaikoData.BlockMetadata, TaikoData.Transition,
    /// TaikoData.TierProof) tuples.
    /// @param _batchProof An abi-encoded TaikoData.TierProof that contains the batch/aggregated
    /// proof for the given blocks.
    function proveBlocks(
        uint64[] calldata _blockIds,
        bytes[] calldata _inputs,
        bytes calldata _batchProof
    )
        external;

    /// @notice Verifies up to a certain number of blocks.
    /// @param _maxBlocksToVerify Max number of blocks to verify.
    function verifyBlocks(uint64 _maxBlocksToVerify) external;

    /// @notice Pause block proving.
    /// @param _pause True if paused.
    function pauseProving(bool _pause) external;

    /// @notice Gets the details of a block.
    /// @param _blockId Index of the block.
    /// @return blk_ The block.
    function getBlockV2(uint64 _blockId) external view returns (TaikoData.BlockV2 memory blk_);

    /// @notice Gets the state transition for a specific block.
    /// @param _blockId Index of the block.
    /// @param _tid The transition id.
    /// @return The state transition data of the block.
    function getTransition(
        uint64 _blockId,
        uint32 _tid
    )
        external
        view
        returns (TaikoData.TransitionState memory);

    /// @notice Deposits Ether to be used as bonds.
    function depositBond() external payable;

    /// @notice Withdraws Ether deposited as bonds.
    /// @param _amount The amount of Ether to withdraw.
    function withdrawBond(uint256 _amount) external;

    /// @notice Gets the prover that actually proved a verified block.
    /// @param _blockId The index of the block.
    /// @return The prover's address. If the block is not verified yet, address(0) will be returned.
    function getVerifiedBlockProver(uint64 _blockId) external view returns (address);

    /// @notice Returns the timestamp since when the verification streak has been going
    /// @return The timestamp since which verification streak has been maintained 
    // Surge: Returns `verificationStreakStartedAt` added as stage-2 requirements for surge
    function getVerificationStreakStartAt() external view returns(uint256);

    /// @notice Gets the configuration of the TaikoL1 contract.
    /// @return Config struct containing configuration parameters.
    // Surge: switch to `view` to allow for dynamic chainid
    function getConfig() external view returns (TaikoData.Config memory);
}
