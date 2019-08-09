package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/rpc/client"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

var (
	//flag key words which is used in different cmd, use variable as key name
	fileFlag          = "file"    //short `f`
	numberAccountFlag = "number"  //short name `n`
	outputFlag        = "output"  //short name `o`
	privkeyFlag       = "privkey" //short name `p`
	addressFlag       = "address"
	urlFlag           = "nodeurl"
	required          = "required"

	//transaction flags
	toFlag     = "to" //short name `t`
	amountFlag = "amount"

	//admin flags
	pubkeyFlag          = "pubkey"
	powerFlag           = "power"
	dcnameFlag          = "dcname"
	nameSpaceFlag       = "namespace" //short name `ns`
	keystoreFlag        = "keystore"  //short name `k`
	valueFlag           = "value"
	adminPrivateKeyFlag = "adminkey"
	permFlag            = "perm"

	//query flags
	heightFlag = "height"
	txidFlag = "txid"
	approveFlag = "approve"
	limitFlag = "limit"
	pageFlag = "page"
	perPageFlag = "perpage"
	meteringFlag = "metering"
	timeStampFlag = "timestamp"
	typeFlag = "type"
	fromFlag = "from"
	creatorFlag = "creator"
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
func addStringFlag(cmd *cobra.Command, bindKeyName, keyName, shortName, defaultValue, description, required string) error {
	cmd.Flags().StringP(keyName, shortName, defaultValue, description)
	err := viper.BindPFlag(bindKeyName, cmd.Flags().Lookup(keyName))
	if err != nil {
		return err
	}
	if required == "required" {
		err = cmd.MarkFlagRequired(keyName)
		if err != nil {
			return err
		}
	}
	return nil
}

//add int type flags
func addIntFlag(cmd *cobra.Command, bindKeyName, keyName, shortName string, defaultValue int, description, requiredFlag string) error {
	cmd.Flags().IntP(keyName, shortName, defaultValue, description)
	err := viper.BindPFlag(bindKeyName, cmd.Flags().Lookup(keyName))
	if err != nil {
		return err
	}
	if requiredFlag == required {
		err := cmd.MarkFlagRequired(keyName)
		return err
	}
	return nil
}

//add int64 flags
func addInt64Flag(cmd *cobra.Command, bindKeyName, keyName, shortName string, defaultValue int64, description, requiredFlag string) error {
	cmd.Flags().Int64P(keyName, shortName, defaultValue, description)
	err := viper.BindPFlag(bindKeyName, cmd.Flags().Lookup(keyName))
	if err != nil {
		return err
	}
	if requiredFlag == required {
		err := cmd.MarkFlagRequired(keyName)
		return err
	}
	return nil
}

//add bool type flags
func addBoolFlag(cmd *cobra.Command, bindKeyName, keyName, shortName string, defaultValue bool, description, requiredFlag string) error {
	cmd.Flags().BoolP(keyName, shortName, defaultValue, description)
	if requiredFlag == required {
		err := cmd.MarkFlagRequired(keyName)
		return err
	}
	err := viper.BindPFlag(bindKeyName, cmd.Flags().Lookup(keyName))
	return err
}
func addPersistentString(cmd *cobra.Command, bindKeyName, keyName, shortName, defaultValue, description, requiredFlag string) error {
	cmd.PersistentFlags().StringP(keyName, shortName, defaultValue, description)
	err := viper.BindPFlag(bindKeyName, cmd.PersistentFlags().Lookup(keyName))
	if err != nil {
		return err
	}
	if requiredFlag == required {
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
func newAnkrHttpClient(url string)  *client.HTTP{
	return client.NewHTTP(url, "/websocket")
}

//display a json struct type with pretty format
func display(v interface{})  {
	if v == nil {
		fmt.Println("[]")
		return
	}
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}

//display information in table
func newTabWriter(out io.Writer) *tabwriter.Writer {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 0, 4, ' ', 0)
	return w
}

