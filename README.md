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
