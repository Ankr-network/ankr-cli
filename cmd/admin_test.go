package cmd

import (
	"github.com/Ankr-network/ankr-chain-cli/mock_cmd"
	"github.com/Ankr-network/dccn-common/wallet"
	"github.com/agiledragon/gomonkey"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	adminPrivate = "0mqsOtVueE7uq/I5J/dAhesumWXTu619xXuRgtj4l0d0ELMH6X9ZjGqT6Lnhrhp13LVeGIgrm3QgBnk4q16BZg=="
)
func TestAdmin(t *testing.T) {
	convey.Convey("test set balance", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock wallet interface
		mockWallet := mock_cmd.NewMockWallet(ctl)
		mockWallet.EXPECT().SetBalance("http://localhost","26657", "95CD00025C3807CEE9804D19B1E410A30A47B303371C12","12000000000000000000", adminPrivate).Return(nil)
		mockWallet.EXPECT().GetBalance("http://localhost","26657", "95CD00025C3807CEE9804D19B1E410A30A47B303371C12").Return("12000000000000000000", nil)

		wbPatch := gomonkey.ApplyFunc(wallet.GetBalance, mockWallet.GetBalance)
		defer wbPatch.Reset()
		wsPatch := gomonkey.ApplyFunc(wallet.SetBalance, mockWallet.SetBalance)
		defer wsPatch.Reset()
		//start test case
		args := []string{"admin", "setbalance", "--address", "95CD00025C3807CEE9804D19B1E410A30A47B303371C12","--amount","12000000000000000000" ,"--nodeurl", "http://localhost:26657", "--privkey", adminPrivate}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestSetValidator(t *testing.T) {
	convey.Convey("test setValidator", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock wallet interface
		mockWallet := mock_cmd.NewMockWallet(ctl)
		mockWallet.EXPECT().SetValidator("http://localhost","26657", "FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY=","20", adminPrivate).Return(nil)

		svPatch := gomonkey.ApplyFunc(wallet.SetValidator, mockWallet.SetValidator)
		defer svPatch.Reset()
		//start test case
		args := []string{"admin", "setvalidator", "--pubkey", "FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY=","--power","20" ,"--nodeurl", "http://localhost:26657", "--privkey", adminPrivate}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestSetCert(t *testing.T) {
	convey.Convey("test setCert", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock wallet interface
		mockWallet := mock_cmd.NewMockWallet(ctl)
		mockWallet.EXPECT().SetMeteringCert("http://localhost","26657", adminPrivate,"dc-name", "perm-string").Return(nil)

		smPatch := gomonkey.ApplyFunc(wallet.SetMeteringCert, mockWallet.SetMeteringCert)
		defer smPatch.Reset()
		//start test case
		args := []string{"admin", "setcert", "--dcname", "dc-name","--perm","perm-string" ,"--nodeurl", "http://localhost:26657", "--privkey", adminPrivate}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestRemoveCert(t *testing.T) {
	convey.Convey("test removeCert", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock wallet interface
		mockWallet := mock_cmd.NewMockWallet(ctl)
		mockWallet.EXPECT().RemoveMeteringCert("http://localhost","26657", adminPrivate, "my-dcname").Return(nil)

		wbPatch := gomonkey.ApplyFunc(wallet.RemoveMeteringCert, mockWallet.RemoveMeteringCert)
		defer wbPatch.Reset()

		args := []string{"admin", "removecert", "--dcname", "my-dcname","--nodeurl", "http://localhost:26657", "--privkey", adminPrivate}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestSetStake(t *testing.T) {
	convey.Convey("test setStake", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock wallet interface
		mockWallet := mock_cmd.NewMockWallet(ctl)
		mockWallet.EXPECT().SetStake("http://localhost","26657", adminPrivate,"99", "FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY=").Return(nil)

		wbPatch := gomonkey.ApplyFunc(wallet.SetStake, mockWallet.SetStake)
		defer wbPatch.Reset()
		//start test case
		args := []string{"admin", "setstake", "--pubkey", "FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY=","--amount","99" ,"--nodeurl", "http://localhost:26657", "--privkey", adminPrivate}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestRemoveValidator(t *testing.T) {
	convey.Convey("test removeValidator", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock wallet interface
		mockWallet := mock_cmd.NewMockWallet(ctl)
		mockWallet.EXPECT().RemoveValidator("http://localhost","26657", "FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY=", adminPrivate).Return(nil)

		rmvPatch := gomonkey.ApplyFunc(wallet.RemoveValidator, mockWallet.RemoveValidator)
		defer rmvPatch.Reset()

		//start test case
		args := []string{"admin", "removevalidator", "--pubkey", `FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY=`,"--nodeurl", localUrl, "--privkey",
			adminPrivate}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}