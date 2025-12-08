#!/usr/bin/env python3
"""
Test EIP-7910: eth_config JSON-RPC method
"""

from web3 import Web3
import json
import sys

# Configuration
RPC_URL = "http://localhost:8547"

def main():
    """Test if eth_config method works"""
    print("Testing EIP-7910: eth_config method")
    print(f"RPC URL: {RPC_URL}\n")
    
    w3 = Web3(Web3.HTTPProvider(RPC_URL))
    
    if not w3.is_connected():
        print("❌ Cannot connect to node!")
        return 1
    
    print(f"✓ Connected: {w3.client_version}\n")
    
    try:
        response = w3.provider.make_request("eth_config", [])
        
        if "result" in response:
            print("✓ eth_config method is SUPPORTED!\n")
            print(json.dumps(response["result"], indent=2))
            return 0
            
        elif "error" in response:
            error = response["error"]
            print(f"✗ Error: {error.get('message')}")
            
            if error.get('code') == -32601:
                print("⚠️  Method not found - EIP-7910 not implemented")
            
            return 1
            
    except Exception as e:
        print(f"✗ Failed: {str(e)}")
        return 1

if __name__ == "__main__":
    sys.exit(main())

if __name__ == "__main__":
    try:
        sys.exit(main())
    except KeyboardInterrupt:
        print("\n\nTest cancelled by user")
        sys.exit(1)
    except Exception as e:
        print(f"\n\nERROR: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
