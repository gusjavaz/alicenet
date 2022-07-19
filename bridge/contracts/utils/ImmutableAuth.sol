// This file is auto-generated by hardhat generate-immutable-auth-contract task. DO NOT EDIT.
// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.11;

import "./DeterministicAddress.sol";
import {ImmutableAuthErrorCodes} from "contracts/libraries/errorCodes/ImmutableAuthErrorCodes.sol";

abstract contract ImmutableFactory is DeterministicAddress {
    address private immutable _factory;

    modifier onlyFactory() {
        require(
            msg.sender == _factory,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_FACTORY))
        );
        _;
    }

    constructor(address factory_) {
        _factory = factory_;
    }

    function _factoryAddress() internal view returns (address) {
        return _factory;
    }
}

abstract contract ImmutableAToken is ImmutableFactory {
    address private immutable _aToken;

    modifier onlyAToken() {
        require(
            msg.sender == _aToken,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ATOKEN))
        );
        _;
    }

    constructor() {
        _aToken = getMetamorphicContractAddress(
            0x41546f6b656e0000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenAddress() internal view returns (address) {
        return _aToken;
    }

    function _saltForAToken() internal pure returns (bytes32) {
        return 0x41546f6b656e0000000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableATokenBurner is ImmutableFactory {
    address private immutable _aTokenBurner;

    modifier onlyATokenBurner() {
        require(
            msg.sender == _aTokenBurner,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ATOKENBURNER))
        );
        _;
    }

    constructor() {
        _aTokenBurner = getMetamorphicContractAddress(
            0x41546f6b656e4275726e65720000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenBurnerAddress() internal view returns (address) {
        return _aTokenBurner;
    }

    function _saltForATokenBurner() internal pure returns (bytes32) {
        return 0x41546f6b656e4275726e65720000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableATokenMinter is ImmutableFactory {
    address private immutable _aTokenMinter;

    modifier onlyATokenMinter() {
        require(
            msg.sender == _aTokenMinter,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ATOKENMINTER))
        );
        _;
    }

    constructor() {
        _aTokenMinter = getMetamorphicContractAddress(
            0x41546f6b656e4d696e7465720000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _aTokenMinterAddress() internal view returns (address) {
        return _aTokenMinter;
    }

    function _saltForATokenMinter() internal pure returns (bytes32) {
        return 0x41546f6b656e4d696e7465720000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableBToken is ImmutableFactory {
    address private immutable _bToken;

    modifier onlyBToken() {
        require(
            msg.sender == _bToken,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_BTOKEN))
        );
        _;
    }

    constructor() {
        _bToken = getMetamorphicContractAddress(
            0x42546f6b656e0000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _bTokenAddress() internal view returns (address) {
        return _bToken;
    }

    function _saltForBToken() internal pure returns (bytes32) {
        return 0x42546f6b656e0000000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableBridgePoolDepositNotifier is ImmutableFactory {
    address private immutable _bridgePoolDepositNotifier;

    modifier onlyBridgePoolDepositNotifier() {
        require(
            msg.sender == _bridgePoolDepositNotifier,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_BRIDGEPOOLDEPOSITNOTIFIER
                )
            )
        );
        _;
    }

    constructor() {
        _bridgePoolDepositNotifier = getMetamorphicContractAddress(
            0x427269646765506f6f6c4465706f7369744e6f74696669657200000000000000,
            _factoryAddress()
        );
    }

    function _bridgePoolDepositNotifierAddress() internal view returns (address) {
        return _bridgePoolDepositNotifier;
    }

    function _saltForBridgePoolDepositNotifier() internal pure returns (bytes32) {
        return 0x427269646765506f6f6c4465706f7369744e6f74696669657200000000000000;
    }
}

abstract contract ImmutableBridgeRouter is ImmutableFactory {
    address private immutable _bridgeRouter;

    modifier onlyBridgeRouter() {
        require(
            msg.sender == _bridgeRouter,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_BRIDGEROUTER))
        );
        _;
    }

    constructor() {
        _bridgeRouter = getMetamorphicContractAddress(
            0x427269646765526f757465720000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _bridgeRouterAddress() internal view returns (address) {
        return _bridgeRouter;
    }

    function _saltForBridgeRouter() internal pure returns (bytes32) {
        return 0x427269646765526f757465720000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableFoundation is ImmutableFactory {
    address private immutable _foundation;

    modifier onlyFoundation() {
        require(
            msg.sender == _foundation,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_FOUNDATION))
        );
        _;
    }

    constructor() {
        _foundation = getMetamorphicContractAddress(
            0x466f756e646174696f6e00000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _foundationAddress() internal view returns (address) {
        return _foundation;
    }

    function _saltForFoundation() internal pure returns (bytes32) {
        return 0x466f756e646174696f6e00000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableGovernance is ImmutableFactory {
    address private immutable _governance;

    modifier onlyGovernance() {
        require(
            msg.sender == _governance,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_GOVERNANCE))
        );
        _;
    }

    constructor() {
        _governance = getMetamorphicContractAddress(
            0x476f7665726e616e636500000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _governanceAddress() internal view returns (address) {
        return _governance;
    }

    function _saltForGovernance() internal pure returns (bytes32) {
        return 0x476f7665726e616e636500000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableLiquidityProviderStaking is ImmutableFactory {
    address private immutable _liquidityProviderStaking;

    modifier onlyLiquidityProviderStaking() {
        require(
            msg.sender == _liquidityProviderStaking,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_LIQUIDITYPROVIDERSTAKING
                )
            )
        );
        _;
    }

    constructor() {
        _liquidityProviderStaking = getMetamorphicContractAddress(
            0x4c697175696469747950726f76696465725374616b696e670000000000000000,
            _factoryAddress()
        );
    }

    function _liquidityProviderStakingAddress() internal view returns (address) {
        return _liquidityProviderStaking;
    }

    function _saltForLiquidityProviderStaking() internal pure returns (bytes32) {
        return 0x4c697175696469747950726f76696465725374616b696e670000000000000000;
    }
}

abstract contract ImmutableLocalERC20BridgePoolV1 is ImmutableFactory {
    address private immutable _localERC20BridgePoolV1;

    modifier onlyLocalERC20BridgePoolV1() {
        require(
            msg.sender == _localERC20BridgePoolV1,
            string(
                abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_LOCALERC20BRIDGEPOOLV1)
            )
        );
        _;
    }

    constructor() {
        _localERC20BridgePoolV1 = getMetamorphicContractAddress(
            0x4c6f63616c4552433230427269646765506f6f6c563100000000000000000000,
            _factoryAddress()
        );
    }

    function _localERC20BridgePoolV1Address() internal view returns (address) {
        return _localERC20BridgePoolV1;
    }

    function _saltForLocalERC20BridgePoolV1() internal pure returns (bytes32) {
        return 0x4c6f63616c4552433230427269646765506f6f6c563100000000000000000000;
    }
}

abstract contract ImmutableLocalERC721BridgePoolV1 is ImmutableFactory {
    address private immutable _localERC721BridgePoolV1;

    modifier onlyLocalERC721BridgePoolV1() {
        require(
            msg.sender == _localERC721BridgePoolV1,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_LOCALERC721BRIDGEPOOLV1
                )
            )
        );
        _;
    }

    constructor() {
        _localERC721BridgePoolV1 = getMetamorphicContractAddress(
            0x4c6f63616c455243373231427269646765506f6f6c5631000000000000000000,
            _factoryAddress()
        );
    }

    function _localERC721BridgePoolV1Address() internal view returns (address) {
        return _localERC721BridgePoolV1;
    }

    function _saltForLocalERC721BridgePoolV1() internal pure returns (bytes32) {
        return 0x4c6f63616c455243373231427269646765506f6f6c5631000000000000000000;
    }
}

abstract contract ImmutablePublicStaking is ImmutableFactory {
    address private immutable _publicStaking;

    modifier onlyPublicStaking() {
        require(
            msg.sender == _publicStaking,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_PUBLICSTAKING))
        );
        _;
    }

    constructor() {
        _publicStaking = getMetamorphicContractAddress(
            0x5075626c69635374616b696e6700000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _publicStakingAddress() internal view returns (address) {
        return _publicStaking;
    }

    function _saltForPublicStaking() internal pure returns (bytes32) {
        return 0x5075626c69635374616b696e6700000000000000000000000000000000000000;
    }
}

abstract contract ImmutableSnapshots is ImmutableFactory {
    address private immutable _snapshots;

    modifier onlySnapshots() {
        require(
            msg.sender == _snapshots,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_SNAPSHOTS))
        );
        _;
    }

    constructor() {
        _snapshots = getMetamorphicContractAddress(
            0x536e617073686f74730000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _snapshotsAddress() internal view returns (address) {
        return _snapshots;
    }

    function _saltForSnapshots() internal pure returns (bytes32) {
        return 0x536e617073686f74730000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableStakingPositionDescriptor is ImmutableFactory {
    address private immutable _stakingPositionDescriptor;

    modifier onlyStakingPositionDescriptor() {
        require(
            msg.sender == _stakingPositionDescriptor,
            string(
                abi.encodePacked(
                    ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_STAKINGPOSITIONDESCRIPTOR
                )
            )
        );
        _;
    }

    constructor() {
        _stakingPositionDescriptor = getMetamorphicContractAddress(
            0x5374616b696e67506f736974696f6e44657363726970746f7200000000000000,
            _factoryAddress()
        );
    }

    function _stakingPositionDescriptorAddress() internal view returns (address) {
        return _stakingPositionDescriptor;
    }

    function _saltForStakingPositionDescriptor() internal pure returns (bytes32) {
        return 0x5374616b696e67506f736974696f6e44657363726970746f7200000000000000;
    }
}

abstract contract ImmutableValidatorPool is ImmutableFactory {
    address private immutable _validatorPool;

    modifier onlyValidatorPool() {
        require(
            msg.sender == _validatorPool,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_VALIDATORPOOL))
        );
        _;
    }

    constructor() {
        _validatorPool = getMetamorphicContractAddress(
            0x56616c696461746f72506f6f6c00000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _validatorPoolAddress() internal view returns (address) {
        return _validatorPool;
    }

    function _saltForValidatorPool() internal pure returns (bytes32) {
        return 0x56616c696461746f72506f6f6c00000000000000000000000000000000000000;
    }
}

abstract contract ImmutableValidatorStaking is ImmutableFactory {
    address private immutable _validatorStaking;

    modifier onlyValidatorStaking() {
        require(
            msg.sender == _validatorStaking,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_VALIDATORSTAKING))
        );
        _;
    }

    constructor() {
        _validatorStaking = getMetamorphicContractAddress(
            0x56616c696461746f725374616b696e6700000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _validatorStakingAddress() internal view returns (address) {
        return _validatorStaking;
    }

    function _saltForValidatorStaking() internal pure returns (bytes32) {
        return 0x56616c696461746f725374616b696e6700000000000000000000000000000000;
    }
}

abstract contract ImmutableCallAny is ImmutableFactory {
    address private immutable _callAny;

    modifier onlyCallAny() {
        require(
            msg.sender == _callAny,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_CALLANY))
        );
        _;
    }

    constructor() {
        _callAny = getMetamorphicContractAddress(
            0x43616c6c416e7900000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _callAnyAddress() internal view returns (address) {
        return _callAny;
    }

    function _saltForCallAny() internal pure returns (bytes32) {
        return 0x43616c6c416e7900000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKGAccusations is ImmutableFactory {
    address private immutable _ethdkgAccusations;

    modifier onlyETHDKGAccusations() {
        require(
            msg.sender == _ethdkgAccusations,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ETHDKGACCUSATIONS))
        );
        _;
    }

    constructor() {
        _ethdkgAccusations = getMetamorphicContractAddress(
            0x455448444b4741636375736174696f6e73000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgAccusationsAddress() internal view returns (address) {
        return _ethdkgAccusations;
    }

    function _saltForETHDKGAccusations() internal pure returns (bytes32) {
        return 0x455448444b4741636375736174696f6e73000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKGPhases is ImmutableFactory {
    address private immutable _ethdkgPhases;

    modifier onlyETHDKGPhases() {
        require(
            msg.sender == _ethdkgPhases,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ETHDKGPHASES))
        );
        _;
    }

    constructor() {
        _ethdkgPhases = getMetamorphicContractAddress(
            0x455448444b475068617365730000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgPhasesAddress() internal view returns (address) {
        return _ethdkgPhases;
    }

    function _saltForETHDKGPhases() internal pure returns (bytes32) {
        return 0x455448444b475068617365730000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableETHDKG is ImmutableFactory {
    address private immutable _ethdkg;

    modifier onlyETHDKG() {
        require(
            msg.sender == _ethdkg,
            string(abi.encodePacked(ImmutableAuthErrorCodes.IMMUTEABLEAUTH_ONLY_ETHDKG))
        );
        _;
    }

    constructor() {
        _ethdkg = getMetamorphicContractAddress(
            0x455448444b470000000000000000000000000000000000000000000000000000,
            _factoryAddress()
        );
    }

    function _ethdkgAddress() internal view returns (address) {
        return _ethdkg;
    }

    function _saltForETHDKG() internal pure returns (bytes32) {
        return 0x455448444b470000000000000000000000000000000000000000000000000000;
    }
}
