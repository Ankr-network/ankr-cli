package cmd

import (
	"fmt"
	"github.com/Ankr-network/ankr-cli/mock_cmd"
	"github.com/Ankr-network/dccn-common/wallet"
	"github.com/agiledragon/gomonkey"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"syscall"
	"testing"
)

//normal cases
func TestGetBalance(t *testing.T) {
	convey.Convey("test get balance", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockWallet := mock_cmd.NewMockWallet(ctl)
		mockWallet.EXPECT().GetBalance("https://chain-01.dccn.ankr.com", "443", "95CD00025C3807CEE9804D19B1E410A30A47B303371C12").Return("50000", nil)
		patch := gomonkey.ApplyFunc(wallet.GetBalance, mockWallet.GetBalance)
		defer patch.Reset()

		args := []string{"account", "getbalance", "--address", "95CD00025C3807CEE9804D19B1E410A30A47B303371C12", "--nodeurl", "https://chain-01.dccn.ankr.com:443"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestGenAccount(t *testing.T) {
	convey.Convey("test generate account command", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock terminal
		mockTerminal := mock_cmd.NewMockTerminal(ctl)
		mockTerminal.EXPECT().ReadPassword(int(syscall.Stdin)).Return([]byte("123"), nil).Times(2)
		patch := gomonkey.ApplyFunc(terminal.ReadPassword, mockTerminal.ReadPassword)
		defer patch.Reset()

		//mock wallet
		mockWallet := mock_cmd.NewMockWallet(ctl)
		privKey := "PmWSb6C8a1dE0mBC3+rSkRHdHUXqQZy73cBc5KNEn3cF8fMkvyiIB1eXCa25D7qIt4vPCay/zwTp4/Jb0aKo+Q=="
		pubKey := "BfHzJL8oiAdXlwmtuQ+6iLeLzwmsv88E6ePyW9GiqPk="
		address := "5D6D5EE541DC521001385F542EC7332416E6565F8A269F"
		mockWallet.EXPECT().GenerateKeys().Return(privKey, pubKey, address)
		walletPatch := gomonkey.ApplyFunc(wallet.GenerateKeys, mockWallet.GenerateKeys)
		defer walletPatch.Reset()

		args := []string{"account", "genaccount","-o", "./tmp"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
		removeKeyStore("./tmp", address)
	})
}

var keyStoreString = `{"address":"B47982CF51CD7718FE25BA96B707C449F4F917949E7A25","publickey":"TokCLsXk6Z4lhyr/v7PkSksuHkDQURrjVn1Pt8JsKNI=","crypto":{"cipher":"aes-128-ctr","ciphertext":"3e0d4968e664d201682fecbba66e9fa2b8f5d4ccd9445c45705a9864743234aed4e2489d26025fd97a3fa7f23b3714e580b7a38bac1b5c734e75a6c5fac9cd922239f54934f7cc073b65f5fb914166f7964f8fc07093a4c6","cipherparams":{"iv":"e9b282f254f2bcdf98e42cb32c127b06"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"f3698a69cb65c64fa9548febbd19174f1af898cd56590420e5a5fe0578bd351c"},"mac":"eaf9d07ff51732eb3c58e561f0ffc37bf099c5570d2a51388cc5b2f3d7420860"},"version":3}`
func TestExportPriv(t *testing.T) {
	convey.Convey("test exporting private key from keystore", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockTerminal := mock_cmd.NewMockTerminal(ctl)
		mockTerminal.EXPECT().ReadPassword(int(syscall.Stdin)).Return([]byte("123"), nil)
		patch := gomonkey.ApplyFunc(terminal.ReadPassword, mockTerminal.ReadPassword)
		defer patch.Reset()

		filePatch := gomonkey.ApplyFunc(ioutil.ReadFile, func(file string) ([]byte, error){
			return []byte(keyStoreString), nil
		})
		defer filePatch.Reset()

		args := []string{"account", "exportprivatekey", "--file", "./keystore"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestResetPWD(t *testing.T) {
	convey.Convey("test exporting private key from keystore", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockTerminal := mock_cmd.NewMockTerminal(ctl)
		gomock.InOrder(
			mockTerminal.EXPECT().ReadPassword(int(syscall.Stdin)).Return([]byte("123"), nil),
			mockTerminal.EXPECT().ReadPassword(int(syscall.Stdin)).Return([]byte("abcd"), nil),
			mockTerminal.EXPECT().ReadPassword(int(syscall.Stdin)).Return([]byte("abcd"), nil),
			)
		patch := gomonkey.ApplyFunc(terminal.ReadPassword, mockTerminal.ReadPassword)
		defer patch.Reset()

		filePatch := gomonkey.ApplyFunc(ioutil.ReadFile, func(file string) ([]byte, error){
			return []byte(keyStoreString), nil
		})
		defer filePatch.Reset()

		args := []string{"account", "resetpwd", "--file", "tmp/keystore"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

//help functions
func removeKeyStore(dir string, address string)  {
	_files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range  _files {
		if !file.IsDir(){
			if strings.HasSuffix(file.Name(), address) {
				err = os.Remove(path.Join(dir, file.Name()))
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}