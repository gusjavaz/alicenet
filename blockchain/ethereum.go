package blockchain

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/bridge/bindings"
	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/golang-collections/go-datastructures/queue"
	"github.com/sirupsen/logrus"
)

//Ethereum contains state information about a connection to Ethereum
type Ethereum interface {
	IsEthereumAccessible() bool

	GetCallOpts(context.Context, accounts.Account) *bind.CallOpts
	GetTransactionOpts(context.Context, accounts.Account) (*bind.TransactOpts, error)

	LoadAccounts(string)
	LoadPasscodes(string) error

	UnlockAccount(accounts.Account) error

	TransferEther(common.Address, common.Address, *big.Int) error

	GetAccount(common.Address) (accounts.Account, error)
	GetAccountKeys(addr common.Address) (*keystore.Key, error)
	GetBalance(common.Address) (*big.Int, error)
	GetGethClient() GethClient
	GetCoinbaseAddress() common.Address
	GetCurrentHeight(context.Context) (uint64, error)
	GetDefaultAccount() accounts.Account
	GetEndpoint() string
	GetEvents(ctx context.Context, firstBlock uint64, lastBlock uint64, addresses []common.Address) ([]types.Log, error)
	GetFinalizedHeight(context.Context) (uint64, error)
	GetPeerCount(context.Context) (uint64, error)
	GetSnapshot() ([]byte, error)
	GetSyncProgress() (bool, *geth.SyncProgress, error)
	GetTimeoutContext() (context.Context, context.CancelFunc)
	GetValidators() ([]common.Address, error)

	WaitForReceipt(context.Context, *types.Transaction) (*types.Receipt, error)

	RetryCount() int
	RetryDelay() time.Duration

	Timeout() time.Duration

	Contracts() *Contracts
}

// Ethereum specific errors
var (
	ErrAccountNotFound  = errors.New("could not find specified account")
	ErrKeysNotFound     = errors.New("account either not found or not unlocked")
	ErrPasscodeNotFound = errors.New("could not find specified passcode")
)

// GethClient is an amalgamation of the geth interfaces used
type GethClient interface {

	// geth.ChainReader
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (geth.Subscription, error)

	// geth.TransactionReader
	TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)

	// geth.ChainStateReader
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)

	// bind.ContractBackend
	// -- bind.ContractCaller
	// CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call geth.CallMsg, blockNumber *big.Int) ([]byte, error)

	// -- bind.ContractTransactor
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, call geth.CallMsg) (gas uint64, err error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error

	// -- bind.ContractFilterer
	FilterLogs(ctx context.Context, query geth.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, query geth.FilterQuery, ch chan<- types.Log) (geth.Subscription, error)
}

type ethereum struct {
	logger         *logrus.Logger
	endpoint       string
	keystore       *keystore.KeyStore
	finalityDelay  uint64
	accounts       map[common.Address]accounts.Account
	coinbase       common.Address
	defaultAccount accounts.Account
	keys           map[common.Address]*keystore.Key
	passcodes      map[common.Address]string
	timeout        time.Duration
	retryCount     int
	retryDelay     time.Duration
	contracts      *Contracts
	client         GethClient
	commit         func()
	chainID        *big.Int
	syncing        func(ctx context.Context) (*geth.SyncProgress, error)
	peerCount      func(ctx context.Context) (uint64, error)
}

// Contracts contains bindings to smart contract system
type Contracts struct {
	eth                 *ethereum
	Crypto              *bindings.Crypto
	CryptoAddress       common.Address
	Deposit             *bindings.Deposit
	DepositAddress      common.Address
	Ethdkg              *bindings.ETHDKG
	EthdkgAddress       common.Address
	Participants        *bindings.ParticipantsFacet
	Registry            *bindings.Registry
	RegistryAddress     common.Address
	Snapshots           *bindings.SnapshotsFacet
	Staking             *bindings.Staking
	StakingAddress      common.Address
	StakingToken        *bindings.Token
	StakingTokenAddress common.Address
	StakingValues       *bindings.StakingValuesFacet
	UtilityToken        *bindings.Token
	UtilityTokenAddress common.Address
	Validators          *bindings.Validators
	ValidatorsAddress   common.Address
}

