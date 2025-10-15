// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/src/Script.sol";
import "forge-std/src/console2.sol";

// Solady for JSON parsing
import "solady/src/utils/JSONParserLib.sol";
import "solady/src/utils/LibString.sol";

// Layer 1 contracts
import "src/layer1/verifiers/AzureTdxVerifier.sol";
import { AzureTDX } from "azure-tdx-verifier/AzureTDX.sol";
import "src/shared/libs/LibStrings.sol";
import "test/shared/DeployCapability.sol";

// TDX Automata interfaces
interface IPcsDao {
    function owner() external view returns (address);
    function transferOwnership(address) external;
    function upsertPcsCertificates(uint8 ca, bytes calldata cert) external returns (bytes32 attestationId);
}

interface IFmspcTcbDao {
    function owner() external view returns (address);
    function transferOwnership(address) external;
    function upsertFmspcTcb(TcbInfoJsonObj memory tcbInfoJson) external;
}

interface IAutomataEnclaveIdentityDao {
    function owner() external view returns (address);
    function transferOwnership(address) external;
    function upsertEnclaveIdentity(
        uint256 id, 
        uint256 isvsvn, 
        EnclaveIdentityJsonObj memory identityJson
    ) external;
}

interface IEnclaveIdentityHelper {
    function parseIdentityString(string memory identityStr) 
        external pure returns (IdentityObj memory identity, string memory success);
}

// TDX data structures
struct TcbInfoJsonObj {
    string tcbInfoStr;
    bytes signature;
}

struct EnclaveIdentityJsonObj {
    string identityStr;
    bytes signature;
}

enum EnclaveId {
    QE,
    QVE,
    TD_QE
}

struct IdentityObj {
    EnclaveId id;
    uint32 version;
    uint64 issueDateTimestamp;
    uint64 nextUpdateTimestamp;
    uint32 tcbEvaluationDataNumber;
    bytes4 miscselect;
    bytes4 miscselectMask;
    bytes16 attributes;
    bytes16 attributesMask;
    bytes32 mrsigner;
    uint16 isvprodid;
    Tcb[] tcb;
}

enum EnclaveIdTcbStatus {
    SGX_ENCLAVE_REPORT_ISVSVN_NOT_SUPPORTED,
    OK,
    SGX_ENCLAVE_REPORT_ISVSVN_REVOKED,
    SGX_ENCLAVE_REPORT_ISVSVN_OUT_OF_DATE
}

struct Tcb {
    uint16 isvsvn;
    uint256 dateTimestamp;
    EnclaveIdTcbStatus status;
}

