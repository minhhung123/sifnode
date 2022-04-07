#!/usr/bin/env bash


# DENOM='rowan'
# CHAIN_ID='localnet'
# RPC='http://localhost:26657/'
# TXFLAG="--gas-prices 0.1$DENOM --gas auto --gas-adjustment 1.3 -y -b block --chain-id $CHAIN_ID --node $RPC"
BINARY=sifnoded

CONTRACT_CODE=$($BINARY tx wasm store artifacts/cw_erc20.wasm --from sif --keyring-backend=test --gas-prices 1rowan --gas auto --gas-adjustment 1.3 -y -b block --chain-id=localnet --node=http://localhost:26657/  --output json | jq -r '.logs[0].events[-1].attributes[0].value')
echo "Stored: $CONTRACT_CODE"

INIT='{"name":"Hung Coin","symbol":"HUNG","decimals":6,"initial_balances":[{"address":"<validator-self-delegate-address>","amount":"12345678000"}]}'

$BINARY tx wasm instantiate $CONTRACT_CODE "$INIT" --amount 50000rowan --label "Hungcoin erc20" --from sif --keyring-backend=test --chain-id=localnet --gas-prices 1rowan --gas auto --gas-adjustment 1.3 -b block -y --no-admin

CONTRACT_ADDRESS=$($BINARY query wasm list-contract-by-code $CONTRACT_CODE --output json | jq -r '.contracts[-1]')
echo "Contract address: $CONTRACT_ADDRESS"
