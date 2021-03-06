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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ankr-network/dccn-common/wallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "account is used to generate new accounts, encrypt privatekey or decrypt privatekey from keystore",
}

//names of sub command bind in viper, which is used to bind flags
//naming notions subCmdNameKey. eg. "account" have a flag named "url" shall named as accountUrl,
var (
	genAccNumber   = "genAccNumber"
	genAccOutput   = "agenAccOutput"
	getUrl         = "getUrl"
	getAddress     = "getAddress"
	genkeyPrivkey  = "genkeyPrivkey"
	genkeyOutput   = "genkeyOutput"
	exportKeystore = "exportKeystore"
	resetKeystore  = "resetKeystore"
)

func init() {
	appendSubCmd(accountCmd, "genaccount", "generate new account.", generateAccounts, addGenAccountFlags)
	appendSubCmd(accountCmd, "genkeystore", "generate keystore file based on private key and user input password.", genKeystore, addGenkeystoreFlags)
	appendSubCmd(accountCmd, "getbalance", "get the balance of an address.", getBalance, addGetBalanceFlags)
	appendSubCmd(accountCmd, "exportprivatekey", "recover private key from keystore.", exportPrivatekey, addExportFlags)
	appendSubCmd(accountCmd,"resetpwd", "reset keystore password.", resetPwd, addResetPWDFlags)
}

type ExeCmd struct {
	Name     string
	Short    string
	Long     string
	Exec     func(cmd *cobra.Command, args []string)
	FlagFunc func(cmd *cobra.Command)
}

type Account struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Address    string `json:"address"`
}

//account genaccount functions
func addGenAccountFlags(cmd *cobra.Command) {
	err := addIntFlag(cmd, genAccNumber, numberAccountFlag, "n", 1, "number of accounts to be generated", "")
	if err != nil {
		panic(err)
	}
	//genAccountCmd.Flags().IntP(numberAccountFlag, "n", 1, "number of accounts to be generated")
	//viper.BindPFlag(numberAccountFlag, genAccountCmd.Flags().Lookup(numberAccountFlag))

	err = addStringFlag(cmd, genAccOutput, outputFlag, "o", "", "output account to file", "")
	if err != nil {
		panic(err)
	}
}

//generate new account, encrypt private key to keystore base on user input password
func generateAccounts(cmd *cobra.Command, args []string) {
	fmt.Println(`please record and backup keystore once it is generated, we don’t store your private key!`)
	fmt.Println("\ngenerating accounts...")
	numberAccount := viper.GetInt(genAccNumber)
	for i := 0; i < numberAccount; i++ {
		//generate single Account
		//input password from terminal
		acc := generateAccount()
		s := fmt.Sprintf("\nAccount_%d", i)
		fmt.Println(s)
		fmt.Println("private key: ", acc.PrivateKey, "\npublic key: ", acc.PublicKey, "\naddress: ", acc.Address)
		path := viper.GetString(genAccOutput)
		if path == "" {
			path = configHome()
		}
		generateKeystore(acc, path)
	}
}

func generateAccount() Account {
	priv, pub, addr := wallet.GenerateKeys()
	return Account{priv, pub, addr}
}

//generate keystore based account and password
func generateKeystore(acc Account, path string) error {
	fmt.Println("\nabout to export to keystore.. ")

InputPassword:
	fmt.Print("please input the keystore encryption password:")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil
	}

	fmt.Print("\nplease input password again: ")
	confirmPassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}

	if string(password) != string(confirmPassword) {
		fmt.Println("\nError:password and confirm password not match!")
		goto InputPassword
		//return errors.New("\npassword and confirm password not match")
	}

	cryptoStruct, err := EncryptDataV3([]byte(acc.PrivateKey), []byte(password), StandardScryptN, StandardScryptP)
	if err != nil {
		return err
	}
	//_ := cryptoStruct

	encryptedKeyJSONV3 := EncryptedKeyJSONV3{
		Address:        acc.Address,
		PublicKey:      acc.PublicKey,
		Crypto:         cryptoStruct,
		KeyJSONVersion: keyJSONVersion,
	}
	jsonKey, err := json.Marshal(encryptedKeyJSONV3)
	if err != nil {
		return err
	}

	fmt.Println("\n\nexporting to keystore...")
	ts := time.Now().UTC()

	kfw, err := KeyFileWriter(path, fmt.Sprintf("UTC--%s--%s", toISO8601(ts), acc.Address))
	if err != nil {
		return err
	}

	defer kfw.Close()

	_, err = kfw.Write(jsonKey)
	if err != nil {
		return errors.New("unable to write keystore")
	}

	fmt.Printf("\ncreated keystore: %s/UTC--%s--%s\n\n", path, toISO8601(ts), acc.Address)
	return nil
}

