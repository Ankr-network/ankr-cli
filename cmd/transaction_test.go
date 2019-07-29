package cmd

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTransaction(t *testing.T) {
	convey.Convey("test transaction transfer function", t, func() {
		//args := []string{"transaction", "transfer", "--keystore", "C:\\Users\\ankr_zhang\\AppData\\Local\\ankr-chain\\config\\UTC--2019-07-23T05-54-59.725333000Z--B47982CF51CD7718FE25BA96B707C449F4F917949E7A25",
		//	"--to", "0x0", "--amount", "jdls"}
		args := []string{"transaction", "transfer", "--help"}
		cmd := RootCmd
		cmd.SetArgs(args)
		cmd.Execute()
	})
}