//NewEthereumSimulator returns a simulator for testing
func NewEthereumSimulator(
	pathKeystore string,
	pathPasscodes string,
	retryCount int,
	retryDelay time.Duration,
	finalityDelay int,
	wei *big.Int,
	addresses ...string) (Ethereum, func(), error) {
	logger := logging.GetLogger("ethsim")

	if len(addresses) < 1 {
		return nil, nil, errors.New("at least 1 account address required")
	}

	defaultAccount := addresses[0]

	genAlloc := make(core.GenesisAlloc)
	for _, address := range addresses {
		addr := common.HexToAddress(address)
		genAlloc[addr] = core.GenesisAccount{Balance: wei}
	}

	eth := &ethereum{
		logger:        logger,
		accounts:      make(map[common.Address]accounts.Account),
		keys:          make(map[common.Address]*keystore.Key),
		passcodes:     make(map[common.Address]string),
		retryCount:    retryCount,
		retryDelay:    retryDelay,
		finalityDelay: uint64(finalityDelay)}
	eth.contracts = &Contracts{eth: eth}

	eth.LoadAccounts(pathKeystore)
	err := eth.LoadPasscodes(pathPasscodes)
	if err != nil {
		logger.Errorf("Error in NewEthereumSimulator at eth.LoadPasscodes: %v", err)
		return nil, nil, err
	}

	eth.defaultAccount, err = eth.GetAccount(common.HexToAddress(defaultAccount))
	if err != nil {
		logger.Errorf("Can't find user to set as default %v: %v", defaultAccount, err)
		return nil, nil, err
	}

	gasLimit := uint64(10000000000000000)
	sim := backends.NewSimulatedBackend(genAlloc, gasLimit)
	eth.client = sim
	eth.chainID = big.NewInt(1337)
	eth.peerCount = func(context.Context) (uint64, error) {
		return 0, nil
	}
	eth.syncing = func(ctx context.Context) (*geth.SyncProgress, error) {
		return nil, nil
	}

	eth.commit = func() {
		sim.Commit()
	}

	return eth, eth.commit, nil
}

//NewEthereum creates a new Ethereum
func NewEthereumEndpoint(
	endpoint string,
	pathKeystore string,
	pathPasscodes string,
	defaultAccount string,
	timeout time.Duration,
	retryCount int,
	retryDelay time.Duration,
	finalityDelay int) (Ethereum, error) {

	logger := logging.GetLogger("ethereum")

	eth := &ethereum{
		endpoint:      endpoint,
		logger:        logger,
		accounts:      make(map[common.Address]accounts.Account),
		keys:          make(map[common.Address]*keystore.Key),
		passcodes:     make(map[common.Address]string),
		finalityDelay: uint64(finalityDelay),
		timeout:       timeout,
		retryCount:    retryCount,
		retryDelay:    retryDelay}

	eth.contracts = &Contracts{eth: eth}

	// Load accounts + passcodes
	eth.LoadAccounts(pathKeystore)
	err := eth.LoadPasscodes(pathPasscodes)
	if err != nil {
		logger.Errorf("Error in NewEthereumEndpoint at eth.LoadPasscodes: %v", err)
		return nil, err
	}

	// Designate accounts
	var acct accounts.Account
	acct, err = eth.GetAccount(common.HexToAddress(defaultAccount))
	if err != nil {
		logger.Errorf("Can't find user to set as default %v: %v", defaultAccount, err)
		return nil, err
	}
	eth.SetDefaultAccount(acct)

	// Low level rpc client
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rpcClient, rpcErr := rpc.DialContext(ctx, endpoint)
	if rpcErr != nil {
		logger.Errorf("Error in NewEthereumEndpoint at rpc.DialContext: %v", err)
		return nil, rpcErr
	}
	ethClient := ethclient.NewClient(rpcClient)
	eth.client = ethClient
	eth.chainID, err = ethClient.ChainID(ctx)
	if err != nil {
		logger.Errorf("Error in NewEthereumEndpoint at ethClient.ChainID: %v", err)
		return nil, err
	}

	eth.peerCount = func(ctx context.Context) (uint64, error) {
		return eth.getPeerCount(ctx, rpcClient)
	}
	eth.syncing = ethClient.SyncProgress

	// Find coinbase
	if e := rpcClient.CallContext(ctx, &eth.coinbase, "eth_coinbase"); e != nil {
		logger.Warnf("Failed to determine coinbase: %v", e)
	} else {
		logger.Infof("Coinbase: %v", eth.coinbase.Hex())
	}

	logger.Debug("Completed initialization")
	eth.commit = func() {}

	return eth, nil
}

func (eth *ethereum) Contracts() *Contracts {
	return eth.contracts
}

func (eth *ethereum) GetPeerCount(ctx context.Context) (uint64, error) {
	return eth.peerCount(ctx)
}

func (eth *ethereum) getPeerCount(ctx context.Context, rpcClient *rpc.Client) (uint64, error) {
	// Let's see how many peers our endpoint has
	var peerCountString string
	if err := rpcClient.CallContext(ctx, &peerCountString, "net_peerCount"); err != nil {
		eth.logger.Warnf("could not get peerCount: %v", err)
		return 0, err
	}

	var peerCount uint64
	_, err := fmt.Sscanf(peerCountString, "0x%x", &peerCount)
	if err != nil {
		eth.logger.Warnf("could not parse peerCount: %v", err)
		return 0, err
	}
	return peerCount, nil
}

//IsEthereumAccessible checks against endpoint to confirm server responds
func (eth *ethereum) IsEthereumAccessible() bool {
	ctx, cancel := eth.GetTimeoutContext()
	defer cancel()
	block, err := eth.client.BlockByNumber(ctx, nil)
	if err == nil && block != nil {
		return true
	}

	eth.logger.Debug("IsEthereumAccessible()...false")
	return false
}

