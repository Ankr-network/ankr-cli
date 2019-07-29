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
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"os"
	"strings"
)


var(
	//flags used by query sub commands
	//persistent flag
	queryUrl = "queryUrl"

	//bind transaction flags
	trxTxid = "trxTxid"
	trxApprove = "trxApprove"

	//bind block flags
	blockHeight = "blockHeight"
	blockResultHeight = "blockResultHeight"
	validatorHeight = "validatorHeight"
	unconfirmedTxLimit = "unconfirmedTxLimit"

	//transaction prefix
     TxPrefix = "trx_send="
     setMeteringPrefix = "set_mtr="
     setBalancePrefix = "set_bal="
     setStakePrefix = "set_stk="
     setCertPrefix = "set_crt="
     removeCertPrefix = "rmv_crt="
     setValidatorPrefix = "val:"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "A brief description of your command",
}

func init() {
	err := addPersistentString(queryCmd, queryUrl, "url", "", "", "validator url", required)
	if err != nil {
		panic(err)
	}
	appendSubCmd(queryCmd, "transaction","transaction allows you to query the transaction results.", transactionInfo, addTransactionInfoFlags)
	appendSubCmd(queryCmd, "block", "Get block at a given height. If no height is provided, it will fetch the latest block.",
		queryBlock, addQueryBlockFlags)
	appendSubCmd(queryCmd, "block-result", "BlockResults gets ABCIResults at a given height. If no height is provided, it will fetch results for the latest block.",
		queryBlockResult, addQueryBlockResultFlags)
	appendSubCmd(queryCmd, "validators", "Get the validator set at the given block height. If no height is provided, it will fetch the current validator set.",
		queryValidator, addQueryValidatorFlags)
	appendSubCmd(queryCmd, "status", "Get Ankr status including node info, pubkey, latest block hash, app hash, block height and time.",
		queryStatus, nil)
	appendSubCmd(queryCmd, "netinfo", "Get network info.", queryNetInfo, nil)
	appendSubCmd(queryCmd, "genesis", "Get genesis file.", queryGenesis, nil)
	appendSubCmd(queryCmd, "consensusstate", "ConsensusState returns a concise summary of the consensus state", queryConsensusState, nil)
	appendSubCmd(queryCmd, "dumpconsensusstate", "dumps consensus state", queryDumpConsensusState, nil)
	appendSubCmd(queryCmd, "unconfirmedtxs", "Get unconfirmed transactions (maximum ?limit entries) including their number",
		queryUnconfirmedTxs, addQueryUncofirmedTxsFlags)
	appendSubCmd(queryCmd, "numunconfirmedtxs","Get number of unconfirmed transactions.", queryNumUnconfiredTxs, nil)
}

func transactionInfo(cmd *cobra.Command, args []string)  {
	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	txid := viper.GetString(trxTxid)
	if len(txid) < 1 {
		fmt.Println("Invalid transaction hash")
		return
	}
	if 	strings.HasPrefix(txid, "0x") {
		txid = strings.TrimLeft(txid, "0x")
	}
	//txid = txid[2:]
	approve := viper.GetBool(trxApprove)
	hash, err := hex.DecodeString(txid)
	if err != nil {
		fmt.Println("Invalid transaction id format!")
		return
	}
	resp, err := cl.Tx(hash, approve )
	if err != nil {
		fmt.Println("Failed to query transaction.")
		fmt.Println(err)
		return
	}
	fmt.Println("testing.....")
	result := parseTx(resp)
	displayTx(result)
	//display(result)
}

func addTransactionInfoFlags(cmd *cobra.Command)  {
	err := addStringFlag(cmd, trxTxid, txidFlag, "", "", "The transaction hash", required)
	if err != nil {
		panic(err)
	}
	err = addBoolFlag(cmd, trxApprove, approveFlag, "", false, "Include a proof of the transaction inclusion in the block", "")
	if err != nil {
		panic(err)
	}
}

//query block
func queryBlock(cmd *cobra.Command, args []string)  {
	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	height := viper.GetInt64(blockHeight)
	heightP := &height
	if height <= 0 {
		heightP = nil
	}
	resp, err :=cl.Block(heightP)
	if err != nil {
		fmt.Println("Query block failed.", err)
		return
	}
	fmt.Println("\nBlock Head:")
	display(resp.BlockMeta)
	fmt.Println("\nTransactions contained in block:")
	if resp.Block.Txs == nil || len(resp.Block.Txs) == 0 {
		fmt.Println("[]")
	}else {
		for _, tx := range resp.Block.Txs {
			fmt.Println(string(tx))
		}
	}
}
func addQueryBlockFlags(cmd *cobra.Command)  {
	err := addInt64Flag(cmd, blockHeight, heightFlag, "", -1, "height of the block to query", "" )
	if err != nil {
		panic(err)
	}
}

