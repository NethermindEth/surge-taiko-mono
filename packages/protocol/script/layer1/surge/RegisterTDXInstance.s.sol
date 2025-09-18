// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

// Foundry
import "forge-std/src/Script.sol";

// Solady for JSON parsing
import "solady/src/utils/JSONParserLib.sol";
import "solady/src/utils/LibString.sol";

// Layer 1 contracts
import "contracts/layer1/verifiers/AzureTdxVerifier.sol";
import { AzureTDX } from "azure-tdx-verifier/AzureTDX.sol";

// TDX Automata interfaces
interface IPcsDao {
    function upsertPcsCertificates(uint8 ca, bytes calldata cert) external returns (bytes32 attestationId);
}

interface IFmspcTcbDao {
    function upsertFmspcTcb(TcbInfoJsonObj memory tcbInfoJson) external;
}

interface IAutomataEnclaveIdentityDao {
    function upsertEnclaveIdentity(uint256 id, uint256 isvsvn, EnclaveIdentityJsonObj memory identityJson) external;
}

interface IEnclaveIdentityHelper {
    function parseIdentityString(string memory identityStr) 
        external pure returns (IdentityObj memory identity, bool success);
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

struct IdentityObj {
    uint256 id;
}

contract RegisterTDXInstance is Script {
    using JSONParserLib for JSONParserLib.Item;
    using LibString for string;

    // Execution configuration
    // ---------------------------------------------------------------------------------------------
    uint256 internal immutable privateKey = vm.envUint("PRIVATE_KEY");

    // TDX configuration
    // ---------------------------------------------------------------------------------------------
    address public azureTdxVerifier = vm.envAddress("AZURE_TDX_VERIFIER");
    
    // TDX Automata contract addresses
    address internal immutable tdxPcsDao = vm.envAddress("TDX_PCS_DAO_ADDRESS");
    address internal immutable tdxFmspcTcbDao = vm.envAddress("TDX_FMSPC_TCB_DAO_ADDRESS");
    address internal immutable tdxEnclaveIdentityDao = vm.envAddress("TDX_ENCLAVE_IDENTITY_DAO_ADDRESS");
    address internal immutable tdxEnclaveIdentityHelper = vm.envAddress("TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS");

    // TDX trusted parameters and quote
    bytes internal tdxTrustedParamsBytes = vm.envBytes("TDX_TRUSTED_PARAMS_BYTES");
    bytes internal tdxQuoteBytes = vm.envBytes("TDX_QUOTE_BYTES");

    modifier broadcast() {
        require(privateKey != 0, "invalid private key");
        vm.startBroadcast(privateKey);
        _;
        vm.stopBroadcast();
    }

    function run() external broadcast {
        AzureTdxVerifier tdxVerifier = AzureTdxVerifier(azureTdxVerifier);

        if (tdxTrustedParamsBytes.length > 0) {
            AzureTdxVerifier.TrustedParams memory params = 
                abi.decode(tdxTrustedParamsBytes, (AzureTdxVerifier.TrustedParams));
            tdxVerifier.setTrustedParams(0, params);
            console2.log("** TDX trusted params configured");
        }

        setupTDXCollaterals();

        if (tdxQuoteBytes.length > 0) {
            vm.writeJson(
                vm.serializeUint(
                    "tdx_instance_ids",
                    "tdx_instance_id",
                    tdxVerifier.nextInstanceId()
                ),
                string.concat(vm.projectRoot(), "/deployments/tdx_instances.json")
            );

            AzureTDX.VerifyParams memory verifyParams = 
                abi.decode(tdxQuoteBytes, (AzureTDX.VerifyParams));
            tdxVerifier.registerInstance(0, verifyParams);
            console2.log("** TDX instance registered with quote");
        }
    }

    function setupTDXCollaterals() internal {
        string memory pcsCertPath = vm.envOr("TDX_PCS_CERT_PATH", string(""));
        if (tdxPcsDao != address(0) && bytes(pcsCertPath).length > 0) {
            bytes memory certBytes = vm.parseBytes(vm.readFile(string.concat(vm.projectRoot(), pcsCertPath)));
            
            // Use CA.SIGNING (0) as default
            IPcsDao(tdxPcsDao).upsertPcsCertificates(0, certBytes);
            console2.log("** TDX PCS certificates configured");
        }

        string memory enclaveIdentityPath = vm.envOr("TDX_QE_IDENTITY_PATH", string(""));
        if (bytes(enclaveIdentityPath).length > 0 && tdxEnclaveIdentityDao != address(0) 
            && tdxEnclaveIdentityHelper != address(0)) {
            string memory enclaveIdentityJson = vm.readFile(string.concat(vm.projectRoot(), enclaveIdentityPath));
            EnclaveIdentityJsonObj memory identityJsonObj = parseEnclaveIdentityJson(enclaveIdentityJson);
            
            (IdentityObj memory identity, bool success) = IEnclaveIdentityHelper(tdxEnclaveIdentityHelper)
                .parseIdentityString(identityJsonObj.identityStr);
            require(success, "RegisterTDXInstance: failed to parse enclave identity");

            // Use isvsvn = 4 for TDX QE
            IAutomataEnclaveIdentityDao(tdxEnclaveIdentityDao).upsertEnclaveIdentity(
                identity.id, 4, identityJsonObj
            );
            console2.log("** TDX enclave identity configured");
        }

        string memory tcbInfoPath = vm.envOr("TDX_TCB_INFO_PATH", string(""));
        if (bytes(tcbInfoPath).length > 0 && tdxFmspcTcbDao != address(0)) {
            string memory tcbInfoJson = vm.readFile(string.concat(vm.projectRoot(), tcbInfoPath));
            TcbInfoJsonObj memory tcbInfoJsonObj = parseTcbInfoJson(tcbInfoJson);
            
            IFmspcTcbDao(tdxFmspcTcbDao).upsertFmspcTcb(tcbInfoJsonObj);
            console2.log("** TDX TCB info configured");
        }
    }

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

    function parseTcbInfoJson(string memory jsonStr) internal pure returns (TcbInfoJsonObj memory result) {
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