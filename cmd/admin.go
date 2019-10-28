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
	client2 "github.com/Ankr-network/ankr-chain/client"
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/tx/metering"
	"github.com/Ankr-network/ankr-chain/tx/serializer"
	"github.com/Ankr-network/ankr-chain/tx/validator"
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
	adminChId = "adminChId"
	adminGasLimt = "adminGasLimt"
	adminGasPrice = "adminGasPrice"
	adminMemo = "adminMemo"
	adminVersion = "adminVersion"


	//sub cmd flags
	setBalAddr          = "setBalAddr"
	setBalAmount        = "setBalAmount"
	setCertDc           = "setCertDc"
	setCertPerm         = "setCertPerm"
	setValidPub         = "setValidPub"
	setValidPower       = "setValidPower"
	setValidAction      = "setValidAction"
	setValidName        = "setValidName"
	setValidStakeAddr   = "setValidStakeAddr"
	setValidStakeAmount = "setValidStakeAmount"
	setValidStakeHeight = "setValidStakeHeight"
	setValidFlag        = "setValidFlag"
	setValidGasUsed     = "setValidGasUsed"
	setStakeAmount      = "setStakeAmount"
	setStakePub         = "setStakePub"
	removeValidPub      = "removeValidPub"
	removeCertDc        = "removeCertDc"
	removeCertNs = "removeCertNs"
)

func init() {
	//init persistent flags and append sub commands
	err := addPersistentString(adminCmd, adminUrl, urlParam, "", "", "url of a validator", required)
	if err != nil {
		panic(err)
	}

	err = addPersistentString(adminCmd, adminPrivateKey, privkeyParam, "", "", "operator private key", required)
	if err != nil {
		panic(err)
	}
	err = addPersistentString(adminCmd, adminChId, chainIDParam, "", "ankr-chain", "block chain id", notRequired)
	if err != nil {
		panic(err)
	}
	err = addPersistentInt(adminCmd, adminGasPrice, gasPriceParam, "", 0, "gas price", notRequired)
	if err != nil {
		panic(err)
	}

	err = addPersistentString(adminCmd, adminMemo, memoParam, "", "", "transaction memo", notRequired)
	if err != nil {
		panic(err)
	}
	err = addPersistentInt(adminCmd, adminGasLimt, gasLimitParam, "", 0, "gas limmit", notRequired)
	if err != nil {
		panic(err)
	}

	err = addPersistentString(adminCmd, adminVersion, versionParam, "", "1.0", "block chain net version", notRequired)
	if err != nil {
		panic(err)
	}

	//add sub cmd to adminCmd
	//appendSubCmd(adminCmd, "setbalance", "set target account with specified amount", setBalance, addSetBalanceFlag)
	appendSubCmd(adminCmd, "setcert", "set metering cert", setCert, addCertFlags)
	appendSubCmd(adminCmd, "validator", "add a new validator", setValidator, addSetValidatorFlags)
	appendSubCmd(adminCmd, "removecert", "remove cert from validator", removeCert, addRemoveCertFlags)
}

//admin setcert --dcname dataCenterName --certPerm certString --url https://validator-url:port
func setCert(cmd *cobra.Command, args []string) {
	client := newAnkrHttpClient(viper.GetString(adminUrl))
	opPrivateKey := viper.GetString(adminPrivateKey)
	if len(opPrivateKey) < 1 {
		fmt.Println("Invalid operator private key!")
		return
	}
	txMsg := new(metering.SetCertMsg)
	txMsg.DCName = viper.GetString(setCertDc)
	txMsg.PemBase64 = viper.GetString(setCertPerm)
	key := crypto.NewSecretKeyEd25519(viper.GetString(adminPrivateKey))
	keyAddr, _ := key.Address()
	txMsg.FromAddr = fmt.Sprintf("%X", keyAddr)
	header := getAdminMsgHeader()
	builder :=client2.NewTxMsgBuilder(*header, txMsg, serializer.NewTxSerializerCDC(), key)
	txHash, cHeight, _, err := builder.BuildAndCommit(client)
	if err != nil {
		fmt.Println("Set Cert failed.")
		fmt.Println(err)
		return
	}

	fmt.Println("Set Cert success.")
	fmt.Println("Transaction hash:",txHash)
	fmt.Println("Block Height:", cHeight)
}

func addCertFlags(cmd *cobra.Command) {
	err := addStringFlag(cmd, setCertDc, dcnameParam, "", "", "data center name", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, setCertPerm, permParam, "", "", "cert perm to be set", required)
	if err != nil {
		panic(err)
	}
}