//genkeystore functions
func addGenkeystoreFlags(cmd *cobra.Command) {
	err := addStringFlag(cmd, genkeyPrivkey, privkeyFlag, "p", "", "private key of an account.", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, genkeyOutput, outputFlag, "o", "", "output file path.", "")
	if err != nil {
		panic(err)
	}
}

func genKeystore(cmd *cobra.Command, args []string) {
	privateKey := viper.GetString(genkeyPrivkey)
	if len(privateKey) == 0 {
		fmt.Println("invalid private key")
		return
	}
	acc, err := getAccountFromPrivatekey(privateKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	path := viper.GetString(genkeyOutput)
	if path == "" {
		path = configHome()
	}
	err = generateKeystore(acc, path)
	if err != nil {
		fmt.Println(err)
	}
}

//get balance functions
func addGetBalanceFlags(cmd *cobra.Command) {
	err := addStringFlag(cmd, getAddress, addressFlag, "a", "", "the address of an account.", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, getUrl, urlFlag, "", "", "the url with an endpoint of an ankr chain validator.", required)
	if err != nil {
		panic(err)
	}
}

func getBalance(cmd *cobra.Command, args []string) {
	httpsUrl := viper.GetString(getUrl)
	address := viper.GetString(getAddress)
	index := strings.LastIndex(httpsUrl, ":")
	if index < 0 {
		fmt.Println("url is not correct, example 'https://chain-01.dccn.ankr.com:443'")
		return
	}

	balance, err := wallet.GetBalance(httpsUrl[:index], httpsUrl[index+1:], address)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(balance) <= 18 {
		balanceDecimalZero := make([]byte, 18-len(balance))
		for i := 0; i < 18-len(balance); i++ {
			balanceDecimalZero = append(balanceDecimalZero, '0')
		}
		balance = "0." + string(balanceDecimalZero) + balance
	} else {
		balance = balance[:len(balance)-18] + "." + balance[len(balance)-18:]
	}
	fmt.Printf("The balance is: %s\n", balance)
}

func addExportFlags(cmd *cobra.Command) {
	err := addStringFlag(cmd, exportKeystore, fileFlag, "f", "", "the path where keystore file is located.", required)
	if err != nil {
		panic(err)
	}
}

//generate private key from keystore and password
func exportPrivatekey(cmd *cobra.Command, args []string) {
	ksf := viper.GetString(exportKeystore)
	privateKey := decryptPrivatekey(ksf)

	if privateKey == "" {
		fmt.Println("Empty privateKey!!")
		return
	}
	fmt.Println("\nPrivate key exported:", privateKey)
}

//decrypt private key from keystore and user input password
func decryptPrivatekey(file string) string {
	ks, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var key EncryptedKeyJSONV3

	err = json.Unmarshal(ks, &key)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Print("\nPlease input the keystore password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	privateKeyDecrypt, err := DecryptDataV3(key.Crypto, string(password))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	privateKey := string(privateKeyDecrypt)
	return privateKey
}

func resetPwd(cmd *cobra.Command, args []string) {
	ksf := viper.GetString(resetKeystore)
	privateKey := decryptPrivatekey(ksf)

	if privateKey == "" {
		fmt.Println("Empty privateKey!!")
		return
	}

	acc, err := getAccountFromPrivatekey(privateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	pwd := readPassword()

	cryptoStruct, err := EncryptDataV3([]byte(acc.PrivateKey), []byte(pwd), StandardScryptN, StandardScryptP)
	if err != nil {
		panic(err)
	}
	//_ := cryptoStruct

	encryptedKeyJSONV3 := EncryptedKeyJSONV3{
		Address:        acc.Address,
		PublicKey:      acc.PublicKey,
		Crypto:         cryptoStruct,
		KeyJSONVersion: keyJSONVersion,
	}
	jsonKey, err := json.Marshal(encryptedKeyJSONV3)
	if err != nil {
		panic(err)
	}

	kfw, err := os.OpenFile(ksf, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("\nUnable to open file:", ksf)
		fmt.Println(err)
		return
	}
	defer kfw.Close()
	_, err = kfw.Write(jsonKey)
	if err != nil {
		fmt.Println("\nUnable to write keystore")
		return
	}

	fmt.Println("\nPassword reset success.")
}

func addResetPWDFlags(cmd *cobra.Command){
	err := addStringFlag(cmd, resetKeystore, fileFlag, "f", "", "the path where keystore file is located.", required)
	if err != nil {
		panic(err)
	}
}

//read password from terminal
func readPassword() []byte {
InputPassword:
	fmt.Print("please input the keystore encryption password:")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}

	fmt.Print("\nplease input password again: ")
	confirmPassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}

	if string(password) != string(confirmPassword) {
		fmt.Println("\nError:password and confirm password not match!")
		goto InputPassword
		//return errors.New("\npassword and confirm password not match")
	}
	return password
}