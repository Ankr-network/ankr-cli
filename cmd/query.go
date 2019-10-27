package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	client2 "github.com/Ankr-network/ankr-chain/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/rpc/client"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	defaultPerPage = 20
	maxPerPage     = 50
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
	trxDetail = "trxDetail"

	//bind block flags
	blockHeight = "blockHeight"
	blockPage = "blockPage"
	blockPerPage ="blockPerPage"
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
	txSearchFlags = []string{meteringParam, timeStampParam, typeParam, fromParam, toParam, heightParam, creatorParam}
	periodRegexp = `((\(|\[)\d\:\d+(\]|\()|\d+)`
	reg, _ = regexp.Compile(periodRegexp)
)
// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query information from ankr chain",
}

func init() {
	err := addPersistentString(queryCmd, queryUrl, urlParam, "", "", "validator url", required)
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
	appendSubCmd(queryCmd, "contract", "get smart contract data", runGetContract, addGetContractFlags)
}

func transactionInfo(cmd *cobra.Command, args []string)  {
	client := newAnkrHttpClient(viper.GetString(queryUrl))
	prove := viper.GetBool(trxApprove)

	// query --txid 0xxxxxx --nodeurl url
	if cmd.Flag(txidParam).Changed {
		txid := viper.GetString(trxTxid)
		txid = strings.TrimLeft(txid,"0x")
		txidByte, err := hex.DecodeString(txid)
		if err != nil {
			fmt.Println("Invalid txid.")
			return
		}
		resp, err := client.Tx(txidByte, prove)
		detail := viper.GetBool(trxDetail)
		displayTxMsg(resp, detail)
		return
	}

	//collectedFlags := make([]string, 0, numFlags)
	collectedFlags := make(map[string] string)
	for _, flag := range txSearchFlags {
		if cmd.Flag(flag).Changed {
			collectedFlags[flag] = cmd.Flag(flag).Value.String()
		}
	}
	//query transaction --txid hash --nodeurl https://xx:xx --prove bool

	//query transaction --type/from/to/metering/timestamp
	query := formatQueryContent(collectedFlags)
	page := viper.GetInt(trxPage)
	perPage := viper.GetInt(trxPerPage)
	resp, err := client.TxSearch(query, prove, page, perPage)
	if err != nil {
		fmt.Println("Transaction search failed.")
		fmt.Println(err)
		return
	}
	detail := viper.GetBool(trxDetail)
	fmt.Println("Total Tx Count:", resp.TotalCount)
	fmt.Println("Transactions search result:")
	for _, tx := range resp.Txs {
		displayTxMsg(tx, detail)
	}
}

func displayTxMsg(txMsg *core_types.ResultTx, detail bool)  {
	displayStruct(txMsg)
	if detail{
		displayTx(txMsg.Tx)
	}
}

func displayTx(data []byte)  {
	decoder := new(client2.TxDecoder)
	tx, err := decoder.Decode(data)
	if err != nil {
		fmt.Println("Decode transaction error!")
		fmt.Println(err)
		return
	}
	jsonByte, err := json.MarshalIndent(tx, "", "\t")
	if err != nil {
		fmt.Println("Marshal Error.")
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonByte))
}