// setvalidator --pubkey jlds --power 21 --url --privkey
func setValidator(cmd *cobra.Command, args []string) {

	client := newAnkrHttpClient(viper.GetString(adminUrl))
	opPrivateKey := viper.GetString(adminPrivateKey)
	if len(opPrivateKey) < 1 {
		fmt.Println("Invalid operator private key!")
		return
	}
	header := getAdminMsgHeader()
	validatorMsg := new(validator.ValidatorMsg)
	validatorMsg.Name = viper.GetString(setValidName)
	validatorMsg.Action = getAction(viper.GetString(setValidAction))
	validatorMsg.StakeAddress = viper.GetString(setValidStakeAddr)
	validatorMsg.StakeAmount.Cur = ankrCurrency
	amount, ok := new(big.Int).SetString(viper.GetString(setValidStakeAmount), 10)
	if !ok {
		fmt.Println("Invalid amount.")
		return
	}
	validatorMsg.StakeAmount.Value = amount.Bytes()
	validatorMsg.SetFlag = getFlagInfo(viper.GetString(setValidFlag))
	validatorMsg.ValidHeight = uint64(viper.GetInt(setValidStakeHeight))
	key := crypto.NewSecretKeyEd25519(opPrivateKey)
	keyAddr, _ := key.Address()
	validatorMsg.FromAddress = fmt.Sprintf("%X", keyAddr)
	builder := client2.NewTxMsgBuilder(*header, validatorMsg,serializer.NewTxSerializerCDC(), key)
	txHash, cHeight, _, err := builder.BuildAndCommit(client)
	if err != nil {
		fmt.Println("Set Validator failed.")
		fmt.Println(err)
		return
	}

	fmt.Println("Set Validator success.")
	fmt.Println("Transaction hash:",txHash)
	fmt.Println("Block Height:", cHeight)
}

func getFlagInfo(flag string) common.ValidatorInfoSetFlag {
	switch flag {
	case "set-name":
		return common.ValidatorInfoSetName
	case "set-val-addr":
		return common.ValidatorInfoSetValAddress
	case "set-pub":
		return common.ValidatorInfoSetPubKey
	case "set-stake-addr":
		return common.ValidatorInfoSetStakeAddress
	case "set-val-height":
		return common.ValidatorInfoSetValidHeight
	case "set-stake-amount":
		return common.ValidatorInfoSetStakeAmount
	}
	return common.ValidatorInfoSetFlag(0)
}

func addSetValidatorFlags(cmd *cobra.Command) {
	err := addStringFlag(cmd, setValidPub, pubkeyParam, "", "", "the public address of the added validator", required)
	if err != nil {
		panic(err)
	}
	//
	//err = addStringFlag(cmd, setValidPower, powerParam, "", "", "the power set to the validator", required)
	//if err != nil {
	//	panic(err)
	//}
	err = addStringFlag(cmd, setValidAction, actionParam, "", "", "update validator action", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, setValidName, nameParam, "", "", "update validator action", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, setValidFlag, flagParam, "", "", "flag of validator tansaction", notRequired)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, setValidStakeAddr, addressParam, "", "", "validator stake address", notRequired)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, setValidStakeAmount, amountParam, "", "", "validator stake amount", notRequired)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, setValidGasUsed, gasUsedParam, "", "", "gas used", notRequired)
	if err != nil {
		panic(err)
	}
	err = addIntFlag(cmd, setValidStakeHeight, heightParam, "", 0, "validator stake height", notRequired)
	if err != nil {
		panic(err)
	}
}

func getAction(action string) uint8 {
	switch action {
	case "create":
		return 1
	case "update":
		return 2
	case "remove":
		return 3
	default:
		return 0
	}
}

// removecert --pubkey string
func removeCert(cmd *cobra.Command, args []string) {
	validatorUrl = viper.GetString(adminUrl)
	client := newAnkrHttpClient(viper.GetString(adminUrl))
	amdinPriv := viper.GetString(adminPrivateKey)
	header := getAdminMsgHeader()
	txMsg := new(metering.RemoveCertMsg)
	txMsg.DCName = viper.GetString(removeCertDc)
	txMsg.NSName = viper.GetString(removeCertNs)
	key := crypto.NewSecretKeyEd25519(amdinPriv)
	keyAddr, _ := key.Address()
	txMsg.FromAddr = fmt.Sprintf("%X", keyAddr)
	builder := client2.NewTxMsgBuilder(*header, txMsg, serializer.NewTxSerializerCDC(), key)
	txHash, cHeight, _, err := builder.BuildAndCommit(client)
	if err != nil {
		fmt.Println("Remove cert failed.")
		fmt.Println(err)
		return
	}

	fmt.Println("Remove cert success.")
	fmt.Println("Transaction hash:",txHash)
	fmt.Println("Block Height:", cHeight)

}

func addRemoveCertFlags(cmd *cobra.Command) {
	err := addStringFlag(cmd, removeCertDc, dcnameParam, "", "", "name of data center name", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, removeCertNs, nameSpaceParam, "", "", "name space", required)
	if err != nil {
		panic(err)
	}
}

// get transaction header .
func getAdminMsgHeader() *client2.TxMsgHeader {
	chainId := viper.GetString(adminChId)
	gasLimit := viper.GetInt(adminGasLimt)
	gasPrice := viper.GetInt(adminGasPrice)

	header := new(client2.TxMsgHeader)
	header.Memo = viper.GetString(adminMemo)
	header.Version = viper.GetString(adminVersion)
	header.GasLimit = new(big.Int).SetUint64(uint64(gasLimit)).Bytes()
	header.GasPrice.Cur = ankrCurrency
	header.GasPrice.Value = new(big.Int).SetUint64(uint64(gasPrice)).Bytes()
	header.ChID = common.ChainID(chainId)
	return header
}
