// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { IBridge } from "../../../shared/bridge/IBridge.sol";
import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

interface ISwapTokenL2 {
    function mint(address _to, uint256 _amount) external;
    function burn(address _from, uint256 _amount) external;
}

interface ISimpleDEX {
    function swapETHForToken(uint256 _minTokenOut) external payable returns (uint256);
    function swapTokenForETH(uint256 _tokenIn, uint256 _minETHOut) external returns (uint256);
    function addLiquidity(uint256 _tokenAmount, address _provider) external payable;
    function removeLiquidity(address _provider) external returns (uint256, uint256);
    function token() external view returns (IERC20);
}

/// @title CrossChainSwapVaultL2
/// @notice L2 counterpart of CrossChainSwapVaultL1. Receives bridge messages and
/// handles minting bridged tokens, DEX swaps, and liquidity provisioning.
/// Also exposes user-facing entrypoints that initiate L2→L1→L2 swaps against
/// the L1-side DEX (paired contract: `CrossChainSwapVaultL1`'s `l1Router`).
/// @dev Has minting authority over the bridged ERC20 (SwapTokenL2).
/// @custom:security-contact security@taiko.xyz
contract CrossChainSwapVaultL2 {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------
    // Enums (must match L1 vault)
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
    uint64 public immutable l1ChainId;
    ISimpleDEX public immutable dex;
    ISwapTokenL2 public immutable swapToken;
    IERC20 public immutable swapTokenERC20;
    address public immutable admin;
    address public l1Vault;

    uint32 public constant GAS_LIMIT = 1_000_000;

    // ---------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------

    event L1VaultSet(address indexed l1Vault);
    event TokensBridged(address indexed recipient, uint256 amount);
    event SwapExecutedETHToToken(address indexed recipient, uint256 ethIn, uint256 tokenOut);
    event SwapExecutedTokenToETH(address indexed recipient, uint256 tokenIn, uint256 ethOut);
    event LiquidityAdded(uint256 ethAmount, uint256 tokenAmount);
    // L2→L1→L2 events
    event L1DexSwapInitiatedETHForToken(
        address indexed user,
        address indexed recipient,
        uint256 ethIn,
        uint256 minTokenOut,
        bytes32 outboundMsgHash
    );
    event L1DexSwapInitiatedTokenForETH(
        address indexed user,
        address indexed recipient,
        uint256 tokenIn,
        uint256 minETHOut,
        bytes32 outboundMsgHash
    );
    event L1DexSwapCompletedETHForToken(address indexed recipient, uint256 tokenOut);
    event L1DexSwapCompletedTokenForETH(address indexed recipient, uint256 ethOut);

    // ---------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------

    error ONLY_ADMIN();
    error ONLY_BRIDGE();
    error INVALID_SENDER();
    error L1_VAULT_NOT_SET();
    error UNKNOWN_ACTION();
    error ZERO_AMOUNT();
    error INVALID_RETURN_MESSAGE();
    error ETH_TRANSFER_FAILED();

    // ---------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------

    constructor(
        address _bridge,
        uint64 _l1ChainId,
        address _dex,
        address _swapToken,
        address _admin
    ) {
        bridge = _bridge;
        l1ChainId = _l1ChainId;
        dex = ISimpleDEX(_dex);
        swapToken = ISwapTokenL2(_swapToken);
        swapTokenERC20 = IERC20(_swapToken);
        admin = _admin;
    }

    // ---------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------

    function setL1Vault(address _l1Vault) external {
        if (msg.sender != admin) revert ONLY_ADMIN();
        l1Vault = _l1Vault;
        emit L1VaultSet(_l1Vault);
    }

    // ---------------------------------------------------------------
    // L2→L1→L2 Entry Points (called directly by user EOA, NOT via UserOp)
    //
    // Both entrypoints expect a Catalyst-simulated `returnMessage` that the
    // bridge will deliver back into this contract within the same L2 tx.
    // ---------------------------------------------------------------

    /// @notice Swap L2 ETH for bUSDC by routing through the L1 DEX.
    /// @param _minTokenOut Minimum bUSDC the user is willing to accept.
    /// @param _recipient Recipient of the resulting bUSDC on L2.
    /// @param _returnMessage Catalyst-simulated L1→L2 return message that mints bUSDC here.
    function swapETHForTokenViaL1(
        uint256 _minTokenOut,
        address _recipient,
        IBridge.Message calldata _returnMessage
    )
        external
        payable
    {
        if (l1Vault == address(0)) revert L1_VAULT_NOT_SET();
        if (msg.value == 0) revert ZERO_AMOUNT();
        // Sanity: the simulated return must target *this* contract or the unlock won't fire.
        if (_returnMessage.to != address(this)) revert INVALID_RETURN_MESSAGE();

        bytes memory data = abi.encode(
            Action.SWAP_ETH_TO_TOKEN_VIA_L1, msg.sender, _recipient, msg.value, _minTokenOut
        );
        bytes32 outboundMsgHash = _sendMessageToL1(data, msg.value);

        emit L1DexSwapInitiatedETHForToken(
            msg.sender, _recipient, msg.value, _minTokenOut, outboundMsgHash
        );

        // Consume the Catalyst-injected fast signal — drives the L1→L2 return into us.
        IBridge(bridge).processMessage(_returnMessage, "");
    }

    /// @notice Swap L2 bUSDC for L2 ETH by routing through the L1 DEX.
    /// @param _amountIn Amount of bUSDC to swap.
    /// @param _minETHOut Minimum L2 ETH the user is willing to accept.
    /// @param _recipient Recipient of the resulting L2 ETH.
    /// @param _returnMessage Catalyst-simulated L1→L2 return message that delivers ETH here.
    function swapTokenForETHViaL1(
        uint256 _amountIn,
        uint256 _minETHOut,
        address _recipient,
        IBridge.Message calldata _returnMessage
    )
        external
    {
        if (l1Vault == address(0)) revert L1_VAULT_NOT_SET();
        if (_amountIn == 0) revert ZERO_AMOUNT();
        if (_returnMessage.to != address(this)) revert INVALID_RETURN_MESSAGE();

        // Pull and burn the user's bUSDC. The L1 vault will draw matching USDC from its
        // inventory and swap it for ETH.
        swapTokenERC20.safeTransferFrom(msg.sender, address(this), _amountIn);
        swapToken.burn(address(this), _amountIn);

        bytes memory data = abi.encode(
            Action.SWAP_TOKEN_TO_ETH_VIA_L1, msg.sender, _recipient, _amountIn, _minETHOut
        );
        bytes32 outboundMsgHash = _sendMessageToL1(data, 0);

        emit L1DexSwapInitiatedTokenForETH(
            msg.sender, _recipient, _amountIn, _minETHOut, outboundMsgHash
        );

        IBridge(bridge).processMessage(_returnMessage, "");
    }

    // ---------------------------------------------------------------
    // Bridge Callback (from L1)
    // ---------------------------------------------------------------

    /// @notice Called by bridge when L1 vault sends a message
    function onMessageInvocation(bytes calldata _data) external payable {
        if (msg.sender != bridge) revert ONLY_BRIDGE();

        IBridge.Context memory ctx = IBridge(bridge).context();
        if (l1Vault == address(0)) revert L1_VAULT_NOT_SET();
        if (ctx.from != l1Vault) revert INVALID_SENDER();

        Action action = abi.decode(_data, (Action));

        if (action == Action.BRIDGE) {
            _handleBridge(_data);
        } else if (action == Action.SWAP_ETH_TO_TOKEN) {
            _handleSwapETHToToken(_data);
        } else if (action == Action.SWAP_TOKEN_TO_ETH) {
            _handleSwapTokenToETH(_data);
        } else if (action == Action.ADD_LIQUIDITY) {
            _handleAddLiquidity(_data);
        } else if (action == Action.REMOVE_LIQUIDITY) {
            _handleRemoveLiquidity(_data);
        } else if (action == Action.SWAP_ETH_TO_TOKEN_VIA_L1) {
            _handleL1DexReturnETHForToken(_data);
        } else if (action == Action.SWAP_TOKEN_TO_ETH_VIA_L1) {
            _handleL1DexReturnTokenForETH(_data);
        } else {
            revert UNKNOWN_ACTION();
        }
    }

    // ---------------------------------------------------------------
    // Internal Handlers (L1→L2→L1 leg)
    // ---------------------------------------------------------------

    function _handleBridge(bytes calldata _data) internal {
        (, address recipient, uint256 amount) = abi.decode(_data, (Action, address, uint256));
        swapToken.mint(recipient, amount);
        emit TokensBridged(recipient, amount);
    }

    function _handleSwapETHToToken(bytes calldata _data) internal {
        (,, address recipient,, uint256 minTokenOut) =
            abi.decode(_data, (Action, address, address, uint256, uint256));

        uint256 tokenOut = dex.swapETHForToken{ value: msg.value }(minTokenOut);
        swapToken.burn(address(this), tokenOut);

        emit SwapExecutedETHToToken(recipient, msg.value, tokenOut);

        bytes memory completionData = abi.encode(Action.SWAP_ETH_TO_TOKEN, recipient, tokenOut);
        _sendMessageToL1(completionData, 0);
    }

    function _handleSwapTokenToETH(bytes calldata _data) internal {
        (,, address recipient, uint256 tokenAmount, uint256 minETHOut) =
            abi.decode(_data, (Action, address, address, uint256, uint256));

        swapToken.mint(address(this), tokenAmount);
        swapTokenERC20.approve(address(dex), tokenAmount);
        uint256 ethOut = dex.swapTokenForETH(tokenAmount, minETHOut);

        emit SwapExecutedTokenToETH(recipient, tokenAmount, ethOut);

        bytes memory completionData = abi.encode(Action.SWAP_TOKEN_TO_ETH, recipient, ethOut);
        _sendMessageToL1(completionData, ethOut);
    }

    function _handleAddLiquidity(bytes calldata _data) internal {
        (, address provider, uint256 tokenAmount) = abi.decode(_data, (Action, address, uint256));

        swapToken.mint(address(this), tokenAmount);
        swapTokenERC20.approve(address(dex), tokenAmount);
        dex.addLiquidity{ value: msg.value }(tokenAmount, provider);

        emit LiquidityAdded(msg.value, tokenAmount);
    }

    function _handleRemoveLiquidity(bytes calldata _data) internal {
        (, address provider) = abi.decode(_data, (Action, address));

        (uint256 ethAmount, uint256 tokenAmount) = dex.removeLiquidity(provider);

        if (tokenAmount > 0) {
            swapToken.burn(address(this), tokenAmount);
        }

        bytes memory completionData = abi.encode(Action.REMOVE_LIQUIDITY, provider, tokenAmount);
        _sendMessageToL1(completionData, ethAmount);
    }

    // ---------------------------------------------------------------
    // Internal Handlers (L2→L1→L2 return leg)
    // ---------------------------------------------------------------

    function _handleL1DexReturnETHForToken(bytes calldata _data) internal {
        (, address recipient, uint256 tokenOut) = abi.decode(_data, (Action, address, uint256));

        // L1 took the user's ETH, swapped on L1 DEX, told us how much bUSDC to mint.
        if (tokenOut > 0) {
            swapToken.mint(recipient, tokenOut);
        }
        emit L1DexSwapCompletedETHForToken(recipient, tokenOut);
    }

    function _handleL1DexReturnTokenForETH(bytes calldata _data) internal {
        (, address recipient, uint256 ethOut) = abi.decode(_data, (Action, address, uint256));

        // L1 used its USDC inventory to swap and bridged the resulting ETH back to us.
        if (msg.value > 0) {
            (bool ok,) = recipient.call{ value: msg.value }("");
            if (!ok) revert ETH_TRANSFER_FAILED();
        }
        emit L1DexSwapCompletedTokenForETH(recipient, ethOut);
    }

    // ---------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------

    function _sendMessageToL1(
        bytes memory _innerData,
        uint256 _ethValue
    )
        internal
        returns (bytes32 msgHash_)
    {
        bytes memory msgData = abi.encodeWithSignature("onMessageInvocation(bytes)", _innerData);

        IBridge.Message memory message = IBridge.Message({
            id: 0,
            fee: 0,
            gasLimit: GAS_LIMIT,
            from: address(0),
            srcChainId: 0,
            srcOwner: address(this),
            destChainId: l1ChainId,
            destOwner: l1Vault,
            to: l1Vault,
            value: _ethValue,
            data: msgData
        });

        (msgHash_,) = IBridge(bridge).sendMessage{ value: _ethValue }(message);
    }

    receive() external payable { }
}
