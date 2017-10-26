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

1. Bittrex
2. Binance
3. Bitfinex
4. Liqui
