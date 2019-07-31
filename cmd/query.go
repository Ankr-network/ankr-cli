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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	//todo replace with dccn-tendermint
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"os"
	"strings"
	"text/tabwriter"
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
	Short: "query information from ankr chain",
}

func init() {
	err := addPersistentString(queryCmd, queryUrl, "url", "", "", "validator url", required)
	if err != nil {
		panic(err)
	}
	appendSubCmd(queryCmd, "transaction","transaction allows you to query the transaction results.", transactionInfo, addTransactionInfoFlags)
	appendSubCmd(queryCmd, "block", "Get block at a given height. If no height is provided, it will fetch the latest block.",
		queryBlock, addQueryBlockFlags)
	//deprecated
	//appendSubCmd(queryCmd, "blockresult", "BlockResults gets ABCIResults at a given height. If no height is provided, it will fetch results for the latest block.",
	//	queryBlockResult, addQueryBlockResultFlags)
	appendSubCmd(queryCmd, "validators", "Get the validator set at the given block height. If no height is provided, it will fetch the current validator set.",
		queryValidator, addQueryValidatorFlags)
	appendSubCmd(queryCmd, "status", "Get Ankr status including node info, pubkey, latest block hash, app hash, block height and time.",
		queryStatus, nil)
	//appendSubCmd(queryCmd, "netinfo", "Get network info.", queryNetInfo, nil)
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
	result := parseTx(resp)
	w := newTabWriter(os.Stdout)
	fmt.Fprintf(w, "tx type\thash\tblock height\tblock index\tdetail\n")
	displayTx(result, w)
	w.Flush()
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
	w := newTabWriter(os.Stdout)
	fmt.Fprintf(w, "\nBlock info:\n")
	outPutHeader(w, resp.Block.Header)
	fmt.Fprintf(w, "\nTransactions contained in block: \n")
	if resp.Block.Txs == nil || len(resp.Block.Txs) == 0 {
		fmt.Fprintf(w, "[]\n")
	}else{
		outPutTransactions(w, resp.Block.Txs)
	}
	w.Flush()
}
func addQueryBlockFlags(cmd *cobra.Command)  {
	err := addInt64Flag(cmd, blockHeight, heightFlag, "", -1, "height of the block to query", "" )
	if err != nil {
		panic(err)
	}
}

//query blockresult
//func queryBlockResult(cmd *cobra.Command, args []string)  {
//	validatorUrl = viper.GetString(queryUrl)
//	if len(validatorUrl) < 1 {
//		fmt.Println("Illegal url is received!")
//		return
//	}
//	cl := newAnkrHttpClient(validatorUrl)
//	height := viper.GetInt64(blockHeight)
//	heightP := &height
//	if height <= 0 {
//		heightP = nil
//	}
//	resp, err := cl.BlockResults(heightP)
//	if err != nil {
//		fmt.Println("Query block result failed.", err)
//		return
//	}
//	display(resp)
//
//}
//func addQueryBlockResultFlags(cmd *cobra.Command)  {
//	err := addInt64Flag(cmd, blockResultHeight, heightFlag, "", -1, "block height", "")
//	if err != nil {
//		panic(err)
//	}
//}

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
	//display(resp)
	w := newTabWriter(os.Stdout)
	outPutValidator(w, resp)
	w.Flush()

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
	w := newTabWriter(os.Stdout)
	outputStatus(w, resp)
	w.Flush()
}