// Scans the directory specified and loads all the accounts found
func (eth *ethereum) LoadAccounts(directoryPath string) {
	logger := eth.logger

	logger.Infof("LoadAccounts(\"%v\")...", directoryPath)
	ks := keystore.NewKeyStore(directoryPath, keystore.StandardScryptN, keystore.StandardScryptP)
	accts := make(map[common.Address]accounts.Account, 10)

	for _, wallet := range ks.Wallets() {
		for _, account := range wallet.Accounts() {
			logger.Infof("... found account %v", account.Address.Hex())
			accts[account.Address] = account
		}
	}

	eth.accounts = accts
	eth.keystore = ks
}

// LoadPasscodes loads the specified passcode file
func (eth *ethereum) LoadPasscodes(filePath string) error {
	logger := eth.logger

	logger.Infof("LoadPasscodes(\"%v\")...", filePath)
	passcodes := make(map[common.Address]string)

	file, err := os.Open(filePath)
	if err != nil {
		logger.Errorf("Failed to open passcode file \"%v\": %s", filePath, err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") {
			components := strings.Split(line, "=")
			if len(components) == 2 {
				address := strings.TrimSpace(components[0])
				passcode := strings.TrimSpace(components[1])

				passcodes[common.HexToAddress(address)] = passcode
			}
		}
	}

	eth.passcodes = passcodes

	return nil
}

// UnlockAccount unlocks the previously loaded account using the previously loaded passcode
func (eth *ethereum) UnlockAccount(acct accounts.Account) error {

	passcode, passcodeFound := eth.passcodes[acct.Address]
	if !passcodeFound {
		return ErrPasscodeNotFound
	}

	err := eth.keystore.Unlock(acct, passcode)
	if err != nil {
		return err
	}

	// Open the account key file
	keyJson, err := ioutil.ReadFile(acct.URL.Path)
	if err != nil {
		return err
	}

	// Get the private key
	key, err := keystore.DecryptKey(keyJson, passcode)
	if err != nil {
		return err
	}

	eth.keys[acct.Address] = key

	return nil
}

// GetGethClient returns an amalgamated geth client interface
func (eth *ethereum) GetGethClient() GethClient {
	return eth.client
}

// GetAccount returns the account specified
func (eth *ethereum) GetAccount(addr common.Address) (accounts.Account, error) {
	acct, accountFound := eth.accounts[addr]
	if !accountFound {
		return acct, ErrAccountNotFound
	}

	return acct, nil
}

func (eth *ethereum) GetAccountKeys(addr common.Address) (*keystore.Key, error) {
	if key, ok := eth.keys[addr]; ok {
		return key, nil
	} else {
		return nil, ErrKeysNotFound
	}
}

// SetDefaultAccount designates the account to be used by default
func (eth *ethereum) SetDefaultAccount(acct accounts.Account) {
	eth.defaultAccount = acct
}

// GetDefaultAccount returns the default account
func (eth *ethereum) GetDefaultAccount() accounts.Account {
	return eth.defaultAccount
}

// GetCoinbaseAddress returns the account to use for contract deploys
func (eth *ethereum) GetCoinbaseAddress() common.Address {
	return eth.coinbase
}

// GetBalance returns the ETHER balance of account specified
func (eth *ethereum) GetBalance(addr common.Address) (*big.Int, error) {
	ctx, cancel := eth.GetTimeoutContext()
	defer cancel()
	balance, err := eth.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (eth *ethereum) GetEndpoint() string {
	return eth.endpoint
}

func (eth *ethereum) GetTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), eth.timeout)
}

// GetSyncProgress returns a flag if we are syncing, a pointer to a struct if we are, or an error
func (eth *ethereum) GetSyncProgress() (bool, *geth.SyncProgress, error) {

	ctx, ctxCancel := eth.GetTimeoutContext()
	progress, err := eth.syncing(ctx)
	defer ctxCancel()

	if err == nil && progress == nil {
		return false, nil, nil
	}

	if err == nil && progress != nil {
		return true, progress, nil
	}

	return false, nil, err
}

func (eth *ethereum) GetEvents(ctx context.Context, firstBlock uint64, lastBlock uint64, addresses []common.Address) ([]types.Log, error) {

	logger := eth.logger

	logger.Debugf("...GetEvents(firstBlock:%v,lastBlock:%v,addresses:%x)", firstBlock, lastBlock, addresses)

	query := geth.FilterQuery{
		FromBlock: new(big.Int).SetUint64(firstBlock),
		ToBlock:   new(big.Int).SetUint64(lastBlock),
		Addresses: addresses}

	logs, err := eth.client.FilterLogs(ctx, query)
	if err != nil {
		logger.Errorf("Could not filter logs: %v", err)
		return nil, err
	}

	for idx, log := range logs {
		logger.Debugf("Log[%v] Block[%v]:%v", idx, log.BlockNumber, log)
		for idx, hash := range log.Topics {
			logger.Debugf("Hash[%v]:%x", idx, hash)
		}
	}

	return logs, nil
}

func (eth *ethereum) RetryCount() int {
	return eth.retryCount
}

