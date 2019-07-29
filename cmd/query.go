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
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
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
	result := parseSendTx(resp)
	display(result)
	//x, _ := new(big.Int).SetString("xxx", 10)
	//parseSendTx(resp)
	//time.Unix(x.Int64(),0)
	//display(resp)
	//fmt.Println(string(resp.Tx))
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
		//display(resp.Block.Txs)
		//txs, err := decodeTxs(resp.Block.Txs)
		//if err != nil {
		//	fmt.Println("base64 decode error!")
		//	return
		//}else {
		//	fmt.Println(txs)
		//}
		for _, tx := range resp.Block.Txs {
			fmt.Printf("%x", tx.Hash())
			//fmt.Println(string(tx.Hash()))
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

func decodeTxs(txs []types.Tx) ([][]byte, error) {
	decodeTxs := make([][]byte,0, len(txs))
	for _, tx := range txs{
		//var dst []byte
		//_, err := base64.StdEncoding.Decode(dst, tx)
		base64.StdEncoding.EncodeToString(tx)
		dst, err :=base64.StdEncoding.DecodeString(string(tx))
		if err != nil {
			return decodeTxs, err
		}
		decodeTxs = append(decodeTxs, dst)
	}
	return decodeTxs, nil
}

type Transaction struct {
	Type string
	Hash string
	From string
	To string
	Nonce string
	Amount string
}


//send coin transaction
type ResultTx struct {
	Type string
	Hash     common.HexBytes           `json:"hash"`
	Height   int64                  `json:"height"`
	Index    uint32                 `json:"index"`
	Data interface{} `json:"data"`
}


//transaction send type
type trxSend struct {
	From string `json:"from"`
	To string `json:"to"`
	Amount string `json:"amount"`
	Nonce string `json:"nonce"`
}

func parseSendTx(tx *core_types.ResultTx) ResultTx {
	var result ResultTx
	result.Hash = tx.Hash
	result.Height = tx.Height
	result.Index = tx.Index
	txString := string(tx.Tx)
	if strings.HasPrefix(txString, TxPrefix) {
		result.Type = TxPrefix
		var trx trxSend
		txString = strings.TrimPrefix(txString, TxPrefix)
		txSlice := strings.Split(txString, ":")
		trx.From = string(txSlice[0])
		trx.To = string(txSlice[1])
		trx.Amount = string(txSlice[2])
		trx.Nonce = string(txSlice[3])
		result.Data = trx
	}
	return result
}
