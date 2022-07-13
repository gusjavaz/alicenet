// Generated by ifacemaker. DO NOT EDIT.

package bindings

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// IBTokenCaller ...
type IBTokenCaller interface {
	// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
	//
	// Solidity: function allowance(address owner, address spender) view returns(uint256)
	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)
	// BTokensToEth is a free data retrieval call binding the contract method 0xc7636071.
	//
	// Solidity: function bTokensToEth(uint256 poolBalance_, uint256 totalSupply_, uint256 numBTK_) pure returns(uint256 numEth)
	BTokensToEth(opts *bind.CallOpts, poolBalance_ *big.Int, totalSupply_ *big.Int, numBTK_ *big.Int) (*big.Int, error)
	// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
	//
	// Solidity: function balanceOf(address account) view returns(uint256)
	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)
	// Decimals is a free data retrieval call binding the contract method 0x313ce567.
	//
	// Solidity: function decimals() view returns(uint8)
	Decimals(opts *bind.CallOpts) (uint8, error)
	// EthToBTokens is a free data retrieval call binding the contract method 0x75a86a7e.
	//
	// Solidity: function ethToBTokens(uint256 poolBalance_, uint256 numEth_) pure returns(uint256)
	EthToBTokens(opts *bind.CallOpts, poolBalance_ *big.Int, numEth_ *big.Int) (*big.Int, error)
	// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
	//
	// Solidity: function getAdmin() view returns(address)
	GetAdmin(opts *bind.CallOpts) (common.Address, error)
	// GetDeposit is a free data retrieval call binding the contract method 0x9f9fb968.
	//
	// Solidity: function getDeposit(uint256 depositID) view returns((uint8,address,uint256))
	GetDeposit(opts *bind.CallOpts, depositID *big.Int) (BTokenDeposit, error)
	// GetMetamorphicContractAddress is a free data retrieval call binding the contract method 0x8653a465.
	//
	// Solidity: function getMetamorphicContractAddress(bytes32 _salt, address _factory) pure returns(address)
	GetMetamorphicContractAddress(opts *bind.CallOpts, _salt [32]byte, _factory common.Address) (common.Address, error)
	// GetPoolBalance is a free data retrieval call binding the contract method 0xabd70aa2.
	//
	// Solidity: function getPoolBalance() view returns(uint256)
	GetPoolBalance(opts *bind.CallOpts) (*big.Int, error)
	// GetStaticPoolContractAddress is a free data retrieval call binding the contract method 0x0ffd7a81.
	//
	// Solidity: function getStaticPoolContractAddress(bytes32 _salt, address _bridgeRouter) pure returns(address)
	GetStaticPoolContractAddress(opts *bind.CallOpts, _salt [32]byte, _bridgeRouter common.Address) (common.Address, error)
	// GetTotalBTokensDeposited is a free data retrieval call binding the contract method 0x5ecef3af.
	//
	// Solidity: function getTotalBTokensDeposited() view returns(uint256)
	GetTotalBTokensDeposited(opts *bind.CallOpts) (*big.Int, error)
	// Name is a free data retrieval call binding the contract method 0x06fdde03.
	//
	// Solidity: function name() view returns(string)
	Name(opts *bind.CallOpts) (string, error)
	// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
	//
	// Solidity: function symbol() view returns(string)
	Symbol(opts *bind.CallOpts) (string, error)
	// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
	//
	// Solidity: function totalSupply() view returns(uint256)
	TotalSupply(opts *bind.CallOpts) (*big.Int, error)
}
