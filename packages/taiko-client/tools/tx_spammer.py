import time
from web3 import Web3
import os
from dotenv import load_dotenv
import argparse

# Load environment variables from .env file
load_dotenv()

# Get required environment variables
private_key = os.getenv('PRIVATE_KEY')
if not private_key:
    raise Exception("Environment variable PRIVATE_KEY not set")

recipient = os.getenv('RECIPIENT_ADDRESS')
if not recipient:
    raise Exception("Environment variable RECIPIENT_ADDRESS not set")

# Get configurable environment variables with defaults
frequency = int(os.getenv('TX_FREQUENCY_SECONDS', 900))  # Default 15 minutes
amount = int(os.getenv('TX_AMOUNT_WEI', 1))  # Default 1 wei
rpc_url = os.getenv('RPC_URL', 'https://l2-rpc.surge.staging-nethermind.xyz')  # Default RPC URL

parser = argparse.ArgumentParser(description='Spam transactions on the Surge network.')
parser.add_argument('--rpc', type=str, default=rpc_url, help='RPC URL for the Surge network')
args = parser.parse_args()

# Connect to the Surge network
w3 = Web3(Web3.HTTPProvider(args.rpc))

# Check if connected
if not w3.is_connected():
    raise Exception("Failed to connect to the Surge network")

# Get the account from the private key
account = w3.eth.account.from_key(private_key)

print("Starting tx spammer...")
print(f"Using RPC URL: {args.rpc}")
print(f"Account address: {account.address}")
print(f"Recipient address: {recipient}")
print(f"Transaction frequency: {frequency} seconds")
print(f"Transaction amount: {amount} wei")

def send_transaction(nonce: int):
    try:
        tx = {
            'nonce': nonce,
            'to': recipient,
            'value': amount,
            'gasPrice': w3.to_wei('10', 'gwei'),
            'chainId': w3.eth.chain_id
        }
        
        print(f'Preparing transaction with nonce {nonce}')
        
        # Let the network estimate the gas
        gas_estimate = w3.eth.estimate_gas({
            'to': recipient,
            'value': amount,
            'from': account.address
        })
        tx['gas'] = gas_estimate
        print(f'Estimated gas: {gas_estimate}')

        print(f'Sending transaction of {amount} wei')
        print(f'From: {account.address} To: {recipient}')
        signed_tx = w3.eth.account.sign_transaction(tx, private_key)
        tx_hash = w3.eth.send_raw_transaction(signed_tx.rawTransaction)
        print(f'Transaction sent: {tx_hash.hex()}')
    except Exception as e:
        print(f'Error in send_transaction: {str(e)}')
        raise e

def spam_transactions():
    print("Starting spam_transactions loop...")
    while True:
        try:
            print("Getting nonce...")
            nonce = w3.eth.get_transaction_count(account.address)
            print(f"Current nonce: {nonce}")
            send_transaction(nonce)
            print(f"Waiting {frequency} seconds before next transaction...")
            time.sleep(frequency)
        except Exception as e:
            print(f"Error occurred: {str(e)}")
            print(f"Retrying in {frequency} seconds...")
            time.sleep(frequency)

if __name__ == '__main__':
    try:
        print("Script started")
        spam_transactions()
    except Exception as e:
        print(f"Fatal error: {str(e)}")