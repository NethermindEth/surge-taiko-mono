// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { ReentrancyGuard } from "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

interface IFlashLoanReceiver {
    function onFlashLoan(address token, uint256 amount) external;
}

/// @title FlashLoanPool
/// @notice Minimal L2 flash loan pool. Standard EIP-3156-style semantics:
///         transfer out, callback, require balance unchanged or greater.
/// @dev The borrower is responsible for returning `amount` tokens. The L1
///      callback adds a fee that the borrower keeps as profit; this pool does
///      not take a fee itself.
/// @custom:security-contact security@nethermind.io
contract FlashLoanPool is ReentrancyGuard {
    using SafeERC20 for IERC20;

    IERC20 public immutable token;

    event FlashLoanIssued(address indexed borrower, uint256 amount);

    error NOT_ENOUGH_LIQUIDITY();
    error NOT_REPAID();

    constructor(address _token) {
        token = IERC20(_token);
    }

    /// @notice Borrow `_amount` tokens. Must be repaid within the same tx.
    function flashLoan(uint256 _amount) external nonReentrant {
        uint256 balanceBefore = token.balanceOf(address(this));
        if (balanceBefore < _amount) revert NOT_ENOUGH_LIQUIDITY();

        token.safeTransfer(msg.sender, _amount);
        emit FlashLoanIssued(msg.sender, _amount);

        IFlashLoanReceiver(msg.sender).onFlashLoan(address(token), _amount);

        uint256 balanceAfter = token.balanceOf(address(this));
        if (balanceAfter < balanceBefore) revert NOT_REPAID();
    }
}