/// @title SetupAzureTDXVerifier
/// @notice Script to setup TDX verifier with attestation configuration and transfer ownership
contract SetupAzureTDXVerifier is Script, DeployCapability {
    using JSONParserLib for JSONParserLib.Item;
    using LibString for string;

    // Configuration
    uint256 internal immutable privateKey = vm.envUint("PRIVATE_KEY");

    // TDX verifier configuration
    address internal immutable tdxVerifierAddress = vm.envAddress("AZURE_TDX_VERIFIER_ADDRESS");
    
    // TDX Automata contract addresses (required)
    address internal immutable tdxPcsDao = vm.envAddress("TDX_PCS_DAO_ADDRESS");
    address internal immutable tdxFmspcTcbDao = vm.envAddress("TDX_FMSPC_TCB_DAO_ADDRESS");
    address internal immutable tdxEnclaveIdentityDao = vm.envAddress("TDX_ENCLAVE_IDENTITY_DAO_ADDRESS");
    address internal immutable tdxEnclaveIdentityHelper = vm.envAddress("TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS");

    // TDX attestation configuration
    bytes internal tdxTrustedParamsBytes = vm.envOr("AZURE_TDX_TRUSTED_PARAMS_BYTES", bytes(""));
    bytes internal tdxQuoteBytes = vm.envOr("AZURE_TDX_QUOTE_BYTES", bytes(""));

    // Ownership transfer
    address internal immutable newOwner = vm.envAddress("NEW_OWNER");

    modifier broadcast() {
        require(privateKey != 0, "invalid private key");
        vm.startBroadcast(privateKey);
        _;
        vm.stopBroadcast();
    }

    function run() external broadcast {
        require(tdxVerifierAddress != address(0), "config: AZURE_TDX_VERIFIER_ADDRESS");
        require(tdxPcsDao != address(0), "config: TDX_PCS_DAO_ADDRESS");
        require(tdxFmspcTcbDao != address(0), "config: TDX_FMSPC_TCB_DAO_ADDRESS");
        require(tdxEnclaveIdentityDao != address(0), "config: TDX_ENCLAVE_IDENTITY_DAO_ADDRESS");
        require(tdxEnclaveIdentityHelper != address(0), "config: TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS");
        require(newOwner != address(0), "config: NEW_OWNER");

        AzureTdxVerifier tdxVerifier = AzureTdxVerifier(tdxVerifierAddress);

        // Verify current ownership
        require(tdxVerifier.owner() == msg.sender, "SetupAzureTDXVerifier: tdx verifier not owner");

        // Setup TDX trusted parameters if provided
        if (tdxTrustedParamsBytes.length > 0) {
            AzureTdxVerifier.TrustedParams memory params = 
                abi.decode(tdxTrustedParamsBytes, (AzureTdxVerifier.TrustedParams));
            tdxVerifier.setTrustedParams(0, params);
            console2.log("** TDX_TRUSTED_PARAMS configured");
        }

        // Setup TDX collaterals
        setupTDXCollaterals();

        // Register TDX instance with quote if provided
        if (tdxQuoteBytes.length > 0) {
            // Log the instance id to Json
            vm.writeJson(
                vm.serializeUint(
                    "tdx_instance_ids",
                    "tdx_instance_id",
                    tdxVerifier.nextInstanceId()
                ),
                string.concat(vm.projectRoot(), "/deployments/tdx_instances.json")
            );

            // Log the instance id to console
            uint256 instanceId = tdxVerifier.nextInstanceId();

            // Register instance
            AzureTDX.VerifyParams memory verifyParams = 
                abi.decode(tdxQuoteBytes, (AzureTDX.VerifyParams));
            tdxVerifier.registerInstance(instanceId, verifyParams);
            console2.log("** TDX instance registered with ID:", instanceId);
        }

        // Transfer ownership
        tdxVerifier.transferOwnership(newOwner);
        console2.log("** AzureTdxVerifier ownership transferred to:", newOwner);

        console2.log("** TDX verifier setup complete **");
    }

    /// @dev Setup TDX collaterals (PCS certificates, enclave identity, TCB info)
    function setupTDXCollaterals() internal {
        // Configure PCS certificates if path provided
        string memory rootPcsCertPath = vm.envOr("AZURE_TDX_ROOT_PCS_CERT_PATH", string(""));
        if (bytes(rootPcsCertPath).length > 0) {
            bytes memory certBytes = vm.parseBytes(
                vm.readFile(string.concat(vm.projectRoot(), rootPcsCertPath))
            );
            
            // Use CA.ROOT (0)
            IPcsDao(tdxPcsDao).upsertPcsCertificates(0, certBytes);
            console2.log("** TDX_ROOT_PCS_CERTIFICATES configured");
        }
        string memory pcsCertPath = vm.envOr("AZURE_TDX_PCS_CERT_PATH", string(""));
        if (bytes(pcsCertPath).length > 0) {
            bytes memory certBytes = vm.parseBytes(
                vm.readFile(string.concat(vm.projectRoot(), pcsCertPath))
            );
            
            // Use CA.SIGNING (3)
            IPcsDao(tdxPcsDao).upsertPcsCertificates(3, certBytes);
            console2.log("** TDX_PCS_CERTIFICATES configured");
        }

        // Configure enclave identity if path provided
        string memory enclaveIdentityPath = vm.envOr("AZURE_TDX_QE_IDENTITY_PATH", string(""));
        if (bytes(enclaveIdentityPath).length > 0) {
            string memory enclaveIdentityJson = vm.readFile(
                string.concat(vm.projectRoot(), enclaveIdentityPath)
            );
            EnclaveIdentityJsonObj memory identityJsonObj = parseEnclaveIdentityJson(enclaveIdentityJson);
            
            (IdentityObj memory identity,) = IEnclaveIdentityHelper(tdxEnclaveIdentityHelper)
                .parseIdentityString(identityJsonObj.identityStr);

            // Use isvsvn = 4 for TDX QE
            IAutomataEnclaveIdentityDao(tdxEnclaveIdentityDao).upsertEnclaveIdentity(
                uint256(identity.id), 4, identityJsonObj
            );
            console2.log("** TDX_QE_IDENTITY configured");
        }

        // Configure TCB info if path provided  
        string memory tcbInfoPath = vm.envOr("AZURE_TDX_TCB_INFO_PATH", string(""));
        if (bytes(tcbInfoPath).length > 0) {
            string memory tcbInfoJson = vm.readFile(
                string.concat(vm.projectRoot(), tcbInfoPath)
            );
            TcbInfoJsonObj memory tcbInfoJsonObj = parseTcbInfoJson(tcbInfoJson);

            IFmspcTcbDao(tdxFmspcTcbDao).upsertFmspcTcb(tcbInfoJsonObj);
            console2.log("** TDX_TCB_INFO configured");
        }
    }

    /// @dev Parse enclave identity JSON to extract identity string and signature
    function parseEnclaveIdentityJson(string memory jsonStr) 
        internal pure returns (EnclaveIdentityJsonObj memory result) {
        JSONParserLib.Item memory root = JSONParserLib.parse(jsonStr);
        JSONParserLib.Item[] memory children = root.children();
        
        for (uint256 i = 0; i < root.size(); i++) {
            string memory key = children[i].key();
            if (LibString.eq(key, "\"enclaveIdentity\"")) {
                result.identityStr = children[i].value();
            } else if (LibString.eq(key, "\"signature\"")) {
                result.signature = vm.parseBytes(JSONParserLib.decodeString(children[i].value()));
            }
        }
    }

    /// @dev Parse TCB info JSON to extract TCB info string and signature  
    function parseTcbInfoJson(string memory jsonStr) 
        internal pure returns (TcbInfoJsonObj memory result) {
        JSONParserLib.Item memory root = JSONParserLib.parse(jsonStr);
        JSONParserLib.Item[] memory children = root.children();
        
        for (uint256 i = 0; i < root.size(); i++) {
            string memory key = children[i].key();
            if (LibString.eq(key, "\"tcbInfo\"")) {
                result.tcbInfoStr = children[i].value();
            } else if (LibString.eq(key, "\"signature\"")) {
                result.signature = vm.parseBytes(JSONParserLib.decodeString(children[i].value()));
            }
        }
    }
}