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
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/rpc/client"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"os"
	"regexp"
	"strconv"
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
	trxPage = "trxPage"
	trxPerPage = "trxPerPage"
	trxMetering = "trxMetering"
	trxTimeStamp = "trxTimeStamp"
	trxType = "trxType"
	trxHeight = "trxHeight"
	trxCreator = "trxCreator"
	trxFrom = "trxFrom"
	trxTo = "trxTo"

	//bind block flags
	blockHeight = "blockHeight"
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

     //error message
     tooManyFlags = "Too many flags, please refer to transaction help"
     tooFewFlags = "Too few flags, requires at least one more flag"

     //types used in query transaction flags
     querySend = "send"
     querySetBalance = "setbalance"
     querySetStake = "setstake"
     queryUpdateValidator = "updatevalidator"
     querySetMetering = "setmetering"
)

var (
	txSearchFlags = []string{txidFlag, meteringFlag, timeStampFlag, typeFlag, fromFlag, toFlag, heightFlag, creatorFlag}
	periodRegexp = `((\(|\[)\d\:\d+(\]|\()|\d+)`
	reg, _ = regexp.Compile(periodRegexp)
)
// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query information from ankr chain",
}

func init() {
	err := addPersistentString(queryCmd, queryUrl, urlFlag, "", "", "validator url", required)
	if err != nil {
		panic(err)
	}
	appendSubCmd(queryCmd, "transaction","transaction allows you to query the transaction results with multiple conditions.", transactionInfo, addTransactionInfoFlags)
	appendSubCmd(queryCmd, "block", "Get block at a given height. If no height is provided, it will fetch the latest block. And you can use \"detail\" to show more information about transactions contained in block",
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

	//check the number of flags cmd received, too few or too many shall raise an error
	numFlags := cmd.Flags().NFlag()
	if numFlags > 5 {
		fmt.Println(tooManyFlags)
		return
	}
	if numFlags <2 {
		fmt.Println(tooFewFlags)
		return
	}
	//collectedFlags := make([]string, 0, numFlags)
	collectedFlags := make(map[string] string, numFlags)
	for _, flag := range txSearchFlags {
		if cmd.Flag(flag).Changed {
			collectedFlags[flag] = cmd.Flag(flag).Value.String()
		}
	}
	if !isValidFlags(collectedFlags) {
		fmt.Println("Invalid flags is received.")
		return
	}

	validatorUrl = viper.GetString(queryUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)

	//query transaction --txid hash --nodeurl https://xx:xx --prove bool
	if txid, ok := collectedFlags[txidFlag]; ok {
		if 	strings.HasPrefix(txid, "0x") {
			txid = strings.TrimLeft(txid, "0x")
		}
		hash, err := hex.DecodeString(txid)
		if err != nil {
			fmt.Println("Invalid transaction id format!")
			return
		}
		rpcTransaction(cl, hash)
		return
	}

	//query transaction --type/from/to/metering/timestamp
	query := formatQueryContent(collectedFlags)
	resp, err := doTxSearch(cl, query)
	if err != nil {
		fmt.Println("Query transaction failed.")
		fmt.Println(err)
		return
	}
	w := newTabWriter(os.Stdout)
	//write txSearchResult header
	fmt.Fprintln(w, "TotalCount:\t",resp.TotalCount)
	fmt.Fprintf(w, "type\thash\theight\tindex\tdetail\n")
	writeTxSearchResult(resp, w)
	w.Flush()
}
func formatQueryContent(flags map[string]string) string {
	result := make([]string, 0, len(flags))
	var query string
	for key, value := range flags {
		switch key {
		case meteringFlag:
			query = fmt.Sprintf("app.metering='%s'",value)
		case timeStampFlag:
			valueSlice := strings.Split(value, ":")
			if len(valueSlice) == 1 {
				//if only one digit is received in interval, trim bracket
				value = strings.TrimPrefix(value,"[")
				value = strings.TrimPrefix(value,"(")
				value = strings.TrimRight(value,"]")
				value = strings.TrimRight(value,")")
				query = fmt.Sprintf("app.timestamp=%s",value)
				break
			}
			interval := formatInterval(value)
			if len(interval) != 2{
				query = fmt.Sprintf("app.timestamp%s", interval[0])
				break
			}
			query = fmt.Sprintf("app.timestamp%s and app.timestamp%s", interval[0], interval[1])
		case typeFlag:
			query = fmt.Sprintf("app.type='%s'",value)
		case fromFlag:
			query = fmt.Sprintf("app.fromaddress='%s'",value)
		case toFlag:
			query = fmt.Sprintf("app.toaddress='%s'",value)
		case creatorFlag:
			query = fmt.Sprintf("app.creator='%s'",value)
		case heightFlag:
			valueSlice := strings.Split(value, ":")
			if len(valueSlice) == 1 {
				value = strings.TrimPrefix(value,"[")
				value = strings.TrimPrefix(value,"(")
				value = strings.TrimRight(value,"]")
				value = strings.TrimRight(value,")")
				query = fmt.Sprintf("tx.height=%s",value)
				break
			}
			interval := formatInterval(value)
			if len(interval) != 2{
				query = fmt.Sprintf("tx.height%s", interval[0])
				break
			}
			query = fmt.Sprintf("tx.height%s and tx.height%s", interval[0], interval[1])
		}
		result =append(result, query)
	}
	return strings.Join(result, " and ")
}
func formatInterval(period string) []string {
	periodSlice := strings.Split(period, ":")
	leftOp := []rune(periodSlice[0])[0]
	length := len(periodSlice[1])
	rightOp := []rune(periodSlice[1])[length-1]
	var leftValue, rightValue string
	switch leftOp {
	case '(':
		leftValue = strings.TrimLeft(periodSlice[0],"(")
		if len(leftValue) > 0 {
			leftValue = fmt.Sprintf(">%s",string(leftValue))
		}
	case '[':
		leftValue = strings.TrimLeft(periodSlice[0],"[")
		if len(leftValue) > 0 {
			leftValue = fmt.Sprintf(">%s",string(leftValue))
		}
	}

	switch rightOp {
	case ')':
		rightValue = strings.TrimRight(periodSlice[1],")")
		if len(rightValue) > 0{
			rightValue = fmt.Sprintf("<%s",rightValue)
		}
	case ']':
		rightValue = strings.TrimRight(periodSlice[1],"]")
		if len(rightValue) > 0{
			rightValue = fmt.Sprintf("<=%s",rightValue)
		}
	}
	result := make([]string, 0 , 2)
	if leftValue != ""{
		result = append(result, leftValue)
	}
	if rightValue != ""{
		result = append(result, rightValue)
	}
	return result
}

//rules of mutiple query condition, check if the cmd flags are correctly set
func isValidFlags(flags map[string]string) bool {
	//if txid is set, only nodeurl is valid
	if _, ok := flags[txidFlag]; ok{
		// after getting ride of prove flag, total flags should be no more than 2
		if len(flags) != 2{
			return false
		}
		return true
	}

	//check time and height format
	if timeStamp, ok := flags[timeStampFlag]; ok {
		if matched := reg.MatchString(timeStamp); !matched {
			return false
		}
	}
	if height, ok := flags[heightFlag]; ok {
		if matched := reg.MatchString(height); !matched {
			return false
		}
	}

	_, existsFrom := flags[fromFlag]
	_, existsTo := flags[toFlag]
	//if transaction type is set, check if other flags is valid
	if value, ok := flags[typeFlag];ok {
		value = strings.ToLower(value)
		switch value {
		case querySetBalance, querySetStake, queryUpdateValidator:
			// after getting ride of page and perpage flag, total flags should be no more than 2
			if len(flags) != 2 {
				return false
			}
		case querySetMetering:
			if existsFrom || existsTo {
				return false
			}
		case querySend:
			return true
		default:   //unknown transaction type
			return false
		}
		return true
	}

	if _, ok := flags[meteringFlag]; ok {
		if existsTo || existsFrom {
			return false
		}
	}
	return true
}

//query transaction and display result
func rpcTransaction(cl *client.HTTP, hash []byte)  {
	prove := viper.GetBool(approveFlag)
	resp, err := cl.Tx(hash, prove)
	if err != nil {
		fmt.Println("Failed to query transaction.")
		fmt.Println(err)
		return
	}
	result := parseTx(resp)
	w := newTabWriter(os.Stdout)
	fmt.Fprintf(w, "type\thash\theight\tindex\tdetail\n")
	writeTx(result, w)
	w.Flush()
}

func writeTxSearchResult(result *core_types.ResultTxSearch, w *tabwriter.Writer)  {
	for _, txResult := range result.Txs {
		txParsed := parseTx(txResult)
		writeTx(txParsed, w)
	}
}

func doTxSearch(cl *client.HTTP, qt string) (*core_types.ResultTxSearch, error) {
	queryPage := viper.GetInt(trxPage)
	queryPerPage := viper.GetInt(trxPerPage)
	prove := viper.GetBool(trxApprove)
	resp, err := cl.TxSearch(qt, prove, queryPage, queryPerPage)
	return resp, err
}

func addTransactionInfoFlags(cmd *cobra.Command)  {
	err := addStringFlag(cmd, trxTxid, txidFlag, "", "", "The transaction hash", "")
	if err != nil {
		panic(err)
	}
	err = addBoolFlag(cmd, trxApprove, approveFlag, "", false, "Include a proof of the transaction inclusion in the block", "")
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, trxFrom, fromFlag, "", "", "the from address contained in a transaction", "")
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, trxTo, toFlag, "", "", "the to address contained in a transaction", "")
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, trxTimeStamp, timeStampFlag, "", "",
		"transaction executed timestamp. Input can be an exactly unix timestamp  or a time interval separate by \":\", and time interval should be enclosed with \"[]\" or \"()\" which is mathematically open interval and close interval." ,
		"")
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, trxMetering, meteringFlag, "", "", "query metering transaction, both datacenter name and namespace should be  provided and separated  by \":\"", "")
	if err != nil {
		panic(err)
	}

	err = addIntFlag(cmd, trxPage, pageFlag, "", 1, "Page number (1 based)", "")
	if err != nil {
		panic(err)
	}

	err = addIntFlag(cmd, trxPerPage, perPageFlag, "", 30, "Number of entries per page(max: 100)", "")
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, trxType, typeFlag, "", "", "Ankr chain predefined types, SetMetering, SetBalance, UpdatValidator, SetStake, Send", "")
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, trxHeight, heightFlag, "", "",
		"block height. Input can be an exactly block height  or a height interval separate by \":\", and height interval should be enclosed with \"[]\" or \"()\" which is mathematically open interval and close interval.", "")
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, trxCreator, creatorFlag, "", "", "app creator", "")
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
	//height := viper.GetString(blockHeight)
	from, to, err := getBlockInterval()
	if err != nil {
		fmt.Println(err)
		return
	}

	for iter := from; iter <= to; iter ++ {
		heightInt := int64(iter)
		heightP := &heightInt
		if heightInt == -1 {
			heightP = nil
		}
		resp, err :=cl.Block(heightP)
		if err != nil {
			fmt.Println("Query block failed.", err)
			return
		}
		detail := false
		if len(args) > 0 && args[0] == "detail"{
			detail = true
		}
		outPutBlockResp(resp, detail)
	}
}

