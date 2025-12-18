package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/atlas/chain/app"
	"github.com/atlas/chain/cmd/atlasd/cmd"
)

func main() {
	encodingConfig := params.MakeTestEncodingConfig()
	app.SetEncodingConfig(encodingConfig)

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}

