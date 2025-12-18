package cli

import (
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "training",
		Short:                      "Training transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	return cmd
}

