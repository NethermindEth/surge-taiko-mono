r"""
This script is used to deposit bonds for the proposer and prover on the Taiko network.
It reads the private keys and addresses from a .env file,
connects to the L1 network, and sends a specified amount of ETH to the deposit contract addresses.

Setup:
1. Create a virtual environment:
   python -m venv venv

2. Activate the virtual environment:
   - On Windows: venv\Scripts\activate
   - On macOS/Linux: source venv/bin/activate

3. Install the required dependencies:
   pip install -r requirements.txt

4. Create a .env file in the tools/tx_spammer directory with the following content:
   L1_PROPOSER_PRIVATE_KEY=<your_private_key>
   L1_PROVER_PRIVATE_KEY=<your_private_key>
   TAIKO_L1_CONTRACT_ADDRESS=<taiko_l1_contract_address>

5. Run the script:
   python deposit_bonds.py [--amount AMOUNT] [--rpc RPC_URL]

CLI Parameters:
--amount: Amount of ETH to send for deposit bonds
--rpc: RPC URL for the L1 network
"""

import json
from web3 import Web3
import os
from dotenv import load_dotenv
import argparse
from eth_account import Account

load_dotenv()

account_private_key = os.getenv('ACCOUNT_PRIVATE_KEY')
if not account_private_key:
    raise Exception("Environment variable ACCOUNT_PRIVATE_KEY not set")

taiko_l1_contract_address = os.getenv('TAIKO_L1_CONTRACT_ADDRESS')
if not taiko_l1_contract_address:
    raise Exception("Environment variable TAIKO_L1_CONTRACT_ADDRESS not set")

def parse_arguments():

    parser = argparse.ArgumentParser(
        description='Deposit bonds for Taiko network proposer and prover',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter
    )
    
    parser.add_argument(
        '--amount',
        type=float,
        default=10.0,
        help='Amount of ETH to send for deposit bonds'
    )
    
    parser.add_argument(
        '--rpc',
        type=str,
        required=True,
        help='RPC URL for the L1 network'
    )
    
    return parser.parse_args()

args = parse_arguments()

w3 = Web3(Web3.HTTPProvider(args.rpc))

# Check if connected
if not w3.is_connected():
    raise Exception("Failed to connect to the L1 network")

# Get the account from the private key
deposit_account = Account.from_key(account_private_key)
deposit_amount = w3.to_wei(args.amount, 'ether')

ABI = "{\"type\":\"function\",\"name\":\"depositBond\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"}"

def send_deposit_bond_transaction(nonce: int, private_key : str, account: Account):
    # Parse the ABI
    contract_abi = json.loads(ABI)

    # Create contract instance
    contract = w3.eth.contract(address=taiko_l1_contract_address, abi=[contract_abi])

    # Prepare the transaction
    tx = contract.functions.depositBond().build_transaction({
        'nonce': nonce,
        'value': deposit_amount,
        'gas': 100000,  # Adjust gas limit as needed
        'gasPrice': w3.to_wei('10', 'gwei'),
        'chainId': w3.eth.chain_id
    })

    print(f'Sending depositBond transaction: {tx} by RPC: {args.rpc}')
    print(f'Sending from: {account.address}')

    # Sign and send the transaction
    signed_tx = w3.eth.account.sign_transaction(tx, private_key)
    tx_hash = w3.eth.send_raw_transaction(signed_tx.raw_transaction)
    print(f'Transaction sent: {tx_hash.hex()}')

def deposit_bonds():
    nonce = w3.eth.get_transaction_count(deposit_account.address)
    send_deposit_bond_transaction(nonce, account_private_key, deposit_account)

deposit_bonds()