//query blockresult
func queryBlockResult(cmd *cobra.Command, args []string)  {
	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	height := viper.GetInt64(blockHeight)
	heightP := &height
	if height <= 0 {
		heightP = nil
	}
	resp, err := cl.BlockResults(heightP)
	if err != nil {
		fmt.Println("Query block result failed.", err)
		return
	}
	display(resp)
}
func addQueryBlockResultFlags(cmd *cobra.Command)  {
	err := addInt64Flag(cmd, blockResultHeight, heightFlag, "", -1, "block height", "")
	if err != nil {
		panic(err)
	}
}

//query validator
func queryValidator(cmd *cobra.Command, args []string)  {
	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	height := viper.GetInt64(validatorHeight)
	heightP := &height
	if height <= 0 {
		heightP = nil
	}
	resp, err := cl.Validators(heightP)
	if err != nil {
		fmt.Println("Query validators failed.", err)
		return
	}
	display(resp)
}
func addQueryValidatorFlags(cmd *cobra.Command)  {
	err := addInt64Flag(cmd, validatorHeight, heightFlag, "", -1, "block height", "")
	if err != nil {
		panic(err)
	}
}

//query status
func queryStatus(cmd *cobra.Command, args []string)  {
	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	resp, err := cl.Status()
	if err != nil{
		fmt.Println("Query status failed.",err)
		return
	}
	display(resp)
}

//query netinfo
func queryNetInfo(cmd *cobra.Command, args []string)  {
	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	resp, err := cl.NetInfo()
	if err != nil {
		fmt.Println("Query net info faild.",err)
		return
	}
	display(resp)
}

//query genesis
func queryGenesis(cmd *cobra.Command, args []string)  {
	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	resp, err := cl.Genesis()
	if err != nil {
		fmt.Println("Query genesis failed.", err)
		return
	}
	display(resp)
}

//query consensus state
func queryConsensusState(cmd *cobra.Command, args []string)  {
	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	resp, err := cl.ConsensusState()
	if err != nil {
		fmt.Println("Query consensus state failed.", err)
		return
	}
	display(resp)
}

//query dump consensus state
func queryDumpConsensusState(cmd *cobra.Command, args []string)  {
	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	resp, err := cl.DumpConsensusState()
	if err != nil {
		fmt.Println("Query dump consensus state failed.", err)
		return
	}
	display(resp)
}

//query unconfirmed transactions
func queryUnconfirmedTxs(cmd *cobra.Command, args []string)  {
	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	limmit := viper.GetInt(unconfirmedTxLimit)
	resp, err := cl.UnconfirmedTxs(limmit)
	if err != nil {
		fmt.Println("Query unconfirmed transactions failed.", err)
		return
	}
	display(resp)
}

func addQueryUncofirmedTxsFlags(cmd *cobra.Command)  {
	err := addIntFlag(cmd, unconfirmedTxLimit, limitFlag, "",30, "number of entries", "")
	if err != nil {
		panic(err)
	}
}

//query number of unconfirmed transactions
func queryNumUnconfiredTxs(cmd *cobra.Command, args []string)  {
	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	resp, err := cl.NumUnconfirmedTxs()
	if err != nil {
		fmt.Println("Query number of unconfirmed transactions failed.", err)
		return
	}
	display(resp)
}

//transaction data structure
type Transaction struct {
	Type string
	Hash string
	From string
	To string
	Nonce string
	Amount string
}


//transaction data structure used in parsing all kinds of transactions
type ResultTx struct {
	Type string `json:"type"`
	Hash     string   `json:"hash"`//common.HexBytes           `json:"hash"`
	Height   int64                  `json:"height"` //block height
	Index    uint32                 `json:"index"` //transaction index in block
	Data map[string] string `json:"data"` //used to store different type of transaction data
}

//parse transaction data from rpc response and write to ResultTx
func (rt *ResultTx)parseSendTx(tx string)  {
	txString := strings.TrimPrefix(tx, TxPrefix)
	txSlice := strings.Split(txString, ":")
	rt.Data["from"] = string(txSlice[0])
	rt.Data["to"] = string(txSlice[1])
	rt.Data["amount"] = string(txSlice[2])
	rt.Data["nonce"] = string(txSlice[3])
}