func formatQueryContent(parameters map[string]string) string {
	result := make([]string, 0, len(parameters))
	var query string
	for key, value := range parameters {
		switch key {
		case meteringParam:
			query = fmt.Sprintf("app.metering='%s'",value)
		case timeStampParam:
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
		case typeParam:
			query = fmt.Sprintf("app.type='%s'",value)
		case fromParam:
			query = fmt.Sprintf("app.fromaddress='%s'",value)
		case toParam:
			query = fmt.Sprintf("app.toaddress='%s'",value)
		case creatorParam:
			query = fmt.Sprintf("app.creator='%s'",value)
		case heightParam:
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

//query transaction and display result
func rpcTransaction(cl *client2.Client, hash []byte)  {
	resp, err := cl.Tx(hash, viper.GetBool(trxApprove))
	if err != nil {
		fmt.Println("Failed to query transaction.")
		fmt.Println(err)
		return
	}
	displayStruct(resp)
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
	err := addStringFlag(cmd, trxTxid, txidParam, "", "", "The transaction hash", notRequired)
	if err != nil {
		panic(err)
	}
	err = addBoolFlag(cmd, trxApprove, approveParam, "", false, "Include a proof of the transaction inclusion in the block", notRequired)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, trxFrom, fromParam, "", "", "the from address contained in a transaction", notRequired)
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, trxTo, toParam, "", "", "the to address contained in a transaction", notRequired)
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, trxTimeStamp, timeStampParam, "", "",
		"transaction executed timestamp. Input can be an exactly unix timestamp  or a time interval separate by \":\", and time interval should be enclosed with \"[]\" or \"()\" which is mathematically open interval and close interval." ,
		notRequired)
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, trxMetering, meteringParam, "", "", "query metering transaction, both datacenter name and namespace should be  provided and separated  by \":\"", notRequired)
	if err != nil {
		panic(err)
	}

	err = addIntFlag(cmd, trxPage, pageParam, "", 1, "Page number (1 based)", notRequired)
	if err != nil {
		panic(err)
	}

	err = addIntFlag(cmd, trxPerPage, perPageParam, "", 30, "Number of entries per page(max: 100)", notRequired)
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, trxType, typeParam, "", "", "Ankr chain predefined types, SetMetering, SetBalance, UpdatValidator, SetStake, Send", notRequired)
	if err != nil {
		panic(err)
	}

	err = addStringFlag(cmd, trxHeight, heightParam, "", "",
		"block height. Input can be an exactly block height  or a height interval separate by \":\", and height interval should be enclosed with \"[]\" or \"()\" which is mathematically open interval and close interval.", notRequired)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(cmd, trxCreator, creatorParam, "", "", "app creator", notRequired)
	if err != nil {
		panic(err)
	}
	err = addBoolFlag(cmd, trxDetail, detailParam, "", false, "display transaction detail", notRequired)
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
	cl := newAnkrHttpClient(viper.GetString(queryUrl))
	from, to, err := getBlockInterval()
	if err != nil {
		fmt.Println(err)
		return
	}

	resps := make([]*core_types.ResultBlock,0, to - from + 1)
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
		resps = append(resps, resp)
	}

	detail := false
	if len(args) > 0 && args[0] == "detail"{
		detail = true
	}
	page := viper.GetInt(blockPage)
	perPage := viper.GetInt(trxPerPage)
	outPutBlockResp(resps, page, perPage, detail)
}

func validatePage(page, perPage, totalCount int) int {
	if perPage < 1 {
		return 1
	}

	pages := ((totalCount - 1) / perPage) + 1
	if page < 1 {
		page = 1
	} else if page > pages {
		page = pages
	}

	return page
}

func validatePerPage(perPage int) int {
	if perPage < 1 {
		return defaultPerPage
	} else if perPage > maxPerPage {
		return maxPerPage
	}
	return perPage
}

