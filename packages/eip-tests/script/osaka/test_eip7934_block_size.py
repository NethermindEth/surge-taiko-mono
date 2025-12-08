#!/usr/bin/env python3
"""
Test EIP-7934: RLP Execution Block Size Limit (10 MiB)

This script attempts to create a block that exceeds the 10 MiB RLP-encoded size limit.
It does this by sending multiple transactions with large calldata payloads.
"""

from web3 import Web3
from eth_account import Account
import time
import sys

# Configuration
RPC_URL = "http://localhost:8547"  # Change to your node's RPC URL
# Use funded private key
PRIVATE_KEY = "0x94eb3102993b41ec55c241060f47daa0f6372e2e3ad7e91612ae36c364042e44"

# EIP-7934 constants
MAX_BLOCK_SIZE_BYTES = 10 * 1024 * 1024  # 10 MiB = 10,485,760 bytes
TARGET_SIZE_BYTES = MAX_BLOCK_SIZE_BYTES + (1 * 1024 * 1024)  # Try to exceed by 1 MiB

def estimate_tx_overhead():
    """
    Estimate the RLP overhead per transaction (without calldata).
    A typical transaction has:
    - nonce, gasPrice, gasLimit, to, value, data, v, r, s
    - Plus RLP encoding overhead
    - Approximately ~200-300 bytes per tx
    """
    return 300

def create_large_calldata_tx(w3, account, nonce, calldata_size):
    """
    Create a transaction with large calldata.
    
    Args:
        w3: Web3 instance
        account: Account object
        nonce: Transaction nonce
        calldata_size: Size of calldata in bytes
    
    Returns:
        Signed transaction
    """
    # Create large calldata (just zeros for simplicity)
    calldata = b'\x00' * calldata_size
    
    # Build transaction
    tx = {
        'chainId': w3.eth.chain_id,
        'from': account.address,
        'to': account.address,  # Send to self
        'value': 0,
        'gas': 21000 + (calldata_size * 16),  # Base gas + calldata cost
        'maxFeePerGas': w3.to_wei(2, 'gwei'),
        'maxPriorityFeePerGas': w3.to_wei(1, 'gwei'),
        'nonce': nonce,
        'data': calldata,
    }
    
    # Sign transaction
    signed_tx = account.sign_transaction(tx)
    return signed_tx

def calculate_optimal_tx_params():
    """
    Calculate optimal number of transactions and calldata size to exceed block limit.
    
    Strategy:
    1. Each transaction has ~300 bytes overhead
    2. We want to fit as many large txs as possible in one block
    3. Block gas limit is typically 30M, each large tx uses ~21k + calldata_gas
    4. Optimize for total RLP size > 10 MiB while staying under gas limit
    """
    tx_overhead = estimate_tx_overhead()
    
    # Try different strategies - adjusted to fit within 30M gas limit and tx size limits
    # Calldata cost: 16 gas per zero byte, so 100KB = 100,000 bytes = 1.6M gas + 21k = ~1.62M gas
    # Max tx size seems to be around 128KB, so use smaller transactions
    strategies = [
        # Strategy 1: Many 50KB transactions to test batching (should fit in one block)
        # {"calldata_size": 50_000, "num_txs": 15, "description": "15 txs × 50KB calldata (fits in one block)"},
        
        # Strategy 2: More 100KB transactions (might span multiple blocks)
        # {"calldata_size": 100_000, "num_txs": 10, "description": "10 txs × 100KB calldata"},
        
        # Strategy 3: Even more smaller transactions
        {"calldata_size": 100_000, "num_txs": 100, "description": "12 txs × 80KB calldata"},
    ]
    
    print("\nStrategy Analysis:")
    print("-" * 80)
    
    best_strategy = None
    best_size = 0
    
    for strategy in strategies:
        calldata_size = strategy["calldata_size"]
        num_txs = strategy["num_txs"]
        
        # Estimate total block size
        total_size = num_txs * (tx_overhead + calldata_size)
        
        # Estimate gas consumption
        gas_per_tx = 21000 + (calldata_size * 16)  # 16 gas per calldata byte (zeros)
        total_gas = num_txs * gas_per_tx
        
        exceeds = "YES ✓" if total_size > MAX_BLOCK_SIZE_BYTES else "NO ✗"
        fits_gas = "YES ✓" if total_gas < 30_000_000 else "NO ✗ (exceeds 30M gas limit)"
        
        print(f"\n{strategy['description']}:")
        print(f"  Estimated RLP size: {total_size:,} bytes ({total_size / (1024*1024):.2f} MiB)")
        print(f"  Exceeds 10 MiB limit: {exceeds}")
        print(f"  Estimated gas: {total_gas:,}")
        print(f"  Fits in gas limit: {fits_gas}")
        
        # Pick best strategy that fits gas limit (prefer one that doesn't exceed size for now)
        if total_gas < 30_000_000:
            if best_strategy is None or total_size > best_size:
                best_size = total_size
                best_strategy = strategy
    
    print("-" * 80)
    
    if best_strategy:
        print(f"\nBest strategy: {best_strategy['description']}")
        print(f"Estimated size: {best_size / (1024*1024):.2f} MiB")
        print(f"\n⚠️  IMPORTANT: All transactions will be sent RAPIDLY without delays")
        print(f"   This maximizes the chance they'll be included in the SAME block!")
        return best_strategy
    else:
        print("\nWARNING: No strategy fits within gas limits!")
        return strategies[0]  # Return first strategy anyway

