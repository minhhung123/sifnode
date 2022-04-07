#!/usr/bin/env bash

BINARY=sifnoded
DENOM='rowan'
CHAIN_ID='localnet'
RPC='http://localhost:26657/'
TXFLAG="--gas-prices 1$DENOM --gas auto --gas-adjustment 1.3 -y -b block --chain-id $CHAIN_ID --node $RPC"


CONTRACT_CODE=6
CONTRACT_ADDRESS="sif12fykm2xhg5ces2vmf4q2aem8c958exv3v0wmvrspa8zucrdwjeds8kfpj8"

# Query contract metadata
CONTRACT_STATE=$($BINARY query wasm contract-state all $CONTRACT_ADDRESS --output json)
echo "Contract state: $CONTRACT_STATE"