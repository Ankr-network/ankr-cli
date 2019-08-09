package cmd

//go:generate mockgen --source=mock.go --destination=../mock_cmd/mock_gen.go

import (
	"github.com/tendermint/tendermint/rpc/client"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"io"
)

type Terminal interface {
	ReadPassword(fd int) ([]byte, error)
}

type Fmt interface {
	Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
	Println(a ...interface{}) (n int, err error)
}

type Client interface {
	Genesis() (*core_types.ResultGenesis, error)
	Tx(hash []byte, prove bool) (*core_types.ResultTx, error)
	Block(height *int64) (*core_types.ResultBlock, error)
	Validators(height *int64) (*core_types.ResultValidators, error)
	Status() (*core_types.ResultStatus, error)
	TxSearch(c *client.HTTP,query string, prove bool, page, perPage int) (*core_types.ResultTxSearch, error)
}

type Wallet interface {
	GenerateKeys() (privateKey, pubKey, address string)
	GetBalance(ip, port, address string) (balance string, err error)
	SetStake(ip, port, privKey, amount, pubKey string) error
	RemoveValidator(ip, port, pubKey, privKey string)  error
	RemoveMeteringCert(ip, port, privKey, dcName string) error
	SetBalance(ip, port, address, amount, privKey string) error
	SetMeteringCert(ip, port, privKey, dc_name, cert_pem string) error
	SetValidator(ip, port, pubKey, power, privKey string) error
	SendCoins(ip, port, privKey, from, to, amount string) (hash string, err error)
	SetMetering(ip, port, privKey, dc, ns, value string) error
}