func validateSkipCount(page, perPage int) int {
	skipCount := (page - 1) * perPage
	if skipCount < 0 {
		return 0
	}

	return skipCount
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

func outPutBlockResp(resps []*core_types.ResultBlock,page int, perPage int, detail bool)  {
	totalCount := len(resps)
	page = validatePage(page, perPage, totalCount)
	perPage = validatePerPage(perPage)
	skipCount := validateSkipCount(page, perPage)
	resultLength := common.MinInt(perPage, totalCount - skipCount)

	fmt.Println("\nTotal ount:", totalCount)
	for i := 0; i < resultLength; i ++{
		resp := resps[i + skipCount]
		fmt.Println( "\nBlock info:")
		outPutHeader( resp.Block.Header)
		fmt.Println( "\nTransactions contained in block:")
		if resp.Block.Txs == nil || len(resp.Block.Txs) == 0 {
			fmt.Println( "[]")
		}else{
			for _, tx := range resp.Block.Txs {
				displayTx(tx)
			}
		}
	}
}

func addQueryBlockFlags(cmd *cobra.Command)  {
	err := addStringFlag(cmd, blockHeight, heightParam, "", "", "height interval of the blocks to query. integer or block interval formatted as [from:to] are accepted ", notRequired )
	if err != nil {
		panic(err)
	}
	err = addIntFlag(cmd, blockPage, pageParam, "", 1, "Page number (1 based)", notRequired)
	if err != nil {
		panic(err)
	}
	err = addIntFlag(cmd, blockPerPage, perPageParam, "", 20, "Page number (1 based)", notRequired)
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
	displayStruct(resp)
	//w := newTabWriter(os.Stdout)
	//outPutValidator(w, resp)
	//w.Flush()

}
func addQueryValidatorFlags(cmd *cobra.Command)  {
	err := addInt64Flag(cmd, validatorHeight, heightParam, "", -1, "block height", notRequired)
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
	displayStruct(resp)
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
	displayStruct(resp)
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
	displayStruct(resp)
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
	displayStruct(resp)
}

//query unconfirmed transactions
func queryUnconfirmedTxs(cmd *cobra.Command, args []string)  {
	cl := newAnkrHttpClient(viper.GetString(queryUrl))
	limmit := viper.GetInt(unconfirmedTxLimit)
	resp, err := cl.UnconfirmedTxs(limmit)
	if err != nil {
		fmt.Println("Query unconfirmed transactions failed.", err)
		return
	}
	outputTxResult(resp)
}

func outputTxResult(txResult *core_types.ResultUnconfirmedTxs)  {
	fmt.Println( "n_tx: ", txResult.Count)
	fmt.Println( "total:", txResult.Total)
	fmt.Println("total_bytes:", txResult.TotalBytes)
	fmt.Println("transactions:")
	if len(txResult.Txs) == 0 {
		fmt.Println("[]\n")
	}else {
		for _, tx := range txResult.Txs {
			displayTx(tx)
		}
	}
}

func addQueryUncofirmedTxsFlags(cmd *cobra.Command)  {
	err := addIntFlag(cmd, unconfirmedTxLimit, limitParam, "",30, "number of entries", notRequired)
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
	outputTxResult(resp)
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

func outPutHeader(header types.Header)  {
	//information to be displayed in the window
	fmt.Println("Version: ", header.Version)
	fmt.Println("Chain-Id:", header.ChainID)
	fmt.Println("Height: ", header.Height)
	fmt.Println("Time:", header.Time)
	fmt.Println("Number-Txs: ", header.NumTxs)
	fmt.Println("Total-Txs:", header.TotalTxs)
	fmt.Println("Last-block-id: ", header.LastBlockID)
	fmt.Println( "Last-commit-hash:",header.LastCommitHash)
	fmt.Println("Data-hash: ", header.DataHash)
	fmt.Println("Validator:", header.ValidatorsHash)
	fmt.Println("Consensus: ", header.ConsensusHash)
	fmt.Println("Version: ", header.Version)
	fmt.Println("App-hash:", header.AppHash)
	fmt.Println("Proposer-Address:", header.ProposerAddress)
}

func outPutValidator(w *tabwriter.Writer, validatorResult *core_types.ResultValidators){
	fmt.Fprintf(w, "Height:%d\n",validatorResult.BlockHeight)
	fmt.Fprintf(w, "\nValidators information: \n")
	fmt.Fprintf(w, "Address\tPubkey\tVoting-Power\tProposer priority\n")
	for _, validator := range validatorResult.Validators {
		fmt.Fprintf(w, "%x\t%s\t%d\t%d\n",validator.Address, base64.StdEncoding.EncodeToString(validator.PubKey.Bytes()), validator.VotingPower, validator.ProposerPriority)
	}
}