// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "../automata-attestation/utils/BytesUtils.sol";
import "src/shared/common/EssentialContract.sol";
import "src/shared/libs/LibStrings.sol";
import "../automata-attestation/lib/QuoteV3Auth/V3Struct.sol";
import "../based/ITaikoInbox.sol";
import "./LibPublicInput.sol";
import "./IVerifier.sol";

interface IAutomataDcapAttestation {
    function verifyAndAttestOnChain(bytes calldata rawQuote)
        external
        payable
        returns (bool, bytes memory);
}

/// @title TdxVerifier
/// @notice This contract is the implementation of verifying TDX signature proofs
/// onchain.
contract TdxVerifier is EssentialContract, IVerifier {
    /// @dev Each public-private key pair (Ethereum address) is generated within
    /// the TDX program when on bootstrap. The off-chain remote attestation
    /// ensures the validity of the program hash and has the capability of
    /// bootstrapping the network with trustworthy instances.
    struct Instance {
        address addr;
        uint64 validSince;
    }

    /// @dev Parameters for trusted TDX instances
    struct TrustedParams {
        bytes16 teeTcbSvn;
        bytes mrSeam;
        bytes mrTd;
        bytes rtMr0;
        bytes rtMr1;
        bytes rtMr2;
        bytes rtMr3;
    }

    /// @dev The parameters required for verification
    struct VerifyParams {
        bytes quote;
        bytes userData;
        bytes nonce;
    }

    /// @notice The expiry time for the TDX instance.
    uint64 public constant INSTANCE_EXPIRY = 365 days;

    /// @notice A security feature, a delay until an instance is enabled when using onchain RA
    /// verification
    uint64 public constant INSTANCE_VALIDITY_DELAY = 0;

    uint64 public immutable taikoChainId;
    address public immutable taikoInbox;
    address public immutable taikoProofVerifier;
    address public immutable automataDcapAttestation;

    /// @dev For gas savings, we shall assign each TDX instance with an id that when we need to
    /// set a new pub key, just write storage once.
    /// Slot 1.
    uint256 public nextInstanceId;

    /// @dev One TDX instance is uniquely identified (on-chain) by it's ECDSA public key
    /// (or rather ethereum address). Once that address is used (by proof verification) it has to be
    /// overwritten by a new one (representing the same instance). This is due to side-channel
    /// protection. Also this public key shall expire after some time.
    /// Slot 2.
    mapping(uint256 instanceId => Instance instance) public instances;

    /// @dev One address shall be registered (during attestation) only once, otherwise it could
    /// bypass this contract's expiry check by always registering with the same attestation and
    /// getting multiple valid instanceIds. While during proving, it is technically possible to
    /// register the old addresses, it is less of a problem, because the instanceId would be the
    /// same for those addresses and if deleted - the attestation cannot be reused anyways.
    /// Slot 3.
    mapping(address instanceAddress => bool alreadyAttested) public addressRegistered;

    /// @dev Indicates whether a quote nonce hash has been used or not.
    /// Slot 4.
    mapping(bytes32 nonceHash => bool isUsed) public nonceUsed;

    /// @dev The trusted parameters for trusted TDX instances
    /// Slot 5.
    mapping(uint256 index => TrustedParams trustedParams) public trustedParams;

    uint256[45] private __gap;

    /// @notice Emitted when a new TDX instance is added to the registry, or replaced.
    /// @param id The ID of the TDX instance.
    /// @param instance The address of the TDX instance.
    /// @param replaced The address of the TDX instance that was replaced. If it is the first
    /// instance, this value is zero address.
    /// @param validSince The time since the instance is valid.
    event InstanceAdded(
        uint256 indexed id, address indexed instance, address indexed replaced, uint256 validSince
    );

    /// @notice Emitted when an TDX instance is deleted from the registry.
    /// @param id The ID of the TDX instance.
    /// @param instance The address of the TDX instance.
    event InstanceDeleted(uint256 indexed id, address indexed instance);

    /// @notice Emitted when trusted params are updated
    /// @param index The index of the trusted params
    /// @param params The trusted params
    event TrustedParamsUpdated(uint256 indexed index, TrustedParams params);

    error TDX_ALREADY_ATTESTED();
    error TDX_INVALID_ATTESTATION();
    error TDX_INVALID_INSTANCE();
    error TDX_INVALID_PROOF();
    error TDX_INVALID_TRUSTED_PARAMS();
    error TDX_INVALID_VERSION_TYPE();
    error TDX_INVALID_TCB_SVN();
    error TDX_INVALID_MR_SEAM();
    error TDX_INVALID_MR_TD();
    error TDX_INVALID_RTMR();
    error TDX_INVALID_REPORT_DATA();

    constructor(
        uint64 _taikoChainId,
        address _taikoInbox,
        address _taikoProofVerifier,
        address _automataDcapAttestation
    )
        EssentialContract(address(0))
    {
        taikoChainId = _taikoChainId;
        taikoInbox = _taikoInbox;
        taikoProofVerifier = _taikoProofVerifier;
        automataDcapAttestation = _automataDcapAttestation;
    }

    /// @notice Initializes the contract.
    /// @param _owner The owner of this contract. msg.sender will be used if this value is zero.
    function init(address _owner) external initializer {
        __Essential_init(_owner);
    }

    /// @notice Adds trusted TDX instances to the registry.
    /// @param _instances The address array of trusted TDX instances.
    /// @return The respective instanceId array per addresses.
    function addInstances(address[] calldata _instances)
        external
        onlyOwner
        returns (uint256[] memory)
    {
        return _addInstances(_instances, true);
    }

    /// @notice Deletes TDX instances from the registry.
    /// @param _ids The ids array of TDX instances.
    function deleteInstances(uint256[] calldata _ids) external onlyOwner {
        uint256 size = _ids.length;
        for (uint256 i; i < size; ++i) {
            uint256 idx = _ids[i];

            require(instances[idx].addr != address(0), TDX_INVALID_INSTANCE());

            emit InstanceDeleted(idx, instances[idx].addr);

            delete instances[idx];
        }
    }

    /// @notice Sets the trusted parameters for quote verification to a specific index
    /// @param index The index of the trusted parameters
    /// @param _params The trusted parameters
    function setTrustedParams(uint256 index, TrustedParams calldata _params) external onlyOwner {
        trustedParams[index] = _params;
        emit TrustedParamsUpdated(index, _params);
    }

    /// @notice Adds an TDX instance after the attestation is verified
    /// @param _trustedParamsIdx The index of the trusted parameters.
    /// @param _attestation The attestation verification parameters.
    /// @return The respective instanceId
    function registerInstance(
        uint256 _trustedParamsIdx,
        VerifyParams memory _attestation
    )
        external
        returns (uint256)
    {
        (bool verified, bytes memory output) = IAutomataDcapAttestation(automataDcapAttestation)
            .verifyAndAttestOnChain(_attestation.quote);
        require(verified, TDX_INVALID_ATTESTATION());

        TrustedParams memory params = trustedParams[_trustedParamsIdx];
        require(params.teeTcbSvn != 0, TDX_INVALID_TRUSTED_PARAMS());

        _validateAttestationOutput(output, _attestation, params);

        bytes32 nonceHash = keccak256(_attestation.nonce);
        require(!nonceUsed[nonceHash], TDX_INVALID_ATTESTATION());
        nonceUsed[nonceHash] = true;

        address[] memory addresses = new address[](1);
        addresses[0] = address(bytes20(_attestation.userData));

        return _addInstances(addresses, false)[0];
    }

    /// @inheritdoc IVerifier
    function verifyProof(
        Context[] calldata _ctxs,
        bytes calldata _proof
    )
        external
        onlyFromEither(taikoInbox, taikoProofVerifier)
    {
        // Size is: 109 bytes
        // 4 bytes + 20 bytes + 20 bytes + 65 bytes (signature) = 109
        require(_proof.length == 109, TDX_INVALID_PROOF());

        address oldInstance = address(bytes20(_proof[4:24]));
        address newInstance = address(bytes20(_proof[24:44]));

        // Collect public inputs
        uint256 size = _ctxs.length;
        bytes32[] memory publicInputs = new bytes32[](size + 2);
        // First public input is the current instance public key
        publicInputs[0] = bytes32(uint256(uint160(oldInstance)));
        publicInputs[1] = bytes32(uint256(uint160(newInstance)));

        // All other inputs are the block program public inputs (a single 32 byte value)
        for (uint256 i; i < size; ++i) {
            // TODO(Yue): For now this assumes the new instance public key to remain the same
            publicInputs[i + 2] = LibPublicInput.hashPublicInputs(
                _ctxs[i].transition, address(this), newInstance, _ctxs[i].metaHash, taikoChainId
            );
        }

        bytes32 signatureHash = keccak256(abi.encodePacked(publicInputs));
        // Verify the blocks
        bytes memory signature = _proof[44:];
        require(oldInstance == ECDSA.recover(signatureHash, signature), TDX_INVALID_PROOF());

        uint32 id = uint32(bytes4(_proof[:4]));
        require(_isInstanceValid(id, oldInstance), TDX_INVALID_INSTANCE());

        if (newInstance != oldInstance && newInstance != address(0)) {
            _replaceInstance(id, oldInstance, newInstance);
        }
    }

    function _addInstances(
        address[] memory _instances,
        bool instantValid
    )
        private
        returns (uint256[] memory ids)
    {
        uint256 size = _instances.length;
        ids = new uint256[](size);

        uint64 validSince = uint64(block.timestamp);

        if (!instantValid) {
            validSince += INSTANCE_VALIDITY_DELAY;
        }

        for (uint256 i; i < size; ++i) {
            require(!addressRegistered[_instances[i]], TDX_ALREADY_ATTESTED());

            addressRegistered[_instances[i]] = true;

            require(_instances[i] != address(0), TDX_INVALID_INSTANCE());

            instances[nextInstanceId] = Instance(_instances[i], validSince);
            ids[i] = nextInstanceId;

            emit InstanceAdded(nextInstanceId, _instances[i], address(0), validSince);

            ++nextInstanceId;
        }
    }

    function _replaceInstance(uint256 id, address oldInstance, address newInstance) private {
        // Replacing an instance means, it went through a cooldown (if added by on-chain RA) so no
        // need to have a cooldown
        instances[id] = Instance(newInstance, uint64(block.timestamp));
        emit InstanceAdded(id, newInstance, oldInstance, block.timestamp);
    }

    function _isInstanceValid(uint256 id, address instance) private view returns (bool) {
        require(instance != address(0), TDX_INVALID_INSTANCE());
        require(instance == instances[id].addr, TDX_INVALID_INSTANCE());
        return instances[id].validSince <= block.timestamp
            && block.timestamp <= instances[id].validSince + INSTANCE_EXPIRY;
    }

    function _validateAttestationOutput(
        bytes memory _attestationOutput,
        VerifyParams memory _attestation,
        TrustedParams memory _params
    )
        private
        pure
    {
        bytes6 teeVersionType = bytes6(BytesUtils.substring(_attestationOutput, 0, 6));

        // TEE Version (0x04) || TEE Type (0x81000000)
        require(teeVersionType == 0x000481000000, TDX_INVALID_VERSION_TYPE());

        bytes16 teeTcbSvn = bytes16(BytesUtils.substring(_attestationOutput, 13, 16));

        require(teeTcbSvn == _params.teeTcbSvn, TDX_INVALID_TCB_SVN());

        bytes memory mrSeam = BytesUtils.substring(_attestationOutput, 29, 48);
        require(mrSeam.length == _params.mrSeam.length, TDX_INVALID_MR_SEAM());
        require(keccak256(mrSeam) == keccak256(_params.mrSeam), TDX_INVALID_MR_SEAM());

        bytes memory mrTd = BytesUtils.substring(_attestationOutput, 149, 48);
        require(mrTd.length == _params.mrTd.length, TDX_INVALID_MR_TD());
        require(keccak256(mrTd) == keccak256(_params.mrTd), TDX_INVALID_MR_TD());

        bytes memory rtMr0 = BytesUtils.substring(_attestationOutput, 341, 48);
        bytes memory rtMr1 = BytesUtils.substring(_attestationOutput, 389, 48);
        bytes memory rtMr2 = BytesUtils.substring(_attestationOutput, 437, 48);
        bytes memory rtMr3 = BytesUtils.substring(_attestationOutput, 485, 48);
        
        require(keccak256(rtMr0) == keccak256(_params.rtMr0), TDX_INVALID_RTMR());
        require(keccak256(rtMr1) == keccak256(_params.rtMr1), TDX_INVALID_RTMR());
        require(keccak256(rtMr2) == keccak256(_params.rtMr2), TDX_INVALID_RTMR());
        require(keccak256(rtMr3) == keccak256(_params.rtMr3), TDX_INVALID_RTMR());
    
        bytes memory reportData = BytesUtils.substring(_attestationOutput, 533, 64);

        bytes32 expectedReportData = sha256(abi.encodePacked(_attestation.userData, _attestation.nonce));
        require(bytes32(reportData) == expectedReportData, TDX_INVALID_REPORT_DATA());
    }
}
