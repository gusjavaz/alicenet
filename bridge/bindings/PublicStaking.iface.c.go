// Generated by ifacemaker. DO NOT EDIT.

package bindings

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// IPublicStakingCaller ...
type IPublicStakingCaller interface {
	// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
	//
	// Solidity: function balanceOf(address owner) view returns(uint256)
	BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error)
	// CircuitBreakerState is a free data retrieval call binding the contract method 0x89465c62.
	//
	// Solidity: function circuitBreakerState() view returns(bool)
	CircuitBreakerState(opts *bind.CallOpts) (bool, error)
	// EstimateEthCollection is a free data retrieval call binding the contract method 0x20ea0d48.
	//
	// Solidity: function estimateEthCollection(uint256 tokenID_) view returns(uint256 payout)
	EstimateEthCollection(opts *bind.CallOpts, tokenID_ *big.Int) (*big.Int, error)
	// EstimateExcessEth is a free data retrieval call binding the contract method 0x905953ed.
	//
	// Solidity: function estimateExcessEth() view returns(uint256 excess)
	EstimateExcessEth(opts *bind.CallOpts) (*big.Int, error)
	// EstimateExcessToken is a free data retrieval call binding the contract method 0x3eed3eff.
	//
	// Solidity: function estimateExcessToken() view returns(uint256 excess)
	EstimateExcessToken(opts *bind.CallOpts) (*big.Int, error)
	// EstimateTokenCollection is a free data retrieval call binding the contract method 0x93c5748d.
	//
	// Solidity: function estimateTokenCollection(uint256 tokenID_) view returns(uint256 payout)
	EstimateTokenCollection(opts *bind.CallOpts, tokenID_ *big.Int) (*big.Int, error)
	// GetAccumulatorScaleFactor is a free data retrieval call binding the contract method 0x99785132.
	//
	// Solidity: function getAccumulatorScaleFactor() pure returns(uint256)
	GetAccumulatorScaleFactor(opts *bind.CallOpts) (*big.Int, error)
	// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
	//
	// Solidity: function getApproved(uint256 tokenId) view returns(address)
	GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error)
	// GetEthAccumulator is a free data retrieval call binding the contract method 0x548652d2.
	//
	// Solidity: function getEthAccumulator() view returns(uint256 accumulator, uint256 slush)
	GetEthAccumulator(opts *bind.CallOpts) (struct {
		Accumulator *big.Int
		Slush       *big.Int
	}, error)
	// GetLatestMintedPositionID is a free data retrieval call binding the contract method 0x79bb0985.
	//
	// Solidity: function getLatestMintedPositionID() view returns(uint256)
	GetLatestMintedPositionID(opts *bind.CallOpts) (*big.Int, error)
	// GetMaxGovernanceLock is a free data retrieval call binding the contract method 0xf44d258b.
	//
	// Solidity: function getMaxGovernanceLock() pure returns(uint256)
	GetMaxGovernanceLock(opts *bind.CallOpts) (*big.Int, error)
	// GetMaxMintLock is a free data retrieval call binding the contract method 0x090f70f0.
	//
	// Solidity: function getMaxMintLock() pure returns(uint256)
	GetMaxMintLock(opts *bind.CallOpts) (*big.Int, error)
	// GetMetamorphicContractAddress is a free data retrieval call binding the contract method 0x8653a465.
	//
	// Solidity: function getMetamorphicContractAddress(bytes32 _salt, address _factory) pure returns(address)
	GetMetamorphicContractAddress(opts *bind.CallOpts, _salt [32]byte, _factory common.Address) (common.Address, error)
	// GetPosition is a free data retrieval call binding the contract method 0xeb02c301.
	//
	// Solidity: function getPosition(uint256 tokenID_) view returns(uint256 shares, uint256 freeAfter, uint256 withdrawFreeAfter, uint256 accumulatorEth, uint256 accumulatorToken)
	GetPosition(opts *bind.CallOpts, tokenID_ *big.Int) (struct {
		Shares            *big.Int
		FreeAfter         *big.Int
		WithdrawFreeAfter *big.Int
		AccumulatorEth    *big.Int
		AccumulatorToken  *big.Int
	}, error)
	// GetStaticPoolContractAddress is a free data retrieval call binding the contract method 0x0ffd7a81.
	//
	// Solidity: function getStaticPoolContractAddress(bytes32 _salt, address _bridgeRouter) pure returns(address)
	GetStaticPoolContractAddress(opts *bind.CallOpts, _salt [32]byte, _bridgeRouter common.Address) (common.Address, error)
	// GetTokenAccumulator is a free data retrieval call binding the contract method 0xc47c6e14.
	//
	// Solidity: function getTokenAccumulator() view returns(uint256 accumulator, uint256 slush)
	GetTokenAccumulator(opts *bind.CallOpts) (struct {
		Accumulator *big.Int
		Slush       *big.Int
	}, error)
	// GetTotalReserveAToken is a free data retrieval call binding the contract method 0x856de8d2.
	//
	// Solidity: function getTotalReserveAToken() view returns(uint256)
	GetTotalReserveAToken(opts *bind.CallOpts) (*big.Int, error)
	// GetTotalReserveEth is a free data retrieval call binding the contract method 0x19b8be2f.
	//
	// Solidity: function getTotalReserveEth() view returns(uint256)
	GetTotalReserveEth(opts *bind.CallOpts) (*big.Int, error)
	// GetTotalShares is a free data retrieval call binding the contract method 0xd5002f2e.
	//
	// Solidity: function getTotalShares() view returns(uint256)
	GetTotalShares(opts *bind.CallOpts) (*big.Int, error)
	// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
	//
	// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
	IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error)
	// Name is a free data retrieval call binding the contract method 0x06fdde03.
	//
	// Solidity: function name() view returns(string)
	Name(opts *bind.CallOpts) (string, error)
	// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
	//
	// Solidity: function ownerOf(uint256 tokenId) view returns(address)
	OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error)
	// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
	//
	// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)
	// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
	//
	// Solidity: function symbol() view returns(string)
	Symbol(opts *bind.CallOpts) (string, error)
	// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
	//
	// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
	TokenByIndex(opts *bind.CallOpts, index *big.Int) (*big.Int, error)
	// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
	//
	// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
	TokenOfOwnerByIndex(opts *bind.CallOpts, owner common.Address, index *big.Int) (*big.Int, error)
	// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
	//
	// Solidity: function tokenURI(uint256 tokenID_) view returns(string)
	TokenURI(opts *bind.CallOpts, tokenID_ *big.Int) (string, error)
	// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
	//
	// Solidity: function totalSupply() view returns(uint256)
	TotalSupply(opts *bind.CallOpts) (*big.Int, error)
}
