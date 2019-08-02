package cmd

//go:generate mockgen --source=mock.go --destination=../mock_cmd/mock_gen.go
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

type Wallet interface {
	GenerateKeys() (privateKey, pubKey, address string)
	GetBalance(ip, port, address string) (balance string, err error)
	SetStake(ip, port, privKey, amount, pubKey string) error
	RemoveValidator(ip, port, pubKey, privKey string) (err error)
	RemoveMeteringCert(ip, port, privKey, dcName string) error
	SetBalance(ip, port, address, amount, privKey string) error
	SetMeteringCert(ip, port, privKey, dc_name, cert_pem string) error
	SetValidator(ip, port, pubKey, power, privKey string) (err_ret error)
	SendCoins(ip, port, privKey, from, to, amount string) (hash string, err error)
	SetMetering(ip, port, privKey, dc, ns, value string) error
}

type WriteCloser interface {
	Write(p []byte) (n int, err error)
	Close() error
}