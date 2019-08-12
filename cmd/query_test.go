package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Ankr-network/ankr-chain-cli/mock_cmd"
	"github.com/agiledragon/gomonkey"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tendermint/tendermint/rpc/client"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	"reflect"
	"testing"
)

var(
	testTxHash = "844D137BD55138EDB05CC83224C47E3D78A7EE3ABAA7EE51DFA3F5203D0E18FC"
	localUrl = "http://localhost:26657"
	remoteUrl = "https://chain-01.dccn.ankr.com:443"
)
func TestQueryBlock(t *testing.T) {
	convey.Convey("test query block", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		//

		var blockResult = &core_types.ResultBlock{}
		err := json.Unmarshal([]byte(blockResultByte),blockResult)
		if err != nil {
			fmt.Println(err)
			t.Error(err)
		}

		//mock method
		var cl = &client.HTTP{}
		clPatch := gomonkey.ApplyMethod(reflect.TypeOf(cl),"Block", func(cl *client.HTTP,height *int64) (*core_types.ResultBlock, error){
			return blockResult, nil
		})
		defer clPatch.Reset()
		args := []string{"query", "block", "--nodeurl", localUrl, "--height", "63", "detail" }
		cmd := RootCmd
		cmd.SetArgs(args)
		err = cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

//72fb3fa4735e2de3e56ab50a5d2ddcdbd019012b34a226dce0b7a3d2e13bddeb
//2048bb58bcadc4a15efae927096f18fc1843e6ceb8f52761d9c6199fef408af5
func TestQueryTxInfo(t *testing.T) {
	convey.Convey("test query transaction", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		var tr = &core_types.ResultTx{}
		err := json.Unmarshal([]byte(txResult),tr)
		if err != nil {
			t.Error(err)
		}
		var cl = &client.HTTP{}
		clPatch := gomonkey.ApplyMethod(reflect.TypeOf(cl),"Tx", func(cl *client.HTTP,hash []byte, prove bool) (*core_types.ResultTx, error){
			return tr, nil
		})
		defer clPatch.Reset()

		args := []string{"query", "transaction", "--nodeurl", "http://localhost:26657", "--txid", "0x72fb3fa4735e2de3e56ab50a5d2ddcdbd019012b34a226dce0b7a3d2e13bddeb"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err = cmd.Execute()
		convey.So(err, convey.ShouldBeNil)

	})
}

func TestQueryValidator(t *testing.T) {
	convey.Convey("test query validator", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock method
		//mock data and marshal from byte
		vr := &core_types.ResultValidators{}
		response := &rpctypes.RPCResponse{}
		err := json.Unmarshal([]byte(validatorResult), response)
		var jsonCl = rpcclient.NewJSONRPCClient(remoteUrl)
		core_types.RegisterAmino(jsonCl.Codec())
		err = jsonCl.Codec().UnmarshalJSON(response.Result, vr)
		convey.So(err, convey.ShouldBeNil)

		var cl = &client.HTTP{}
		clPatch := gomonkey.ApplyMethod(reflect.TypeOf(cl),"Validators", func(cl *client.HTTP, height *int64) (*core_types.ResultValidators, error){
			return vr, nil
		})
		defer clPatch.Reset()

		args := []string{"query", "validators", "--nodeurl", remoteUrl}
		cmd := RootCmd
		cmd.SetArgs(args)
		err = cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestQueryStatus(t *testing.T) {
	convey.Convey("test query status", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock method
		//mock data and marshal from byte
		sr := &core_types.ResultStatus{}
		response := &rpctypes.RPCResponse{}
		err := json.Unmarshal([]byte(statusResult), response)
		var jsonCl = rpcclient.NewJSONRPCClient(remoteUrl)
		core_types.RegisterAmino(jsonCl.Codec())
		err = jsonCl.Codec().UnmarshalJSON(response.Result, sr)
		convey.So(err, convey.ShouldBeNil)

		//stub function
		var cl = &client.HTTP{}
		clPatch := gomonkey.ApplyMethod(reflect.TypeOf(cl),"Status", func(cl *client.HTTP) (*core_types.ResultStatus, error){
			return sr, nil
		})
		defer clPatch.Reset()

		args := []string{"query", "status", "--nodeurl", "https://chain-01.dccn.ankr.com:443"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err = cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestQueryGenesis(t *testing.T)  {
	convey.Convey("test query genesis", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock method
		//mock data and marshal from byte
		gr := &core_types.ResultGenesis{}
		response := &rpctypes.RPCResponse{}
		err := json.Unmarshal([]byte(genesisResult), response)
		var jsonCl = rpcclient.NewJSONRPCClient(remoteUrl)
		core_types.RegisterAmino(jsonCl.Codec())
		err = jsonCl.Codec().UnmarshalJSON(response.Result, gr)
		convey.So(err, convey.ShouldBeNil)

		var cl = &client.HTTP{}
		clPatch := gomonkey.ApplyMethod(reflect.TypeOf(cl),"Genesis", func(cl *client.HTTP) (*core_types.ResultGenesis, error){
			return gr, nil
		})
		defer clPatch.Reset()

		//start test case
		args := []string{"query", "genesis", "--nodeurl", "https://chain-01.dccn.ankr.com:443"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err = cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestQueryConsensusState(t *testing.T)  {
	convey.Convey("test query consensus state", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock method
		//mock data and marshal from byte
		cr := &core_types.ResultConsensusState{}
		response := &rpctypes.RPCResponse{}
		err := json.Unmarshal([]byte(consensusStateResult), response)
		var jsonCl = rpcclient.NewJSONRPCClient(remoteUrl)
		core_types.RegisterAmino(jsonCl.Codec())
		err = jsonCl.Codec().UnmarshalJSON(response.Result, cr)
		convey.So(err, convey.ShouldBeNil)

		var cl = &client.HTTP{}
		clPatch := gomonkey.ApplyMethod(reflect.TypeOf(cl),"ConsensusState", func(cl *client.HTTP) (*core_types.ResultConsensusState, error){
			return cr, nil
		})
		defer clPatch.Reset()

		//start test
		args := []string{"query", "consensusstate", "--nodeurl", "https://chain-01.dccn.ankr.com:443"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err = cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestQueryUnconfirmedTxs(t *testing.T)  {
	convey.Convey("test query unconfirmed transactions", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		//mock method
		//mock data and marshal from byte
		ur := &core_types.ResultUnconfirmedTxs{}
		response := &rpctypes.RPCResponse{}
		err := json.Unmarshal([]byte(unconfirmedTxResult), response)
		var jsonCl = rpcclient.NewJSONRPCClient(remoteUrl)
		core_types.RegisterAmino(jsonCl.Codec())
		err = jsonCl.Codec().UnmarshalJSON(response.Result, ur)
		convey.So(err, convey.ShouldBeNil)

		var cl = &client.HTTP{}
		clPatch := gomonkey.ApplyMethod(reflect.TypeOf(cl),"UnconfirmedTxs", func(cl *client.HTTP,limit int) (*core_types.ResultUnconfirmedTxs, error){
			return ur, nil
		})
		defer clPatch.Reset()

		//start test
		args := []string{"query", "unconfirmedtxs", "--nodeurl", "https://chain-01.dccn.ankr.com:443"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err = cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}


func TestTxSearch(t *testing.T) {
	convey.Convey("test txSearch function", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		txSearchR1 := &core_types.ResultTxSearch{}
		txSearchR2 := &core_types.ResultTxSearch{}
		response := &rpctypes.RPCResponse{}
		var jsonCl = rpcclient.NewJSONRPCClient(localUrl)
		core_types.RegisterAmino(jsonCl.Codec())
		err := json.Unmarshal([]byte(txSearchResult), response)
		err = jsonCl.Codec().UnmarshalJSON(response.Result, txSearchR1)
		err = json.Unmarshal([]byte(txSearchResultMoreCondition), response)
		err = jsonCl.Codec().UnmarshalJSON(response.Result, txSearchR2)
		mockClient := mock_cmd.NewMockClient(ctl)
		//var cl = &client.HTTP{}
		gomock.InOrder(
			mockClient.EXPECT().TxSearch(gomock.Any(),"app.type='Send'",false,1,10).Return(txSearchR1, nil),
			mockClient.EXPECT().TxSearch(gomock.Any(),"app.type='Send' and app.fromaddress='B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67'",false,1,10).Return(txSearchR2, nil),
			)
		var cl = &client.HTTP{}
		patch := gomonkey.ApplyMethod(reflect.TypeOf(cl),"TxSearch", mockClient.TxSearch)
		defer patch.Reset()
		args1 := []string{"query", "transaction", "--nodeurl", localUrl, "--type","Send", "--perpage", "10"}
		args2 := []string{"query", "transaction", "--nodeurl", localUrl, "--type","Send", "--perpage", "10", "--from", "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67"}
		cmd := RootCmd
		cmd.SetArgs(args1)
		err = cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
		cmd.SetArgs(args2)
		err = cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}
