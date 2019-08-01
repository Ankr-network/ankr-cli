# Overview
ankr-chain-cli is used to interacting with ankr blockchain.     
ankr-chain-cli implements multiple functions including accounts operations, admin operations, sending different type transactions and query data from ankr blockchain.    

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
## Help information
```  
   $  ankr-chain-cli
   ankr-chain-cli is used to interacting with ankr blockchain
   
   Usage:
     ankr-chain-cli [command]
   
   Available Commands:
     account     account is used to generate new accounts, encrypt privatekey or decrypt privatekey from keystore
     admin       admin is used to do admin operations
     help        Help about any command
     query       query information from ankr chain
     transaction transaction is used to send coins to specified address or send metering
   
   Flags:
     -h, --help   help for ankr_cli
   
   Use "ankr_cli [command] --help" for more information about a command. 
   ```        
# See also

* [ankr_cli account](doc/ankr_cli_account.md)	 - account is used to generate new accounts, encrypt privatekey or decrypt privatekey from keystore
* [ankr_cli admin](doc/ankr_cli_admin.md)	 - admin is used to do admin operations 
* [ankr_cli query](doc/ankr_cli_query.md)	 - A brief description of your command
* [ankr_cli transaction](doc/ankr_cli_transaction.md)	 - transaction is used to send coins to specified address or send metering
