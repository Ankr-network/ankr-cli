package cmd

import (
	"encoding/json"
	"fmt"
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
		args := []string{"query", "block", "--nodeurl", localUrl, "--height", "631" }
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

func TestQueryNewFunctions(t *testing.T) {
	args := []string{""}
	cmd := RootCmd
	cmd.SetArgs(args)

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
		args := []string{"query", "transaction", "--nodeurl", localUrl, "--type","Send", "--height", "[34:]"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

var (
	//mock client responses
	blockResultByte = `{
    "block_meta":{
        "block_id":{
            "hash":"0969BDF2ED9AE7FBAA099861CC8B2B23130547EECF7098B3F23B243645BA51F1",
            "parts":{
                "total":1,
                "hash":"0574B03F5E08C4CA593640B8EC668EC2D4A27C3789E47E82E6EF0C9F672C0AED"
            }
        },
        "header":{
            "version":{
                "block":10,
                "app":1
            },
            "chain_id":"test-chain-0bOrck",
            "height":631,
            "time":"2019-07-26T02:00:32.4469545Z",
            "num_txs":16,
            "total_txs":173,
            "last_block_id":{
                "hash":"52D15AD4A3CC76E0D59EC0E6C3B4BAB350EFB9ED4C2C9D092964C00C4B5FAEEC",
                "parts":{
                    "total":1,
                    "hash":"97AEDF088ADBA309A36DB3ED04E642B24995742C41F5D36A379E224DEAE5380F"
                }
            },
            "last_commit_hash":"BE3401D8483E8B823B1A9DCE94F9BA9242ED10F6172C4F1A033B46017CDFCABE",
            "data_hash":"E7EB0D21A277E0D9BFAEF6F03DD159954F1FA48B7D3B5A53CC6DFBE261922EF8",
            "validators_hash":"715E8369F9C411F78E586FF20D656AF2383764873F860538D218166D1DC386F2",
            "next_validators_hash":"715E8369F9C411F78E586FF20D656AF2383764873F860538D218166D1DC386F2",
            "consensus_hash":"048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
            "app_hash":"BA02000000000000",
            "last_results_hash":"",
            "evidence_hash":"",
            "proposer_address":"B3584BE04E33B0F10516EC21BF98F91BEB5B0E1B75DB7B"
        }
    },
    "block":{
        "header":{
            "version":{
                "block":10,
                "app":1
            },
            "chain_id":"test-chain-0bOrck",
            "height":631,
            "time":"2019-07-26T02:00:32.4469545Z",
            "num_txs":16,
            "total_txs":173,
            "last_block_id":{
                "hash":"52D15AD4A3CC76E0D59EC0E6C3B4BAB350EFB9ED4C2C9D092964C00C4B5FAEEC",
                "parts":{
                    "total":1,
                    "hash":"97AEDF088ADBA309A36DB3ED04E642B24995742C41F5D36A379E224DEAE5380F"
                }
            },
            "last_commit_hash":"BE3401D8483E8B823B1A9DCE94F9BA9242ED10F6172C4F1A033B46017CDFCABE",
            "data_hash":"E7EB0D21A277E0D9BFAEF6F03DD159954F1FA48B7D3B5A53CC6DFBE261922EF8",
            "validators_hash":"715E8369F9C411F78E586FF20D656AF2383764873F860538D218166D1DC386F2",
            "next_validators_hash":"715E8369F9C411F78E586FF20D656AF2383764873F860538D218166D1DC386F2",
            "consensus_hash":"048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
            "app_hash":"BA02000000000000",
            "last_results_hash":"",
            "evidence_hash":"",
            "proposer_address":"B3584BE04E33B0F10516EC21BF98F91BEB5B0E1B75DB7B"
        },
        "data":{
            "txs":[
                "dHJ4X3NlbmQ9QzYwNjdDQzU3ODI4MUMyRDgxRjI3NTk4N0E1OUZBOTAyQjIwNUFGRDIyQzFGQjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6cUY5amdCMHh6Y2IwblRqWGNyMTd3bzhGNjlwRSs4S09pSzRUVFl1czMyRT06eVYxaVI4Wngvcm5YVG9KaUVXcithRU5oaXZjWjFtZVhHcVVJZkpFZXVaeXhHR1pjdFBJOFRSajA5eEx1cGlEYWE2U1NTdG1rK29tREtIL3VZQXZVREE9PQ==",
                "dHJ4X3NlbmQ9RTdDOEMwNzFEN0ZFRTI0OTVGQTFGQzU1NEJCN0M2MTU1QzBDMkJBNjQzMTJENzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6bm5sQUgyZFd1TVRaUEhPd2JaZ0ptSXl2MWFScVBTL3k5dmpEU3BzcDdKdz06ckFhM2kzenliazRnUGJDQUZka2dBSGZuRHRjQXI1Zm1TdGZQbnNkaXZMNGdCVVNYYkhkVVY5QzlSUkZHNy9BUWQ3NDlKTDNQT0drcEJoczlBVGtCQ1E9PQ==",
                "dHJ4X3NlbmQ9NjEzRjIxNjQwQTI3RjRGMTU2QUMwMzA0QzlGNzIyOTVFM0VDQzgwMjhEOUZEMTo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6d28zTG1HUUxFYkpZVEtEejlVNUhEQStab1ljUW00N3ROUmhrUmx1Rk92VT06dEJBdWdOeVkvMnlvbDFleVhCbGVGejlBSU1vWWJiQnN3N2dYbkhBR3VaNUU2cS92NDMvSDZ4NUZpdU02OXRlZUNtc0hHT0R0UWNlT3EwZGd3Ly84RGc9PQ==",
                "dHJ4X3NlbmQ9RjQ2NTY5NDlCRDc0NzA1N0E1OURERjkwQTIxOEVDMzUyRTM5MTZBMDk2OTI0RDo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6dDVaSmp0WlFDbDY3QTRVVWt4UWc3SDR2blZBNklGeGZ0OUUvUCtva09lWT06WEN6T1BDWndraGxvZEJyMmNzMUJkZklJdnczamNJSEhNL1FmUENKamNIZ25mRGVBNWVxT3lEb29qeXNWRURrMDQ5NS8vZTdtVTROMlkzd2NQaTVRQnc9PQ==",
                "dHJ4X3NlbmQ9QTM1OTlGQjYyNTE2NzMzMEI0ODYxMDRFMTE2MTExQzMxNTJCMzc1MUE3N0ZCNzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6MWRYRmZaN0J0TkQzdndrVjBBb1V2Rm92RTVrTERjY0dlalZvU0ZTT2ZxUT06YTM5UU1YWkplQjdtaTliK1FLUHBzWm1hUG5GTlZjLzBtUE5NeFI0UnJXSDZvSUg1SjVWcWtKWU8wVW4yeUttSnhwcGdEZU1VRVhKOTJPZlA5LzZpQmc9PQ==",
                "dHJ4X3NlbmQ9NjkwQjg3NTkyMTlBQkNBMjg0RDU3QTczQjA1NjY5NEY3QkQyMEZEQzA2NzE4Mzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6WXRwNi9idlBNZHM2V0NkVE5Ia0kzYTQ2ankzU2E2RDgra0FRSWRsQ1dMdz06aDQ0OWptQ0pnYXdUeVhoL0l5blpNbjdGVEk5QzV3UHE3WjA0RmJSdnZlR1Qzb2RjU2NKaXhSRlFBTWlZR3RMMFBiS1hidDVOcDRvUkVBeGVPVlVZQ0E9PQ==",
                "dHJ4X3NlbmQ9OUVERjA4OTk0MTBEQUNDNjMwRjY4RTc5QzEzQTNGOUIwQTkxOUU3REU2RjdGMjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjY6dDJYaG40Tm9XNlkwbnUxUzl4N2h3YmZ6NTZnUXNwMzhsTDFZMkYzTDAzTT06RDRtUU5ITm5HZzEzM20rYmhKQ2QrRXJaSlh6R1Q4NFpQSEJ0TWo0bVJ6SkpVRktwcEJlSTJ0T3YvRHBKS2V6TEU4R2lrQXZYOFk4WTVFZm1qOG1zQVE9PQ==",
                "dHJ4X3NlbmQ9MUQwNEREN0MyQ0JGMEJBREU5OEIwMkUwRkU4ODRBNEIxRkJCODQxRDUyRkRGNzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6THNyTDdxeXM5UFRPMC9hdW81bnpIVVpPcUxYeUUzRzlVa2FxN3BGSVpKND06b2tuSkMvaStxQks2a1R0blhvMlN1L1pXWHpSTVBkODFzeVh5M1kzVGhhaUlPalNlVXpFT2Z6YmFCNFpxUXNJWm9nSDZPT3hodlNUTUxZY0FlRzQ3Q1E9PQ==",
                "dHJ4X3NlbmQ9ODU3ODI3ODAzQUJDQzc3RkM4OEJEQjA4M0RDN0I5REI1MjFGRUUxRDBBNUIwOTo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjU6UVdSRlVXMkF2c2M2dTdTMDRkRytUdXMrRnFQQVpaeGxWTytEc3ZyUEtnUT06LzY1MC83a2hGd0l2a3RvUXJvM0UrUUdiR2hqa1E1YnZxQm10MFNXUzBkdjNWd0xITkU0a3YxODdBWjBvL2FoWi9WVVhQYnlTRnFDaHJnZjF4cFRjQmc9PQ==",
                "dHJ4X3NlbmQ9M0M2QUNCQUE5QkY4RkJDNTNCOTUwQzU5MDdGRURGMUIxRTY5QzBEQjM0QjdBQjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6RitOUytuTEtxaXczY1ZEdTg1YzV2VTJvSFFQMkJIR0dIeExEWXE0alVqMD06TVk1VXBLMXRJb1V3QzZNS25YQ0hnVTdiYmpETmF5RVZTdmdqVmF3RlBWdDNVM1NOQzJVa1FyMkJ1ZndBOHRkQ3Z6NEtsSlJZb0hYVWFGTVVERjR0QkE9PQ==",
                "dHJ4X3NlbmQ9NTc2MUE2M0VBNTczMkU0OTRCMDBBNUFFNDE2RTMzQUQ2OUNCRjA5OEU2OUI4Mzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6UVV1bTNTRURFaUxURWgrMDRDNndId0JxdUxSeU9ISWdSYTU3S0Iybyt3ND06T2cvdW9zK2VmZTcwZ213VkNMN3dZVGQrSHhhRGNpbzlCZHJwQlk3T3B5QW0zTWdNMWVpMkpQOC9PZDJZUW8wR3NscDh4TWNTNzhIVTBQcmNZMi9zQXc9PQ==",
                "dHJ4X3NlbmQ9NkFBQTU0OEJFQzkyOTE1NDNGOUI2MDc2NUM5REM3NjZEOTNGQ0ZCNEE2MzMwRjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6ZWxvaUhtV0xxWWtsSVQvWXVGSWdyYzYxNWNEUS8zWERmZEVHMmQ5azI2WT06QWF6KzFRcFlIb2phaVZWTWxNbzIwajRRQjlFVTMxNXlXYlIzSTk1QkZBRU13VmUvdm9BNDl3bWN1ejhMM0htSGpmQ3pET0dzc1Y4NW1kT0Y3ZGNWQ2c9PQ==",
                "dHJ4X3NlbmQ9NERDMTE0MDM2QzVENUQ5MThFM0REQTJDRTU2Q0QzNTcyMENDQkExRDk0MkE5Qjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6bDJiU2FpSHZSd2lnVVZRd2pIVi8xNm0wNnQvZ3lGTkpqTWZ2RG8vTFBuWT06YUN3WC9QK2F3dkYwVExJTStYTC9sWlBrOWZBc3Q0VWdEdFVJdzNaZVdvT2xCTExaTE5SZGd4QzBFeGNHbWg4aXlpd0JzQi9SQTk1NFJvUXd5NU9KQ1E9PQ==",
                "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjIxOnd2SEczRWRkQmJYUUhjeUphbDBDUy9ZUWNOWXRFYkZZeGVqbnFmOU9oTTQ9Okc4YTVtT0czR2J6OW9vb3JCQXdoRnFoS1dNemwzS0RST1IzVy9pU1lYTGFGODIxZUlKeERMNHRCQitQK2RkdzNKZkt6K0dKWnJrOWlGVHZHeFh2WUF3PT0=",
                "dHJ4X3NlbmQ9RDZENUExMjkwOTQ0OTc2REMyM0IzMzEyNjZFNURDNEJENDg2NjY1RjcwM0QxRjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjU6T2FZZ1RZNnAwam9XSCt6c3lNK2hHL1U1OHRWZm1POUdwQzJwTGNTVWROcz06Y28vTW1EMSsrcXl3R1ZqWlJsWXoyYVZlSlJpbVNTZVk0aUx4VEhFVHh2bGhpaDBNZE1VZ002QmFNc1VzL3dNOFRWNWlNbkt4R3RMREpPQUpGTGUvRFE9PQ==",
                "dHJ4X3NlbmQ9MDlDMjk1MjQ4MTZCN0I1QzVGQTJDMEFGRDRENDM2REI2QURFNEJBOTg0NkQ3Nzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6dGZwVHgyRUswZUw3bDlKN0FjUFJVZCtpdEpTM0xvdms5dFN0Y2RXRzBXND06bGdIUGZXd0FtWURUeWFhb1BBMm5scE5NdWdWaDJxUWgwdFBHeFFmcHM2RFNFOWxkcmMvM0tPK0xCRGpTTXpNQmd4V09vTWk3VmpNVSs0R09kY0RUQUE9PQ=="
            ]
        },
        "evidence":{
            "evidence":null
        },
        "last_commit":{
            "block_id":{
                "hash":"52D15AD4A3CC76E0D59EC0E6C3B4BAB350EFB9ED4C2C9D092964C00C4B5FAEEC",
                "parts":{
                    "total":1,
                    "hash":"97AEDF088ADBA309A36DB3ED04E642B24995742C41F5D36A379E224DEAE5380F"
                }
            },
            "precommits":[
                {
                    "type":2,
                    "height":630,
                    "round":0,
                    "block_id":{
                        "hash":"52D15AD4A3CC76E0D59EC0E6C3B4BAB350EFB9ED4C2C9D092964C00C4B5FAEEC",
                        "parts":{
                            "total":1,
                            "hash":"97AEDF088ADBA309A36DB3ED04E642B24995742C41F5D36A379E224DEAE5380F"
                        }
                    },
                    "timestamp":"2019-07-26T02:00:32.4469545Z",
                    "validator_address":"B3584BE04E33B0F10516EC21BF98F91BEB5B0E1B75DB7B",
                    "validator_index":0,
                    "signature":"YovsP7hovRva1WBKeg4lGRvMcOyBnWAi/D3nYyJCeu+k8Gx/BiBPJNq1I2dgGCltIWTGcsn3NY8IWYXTzsr7Dw=="
                }
            ]
        }
    }
}`
	txResult = `{
    "hash":"72FB3FA4735E2DE3E56AB50A5D2DDCDBD019012B34A226DCE0B7A3D2E13BDDEB",
    "height":85403,
    "index":0,
    "tx_result":{
        "tags":[
            {
                "key":"YXBwLnR5cGU=",
                "value":"U2V0QmFsYW5jZQ=="
            }
        ]
    },
    "tx":"c2V0X2JhbD05NUNEMDAwMjVDMzgwN0NFRTk4MDREMTlCMUU0MTBBMzBBNDdCMzAzMzcxQzEyOjEyMDAwMDAwMDAwMDAwMDAwMDAwOjM6ZEJDekIrbC9XWXhxaytpNTRhNGFkZHkxWGhpSUs1dDBJQVo1T0t0ZWdXWT06S0EvWGxSaFRSUDJPN0hOdDU4V0pad0t6VW9rdE92WHBBYzU1aWNrM3o4MW96anZrYk9MQTcyMUE1Q3NlUmxibXpEdGx2bWFOMjNJY2VSWkViMnhOQ1E9PQ==",
    "proof":{
        "RootHash":"",
        "Data":null,
        "Proof":{
            "total":0,
            "index":0,
            "leaf_hash":null,
            "aunts":null
        }
    }
}`

	validatorResult = `{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "block_height": "221366",
    "validators": [
      {
        "address": "CAB4B5F7B144C66E51B530AB9F970E43DB2EFC04BCCCDC",
        "pub_key": {
          "type": "tendermint/PubKeyEd25519",
          "value": "4+OV+7egfTTYVb3ZRMBZVLyDhwrAZf9vhItzfHs1M34="
        },
        "voting_power": "10",
        "proposer_priority": "-20"
      },
      {
        "address": "D88484B64890E278E42DAE6245CB76AC193CC78F4321B0",
        "pub_key": {
          "type": "tendermint/PubKeyEd25519",
          "value": "Z4DHutHH+rUzVNN6ovOd8crfdKiz+YqVosF6mdweRc8="
        },
        "voting_power": "10",
        "proposer_priority": "-20"
      },
      {
        "address": "DC2CA6C72E1ECC9D748714F5AFB1EAB516B36AF9D64694",
        "pub_key": {
          "type": "tendermint/PubKeyEd25519",
          "value": "cOShNUweHNXQaiNtbeqr5VbqE7Cl5X7zu80JgWqRjFE="
        },
        "voting_power": "10",
        "proposer_priority": "20"
      },
      {
        "address": "F76BBCCE8E3C764A8CD66780285E88EE3BF45C4017DA28",
        "pub_key": {
          "type": "tendermint/PubKeyEd25519",
          "value": "s1Ho4u5je8Cj0Aj37ZpUbbWY+yLx4B1gxG5HnkPGEWU="
        },
        "voting_power": "10",
        "proposer_priority": "20"
      }
    ]
  }
}`
	statusResult = `{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "node_info": {
      "protocol_version": {
        "p2p": "7",
        "block": "10",
        "app": "1"
      },
      "id": "a3eac426d64cd502d51286d3c41ab88b1d51f62ec1dce7",
      "listen_addr": "tcp://0.0.0.0:26656",
      "network": "Ankr-chain",
      "version": "0.31.5",
      "channels": "4020212223303800",
      "moniker": "dccn-tendermint-58bbbcf6c4-6bgst",
      "other": {
        "tx_index": "on",
        "rpc_address": "tcp://0.0.0.0:26657"
      }
    },
    "sync_info": {
      "latest_block_hash": "B8CA14D9A4BCD516277B5727CB6955DF9F9718F2F151A6CE14028A7A7A979FD4",
      "latest_app_hash": "9C80090000000000",
      "latest_block_height": "221602",
      "latest_block_time": "2019-08-05T08:54:38.608468018Z",
      "catching_up": false
    },
    "validator_info": {
      "address": "18014FC41AA87A8756A270AD0479C4303FD71E19B264CB",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "30xR2KUPfR0uWJN6CjGykgufH4BkXbRKX2eCOJMQhjg="
      },
      "voting_power": "0"
    }
  }
}`
	genesisResult = `{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "genesis": {
      "genesis_time": "2019-02-14T11:04:07.552849Z",
      "chain_id": "Ankr-chain",
      "consensus_params": {
        "block": {
          "max_bytes": "22020096",
          "max_gas": "-1",
          "time_iota_ms": "1000"
        },
        "evidence": {
          "max_age": "100000"
        },
        "validator": {
          "pub_key_types": [
            "ed25519"
          ]
        }
      },
      "validators": [
        {
          "address": "F76BBCCE8E3C764A8CD66780285E88EE3BF45C4017DA28",
          "pub_key": {
            "type": "tendermint/PubKeyEd25519",
            "value": "s1Ho4u5je8Cj0Aj37ZpUbbWY+yLx4B1gxG5HnkPGEWU="
          },
          "power": "10",
          "name": "arthur"
        },
        {
          "address": "DC2CA6C72E1ECC9D748714F5AFB1EAB516B36AF9D64694",
          "pub_key": {
            "type": "tendermint/PubKeyEd25519",
            "value": "cOShNUweHNXQaiNtbeqr5VbqE7Cl5X7zu80JgWqRjFE="
          },
          "power": "10",
          "name": "berkeley"
        },
        {
          "address": "CAB4B5F7B144C66E51B530AB9F970E43DB2EFC04BCCCDC",
          "pub_key": {
            "type": "tendermint/PubKeyEd25519",
            "value": "4+OV+7egfTTYVb3ZRMBZVLyDhwrAZf9vhItzfHs1M34="
          },
          "power": "10",
          "name": "cody"
        },
        {
          "address": "D88484B64890E278E42DAE6245CB76AC193CC78F4321B0",
          "pub_key": {
            "type": "tendermint/PubKeyEd25519",
            "value": "Z4DHutHH+rUzVNN6ovOd8crfdKiz+YqVosF6mdweRc8="
          },
          "power": "10",
          "name": "duke"
        }
      ],
      "app_hash": ""
    }
  }
}`
	consensusStateResult = `{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "round_state": {
      "height/round/step": "221654/0/2",
      "start_time": "2019-08-05T09:04:14.789680982Z",
      "proposal_block_hash": "",
      "locked_block_hash": "",
      "valid_block_hash": "",
      "height_vote_set": [
        {
          "round": "0",
          "prevotes": [
            "nil-Vote",
            "nil-Vote",
            "nil-Vote",
            "nil-Vote"
          ],
          "prevotes_bit_array": "BA{4:____} 0/40 = 0.00",
          "precommits": [
            "nil-Vote",
            "nil-Vote",
            "nil-Vote",
            "nil-Vote"
          ],
          "precommits_bit_array": "BA{4:____} 0/40 = 0.00"
        },
        {
          "round": "1",
          "prevotes": [
            "nil-Vote",
            "nil-Vote",
            "nil-Vote",
            "nil-Vote"
          ],
          "prevotes_bit_array": "BA{4:____} 0/40 = 0.00",
          "precommits": [
            "nil-Vote",
            "nil-Vote",
            "nil-Vote",
            "nil-Vote"
          ],
          "precommits_bit_array": "BA{4:____} 0/40 = 0.00"
        }
      ]
    }
  }
}`
	unconfirmedTxResult = `{
    "jsonrpc":"2.0",
    "id":"jsonrpc-client",
    "result":{
        "n_txs":"16",
        "total":"16",
        "total_bytes":"1024",
        "txs":[
            "dHJ4X3NlbmQ9QzYwNjdDQzU3ODI4MUMyRDgxRjI3NTk4N0E1OUZBOTAyQjIwNUFGRDIyQzFGQjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6cUY5amdCMHh6Y2IwblRqWGNyMTd3bzhGNjlwRSs4S09pSzRUVFl1czMyRT06eVYxaVI4Wngvcm5YVG9KaUVXcithRU5oaXZjWjFtZVhHcVVJZkpFZXVaeXhHR1pjdFBJOFRSajA5eEx1cGlEYWE2U1NTdG1rK29tREtIL3VZQXZVREE9PQ==",
            "dHJ4X3NlbmQ9RTdDOEMwNzFEN0ZFRTI0OTVGQTFGQzU1NEJCN0M2MTU1QzBDMkJBNjQzMTJENzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6bm5sQUgyZFd1TVRaUEhPd2JaZ0ptSXl2MWFScVBTL3k5dmpEU3BzcDdKdz06ckFhM2kzenliazRnUGJDQUZka2dBSGZuRHRjQXI1Zm1TdGZQbnNkaXZMNGdCVVNYYkhkVVY5QzlSUkZHNy9BUWQ3NDlKTDNQT0drcEJoczlBVGtCQ1E9PQ==",
            "dHJ4X3NlbmQ9NjEzRjIxNjQwQTI3RjRGMTU2QUMwMzA0QzlGNzIyOTVFM0VDQzgwMjhEOUZEMTo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6d28zTG1HUUxFYkpZVEtEejlVNUhEQStab1ljUW00N3ROUmhrUmx1Rk92VT06dEJBdWdOeVkvMnlvbDFleVhCbGVGejlBSU1vWWJiQnN3N2dYbkhBR3VaNUU2cS92NDMvSDZ4NUZpdU02OXRlZUNtc0hHT0R0UWNlT3EwZGd3Ly84RGc9PQ==",
            "dHJ4X3NlbmQ9RjQ2NTY5NDlCRDc0NzA1N0E1OURERjkwQTIxOEVDMzUyRTM5MTZBMDk2OTI0RDo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6dDVaSmp0WlFDbDY3QTRVVWt4UWc3SDR2blZBNklGeGZ0OUUvUCtva09lWT06WEN6T1BDWndraGxvZEJyMmNzMUJkZklJdnczamNJSEhNL1FmUENKamNIZ25mRGVBNWVxT3lEb29qeXNWRURrMDQ5NS8vZTdtVTROMlkzd2NQaTVRQnc9PQ==",
            "dHJ4X3NlbmQ9QTM1OTlGQjYyNTE2NzMzMEI0ODYxMDRFMTE2MTExQzMxNTJCMzc1MUE3N0ZCNzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6MWRYRmZaN0J0TkQzdndrVjBBb1V2Rm92RTVrTERjY0dlalZvU0ZTT2ZxUT06YTM5UU1YWkplQjdtaTliK1FLUHBzWm1hUG5GTlZjLzBtUE5NeFI0UnJXSDZvSUg1SjVWcWtKWU8wVW4yeUttSnhwcGdEZU1VRVhKOTJPZlA5LzZpQmc9PQ==",
            "dHJ4X3NlbmQ9NjkwQjg3NTkyMTlBQkNBMjg0RDU3QTczQjA1NjY5NEY3QkQyMEZEQzA2NzE4Mzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6WXRwNi9idlBNZHM2V0NkVE5Ia0kzYTQ2ankzU2E2RDgra0FRSWRsQ1dMdz06aDQ0OWptQ0pnYXdUeVhoL0l5blpNbjdGVEk5QzV3UHE3WjA0RmJSdnZlR1Qzb2RjU2NKaXhSRlFBTWlZR3RMMFBiS1hidDVOcDRvUkVBeGVPVlVZQ0E9PQ==",
            "dHJ4X3NlbmQ9OUVERjA4OTk0MTBEQUNDNjMwRjY4RTc5QzEzQTNGOUIwQTkxOUU3REU2RjdGMjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjY6dDJYaG40Tm9XNlkwbnUxUzl4N2h3YmZ6NTZnUXNwMzhsTDFZMkYzTDAzTT06RDRtUU5ITm5HZzEzM20rYmhKQ2QrRXJaSlh6R1Q4NFpQSEJ0TWo0bVJ6SkpVRktwcEJlSTJ0T3YvRHBKS2V6TEU4R2lrQXZYOFk4WTVFZm1qOG1zQVE9PQ==",
            "dHJ4X3NlbmQ9MUQwNEREN0MyQ0JGMEJBREU5OEIwMkUwRkU4ODRBNEIxRkJCODQxRDUyRkRGNzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6THNyTDdxeXM5UFRPMC9hdW81bnpIVVpPcUxYeUUzRzlVa2FxN3BGSVpKND06b2tuSkMvaStxQks2a1R0blhvMlN1L1pXWHpSTVBkODFzeVh5M1kzVGhhaUlPalNlVXpFT2Z6YmFCNFpxUXNJWm9nSDZPT3hodlNUTUxZY0FlRzQ3Q1E9PQ==",
            "dHJ4X3NlbmQ9ODU3ODI3ODAzQUJDQzc3RkM4OEJEQjA4M0RDN0I5REI1MjFGRUUxRDBBNUIwOTo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjU6UVdSRlVXMkF2c2M2dTdTMDRkRytUdXMrRnFQQVpaeGxWTytEc3ZyUEtnUT06LzY1MC83a2hGd0l2a3RvUXJvM0UrUUdiR2hqa1E1YnZxQm10MFNXUzBkdjNWd0xITkU0a3YxODdBWjBvL2FoWi9WVVhQYnlTRnFDaHJnZjF4cFRjQmc9PQ==",
            "dHJ4X3NlbmQ9M0M2QUNCQUE5QkY4RkJDNTNCOTUwQzU5MDdGRURGMUIxRTY5QzBEQjM0QjdBQjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6RitOUytuTEtxaXczY1ZEdTg1YzV2VTJvSFFQMkJIR0dIeExEWXE0alVqMD06TVk1VXBLMXRJb1V3QzZNS25YQ0hnVTdiYmpETmF5RVZTdmdqVmF3RlBWdDNVM1NOQzJVa1FyMkJ1ZndBOHRkQ3Z6NEtsSlJZb0hYVWFGTVVERjR0QkE9PQ==",
            "dHJ4X3NlbmQ9NTc2MUE2M0VBNTczMkU0OTRCMDBBNUFFNDE2RTMzQUQ2OUNCRjA5OEU2OUI4Mzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6UVV1bTNTRURFaUxURWgrMDRDNndId0JxdUxSeU9ISWdSYTU3S0Iybyt3ND06T2cvdW9zK2VmZTcwZ213VkNMN3dZVGQrSHhhRGNpbzlCZHJwQlk3T3B5QW0zTWdNMWVpMkpQOC9PZDJZUW8wR3NscDh4TWNTNzhIVTBQcmNZMi9zQXc9PQ==",
            "dHJ4X3NlbmQ9NkFBQTU0OEJFQzkyOTE1NDNGOUI2MDc2NUM5REM3NjZEOTNGQ0ZCNEE2MzMwRjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6ZWxvaUhtV0xxWWtsSVQvWXVGSWdyYzYxNWNEUS8zWERmZEVHMmQ5azI2WT06QWF6KzFRcFlIb2phaVZWTWxNbzIwajRRQjlFVTMxNXlXYlIzSTk1QkZBRU13VmUvdm9BNDl3bWN1ejhMM0htSGpmQ3pET0dzc1Y4NW1kT0Y3ZGNWQ2c9PQ==",
            "dHJ4X3NlbmQ9NERDMTE0MDM2QzVENUQ5MThFM0REQTJDRTU2Q0QzNTcyMENDQkExRDk0MkE5Qjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6bDJiU2FpSHZSd2lnVVZRd2pIVi8xNm0wNnQvZ3lGTkpqTWZ2RG8vTFBuWT06YUN3WC9QK2F3dkYwVExJTStYTC9sWlBrOWZBc3Q0VWdEdFVJdzNaZVdvT2xCTExaTE5SZGd4QzBFeGNHbWg4aXlpd0JzQi9SQTk1NFJvUXd5NU9KQ1E9PQ==",
            "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjIxOnd2SEczRWRkQmJYUUhjeUphbDBDUy9ZUWNOWXRFYkZZeGVqbnFmOU9oTTQ9Okc4YTVtT0czR2J6OW9vb3JCQXdoRnFoS1dNemwzS0RST1IzVy9pU1lYTGFGODIxZUlKeERMNHRCQitQK2RkdzNKZkt6K0dKWnJrOWlGVHZHeFh2WUF3PT0=",
            "dHJ4X3NlbmQ9RDZENUExMjkwOTQ0OTc2REMyM0IzMzEyNjZFNURDNEJENDg2NjY1RjcwM0QxRjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjU6T2FZZ1RZNnAwam9XSCt6c3lNK2hHL1U1OHRWZm1POUdwQzJwTGNTVWROcz06Y28vTW1EMSsrcXl3R1ZqWlJsWXoyYVZlSlJpbVNTZVk0aUx4VEhFVHh2bGhpaDBNZE1VZ002QmFNc1VzL3dNOFRWNWlNbkt4R3RMREpPQUpGTGUvRFE9PQ==",
            "dHJ4X3NlbmQ9MDlDMjk1MjQ4MTZCN0I1QzVGQTJDMEFGRDRENDM2REI2QURFNEJBOTg0NkQ3Nzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6dGZwVHgyRUswZUw3bDlKN0FjUFJVZCtpdEpTM0xvdms5dFN0Y2RXRzBXND06bGdIUGZXd0FtWURUeWFhb1BBMm5scE5NdWdWaDJxUWgwdFBHeFFmcHM2RFNFOWxkcmMvM0tPK0xCRGpTTXpNQmd4V09vTWk3VmpNVSs0R09kY0RUQUE9PQ=="
        ]
    }
}`
)