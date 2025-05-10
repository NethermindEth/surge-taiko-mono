#!/bin/bash

# Chain ID for Hoodoo testnet
CHAIN_ID=560048

# Function to verify a contract
verify_contract() {
  local contract_path=$1
  local contract_name=$2
  local address=$3
  
  echo "Verifying $contract_name at $address..."
  
  # Run the verification command
  forge verify-contract $address \
    $contract_path:$contract_name \
    --watch \
    --chain-id $CHAIN_ID
  
  # Check if verification succeeded
  if [ $? -eq 0 ]; then
    echo "✅ Verification of $contract_name complete!"
  else
    echo "❌ Verification of $contract_name failed!"
  fi
  echo "----------------------------------------"
}

# Define contract paths using plain variables instead of associative array
# This avoids bash version compatibility issues
CONTRACT_PATH_shared_address_manager="contracts/shared/common/AddressManager.sol:AddressManager"
CONTRACT_PATH_signal_service="contracts/shared/signal/SignalService.sol:SignalService"
CONTRACT_PATH_bridge="contracts/shared/bridge/Bridge.sol:Bridge"
CONTRACT_PATH_erc20_vault="contracts/shared/tokenvault/ERC20Vault.sol:ERC20Vault"
CONTRACT_PATH_erc721_vault="contracts/shared/tokenvault/ERC721Vault.sol:ERC721Vault"
CONTRACT_PATH_erc1155_vault="contracts/shared/tokenvault/ERC1155Vault.sol:ERC1155Vault"
CONTRACT_PATH_bridged_erc20="contracts/shared/tokenvault/BridgedERC20.sol:BridgedERC20"
CONTRACT_PATH_bridged_erc721="contracts/shared/tokenvault/BridgedERC721.sol:BridgedERC721"
CONTRACT_PATH_bridged_erc1155="contracts/shared/tokenvault/BridgedERC1155.sol:BridgedERC1155"
CONTRACT_PATH_rollup_address_manager="contracts/shared/common/AddressManager.sol:AddressManager"
CONTRACT_PATH_taiko="contracts/layer1/surge/SurgeHoodiTaikoL1.sol:SurgeHoodiTaikoL1"
CONTRACT_PATH_tier_sgx="contracts/layer1/verifiers/SgxVerifier.sol:SgxVerifier"
CONTRACT_PATH_tier_router="contracts/layer1/surge/common/SurgeTierRouter.sol:SurgeTierRouter"
CONTRACT_PATH_automata_dcap_attestation="contracts/layer1/automata-attestation/AutomataDcapV3Attestation.sol:AutomataDcapV3Attestation"
CONTRACT_PATH_risc0_groth16_verifier="node_modules/risc0-ethereum/contracts/src/groth16/RiscZeroGroth16Verifier.sol:RiscZeroGroth16Verifier"
CONTRACT_PATH_tier_zkvm_risc0="contracts/layer1/verifiers/Risc0Verifier.sol:Risc0Verifier"
CONTRACT_PATH_sp1_remote_verifier="node_modules/sp1-contracts/contracts/src/v4.0.0-rc.3/SP1VerifierPlonk.sol:SP1Verifier"
CONTRACT_PATH_tier_zkvm_sp1="contracts/layer1/verifiers/SP1Verifier.sol:SP1Verifier"
CONTRACT_PATH_tier_two_of_three="contracts/layer1/verifiers/compose/TwoOfThreeVerifier.sol:TwoOfThreeVerifier"
CONTRACT_PATH_SigVerifyLib="contracts/layer1/automata-attestation/utils/SigVerifyLib.sol:SigVerifyLib"
CONTRACT_PATH_PemCertChainLib="contracts/layer1/automata-attestation/lib/PEMCertChainLib.sol:PEMCertChainLib"

# Function to get contract path
get_contract_path() {
  local contract_name=$1
  local var_name="CONTRACT_PATH_${contract_name}"
  
  # Use indirect reference to get the value
  if [ -n "${!var_name}" ]; then
    echo "${!var_name}"
  else
    echo ""
  fi
}

