package cmd

import (
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"io"
)

type Terminal interface {
	ReadPassword(fd int) ([]byte, error)
}

type Fmt interface {
	Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
	Println(a ...interface{}) (n int, err error)
}

type Client interface {
	Genesis() (*core_types.ResultGenesis, error)
	Tx(hash []byte, prove bool) (*core_types.ResultTx, error)
	Block(height *int64) (*core_types.ResultBlock, error)
	Validators(height *int64) (*core_types.ResultValidators, error)
	Status() (*core_types.ResultStatus, error)
}