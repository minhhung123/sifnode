#!/bin/env sh

# fund registry account
# set registry
# margin param change
# create pools

sifnoded tx clp create-pool --from tempnet2 --symbol ceth --nativeAmount 1000000000000000000 --externalAmount 1000000000000000000 --keyring-backend test --chain-id $marginchain2 --node $marginnet2 --fees 100000000000000000rowan
sifnoded tx clp create-pool --from tempnet2 --symbol cusdt --nativeAmount 1000000000000000000 --externalAmount 1000000000000000000 --keyring-backend test --chain-id $marginchain2 --node $marginnet2 --fees 100000000000000000rowan
sifnoded tx clp create-pool --from tempnet2 --symbol stake --nativeAmount 1000000000000000000 --externalAmount 1000000000000000000 --keyring-backend test --chain-id $marginchain2 --node $marginnet2 --fees 100000000000000000rowan