# First part: Extract contract-address mappings
extract_contracts() {
  local input=$1
  local output_file="contract_addresses.txt"
  
  echo "Extracting contract addresses to $output_file..." > $output_file
  echo "=================================================" >> $output_file
  echo "" >> $output_file
  
  local current_contract=""
  
  while IFS= read -r line; do
    # Implementation addresses (proxy contracts)
    if [[ $line =~ ">"[[:space:]]+([^@]+)[[:space:]]+"@"[[:space:]]+([^[:space:]]+) ]]; then
      current_contract=${BASH_REMATCH[1]// /}
      echo "Contract: $current_contract" >> $output_file
    elif [[ $line =~ "impl"[[:space:]]*:[[:space:]]*([0-9xA-Fa-f]+) ]]; then
      local addr=${BASH_REMATCH[1]}
      echo "  Implementation address: $addr" >> $output_file
      
      local path=$(get_contract_path "$current_contract")
      if [[ -n "$path" ]]; then
        echo "  Path: $path" >> $output_file
      else
        echo "  Path: ⚠️ No path defined for $current_contract" >> $output_file
      fi
      echo "" >> $output_file
    
    # Direct addresses
    elif [[ $line =~ ">"[[:space:]]+([^@]+)[[:space:]]+"@"[[:space:]]+([^[:space:]]+)[[:space:]]+addr[[:space:]]*:[[:space:]]*([0-9xA-Fa-f]+) ]]; then
      current_contract=${BASH_REMATCH[1]// /}
      local addr=${BASH_REMATCH[3]}
      
      echo "Contract: $current_contract" >> $output_file
      if [[ ! $addr =~ ^0x7633750* ]]; then
        echo "  Address: $addr" >> $output_file
        
        local path=$(get_contract_path "$current_contract")
        if [[ -n "$path" ]]; then
          echo "  Path: $path" >> $output_file
        else
          echo "  Path: ⚠️ No path defined for $current_contract" >> $output_file
        fi
      else
        echo "  Address: $addr (skipped - precompile)" >> $output_file
      fi
      echo "" >> $output_file
      
    # Regular address lines
    elif [[ $line =~ "addr"[[:space:]]*:[[:space:]]*([0-9xA-Fa-f]+) ]]; then
      local addr=${BASH_REMATCH[1]}
      
      if [[ ! -z "$current_contract" ]]; then
        if [[ ! $addr =~ ^0x7633750* ]]; then
          echo "  Address: $addr" >> $output_file
          
          local path=$(get_contract_path "$current_contract")
          if [[ -n "$path" ]]; then
            echo "  Path: $path" >> $output_file
          else
            echo "  Path: ⚠️ No path defined for $current_contract" >> $output_file
          fi
        else
          echo "  Address: $addr (skipped - precompile)" >> $output_file
        fi
        echo "" >> $output_file
      fi
    
    # Special standalone libraries
    elif [[ $line =~ (SigVerifyLib|PemCertChainLib|HorseToken|BullToken)[[:space:]]+([0-9xA-Fa-f]+) ]]; then
      local lib_name=${BASH_REMATCH[1]}
      local lib_address=${BASH_REMATCH[2]}
      
      echo "Contract: $lib_name" >> $output_file
      echo "  Address: $lib_address" >> $output_file
      
      local path=$(get_contract_path "$lib_name")
      if [[ -n "$path" ]]; then
        echo "  Path: $path" >> $output_file
      else
        echo "  Path: ⚠️ No path defined for $lib_name" >> $output_file
      fi
      echo "" >> $output_file
    fi
    
  done <<< "$input"
  
  echo "Extraction complete. Addresses saved to $output_file"
  cat $output_file
}

# Second part: Verify the contracts
verify_extracted_contracts() {
  local input_file="contract_addresses.txt"
  
  if [[ ! -f $input_file ]]; then
    echo "Error: $input_file not found. Run extraction first."
    return 1
  fi
  
  echo "Starting verification based on extracted addresses..."
  echo "====================================================="
  
  local current_contract=""
  local address=""
  local path=""
  local contract_name=""
  
  while IFS= read -r line; do
    if [[ $line =~ ^Contract:[[:space:]]+(.*)$ ]]; then
      current_contract=${BASH_REMATCH[1]}
    elif [[ $line =~ ^[[:space:]]+Implementation[[:space:]]address:[[:space:]]+(.*)$ || $line =~ ^[[:space:]]+Address:[[:space:]]+(.*)$ ]]; then
      address=${BASH_REMATCH[1]}
      if [[ $address =~ "skipped" ]]; then
        continue  # Skip precompile addresses
      fi
    elif [[ $line =~ ^[[:space:]]+Path:[[:space:]]+([^:]+):([^[:space:]]+)$ ]]; then
      path=${BASH_REMATCH[1]}
      contract_name=${BASH_REMATCH[2]}
      # Verify when we have all needed information
      if [[ ! -z "$current_contract" && ! -z "$address" && ! -z "$path" && ! -z "$contract_name" ]]; then
        verify_contract "$path" "$contract_name" "$address"
        # Reset variables for next contract
        current_contract=""
        address=""
        path=""
        contract_name=""
      fi
    fi
  done < "$input_file"
  
  echo "All verification attempts completed!"
}

# Main script
echo "Contract Verification Script"
echo "============================"
echo "Chain ID: $CHAIN_ID"
echo "============================"

# Check if a deployment log file was provided
if [ $# -ne 1 ]; then
  echo "Usage: $0 <deployment-log-file>"
  echo "Please provide the path to the deployment log file."
  exit 1
fi

DEPLOYMENT_LOG_FILE="$1"

# Check if the file exists
if [ ! -f "$DEPLOYMENT_LOG_FILE" ]; then
  echo "Error: Deployment log file '$DEPLOYMENT_LOG_FILE' not found."
  exit 1
fi

echo "Using deployment log from: $DEPLOYMENT_LOG_FILE"
echo ""

# Read the deployment output from file
DEPLOYMENT_OUTPUT=$(cat "$DEPLOYMENT_LOG_FILE")

# Start verification
extract_contracts "$DEPLOYMENT_OUTPUT"
echo ""
echo "Press Enter to continue with verification, or Ctrl+C to cancel..."
read

verify_extracted_contracts