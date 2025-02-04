// Generated by ifacemaker. DO NOT EDIT.

package bindings

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// IETHDKGCaller ...
type IETHDKGCaller interface {
	// GetBadParticipants is a free data retrieval call binding the contract method 0x32d4d570.
	//
	// Solidity: function getBadParticipants() view returns(uint256)
	GetBadParticipants(opts *bind.CallOpts) (*big.Int, error)
	// GetConfirmationLength is a free data retrieval call binding the contract method 0x8c848d32.
	//
	// Solidity: function getConfirmationLength() view returns(uint256)
	GetConfirmationLength(opts *bind.CallOpts) (*big.Int, error)
	// GetETHDKGPhase is a free data retrieval call binding the contract method 0x2958e81c.
	//
	// Solidity: function getETHDKGPhase() view returns(uint8)
	GetETHDKGPhase(opts *bind.CallOpts) (uint8, error)
	// GetLastRoundParticipantIndex is a free data retrieval call binding the contract method 0x775694e1.
	//
	// Solidity: function getLastRoundParticipantIndex(address participant) view returns(uint256)
	GetLastRoundParticipantIndex(opts *bind.CallOpts, participant common.Address) (*big.Int, error)
	// GetMasterPublicKey is a free data retrieval call binding the contract method 0xe146372a.
	//
	// Solidity: function getMasterPublicKey() view returns(uint256[4])
	GetMasterPublicKey(opts *bind.CallOpts) ([4]*big.Int, error)
	// GetMasterPublicKeyHash is a free data retrieval call binding the contract method 0x1c67d576.
	//
	// Solidity: function getMasterPublicKeyHash() view returns(bytes32)
	GetMasterPublicKeyHash(opts *bind.CallOpts) ([32]byte, error)
	// GetMetamorphicContractAddress is a free data retrieval call binding the contract method 0x8653a465.
	//
	// Solidity: function getMetamorphicContractAddress(bytes32 _salt, address _factory) pure returns(address)
	GetMetamorphicContractAddress(opts *bind.CallOpts, _salt [32]byte, _factory common.Address) (common.Address, error)
	// GetMinValidators is a free data retrieval call binding the contract method 0xecbadb36.
	//
	// Solidity: function getMinValidators() pure returns(uint256)
	GetMinValidators(opts *bind.CallOpts) (*big.Int, error)
	// GetNonce is a free data retrieval call binding the contract method 0xd087d288.
	//
	// Solidity: function getNonce() view returns(uint256)
	GetNonce(opts *bind.CallOpts) (*big.Int, error)
	// GetNumParticipants is a free data retrieval call binding the contract method 0xfd478ca9.
	//
	// Solidity: function getNumParticipants() view returns(uint256)
	GetNumParticipants(opts *bind.CallOpts) (*big.Int, error)
	// GetParticipantInternalState is a free data retrieval call binding the contract method 0xbf7786b6.
	//
	// Solidity: function getParticipantInternalState(address participant) view returns((uint256[2],uint64,uint64,uint8,bytes32,uint256[2],uint256[2],uint256[4]))
	GetParticipantInternalState(opts *bind.CallOpts, participant common.Address) (Participant, error)
	// GetParticipantsInternalState is a free data retrieval call binding the contract method 0xc016baee.
	//
	// Solidity: function getParticipantsInternalState(address[] participantAddresses) view returns((uint256[2],uint64,uint64,uint8,bytes32,uint256[2],uint256[2],uint256[4])[])
	GetParticipantsInternalState(opts *bind.CallOpts, participantAddresses []common.Address) ([]Participant, error)
	// GetPhaseLength is a free data retrieval call binding the contract method 0x106da57d.
	//
	// Solidity: function getPhaseLength() view returns(uint256)
	GetPhaseLength(opts *bind.CallOpts) (*big.Int, error)
	// GetPhaseStartBlock is a free data retrieval call binding the contract method 0xa2bc9c78.
	//
	// Solidity: function getPhaseStartBlock() view returns(uint256)
	GetPhaseStartBlock(opts *bind.CallOpts) (*big.Int, error)
	// IsETHDKGCompleted is a free data retrieval call binding the contract method 0x2b7c6724.
	//
	// Solidity: function isETHDKGCompleted() view returns(bool)
	IsETHDKGCompleted(opts *bind.CallOpts) (bool, error)
	// IsETHDKGHalted is a free data retrieval call binding the contract method 0x43ced534.
	//
	// Solidity: function isETHDKGHalted() view returns(bool)
	IsETHDKGHalted(opts *bind.CallOpts) (bool, error)
	// IsETHDKGRunning is a free data retrieval call binding the contract method 0x747b217c.
	//
	// Solidity: function isETHDKGRunning() view returns(bool)
	IsETHDKGRunning(opts *bind.CallOpts) (bool, error)
	// IsMasterPublicKeySet is a free data retrieval call binding the contract method 0x08efcf16.
	//
	// Solidity: function isMasterPublicKeySet() view returns(bool)
	IsMasterPublicKeySet(opts *bind.CallOpts) (bool, error)
}