def send_large_transactions(w3, account, strategy):
    """
    Send transactions according to the strategy.
    """
    calldata_size = strategy["calldata_size"]
    num_txs = strategy["num_txs"]
    
    print(f"\n{'=' * 80}")
    print(f"Sending {num_txs} transactions with {calldata_size:,} byte calldata each")
    print(f"{'=' * 80}\n")
    
    nonce = w3.eth.get_transaction_count(account.address)
    
    # IMPORTANT: Sign all transactions first, THEN send them rapidly
    # This maximizes the chance they'll be included in the same block
    print(f"Step 1: Signing {num_txs} transactions...")
    signed_txs = []
    
    for i in range(num_txs):
        try:
            signed_tx = create_large_calldata_tx(w3, account, nonce + i, calldata_size)
            signed_txs.append(signed_tx)
            print(f"  ✓ Signed tx {i+1}/{num_txs}")
        except Exception as e:
            print(f"  ✗ Failed to sign tx {i+1}: {str(e)[:100]}")
    
    print(f"\nStep 2: Broadcasting {len(signed_txs)} transactions rapidly (no delays)...")
    print("This maximizes the chance they'll be included in the same block!\n")
    
    tx_hashes = []
    start_time = time.time()
    
    for i, signed_tx in enumerate(signed_txs):
        try:
            tx_hash = w3.eth.send_raw_transaction(signed_tx.raw_transaction)
            tx_hashes.append(tx_hash)
            print(f"✓ Sent tx {i+1}/{len(signed_txs)}: {tx_hash.hex()[:20]}...")
            
            # NO DELAY - send as fast as possible to get them in one block!
            
        except Exception as e:
            print(f"✗ Failed to send tx {i+1}: {str(e)[:100]}")
            
            # Check if it's a block size limit error
            error_msg = str(e).lower()
            if "size" in error_msg or "limit" in error_msg or "too large" in error_msg:
                print("\n" + "!" * 80)
                print("BLOCK SIZE LIMIT DETECTED!")
                print("!" * 80)
                print(f"Error message: {str(e)}")
                return tx_hashes
            
            continue
    
    elapsed = time.time() - start_time
    print(f"\n✓ Sent {len(tx_hashes)} transactions in {elapsed:.3f} seconds")
    
    print(f"\nSent {len(tx_hashes)} transactions successfully")
    return tx_hashes

def wait_for_txs(w3, tx_hashes, timeout=300):
    """
    Wait for transactions to be mined and analyze the resulting block(s).
    """
    if not tx_hashes:
        print("No transactions to wait for")
        return
    
    print(f"\nWaiting for transactions to be mined (timeout: {timeout}s)...")
    
    start_time = time.time()
    mined_blocks = set()
    
    for i, tx_hash in enumerate(tx_hashes):
        print(f"Waiting for tx {i+1}/{len(tx_hashes)}: {tx_hash.hex()[:20]}...", end=" ")
        
        try:
            receipt = w3.eth.wait_for_transaction_receipt(tx_hash, timeout=timeout)
            print(f"✓ Block {receipt['blockNumber']}")
            mined_blocks.add(receipt['blockNumber'])
            
        except Exception as e:
            print(f"✗ Timeout or error: {str(e)[:50]}")
            continue
    
    elapsed = time.time() - start_time
    print(f"\nMining completed in {elapsed:.1f}s")
    print(f"Transactions mined in {len(mined_blocks)} block(s): {sorted(mined_blocks)}")
    
    # Analyze blocks
    print(f"\n{'=' * 80}")
    print("Block Analysis:")
    print(f"{'=' * 80}\n")
    
    for block_num in sorted(mined_blocks):
        analyze_block(w3, block_num)

