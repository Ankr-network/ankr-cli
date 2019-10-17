# overview
ankr chain cli is used to interacting with ankr blockchain. It implements accounts operations, admin operations, tranfer operation and query operation.
Account operations are generating new account, exporting private key from keystore or otherwise, and query target address balance. 
Admin operations are sending ankr tokens to specified address using admin privatekey, set validator or metering of ankr-chain.
Transfer operation is used to send different type of transactions.
Query is used to query different kind of data from ankr-chain.


# How to install
ankr-cli is written in Go with support for multiple platforms.   
There are two ways to install ankr-cli, both of them depends on the `go mod` tool, make sure you have already enabled go module before installing ankr-cli.  
Open a terminal and type `export GO111MODULE=on`  to activate go module 
1. install ankr-cli using go commands    
```$xslt
go get github.com/Ankr-network/ankr-cli    
go install github.com/Ankr-network/ankr-cli    
ankr-cli <sub-commands/--help>
```   

2. buid from source    
download and build from the source code  
```
git clone https://github.com/Ankr-network/ankr-cli.git    
cd ankr-cli    
go build -o ankr-cli main.go    
./ankr-cli <sub-commands/--help>
```
## help information
```  
   $  ankr-cli
   ankr-cli is used to interacting with ankr blockchain
   
   Usage:
     ankr-cli [command]
   
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


### SEE ALSO

* [ankr_cli account](doc/ankr_cli_account.md)	 - account is used to generate new accounts, encrypt privatekey or decrypt privatekey from keystore
* [ankr_cli admin](doc/ankr_cli_admin.md)	 - admin is used to do admin operations 
* [ankr_cli query](doc/ankr_cli_query.md)	 - query information from ankr chain
* [ankr_cli transaction](doc/ankr_cli_transaction.md)	 - transaction is used to send coins to specified address or send metering
