## ankr_cli admin

admin is used to do admin operations including set validator, set cert, set stake, set balance to target address,  remove validator and remove cert    
### Sub commands
```
PS D:\> ankr-cli admin --help
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
      --nodeurl string       url of a validator

Use "ankr_cli admin [command] --help" for more information about a command.
```

### usage
    global options 
        --privkey string   operator private key
        --nodeurl string       url of a validator 
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
    ```
    PS D:\> ankr-cli admin removecert --nodeurl http://localhost:26657 --privkey wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg== --dcname my-dcname
    Remove cert success. 
    ```
+ removevalidator 
    ```
    PS D:\> ankr-cli admin removevalidator --nodeurl http://localhost:26657 --privkey Q5P4l16P+/Cyxq3BvavuWnQPkmeHNYPFkjfuWyQoNyK2vCvT1jyyoh2DYfu+EkWx/hoGjAHOqQw6PMAa7ZkXoQ== --pubkey FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY=
    152CAAFE64D53CEFD677130D08C10A880E5454115754BF0E0270ECA4EF9BB996
    Remove validator success.
    ```
+ setbalance 
    ``` 
    PS D:\> ankr-cli admin setbalance --address E1403CA0DC201F377E820CFA62117A48D4D612400C20D3 --amount 50000000000000
    000000 --nodeurl http://localhost:26657 --privkey 0mqsOtVueE7uq/I5J/dAhesumWXTu619xXuRgtj4l0d0ELMH6X9ZjGqT6Lnhrhp13LVeGIgrm3
    QgBnk4q16BZg==
    Set balance Success.
    Address: E1403CA0DC201F377E820CFA62117A48D4D612400C20D3
    Balance: 50
    ```
+ setcert 
    ```
    PS D:\> ankr-cli admin setcert --nodeurl http://localhost:26657 --privkey wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg== --dcname my-dcname --perm `signature perm`
    c2lnbmF0dXJl
    set_crt=my-dcname:c2lnbmF0dXJl:4:MKR/hOyrYKS85sjl1Je3t4DO358hx0i75NAsjV4ot/dXoo5nGDnUj4tS6KRYyEGiIk1kKL5Hf7fAqDdqb74aAQ==
    Set cert success. 
    ```    
+ setstake 
    ``` 
    PS D:\> ankr-cli admin setstake --nodeurl http://localhost:26657 --privkey Q5P4l16P+/Cyxq3BvavuWnQPkmeHNYPFkjfuWyQoNyK2vCvT1jyyoh2DYfu+EkWx/hoGjAHOqQw6PMAa7ZkXoQ== --pubkey FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY= --amount 99
    Set Stake success.
    ```
+ setvalidator
    ``` 
    PS D:\> ankr-cli admin setvalidator --nodeurl http://localhost:26657 --privkey Q5P4l16P+/Cyxq3BvavuWnQPkmeHNYPFkjfuWyQoNyK2vCvT1jyyoh2DYfu+EkWx/hoGjAHOqQw6PMAa7ZkXoQ== --pubkey FSyq/mTVPO/WdxMNCMEKiA5UVBFXVL8OAnDspO+buZY= --power 3
    152CAAFE64D53CEFD677130D08C10A880E5454115754BF0E0270ECA4EF9BB996
    Set validator success.
    ```
  
