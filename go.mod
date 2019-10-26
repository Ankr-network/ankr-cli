module github.com/Ankr-network/ankr-cli

go 1.12

replace github.com/tendermint/tendermint => github.com/Ankr-network/tendermint v0.31.5-0.20190719093344-1f8077fcd482

replace github.com/go-interpreter/wagon => github.com/Ankr-network/wagon v0.9.0-0.20191015152132-a57bd86fecb0

require (
	github.com/Ankr-network/ankr-chain v0.0.0-20191024102123-ec249349c2fc
	github.com/Ankr-network/dccn-common v0.0.0-20191014090437-9fa44d3777fe
	github.com/Ankr-network/dccn-tendermint v0.28.1
	github.com/agiledragon/gomonkey v0.0.0-20190517145658-8fa491f7b918
	github.com/go-interpreter/wagon v0.6.0
	github.com/golang/mock v1.3.1
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/tendermint/tendermint v0.32.6
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/net v0.0.0-20191021144547-ec77196f6094 // indirect
	golang.org/x/sys v0.0.0-20191022100944-742c48ecaeb7 // indirect
	golang.org/x/text v0.3.2 // indirect
)
