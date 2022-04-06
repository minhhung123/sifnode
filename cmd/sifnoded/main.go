package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/cmd/sifnoded/cmd"
)

func main() {
	cmdOptions := GetWasmCmdOptions()
	rootCmd, _ := cmd.NewRootCmd(cmdOptions...)

	app.SetConfig(true)

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
