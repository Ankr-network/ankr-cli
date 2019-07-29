package cmd

import (
	"encoding/hex"
	"fmt"
	"testing"
)

var(
	testTxHash = "844D137BD55138EDB05CC83224C47E3D78A7EE3ABAA7EE51DFA3F5203D0E18FC"
)
func TestQueryNetInfo(t *testing.T) {
	args := []string{"query", "block", "--url", "http://localhost:26657", "--height", "631" }
	cmd := RootCmd
	cmd.SetArgs(args)
	cmd.Execute()
}

func TestQueryTxInfo(t *testing.T) {
	args := []string{"query", "transaction", "--url", "http://localhost:26657", "--txid", "0x2048bb58bcadc4a15efae927096f18fc1843e6ceb8f52761d9c6199fef408af5"}
	cmd := RootCmd
	cmd.SetArgs(args)
	cmd.Execute()
}

func TestHexString(t *testing.T) {
	str1 := "0x2048bb58bcadc4a15efae927096f18fc1843e6ceb8f52761d9c6199fef408af5"
	var dst []byte
	fmt.Scan(str1, "%x", &dst)
	fmt.Printf("str1 %x", dst)

	fmt.Println()
	str2 := "2048bb58bcadc4a15efae927096f18fc1843e6ceb8f52761d9c6199fef408af5"
	hash, err := hex.DecodeString(str2)
	if err != nil {
		fmt.Println("Decode error!")
	}else {
		fmt.Printf("str2 %x", hash)
	}
}
