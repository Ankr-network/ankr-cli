package main

import (
	"fmt"
	"github.com/tendermint/tendermint/cmd/ankr_cli/cmd"
)

func main()  {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
