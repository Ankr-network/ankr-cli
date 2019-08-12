## ankr_cli account

account is used to generate new accounts, encrypt privatekey or decrypt privatekey from keystore

### Sub commands

```
  PS D:\> ankr-cli account --help
  account is used to generate new accounts, encrypt privatekey or decrypt privatekey from keystore
  
  Usage:
    ankr_cli account [command]
  
  Available Commands:
    exportprivatekey recover private key from keystore.
    genaccount       generate new account.
    genkeystore      generate keystore file based on private key and user input password.
    getbalance       get the balance of an address.
    resetpwd         reset keystore password.
  
  Flags:
    -h, --help   help for account
  
  Use "ankr_cli account [command] --help" for more information about a command.
```

### usage
    * exportprivatekey, recover private key from keystore based on user input password, private key or fail message is printed on the screen.  
        options: 
            -f, --file string <required> the path where keystore is located
            
    * genaccount, generate number of new accounts, private key are printed on the screen and keystore is saved to file after user input password .   
        options:
            -n, --number int [optional] the number of accounts to be generated, default number is 1
            -o, --output string [optional] the path to store keystore, keystore is stored $home/user/AppData/Local/ankr-chain/config/
    * genkeystore,  encrypt privatekey into keystore file.    
        options:
            -p, --privkey string <requried> private key used to generate a keystore file
            -o, --output string [optional] the path to save keystore file
    * getbalance query target account balance    
        options:
            -a, --address string <requried> the address of the target account
                --nodeurl string <requried> the url of an ankr chain validator     
    * resetpwd reset keystore password.    
        options:
            -f, --file string   the path where keystore file is located.      

## examples  
+ exportprivatekey     
    ``` 
    PS D:\> ankr-cli account exportprivatekey -f .\UTC--2019-08-01T02-25-01.685454800Z--E1403CA0DC201F377E820CFA62117A48D4D612400C20D3
    
    Please input the keystore password:
    
    Private key exported: 1gEcOfgXL/rmMvDJPAyL48CanFTeLMU5yASNA9KXmrEVLKr+ZNU879Z3Ew0IwQqIDlRUEVdUvw4CcOyk75u5lg==
    ```
+ genaccount     
    ```
    PS D:\> ankr-cli account genaccount -o ./
    
    generating accounts...
    
    Account_0
    private key:  1gEcOfgXL/rmMvDJPAyL48CanFTeLMU5yASNA9KXmrEVLKr+ZNU879Z3Ew0IwQqIDlRUEVdUvw4CcOyk75u5lg==
    public key:  FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY=
    address:  E1403CA0DC201F377E820CFA62117A48D4D612400C20D3
    
    about to export to keystore..
    please input the keystore encryption password:
    please input password again:
    
    exporting to keystore...
    
    created keystore: .//UTC--2019-08-01T02-25-01.685454800Z--E1403CA0DC201F377E820CFA62117A48D4D612400C20D3
     ```    
+ genkeystore     
    ``` 
    PS D:\> ankr-cli account genkeystore -p 1gEcOfgXL/rmMvDJPAyL48CanFTeLMU5yASNA9KXmrEVLKr+ZNU879Z3Ew0IwQqIDlRUEVdUvw
    4CcOyk75u5lg== -o ./
    
    about to export to keystore..
    please input the keystore encryption password:
    please input password again:
    
    exporting to keystore...
    
    created keystore: .//UTC--2019-08-01T02-35-08.700515700Z--E1403CA0DC201F377E820CFA62117A48D4D612400C20D3
    ```   
+ getbalance    
    ```
     PS D:\> ankr-cli account getbalance --nodeurl http://localhost:26657 --address E1403CA0DC201F377E820CFA62117A48D4D6124
     00C20D3
     The balance is: 50.000000000000000000
    ```    
+ resetpwd    
    ```
    PS D:\> ankr-cli account resetpwd -f ./UTC--2019-08-01T02-34-55.230477100Z--E1403CA0DC201F377E820CFA62117A48D4D612400C20D3
    
    Please input the keystore password:
    please input the keystore encryption password:
    please input password again:
    Password reset success.
    ```    
    
