## ankr_cli query

A brief description of your command
### Sub commands

```
  PS C:\Users\ankr> ankr-chain-cli query
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
        --url string   validator url
  
  Use "ankr_cli query [command] --help" for more information about a command.
```

### usage

