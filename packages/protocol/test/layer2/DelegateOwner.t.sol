// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "../shared/thirdparty/Multicall3.sol";
import "./TaikoL2Test.sol";

contract Target is EssentialContract {
    function init(address _owner) external initializer {
        __Essential_init(_owner);
    }
}

contract TestDelegateOwner is TaikoL2Test {
    address public owner;
    address public remoteOwner;
    Bridge public bridge;
    SignalService public signalService;
    AddressManager public addressManager;
    DelegateOwner public delegateOwner;
    Multicall3 public multicall;

    uint64 remoteChainId = uint64(block.chainid + 1);
    address remoteBridge = vm.addr(0x2000);

    function setUp() public {
        owner = vm.addr(0x1000);
        vm.deal(owner, 100 ether);

        remoteOwner = vm.addr(0x2000);

        vm.startPrank(owner);

        multicall = new Multicall3();

        addressManager = AddressManager(
            deployProxy({
                name: "address_manager",
                impl: address(new AddressManager()),
                data: abi.encodeCall(AddressManager.init, (address(0)))
            })
        );

        delegateOwner = DelegateOwner(
            deployProxy({
                name: "delegate_owner",
                impl: address(new DelegateOwner()),
                data: abi.encodeCall(
                    DelegateOwner.init,
                    (remoteOwner, address(addressManager), remoteChainId, address(0))
                ),
                registerTo: address(addressManager)
            })
        );

        signalService = SkipProofCheckSignal(
            deployProxy({
                name: "signal_service",
                impl: address(new SkipProofCheckSignal()),
                data: abi.encodeCall(SignalService.init, (address(0), address(addressManager))),
                registerTo: address(addressManager)
            })
        );

        bridge = Bridge(
            payable(
                deployProxy({
                    name: "bridge",
                    impl: address(new Bridge()),
                    data: abi.encodeCall(Bridge.init, (address(0), address(addressManager))),
                    registerTo: address(addressManager)
                })
            )
        );

        addressManager.setAddress(remoteChainId, "bridge", remoteBridge);
        vm.stopPrank();
    }

    function test_delegate_owner_single_non_delegatecall_self() public {
        address delegateOwnerImpl2 = address(new DelegateOwner());

        bytes memory data = abi.encode(
            DelegateOwner.Call(
                uint64(0),
                address(delegateOwner),
                false, // CALL
                abi.encodeCall(UUPSUpgradeable.upgradeTo, (delegateOwnerImpl2))
            )
        );

        vm.expectRevert(DelegateOwner.DO_DRYRUN_SUCCEEDED.selector);
        delegateOwner.dryrunInvocation(data);

        IBridge.Message memory message;
        message.from = remoteOwner;
        message.destChainId = uint64(block.chainid);
        message.srcChainId = remoteChainId;
        message.destOwner = Bob;
        message.data = abi.encodeCall(DelegateOwner.onMessageInvocation, (data));
        message.to = address(delegateOwner);

        vm.prank(Bob);
        bridge.processMessage(message, "");

        bytes32 hash = bridge.hashMessage(message);
        assertTrue(bridge.messageStatus(hash) == IBridge.Status.DONE);

        assertEq(delegateOwner.nextTxId(), 1);
        assertEq(delegateOwner.impl(), delegateOwnerImpl2);
    }
}
