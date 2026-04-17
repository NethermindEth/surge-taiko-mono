// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

/// @title IUniswapV2Router02 (minimal)
/// @notice Subset of the Uniswap V2 Router02 interface used by `CrossChainSwapVaultL1`
/// when performing the L1 leg of an L2→L1→L2 swap. Both the live Uniswap V2 router and
/// the demo `SimpleDEXL1` implement this interface.
/// @custom:security-contact security@taiko.xyz
interface IUniswapV2Router02 {
    function WETH() external view returns (address);

    function swapExactETHForTokens(
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256 deadline
    )
        external
        payable
        returns (uint256[] memory amounts);

    function swapExactTokensForETH(
        uint256 amountIn,
        uint256 amountOutMin,
        address[] calldata path,
        address to,
        uint256 deadline
    )
        external
        returns (uint256[] memory amounts);

    function getAmountsOut(
        uint256 amountIn,
        address[] calldata path
    )
        external
        view
        returns (uint256[] memory amounts);
}
