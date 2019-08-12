## ankr_cli transaction

transaction is used to send coins to specified address or send metering
### Sub commands

```
PS D:\> ankr-cli transaction --help
transaction is used to send coins to specified address or send metering

Usage:
  ankr_cli transaction [flags]
  ankr_cli transaction [command]

Available Commands:
  metering    send metering transaction
  transfer    send coins to another account

Flags:
  -h, --help         help for transaction
      --nodeurl string   the url of a validator

Use "ankr_cli transaction [command] --help" for more information about a command.
```

### usage

```
    global options 
        --nodeurl string       url of a validator 
    * metering, send metering transaction.  
        options: 
            --dcname string      data center name
            --namespace string   namespace
            --privkey string     admin private key
            --value string       the value to be set
    * transfer, send coins to another account.   
        options: 
            --amount string     amount of ankr token to be transfered
            --keystore string   keystore of the transfer from account
            --to string         receive ankr token address
            
```
### example    
+ metering  
    ```
    PS D:\> ankr-cli transaction metering --nodeurl http://localhost:26657 --dcname datacenter_name --namespace test --value test-value --privkey privkey
    Set metering success.
    ```  
+ transfer    
    ```
    PS D:\> ankr-cli transaction transfer --to F4656949BD747057A59DDF90A218EC352E3916A096924D --amount 20000000000000000000 --keystore .\UTC--2019-08-01T02-25-01.685454800Z--E1403CA0DC201F377E820CFA62117A48D4D612400C20D3 --nodeurl http://localhost:26657
    
    Please input the keystore password:
    
    Transaction sent. Tx hash: 210AEB37AD654AE04CC7A5FC650C23CD4E03A12CC4D2A63A1288D534A8475C31
    ``` 
