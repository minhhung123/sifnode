package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

const (
	DefaultSifInstanceCost uint64 = 60_000
	DefaultSifCompileCost uint64 = 100
)

func SifGasRegisterConfig() wasmkeeper.WasmGasRegisterConfig {
	gasConfig := wasmkeeper.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultSifInstanceCost
	gasConfig.CompileCost = DefaultSifCompileCost

	return gasConfig
}

func NewSifWasmGasRegister() wasmkeeper.WasmGasRegister {
	return wasmkeeper.NewWasmGasRegister(SifGasRegisterConfig())
}