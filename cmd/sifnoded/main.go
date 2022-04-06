package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/cmd/sifnoded/cmd"
)

func main() {
	cmdOptions := GetWasmCmdOptions()
	rootCmd, _ := cmd.NewRootCmd(cmdOptions...)

	app.SetConfig(true)

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
