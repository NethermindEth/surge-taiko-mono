// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { IBridge } from "../../../shared/bridge/IBridge.sol";
import { FlashLoanToken } from "./FlashLoanToken.sol";

/// @title FlashLoanCallbackL1
/// @notice L1 counterpart of the L2→L1→L2 flash loan demo. Receives a bridge
///         message from the L2 executor representing tokens sent from L2, adds
///         a fee (funded by freely minting on L1), and bridges the total back
///         to L2 — producing the L1→L2 return signal that the L2 executor's
///         pre-loaded `processMessage` call expects to consume as a fast signal.
/// @dev Supply bookkeeping on L1: mint on receive, mint the fee, burn before
///      sending back. Net change on L1 is zero every flash loan.
/// @custom:security-contact security@nethermind.io
contract FlashLoanCallbackL1 {
    // ---------------------------------------------------------------
    // Immutable Configuration
    // ---------------------------------------------------------------

    address public immutable bridge;
    uint64 public immutable l2ChainId;
    FlashLoanToken public immutable token;

    /// @notice Fee in basis points (1 bp = 0.01%). 100 bp = 1%.
    uint256 public constant FEE_BPS = 100;
    uint256 public constant BPS_DENOM = 10_000;

    /// @notice Gas limit used for the return bridge message.
    uint32 public constant RETURN_GAS_LIMIT = 1_000_000;

    // ---------------------------------------------------------------
    // State
    // ---------------------------------------------------------------

    address public l2Executor;
    address public admin;

    // ---------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------

    event L2ExecutorSet(address indexed l2Executor);
    event FlashLoanServed(uint256 amountIn, uint256 amountOut, bytes32 returnMsgHash);

    // ---------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------

    error ONLY_BRIDGE();
    error ONLY_ADMIN();
    error INVALID_SENDER();
    error L2_EXECUTOR_NOT_SET();

    // ---------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------

    constructor(address _bridge, uint64 _l2ChainId, address _token, address _admin) {
        bridge = _bridge;
        l2ChainId = _l2ChainId;
        token = FlashLoanToken(_token);
        admin = _admin;
    }

    // ---------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------

    function setL2Executor(address _l2Executor) external {
        if (msg.sender != admin) revert ONLY_ADMIN();
        l2Executor = _l2Executor;
        emit L2ExecutorSet(_l2Executor);
    }

    // ---------------------------------------------------------------
    // Bridge Callback (invoked by L1 Bridge.processMessage of the L2→L1 signal)
    // ---------------------------------------------------------------

    /// @notice Called when the L2 executor's outbound bridge message is processed on L1.
    ///         The payload is the amount of tokens the L2 side "sent". We add a fee,
    ///         then bridge `amount + fee` back to the L2 executor.
    /// @param _data abi.encode(uint256 amount, address l2ExecutorRecipient)
    function onMessageInvocation(bytes calldata _data) external payable {
        if (msg.sender != bridge) revert ONLY_BRIDGE();
        if (l2Executor == address(0)) revert L2_EXECUTOR_NOT_SET();

        IBridge.Context memory ctx = IBridge(bridge).context();
        if (ctx.from != l2Executor) revert INVALID_SENDER();

        (uint256 amount, address l2Recipient) = abi.decode(_data, (uint256, address));

        uint256 fee = (amount * FEE_BPS) / BPS_DENOM;
        uint256 total = amount + fee;

        // Supply bookkeeping on L1:
        //   1. Mint `amount` to represent tokens arriving from L2
        //   2. Mint `fee` (the profit the L1 side is choosing to add)
        //   3. Burn `total` to represent tokens leaving back to L2
        // Net effect on L1 supply is zero.
        token.mint(address(this), total);
        token.burn(total);

        // Send return message to L2 executor. Payload carries the total amount.
        bytes memory returnPayload = abi.encode(total, l2Recipient);
        bytes memory msgData = abi.encodeWithSignature("onMessageInvocation(bytes)", returnPayload);

        IBridge.Message memory message = IBridge.Message({
            id: 0,
            fee: 0,
            gasLimit: RETURN_GAS_LIMIT,
            from: address(0),
            srcChainId: 0,
            srcOwner: address(this),
            destChainId: l2ChainId,
            destOwner: l2Executor,
            to: l2Executor,
            value: 0,
            data: msgData
        });

        (bytes32 msgHash,) = IBridge(bridge).sendMessage(message);
        emit FlashLoanServed(amount, total, msgHash);
    }

    // ---------------------------------------------------------------
    // Receive ETH (bridge messages may carry value; unused here)
    // ---------------------------------------------------------------

    receive() external payable { }
}
