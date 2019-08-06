module github.com/Ankr-network/ankr-chain-cli

go 1.12

replace github.com/tendermint/tendermint => github.com/Ankr-network/tendermint v0.31.5-0.20190719093344-1f8077fcd482

require (
	github.com/Ankr-network/dccn-common v0.0.0-20190729064917-c6a667db8f77
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/tendermint/tendermint v0.32.2
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
)