func getBlockInterval() (from int, to int,err  error) {
	from = -1
	to = -1
	heightStr := viper.GetString(blockHeight)

	//height flag is not set, get the latest block
	if heightStr == ""{
		return from, to, nil
	}

	//if height flags is not set properly, return error
	if matched := reg.MatchString(heightStr); !matched {
		return from, to , errors.New("Invalid Height format, should be \"[from:to]\". ")
	}

	//strictly flow the rule [from:to]
	heightStr = strings.TrimLeft(heightStr, "[")
	heightStr = strings.TrimRight(heightStr, "]")
	height, err := strconv.Atoi(heightStr)
	if err == nil {
		from = height
		to = height
		return from, to, nil
	}
	heightSlice := strings.Split(heightStr, ":")
	if len(heightSlice) != 2 {
		return from, to, errors.New("input both from and to separated with \":\". ")
	}
	fromStr, toStr := heightSlice[0],heightSlice[1]
	from, err = strconv.Atoi(fromStr)
	if err != nil {
		return from, to, errors.New("from is not an integer. ")
	}
	to, err = strconv.Atoi(toStr)
	if err != nil {
		return from, to, errors.New("to is not an integer. ")
	}
	if from >to {
		return from, to, errors.New("from should be less or equal than to")
	}
	return from, to, nil
}