def analyze_block(w3, block_number):
    """
    Analyze a block's size and transactions.
    """
    block = w3.eth.get_block(block_number, full_transactions=True)
    
    # Estimate RLP size (rough approximation)
    # In practice, you'd need to RLP-encode the entire block
    estimated_size = 0
    
    # Block header overhead (~500-1000 bytes)
    estimated_size += 1000
    
    # Transaction data
    for tx in block.transactions:
        # TX overhead (~200-300 bytes for signature, nonce, gas params, etc.)
        estimated_size += 300
        
        # Calldata
        if hasattr(tx, 'input') and tx.input:
            estimated_size += len(tx.input)
    
    size_mb = estimated_size / (1024 * 1024)
    exceeds = estimated_size > MAX_BLOCK_SIZE_BYTES
    
    print(f"Block #{block_number}:")
    print(f"  Transactions: {len(block.transactions)}")
    print(f"  Gas used: {block.gasUsed:,} / {block.gasLimit:,}")
    print(f"  Estimated RLP size: {estimated_size:,} bytes ({size_mb:.2f} MiB)")
    print(f"  Exceeds 10 MiB limit: {'YES ✓' if exceeds else 'NO ✗'}")
    
    if exceeds:
        print(f"  ⚠️  This block EXCEEDS the EIP-7934 limit!")
        print(f"  ⚠️  On a network with EIP-7934 active, this block would be REJECTED")
    
    print()

def main():
    """
    Main test execution.
    """
    print("\n")
    print("=" * 80)
    print("EIP-7934: RLP Block Size Limit Test")
    print("=" * 80)
    print(f"Target RPC: {RPC_URL}")
    print(f"Max block size (EIP-7934): {MAX_BLOCK_SIZE_BYTES:,} bytes ({MAX_BLOCK_SIZE_BYTES / (1024*1024):.1f} MiB)")
    print(f"Target size to exceed: {TARGET_SIZE_BYTES:,} bytes ({TARGET_SIZE_BYTES / (1024*1024):.1f} MiB)")
    print("=" * 80)
    
    try:
        # Connect to node
        w3 = Web3(Web3.HTTPProvider(RPC_URL))
        account = Account.from_key(PRIVATE_KEY)
        
        # Check connection
        if not w3.is_connected():
            print("\n❌ ERROR: Cannot connect to node!")
            print(f"Please ensure a node is running at {RPC_URL}")
            print("\nTo start a local test node, you can use:")
            print("  anvil")
            print("  or")
            print("  hardhat node")
            print("  or")
            print("  your surge node")
            return 1
        
        print(f"\n✓ Connected successfully!")
        print(f"Chain ID: {w3.eth.chain_id}")
        print(f"Account: {account.address}")
        balance = w3.eth.get_balance(account.address)
        print(f"Balance: {w3.from_wei(balance, 'ether')} ETH")
        
        if balance == 0:
            print("\n⚠️  WARNING: Account has zero balance!")
            print("   You need ETH to send transactions")
            print("   Fund this account or use a different private key")
            return 1
        
        # Calculate strategy
        strategy = calculate_optimal_tx_params()
        
        # Ask for confirmation (skip if input is piped/redirected)
        if sys.stdin.isatty():
            print(f"\nReady to send {strategy['num_txs']} transactions.")
            print("Press Enter to continue, or Ctrl+C to cancel...")
            input()
        else:
            print(f"\nAuto-starting: Sending {strategy['num_txs']} transactions...")
            time.sleep(1)
        
        # Send transactions
        tx_hashes = send_large_transactions(w3, account, strategy)
        
        # Wait and analyze
        if tx_hashes:
            wait_for_txs(w3, tx_hashes)
        
        print("\n" + "=" * 80)
        print("Test completed!")
        print("=" * 80)
        
        return 0
        
    except KeyboardInterrupt:
        print("\n\nTest cancelled by user")
        return 1
    except Exception as e:
        print(f"\n\n❌ ERROR: {e}")
        import traceback
        traceback.print_exc()
        return 1

if __name__ == "__main__":
    sys.exit(main())
