#!/usr/bin/env python3
"""
Minimal test for EIP-7702: Set EOA account code
Tests whether an EOA can have code set via authorization lists
"""

import rlp
from web3 import Web3
from eth_account import Account
from eth_keys import keys

# Configuration
RPC_URL = "http://localhost:8547"
PRIVATE_KEY = "0x94eb3102993b41ec55c241060f47daa0f6372e2e3ad7e91612ae36c364042e44"
CHAIN_ID = 763374  # Taiko chain ID

# Connect to node
w3 = Web3(Web3.HTTPProvider(RPC_URL))
account = Account.from_key(PRIVATE_KEY)

print("=" * 60)
print("EIP-7702: Set EOA Account Code Test")
print("=" * 60)
print(f"Connected to: {RPC_URL}")
print(f"Chain ID: {w3.eth.chain_id}")
print(f"Test account: {account.address}")
print(f"Account balance: {w3.from_wei(w3.eth.get_balance(account.address), 'ether')} ETH")
print()

print("Step 1: Deploy delegate contract")
# Simple contract that returns a constant value and has a state variable
# contract Test { 
#   uint256 public value;
#   function getValue() public view returns (uint256) { return value; }
#   function setValue(uint256 _value) public { value = _value; }
#   function getConstant() public pure returns (uint256) { return 12345; }
# }
delegate_bytecode = "0x6080604052348015600e575f5ffd5b506101a38061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061004a575f3560e01c8063209652551461004e5780633fa4f2451461006c578063552410771461008a578063f13a38a6146100a6575b5f5ffd5b6100566100c4565b60405161006391906100fb565b60405180910390f35b6100746100cc565b60405161008191906100fb565b60405180910390f35b6100a4600480360381019061009f9190610142565b6100d1565b005b6100ae6100da565b6040516100bb91906100fb565b60405180910390f35b5f5f54905090565b5f5481565b805f8190555050565b5f613039905090565b5f819050919050565b6100f5816100e3565b82525050565b5f60208201905061010e5f8301846100ec565b92915050565b5f5ffd5b610121816100e3565b811461012b575f5ffd5b50565b5f8135905061013c81610118565b92915050565b5f6020828403121561015757610156610114565b5b5f6101648482850161012e565b9150509291505056fea264697066735822122011dc001a28350260e7b3193ee7e82fa0d38cd9ef8557ae1c7149ad397e78402864736f6c634300081b0033"

tx_deploy = {
    'from': account.address,
    'data': delegate_bytecode,
    'gas': 500000,
    'gasPrice': w3.eth.gas_price,
    'nonce': w3.eth.get_transaction_count(account.address),
    'chainId': w3.eth.chain_id
}

signed_tx = w3.eth.account.sign_transaction(tx_deploy, PRIVATE_KEY)
tx_hash = w3.eth.send_raw_transaction(signed_tx.raw_transaction)
tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
delegate_address = tx_receipt['contractAddress']

print(f"✓ Delegate contract deployed at: {delegate_address}")
print(f"  Gas used: {tx_receipt['gasUsed']}")
print()

# Create a new EOA that will be delegated
eoa_account = Account.create()
print(f"Step 2: Created new EOA: {eoa_account.address}")

# Check initial code (should be empty)
initial_code = w3.eth.get_code(eoa_account.address)
print(f"✓ Initial code length: {len(initial_code)} bytes")
assert len(initial_code) == 0, "EOA should have no code initially"
print()

# Fund the new EOA
print("Step 3: Fund the new EOA")
tx_fund = {
    'from': account.address,
    'to': eoa_account.address,
    'value': w3.to_wei(1, 'ether'),
    'gas': 21000,
    'gasPrice': w3.eth.gas_price,
    'nonce': w3.eth.get_transaction_count(account.address),
    'chainId': w3.eth.chain_id
}
signed_tx = w3.eth.account.sign_transaction(tx_fund, PRIVATE_KEY)
tx_hash = w3.eth.send_raw_transaction(signed_tx.raw_transaction)
w3.eth.wait_for_transaction_receipt(tx_hash)
print(f"✓ Funded EOA with 1 ETH")
print(f"  EOA balance: {w3.from_wei(w3.eth.get_balance(eoa_account.address), 'ether')} ETH")
print()

# Create EIP-7702 authorization
print("Step 4: Create EIP-7702 authorization")

# EIP-7702 authorization tuple format:
# [chain_id, address, nonce]
# MAGIC = 0x05 as per EIP-7702
MAGIC = b"\x05"

# Get the EOA's current nonce
eoa_nonce = w3.eth.get_transaction_count(eoa_account.address)+1

# Encode the authorization: [chain_id, address, nonce]
encoded = rlp.encode([w3.eth.chain_id, bytes.fromhex(delegate_address[2:]), eoa_nonce])

# Create the message hash with MAGIC prefix
msg_hash = w3.keccak(MAGIC + encoded)

# Sign with the EOA's private key (the account that will delegate)
private_key = keys.PrivateKey(eoa_account.key)
signature = private_key.sign_msg_hash(msg_hash)

# Extract r, s, v from signature
# eth_keys returns v as 0 or 1 (y_parity), not 27/28
r = signature.r
s = signature.s
y_parity = signature.v  # Already 0 or 1

