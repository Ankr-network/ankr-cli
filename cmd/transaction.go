/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/Ankr-network/dccn-common/wallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var (
	validatorUrl string
	privateKey   string
	txUrlFlag    = "txUrlFlag"

	//names of flags used in viper to bind keys
	transferTo      = "transferTo"
	transferAmount  = "transferAmount"
	transferKeyfile = "transferKeyfile"
	meteringDc      = "meteringDc"
	meteringNs      = "meteringNs"
	meteringValue   = "meteringValue"
	meteringPriv    = "meteringPriv"
)

// transactionCmd represents the transaction command
var transactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "transaction is used to send coins to specified address or send metering",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("transaction called")
	},
}

func init() {
	err := addPersistentString(transactionCmd, txUrlFlag, urlFlag, "", "", "the url of a validator", required)
	if err != nil {
		panic(err)
	}
	appendSubCmd(transactionCmd, "transfer", "send coins to another account", transfer, addTransferFlag)
	appendSubCmd(transactionCmd, "metering", "send metering transaction", sendMetering, addMeteringFlags)
}

//transaction transfer functions
func transfer(cmd *cobra.Command, args []string) {
	keystorePath := viper.GetString(transferKeyfile)
	for i, arg := range args {
		fmt.Println("arg", i, ":", arg)
	}
	_, err := os.Stat(keystorePath)
	if err != nil {
		fmt.Println("Error: Keystore does not exist!")
		return
	}
	privateKey := decryptPrivatekey(keystorePath)
	if privateKey == "" {
		fmt.Println("Error: Wrong keystore or password!")
		return
	}
	acc, err := getAccountFromPrivatekey(privateKey)
	if err != nil {
		fmt.Println("Error: generate account from private key", err)
		return
	}
	if privateKey == "" {
		fmt.Println("Error: Wrong private keystore or password!")
		return
	}

	validatorUrl = viper.GetString(txUrlFlag)
	index := strings.LastIndex(validatorUrl, ":")
	if index < 0 {
		fmt.Println("Error: url is not correct, example 'https://chain-01.dccn.ankr.com:443'")
		return
	}

	to := viper.GetString(transferTo)
	amount := viper.GetString(transferAmount)
	hash, err := wallet.SendCoins(validatorUrl[:index], validatorUrl[index+1:], privateKey, acc.Address, to, amount)
	if err != nil {
		fmt.Println("\nTransfer encountered an error:", err)
		return
	}
	fmt.Println("\nTransaction sent. Tx hash:", hash)

}

func addTransferFlag(cmd *cobra.Command) {
	err := addStringFlag(cmd, transferTo, toFlag, "", "", "", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, transferAmount, amountFlag, "", "", "", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, transferKeyfile, keystoreFlag, "", "", "", required)
	if err != nil {
		panic(err)
	}
}

//transaction metering function
func sendMetering(cmd *cobra.Command, args []string) {
	validatorUrl = viper.GetString(txUrlFlag)
	index := strings.LastIndex(validatorUrl, ":")
	if index < 0 {
		fmt.Println("Error: url is not correct, example 'https://chain-01.dccn.ankr.com:443'")
		return
	}
	dc := viper.GetString(meteringDc)
	ns := viper.GetString(meteringNs)
	value := viper.GetString(meteringValue)
	privateKey = viper.GetString(meteringPriv)
	err := wallet.SetMetering(validatorUrl[:index], validatorUrl[index+1:], privateKey, dc, ns, value)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Set metering success.")
}

func addMeteringFlags(cmd *cobra.Command) {
	//cmd.Flags().StringVarP(&privateKey, "privkey", "p", "", "admin private key")
	err := addStringFlag(cmd, meteringDc, dcnameFlag, "", "", "data center name", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, meteringNs, nameSpaceFlag, "", "", "namespace", required)
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, meteringValue, valueFlag, "", "", "the value to be set", required)
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, meteringPriv, privkeyFlag, "", "", "admin private key", required)
	if err != nil {
		panic(err)
	}
}
