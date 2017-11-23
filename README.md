# Data fetcher for KyberNetwork reserve

## APIs

### Get prices for specific base-quote pair

```
<host>:8000/prices/<base>/<quote>
```

Where *<base>* is symbol of the base token and *<quote>* is symbol of the quote token

eg:
```
curl -X GET "http://13.229.54.28:8000/prices/omg/eth"
```

### Get prices for all base-quote pairs
```
<host>:8000/prices
```

eg:
```
curl -X GET "http://13.229.54.28:8000/prices"
```

### Get balances for all tokens of the reserve manager on blockchain
```
<host>:8000/balances
```

eg:
```
curl -X GET "http://13.229.54.28:8000/balances"
```

### Get balances for all tokens of the reserve manager on supported changes
```
<host>:8000/ebalances
```

eg:
```
curl -X GET "http://13.229.54.28:8000/ebalances"
```
### Deposit to exchanges
```
<host>:8000/deposit/:exchange_id
POST request
Form params:
  - amount: little endian hex string (must starts with 0x), eg: 0xde0b6b3a7640000
  - token: token id string, eg: ETH, EOS...
```

eg:
```
curl -X POST \
  http://localhost:8000/deposit/liqui \
  -H 'content-type: multipart/form-data' \
  -F token=EOS \
  -F amount=0xde0b6b3a7640000
```
Response:

```json
{
    "hash": "0x1b0c09f059904f1a9587641f2357c16c1c9fe43dfea161db31607f9221b0cfbb",
    "success": true
}
```
Where `hash` is the transaction hash

### Withdraw from exchanges
```
<host>:8000/withdraw/:exchange_id
POST request
Form params:
  - amount: little endian hex string (must starts with 0x), eg: 0xde0b6b3a7640000
  - token: token id string, eg: ETH, EOS...
```

eg:
```
curl -X POST \
  http://localhost:8000/withdraw/liqui \
  -H 'content-type: multipart/form-data' \
  -F token=EOS \
  -F amount=0xde0b6b3a7640000
```
Response:

```json
{
    "success": true
}
```
Where `hash` is the transaction hash

### Setting rates
```
<host>:8000/setrates
POST request
Form params:
  - sources: string, represent all base token IDs separated by "-", eg: "ETH-ETH"
  - dests: string, represent all quote token IDs separated by "-", eg: "KNC-EOS"
  - rates: string, represent all the rates in little endian hex string, rates are separated by "-", eg: "0x5-0x7"
  - expiries: string, represent all the expiry blocks in little endian hex string, they are separated by "-", eg: "0x989680-0x989680"
```
eg:
```
curl -X POST \
  http://localhost:8000/setrates \
  -H 'content-type: multipart/form-data' \
  -F sources=ETH-ETH \
  -F dests=KNC-EOS \
  -F rates=0x5-0x7 \
  -F expiries=0x989680-0x989680
```
Response:

```json
{
    "hash": "0x8004f8613b9944fc73c59b7a70b0a491c9d190e7d3703488423855ac8dada239",
    "success": true
}
```
Where `hash` is the transaction hash

### Trade
```
<host>:8000/trade/:exchange_id
POST request
Form params:
  - base: token id string, eg: ETH, EOS...
  - quote: token id string, eg: ETH, EOS...
  - amount: float
  - rate: float
  - type: "buy" or "sell"
```

eg:
```
curl -X POST \
  http://localhost:8000/trade/liqui \
  -F base=ETH \
  -F quote=KNC \
  -F rate=300 \
  -F type=buy \
  -F amount=0.01
```
Response:

```json
{
    "success": true,
    "done": 0,
    "remaining": 0.01,
    "finished": false,
}
```
Where `hash` is the transaction hash

### Get open orders
```
<host>:8000/orders
GET request
```

Response:
```json
{
	"data": {
		"binance": {
			"Valid": true,
			"Error": "",
			"Timestamp": "1511426133904",
			"ReturnTime": "1511426134053",
			"Data": [{
				"Base": "KNC",
				"Quote": "ETH",
				"OrderId": "2025775",
				"Price": 0.002,
				"OrigQty": 10,
				"ExecutedQty": 0,
				"TimeInForce": "GTC",
				"Type": "LIMIT",
				"Side": "BUY",
				"StopPrice": "0.00000000",
				"IcebergQty": "0.00000000",
				"Time": 1511426052681
			}]
		},
		"liqui": {
			"Valid": true,
			"Error": "",
			"Timestamp": "1511426133904",
			"ReturnTime": "1511426134159",
			"Data": []
		}
	},
	"success": true,
	"timestamp": "1511426136158",
	"version": 26
}
```

### Cancel order
```
<host>:8000/cancelorder/:exchange
POST request
Form params:
  - base: token id string, eg: ETH, EOS...
  - quote: token id string, eg: ETH, EOS...
  - order_id: string
```

response:
```json
{
    "reason": "UNKNOWN_ORDER",
    "success": false
}
```



## Supported tokens

1. eth (ETH)
2. bat (BAT)
3. civic (CVC)
4. digix (DGD)
5. eos (EOS)
6. adex (ADX)
7. funfair (FUN)
8. golem (GNT)
9. kybernetwork (KNC)
10. link (LINK)
11. monaco (MCO)
12. omisego (OMG)
13. tenx (PAY)

## Supported exchanges

1. Bittrex (bittrex)
2. Binance (binance)
3. Bitfinex (bitfinex)
4. Liqui (liqui)
