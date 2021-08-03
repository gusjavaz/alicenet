#!/bin/sh
DATADIR=./local-geth/

geth --miner.threads 1 --miner.gasprice 1 --miner.gaslimit 10000000 --miner.etherbase 546f99f244b7b58b855330ae0e2bc1b30b41302f --nodiscover --mine --txpool.nolocals --maxpeers 0 --ws --ws.addr=0.0.0.0 --ws.port=8546 --ws.api="eth,net,web3" --rpc --rpcaddr=0.0.0.0 --rpcport=8545 --datadir=$DATADIR --nousb --rpcapi="admin,eth,net,web3,personal,miner" --rpc.allow-unprotected-txs
