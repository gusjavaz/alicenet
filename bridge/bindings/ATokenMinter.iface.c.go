// Generated by ifacemaker. DO NOT EDIT.

package bindings

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// IATokenMinterCaller ...
type IATokenMinterCaller interface {
	// GetMetamorphicContractAddress is a free data retrieval call binding the contract method 0x8653a465.
	//
	// Solidity: function getMetamorphicContractAddress(bytes32 _salt, address _factory) pure returns(address)
	GetMetamorphicContractAddress(opts *bind.CallOpts, _salt [32]byte, _factory common.Address) (common.Address, error)
	// GetStaticPoolContractAddress is a free data retrieval call binding the contract method 0x0ffd7a81.
	//
	// Solidity: function getStaticPoolContractAddress(bytes32 _salt, address _bridgeRouter) pure returns(address)
	GetStaticPoolContractAddress(opts *bind.CallOpts, _salt [32]byte, _bridgeRouter common.Address) (common.Address, error)
}