func outPutBlockResp(resp *core_types.ResultBlock, detail bool)  {
	w := newTabWriter(os.Stdout)
	fmt.Fprintf(w, "\nBlock info:\n")
	outPutHeader(w, resp.Block.Header)
	fmt.Fprintf(w, "\nTransactions contained in block: \n")
	if resp.Block.Txs == nil || len(resp.Block.Txs) == 0 {
		fmt.Fprintf(w, "[]\n")
	}else{
		if detail {
			outPutTransactions(w, resp.Block.Txs)
		}else {
			outPutTransactionsSimple(w, resp.Block.Txs)
		}
	}
	w.Flush()
}

func addQueryBlockFlags(cmd *cobra.Command)  {
	err := addStringFlag(cmd, blockHeight, heightFlag, "", "", "height interval of the blocks to query. integer or block interval formatted as [from:to] are accepted ", "" )
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
func writeTx(rt ResultTx, w *tabwriter.Writer)  {
	//table header
	fmt.Fprintf(w, "%s\t%s\t%d\t%d\t", rt.Type, rt.Hash, rt.Height, rt.Index)
	//table contents
	switch rt.Type {
	case "transfer":
		fmt.Fprintf(w, "from: %s\tto:%s\tamount:%s\tnonce:%s\n",rt.Data["from"],rt.Data["to"],rt.Data["amount"],rt.Data["nonce"])
		//fmt.Fprintf(w, "from: %s\n",rt.Data["from"])
		//fmt.Fprintf(w, "\t\t\t\tnonce: %s\n",rt.Data["from"])
		//fmt.Fprintf(w, "\t\t\t\tto: %s\n",rt.Data["from"])
		//fmt.Fprintf(w, "\t\t\t\tamount: %s\n",rt.Data["from"])
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

func outPutTransactionsSimple(w *tabwriter.Writer,txs types.Txs)  {
	fmt.Fprintf(w, "hash\n")
	for _, tx := range txs {
		hash := fmt.Sprintf("0x%x",tx.Hash())
		fmt.Fprintf(w,"%s\n", hash)
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