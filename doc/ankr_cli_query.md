## ankr_cli query

query information from ankr chain
### Sub commands

```
  PS D:\> ankr-cli query
  query information from ankr chain
  
  Usage:
    ankr_cli query [command]
  
  Available Commands:
    block              Get block at a given height. If no height is provided, it will fetch the latest block.
    consensusstate     ConsensusState returns a concise summary of the consensus state
    dumpconsensusstate dumps consensus state
    genesis            Get genesis file.
    numunconfirmedtxs  Get number of unconfirmed transactions.
    status             Get Ankr status including node info, pubkey, latest block hash, app hash, block height and time.
    transaction        transaction allows you to query the transaction results.
    unconfirmedtxs     Get unconfirmed transactions (maximum ?limit entries) including their number
    validators         Get the validator set at the given block height. If no height is provided, it will fetch the current validator set.
  
  Flags:
    -h, --help         help for query
        --nodeurl string   validator url
  
  Use "ankr_cli query [command] --help" for more information about a command.
```

### usage  
    global options 
        --nodeurl string       url of a validator 
    * block,  Get block at a given height. If no height is provided, it will fetch the latest block. And you can use "detail" to show more information about transactions contained in block.  
        options: 
            --height int   height interval of the blocks to query. integer or block interval formatted as [from:to] are accepted 
    * consensusstate, get the summary of the consensus state.   
        options: 
            NULL
    * dumpconsensusstate,  dumps consensus state. 
        options: 
            NULL
    * genesis,  Get genesis file. 
        options: 
            NULL
    * numunconfirmedtxs,  Get number of unconfirmed transactions. 
        options: 
            NULL
    * status,  Get Ankr status including node info, pubkey, latest block hash, app hash, block height and time. 
        options: 
            NULL
    * transaction,  transaction allows you to query the transaction results with multiple conditions. 
        options:
            --approve    bool      Include a proof of the transaction inclusion in the block
            --txid       string    The transaction hash
            --creator    string    app creator
            --from       string    the from address contained in a transaction
            --height     string    block height. Input can be an exactly block height  or a height interval separate by ":", and height interval should be enclosed with "[]" or "()" which is mathematically open interval and close interval.
            --metering   string    query metering transaction, both datacenter name and namespace should be  provided and separated  by ":"
            --page       int       Page number (1 based) (default 1)
            --perpage    int       Number of entries per page(max: 100) (default 30)
            --timestamp  string    transaction executed timestamp. Input can be an exactly unix timestamp  or a time interval separate by ":", and time interval should be enclosed with "[]" or "()" which is mathematically open interval and close interval.
            --to         string    the to address contained in a transaction
            --txid       string    The transaction hash
            --type       string    Ankr chain predefined types, SetMetering, SetBalance, UpdatValidator, SetStake, Send
    * unconfirmedtxs,  unconfirmed transactions including their number.
        options: 
            --limit int   number of entries (default 30)
    * validators,  Get the validator set at the given block heigh. 
        options:
            --height int   block height 
