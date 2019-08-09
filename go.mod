module github.com/Ankr-network/ankr-chain-cli

go 1.12

replace github.com/tendermint/tendermint => github.com/Ankr-network/tendermint v0.31.5-0.20190719093344-1f8077fcd482

require (
	github.com/Ankr-network/dccn-common v0.0.0-20190729064917-c6a667db8f77
	github.com/btcsuite/btcd v0.0.0-20190807005414-4063feeff79a // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190706150252-9beb055b7962 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tendermint/go-amino v0.15.0 // indirect
	github.com/tendermint/tendermint v0.32.2
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	google.golang.org/genproto v0.0.0-20180831171423-11092d34479b // indirect
)