// WaitForReceipt
func (eth *ethereum) WaitForReceipt(ctx context.Context, txn *types.Transaction) (*types.Receipt, error) {

	count := 1
	receipt, err := eth.client.TransactionReceipt(ctx, txn.Hash())

	// Ugly condition, because
	// -- Real endpoint returns err==geth.NotFound if receipt is nil
	// -- Simulated endpoint returns err==nil and receipt==nil until commit() is called
	for err == geth.NotFound || (err == nil && receipt == nil) {
		eth.logger.Debugf("Retry #%d getting receipt for %v ...", count, txn.Hash().Hex())
		count++
		SleepWithContext(ctx, eth.retryDelay)
		receipt, err = eth.client.TransactionReceipt(ctx, txn.Hash())
	}

	return receipt, err
}

func (eth *ethereum) RetryDelay() time.Duration {
	return eth.retryDelay
}

func (eth *ethereum) Timeout() time.Duration {
	return eth.timeout
}

func (eth *ethereum) GetTransactionOpts(ctx context.Context, account accounts.Account) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyStoreTransactor(eth.keystore, account)
	if err != nil {
		eth.logger.Errorf("could not create transactor for %v: %v", account.Address.Hex(), err)
	} else {
		opts.Context = ctx
		opts.Nonce = nil
		opts.Value = big.NewInt(0)
		opts.GasLimit = uint64(0)
		opts.GasPrice = nil
	}

	return opts, err
}

func (eth *ethereum) GetCallOpts(ctx context.Context, account accounts.Account) *bind.CallOpts { // TODO provide and use context
	return &bind.CallOpts{
		BlockNumber: nil,
		Context:     ctx,
		Pending:     false,
		From:        account.Address}
}

// TransferEther transfer's ether from one account to another, assumes from is unlocked
func (eth *ethereum) TransferEther(from common.Address, to common.Address, wei *big.Int) error {

	nonce, err := eth.client.PendingNonceAt(context.Background(), from)
	if err != nil {
		return err
	}

	gasPrice, err := eth.client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	var data []byte
	gasLimit := uint64(21000)
	tx := types.NewTransaction(nonce, to, wei, gasLimit, gasPrice, data)

	eth.logger.Debugf("TransferEther => chainID:%v from:%v nonce:%v, to:%v, wei:%v, gasLimit:%v, gasPrice:%v",
		eth.chainID, from.Hex(), nonce, to.Hex(), wei, gasLimit, gasPrice)

	signer := types.NewEIP155Signer(eth.chainID)

	signedTx, err := types.SignTx(tx, signer, eth.keys[from].PrivateKey)
	if err != nil {
		eth.logger.Error(err)
	}
	ctx, cancel := eth.GetTimeoutContext()
	defer cancel()
	err = eth.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return err
	}

	return nil
}

// GetCurrentHeight gets the height of the endpoints chain
func (eth *ethereum) GetCurrentHeight(ctx context.Context) (uint64, error) {
	header, err := eth.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	return header.Number.Uint64(), nil
}

// GetFinalizedHeight gets the height of the endpoints chain at which is is considered finalized
func (eth *ethereum) GetFinalizedHeight(ctx context.Context) (uint64, error) {
	height, err := eth.GetCurrentHeight(ctx)
	if err != nil {
		return height, err
	}

	if eth.finalityDelay >= height {
		return 0, nil
	}
	return height - eth.finalityDelay, nil

}

func (eth *ethereum) GetSnapshot() ([]byte, error) {
	return nil, nil
}

func (eth *ethereum) GetValidators() ([]common.Address, error) {
	c := eth.contracts
	validatorAddresses, err := c.Validators.GetValidators(eth.GetCallOpts(context.TODO(), eth.defaultAccount))
	if err != nil {
		eth.logger.Warnf("Could not call contract:%v", err)
		return nil, err
	}

	return validatorAddresses, nil
}

func (eth *ethereum) Clone(defaultAccount accounts.Account) Ethereum {
	nEth := *eth

	nEth.defaultAccount = defaultAccount

	return &nEth
}