func outputStatus(w *tabwriter.Writer, st *core_types.ResultStatus)  {
	b, _ := json.MarshalIndent(st.NodeInfo, "", "    ")
	fmt.Fprintf(w,"node_info:%s\n", string(b))
	b, _ = json.MarshalIndent(st.SyncInfo, "", "    ")
	fmt.Fprintf(w, "sync_info:%s\n", string(b))
	fmt.Fprintf(w, "validator_info:{ \n", )
	fmt.Fprintf(w, "\t\"address\":%s\n", st.ValidatorInfo.Address )
	fmt.Fprintf(w, "\t\"pub_key\":%s\n", base64.StdEncoding.EncodeToString(st.ValidatorInfo.PubKey.Bytes()) )
	fmt.Fprintf(w, "\t\"voting_power\":%d\n}\n", st.ValidatorInfo.VotingPower )
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
	//display(resp)
	w := newTabWriter(os.Stdout)
	outputGenesis(w, resp)
	w.Flush()
}

//format output result genesis
func outputGenesis(w *tabwriter.Writer,resultGenesis * core_types.ResultGenesis )  {
	fmt.Fprintf(w, "Genesis:{\n")
	fmt.Fprintf(w, "\t\"genesis_time\": %v\n", resultGenesis.Genesis.GenesisTime)
	fmt.Fprintf(w, "\t\"chain_id\": %v\n", resultGenesis.Genesis.ChainID)
	b, _ := json.MarshalIndent(resultGenesis.Genesis.ConsensusParams, "", "\t")
	fmt.Fprintf(w, "\t\"consensus_params\":%s\n",string(b))
	fmt.Fprintf(w, "\t\"validators\":{\n")
	fmt.Fprintf(w, "\t\taddress\tpub_key\tpower\tname\n")
	for _, valid := range resultGenesis.Genesis.Validators {
		fmt.Fprintf(w, "\t\t%v\t%v\t%v\t%v\n", valid.Address, base64.StdEncoding.EncodeToString(valid.PubKey.Bytes()), valid.Power, valid.Name)
	}
	fmt.Fprintf(w, "\t}\n")
	fmt.Fprintf(w, "}\n")
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
	//todo resp.Txs[0 - n-1]
	//display(resp)
	w := newTabWriter(os.Stdout)
	outputTxResult(w, resp)
	w.Flush()
}

