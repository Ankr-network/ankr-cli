package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd is the root command for Tendermint core.
var RootCmd = &cobra.Command{
	Use:   "ankr-cli",
	Short: "ankr-cli is used to interacting with ankr blockchain",
	ValidArgs:nil,
}

func init() {
	//add sub commands for ankr_cli
	RootCmd.AddCommand(accountCmd)
	RootCmd.AddCommand(transactionCmd)
	RootCmd.AddCommand(adminCmd)
	RootCmd.AddCommand(queryCmd)
	//RootCmd.GenBashCompletionFile("ankr-chain-cli.sh")
}
