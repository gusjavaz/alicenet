// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.16;

library BridgePoolFactoryErrors {
    error FailedToDeployLogic();
    error PoolLogicNotSupported();
    error LogicVersionDoesNotExist(uint8 poolType, uint8 tokenType);
    error StaticPoolDeploymentFailed(bytes32 salt_);
    error UnexistentBridgePoolImplementationVersion(uint16 version);
    error UnableToDeployBridgePool(bytes32 salt_);
    error InsufficientFunds();
    error PublicPoolDeploymentTemporallyDisabled();
}
