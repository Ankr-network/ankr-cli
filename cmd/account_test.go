package cmd

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

//normal cases
func TestGetBalance(t *testing.T) {
	convey.Convey("test get balance", t, func() {
		args := []string{"account", "getbalance", "--address", "95CD00025C3807CEE9804D19B1E410A30A47B303371C12", "--url", "https://chain-01.dccn.ankr.com:443"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestGenAccount(t *testing.T) {
	convey.Convey("test generate account command", t, func() {
		args := []string{"account", "genaccount"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestGenKeyStore(t *testing.T) {
	convey.Convey("test generate keystore", t, func() {
		args := []string{"account", "genkeystore", "--privkey", "ntIlBuH4gTHQFQqixN+meQ7yjnoohzd3KF6P5Q+QoLVOiQIuxeTpniWHKv+/s+RKSy4eQNBRGuNWfU+3wmwo0g=="}
		cmd := RootCmd
		cmd.SetArgs(args)
		cmd.Execute()
	})
}

func TestExportPriv(t *testing.T) {
	convey.Convey("test exporting private key from keystore", t, func() {
		args := []string{"account", "exportprivatekey", "--file", "C:\\Users\\ankr_zhang\\AppData\\Local\\ankr-chain\\config\\UTC--2019-07-23T05-54-59.725333000Z--B47982CF51CD7718FE25BA96B707C449F4F917949E7A25"}
		cmd := RootCmd
		cmd.SetArgs(args)
		cmd.Execute()
	})
}
