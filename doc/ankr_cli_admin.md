## ankr_cli admin

admin is used to do admin operations including set validator, set cert, set stake, set balance to target address,  remove validator and remove cert    
### Sub commands
```
PS C:\Users\ankr> ankr-chain-cli admin --help
admin is used to do admin operations

Usage:
  ankr_cli admin [flags]
  ankr_cli admin [command]

Available Commands:
  removecert      remove cert from validator
  removevalidator remove a validator
  setbalance      set target account with specified amount
  setcert         set metering cert
  setstake        set stake
  setvalidator    add a new validator

Flags:
  -h, --help             help for admin
      --privkey string   operator private key
      --url string       url of a validator

Use "ankr_cli admin [command] --help" for more information about a command.
```

### usage
    global options 
        --privkey string   operator private key
        --url string       url of a validator 
    * removecert, remove cert from validator.  
        options: 
            --dcname string   name of data center name
    * removevalidator, remove a validator.   
        options:
            --pubkey string   public key of the to be removed validator
    * setbalance,  set target account with specified amount.    
        options:
            --address string   the address of the target account to receive ankr token
            --amount string    the amount to set to the target address
    * setcert, set metering cert   
        options:
            --perm string     cert perm to be set
            --dcname string   data center name
    * setstake, set stake   
             options:
                 --amount string   set stake amount
                 --pubkey string   public key
    * setvalidator, add a new validator    
            options:
                --power string    the power set to the validator
                --pubkey string   the public address of the added validator    
                
### example 
+ removecert 
+ removevalidator 
+ setbalance 
``` ankr-chain-cli```
+ setcert 
+ setstake 
+ setvalidator
  
