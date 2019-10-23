package cmd

import (
	"github.com/Ankr-network/ankr-cli/mock_cmd"
	"github.com/Ankr-network/dccn-common/wallet"
	"github.com/agiledragon/gomonkey"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"syscall"
	"testing"
)

var(
	ksString = `{
    "address":"E1403CA0DC201F377E820CFA62117A48D4D612400C20D3",
    "publickey":"FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY=",
    "crypto":{
        "cipher":"aes-128-ctr",
        "ciphertext":"f5f7dc4672fb0de7e402b2416c6f0bc00dca80426d4b583c68e964ffff6a6861b56c915061276ddc0e24f1c03208346d6877129d38eff993ca53e5d0f40918c26b80435fde9e341680b285cb7c3e33bf60ceb964040fe316",
        "cipherparams":{
            "iv":"684328737d493bf7a2c6f2aa56956ffe"
        },
        "kdf":"scrypt",
        "kdfparams":{
            "dklen":32,
            "n":262144,
            "p":1,
            "r":8,
            "salt":"cdedb74a765facd1701705a9b8310b3fb93d3c05b6a5f3ed3f6c570d1ee1076b"
        },
        "mac":"fa152097d949a924bac45e03d31567baacd61a9bf3258601278e59fb70f25804"
    },
    "version":3
}`
)
func TestTransaction(t *testing.T) {
	convey.Convey("test transaction transfer function", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockWallet := mock_cmd.NewMockWallet(ctl)
		mockWallet.EXPECT().SendCoins("http://localhost",
			"26657",
			"1gEcOfgXL/rmMvDJPAyL48CanFTeLMU5yASNA9KXmrEVLKr+ZNU879Z3Ew0IwQqIDlRUEVdUvw4CcOyk75u5lg==",
			"E1403CA0DC201F377E820CFA62117A48D4D612400C20D3",
			"F4656949BD747057A59DDF90A218EC352E3916A096924D",
			"20000000000000000000").Return(
				"0x2048bb58bcadc4a15efae927096f18fc1843e6ceb8f52761d9c6199fef408af5", nil)
		mwPatch := gomonkey.ApplyFunc(wallet.SendCoins, mockWallet.SendCoins)
		defer mwPatch.Reset()
		mockTerminal := mock_cmd.NewMockTerminal(ctl)
		mockTerminal.EXPECT().ReadPassword(int(syscall.Stdin)).Return([]byte("123"), nil)
		mtPatch := gomonkey.ApplyFunc(terminal.ReadPassword, mockTerminal.ReadPassword)
		defer mtPatch.Reset()

		sysPatch := gomonkey.ApplyFunc(os.Stat, func(name string) (os.FileInfo, error) {
			return  nil ,nil
		})
		defer sysPatch.Reset()
		ioPatch := gomonkey.ApplyFunc(ioutil.ReadFile, func(filename string) ([]byte, error){
			return []byte(ksString), nil
		})
		defer ioPatch.Reset()


		args := []string{"transaction", "transfer", "--keystore",
			"C:\\Users\\ankr_zhang\\AppData\\Local\\ankr-chain\\config\\UTC--2019-07-23T05-54-59.725333000Z--E1403CA0DC201F377E820CFA62117A48D4D612400C20D3",
			"--to", "F4656949BD747057A59DDF90A218EC352E3916A096924D",
			"--amount", "20000000000000000000",
			"--nodeurl", localUrl}
		//args := []string{"transaction", "transfer", "--help"}
		cmd := RootCmd
		cmd.SetArgs(args)
		cmd.Execute()
	})
}