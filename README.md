# overview
ankr chain cli is used to interacting with ankr blockchain. It implements accounts operations like generating new account, exporting private key from keystore or otherwise. 
Admin operations are also provided like send ankr tokens to specified address, etc.sending and querying transactions from ankr chain, etc.


# How to install
ankr-chain-cli is wriitten in Go with support for multiple platforms.   
There are two ways to install ankr chain cli    
1. install ankr-chain-cli using go commands    
```$xslt
go get github.com/Ankr-network/ankr-chain-cli    
cd $GOPATH/src/github.com/Ankr-network/ankr-chain-cli    
dep ensure --vendor-only    
go install github.com/Ankr-network/ankr-chain-cli    
ankr-chain-cli <sub-commands/--help>
```   

2. buid from source    
download and build from the source code  
```
git clone https://github.com/Ankr-network/ankr-chain-cli.git $GOPATH/src/github.com/Ankr-network/ankr-chain-cli    
cd $GOPATH/src/github.com/Ankr-network/ankr-chain-cli    
dep ensure --vendor-only    
go build -o ankr-chain-cli main.go    
./ankr-chain-cli <sub-commands/--help>
```
## help information
```  
   $  ankr-chain-cli
   ankr-chain-cli is used to interacting with ankr blockchain
   
   Usage:
     ankr-chain-cli [command]
   
   Available Commands:
     account     account is used to generate new accounts, encrypt privatekey or decrypt privatekey from keystore
     admin       admin is used to do admin operations
     help        Help about any command
     query       A brief description of your command
     transaction transaction is used to send coins to specified address or send metering
   
   Flags:
     -h, --help   help for ankr_cli
   
   Use "ankr_cli [command] --help" for more information about a command. ```


### SEE ALSO

* [ankr_cli account](ankr_cli_account.md)	 - account is used to generate new accounts, encrypt privatekey or decrypt privatekey from keystore
* [ankr_cli admin](ankr_cli_admin.md)	 - admin is used to do admin operations 
* [ankr_cli query](ankr_cli_query.md)	 - A brief description of your command
* [ankr_cli transaction](ankr_cli_transaction.md)	 - transaction is used to send coins to specified address or send metering
