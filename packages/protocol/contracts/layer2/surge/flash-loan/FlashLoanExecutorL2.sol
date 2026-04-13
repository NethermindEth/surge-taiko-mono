// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { IBridge } from "../../../shared/bridge/IBridge.sol";
import { FlashLoanPool, IFlashLoanReceiver } from "./FlashLoanPool.sol";
import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { ERC20Burnable } from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

interface IMintableBurnable {
    function mint(address to, uint256 amount) external;
    function burn(uint256 amount) external;
}

/// @title FlashLoanExecutorL2
/// @notice Orchestrator on L2 for the L2→L1→L2 synchronous flash loan.
/// @dev Execution flow (all inside one L2 transaction):
///      1. User calls `execute(amount, returnMessage)` — `returnMessage` is the
///         pre-computed L1→L2 bridge return Message (produced by the builder
///         via simulation). It hashes to the fast signal slot that the anchor
///         injects and the inbox records as a required return signal.
///      2. We store `returnMessage` in transient storage, call pool.flashLoan.
///      3. Pool transfers tokens to us, invokes `onFlashLoan`.
///      4. Inside `onFlashLoan`:
///         a. Burn borrowed tokens (simulates leaving L2)
///         b. `Bridge.sendMessage` to L1 callback with the borrowed amount
///            → creates the L2→L1 signal S_out that L1's `processMessage` of
///              the outbound leg consumes during the L1 multicall
///         c. `Bridge.processMessage(returnMessage, "")` — fast signal path.
///            This invokes our `onMessageInvocation` callback with the L1-supplied
///            amount (original + fee). We mint that total back.
///         d. Repay the pool (original `amount`); fee stays as profit.
/// @custom:security-contact security@nethermind.io
contract FlashLoanExecutorL2 is IFlashLoanReceiver {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------
    // Immutable Configuration
    // ---------------------------------------------------------------

    address public immutable bridge;
    uint64 public immutable l1ChainId;
    IERC20 public immutable tokenERC20;
    IMintableBurnable public immutable token;
    FlashLoanPool public immutable pool;

    /// @notice Gas limit for the outbound L2→L1 bridge message.
    uint32 public constant OUTBOUND_GAS_LIMIT = 1_000_000;

    // ---------------------------------------------------------------
    // State
    // ---------------------------------------------------------------

    address public l1Callback;
    address public admin;

    // ---------------------------------------------------------------
    // Per-Flash-Loan Storage
    // ---------------------------------------------------------------
    //
    // The pre-computed L1→L2 return Message is passed into `execute` by the
    // builder (via calldata) and stashed here for `onFlashLoan` to consume
    // via `Bridge.processMessage`. Cleared at the end of `execute`.
    // (L2 EVM version may be pre-Cancun — no transient storage available.)
    bytes private _pendingReturnMessageData;
    bool private _pendingReturnMessageSet;

    // Active flash loan bookkeeping — tracks the amount/recipient per tx so
    // `onFlashLoan` knows what to send to L1, and `onMessageInvocation` knows
    // how much to mint back. Cleared at the end of `execute`.
    address private _activeBeneficiary;
    uint256 private _activeAmount;

    // ---------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------

    event L1CallbackSet(address indexed l1Callback);
    event FlashLoanInitiated(address indexed beneficiary, uint256 amount);
    event FlashLoanCompleted(address indexed beneficiary, uint256 amount, uint256 profit);
    event OutboundMessageSent(bytes32 msgHash, uint256 amount);
    event ReturnProcessed(uint256 totalMinted);

    // ---------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------

    error ONLY_ADMIN();
    error ONLY_POOL();
    error ONLY_BRIDGE();
    error INVALID_SENDER();
    error L1_CALLBACK_NOT_SET();
    error NO_PENDING_RETURN_MESSAGE();
    error NO_ACTIVE_LOAN();
    error INSUFFICIENT_RETURN_AMOUNT();

    // ---------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------

    constructor(
        address _bridge,
        uint64 _l1ChainId,
        address _token,
        address _pool,
        address _admin
    ) {
        bridge = _bridge;
        l1ChainId = _l1ChainId;
        tokenERC20 = IERC20(_token);
        token = IMintableBurnable(_token);
        pool = FlashLoanPool(_pool);
        admin = _admin;
    }

    // ---------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------

    function setL1Callback(address _l1Callback) external {
        if (msg.sender != admin) revert ONLY_ADMIN();
        l1Callback = _l1Callback;
        emit L1CallbackSet(_l1Callback);
    }

    // ---------------------------------------------------------------
    // Entry Point
    // ---------------------------------------------------------------

    /// @notice Initiates a synchronous L2→L1→L2 flash loan.
    /// @param _amount Amount of tokens to borrow from the L2 pool.
    /// @param _beneficiary Address that receives the fee profit.
    /// @param _returnMessage Pre-computed L1→L2 Bridge.Message that the L1 callback
    ///        will produce. Its hash must match the fast-signal slot injected by the
    ///        anchor and committed as a required return signal in the inbox.
    function execute(
        uint256 _amount,
        address _beneficiary,
        IBridge.Message calldata _returnMessage
    )
        external
    {
        if (l1Callback == address(0)) revert L1_CALLBACK_NOT_SET();

        // Stash the pre-computed return message and active bookkeeping for the callback.
        _pendingReturnMessageData = abi.encode(_returnMessage);
        _pendingReturnMessageSet = true;
        _activeBeneficiary = _beneficiary;
        _activeAmount = _amount;

        emit FlashLoanInitiated(_beneficiary, _amount);

        // Borrow from the pool — triggers onFlashLoan below.
        pool.flashLoan(_amount);

        // After return, forward any remaining balance (the fee profit) to beneficiary.
        uint256 profit = tokenERC20.balanceOf(address(this));
        if (profit > 0) {
            tokenERC20.safeTransfer(_beneficiary, profit);
        }

        emit FlashLoanCompleted(_beneficiary, _amount, profit);

        // Cleanup (since this L2 profile may lack transient storage).
        delete _pendingReturnMessageData;
        _pendingReturnMessageSet = false;
        _activeBeneficiary = address(0);
        _activeAmount = 0;
    }

    // ---------------------------------------------------------------
    // Flash Loan Callback (from pool)
    // ---------------------------------------------------------------

    function onFlashLoan(
        address,
        /*_token*/
        uint256 _amount
    )
        external
        override
    {
        if (msg.sender != address(pool)) revert ONLY_POOL();
        if (!_pendingReturnMessageSet) revert NO_PENDING_RETURN_MESSAGE();

        // (a) Burn borrowed tokens — conceptually "sending" them to L1.
        token.burn(_amount);

        // (b) Send the outbound bridge message to the L1 callback. This writes
        //     an L2→L1 signal slot; the L1 multicall's processMessage consumes it
        //     to invoke FlashLoanCallbackL1, which produces the required return signal.
        bytes memory outboundPayload = abi.encode(_amount, address(this));
        bytes memory outboundData =
            abi.encodeWithSignature("onMessageInvocation(bytes)", outboundPayload);

        IBridge.Message memory outboundMsg = IBridge.Message({
            id: 0,
            fee: 0,
            gasLimit: OUTBOUND_GAS_LIMIT,
            from: address(0),
            srcChainId: 0,
            srcOwner: address(this),
            destChainId: l1ChainId,
            destOwner: l1Callback,
            to: l1Callback,
            value: 0,
            data: outboundData
        });

        (bytes32 outHash,) = IBridge(bridge).sendMessage(outboundMsg);
        emit OutboundMessageSent(outHash, _amount);

        // (c) Consume the L1→L2 return message via fast signal (empty proof).
        //     The pre-computed Message is decoded and its hash verified against
        //     the fast signals set by the anchor. processMessage invokes our
        //     own onMessageInvocation, which mints the total back.
        IBridge.Message memory returnMsg = abi.decode(_pendingReturnMessageData, (IBridge.Message));
        IBridge(bridge).processMessage(returnMsg, "");

        // (d) Repay the pool. Only the borrowed amount goes back; fee profit
        //     remains in this contract and is forwarded to the beneficiary
        //     by `execute` after the flashLoan call returns.
        tokenERC20.safeTransfer(address(pool), _amount);
    }

    // ---------------------------------------------------------------
    // Bridge Callback (from L2 Bridge.processMessage of the return message)
    // ---------------------------------------------------------------

    /// @notice Called by the L2 bridge when the L1→L2 return Message is processed.
    /// @dev msg.sender is the L2 bridge; ctx.from must be the L1 callback contract.
    function onMessageInvocation(bytes calldata _data) external payable {
        if (msg.sender != bridge) revert ONLY_BRIDGE();
        IBridge.Context memory ctx = IBridge(bridge).context();
        if (ctx.from != l1Callback) revert INVALID_SENDER();

        if (_activeAmount == 0) revert NO_ACTIVE_LOAN();

        (uint256 total,) = abi.decode(_data, (uint256, address));
        if (total < _activeAmount) revert INSUFFICIENT_RETURN_AMOUNT();

        // Mint the returned tokens back on L2 (including the fee from L1).
        token.mint(address(this), total);

        emit ReturnProcessed(total);
    }
}
