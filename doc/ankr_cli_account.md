## ankr_cli account

account is used to generate new accounts, encrypt privatekey or decrypt privatekey from keystore

### Sub commands

```
  $ ankr-chain-cli account --help
  account is used to generate new accounts, encrypt privatekey or decrypt privatekey from keystore
  
  Usage:
    ankr_cli account [command]
  
  Available Commands:
    exportprivatekey recover private key from keystore.
    genaccount       generate new account.
    genkeystore      generate keystore file based on private key and user input password.
    getbalance       get the balance of an address.
  
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
                --url string <requried> the url of an ankr chain validator            

## examples  
+ exportprivatekey     
    ``` 
    PS C:\Users\ankr_zhang> ankr-chain-cli account exportprivatekey -f .\UTC--2019-07-31T08-41-55.912933900Z--C02CA68DB76A133704DDC3A1CFE84F8C8AC9F666
    
    Please input the keystore password:
    
    Private key exported: O7/HdfLB4/+//7oiwQypnl034j4m39unfFWXONKS+Aoc9dMGg+fBKcesdBNIeDvDP8i5CLCmV18EQTY6ErP5mw==
    ```
+ genaccount     
    ```
    PS C:\Users\ankr> ankr-chain-cli account genaccount -o ./
    
    generating accounts...
    
    Account_0
    private key:  O7/HdfLB4/+//7oiwQypnl034j4m39unfFWXONKS+Aoc9dMGg+fBKcesdBNIeDvDP8i5CLCmV18EQTY6ErP5mw==
    public key:  HPXTBoPnwSnHrHQTSHg7wz/IuQiwpldfBEE2OhKz+Zs=
    address:  C02CA68DB76A133704DDC3A1CFE84F8C8AC9F666
    
    about to export to keystore..
    please input the keystore encryption password:
    please input password again:
    
    exporting to keystore...
    
    created keystore: .//UTC--2019-07-31T08-38-24.681212800Z--C02CA68DB76A133704DDC3A1CFE84F8C8AC9F666
     ```    
+ genkeystore     
    ``` 
    PS C:\Users\ankr> ankr-chain-cli account genkeystore -p O7/HdfLB4/+//7oiwQypnl034j4m39unfFWXONKS+Aoc9dMGg+fBKcesdBNIeDvDP8i5CLCmV18EQTY6ErP5mw==  -o ./
    
    about to export to keystore..
    please input the keystore encryption password:
    please input password again:
    
    exporting to keystore...
    
    created keystore: .//UTC--2019-07-31T08-41-55.912933900Z--C02CA68DB76A133704DDC3A1CFE84F8C8AC9F666
    ```   
+ getbalance    
    ```
     PS C:\Users\ankr> ankr-chain-cli account getbalance --address D6D5A1290944976DC23B331266E5DC4BD486665F703D1F --url http://localhost:26657
     The balance is: 210.898437500000000000
    ```    
    
