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
	"math/big"
	"strings"
	"github.com/Ankr-network/dccn-common/wallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

)

// adminCmd represents the admin command
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "admin is used to do admin operations ",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var (
	//names of sub command bind in viper, which is used to bind flags
	// naming notions subCmdNameKey. eg. "account" have a flag named "url" shall named as accountUrl,
	//persistent flags
	adminUrl        = "adminUrl"
	adminPrivateKey = "adminPrivateKey"

	//sub cmd flags
	setBalAddr     = "setBalAddr"
	setBalAmount   = "setBalAmount"
	setCertDc      = "setCertDc"
	setCertPerm    = "setCertPerm"
	setValidPub    = "setValidPub"
	setValidPower  = "setValidPower"
	setStakeAmount = "setStakeAmount"
	setStakePub    = "setStakePub"
	removeValidPub = "removeValidPub"
	removeCertDc   = "removeCertDc"
)

func init() {
	//init persistent flags and append sub commands
	err := addPersistentString(adminCmd, adminUrl, urlFlag, "", "", "url of a validator", required)
	if err != nil {
		panic(err)
	}

	err = addPersistentString(adminCmd, adminPrivateKey, privkeyFlag, "", "", "operator private key", required)
	if err != nil {
		panic(err)
	}

	//add sub cmd to adminCmd
	appendSubCmd(adminCmd, "setbalance", "set target account with specified amount", setBalance, addSetBalanceFlag)
	appendSubCmd(adminCmd, "setcert", "set metering cert", setCert, addCertFlags)
	appendSubCmd(adminCmd, "setvalidator", "add a new validator", setValidator, addSetValidatorFlags)
	appendSubCmd(adminCmd, "setstake", "set stake", setStake, addSetStakeFlags)
	appendSubCmd(adminCmd, "removevalidator", "remove a validator", removeValidator, addRemoveValidatorFlags)
	appendSubCmd(adminCmd, "removecert", "remove cert from validator", removeCert, addRemoveCertFlags)
}

//admin setbalance --address 0xXXX --amount 12313 --url https://validator-url:port
func setBalance(cmd *cobra.Command, args []string) {
	validatorUrl = viper.GetString(adminUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	address := viper.GetString(setBalAddr)
	amount := viper.GetString(setBalAmount)
	operatorPriv := viper.GetString(adminPrivateKey)
	index := strings.LastIndex(validatorUrl, ":")
	if index < 0 {
		fmt.Println("Error: url is not correct, example 'https://chain-01.dccn.ankr.com:443'")
		return
	}
	err := wallet.SetBalance(validatorUrl[:index], validatorUrl[index+1:], address, amount, operatorPriv)
	if err != nil {
		fmt.Println("Error: Set balance failed.", err)
		return
	}
	bal, err := wallet.GetBalance(validatorUrl[:index], validatorUrl[index+1:], address)
	if err != nil {
		fmt.Println("Error: Set balance failed", err)
		return
	}
	balance, bo := new(big.Int).SetString(bal, 10)
	if bo != true {
		fmt.Println("Get Balance error:", err)
		return
	}
	fmt.Println("Set balance Success.")
	fmt.Println("Address:", address)
	fmt.Println("Balance:", balance.Div(balance, AnkrBase))

}

func addSetBalanceFlag(cmd *cobra.Command) {
	err := addStringFlag(cmd, setBalAddr, addressFlag, "", "", "the address of the target account to receive ankr token", required)
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, setBalAmount, amountFlag, "", "", "the amount to set to the target address", required)
	if err != nil {
		panic(err)
	}
}