### example  
+ block     
    ``` 
    PS D:\> ankr-cli query block --nodeurl http://localhost:26657 --height 631 detail
    
    Block info:
    Version: {10 1}
    Chain-Id: test-chain-0bOrck
    Height: 631
    Time: 2019-07-26 02:00:32.4469545 +0000 UTC
    Number-Txs: 16
    Total-Txs: 173
    Last-block-id: 52D15AD4A3CC76E0D59EC0E6C3B4BAB350EFB9ED4C2C9D092964C00C4B5FAEEC:1:97AEDF088ADB
    Last-commit-hash:BE3401D8483E8B823B1A9DCE94F9BA9242ED10F6172C4F1A033B46017CDFCABE
    Data-hash: E7EB0D21A277E0D9BFAEF6F03DD159954F1FA48B7D3B5A53CC6DFBE261922EF8
    Validator: 715E8369F9C411F78E586FF20D656AF2383764873F860538D218166D1DC386F2
    Consensus: 048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F
    Version: {10 1}
    App-hash: BA02000000000000
    Proposer-Address: B3584BE04E33B0F10516EC21BF98F91BEB5B0E1B75DB7B
    
    Transactions contained in block:
    tx type     hash                                                                  detail
    transfer    0x6c2f1c85ff52f24b37105ce9df71a3a487b2f250a96580b5f879c15725c78a72    from:C6067CC578281C2D81F275987A59FA902B205AFD22C1FB    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:4
    transfer    0xafa157bd6ea99aed5b57cb228246cee2a7820c3d59679c0593995fcc00243f77    from:E7C8C071D7FEE2495FA1FC554BB7C6155C0C2BA64312D7    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:4
    transfer    0x20d8d3f5be299c5f200dae38f48b320ac780289c57352d99c04e3f6fe9eb30ca    from:613F21640A27F4F156AC0304C9F72295E3ECC8028D9FD1    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:4
    transfer    0x545936a089ec2dbacf54cdcae6d0e412982f9a195abd71951ef018090cde0bc3    from:F4656949BD747057A59DDF90A218EC352E3916A096924D    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:3
    transfer    0xff22359fb4ab3d42acbb545c25cb76d6200f4c26d7c67d5775281383ff93ffd1    from:A3599FB625167330B486104E116111C3152B3751A77FB7    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:3
    transfer    0x0c78fff7d6a653fd009f8edf4c6217d1fe9767ddde25b529036cb077b2388e42    from:690B8759219ABCA284D57A73B056694F7BD20FDC067183    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:3
    transfer    0x05a2c0744da1ae74ed864be1d080d156a490f6d9557920cba2e87a1c559aa4b4    from:9EDF0899410DACC630F68E79C13A3F9B0A919E7DE6F7F2    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:6
    transfer    0xeaff3380bd9fe047ef6b6316b54ccf08cd3fe44a94a6013631ff21b1f2d00bc9    from:1D04DD7C2CBF0BADE98B02E0FE884A4B1FBB841D52FDF7    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:3
    transfer    0xc639ea114a575ef963a30942c3067d829cd91243173a299f385284f379d3dea6    from:857827803ABCC77FC88BDB083DC7B9DB521FEE1D0A5B09    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:5
    transfer    0xeddf71155568ff6ddd30738ae2c5bf0f16351a9a5f113993be5ea2f7d4bb1751    from:3C6ACBAA9BF8FBC53B950C5907FEDF1B1E69C0DB34B7AB    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:3
    transfer    0x42b393a27ac102be241b26bf50d75e98308e149c44492365065afbea02437d8c    from:5761A63EA5732E494B00A5AE416E33AD69CBF098E69B83    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:3
    transfer    0x9bb9802487ea47c83d4f9116f639f48e815b039674c2818626d39e3b1894a652    from:6AAA548BEC9291543F9B60765C9DC766D93FCFB4A6330F    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:3
    transfer    0x27e856820dcaa24594daffc292df7c01689cef6de0d147119086c66cd4e92bad    from:4DC114036C5D5D918E3DDA2CE56CD35720CCBA1D942A9B    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:3
    transfer    0x7b843b1551ce61c2c1a4d71be47c7c0e562d00540334a271375fbb51102244c2    from:B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:21
    transfer    0x7030dd72aba5391bf9a33e94a20389787dd9f1490926d12866ed153b97838b45    from:D6D5A1290944976DC23B331266E5DC4BD486665F703D1F    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:5
    transfer    0x2048bb58bcadc4a15efae927096f18fc1843e6ceb8f52761d9c6199fef408af5    from:09C29524816B7B5C5FA2C0AFD4D436DB6ADE4BA9846D77    to:4CEEA87559E1F1C206948ABA668925CBE04FCB9C39EA91    amount:10000000000000000000    nonce:4
    ```
+ consensusstate     
    ``` 
    PS D:\> ankr-cli query consensusstate --nodeurl http://localhost:26657
    {
        "round_state": {
            "height/round/step": "117983/0/1",
            "start_time": "2019-08-01T05:35:21.2209774Z",
            "proposal_block_hash": "",
            "locked_block_hash": "",
            "valid_block_hash": "",
            "height_vote_set": [
                {
                    "round": "0",
                    "prevotes": [
                        "nil-Vote"
                    ],
                    "prevotes_bit_array": "BA{1:_} 0/10 = 0.00",
                    "precommits": [
                        "nil-Vote"
                    ],
                    "precommits_bit_array": "BA{1:_} 0/10 = 0.00"
                }
            ]
        }
    }
    ```
+ dumpconsensusstate     
    ``` 
    PS D:\> ankr-cli query consensusstate --nodeurl http://localhost:26657
    {
        "round_state": {
            "height/round/step": "118015/0/1",
            "start_time": "2019-08-01T05:35:54.1263101Z",
            "proposal_block_hash": "",
            "locked_block_hash": "",
            "valid_block_hash": "",
            "height_vote_set": [
                {
                    "round": "0",
                    "prevotes": [
                        "nil-Vote"
                    ],
                    "prevotes_bit_array": "BA{1:_} 0/10 = 0.00",
                    "precommits": [
                        "nil-Vote"
                    ],
                    "precommits_bit_array": "BA{1:_} 0/10 = 0.00"
                }
            ]
        }
    }
    ```
+ genesis     
    ``` 
    PS D:\> ankr-cli query genesis --nodeurl http://localhost:26657
    Genesis:{
        "genesis_time": 2019-07-24 10:44:03.9174995 +0000 UTC
        "chain_id": test-chain-0bOrck
        "consensus_params":{
        "block": {
            "max_bytes": 22020096,
            "max_gas": -1,
            "time_iota_ms": 1000
        },
        "evidence": {
            "max_age": 100000
        },
        "validator": {
            "pub_key_types": [
                "ed25519"
            ]
        }
    }
        "validators":{
            address                                           pub_key                                                 power    name
            B3584BE04E33B0F10516EC21BF98F91BEB5B0E1B75DB7B    FiTeZAog7fm981JqrKH+tAStvbh8aoHNq0J91llxDsNsQ4EdrHY=    10
        }
    }
    ```
