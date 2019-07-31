package cmd

import (
	"testing"
)

var(
	testTxHash = "844D137BD55138EDB05CC83224C47E3D78A7EE3ABAA7EE51DFA3F5203D0E18FC"
	localUrl = "http://localhost:26657"
	remoteUrl = "https://chain-01.dccn.ankr.com:443"
)
func TestQueryBlock(t *testing.T) {
	args := []string{"query", "block", "--url", "http://localhost:26657", "--height", "631" }
	cmd := RootCmd
	cmd.SetArgs(args)
	cmd.Execute()
}

//72fb3fa4735e2de3e56ab50a5d2ddcdbd019012b34a226dce0b7a3d2e13bddeb
//2048bb58bcadc4a15efae927096f18fc1843e6ceb8f52761d9c6199fef408af5
func TestQueryTxInfo(t *testing.T) {
	args := []string{"query", "transaction", "--url", "http://localhost:26657", "--txid", "0x72fb3fa4735e2de3e56ab50a5d2ddcdbd019012b34a226dce0b7a3d2e13bddeb"}
	cmd := RootCmd
	cmd.SetArgs(args)
	cmd.Execute()
}

//func TestQueryBlockResult(t *testing.T) {
//	args := []string{"query", "blockresult", "--url", localUrl,"--height", "630"}
//	cmd := RootCmd
//	cmd.SetArgs(args)
//	cmd.Execute()
//}

func TestQueryValidator(t *testing.T) {
	args := []string{"query", "validators", "--url", remoteUrl}
	cmd := RootCmd
	cmd.SetArgs(args)
	cmd.Execute()
}

func TestQueryStatus(t *testing.T) {
	args := []string{"query", "status", "--url", "https://chain-01.dccn.ankr.com:443"}
	cmd := RootCmd
	cmd.SetArgs(args)
	cmd.Execute()
}

func TestQueryGenesis(t *testing.T)  {
	args := []string{"query", "genesis", "--url", "https://chain-01.dccn.ankr.com:443"}
	cmd := RootCmd
	cmd.SetArgs(args)
	cmd.Execute()
}

func TestQueryConsensusState(t *testing.T)  {
	args := []string{"query", "consensusstate", "--url", "https://chain-01.dccn.ankr.com:443"}
	cmd := RootCmd
	cmd.SetArgs(args)
	cmd.Execute()
}

func TestQueryDumpConsensus(t *testing.T)  {
	args := []string{"query", "dumpconsensusstate", "--url", "https://chain-01.dccn.ankr.com:443"}
	cmd := RootCmd
	cmd.SetArgs(args)
	cmd.Execute()
}

func TestQueryUnconfirmedTxs(t *testing.T)  {
	args := []string{"query", "unconfirmedtxs", "--url", "https://chain-01.dccn.ankr.com:443"}
	cmd := RootCmd
	cmd.SetArgs(args)
	cmd.Execute()
}

func TestQueryNumUnconfirmed(t *testing.T)  {
	args := []string{"query", "numunconfirmedtxs", "--url", "https://chain-01.dccn.ankr.com:443"}
	cmd := RootCmd
	cmd.SetArgs(args)
	cmd.Execute()
}