type WriteCloser interface {
	Write(p []byte) (n int, err error)
	Close() error
}
var (
	//mock client responses
	blockResultByte = `{
    "block_meta":{
        "block_id":{
            "hash":"0969BDF2ED9AE7FBAA099861CC8B2B23130547EECF7098B3F23B243645BA51F1",
            "parts":{
                "total":1,
                "hash":"0574B03F5E08C4CA593640B8EC668EC2D4A27C3789E47E82E6EF0C9F672C0AED"
            }
        },
        "header":{
            "version":{
                "block":10,
                "app":1
            },
            "chain_id":"test-chain-0bOrck",
            "height":631,
            "time":"2019-07-26T02:00:32.4469545Z",
            "num_txs":16,
            "total_txs":173,
            "last_block_id":{
                "hash":"52D15AD4A3CC76E0D59EC0E6C3B4BAB350EFB9ED4C2C9D092964C00C4B5FAEEC",
                "parts":{
                    "total":1,
                    "hash":"97AEDF088ADBA309A36DB3ED04E642B24995742C41F5D36A379E224DEAE5380F"
                }
            },
            "last_commit_hash":"BE3401D8483E8B823B1A9DCE94F9BA9242ED10F6172C4F1A033B46017CDFCABE",
            "data_hash":"E7EB0D21A277E0D9BFAEF6F03DD159954F1FA48B7D3B5A53CC6DFBE261922EF8",
            "validators_hash":"715E8369F9C411F78E586FF20D656AF2383764873F860538D218166D1DC386F2",
            "next_validators_hash":"715E8369F9C411F78E586FF20D656AF2383764873F860538D218166D1DC386F2",
            "consensus_hash":"048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
            "app_hash":"BA02000000000000",
            "last_results_hash":"",
            "evidence_hash":"",
            "proposer_address":"B3584BE04E33B0F10516EC21BF98F91BEB5B0E1B75DB7B"
        }
    },
    "block":{
        "header":{
            "version":{
                "block":10,
                "app":1
            },
            "chain_id":"test-chain-0bOrck",
            "height":631,
            "time":"2019-07-26T02:00:32.4469545Z",
            "num_txs":16,
            "total_txs":173,
            "last_block_id":{
                "hash":"52D15AD4A3CC76E0D59EC0E6C3B4BAB350EFB9ED4C2C9D092964C00C4B5FAEEC",
                "parts":{
                    "total":1,
                    "hash":"97AEDF088ADBA309A36DB3ED04E642B24995742C41F5D36A379E224DEAE5380F"
                }
            },
            "last_commit_hash":"BE3401D8483E8B823B1A9DCE94F9BA9242ED10F6172C4F1A033B46017CDFCABE",
            "data_hash":"E7EB0D21A277E0D9BFAEF6F03DD159954F1FA48B7D3B5A53CC6DFBE261922EF8",
            "validators_hash":"715E8369F9C411F78E586FF20D656AF2383764873F860538D218166D1DC386F2",
            "next_validators_hash":"715E8369F9C411F78E586FF20D656AF2383764873F860538D218166D1DC386F2",
            "consensus_hash":"048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
            "app_hash":"BA02000000000000",
            "last_results_hash":"",
            "evidence_hash":"",
            "proposer_address":"B3584BE04E33B0F10516EC21BF98F91BEB5B0E1B75DB7B"
        },
        "data":{
            "txs":[
                "dHJ4X3NlbmQ9QzYwNjdDQzU3ODI4MUMyRDgxRjI3NTk4N0E1OUZBOTAyQjIwNUFGRDIyQzFGQjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6cUY5amdCMHh6Y2IwblRqWGNyMTd3bzhGNjlwRSs4S09pSzRUVFl1czMyRT06eVYxaVI4Wngvcm5YVG9KaUVXcithRU5oaXZjWjFtZVhHcVVJZkpFZXVaeXhHR1pjdFBJOFRSajA5eEx1cGlEYWE2U1NTdG1rK29tREtIL3VZQXZVREE9PQ==",
                "dHJ4X3NlbmQ9RTdDOEMwNzFEN0ZFRTI0OTVGQTFGQzU1NEJCN0M2MTU1QzBDMkJBNjQzMTJENzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6bm5sQUgyZFd1TVRaUEhPd2JaZ0ptSXl2MWFScVBTL3k5dmpEU3BzcDdKdz06ckFhM2kzenliazRnUGJDQUZka2dBSGZuRHRjQXI1Zm1TdGZQbnNkaXZMNGdCVVNYYkhkVVY5QzlSUkZHNy9BUWQ3NDlKTDNQT0drcEJoczlBVGtCQ1E9PQ==",
                "dHJ4X3NlbmQ9NjEzRjIxNjQwQTI3RjRGMTU2QUMwMzA0QzlGNzIyOTVFM0VDQzgwMjhEOUZEMTo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6d28zTG1HUUxFYkpZVEtEejlVNUhEQStab1ljUW00N3ROUmhrUmx1Rk92VT06dEJBdWdOeVkvMnlvbDFleVhCbGVGejlBSU1vWWJiQnN3N2dYbkhBR3VaNUU2cS92NDMvSDZ4NUZpdU02OXRlZUNtc0hHT0R0UWNlT3EwZGd3Ly84RGc9PQ==",
                "dHJ4X3NlbmQ9RjQ2NTY5NDlCRDc0NzA1N0E1OURERjkwQTIxOEVDMzUyRTM5MTZBMDk2OTI0RDo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6dDVaSmp0WlFDbDY3QTRVVWt4UWc3SDR2blZBNklGeGZ0OUUvUCtva09lWT06WEN6T1BDWndraGxvZEJyMmNzMUJkZklJdnczamNJSEhNL1FmUENKamNIZ25mRGVBNWVxT3lEb29qeXNWRURrMDQ5NS8vZTdtVTROMlkzd2NQaTVRQnc9PQ==",
                "dHJ4X3NlbmQ9QTM1OTlGQjYyNTE2NzMzMEI0ODYxMDRFMTE2MTExQzMxNTJCMzc1MUE3N0ZCNzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6MWRYRmZaN0J0TkQzdndrVjBBb1V2Rm92RTVrTERjY0dlalZvU0ZTT2ZxUT06YTM5UU1YWkplQjdtaTliK1FLUHBzWm1hUG5GTlZjLzBtUE5NeFI0UnJXSDZvSUg1SjVWcWtKWU8wVW4yeUttSnhwcGdEZU1VRVhKOTJPZlA5LzZpQmc9PQ==",
                "dHJ4X3NlbmQ9NjkwQjg3NTkyMTlBQkNBMjg0RDU3QTczQjA1NjY5NEY3QkQyMEZEQzA2NzE4Mzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6WXRwNi9idlBNZHM2V0NkVE5Ia0kzYTQ2ankzU2E2RDgra0FRSWRsQ1dMdz06aDQ0OWptQ0pnYXdUeVhoL0l5blpNbjdGVEk5QzV3UHE3WjA0RmJSdnZlR1Qzb2RjU2NKaXhSRlFBTWlZR3RMMFBiS1hidDVOcDRvUkVBeGVPVlVZQ0E9PQ==",
                "dHJ4X3NlbmQ9OUVERjA4OTk0MTBEQUNDNjMwRjY4RTc5QzEzQTNGOUIwQTkxOUU3REU2RjdGMjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjY6dDJYaG40Tm9XNlkwbnUxUzl4N2h3YmZ6NTZnUXNwMzhsTDFZMkYzTDAzTT06RDRtUU5ITm5HZzEzM20rYmhKQ2QrRXJaSlh6R1Q4NFpQSEJ0TWo0bVJ6SkpVRktwcEJlSTJ0T3YvRHBKS2V6TEU4R2lrQXZYOFk4WTVFZm1qOG1zQVE9PQ==",
                "dHJ4X3NlbmQ9MUQwNEREN0MyQ0JGMEJBREU5OEIwMkUwRkU4ODRBNEIxRkJCODQxRDUyRkRGNzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6THNyTDdxeXM5UFRPMC9hdW81bnpIVVpPcUxYeUUzRzlVa2FxN3BGSVpKND06b2tuSkMvaStxQks2a1R0blhvMlN1L1pXWHpSTVBkODFzeVh5M1kzVGhhaUlPalNlVXpFT2Z6YmFCNFpxUXNJWm9nSDZPT3hodlNUTUxZY0FlRzQ3Q1E9PQ==",
                "dHJ4X3NlbmQ9ODU3ODI3ODAzQUJDQzc3RkM4OEJEQjA4M0RDN0I5REI1MjFGRUUxRDBBNUIwOTo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjU6UVdSRlVXMkF2c2M2dTdTMDRkRytUdXMrRnFQQVpaeGxWTytEc3ZyUEtnUT06LzY1MC83a2hGd0l2a3RvUXJvM0UrUUdiR2hqa1E1YnZxQm10MFNXUzBkdjNWd0xITkU0a3YxODdBWjBvL2FoWi9WVVhQYnlTRnFDaHJnZjF4cFRjQmc9PQ==",
                "dHJ4X3NlbmQ9M0M2QUNCQUE5QkY4RkJDNTNCOTUwQzU5MDdGRURGMUIxRTY5QzBEQjM0QjdBQjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6RitOUytuTEtxaXczY1ZEdTg1YzV2VTJvSFFQMkJIR0dIeExEWXE0alVqMD06TVk1VXBLMXRJb1V3QzZNS25YQ0hnVTdiYmpETmF5RVZTdmdqVmF3RlBWdDNVM1NOQzJVa1FyMkJ1ZndBOHRkQ3Z6NEtsSlJZb0hYVWFGTVVERjR0QkE9PQ==",
                "dHJ4X3NlbmQ9NTc2MUE2M0VBNTczMkU0OTRCMDBBNUFFNDE2RTMzQUQ2OUNCRjA5OEU2OUI4Mzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6UVV1bTNTRURFaUxURWgrMDRDNndId0JxdUxSeU9ISWdSYTU3S0Iybyt3ND06T2cvdW9zK2VmZTcwZ213VkNMN3dZVGQrSHhhRGNpbzlCZHJwQlk3T3B5QW0zTWdNMWVpMkpQOC9PZDJZUW8wR3NscDh4TWNTNzhIVTBQcmNZMi9zQXc9PQ==",
                "dHJ4X3NlbmQ9NkFBQTU0OEJFQzkyOTE1NDNGOUI2MDc2NUM5REM3NjZEOTNGQ0ZCNEE2MzMwRjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6ZWxvaUhtV0xxWWtsSVQvWXVGSWdyYzYxNWNEUS8zWERmZEVHMmQ5azI2WT06QWF6KzFRcFlIb2phaVZWTWxNbzIwajRRQjlFVTMxNXlXYlIzSTk1QkZBRU13VmUvdm9BNDl3bWN1ejhMM0htSGpmQ3pET0dzc1Y4NW1kT0Y3ZGNWQ2c9PQ==",
                "dHJ4X3NlbmQ9NERDMTE0MDM2QzVENUQ5MThFM0REQTJDRTU2Q0QzNTcyMENDQkExRDk0MkE5Qjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6bDJiU2FpSHZSd2lnVVZRd2pIVi8xNm0wNnQvZ3lGTkpqTWZ2RG8vTFBuWT06YUN3WC9QK2F3dkYwVExJTStYTC9sWlBrOWZBc3Q0VWdEdFVJdzNaZVdvT2xCTExaTE5SZGd4QzBFeGNHbWg4aXlpd0JzQi9SQTk1NFJvUXd5NU9KQ1E9PQ==",
                "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjIxOnd2SEczRWRkQmJYUUhjeUphbDBDUy9ZUWNOWXRFYkZZeGVqbnFmOU9oTTQ9Okc4YTVtT0czR2J6OW9vb3JCQXdoRnFoS1dNemwzS0RST1IzVy9pU1lYTGFGODIxZUlKeERMNHRCQitQK2RkdzNKZkt6K0dKWnJrOWlGVHZHeFh2WUF3PT0=",
                "dHJ4X3NlbmQ9RDZENUExMjkwOTQ0OTc2REMyM0IzMzEyNjZFNURDNEJENDg2NjY1RjcwM0QxRjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjU6T2FZZ1RZNnAwam9XSCt6c3lNK2hHL1U1OHRWZm1POUdwQzJwTGNTVWROcz06Y28vTW1EMSsrcXl3R1ZqWlJsWXoyYVZlSlJpbVNTZVk0aUx4VEhFVHh2bGhpaDBNZE1VZ002QmFNc1VzL3dNOFRWNWlNbkt4R3RMREpPQUpGTGUvRFE9PQ==",
                "dHJ4X3NlbmQ9MDlDMjk1MjQ4MTZCN0I1QzVGQTJDMEFGRDRENDM2REI2QURFNEJBOTg0NkQ3Nzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6dGZwVHgyRUswZUw3bDlKN0FjUFJVZCtpdEpTM0xvdms5dFN0Y2RXRzBXND06bGdIUGZXd0FtWURUeWFhb1BBMm5scE5NdWdWaDJxUWgwdFBHeFFmcHM2RFNFOWxkcmMvM0tPK0xCRGpTTXpNQmd4V09vTWk3VmpNVSs0R09kY0RUQUE9PQ=="
            ]
        },
        "evidence":{
            "evidence":null
        },
        "last_commit":{
            "block_id":{
                "hash":"52D15AD4A3CC76E0D59EC0E6C3B4BAB350EFB9ED4C2C9D092964C00C4B5FAEEC",
                "parts":{
                    "total":1,
                    "hash":"97AEDF088ADBA309A36DB3ED04E642B24995742C41F5D36A379E224DEAE5380F"
                }
            },
            "precommits":[
                {
                    "type":2,
                    "height":630,
                    "round":0,
                    "block_id":{
                        "hash":"52D15AD4A3CC76E0D59EC0E6C3B4BAB350EFB9ED4C2C9D092964C00C4B5FAEEC",
                        "parts":{
                            "total":1,
                            "hash":"97AEDF088ADBA309A36DB3ED04E642B24995742C41F5D36A379E224DEAE5380F"
                        }
                    },
                    "timestamp":"2019-07-26T02:00:32.4469545Z",
                    "validator_address":"B3584BE04E33B0F10516EC21BF98F91BEB5B0E1B75DB7B",
                    "validator_index":0,
                    "signature":"YovsP7hovRva1WBKeg4lGRvMcOyBnWAi/D3nYyJCeu+k8Gx/BiBPJNq1I2dgGCltIWTGcsn3NY8IWYXTzsr7Dw=="
                }
            ]
        }
    }
}`
	txResult = `{
    "hash":"72FB3FA4735E2DE3E56AB50A5D2DDCDBD019012B34A226DCE0B7A3D2E13BDDEB",
    "height":85403,
    "index":0,
    "tx_result":{
        "tags":[
            {
                "key":"YXBwLnR5cGU=",
                "value":"U2V0QmFsYW5jZQ=="
            }
        ]
    },
    "tx":"c2V0X2JhbD05NUNEMDAwMjVDMzgwN0NFRTk4MDREMTlCMUU0MTBBMzBBNDdCMzAzMzcxQzEyOjEyMDAwMDAwMDAwMDAwMDAwMDAwOjM6ZEJDekIrbC9XWXhxaytpNTRhNGFkZHkxWGhpSUs1dDBJQVo1T0t0ZWdXWT06S0EvWGxSaFRSUDJPN0hOdDU4V0pad0t6VW9rdE92WHBBYzU1aWNrM3o4MW96anZrYk9MQTcyMUE1Q3NlUmxibXpEdGx2bWFOMjNJY2VSWkViMnhOQ1E9PQ==",
    "proof":{
        "RootHash":"",
        "Data":null,
        "Proof":{
            "total":0,
            "index":0,
            "leaf_hash":null,
            "aunts":null
        }
    }
}`

	validatorResult = `{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "block_height": "221366",
    "validators": [
      {
        "address": "CAB4B5F7B144C66E51B530AB9F970E43DB2EFC04BCCCDC",
        "pub_key": {
          "type": "tendermint/PubKeyEd25519",
          "value": "4+OV+7egfTTYVb3ZRMBZVLyDhwrAZf9vhItzfHs1M34="
        },
        "voting_power": "10",
        "proposer_priority": "-20"
      },
      {
        "address": "D88484B64890E278E42DAE6245CB76AC193CC78F4321B0",
        "pub_key": {
          "type": "tendermint/PubKeyEd25519",
          "value": "Z4DHutHH+rUzVNN6ovOd8crfdKiz+YqVosF6mdweRc8="
        },
        "voting_power": "10",
        "proposer_priority": "-20"
      },
      {
        "address": "DC2CA6C72E1ECC9D748714F5AFB1EAB516B36AF9D64694",
        "pub_key": {
          "type": "tendermint/PubKeyEd25519",
          "value": "cOShNUweHNXQaiNtbeqr5VbqE7Cl5X7zu80JgWqRjFE="
        },
        "voting_power": "10",
        "proposer_priority": "20"
      },
      {
        "address": "F76BBCCE8E3C764A8CD66780285E88EE3BF45C4017DA28",
        "pub_key": {
          "type": "tendermint/PubKeyEd25519",
          "value": "s1Ho4u5je8Cj0Aj37ZpUbbWY+yLx4B1gxG5HnkPGEWU="
        },
        "voting_power": "10",
        "proposer_priority": "20"
      }
    ]
  }
}`
	statusResult = `{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "node_info": {
      "protocol_version": {
        "p2p": "7",
        "block": "10",
        "app": "1"
      },
      "id": "a3eac426d64cd502d51286d3c41ab88b1d51f62ec1dce7",
      "listen_addr": "tcp://0.0.0.0:26656",
      "network": "Ankr-chain",
      "version": "0.31.5",
      "channels": "4020212223303800",
      "moniker": "dccn-tendermint-58bbbcf6c4-6bgst",
      "other": {
        "tx_index": "on",
        "rpc_address": "tcp://0.0.0.0:26657"
      }
    },
    "sync_info": {
      "latest_block_hash": "B8CA14D9A4BCD516277B5727CB6955DF9F9718F2F151A6CE14028A7A7A979FD4",
      "latest_app_hash": "9C80090000000000",
      "latest_block_height": "221602",
      "latest_block_time": "2019-08-05T08:54:38.608468018Z",
      "catching_up": false
    },
    "validator_info": {
      "address": "18014FC41AA87A8756A270AD0479C4303FD71E19B264CB",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "30xR2KUPfR0uWJN6CjGykgufH4BkXbRKX2eCOJMQhjg="
      },
      "voting_power": "0"
    }
  }
}`
	genesisResult = `{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "genesis": {
      "genesis_time": "2019-02-14T11:04:07.552849Z",
      "chain_id": "Ankr-chain",
      "consensus_params": {
        "block": {
          "max_bytes": "22020096",
          "max_gas": "-1",
          "time_iota_ms": "1000"
        },
        "evidence": {
          "max_age": "100000"
        },
        "validator": {
          "pub_key_types": [
            "ed25519"
          ]
        }
      },
      "validators": [
        {
          "address": "F76BBCCE8E3C764A8CD66780285E88EE3BF45C4017DA28",
          "pub_key": {
            "type": "tendermint/PubKeyEd25519",
            "value": "s1Ho4u5je8Cj0Aj37ZpUbbWY+yLx4B1gxG5HnkPGEWU="
          },
          "power": "10",
          "name": "arthur"
        },
        {
          "address": "DC2CA6C72E1ECC9D748714F5AFB1EAB516B36AF9D64694",
          "pub_key": {
            "type": "tendermint/PubKeyEd25519",
            "value": "cOShNUweHNXQaiNtbeqr5VbqE7Cl5X7zu80JgWqRjFE="
          },
          "power": "10",
          "name": "berkeley"
        },
        {
          "address": "CAB4B5F7B144C66E51B530AB9F970E43DB2EFC04BCCCDC",
          "pub_key": {
            "type": "tendermint/PubKeyEd25519",
            "value": "4+OV+7egfTTYVb3ZRMBZVLyDhwrAZf9vhItzfHs1M34="
          },
          "power": "10",
          "name": "cody"
        },
        {
          "address": "D88484B64890E278E42DAE6245CB76AC193CC78F4321B0",
          "pub_key": {
            "type": "tendermint/PubKeyEd25519",
            "value": "Z4DHutHH+rUzVNN6ovOd8crfdKiz+YqVosF6mdweRc8="
          },
          "power": "10",
          "name": "duke"
        }
      ],
      "app_hash": ""
    }
  }
}`
	consensusStateResult = `{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "round_state": {
      "height/round/step": "221654/0/2",
      "start_time": "2019-08-05T09:04:14.789680982Z",
      "proposal_block_hash": "",
      "locked_block_hash": "",
      "valid_block_hash": "",
      "height_vote_set": [
        {
          "round": "0",
          "prevotes": [
            "nil-Vote",
            "nil-Vote",
            "nil-Vote",
            "nil-Vote"
          ],
          "prevotes_bit_array": "BA{4:____} 0/40 = 0.00",
          "precommits": [
            "nil-Vote",
            "nil-Vote",
            "nil-Vote",
            "nil-Vote"
          ],
          "precommits_bit_array": "BA{4:____} 0/40 = 0.00"
        },
        {
          "round": "1",
          "prevotes": [
            "nil-Vote",
            "nil-Vote",
            "nil-Vote",
            "nil-Vote"
          ],
          "prevotes_bit_array": "BA{4:____} 0/40 = 0.00",
          "precommits": [
            "nil-Vote",
            "nil-Vote",
            "nil-Vote",
            "nil-Vote"
          ],
          "precommits_bit_array": "BA{4:____} 0/40 = 0.00"
        }
      ]
    }
  }
}`
	unconfirmedTxResult = `{
    "jsonrpc":"2.0",
    "id":"jsonrpc-client",
    "result":{
        "n_txs":"16",
        "total":"16",
        "total_bytes":"1024",
        "txs":[
            "dHJ4X3NlbmQ9QzYwNjdDQzU3ODI4MUMyRDgxRjI3NTk4N0E1OUZBOTAyQjIwNUFGRDIyQzFGQjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6cUY5amdCMHh6Y2IwblRqWGNyMTd3bzhGNjlwRSs4S09pSzRUVFl1czMyRT06eVYxaVI4Wngvcm5YVG9KaUVXcithRU5oaXZjWjFtZVhHcVVJZkpFZXVaeXhHR1pjdFBJOFRSajA5eEx1cGlEYWE2U1NTdG1rK29tREtIL3VZQXZVREE9PQ==",
            "dHJ4X3NlbmQ9RTdDOEMwNzFEN0ZFRTI0OTVGQTFGQzU1NEJCN0M2MTU1QzBDMkJBNjQzMTJENzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6bm5sQUgyZFd1TVRaUEhPd2JaZ0ptSXl2MWFScVBTL3k5dmpEU3BzcDdKdz06ckFhM2kzenliazRnUGJDQUZka2dBSGZuRHRjQXI1Zm1TdGZQbnNkaXZMNGdCVVNYYkhkVVY5QzlSUkZHNy9BUWQ3NDlKTDNQT0drcEJoczlBVGtCQ1E9PQ==",
            "dHJ4X3NlbmQ9NjEzRjIxNjQwQTI3RjRGMTU2QUMwMzA0QzlGNzIyOTVFM0VDQzgwMjhEOUZEMTo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6d28zTG1HUUxFYkpZVEtEejlVNUhEQStab1ljUW00N3ROUmhrUmx1Rk92VT06dEJBdWdOeVkvMnlvbDFleVhCbGVGejlBSU1vWWJiQnN3N2dYbkhBR3VaNUU2cS92NDMvSDZ4NUZpdU02OXRlZUNtc0hHT0R0UWNlT3EwZGd3Ly84RGc9PQ==",
            "dHJ4X3NlbmQ9RjQ2NTY5NDlCRDc0NzA1N0E1OURERjkwQTIxOEVDMzUyRTM5MTZBMDk2OTI0RDo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6dDVaSmp0WlFDbDY3QTRVVWt4UWc3SDR2blZBNklGeGZ0OUUvUCtva09lWT06WEN6T1BDWndraGxvZEJyMmNzMUJkZklJdnczamNJSEhNL1FmUENKamNIZ25mRGVBNWVxT3lEb29qeXNWRURrMDQ5NS8vZTdtVTROMlkzd2NQaTVRQnc9PQ==",
            "dHJ4X3NlbmQ9QTM1OTlGQjYyNTE2NzMzMEI0ODYxMDRFMTE2MTExQzMxNTJCMzc1MUE3N0ZCNzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6MWRYRmZaN0J0TkQzdndrVjBBb1V2Rm92RTVrTERjY0dlalZvU0ZTT2ZxUT06YTM5UU1YWkplQjdtaTliK1FLUHBzWm1hUG5GTlZjLzBtUE5NeFI0UnJXSDZvSUg1SjVWcWtKWU8wVW4yeUttSnhwcGdEZU1VRVhKOTJPZlA5LzZpQmc9PQ==",
            "dHJ4X3NlbmQ9NjkwQjg3NTkyMTlBQkNBMjg0RDU3QTczQjA1NjY5NEY3QkQyMEZEQzA2NzE4Mzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6WXRwNi9idlBNZHM2V0NkVE5Ia0kzYTQ2ankzU2E2RDgra0FRSWRsQ1dMdz06aDQ0OWptQ0pnYXdUeVhoL0l5blpNbjdGVEk5QzV3UHE3WjA0RmJSdnZlR1Qzb2RjU2NKaXhSRlFBTWlZR3RMMFBiS1hidDVOcDRvUkVBeGVPVlVZQ0E9PQ==",
            "dHJ4X3NlbmQ9OUVERjA4OTk0MTBEQUNDNjMwRjY4RTc5QzEzQTNGOUIwQTkxOUU3REU2RjdGMjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjY6dDJYaG40Tm9XNlkwbnUxUzl4N2h3YmZ6NTZnUXNwMzhsTDFZMkYzTDAzTT06RDRtUU5ITm5HZzEzM20rYmhKQ2QrRXJaSlh6R1Q4NFpQSEJ0TWo0bVJ6SkpVRktwcEJlSTJ0T3YvRHBKS2V6TEU4R2lrQXZYOFk4WTVFZm1qOG1zQVE9PQ==",
            "dHJ4X3NlbmQ9MUQwNEREN0MyQ0JGMEJBREU5OEIwMkUwRkU4ODRBNEIxRkJCODQxRDUyRkRGNzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6THNyTDdxeXM5UFRPMC9hdW81bnpIVVpPcUxYeUUzRzlVa2FxN3BGSVpKND06b2tuSkMvaStxQks2a1R0blhvMlN1L1pXWHpSTVBkODFzeVh5M1kzVGhhaUlPalNlVXpFT2Z6YmFCNFpxUXNJWm9nSDZPT3hodlNUTUxZY0FlRzQ3Q1E9PQ==",
            "dHJ4X3NlbmQ9ODU3ODI3ODAzQUJDQzc3RkM4OEJEQjA4M0RDN0I5REI1MjFGRUUxRDBBNUIwOTo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjU6UVdSRlVXMkF2c2M2dTdTMDRkRytUdXMrRnFQQVpaeGxWTytEc3ZyUEtnUT06LzY1MC83a2hGd0l2a3RvUXJvM0UrUUdiR2hqa1E1YnZxQm10MFNXUzBkdjNWd0xITkU0a3YxODdBWjBvL2FoWi9WVVhQYnlTRnFDaHJnZjF4cFRjQmc9PQ==",
            "dHJ4X3NlbmQ9M0M2QUNCQUE5QkY4RkJDNTNCOTUwQzU5MDdGRURGMUIxRTY5QzBEQjM0QjdBQjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6RitOUytuTEtxaXczY1ZEdTg1YzV2VTJvSFFQMkJIR0dIeExEWXE0alVqMD06TVk1VXBLMXRJb1V3QzZNS25YQ0hnVTdiYmpETmF5RVZTdmdqVmF3RlBWdDNVM1NOQzJVa1FyMkJ1ZndBOHRkQ3Z6NEtsSlJZb0hYVWFGTVVERjR0QkE9PQ==",
            "dHJ4X3NlbmQ9NTc2MUE2M0VBNTczMkU0OTRCMDBBNUFFNDE2RTMzQUQ2OUNCRjA5OEU2OUI4Mzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6UVV1bTNTRURFaUxURWgrMDRDNndId0JxdUxSeU9ISWdSYTU3S0Iybyt3ND06T2cvdW9zK2VmZTcwZ213VkNMN3dZVGQrSHhhRGNpbzlCZHJwQlk3T3B5QW0zTWdNMWVpMkpQOC9PZDJZUW8wR3NscDh4TWNTNzhIVTBQcmNZMi9zQXc9PQ==",
            "dHJ4X3NlbmQ9NkFBQTU0OEJFQzkyOTE1NDNGOUI2MDc2NUM5REM3NjZEOTNGQ0ZCNEE2MzMwRjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6ZWxvaUhtV0xxWWtsSVQvWXVGSWdyYzYxNWNEUS8zWERmZEVHMmQ5azI2WT06QWF6KzFRcFlIb2phaVZWTWxNbzIwajRRQjlFVTMxNXlXYlIzSTk1QkZBRU13VmUvdm9BNDl3bWN1ejhMM0htSGpmQ3pET0dzc1Y4NW1kT0Y3ZGNWQ2c9PQ==",
            "dHJ4X3NlbmQ9NERDMTE0MDM2QzVENUQ5MThFM0REQTJDRTU2Q0QzNTcyMENDQkExRDk0MkE5Qjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjM6bDJiU2FpSHZSd2lnVVZRd2pIVi8xNm0wNnQvZ3lGTkpqTWZ2RG8vTFBuWT06YUN3WC9QK2F3dkYwVExJTStYTC9sWlBrOWZBc3Q0VWdEdFVJdzNaZVdvT2xCTExaTE5SZGd4QzBFeGNHbWg4aXlpd0JzQi9SQTk1NFJvUXd5NU9KQ1E9PQ==",
            "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjIxOnd2SEczRWRkQmJYUUhjeUphbDBDUy9ZUWNOWXRFYkZZeGVqbnFmOU9oTTQ9Okc4YTVtT0czR2J6OW9vb3JCQXdoRnFoS1dNemwzS0RST1IzVy9pU1lYTGFGODIxZUlKeERMNHRCQitQK2RkdzNKZkt6K0dKWnJrOWlGVHZHeFh2WUF3PT0=",
            "dHJ4X3NlbmQ9RDZENUExMjkwOTQ0OTc2REMyM0IzMzEyNjZFNURDNEJENDg2NjY1RjcwM0QxRjo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjU6T2FZZ1RZNnAwam9XSCt6c3lNK2hHL1U1OHRWZm1POUdwQzJwTGNTVWROcz06Y28vTW1EMSsrcXl3R1ZqWlJsWXoyYVZlSlJpbVNTZVk0aUx4VEhFVHh2bGhpaDBNZE1VZ002QmFNc1VzL3dNOFRWNWlNbkt4R3RMREpPQUpGTGUvRFE9PQ==",
            "dHJ4X3NlbmQ9MDlDMjk1MjQ4MTZCN0I1QzVGQTJDMEFGRDRENDM2REI2QURFNEJBOTg0NkQ3Nzo0Q0VFQTg3NTU5RTFGMUMyMDY5NDhBQkE2Njg5MjVDQkUwNEZDQjlDMzlFQTkxOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjQ6dGZwVHgyRUswZUw3bDlKN0FjUFJVZCtpdEpTM0xvdms5dFN0Y2RXRzBXND06bGdIUGZXd0FtWURUeWFhb1BBMm5scE5NdWdWaDJxUWgwdFBHeFFmcHM2RFNFOWxkcmMvM0tPK0xCRGpTTXpNQmd4V09vTWk3VmpNVSs0R09kY0RUQUE9PQ=="
        ]
    }
}`
	txSearchResult = `{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "txs": [
      {
        "hash": "1142188B5FFDD69AB892B47D748406DC4A4C41F7059DDB573639C14DA20701F8",
        "height": "25",
        "index": 0,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "QTk5NjNGQTg3NEI2QjFDOTRBMTQwMUYyOTYzMEIzNTI5OEU0N0Y3MEEyQkE2NQ=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjMzNjQ4MDYxODcwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2NzpBOTk2M0ZBODc0QjZCMUM5NEExNDAxRjI5NjMwQjM1Mjk4RTQ3RjcwQTJCQTY1OjUwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDoyOnd2SEczRWRkQmJYUUhjeUphbDBDUy9ZUWNOWXRFYkZZeGVqbnFmOU9oTTQ9OlFGOEltaVc1bVBxOWREbzAxWjE4dTZma1g2L2lLVGhITlBlSXJHVFRmM0ROdnEzbVV0OXViVTRpV2NwMFpzWWxURXE1RXZwcVRDQjFGU1B4ajNIWEFnPT0="
      },
      {
        "hash": "4CC298E26A2E8751CD6E6980C5FBB03EE41DAC6265529514362D94BE6C580F78",
        "height": "27",
        "index": 0,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "NDA3Mjc5QUJDRjdBQzhBQzM3NzY0RUJCMDNCQzc3M0NGRkU4NzVFQTU4RTVFQw=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjMzODUzMjE5MDQwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nzo0MDcyNzlBQkNGN0FDOEFDMzc3NjRFQkIwM0JDNzczQ0ZGRTg3NUVBNThFNUVDOjI1MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDozOnd2SEczRWRkQmJYUUhjeUphbDBDUy9ZUWNOWXRFYkZZeGVqbnFmOU9oTTQ9OllLdE1tQ3cyTjh1ODlvelhIUHl5ei82TUxRS1ZORmFOWTBSSXdYWmYrejhUaFNLRDhDd2RoenMxWXIybDRic3YvTTJWelFHRFBaTGtzUTJNc1ZFeEFnPT0="
      },
      {
        "hash": "464C308AA56F7D93DE9ADBD3380EAF72ABB6D9A99D7902E33190CFBFF90929E3",
        "height": "27",
        "index": 1,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QTk5NjNGQTg3NEI2QjFDOTRBMTQwMUYyOTYzMEIzNTI5OEU0N0Y3MEEyQkE2NQ=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "N0EyODNBRDNCOEY2OUQwM0M3OTkyQjI4RDBGQzUyNUI3M0E2NjQ3NjdBMDczMg=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjMzODUzMjE5MDQwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QTk5NjNGQTg3NEI2QjFDOTRBMTQwMUYyOTYzMEIzNTI5OEU0N0Y3MEEyQkE2NTo3QTI4M0FEM0I4RjY5RDAzQzc5OTJCMjhEMEZDNTI1QjczQTY2NDc2N0EwNzMyOjI0OTk5NzUwMDAwMDAwMDAwMDAwMDAwMDoyOm1XZ3hUekVSc1RnT0NvbFN2bHB4eFJkOHJQM2t6M2I5NFlzSG1sTVF5eEU9Ok5rWnZyR3VyQytsZDFtY1c1bkJldEkzMDllTnVHZENDMVk1Rk1TTisyS1FZV09PeFVlMXVROWVkbG9sbTIwRlRISGNOQmxOeXVITG1CQXZxQWUrZERBPT0="
      },
      {
        "hash": "0E18D23F1AF721CF2410FE7F9EA26A5E49D251A4CE581EE5648D14AFD62099E7",
        "height": "29",
        "index": 0,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QTk5NjNGQTg3NEI2QjFDOTRBMTQwMUYyOTYzMEIzNTI5OEU0N0Y3MEEyQkE2NQ=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "OTcyQTkzNTBEMjcxQ0I4MUIxQzY5RUM0NzA0NjU0QjY3OUY3QzAwQTc2QjI5Mg=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM0MDU3NTc4MTcwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QTk5NjNGQTg3NEI2QjFDOTRBMTQwMUYyOTYzMEIzNTI5OEU0N0Y3MEEyQkE2NTo5NzJBOTM1MEQyNzFDQjgxQjFDNjlFQzQ3MDQ2NTRCNjc5RjdDMDBBNzZCMjkyOjEyNDk5ODc1MDAwMDAwMDAwMDAwMDAwMDozOm1XZ3hUekVSc1RnT0NvbFN2bHB4eFJkOHJQM2t6M2I5NFlzSG1sTVF5eEU9Ojhmdmt4YkVNakZYOUtMRi9RbHBYRnFiRlNsVlhOZGtMR3BoS09ydy8xZEkzNXZJSEhNcDdhNEYvTm9jaitMdTZNa3FzMERqNlgrTGJWSzJxcGdLcENBPT0="
      },
      {
        "hash": "EFFDA65DDEF7DDF44DCE34D5883D255FFA24EFD8FA894B7D7C6ACAABEC3228F0",
        "height": "29",
        "index": 1,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "QjlGQTQ1NEYxNUEzQkZBOTdBNTA0MzMzMUNBNDE1QjRDQ0MxNTU1QUEwM0FDNQ=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM0MDU3NTc4MTcwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2NzpCOUZBNDU0RjE1QTNCRkE5N0E1MDQzMzMxQ0E0MTVCNENDQzE1NTVBQTAzQUM1OjEyNTAwMDAwMDAwMDAwMDAwMDAwMDAwMDo0Ond2SEczRWRkQmJYUUhjeUphbDBDUy9ZUWNOWXRFYkZZeGVqbnFmOU9oTTQ9Ojhqdmk5M2k1V1Z1d3cvZk5neGszZm5mWVZna3NkS1pDY2IzNTJlZWpPWW10WUdYRmFmUHcrZEZmdGpYWTBoeUlzblZqSEZ5RTRTdEdGRDJkdUg3NUNRPT0="
      },
      {
        "hash": "ABA508D0E4B80AD603ABB6BE30F454478006FB0F1AFCAB39306B7966F524EDA2",
        "height": "29",
        "index": 2,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "NDA3Mjc5QUJDRjdBQzhBQzM3NzY0RUJCMDNCQzc3M0NGRkU4NzVFQTU4RTVFQw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "QzQ4MjY4Q0U5RjA1OTZEMjRENTQ1NDc5M0RBMTg4NDk2REMzNDZGRkQ1RDE4RA=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM0MDU3NTc4MTcwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9NDA3Mjc5QUJDRjdBQzhBQzM3NzY0RUJCMDNCQzc3M0NGRkU4NzVFQTU4RTVFQzpDNDgyNjhDRTlGMDU5NkQyNEQ1NDU0NzkzREExODg0OTZEQzM0NkZGRDVEMThEOjEyNDk5NzUwMDAwMDAwMDAwMDAwMDAwMDoyOmprV3poeUF1QmEreGFGVjl6STNuSkdneHYrS1V4YzcrVndEc0lxWUdjRTA9Oi9KeG9rVndvUHdBTDdXNjFvSHVtMTFrUTNIQ3hFTVlBZldsc1QyQTRza2Z5NUgxUExyS3grS1dXVlRFUC9TZjRjRHBBdWphZ2NOUTZUWDJid1BUckNBPT0="
      },
      {
        "hash": "16C29EE0618B77AECCAA9AE0A6E646CAFBAA92C6F4603B5BFDC1446349F060FE",
        "height": "29",
        "index": 3,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "N0EyODNBRDNCOEY2OUQwM0M3OTkyQjI4RDBGQzUyNUI3M0E2NjQ3NjdBMDczMg=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "NDI4MEE0RUMzRkRCOTNEQUI5OURCNkI0NEMzMkZDRjNBNTZDRkYyNUI4NUJDMw=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM0MDU3NjY2ODMwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9N0EyODNBRDNCOEY2OUQwM0M3OTkyQjI4RDBGQzUyNUI3M0E2NjQ3NjdBMDczMjo0MjgwQTRFQzNGREI5M0RBQjk5REI2QjQ0QzMyRkNGM0E1NkNGRjI1Qjg1QkMzOjEyNDk5NjI1MDAwMDAwMDAwMDAwMDAwMDoyOkovK3Q2MVZxd2FyRTVoV1krTTdPNitCaDVIeU1uTHVnV3BBVEpncC9ETm89OkpndTZ3bmdkVTAraGhWWmFiaXdoRUdLbmNCbi9lNStqNEtJRGxVWnhOM1FtYVhMYzkvQVBDZnR4NVh6WmRFWml2cVo3enB5eFhUVERqVkh1MEtZR0RBPT0="
      },
      {
        "hash": "793374B46252ABFAC7D4F6F2DF2D3454F0E21838B62E630EE5809A84B86E54BD",
        "height": "31",
        "index": 0,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QTk5NjNGQTg3NEI2QjFDOTRBMTQwMUYyOTYzMEIzNTI5OEU0N0Y3MEEyQkE2NQ=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "MDRBNUVBRkJFNUIyODYxODMxRDBBN0NDNjk3M0FGN0EwMDFFMkNERTA1NTI2Qg=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM0MjYyMzE5NzYwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QTk5NjNGQTg3NEI2QjFDOTRBMTQwMUYyOTYzMEIzNTI5OEU0N0Y3MEEyQkE2NTowNEE1RUFGQkU1QjI4NjE4MzFEMEE3Q0M2OTczQUY3QTAwMUUyQ0RFMDU1MjZCOjYyNDk5Mzc1MDAwMDAwMDAwMDAwMDAwOjQ6bVdneFR6RVJzVGdPQ29sU3ZscHh4UmQ4clAza3ozYjk0WXNIbWxNUXl4RT06WHl1UHBaNlMwa1ordnNZUERKY2F0OEUvMUQ3Vk1zOTY2U2NMSkc4bGc2cTlTQnJwcDBvZnE5c3RHZkhHNGJYQkZ5TU5nUm5GYm52ZjV2MThwSEVKQkE9PQ=="
      },
      {
        "hash": "D89D913B82E3A46748206450169065D10DEBA2AE26121152450D424B2EDE86C8",
        "height": "31",
        "index": 1,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "N0EyODNBRDNCOEY2OUQwM0M3OTkyQjI4RDBGQzUyNUI3M0E2NjQ3NjdBMDczMg=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "Q0Y2MUE0RUQxMDNCMkQxRThBN0VCQjVCOTk2QkRCQUMyNzMzRDhBNkU0RTY0Mw=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM0MjYyMzE5NzYwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9N0EyODNBRDNCOEY2OUQwM0M3OTkyQjI4RDBGQzUyNUI3M0E2NjQ3NjdBMDczMjpDRjYxQTRFRDEwM0IyRDFFOEE3RUJCNUI5OTZCREJBQzI3MzNEOEE2RTRFNjQzOjYyNDk4MTI1MDAwMDAwMDAwMDAwMDAwOjM6Si8rdDYxVnF3YXJFNWhXWStNN082K0JoNUh5TW5MdWdXcEFUSmdwL0RObz06dlgyUW81azAvcm5jRkhGZXhpV0VwbVhMTHRiUlVUSWhiYTVXS3VOOVNaWGRZU1hQVDZ2MW5qa0lQQTFoY3BoMCs1QmUxd1hzOGFmSUY3amlJQ1RmRHc9PQ=="
      },
      {
        "hash": "4381EA6E6B3EAA6F07B638391D3347EF03200A8C5FC0DE732142E829A6BDAEEB",
        "height": "31",
        "index": 2,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "QTIzNzA1OUE2NkNFMzU2RUMxMzI5NUNGNURFN0YyNTYzQjVDNzc0NkFGOTczMA=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM0MjYyMzE5NzYwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2NzpBMjM3MDU5QTY2Q0UzNTZFQzEzMjk1Q0Y1REU3RjI1NjNCNUM3NzQ2QUY5NzMwOjYyNTAwMDAwMDAwMDAwMDAwMDAwMDAwOjU6d3ZIRzNFZGRCYlhRSGN5SmFsMENTL1lRY05ZdEViRll4ZWpucWY5T2hNND06MHUrczk5cG5HYnlZSngycGJrNzZ3a1hOR2x5M0NTcEU5dGQ4elFSbmhtd0VuU2JUdTF5MStNKzlXZUN0NVZSdi9WMTFhYXA0bEFYQi82U3Y4YVRrQ3c9PQ=="
      }
    ],
    "total_count": "63"
  }
}`
	txSearchResultMoreCondition = `{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "txs": [
      {
        "hash": "1142188B5FFDD69AB892B47D748406DC4A4C41F7059DDB573639C14DA20701F8",
        "height": "25",
        "index": 0,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "QTk5NjNGQTg3NEI2QjFDOTRBMTQwMUYyOTYzMEIzNTI5OEU0N0Y3MEEyQkE2NQ=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjMzNjQ4MDYxODcwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2NzpBOTk2M0ZBODc0QjZCMUM5NEExNDAxRjI5NjMwQjM1Mjk4RTQ3RjcwQTJCQTY1OjUwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDoyOnd2SEczRWRkQmJYUUhjeUphbDBDUy9ZUWNOWXRFYkZZeGVqbnFmOU9oTTQ9OlFGOEltaVc1bVBxOWREbzAxWjE4dTZma1g2L2lLVGhITlBlSXJHVFRmM0ROdnEzbVV0OXViVTRpV2NwMFpzWWxURXE1RXZwcVRDQjFGU1B4ajNIWEFnPT0="
      },
      {
        "hash": "4CC298E26A2E8751CD6E6980C5FBB03EE41DAC6265529514362D94BE6C580F78",
        "height": "27",
        "index": 0,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "NDA3Mjc5QUJDRjdBQzhBQzM3NzY0RUJCMDNCQzc3M0NGRkU4NzVFQTU4RTVFQw=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjMzODUzMjE5MDQwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nzo0MDcyNzlBQkNGN0FDOEFDMzc3NjRFQkIwM0JDNzczQ0ZGRTg3NUVBNThFNUVDOjI1MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDozOnd2SEczRWRkQmJYUUhjeUphbDBDUy9ZUWNOWXRFYkZZeGVqbnFmOU9oTTQ9OllLdE1tQ3cyTjh1ODlvelhIUHl5ei82TUxRS1ZORmFOWTBSSXdYWmYrejhUaFNLRDhDd2RoenMxWXIybDRic3YvTTJWelFHRFBaTGtzUTJNc1ZFeEFnPT0="
      },
      {
        "hash": "EFFDA65DDEF7DDF44DCE34D5883D255FFA24EFD8FA894B7D7C6ACAABEC3228F0",
        "height": "29",
        "index": 1,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "QjlGQTQ1NEYxNUEzQkZBOTdBNTA0MzMzMUNBNDE1QjRDQ0MxNTU1QUEwM0FDNQ=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM0MDU3NTc4MTcwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2NzpCOUZBNDU0RjE1QTNCRkE5N0E1MDQzMzMxQ0E0MTVCNENDQzE1NTVBQTAzQUM1OjEyNTAwMDAwMDAwMDAwMDAwMDAwMDAwMDo0Ond2SEczRWRkQmJYUUhjeUphbDBDUy9ZUWNOWXRFYkZZeGVqbnFmOU9oTTQ9Ojhqdmk5M2k1V1Z1d3cvZk5neGszZm5mWVZna3NkS1pDY2IzNTJlZWpPWW10WUdYRmFmUHcrZEZmdGpYWTBoeUlzblZqSEZ5RTRTdEdGRDJkdUg3NUNRPT0="
      },
      {
        "hash": "4381EA6E6B3EAA6F07B638391D3347EF03200A8C5FC0DE732142E829A6BDAEEB",
        "height": "31",
        "index": 2,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "QTIzNzA1OUE2NkNFMzU2RUMxMzI5NUNGNURFN0YyNTYzQjVDNzc0NkFGOTczMA=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM0MjYyMzE5NzYwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2NzpBMjM3MDU5QTY2Q0UzNTZFQzEzMjk1Q0Y1REU3RjI1NjNCNUM3NzQ2QUY5NzMwOjYyNTAwMDAwMDAwMDAwMDAwMDAwMDAwOjU6d3ZIRzNFZGRCYlhRSGN5SmFsMENTL1lRY05ZdEViRll4ZWpucWY5T2hNND06MHUrczk5cG5HYnlZSngycGJrNzZ3a1hOR2x5M0NTcEU5dGQ4elFSbmhtd0VuU2JUdTF5MStNKzlXZUN0NVZSdi9WMTFhYXA0bEFYQi82U3Y4YVRrQ3c9PQ=="
      },
      {
        "hash": "7D807E82D45960E6952F2812930896B1C676FCB40509CB6AB0C261D719E42A88",
        "height": "35",
        "index": 5,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "REQ3MDM1Mjk3N0YwQUIyRkY4OUJERTk0OTdBN0JBQTZEQThEOTI2NDY4OEE1Qw=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM0Njc0ODE3MDkwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2NzpERDcwMzUyOTc3RjBBQjJGRjg5QkRFOTQ5N0E3QkFBNkRBOEQ5MjY0Njg4QTVDOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjY6d3ZIRzNFZGRCYlhRSGN5SmFsMENTL1lRY05ZdEViRll4ZWpucWY5T2hNND06UHF6bmF0aUE2dVQrNHAvK01keFVjRHZ1WE4wdStrdEdqRnF4RE5Cd0Qvdm41TWFjTTJiS3RGR3F4SC8veUlrbzd5MXJyL0IyazVISFBiV3ZWUkZhQlE9PQ=="
      },
      {
        "hash": "BDB8CE1966F7755D66BC48FA06E3CCB0FABE1F3E2B6EE7D79BE993DD12D4B67F",
        "height": "37",
        "index": 5,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "REQ3MDM1Mjk3N0YwQUIyRkY4OUJERTk0OTdBN0JBQTZEQThEOTI2NDY4OEE1Qw=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM0ODgxODY2NDQwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2NzpERDcwMzUyOTc3RjBBQjJGRjg5QkRFOTQ5N0E3QkFBNkRBOEQ5MjY0Njg4QTVDOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjc6d3ZIRzNFZGRCYlhRSGN5SmFsMENTL1lRY05ZdEViRll4ZWpucWY5T2hNND06azB4WlRGUjVPVEtzTlR4TW5HbUk5MlB5R3ZobkJZcGpwQ3J4ME1pMHFWc2REVXc5RTJhNVBJT1ZDQzZnV2M5WGZ3MnNvQUpJSnJOMFp0TjVjb2VaRGc9PQ=="
      },
      {
        "hash": "C8C4340B450213D0969B87DF9DE580BE235C4BAFD69DC34172662653B38614DA",
        "height": "39",
        "index": 10,
        "tx_result": {
          "tags": [
            {
              "key": "YXBwLmZyb21hZGRyZXNz",
              "value": "QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2Nw=="
            },
            {
              "key": "YXBwLnRvYWRkcmVzcw==",
              "value": "REQ3MDM1Mjk3N0YwQUIyRkY4OUJERTk0OTdBN0JBQTZEQThEOTI2NDY4OEE1Qw=="
            },
            {
              "key": "YXBwLnRpbWVzdGFtcA==",
              "value": "MTU2NTE3NjM1MDg3OTEzMDgwMA=="
            },
            {
              "key": "YXBwLnR5cGU=",
              "value": "U2VuZA=="
            }
          ]
        },
        "tx": "dHJ4X3NlbmQ9QjUwOEVEMEQ1NDU5N0Q1MTZBNjgwRTc5NTFGMThDQUQyNEM3RUM5RkNGQ0Q2NzpERDcwMzUyOTc3RjBBQjJGRjg5QkRFOTQ5N0E3QkFBNkRBOEQ5MjY0Njg4QTVDOjEwMDAwMDAwMDAwMDAwMDAwMDAwOjg6d3ZIRzNFZGRCYlhRSGN5SmFsMENTL1lRY05ZdEViRll4ZWpucWY5T2hNND06MVVENkIyV3pQMVdmTWp5MVQ3MHlhVlNPNk91bWVLNER0dkdOS1FGRXR3d0ZJbTZiSmcxcnpydUs4aGFsUklqazFhOUhVWkRvenJTdjVtMWFDdzZzQ0E9PQ=="
      }
    ],
    "total_count": "7"
  }
}`
)
