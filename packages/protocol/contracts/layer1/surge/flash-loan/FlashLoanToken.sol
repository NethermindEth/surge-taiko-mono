// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import { ERC20 } from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import { ERC20Burnable } from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";

/// @title FlashLoanToken
/// @notice Freely-mintable ERC20 used by the L2→L1→L2 flash loan demo. Deployed
///         identically on L1 and L2. Anyone can mint; holders can burn their own
///         balance.
/// @dev Decimals are configurable via the constructor (6 for the demo).
/// @custom:security-contact security@nethermind.io
contract FlashLoanToken is ERC20, ERC20Burnable {
    uint8 private immutable _tokenDecimals;

    constructor(
        string memory _name,
        string memory _symbol,
        uint8 _decimals
    )
        ERC20(_name, _symbol)
    {
        _tokenDecimals = _decimals;
    }

    function decimals() public view override returns (uint8) {
        return _tokenDecimals;
    }

    /// @notice Open faucet — anyone can mint to any address (demo only).
    function mint(address _to, uint256 _amount) external {
        _mint(_to, _amount);
    }
}
