// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/// @title SimpleDEXL1
/// @notice Test-mode L1 DEX for the L2→L1→L2 swap pipeline. Implements the subset
/// of the Uniswap V2 router02 ABI that `CrossChainSwapVaultL1` calls
/// (`swapExactETHForTokens`, `swapExactTokensForETH`, `WETH`, `getAmountsOut`),
/// so the vault is unaware whether it is wired to a live Uniswap deployment or
/// to this contract.
/// @dev Constant product (x·y=k) with a 0.3% fee. Holds raw ETH (no WETH wrapping
/// internally) — the WETH9 reference is only used for path[] validation so the
/// Uniswap V2 calling convention is preserved.
/// @custom:security-contact security@taiko.xyz
contract SimpleDEXL1 {
    using SafeERC20 for IERC20;

    uint256 public constant FEE_NUMERATOR = 3;
    uint256 public constant FEE_DENOMINATOR = 1000;

    /// @notice The ERC20 token paired with ETH
    IERC20 public immutable token;

    /// @notice WETH9 address used as a marker in path[] (Uniswap V2 router compatibility)
    address public immutable wethAddress;

    /// @notice Admin who can seed the pool
    address public immutable admin;

    /// @notice ETH reserve held by this contract
    uint256 public reserveETH;
    /// @notice Token reserve held by this contract
    uint256 public reserveToken;

    event LiquiditySeeded(uint256 ethAmount, uint256 tokenAmount);
    event SwapETHForToken(address indexed to, uint256 ethIn, uint256 tokenOut);
    event SwapTokenForETH(address indexed to, uint256 tokenIn, uint256 ethOut);

    error ONLY_ADMIN();
    error ZERO_AMOUNT();
    error INVALID_PATH();
    error EXPIRED();
    error INSUFFICIENT_OUTPUT();
    error INSUFFICIENT_LIQUIDITY();
    error ETH_TRANSFER_FAILED();

    constructor(address _token, address _weth, address _admin) {
        token = IERC20(_token);
        wethAddress = _weth;
        admin = _admin;
    }

    // ---------------------------------------------------------------
    // Uniswap V2 router02 surface
    // ---------------------------------------------------------------

    function WETH() external view returns (address) {
        return wethAddress;
    }

    function swapExactETHForTokens(
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256 deadline
    )
        external
        payable
        returns (uint256[] memory amounts)
    {
        if (block.timestamp > deadline) revert EXPIRED();
        if (msg.value == 0) revert ZERO_AMOUNT();
        if (path.length != 2 || path[0] != wethAddress || path[1] != address(token)) {
            revert INVALID_PATH();
        }

        uint256 tokenOut = _getAmountOut(msg.value, reserveETH, reserveToken);
        if (tokenOut < amountOutMin) revert INSUFFICIENT_OUTPUT();
        if (tokenOut > reserveToken) revert INSUFFICIENT_LIQUIDITY();

        reserveETH += msg.value;
        reserveToken -= tokenOut;

        token.safeTransfer(to, tokenOut);

        amounts = new uint256[](2);
        amounts[0] = msg.value;
        amounts[1] = tokenOut;

        emit SwapETHForToken(to, msg.value, tokenOut);
    }

    function swapExactTokensForETH(
        uint256 amountIn,
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256 deadline
    )
        external
        returns (uint256[] memory amounts)
    {
        if (block.timestamp > deadline) revert EXPIRED();
        if (amountIn == 0) revert ZERO_AMOUNT();
        if (path.length != 2 || path[0] != address(token) || path[1] != wethAddress) {
            revert INVALID_PATH();
        }

        uint256 ethOut = _getAmountOut(amountIn, reserveToken, reserveETH);
        if (ethOut < amountOutMin) revert INSUFFICIENT_OUTPUT();
        if (ethOut > reserveETH) revert INSUFFICIENT_LIQUIDITY();

        token.safeTransferFrom(msg.sender, address(this), amountIn);

        reserveToken += amountIn;
        reserveETH -= ethOut;

        (bool ok,) = to.call{ value: ethOut }("");
        if (!ok) revert ETH_TRANSFER_FAILED();

        amounts = new uint256[](2);
        amounts[0] = amountIn;
        amounts[1] = ethOut;

        emit SwapTokenForETH(to, amountIn, ethOut);
    }

    function getAmountsOut(
        uint256 amountIn,
        address[] calldata path
    )
        external
        view
        returns (uint256[] memory amounts)
    {
        if (path.length != 2) revert INVALID_PATH();

        amounts = new uint256[](2);
        amounts[0] = amountIn;

        if (path[0] == wethAddress && path[1] == address(token)) {
            amounts[1] = _getAmountOut(amountIn, reserveETH, reserveToken);
        } else if (path[0] == address(token) && path[1] == wethAddress) {
            amounts[1] = _getAmountOut(amountIn, reserveToken, reserveETH);
        } else {
            revert INVALID_PATH();
        }
    }

    // ---------------------------------------------------------------
    // Demo helpers
    // ---------------------------------------------------------------

    /// @notice One-shot liquidity seeding — admin pulls tokens from itself, sends ETH along.
    function seedLiquidity(uint256 _tokenAmount) external payable {
        if (msg.sender != admin) revert ONLY_ADMIN();
        if (msg.value == 0 || _tokenAmount == 0) revert ZERO_AMOUNT();

        token.safeTransferFrom(msg.sender, address(this), _tokenAmount);

        reserveETH += msg.value;
        reserveToken += _tokenAmount;

        emit LiquiditySeeded(msg.value, _tokenAmount);
    }

    function getReserves() external view returns (uint256 ethReserve_, uint256 tokenReserve_) {
        return (reserveETH, reserveToken);
    }

    // ---------------------------------------------------------------
    // Internal
    // ---------------------------------------------------------------

    function _getAmountOut(
        uint256 _amountIn,
        uint256 _reserveIn,
        uint256 _reserveOut
    )
        internal
        pure
        returns (uint256)
    {
        if (_reserveIn == 0 || _reserveOut == 0) return 0;
        uint256 amountInWithFee = _amountIn * (FEE_DENOMINATOR - FEE_NUMERATOR);
        uint256 numerator = amountInWithFee * _reserveOut;
        uint256 denominator = (_reserveIn * FEE_DENOMINATOR) + amountInWithFee;
        return numerator / denominator;
    }

    receive() external payable { }
}