func (c *Contracts) LookupContracts(registryAddress common.Address) error {

	eth := c.eth
	logger := eth.logger

	// Load the registry first
	registry, err := bindings.NewRegistry(registryAddress, eth.client)
	if err != nil {
		return err
	}
	c.Registry = registry
	c.RegistryAddress = registryAddress

	// Just a help for looking up other contracts
	lookup := func(name string) (common.Address, error) {
		addr, err := registry.Lookup(eth.GetCallOpts(context.TODO(), eth.defaultAccount), name)
		if err != nil {
			logger.Errorf("Failed lookup of \"%v\": %v", name, err)
		} else {
			logger.Infof("Lookup up of \"%v\" is 0x%x", name, addr)
		}
		return addr, err
	}

	c.CryptoAddress, err = lookup("crypto/v1")
	logAndEat(logger, err)

	c.Crypto, err = bindings.NewCrypto(c.CryptoAddress, eth.client)
	logAndEat(logger, err)

	c.DepositAddress, err = lookup("deposit/v1")
	logAndEat(logger, err)

	c.Deposit, err = bindings.NewDeposit(c.DepositAddress, eth.client)
	logAndEat(logger, err)

	c.EthdkgAddress, err = lookup("ethdkg/v1")
	logAndEat(logger, err)

	c.Ethdkg, err = bindings.NewETHDKG(c.EthdkgAddress, eth.client)
	logAndEat(logger, err)

	_, err = lookup("ethdkgCompletion/v1")
	logAndEat(logger, err)

	_, err = lookup("ethdkgGroupAccusation/v1")
	logAndEat(logger, err)

	_, err = lookup("ethdkgSubmitMPK/v1")
	logAndEat(logger, err)

	c.StakingAddress, err = lookup("staking/v1")
	logAndEat(logger, err)

	c.Staking, err = bindings.NewStaking(c.StakingAddress, eth.client)
	logAndEat(logger, err)

	c.StakingTokenAddress, err = lookup("stakingToken/v1")
	logAndEat(logger, err)

	c.StakingToken, err = bindings.NewToken(c.StakingTokenAddress, eth.client)
	logAndEat(logger, err)

	c.UtilityTokenAddress, err = lookup("utilityToken/v1")
	logAndEat(logger, err)

	c.UtilityToken, err = bindings.NewToken(c.UtilityTokenAddress, eth.client)
	logAndEat(logger, err)

	c.ValidatorsAddress, err = lookup("validators/v1")
	logAndEat(logger, err)

	// These all call the ValidatorsDiamond contract but we need facet specific bindings for bindings
	c.Validators, err = bindings.NewValidators(c.ValidatorsAddress, eth.client)
	logAndEat(logger, err)

	c.StakingValues, err = bindings.NewStakingValuesFacet(c.ValidatorsAddress, eth.client)
	logAndEat(logger, err)

	c.Participants, err = bindings.NewParticipantsFacet(c.ValidatorsAddress, eth.client)
	logAndEat(logger, err)

	c.Snapshots, err = bindings.NewSnapshotsFacet(c.ValidatorsAddress, eth.client)
	logAndEat(logger, err)

	return nil
}