if y_parity not in (0, 1):
    raise ValueError(f"Unexpected y_parity={y_parity}")

# Create the authorization tuple for the transaction
authorization_tuple = {
    'chainId': w3.eth.chain_id,
    'address': delegate_address,
    'nonce': eoa_nonce,
    'yParity': y_parity,
    'r': r,
    's': s
}

print(f"✓ Authorization signed by EOA")
print(f"  EOA address: {eoa_account.address}")
print(f"  Delegate to: {delegate_address}")
print(f"  Chain ID: {w3.eth.chain_id}")
print(f"  Nonce: {eoa_nonce}")
print(f"  y_parity: {y_parity}")
print(f"  r: {hex(r)[:18]}...")
print(f"  s: {hex(s)[:18]}...")
print()

# Send transaction with authorization list (EIP-7702)
print("Step 5: Send transaction with EIP-7702 authorization")
print("Note: If this fails, EIP-7702 may not be activated on this network")

try:
    # Try to send a transaction with authorization list
    # This is the EIP-7702 specific format (Type 4 transaction)
    # EIP-7702 transactions must use EIP-1559 format (maxFeePerGas, maxPriorityFeePerGas)
    base_fee = w3.eth.get_block('latest')['baseFeePerGas']
    max_priority_fee = w3.to_wei(2, 'gwei')
    max_fee = base_fee * 2 + max_priority_fee
    
    tx_with_auth = {
        'type': 4,  # EIP-7702 transaction type
        'from': eoa_account.address,
        'to': eoa_account.address,  # Self-transaction to trigger delegation
        'value': 0,
        'gas': 500000,
        'maxFeePerGas': max_fee,
        'maxPriorityFeePerGas': max_priority_fee,
        'nonce': w3.eth.get_transaction_count(eoa_account.address),  # Use the same nonce we signed in the authorization
        'chainId': w3.eth.chain_id,
        'authorizationList': [
            authorization_tuple
        ]
    }

    print("Transaction with authorization:")
    print(f"  Type: {tx_with_auth['type']}")
    print(f"  From: {tx_with_auth['from']}")
    print(f"  Nonce: {tx_with_auth['nonce']}")
    print(f"  Max Fee: {w3.from_wei(max_fee, 'gwei')} gwei")
    print()
    
    signed_tx = w3.eth.account.sign_transaction(tx_with_auth, eoa_account.key)
    tx_hash = w3.eth.send_raw_transaction(signed_tx.raw_transaction)
    tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
    
    print(f"✓ Transaction sent: {tx_hash.hex()}")
    print(f"  Gas used: {tx_receipt['gasUsed']}")
    print(f"  Status: {tx_receipt['status']}")
    
    # Check the transaction details
    tx = w3.eth.get_transaction(tx_hash)
    print(f"  Transaction type: {tx.get('type', 'unknown')}")
    print()
    
    # Check if delegation is working
    print("Step 6: Verify delegation is working")
    print("Note: EOA delegates calls to the contract, storage is in EOA's space")
    
    # The delegate contract has functions:
    # - getConstant() -> returns 12345 (pure function, no state)
    # - getValue() -> returns stored value at slot 0
    # - setValue(uint256) -> stores a value at slot 0
    
    # First, try calling a pure function to prove delegation works
    # Function signature: getConstant() -> 0x20965255
    get_constant_sig = w3.keccak(text='getConstant()')[:4]
    
    try:
        # Call getConstant() on the EOA address
        result = w3.eth.call({
            'to': eoa_account.address,
            'data': get_constant_sig.hex()
        })
        
        constant_value = int.from_bytes(result, 'big')
        print(f"✓ Called getConstant() on EOA: {constant_value}")
        
        if constant_value == 12345:
            print(f"  ✓ Returned correct constant (12345)")
            print(f"  ✓ This proves EIP-7702 delegation is working!")
            print("\n" + "=" * 60)
            print("✓✓✓ SUCCESS: EIP-7702 IS WORKING! ✓✓✓")
            print("=" * 60)
            print(f"EOA {eoa_account.address} successfully delegates to {delegate_address}")
            print(f"Pure function call returned correct value: {constant_value}")
            exit(0)
        else:
            print(f"  ✗ Expected 12345 but got {constant_value}")
            print("\n" + "=" * 60)
            print("⚠ PARTIAL SUCCESS")
            print("=" * 60)
            print("Delegation works but returned unexpected value")
            exit(1)
            
    except Exception as e:
        print(f"✗ Delegation call failed: {e}")
        print("\n" + "=" * 60)
        print("✗✗✗ FAILURE: EIP-7702 NOT WORKING ✗✗✗")
        print("=" * 60)
        print("EOA does not delegate calls to the contract")
        exit(1)
        
except Exception as e:
    print(f"\n✗ Transaction failed: {e}")
    print("\nThis likely means:")
    print("1. EIP-7702 is not activated on this network")
    print("2. The authorization list format is not supported")
    print("3. The node does not implement EIP-7702 yet")
    print("\n" + "=" * 60)
    print("✗✗✗ EIP-7702 NOT AVAILABLE ✗✗✗")
    print("=" * 60)
    exit(1)