func outputTxResult(w *tabwriter.Writer, txResult *core_types.ResultUnconfirmedTxs)  {
	fmt.Fprintf(w, "n_tx: %d\n", txResult.Count)
	fmt.Fprintf(w, "total: %d\n", txResult.Total)
	fmt.Fprintf(w, "total_bytes: %d\n", txResult.TotalBytes)
	fmt.Fprintf(w, "transactions:\n")
	if len(txResult.Txs) == 0 {
		fmt.Fprintf(w, "[]\n")
	}else {
		outPutTransactions(w, txResult.Txs)
	}
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
	//todo resp.Txs
	//display(resp)
	w := newTabWriter(os.Stdout)
	outputTxResult(w, resp)
	w.Flush()
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
func parseSendTx(tx string) map[string]string {
	data := make(map[string]string)
	txString := strings.TrimPrefix(tx, TxPrefix)
	txSlice := strings.Split(txString, ":")
	data["from"] = string(txSlice[0])
	data["to"] = string(txSlice[1])
	data["amount"] = string(txSlice[2])
	data["nonce"] = string(txSlice[3])
	return data
}

func parseSetMeteringTx(tx string)  map[string]string{
	data := make(map[string]string)
	txString := strings.TrimPrefix(tx, setMeteringPrefix)
	txSlice := strings.Split(txString, ":")
	data["dc name"] = string(txSlice[0])
	data["name space"] = string(txSlice[1])
	//data["nonce"] = string(txSlice[4])
	data["value"] = string(txSlice[5])
	return data
}

func parseSetBalanceTx(tx string)  map[string]string{
	data := make(map[string]string)
	txString := strings.TrimPrefix(tx, setBalancePrefix)
	txSlice := strings.Split(txString, ":")
	data["address"] = string(txSlice[0])
	data["amount"] = string(txSlice[1])
	//data["nonce"] = string(txSlice[2])
	return data
}
func parseSetStakeTx(tx string) map[string]string {
	data := make(map[string]string)
	txString := strings.TrimPrefix(tx, setStakePrefix)
	txSlice := strings.Split(txString, ":")
	data["amount"] = string(txSlice[0])
	return data
}

func parseSetCertTx(tx string) map[string]string {
	data := make(map[string]string)
	txString := strings.TrimPrefix(tx, setCertPrefix)
	txSlice := strings.Split(txString, ":")
	data["dc name"] = string(txSlice[0])
	data["cert perm"] = string(txSlice[1])
	return data
}

func parseRemoveCertTx(tx string) map[string]string {
	data := make(map[string]string)
	txString := strings.TrimPrefix(tx, removeCertPrefix)
	txSlice := strings.Split(txString, ":")
	data["dc name"] = string(txSlice[0])
	return data
}

func parseSetValidatorTx(tx string)map[string]string  {
	data := make(map[string]string)
	txString := strings.TrimPrefix(tx, setValidatorPrefix)
	txSlice := strings.Split(txString, ":")
	data["public key"] = string(txSlice[0])
	data["power"] = string(txSlice[1])
	return data
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
		result.Data = parseSendTx(txString)
		return  result
	}else if strings.HasPrefix(txString, setMeteringPrefix) {
		result.Type = "set metering"
		result.Data = parseSetMeteringTx(txString)
		return result
	}else if strings.HasPrefix(txString, setBalancePrefix) {
		result.Type = "set balance"
		result.Data = parseSetBalanceTx(txString)
		return result
	}else if strings.HasPrefix(txString, setStakePrefix) {
		result.Type = "set stake"
		result.Data = parseSetStakeTx(txString)
		return result
	}else if strings.HasPrefix(txString, setCertPrefix) {
		result.Type = "set cert"
		result.Data = parseSetCertTx(txString)
		return result
	}else if strings.HasPrefix(txString, removeCertPrefix) {
		result.Type = "remove cert"
		result.Data = parseRemoveCertTx(txString)
		return result
	}else if strings.HasPrefix(txString, setValidatorPrefix) {
		result.Type = "set validator"
		result.Data = parseSetValidatorTx(txString)
		return result
	}else {
		fmt.Println("Can not parse Transaction data:", string(tx.Tx))
		return result
	}
}

//display transaction information
func displayTx(rt ResultTx, w *tabwriter.Writer)  {
	//table header
	fmt.Fprintf(w, "%s\t%s\t%d\t%d\t", rt.Type, rt.Hash, rt.Height, rt.Index)
	//table contents
	switch rt.Type {
	case "transfer":
		fmt.Fprintf(w, "from: %s\tto:%s\tamount:%s\tnonce:%s\n:",rt.Data["from"],rt.Data["to"],rt.Data["amount"],rt.Data["nonce"])
	case "set metering":
		fmt.Fprintf(w,"dc-name:%s\tname-space:%s\tvalue:%s\n",rt.Data["dc name"],rt.Data["name space"],rt.Data["value"])
	case "set balance":
		fmt.Fprintf(w,"address:%s\tamount:%s\n",rt.Data["address"],rt.Data["amount"])
	case "set stake":
		fmt.Fprintf(w,"tamount:%s\n", rt.Data["amount"])
	case "set cert":
		fmt.Fprintf(w,"dc name:%s\tcert perm:%s\n", rt.Data["dc name"],rt.Data["cert perm"])
	case "remove cert":
		fmt.Fprintf(w,"tdc name:%s\n",rt.Data["dc name"])
	case "set validator":
		fmt.Fprintf(w,"public key:%s\tpower:%s\n", rt.Data["public key"],rt.Data["power"])
	default :
		fmt.Fprintf(w, "unrecognized transaction!",)
	}
}