//admin setcert --dcname dataCenterName --certPerm certString --url https://validator-url:port
func setCert(cmd *cobra.Command, args []string) {
	validatorUrl = viper.GetString(adminUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	index := strings.LastIndex(validatorUrl, ":")
	if index < 0 {
		fmt.Println("Error: url is not correct, example 'https://chain-01.dccn.ankr.com:443'")
		return
	}
	opPrivateKey := viper.GetString(adminPrivateKey)
	if len(opPrivateKey) < 1 {
		fmt.Println("Invalid operator private key!")
		return
	}
	dcName := viper.GetString(setCertDc)
	certPerm := viper.GetString(setCertPerm)
	if len(certPerm) < 1 {
		fmt.Println("Invalid cert perm!")
		return
	}
	err := wallet.SetMeteringCert(validatorUrl[:index], validatorUrl[index+1:], opPrivateKey, dcName, certPerm)
	if err != nil {
		fmt.Println("Set metering cert failed", err)
		return
	}
}

func addCertFlags(cmd *cobra.Command) {
	err := addStringFlag(cmd, setCertDc, dcnameFlag, "", "", "data center name", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, setCertPerm, permFlag, "", "", "cert perm to be set", required)
	if err != nil {
		panic(err)
	}
}

// setvalidator --pubkey jlds --power 21 --url --privkey
func setValidator(cmd *cobra.Command, args []string) {
	validatorUrl = viper.GetString(adminUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	index := strings.LastIndex(validatorUrl, ":")
	if index < 0 {
		fmt.Println("Error: url is not correct, example 'https://chain-01.dccn.ankr.com:443'")
		return
	}

	validatorPub := viper.GetString(setValidPub)
	if len(validatorPub) < 1 {
		fmt.Println("Invalid validator public key!")
		return
	}

	validatorPower := viper.GetString(setValidPower)
	if len(validatorPower) < 1 {
		fmt.Println("Invalid validator power!")
		return
	}

	opPrivateKey := viper.GetString(adminPrivateKey)
	if len(opPrivateKey) < 1 {
		fmt.Println("Invalid operator private key!")
		return
	}
	err := wallet.SetValidator(validatorUrl[:index], validatorUrl[index+1:], validatorPub, validatorPower, opPrivateKey)
	if err != nil {
		fmt.Println("Set validator failed!", err)
		return
	}
	fmt.Println("Set validator success.")

}

func addSetValidatorFlags(cmd *cobra.Command) {
	err := addStringFlag(cmd, setValidPub, pubkeyFlag, "", "", "the public address of the added validator", required)
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, setValidPower, powerFlag, "", "", "the power set to the validator", required)
	if err != nil {
		panic(err)
	}
}

// setstake --amount 3
func setStake(cmd *cobra.Command, args []string) {
	validatorUrl = viper.GetString(adminUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	index := strings.LastIndex(validatorUrl, ":")
	if index < 0 {
		fmt.Println("Error: url is not correct, example 'https://chain-01.dccn.ankr.com:443'")
		return
	}
	opPrivateKey := viper.GetString(adminPrivateKey)
	if len(opPrivateKey) < 1 {
		fmt.Println("Invalid operator private key!")
		return
	}

	amount := viper.GetString(setStakeAmount)
	_, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		fmt.Println("Invalid Stake Amount!")
		return
	}
	stakePub := viper.GetString(setStakePub)
	if len(stakePub) < 1 {
		fmt.Println("Invalid public key!")
		return
	}
	err := wallet.SetStake(validatorUrl[:index], validatorUrl[index+1:], opPrivateKey, amount, stakePub)
	if err != nil {
		fmt.Println("Set stake failed.", err)
		return
	}
	fmt.Println("Set Stake success.")
}

func addSetStakeFlags(cmd *cobra.Command) {
	err := addStringFlag(cmd, setStakeAmount, amountFlag, "", "", "set stake amount", required)
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, setStakePub, pubkeyFlag, "", "", "public key", required)
	if err != nil {
		panic(err)
	}
}

// removevalidator --pubkey string
func removeValidator(cmd *cobra.Command, args []string) {
	validatorUrl = viper.GetString(adminUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	index := strings.LastIndex(validatorUrl, ":")
	if index < 0 {
		fmt.Println("Error: url is not correct, example 'https://chain-01.dccn.ankr.com:443'")
		return
	}
	opPrivateKey := viper.GetString(adminPrivateKey)
	if len(opPrivateKey) < 1 {
		fmt.Println("Invalid operator private key!")
		return
	}

	validatorPub := viper.GetString(removeValidPub)
	if len(validatorPub) < 1 {
		fmt.Println("Invalid validator public key!")
		return
	}

	err := wallet.RemoveValidator(validatorUrl[:index], validatorUrl[index+1:], validatorPub, opPrivateKey)
	if err != nil {
		fmt.Println("Remove validator failed:", err)
		return
	}
	fmt.Println("Remove validator success.")
}

func addRemoveValidatorFlags(cmd *cobra.Command) {
	err := addStringFlag(cmd, removeValidPub, pubkeyFlag, "", "", "public key of the to be removed validator", required)
	if err != nil {
		panic(err)
	}
}

// removecert --pubkey string
func removeCert(cmd *cobra.Command, args []string) {
	validatorUrl = viper.GetString(adminUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	index := strings.LastIndex(validatorUrl, ":")
	if index < 0 {
		fmt.Println("Error: url is not correct, example 'https://chain-01.dccn.ankr.com:443'")
		return
	}
	opPrivateKey := viper.GetString(adminPrivateKey)
	if len(opPrivateKey) < 1 {
		fmt.Println("Invalid operator private key!")
		return
	}
	removeDcName := viper.GetString(removeCertDc)
	if len(removeCertDc) < 1 {
		fmt.Println("Invalid data center name!")
		return
	}
	err := wallet.RemoveMeteringCert(validatorUrl[:index], validatorUrl[index+1:], opPrivateKey, removeDcName)
	if err != nil {
		fmt.Println("Failed to remove cert:", err)
		return
	}
	fmt.Println("Remove cert success.")
}

func addRemoveCertFlags(cmd *cobra.Command) {
	err := addStringFlag(cmd, removeCertDc, dcnameFlag, "", "", "name of data center name", required)
	if err != nil {
		panic(err)
	}
}
