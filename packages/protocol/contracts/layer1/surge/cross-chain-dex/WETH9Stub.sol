// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

/// @title WETH9Stub
/// @notice Minimal WETH9-style wrapper used in test mode for `SimpleDEXL1`.
/// @dev Mirrors the canonical WETH9 surface (`deposit`, `withdraw`, ERC20)
/// so that contracts (and front-ends) speaking the Uniswap V2 router ABI can
/// reference its address as `path[0]` / `path[1]`.
/// @custom:security-contact security@taiko.xyz
contract WETH9Stub {
    string public constant name = "Wrapped Ether (Stub)";
    string public constant symbol = "WETH";
    uint8 public constant decimals = 18;

    event Deposit(address indexed dst, uint256 wad);
    event Withdrawal(address indexed src, uint256 wad);
    event Transfer(address indexed src, address indexed dst, uint256 wad);
    event Approval(address indexed src, address indexed guy, uint256 wad);

    error INSUFFICIENT_BALANCE();
    error INSUFFICIENT_ALLOWANCE();
    error ETH_TRANSFER_FAILED();

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    receive() external payable {
        deposit();
    }

    function deposit() public payable {
        balanceOf[msg.sender] += msg.value;
        emit Deposit(msg.sender, msg.value);
    }

    function withdraw(uint256 wad) external {
        if (balanceOf[msg.sender] < wad) revert INSUFFICIENT_BALANCE();
        balanceOf[msg.sender] -= wad;
        (bool ok,) = msg.sender.call{ value: wad }("");
        if (!ok) revert ETH_TRANSFER_FAILED();
        emit Withdrawal(msg.sender, wad);
    }

    function totalSupply() external view returns (uint256) {
        return address(this).balance;
    }

    function approve(address guy, uint256 wad) external returns (bool) {
        allowance[msg.sender][guy] = wad;
        emit Approval(msg.sender, guy, wad);
        return true;
    }

    function transfer(address dst, uint256 wad) external returns (bool) {
        return _transferFrom(msg.sender, dst, wad);
    }

    function transferFrom(address src, address dst, uint256 wad) external returns (bool) {
        if (src != msg.sender && allowance[src][msg.sender] != type(uint256).max) {
            if (allowance[src][msg.sender] < wad) revert INSUFFICIENT_ALLOWANCE();
            allowance[src][msg.sender] -= wad;
        }
        return _transferFrom(src, dst, wad);
    }

    function _transferFrom(address src, address dst, uint256 wad) internal returns (bool) {
        if (balanceOf[src] < wad) revert INSUFFICIENT_BALANCE();
        balanceOf[src] -= wad;
        balanceOf[dst] += wad;
        emit Transfer(src, dst, wad);
        return true;
    }
}