+ numunconfirmedtxs     
    ``` 
    PS D:\> ankr-cli query numunconfirmedtxs --nodeurl http://localhost:26657
    n_tx: 0
    total: 0
    total_bytes: 0
    transactions:
    []
    ```
+ status     
    ``` 
    PS D:\> ankr-cli query status --nodeurl http://localhost:26657
    node_info:{
        "protocol_version": {
            "p2p": 7,
            "block": 10,
            "app": 1
        },
        "id": "26a0d2ecfc37c80a2fb46e64e707874b1803f48aece8cc",
        "listen_addr": "tcp://0.0.0.0:26656",
        "network": "test-chain-0bOrck",
        "version": "0.31.5",
        "channels": "4020212223303800",
        "moniker": "DESKTOP-NF0AS58",
        "other": {
            "tx_index": "on",
            "rpc_address": "tcp://0.0.0.0:26657"
        }
    }
    sync_info:{
        "latest_block_hash": "B61780C55A89BB4632BCCB935AF81718CD1A7F437311B918BE3970826F119D4E",
        "latest_app_hash": "9403000000000000",
        "latest_block_height": 118350,
        "latest_block_time": "2019-08-01T05:41:38.0941923Z",
        "catching_up": false
    }
    validator_info:{
        "address":B3584BE04E33B0F10516EC21BF98F91BEB5B0E1B75DB7B
        "pub_key":FiTeZAog7fm981JqrKH+tAStvbh8aoHNq0J91llxDsNsQ4EdrHY=
        "voting_power":10
    }
    ```
+ transaction     
    ``` 
    PS D:\> ankr-cli query transaction --nodeurl http://localhost:26657 --txid 0x72fb3fa4735e2de3e56ab50a5d2ddcdbd019012b34a226dce0b7a3d2e13bddeb
    tx type        hash                                                                  block height    block index    detail
    set balance    0x72FB3FA4735E2DE3E56AB50A5D2DDCDBD019012B34A226DCE0B7A3D2E13BDDEB    85403           0              address:95CD00025C3807CEE9804D19B1E410A30A47B303371C12    amount:12000000000000000000
    ```
    ```
    PS D:\> ankr-cli query transaction --type Send --nodeurl http://localhost:26657 --from  B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67
    TotalCount:     7
    type           hash                                                                  height    index    detail
    transfer       0x1142188B5FFDD69AB892B47D748406DC4A4C41F7059DDB573639C14DA20701F8    25        0        from: B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67    to:A9963FA874B6B1C94A1401F29630B35298E47F70A2BA65    amount:500000000000000000000000    nonce:2
    transfer       0x4CC298E26A2E8751CD6E6980C5FBB03EE41DAC6265529514362D94BE6C580F78    27        0        from: B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67    to:407279ABCF7AC8AC37764EBB03BC773CFFE875EA58E5EC    amount:250000000000000000000000    nonce:3
    transfer       0xEFFDA65DDEF7DDF44DCE34D5883D255FFA24EFD8FA894B7D7C6ACAABEC3228F0    29        1        from: B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67    to:B9FA454F15A3BFA97A5043331CA415B4CCC1555AA03AC5    amount:125000000000000000000000    nonce:4
    transfer       0x4381EA6E6B3EAA6F07B638391D3347EF03200A8C5FC0DE732142E829A6BDAEEB    31        2        from: B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67    to:A237059A66CE356EC13295CF5DE7F2563B5C7746AF9730    amount:62500000000000000000000     nonce:5
    transfer       0x7D807E82D45960E6952F2812930896B1C676FCB40509CB6AB0C261D719E42A88    35        5        from: B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67    to:DD70352977F0AB2FF89BDE9497A7BAA6DA8D9264688A5C    amount:10000000000000000000        nonce:6
    transfer       0xBDB8CE1966F7755D66BC48FA06E3CCB0FABE1F3E2B6EE7D79BE993DD12D4B67F    37        5        from: B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67    to:DD70352977F0AB2FF89BDE9497A7BAA6DA8D9264688A5C    amount:10000000000000000000        nonce:7
    transfer       0xC8C4340B450213D0969B87DF9DE580BE235C4BAFD69DC34172662653B38614DA    39        10       from: B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67    to:DD70352977F0AB2FF89BDE9497A7BAA6DA8D9264688A5C    amount:10000000000000000000        nonce:8 
    ```
+ unconfirmedtxs     
    ``` 
   PS D:\> ankr-cli query unconfirmedtxs --nodeurl http://localhost:26657
   n_tx: 0
   total: 0
   total_bytes: 0
   transactions:
   []
    ```
+ validators
    ```
     PS D:\> ankr-cli query validators --nodeurl http://localhost:26657
     Height:118491
     
     Validators information:
     Address                                           Pubkey                                                  Voting-Power    Proposer priority
     B3584BE04E33B0F10516EC21BF98F91BEB5B0E1B75DB7B    FiTeZAog7fm981JqrKH+tAStvbh8aoHNq0J91llxDsNsQ4EdrHY=    10              0
    ```
