// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { IBridge } from "../../../shared/bridge/IBridge.sol";
import { IUniswapV2Router02 } from "./IUniswapV2Router02.sol";
import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/// @title CrossChainSwapVaultL1
/// @notice Unified vault on L1 that handles token bridging, cross-chain swaps,
/// and L2 liquidity provisioning — all in a single message per hop.
/// @dev Holds canonical ERC20 tokens. Implements IMessageInvocable pattern.
/// Supports two pipelines:
///   1. L1→L2→L1 swaps that route through the L2 `SimpleDEX` (entry points on this contract).
///   2. L2→L1→L2 swaps that route through an L1 DEX (`l1Router`, Uniswap V2 ABI). The entry
///      point lives on `CrossChainSwapVaultL2`; this contract handles the L1-side leg via
///      `onMessageInvocation`.
/// @custom:security-contact security@taiko.xyz
contract CrossChainSwapVaultL1 {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------
    // Enums (must match L2 vault)
    // ---------------------------------------------------------------

    enum Action {
        BRIDGE,
        SWAP_ETH_TO_TOKEN, // L1→L2→L1, swap on L2 DEX
        SWAP_TOKEN_TO_ETH, // L1→L2→L1, swap on L2 DEX
        ADD_LIQUIDITY,
        REMOVE_LIQUIDITY,
        SWAP_ETH_TO_TOKEN_VIA_L1, // L2→L1→L2, swap on L1 DEX
        SWAP_TOKEN_TO_ETH_VIA_L1 // L2→L1→L2, swap on L1 DEX
    }

    // ---------------------------------------------------------------
    // State Variables
    // ---------------------------------------------------------------

    address public immutable bridge;
    uint64 public immutable l2ChainId;
    IERC20 public immutable swapToken;
    address public immutable admin;
    IUniswapV2Router02 public immutable l1Router;
    address public immutable weth;
    address public l2Vault;

    uint32 public constant GAS_LIMIT = 2_000_000;
    uint32 public constant RETURN_GAS_LIMIT = 2_000_000;

    // ---------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------

    event L2VaultSet(address indexed l2Vault);
    event TokensBridgedToL2(
        address indexed from, address indexed recipient, uint256 amount, bytes32 msgHash
    );
    event SwapETHForTokenInitiated(
        address indexed user, uint256 ethAmount, uint256 minTokenOut, bytes32 msgHash
    );
    event SwapTokenForETHInitiated(
        address indexed user, uint256 tokenAmount, uint256 minETHOut, bytes32 msgHash
    );
    event LiquidityAddedToL2(
        address indexed user, uint256 ethAmount, uint256 tokenAmount, bytes32 msgHash
    );
    event LiquidityRemovedFromL2(address indexed user, bytes32 msgHash);
    event LiquidityRemovalCompleted(
        address indexed recipient, uint256 ethAmount, uint256 tokenAmount
    );
    event SwapETHForTokenCompleted(address indexed recipient, uint256 tokenAmount);
    event SwapTokenForETHCompleted(address indexed recipient, uint256 ethAmount);
    // L2→L1→L2 leg events
    event L1DexSwapETHForTokenExecuted(
        address indexed recipient, uint256 ethIn, uint256 tokenOut, bytes32 returnMsgHash
    );
    event L1DexSwapTokenForETHExecuted(
        address indexed recipient, uint256 tokenIn, uint256 ethOut, bytes32 returnMsgHash
    );

    // ---------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------

    error ONLY_ADMIN();
    error ONLY_BRIDGE();
    error INVALID_SENDER();
    error L2_VAULT_NOT_SET();
    error ZERO_AMOUNT();
    error INSUFFICIENT_TOKEN_BALANCE();
    error INSUFFICIENT_ETH_BALANCE();
    error ETH_TRANSFER_FAILED();
    error UNKNOWN_ACTION();

    // ---------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------

    constructor(
        address _bridge,
        uint64 _l2ChainId,
        address _swapToken,
        address _l1Router,
        address _weth,
        address _admin
    ) {
        bridge = _bridge;
        l2ChainId = _l2ChainId;
        swapToken = IERC20(_swapToken);
        l1Router = IUniswapV2Router02(_l1Router);
        weth = _weth;
        admin = _admin;
    }

    // ---------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------

    function setL2Vault(address _l2Vault) external {
        if (msg.sender != admin) revert ONLY_ADMIN();
        l2Vault = _l2Vault;
        emit L2VaultSet(_l2Vault);
    }

    // ---------------------------------------------------------------
    // Bridge: L1 → L2 (1 message)
    // ---------------------------------------------------------------

    /// @notice Bridge canonical tokens from L1 to L2 (mints bridged tokens on L2)
    function bridgeTokenToL2(uint256 _amount, address _recipient) external {
        if (l2Vault == address(0)) revert L2_VAULT_NOT_SET();
        if (_amount == 0) revert ZERO_AMOUNT();

        // Lock canonical tokens in this vault
        swapToken.safeTransferFrom(msg.sender, address(this), _amount);

        bytes memory data = abi.encode(Action.BRIDGE, _recipient, _amount);
        bytes32 msgHash = _sendMessageToL2(data, 0, msg.sender, GAS_LIMIT);

        emit TokensBridgedToL2(msg.sender, _recipient, _amount, msgHash);
    }

    // ---------------------------------------------------------------
    // Swap: ETH → Token (L1→L2→L1, swap on L2 DEX)
    // ---------------------------------------------------------------

    function swapETHForToken(uint256 _minTokenOut, address _recipient) external payable {
        if (l2Vault == address(0)) revert L2_VAULT_NOT_SET();
        if (msg.value == 0) revert ZERO_AMOUNT();

        bytes memory data =
            abi.encode(Action.SWAP_ETH_TO_TOKEN, msg.sender, _recipient, msg.value, _minTokenOut);
        bytes32 msgHash = _sendMessageToL2(data, msg.value, msg.sender, GAS_LIMIT);

        emit SwapETHForTokenInitiated(msg.sender, msg.value, _minTokenOut, msgHash);
    }

    // ---------------------------------------------------------------
    // Swap: Token → ETH (L1→L2→L1, swap on L2 DEX)
    // ---------------------------------------------------------------

    function swapTokenForETH(
        uint256 _tokenAmount,
        uint256 _minETHOut,
        address _recipient
    )
        external
    {
        if (l2Vault == address(0)) revert L2_VAULT_NOT_SET();
        if (_tokenAmount == 0) revert ZERO_AMOUNT();

        // Lock canonical tokens
        swapToken.safeTransferFrom(msg.sender, address(this), _tokenAmount);

        bytes memory data = abi.encode(
            Action.SWAP_TOKEN_TO_ETH, msg.sender, _recipient, _tokenAmount, _minETHOut
        );
        bytes32 msgHash = _sendMessageToL2(data, 0, msg.sender, GAS_LIMIT);

        emit SwapTokenForETHInitiated(msg.sender, _tokenAmount, _minETHOut, msgHash);
    }

    // ---------------------------------------------------------------
    // Add Liquidity to L2 DEX from L1 (1 message)
    // ---------------------------------------------------------------

    function addLiquidityToL2(uint256 _tokenAmount) external payable {
        if (l2Vault == address(0)) revert L2_VAULT_NOT_SET();
        if (msg.value == 0 || _tokenAmount == 0) revert ZERO_AMOUNT();

        // Lock canonical tokens
        swapToken.safeTransferFrom(msg.sender, address(this), _tokenAmount);

        bytes memory data = abi.encode(Action.ADD_LIQUIDITY, msg.sender, _tokenAmount);
        bytes32 msgHash = _sendMessageToL2(data, msg.value, msg.sender, GAS_LIMIT);

        emit LiquidityAddedToL2(msg.sender, msg.value, _tokenAmount, msgHash);
    }

    function removeLiquidityFromL2() external {
        if (l2Vault == address(0)) revert L2_VAULT_NOT_SET();

        bytes memory data = abi.encode(Action.REMOVE_LIQUIDITY, msg.sender);
        bytes32 msgHash = _sendMessageToL2(data, 0, msg.sender, GAS_LIMIT);

        emit LiquidityRemovedFromL2(msg.sender, msgHash);
    }

    // ---------------------------------------------------------------
    // Bridge Callback (from L2)
    // ---------------------------------------------------------------

    function onMessageInvocation(bytes calldata _data) external payable {
        if (msg.sender != bridge) revert ONLY_BRIDGE();

        IBridge.Context memory ctx = IBridge(bridge).context();
        if (ctx.from != l2Vault) revert INVALID_SENDER();

        Action action = abi.decode(_data, (Action));

        if (action == Action.SWAP_ETH_TO_TOKEN) {
            // L1→L2→L1 completion: release canonical tokens to recipient
            (, address recipient, uint256 tokenAmount) =
                abi.decode(_data, (Action, address, uint256));
            if (swapToken.balanceOf(address(this)) < tokenAmount) {
                revert INSUFFICIENT_TOKEN_BALANCE();
            }
            swapToken.safeTransfer(recipient, tokenAmount);
            emit SwapETHForTokenCompleted(recipient, tokenAmount);
        } else if (action == Action.SWAP_TOKEN_TO_ETH) {
            // L1→L2→L1 completion: forward ETH to recipient
            (, address recipient,) = abi.decode(_data, (Action, address, uint256));
            if (msg.value > 0) {
                (bool success,) = recipient.call{ value: msg.value }("");
                if (!success) revert ETH_TRANSFER_FAILED();
            }
            emit SwapTokenForETHCompleted(recipient, msg.value);
        } else if (action == Action.REMOVE_LIQUIDITY) {
            (, address recipient, uint256 tokenAmount) =
                abi.decode(_data, (Action, address, uint256));
            if (tokenAmount > 0) {
                if (swapToken.balanceOf(address(this)) < tokenAmount) {
                    revert INSUFFICIENT_TOKEN_BALANCE();
                }
                swapToken.safeTransfer(recipient, tokenAmount);
            }
            if (msg.value > 0) {
                (bool success,) = recipient.call{ value: msg.value }("");
                if (!success) revert ETH_TRANSFER_FAILED();
            }
            emit LiquidityRemovalCompleted(recipient, msg.value, tokenAmount);
        } else if (action == Action.SWAP_ETH_TO_TOKEN_VIA_L1) {
            _handleL1DexSwapETHForToken(_data);
        } else if (action == Action.SWAP_TOKEN_TO_ETH_VIA_L1) {
            _handleL1DexSwapTokenForETH(_data);
        } else {
            revert UNKNOWN_ACTION();
        }
        // BRIDGE and ADD_LIQUIDITY don't have L2→L1 completions
    }

    // ---------------------------------------------------------------
    // L2→L1→L2 leg handlers (L1 DEX side)
    // ---------------------------------------------------------------

    /// @dev Receives bridged ETH from L2, swaps it on the L1 router for tokens, then
    /// sends a return message instructing L2Vault to mint mirror tokens to recipient.
    function _handleL1DexSwapETHForToken(bytes calldata _data) internal {
        (,, address recipient,, uint256 minTokenOut) =
            abi.decode(_data, (Action, address, address, uint256, uint256));

        if (msg.value == 0) revert ZERO_AMOUNT();

        address[] memory path = new address[](2);
        path[0] = weth;
        path[1] = address(swapToken);

        uint256[] memory amounts = l1Router.swapExactETHForTokens{
            value: msg.value
        }(minTokenOut, path, address(this), block.timestamp + 60);
        uint256 tokenOut = amounts[1];

        // Send return message to L2: "mint tokenOut bUSDC to recipient"
        bytes memory returnData = abi.encode(Action.SWAP_ETH_TO_TOKEN_VIA_L1, recipient, tokenOut);
        bytes32 returnHash = _sendMessageToL2(returnData, 0, address(this), RETURN_GAS_LIMIT);

        emit L1DexSwapETHForTokenExecuted(recipient, msg.value, tokenOut, returnHash);
    }

    /// @dev Pulls swapToken from this vault's inventory, swaps it on the L1 router for ETH,
    /// then bridges the ETH back to L2Vault for delivery to recipient.
    function _handleL1DexSwapTokenForETH(bytes calldata _data) internal {
        (,, address recipient, uint256 amountIn, uint256 minETHOut) =
            abi.decode(_data, (Action, address, address, uint256, uint256));

        if (amountIn == 0) revert ZERO_AMOUNT();
        if (swapToken.balanceOf(address(this)) < amountIn) revert INSUFFICIENT_TOKEN_BALANCE();

        // Approve the router. Set to 0 first to be safe with non-standard ERC20s.
        swapToken.forceApprove(address(l1Router), amountIn);

        address[] memory path = new address[](2);
        path[0] = address(swapToken);
        path[1] = weth;

        uint256[] memory amounts = l1Router.swapExactTokensForETH(
            amountIn, minETHOut, path, address(this), block.timestamp + 60
        );
        uint256 ethOut = amounts[1];

        if (address(this).balance < ethOut) revert INSUFFICIENT_ETH_BALANCE();

        bytes memory returnData = abi.encode(Action.SWAP_TOKEN_TO_ETH_VIA_L1, recipient, ethOut);
        bytes32 returnHash = _sendMessageToL2(returnData, ethOut, address(this), RETURN_GAS_LIMIT);

        emit L1DexSwapTokenForETHExecuted(recipient, amountIn, ethOut, returnHash);
    }

    // ---------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------

    function _sendMessageToL2(
        bytes memory _innerData,
        uint256 _ethValue,
        address _srcOwner,
        uint32 _gasLimit
    )
        internal
        returns (bytes32)
    {
        bytes memory msgData = abi.encodeWithSignature("onMessageInvocation(bytes)", _innerData);

        IBridge.Message memory message = IBridge.Message({
            id: 0,
            fee: 0,
            gasLimit: _gasLimit,
            from: address(0),
            srcChainId: 0,
            srcOwner: _srcOwner,
            destChainId: l2ChainId,
            destOwner: l2Vault,
            to: l2Vault,
            value: _ethValue,
            data: msgData
        });

        (bytes32 msgHash,) = IBridge(bridge).sendMessage{ value: _ethValue }(message);
        return msgHash;
    }

    receive() external payable { }
}