func outPutTransactions(w *tabwriter.Writer,txs types.Txs)  {
	fmt.Fprintf(w, "tx type\thash\tdetail\n")
	for _, tx := range txs {
		hash := fmt.Sprintf("0x%x",tx.Hash())
		txString := string(tx)
		if strings.HasPrefix(txString, TxPrefix) {
			data := parseSendTx(txString)
			fmt.Fprintf(w,"transfer\t%s\tfrom:%s\tto:%s\tamount:%s\tnonce:%s\n",hash,data["from"],data["to"],data["amount"],data["nonce"])
		}else if strings.HasPrefix(txString, setMeteringPrefix) {
			data := parseSetMeteringTx(txString)
			fmt.Fprintf(w,"set metering\t%s\tdc-name:%s\tname-space:%s\tvalue:%s\n",hash,data["dc name"],data["name space"],data["value"])
		}else if strings.HasPrefix(txString, setBalancePrefix) {
			data := parseSetBalanceTx(txString)
			fmt.Fprintf(w,"set balance\t%s\taddress:%s\tamount:%s\n",hash,data["address"],data["amount"])
		}else if strings.HasPrefix(txString, setStakePrefix) {
			data := parseSetStakeTx(txString)
			fmt.Fprintf(w,"set stake\t%s\tamount:%s\n",hash, data["amount"])
		}else if strings.HasPrefix(txString, setCertPrefix) {
			data := parseSetCertTx(txString)
			fmt.Fprintf(w,"set cert\t%s\tdc name:%s\tcert perm:%s\n",hash,data["dc name"],data["cert perm"])
		}else if strings.HasPrefix(txString, removeCertPrefix) {
			data := parseRemoveCertTx(txString)
			fmt.Fprintf(w,"remove cert\t%s\tdc name:%s\n",hash,data["dc name"])
		}else if strings.HasPrefix(txString, setValidatorPrefix) {
			data := parseSetValidatorTx(txString)
			fmt.Fprintf(w,"set validator\t%s\tpublic key:%s\tpower:%s\n",hash,data["public key"],data["power"])
		}else {
			fmt.Fprintf(w,"unrecognized  transaction %s ", txString)
		}
	}
}
func outPutHeader(w *tabwriter.Writer, header types.Header)  {
	//information to be displayed in the window
	fmt.Fprintf(w,"Version: %v\n", header.Version)
	fmt.Fprintf(w,"Chain-Id: %v\n", header.ChainID)
	fmt.Fprintf(w,"Height: %v\n", header.Height)
	fmt.Fprintf(w,"Time: %v\n", header.Time)
	fmt.Fprintf(w,"Number-Txs: %v\n", header.NumTxs)
	fmt.Fprintf(w,"Total-Txs: %v\n", header.TotalTxs)
	fmt.Fprintf(w,"Last-block-id: %v\n", header.LastBlockID)
	fmt.Fprintf(w, "Last-commit-hash:%v\n",header.LastCommitHash)
	fmt.Fprintf(w,"Data-hash: %v\n", header.DataHash)
	fmt.Fprintf(w,"Validator: %v\n", header.ValidatorsHash)
	fmt.Fprintf(w,"Consensus: %v\n", header.ConsensusHash)
	fmt.Fprintf(w,"Version: %v\n", header.Version)
	fmt.Fprintf(w,"App-hash: %v\n", header.AppHash)
	fmt.Fprintf(w,"Proposer-Address: %v\n", header.ProposerAddress)
}

func outPutValidator(w *tabwriter.Writer, validatorResult *core_types.ResultValidators){
	fmt.Fprintf(w, "Height:%d\n",validatorResult.BlockHeight)
	fmt.Fprintf(w, "\nValidators information: \n")
	fmt.Fprintf(w, "Address\tPubkey\tVoting-Power\tProposer priority\n")
	for _, validator := range validatorResult.Validators {
		fmt.Fprintf(w, "%x\t%s\t%d\t%d\n",validator.Address, base64.StdEncoding.EncodeToString(validator.PubKey.Bytes()), validator.VotingPower, validator.ProposerPriority)
	}
}