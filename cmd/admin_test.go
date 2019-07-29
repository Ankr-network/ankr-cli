package cmd

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAdmin(t *testing.T) {
	convey.Convey("test get balance", t, func() {
		args := []string{"account", "getbalance", "--address", "95CD00025C3807CEE9804D19B1E410A30A47B303371C12", "--url", "https://chain-01.dccn.ankr.com:443"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}
