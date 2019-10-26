package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	client2 "github.com/Ankr-network/ankr-chain/client"
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/tx/serializer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

var (
	//flag key words which is used in different cmd, use variable as key name
	fileParam         = "file"    //short `f`
	numberAccountParam = "number"  //short name `n`
	outputParam       = "output"  //short name `o`
	privkeyParam      = "privkey" //short name `p`
	addressParam      = "address"
	urlParam          = "nodeurl"
	symbolParam       = "symbol"
	required          = true
	notRequired       = false
	chainIDParam      = "chain-id"
	gasPriceParam     = "gas-price"
	gasLimitParam     = "gas-limit"

	//transaction flags
	toParam     = "to" //short name `t`
	amountParam = "amount"
	memoParam = "memo"
	versionParam = "version"
	abiParam = "abi"
	nameParam = "name"
	methodParam = "method"
	argsParam = "args"
	returnParam = "rtn-type"


	//admin flags
	pubkeyParam          = "pubkey"
	powerParam           = "power"
	dcnameParam          = "dcname"
	nameSpaceParam       = "namespace" //short name `ns`
	keystoreParam        = "keystore"  //short name `k`
	valueParam           = "value"
	adminPrivateKeyParam = "adminkey"
	permParam            = "perm"
	pubKeyParam = "pub-key"
	stakeAddrParam = "stake-addr"
	flagParam = "flag"
	gasUsedParam = "gas-used"
	actionParam = "action"


	//query flags
	heightParam    = "height"
	txidParam      = "txid"
	approveParam   = "approve"
	limitParam     = "limit"
	pageParam      = "page"
	perPageParam   = "perpage"
	meteringParam  = "metering"
	timeStampParam = "timestamp"
	typeParam      = "type"
	fromParam      = "from"
	creatorParam   = "creator"
	detailParam    = "detail"
)

var (
	TxSerializer = serializer.NewTxSerializerCDC()
	ankrCurrency = common.Currency{"ANKR", 18}
)

// retriveUserInput is a function that can retrive user input in form of string. By default,
// it will prompt the user. In test, you can replace this with code that returns the appropriate response.
func retrieveUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	answer = strings.Replace(answer, "\r", "", 1)
	answer = strings.Replace(answer, "\n", "", 1)

	return answer, nil
}

//get the home directory
func configHome() string {
	configHome := os.Getenv("LOCALAPPDATA")
	if configHome == "" {
		// Resort to APPDATA for Windows XP users.
		configHome = os.Getenv("APPDATA")
		if configHome == "" {
			// If still empty, use the default path
			userName := os.Getenv("USERNAME")
			configHome = filepath.Join("C:/", "Users", userName, "AppData", "Local")
		}
	}
	return filepath.Join(configHome, "ankr-chain", "config")
}


//helper functions used in most commands
//add string type flags
func addStringFlag(cmd *cobra.Command, bindKeyName, keyName, shortName, defaultValue, description string, required bool) error {
	cmd.Flags().StringP(keyName, shortName, defaultValue, description)
	err := viper.BindPFlag(bindKeyName, cmd.Flags().Lookup(keyName))
	if err != nil {
		return err
	}
	if required {
		err = cmd.MarkFlagRequired(keyName)
		if err != nil {
			return err
		}
	}
	return nil
}

//add int type flags
func addIntFlag(cmd *cobra.Command, bindKeyName, keyName, shortName string, defaultValue int, description string, requiredFlag bool) error {
	cmd.Flags().IntP(keyName, shortName, defaultValue, description)
	err := viper.BindPFlag(bindKeyName, cmd.Flags().Lookup(keyName))
	if err != nil {
		return err
	}
	if requiredFlag{
		err := cmd.MarkFlagRequired(keyName)
		return err
	}
	return nil
}

//add int64 flags
func addInt64Flag(cmd *cobra.Command, bindKeyName, keyName, shortName string, defaultValue int64, description string, requiredFlag bool) error {
	cmd.Flags().Int64P(keyName, shortName, defaultValue, description)
	err := viper.BindPFlag(bindKeyName, cmd.Flags().Lookup(keyName))
	if err != nil {
		return err
	}
	if requiredFlag {
		err := cmd.MarkFlagRequired(keyName)
		return err
	}
	return nil
}

//add bool type flags
func addBoolFlag(cmd *cobra.Command, bindKeyName, keyName, shortName string, defaultValue bool, description string, requiredFlag bool) error {
	cmd.Flags().BoolP(keyName, shortName, defaultValue, description)
	if requiredFlag {
		err := cmd.MarkFlagRequired(keyName)
		return err
	}
	err := viper.BindPFlag(bindKeyName, cmd.Flags().Lookup(keyName))
	return err
}
func addPersistentString(cmd *cobra.Command, bindKeyName, keyName, shortName, defaultValue, description string, requiredFlag bool) error {
	cmd.PersistentFlags().StringP(keyName, shortName, defaultValue, description)
	err := viper.BindPFlag(bindKeyName, cmd.PersistentFlags().Lookup(keyName))
	if err != nil {
		return err
	}
	if requiredFlag {
		err = cmd.MarkPersistentFlagRequired(keyName)
		if err != nil {
			return err
		}
	}
	return nil
}

func addPersistentInt(cmd *cobra.Command, bindKeyName, keyName, shortName string, defaultValue int, description string, requiredFlag bool) error {
	cmd.PersistentFlags().IntP(keyName, shortName, defaultValue, description)
	err := viper.BindPFlag(bindKeyName, cmd.PersistentFlags().Lookup(keyName))
	if err != nil {
		return err
	}
	if required {
		err = cmd.MarkPersistentFlagRequired(keyName)
		if err != nil {
			return err
		}
	}
	return nil
}

func appendSubCmd(parent *cobra.Command, cmdName, desc string, exec func(cmd *cobra.Command, args []string), flagFunc func(cmd *cobra.Command)) {
	cmd := &cobra.Command{
		Use:   cmdName,
		Short: desc,
		Run:   exec,
	}

	if flagFunc != nil {
		flagFunc(cmd)
	}
	parent.AddCommand(cmd)
}
func newAnkrHttpClient(url string)  *client2.Client{
	return client2.NewClient(url)
}

// display response in json format
func displayStruct(stru interface{})  {
	if stru == nil {
		fmt.Println("[]")
		return
	}
	jsonByte, err := json.MarshalIndent(stru, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonByte))
}

func newTabWriter(out io.Writer) *tabwriter.Writer {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 0, 4, ' ', 0)
	return w
}

