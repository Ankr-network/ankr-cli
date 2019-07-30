package cmd

import (
	"github.com/Ankr-network/dccn-common/wallet"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	adminPrivate = "0mqsOtVueE7uq/I5J/dAhesumWXTu619xXuRgtj4l0d0ELMH6X9ZjGqT6Lnhrhp13LVeGIgrm3QgBnk4q16BZg=="
)
func TestAdmin(t *testing.T) {
	convey.Convey("test get balance", t, func() {
		args := []string{"admin", "setbalance", "--address", "95CD00025C3807CEE9804D19B1E410A30A47B303371C12","--amount","12000000000000000000" ,"--url", "http://localhost:26657", "--privkey",adminPrivate}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
		bal, err := wallet.GetBalance("http://localhost","26657","95CD00025C3807CEE9804D19B1E410A30A47B303371C12")
		if err != nil {
			t.Error(err)
		}else {
			t.Log(bal)
		}
	})
}
