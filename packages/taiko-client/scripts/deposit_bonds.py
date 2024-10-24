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
   PROPOSER_DEPOSIT_CONTRACT_ADDRESS=<proposer_deposit_contract_address>
   PROVER_DEPOSIT_CONTRACT_ADDRESS=<prover_deposit_contract_address>
   RECIPIENT_ADDRESS=<recipient_address>

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

# Load environment variables from .env file
load_dotenv()

l1_proposer_private_key = os.getenv('L1_PROPOSER_PRIVATE_KEY')
if not l1_proposer_private_key:
    raise Exception("Environment variable L1_PROPOSER_PRIVATE_KEY not set")

l1_prover_private_key = os.getenv('L1_PROVER_PRIVATE_KEY')
if not l1_prover_private_key:
    raise Exception("Environment variable L1_PROVER_PRIVATE_KEY not set")

proposer_deposit_contract_address = os.getenv('PROPOSER_DEPOSIT_CONTRACT_ADDRESS')
if not proposer_deposit_contract_address:
    raise Exception("Environment variable PROPOSER_DEPOSIT_CONTRACT_ADDRESS not set")

prover_deposit_contract_address = os.getenv('PROVER_DEPOSIT_CONTRACT_ADDRESS')
if not prover_deposit_contract_address:
    raise Exception("Environment variable PROVER_DEPOSIT_CONTRACT_ADDRESS not set")

recipient = os.getenv('RECIPIENT_ADDRESS')
if not recipient:
    raise Exception("Environment variable RECIPIENT_ADDRESS not set")

parser = argparse.ArgumentParser(description='Spam transactions on the Taiko network.')
parser.add_argument('--amount', type=float, default=1000, help='Amount of ETH to send for deposit bonds')
parser.add_argument('--rpc', type=str, help='RPC URL for the L1 network')
args = parser.parse_args()

import pdb; pdb.set_trace()
w3 = Web3(Web3.HTTPProvider(args.rpc))

# Check if connected
if not w3.is_connected():
    raise Exception("Failed to connect to the L1 network")

# Get the account from the private key
proposer_account = Account.from_key(l1_proposer_private_key)
prover_account = Account.from_key(l1_prover_private_key)
deposit_amount = w3.to_wei(args.amount, 'ether')

ABI = "{\"type\":\"function\",\"name\":\"depositBond\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"}"

def send_deposit_bond_transaction(nonce: int, private_key : str, account: Account):
    # Parse the ABI
    contract_abi = json.loads(ABI)

    # Create contract instance
    contract = w3.eth.contract(address=recipient, abi=[contract_abi])

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
    tx_hash = w3.eth.send_raw_transaction(signed_tx.rawTransaction)
    print(f'Transaction sent: {tx_hash.hex()}')

def send_transaction(nonce : int, private_key : str, account: Account, contract_function: str):
    tx = {
        'nonce': nonce,
        'to': recipient,
        'value': deposit_amount,
        'gas': 21000,
        'gasPrice': w3.to_wei('10', 'gwei'),
        'chainId': w3.eth.chain_id
    }
    print(f'Sending transaction: {tx} by RPC: {args.rpc}')
    print(f'Sending from: {account.address}')
    signed_tx = w3.eth.account.sign_transaction(tx, private_key)
    tx_hash = w3.eth.send_raw_transaction(signed_tx.raw_transaction)
    print(f'Transaction sent: {tx_hash.hex()}')

def deposit_for_proposer():
    nonce = w3.eth.get_transaction_count(proposer_account.address)
    send_deposit_bond_transaction(nonce, l1_proposer_private_key, proposer_account)

def deposit_for_prover():
    nonce = w3.eth.get_transaction_count(prover_account.address)
    send_deposit_bond_transaction(nonce, l1_prover_private_key, prover_account)

def deposit_bonds():
    deposit_for_proposer()
    deposit_for_prover()

deposit_bonds()