func (rt *ResultTx)parseSetMeteringTx(tx string)  {
	txString := strings.TrimPrefix(tx, setMeteringPrefix)
	txSlice := strings.Split(txString, ":")
	rt.Data["dc name"] = string(txSlice[0])
	rt.Data["name space"] = string(txSlice[1])
	rt.Data["nonce"] = string(txSlice[4])
	rt.Data["value"] = string(txSlice[5])
}

func (rt *ResultTx)parseSetBalanceTx(tx string)  {
	txString := strings.TrimPrefix(tx, setBalancePrefix)
	txSlice := strings.Split(txString, ":")
	rt.Data["address"] = string(txSlice[0])
	rt.Data["amount"] = string(txSlice[1])
	rt.Data["nonce"] = string(txSlice[2])
}
func (rt *ResultTx)parseSetStakeTx(tx string)  {
	txString := strings.TrimPrefix(tx, setStakePrefix)
	txSlice := strings.Split(txString, ":")
	rt.Data["amount"] = string(txSlice[0])
}

func (rt *ResultTx)parseSetCertTx(tx string)  {
	txString := strings.TrimPrefix(tx, setCertPrefix)
	txSlice := strings.Split(txString, ":")
	rt.Data["dc name"] = string(txSlice[0])
	rt.Data["cert perm"] = string(txSlice[1])
}

func (rt *ResultTx)parseRemoveCertTx(tx string)  {
	txString := strings.TrimPrefix(tx, removeCertPrefix)
	txSlice := strings.Split(txString, ":")
	rt.Data["dc name"] = string(txSlice[0])
}

func (rt *ResultTx)parseSetValidatorTx(tx string)  {
	txString := strings.TrimPrefix(tx, setValidatorPrefix)
	txSlice := strings.Split(txString, ":")
	rt.Data["public key"] = string(txSlice[0])
	rt.Data["power"] = string(txSlice[1])
}

//parse transaction into ResultTx struct
func parseTx(tx *core_types.ResultTx) ResultTx {
	var result ResultTx
	result.Data = make(map[string]string)
	result.Hash = fmt.Sprintf("0x%x",tx.Hash)
	result.Height = tx.Height
	result.Index = tx.Index
	txString := string(tx.Tx)
	if strings.HasPrefix(txString, TxPrefix) {
		result.Type = "transfer"
		result.parseSendTx(txString)
		return  result
	}else if strings.HasPrefix(txString, setMeteringPrefix) {
		result.Type = "set metering"
		result.parseSetMeteringTx(txString)
		return result
	}else if strings.HasPrefix(txString, setBalancePrefix) {
		result.Type = "set balance"
		result.parseSetBalanceTx(txString)
		return result
	}else if strings.HasPrefix(txString, setStakePrefix) {
		result.Type = "set stake"
		result.parseSetStakeTx(txString)
		return result
	}else if strings.HasPrefix(txString, setCertPrefix) {
		result.Type = "set cert"
		result.parseSetCertTx(txString)
		return result
	}else if strings.HasPrefix(txString, removeCertPrefix) {
		result.Type = "remove cert"
		result.parseRemoveCertTx(txString)
		return result
	}else if strings.HasPrefix(txString, setValidatorPrefix) {
		result.Type = "set validator"
		result.parseSetValidatorTx(txString)
		return result
	}else {
		fmt.Println("Can not parse Transaction data:", string(tx.Tx))
		return result
	}
}

//display transaction information
func displayTx(rt ResultTx)  {
	w := newTabWriter(os.Stdout)
	fmt.Fprintf(w, "tx type\thash\tblock height\tblock index\tdetail\n")
	fmt.Fprintf(w, "%s\t%s\t%d\t%d\t", rt.Type, rt.Hash, rt.Height, rt.Index)
	switch rt.Type {
	case "transfer":
		fmt.Fprintf(w, "from: %s\tto:%s\tamount:%s\tnonce:%s\n:",rt.Data["from"],rt.Data["to"],rt.Data["amount"],rt.Data["nonce"])
	case "":
		fmt.Fprintf(w, "from: %s\tto%s\tamount:%s\tnonce:%s\n",rt.Data["from"],rt.Data["to"],rt.Data["amount"],rt.Data["nonce"])
	}
	w.Flush()
}