func (c *Contracts) DeployContracts(ctx context.Context, account accounts.Account) (*bindings.Registry, common.Address, error) {
	eth := c.eth
	logger := eth.logger

	txnOpts, err := eth.GetTransactionOpts(ctx, account)
	if err != nil {
		return nil, common.Address{}, err
	}

	txnQueue := queue.New(10)
	q := func(tx *types.Transaction) {
		if tx != nil {
			logger.Infof("Queueing transaction %v", tx.Hash().String())
			txnQueue.Put(tx)
		} else {
			logger.Warn("Ignoring nil transaction")
		}
	}

	flushQ := func(queue *queue.Queue) {
		logger.Infof("waiting for txns...")
		for txns, err := queue.Get(1); !queue.Empty(); txns, err = queue.Get(1) {
			if err != nil {
				logger.Infof("failure: %v", err)
			}
			tx := txns[0].(*types.Transaction)
			logger.Infof("waiting for txn: %v", tx.Hash().String())
			eth.WaitForReceipt(ctx, tx)
		}
	}

	var txn *types.Transaction
	c.RegistryAddress, txn, c.Registry, err = bindings.DeployRegistry(txnOpts, eth.client)
	if err != nil {
		logger.Errorf("Failed to deploy registry...")
		return nil, common.Address{}, err
	}
	q(txn)
	logger.Infof("* registryAddress = \"0x%0.40x\"", c.RegistryAddress)

	c.CryptoAddress, txn, c.Crypto, err = bindings.DeployCrypto(txnOpts, eth.client)
	if err != nil {
		logger.Errorf("Failed to deploy crypto...")
		return nil, common.Address{}, err
	}
	q(txn)
	logger.Infof("  cryptoContract = \"0x%0.40x\"", c.CryptoAddress)

	c.StakingTokenAddress, txn, c.StakingToken, err = bindings.DeployToken(txnOpts, eth.client, StringToBytes32("STK"), StringToBytes32("MadNet Staking"))
	if err != nil {
		logger.Errorf("Failed to deploy stakingToken...")
		return nil, common.Address{}, err
	}
	q(txn)
	logger.Infof("  stakingTokenAddress = \"0x%0.40x\"", c.StakingTokenAddress)

	c.UtilityTokenAddress, txn, c.UtilityToken, err = bindings.DeployToken(txnOpts, eth.client, StringToBytes32("UTL"), StringToBytes32("MadNet Utility"))
	if err != nil {
		logger.Errorf("Failed to deploy utilityToken...")
		return nil, common.Address{}, err
	}
	q(txn)
	logger.Infof("  utilityTokenAddress = \"0x%0.40x\"", c.UtilityTokenAddress)

	c.DepositAddress, txn, c.Deposit, err = bindings.DeployDeposit(txnOpts, eth.client, c.RegistryAddress)
	if err != nil {
		logger.Errorf("Failed to deploy deposit...")
		return nil, common.Address{}, err
	}
	q(txn)
	logger.Infof("  depositAddress = \"0x%0.40x\"", c.DepositAddress)

	c.StakingAddress, txn, c.Staking, err = bindings.DeployStaking(txnOpts, eth.client, c.RegistryAddress)
	if err != nil {
		logger.Errorf("Failed to deploy staking...")
		return nil, common.Address{}, err
	}
	q(txn)
	logger.Infof("  stakingAddress = \"0x%0.40x\"", c.StakingAddress)

	// Deploy ValidatorsDiamond
	c.ValidatorsAddress, txn, _, err = bindings.DeployValidatorsDiamond(txnOpts, eth.client) // Deploy the core diamond
	if err != nil {
		logger.Errorf("Failed to deploy validators diamond...")
		return nil, common.Address{}, err
	}
	q(txn)

	// Deploy validators facets
	participantsFacet, txn, _, err := bindings.DeployParticipantsFacet(txnOpts, eth.client)
	if err != nil {
		logger.Error("Failed to deploy participants facet...")
		return nil, common.Address{}, err
	}
	q(txn)

	snapshotsFacet, txn, _, err := bindings.DeploySnapshotsFacet(txnOpts, eth.client)
	if err != nil {
		logger.Error("Failed to deploy snapshots facet...")
		return nil, common.Address{}, err
	}
	q(txn)

	stakingValuesFacet, txn, _, err := bindings.DeployStakingValuesFacet(txnOpts, eth.client)
	if err != nil {
		logger.Error("Failed to deploy staking values facet...")
		return nil, common.Address{}, err
	}
	q(txn)

	// validatorsUpdateFacet, txn, _, err := bindings.DeployValidatorsUpdateFacet(txnOpts, eth.client)
	// if err != nil {
	// 	logger.Error("Failed to deploy validators update facet...")
	// 	return nil, common.Address{}, err
	// }
	// q(txn)

	// Bind diamond to interfaces
	c.Validators, err = bindings.NewValidators(c.ValidatorsAddress, eth.client)
	logAndEat(logger, err)

	c.StakingValues, err = bindings.NewStakingValuesFacet(c.ValidatorsAddress, eth.client)
	logAndEat(logger, err)

	c.Participants, err = bindings.NewParticipantsFacet(c.ValidatorsAddress, eth.client)
	logAndEat(logger, err)

	c.Snapshots, err = bindings.NewSnapshotsFacet(c.ValidatorsAddress, eth.client)
	logAndEat(logger, err)

	c.Validators, err = bindings.NewValidators(c.ValidatorsAddress, eth.client) // Validators is just an interface
	if err != nil {
		logger.Errorf("Failed to deploy validators...")
		return nil, common.Address{}, err
	}
	logger.Infof("  validatorsAddress = \"0x%0.40x\"", c.ValidatorsAddress)

	validatorsUpdate, err := bindings.NewValidatorsUpdateFacet(c.ValidatorsAddress, eth.client)
	if err != nil {
		logger.Errorf("Failed to deploy validators update  ..")
		return nil, common.Address{}, err
	}

	// Wait for all the deploys to finish
	eth.commit()
	flushQ(txnQueue)

	// Register all the facets
	vu := &updater{vu: validatorsUpdate, txnOpts: txnOpts}

	// Value maintenance
	q(vu.add("majorStakeFine()", stakingValuesFacet))
	q(vu.add("minimumStake()", stakingValuesFacet))
	q(vu.add("minorStakeFine()", stakingValuesFacet))
	q(vu.add("rewardAmount()", stakingValuesFacet))
	q(vu.add("rewardBonus()", stakingValuesFacet))
	q(vu.add("setMajorStakeFine(uint256)", stakingValuesFacet))
	q(vu.add("setMinimumStake(uint256)", stakingValuesFacet))
	q(vu.add("setMinorStakeFine(uint256)", stakingValuesFacet))
	q(vu.add("setRewardAmount(uint256)", stakingValuesFacet))
	q(vu.add("setRewardBonus(uint256)", stakingValuesFacet))

	// Snapshot maintenance
	q(vu.add("initializeSnapshots(address)", snapshotsFacet))
	q(vu.add("snapshot(bytes,bytes)", snapshotsFacet))
	q(vu.add("setNextSnapshot(uint256)", snapshotsFacet))
	q(vu.add("nextSnapshot()", snapshotsFacet))

	q(vu.add("getChainIdFromSnapshot(uint256)", snapshotsFacet))
	q(vu.add("getRawBlockClaimsSnapshot(uint256)", snapshotsFacet))
	q(vu.add("getRawSignatureSnapshot(uint256)", snapshotsFacet))
	q(vu.add("getHeightFromSnapshot(uint256)", snapshotsFacet))
	q(vu.add("getMadHeightFromSnapshot(uint256)", snapshotsFacet))

	// Validator maintenance
	q(vu.add("initializeParticipants(address)", participantsFacet))
	q(vu.add("addValidator(address,uint256[2])", participantsFacet))
	q(vu.add("removeValidator(address,uint256[2])", participantsFacet))
	q(vu.add("queueValidator(address,uint256[2])", participantsFacet))
	q(vu.add("isValidator(address)", participantsFacet))
	q(vu.add("getValidatorPublicKey(address)", participantsFacet))
	q(vu.add("confirmValidators()", participantsFacet))
	q(vu.add("validatorMaxCount()", participantsFacet))
	q(vu.add("validatorCount()", participantsFacet))
	q(vu.add("setValidatorMaxCount(uint8)", participantsFacet))

	c.EthdkgAddress, txn, c.Ethdkg, err = bindings.DeployETHDKG(txnOpts, eth.client, c.RegistryAddress)
	if err != nil {
		logger.Errorf("Failed to deploy ethdkg...")
		return nil, common.Address{}, err
	}
	q(txn)
	logger.Infof("  ethdkgAddress = \"0x%0.40x\"", c.EthdkgAddress)

	var ethdkgCompletionAddress common.Address
	ethdkgCompletionAddress, txn, _, err = bindings.DeployETHDKGCompletion(txnOpts, eth.client)
	if err != nil {
		logger.Errorf("Failed to deploy ethdkgCompletion...")
		return nil, common.Address{}, err
	}
	q(txn)
	logger.Infof("  ethdkgCompletion = \"0x%0.40x\"", ethdkgCompletionAddress)

	var ethdkgGroupAccusationAddress common.Address
	ethdkgGroupAccusationAddress, txn, _, err = bindings.DeployETHDKGGroupAccusation(txnOpts, eth.client)
	if err != nil {
		logger.Errorf("Failed to deploy ethdkgGroupAccusation...")
		return nil, common.Address{}, err
	}
	q(txn)
	logger.Infof("  ethdkgGroupAccusation = \"0x%0.40x\"", ethdkgGroupAccusationAddress)

	var ethdkgSubmitMPKAddress common.Address
	ethdkgSubmitMPKAddress, txn, _, err = bindings.DeployETHDKGSubmitMPK(txnOpts, eth.client)
	if err != nil {
		logger.Errorf("Failed to deploy ethdkgSubmitMPKAddress...")
		return nil, common.Address{}, err
	}
	q(txn)
	logger.Infof(" ethdkgSubmitMPKAddress = \"0x%0.40x\"", ethdkgSubmitMPKAddress)

	eth.contracts = c
	eth.commit()

	txn, err = c.Registry.Register(txnOpts, "crypto/v1", c.CryptoAddress)
	logAndEat(logger, err)
	q(txn)

	txn, err = c.Registry.Register(txnOpts, "deposit/v1", c.DepositAddress)
	logAndEat(logger, err)
	q(txn)

	txn, err = c.Registry.Register(txnOpts, "ethdkg/v1", c.EthdkgAddress)
	logAndEat(logger, err)
	q(txn)

	txn, err = c.Registry.Register(txnOpts, "ethdkgCompletion/v1", ethdkgCompletionAddress)
	logAndEat(logger, err)
	q(txn)

	txn, err = c.Registry.Register(txnOpts, "ethdkgGroupAccusation/v1", ethdkgGroupAccusationAddress)
	logAndEat(logger, err)
	q(txn)

	txn, err = c.Registry.Register(txnOpts, "ethdkgSubmitMPK/v1", ethdkgSubmitMPKAddress)
	logAndEat(logger, err)
	q(txn)

	txn, err = c.Registry.Register(txnOpts, "staking/v1", c.StakingAddress)
	logAndEat(logger, err)
	q(txn)

	txn, err = c.Registry.Register(txnOpts, "stakingToken/v1", c.StakingTokenAddress)
	logAndEat(logger, err)
	q(txn)

	txn, err = c.Registry.Register(txnOpts, "utilityToken/v1", c.UtilityTokenAddress)
	logAndEat(logger, err)
	q(txn)

	txn, err = c.Registry.Register(txnOpts, "validators/v1", c.ValidatorsAddress)
	logAndEat(logger, err)
	q(txn)

	eth.commit()

	// Wait for all the deploys to finish
	flushQ(txnQueue)

	// Initialize Snapshots facet
	tx, err := c.Snapshots.InitializeSnapshots(txnOpts, c.RegistryAddress)
	if err != nil {
		logger.Errorf("Failed to initialize SnapshotsFacet: %v", err)
		return nil, common.Address{}, err
	}
	eth.commit()
	rcpt, err := eth.WaitForReceipt(ctx, tx)
	if err != nil {
		logger.Errorf("Failed to get receipt for initializing Snapshots facet: %v", err)
		return nil, common.Address{}, err
	}
	if rcpt != nil {
		logger.Infof("Snapshots update status: %v", rcpt.Status)
	} else {
		logger.Errorf("Snapshots update receipt is nil")
	}

	tx, err = c.Snapshots.SetNextSnapshot(txnOpts, big.NewInt(1))
	if err != nil {
		logger.Errorf("Failed to initialize Snapshots facet next snapshot: %v", err)
		return nil, common.Address{}, err
	}
	q(tx)

	// Default staking values
	tx, err = c.StakingValues.SetMinimumStake(txnOpts, big.NewInt(1000000))
	logAndEat(logger, err)
	q(tx)

	tx, err = c.StakingValues.SetMajorStakeFine(txnOpts, big.NewInt(200000))
	logAndEat(logger, err)
	q(tx)

	tx, err = c.StakingValues.SetMinorStakeFine(txnOpts, big.NewInt(50000))
	logAndEat(logger, err)
	q(tx)

	tx, err = c.StakingValues.SetRewardAmount(txnOpts, big.NewInt(1000))
	logAndEat(logger, err)
	q(tx)

	tx, err = c.StakingValues.SetRewardBonus(txnOpts, big.NewInt(1000))
	logAndEat(logger, err)
	q(tx)

	eth.commit()

	flushQ(txnQueue)

	// Initialize Participants facet
	tx, err = c.Participants.InitializeParticipants(txnOpts, c.RegistryAddress)
	if err != nil {
		logger.Errorf("Failed to initialize Participants facet: %v", err)
		return nil, common.Address{}, err
	}
	eth.commit()

	rcpt, err = eth.WaitForReceipt(ctx, tx)
	if err != nil {
		logger.Errorf("Failed to get receipt for initializing Participants facet: %v", err)
		return nil, common.Address{}, err
	}
	if rcpt != nil {
		logger.Infof("Snapshots update status: %v", rcpt.Status)
	} else {
		logger.Errorf("Snapshots update receipt is nil")
	}

	tx, err = c.Participants.SetValidatorMaxCount(txnOpts, 10)
	if err != nil {
		logger.Errorf("Failed to initialize Participants facet: %v", err)
		return nil, common.Address{}, err
	}
	q(tx)
	eth.commit()
	flushQ(txnQueue)

	// Staking updates
	tx, err = c.Staking.ReloadRegistry(txnOpts)
	if err != nil {
		logger.Errorf("Failed to update staking contract references: %v", err)
		return nil, common.Address{}, err
	}
	eth.commit()
	rcpt, err = eth.WaitForReceipt(ctx, tx)
	if err != nil {
		logger.Errorf("Failed to get receipt for staking update: %v", err)
		return nil, common.Address{}, err

	}
	if rcpt != nil {
		logger.Infof("staking update status: %v", rcpt.Status)
	} else {
		logger.Errorf("staking receipt is nil")
	}

	// Deposit updates
	tx, err = c.Deposit.ReloadRegistry(txnOpts)
	if err != nil {
		logger.Errorf("Failed to update deposit contract references: %v", err)
		return nil, common.Address{}, err
	}
	eth.commit()
	rcpt, err = eth.WaitForReceipt(ctx, tx)
	if err != nil {
		logger.Errorf("Failed to get receipt for deposit update: %v", err)
		return nil, common.Address{}, err
	} else if rcpt != nil {
		logger.Infof("deposit update status: %v", rcpt.Status)
	}

	// Validator updates
	// tx, err = c.Validators.ReloadRegistry(txnOpts)
	// if err != nil {
	// 	logger.Errorf("Failed to update validators contract references: %v", err)
	// 	return nil, common.Address{}, err
	// }
	eth.commit()
	rcpt, err = eth.WaitForReceipt(ctx, tx)
	if err != nil {
		logger.Errorf("Failed to get receipt for validators update: %v", err)
		return nil, common.Address{}, err
	} else if rcpt != nil {
		logger.Infof("validators update status: %v", rcpt.Status)
	}

	// ETHDKG updates
	tx, err = c.Ethdkg.ReloadRegistry(txnOpts)
	if err != nil {
		logger.Errorf("Failed to update ethdkg contract references: %v", err)
		return nil, common.Address{}, err
	}
	eth.commit()
	rcpt, err = eth.WaitForReceipt(ctx, tx)
	if err != nil {
		logger.Errorf("Failed to get receipt for ethdkg update: %v", err)
		return nil, common.Address{}, err
	} else if rcpt != nil {
		logger.Infof("ethdkg update status: %v", rcpt.Status)
	}

	return c.Registry, c.RegistryAddress, nil
}

// StringToBytes32 is useful for convert a Go string into a bytes32 useful calling Solidity
func StringToBytes32(str string) (b [32]byte) {
	copy(b[:], []byte(str)[0:32])
	return
}

// CalculateSelector calculates the hash of the supplied function signature
func CalculateSelector(functionSignature string) [4]byte {
	var selector [4]byte

	selectorSlice := crypto.Keccak256([]byte(functionSignature))[:4]
	selector[0] = selectorSlice[0]
	selector[1] = selectorSlice[1]
	selector[2] = selectorSlice[2]
	selector[3] = selectorSlice[3]

	return selector
}

func logAndEat(logger *logrus.Logger, err error) {
	if err != nil {
		logger.Error(err)
	}
}

type updater struct {
	err     error
	vu      *bindings.ValidatorsUpdateFacet
	txnOpts *bind.TransactOpts
}

//
func (u *updater) add(signature string, facet common.Address) *types.Transaction {
	if u.err != nil {
		return nil
	}

	selector := CalculateSelector(signature)
	txn, err := u.vu.AddFacet(u.txnOpts, selector, facet)
	u.err = err
	return txn
}
