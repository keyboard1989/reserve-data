# 0.4.1 (2018-02-19)
## Features:
- Listed more 4 tokens (eng, salt, appc, rdn)
- Added more tools for monitoring and testing such as deposit/withdraw trigger, rate validator

## Bug fixes:
- Fixed submit empty setrate for the first one
- Fixed bug in rare case that panics when core couldn't get mined nonce
- Fixed incompatibility between geth and parity in tx receipt data
- Enable microsecond info in log

## improvements:
- Separated cex token pairs to config
- Separated cex fee to config
- Added sanity checks on setrates, deposit, withdraw and trade
- Added env tag to sentry

## Compatability:
- This version only works with KyberNetwork smart contracts version 0.3.0 or later

# 0.4.0 (2018-02-08)

## Features:
- Support rebalance toggle, dynamic target qty with set/confirm key model
- Support multiple keys for different roles

## Bug fixes:
- Fixed minor bugs
- Detect throwing txs

## Improvements:
- Done sanity check in with setrate api
- Rebroadcasting tx to multiple node to improve tx propagation
- Replace staled/long mining set rate txs
- Made improvements to the code base
- Applied timeout to communication to nodes to ensure analytic doesn't have to wait for too long to set another rate

## Compatability:
- This version only works with KyberNetwork smart contracts version 0.3.0 or later

# 0.3.0 (2018-01-31)

## Features:
- Introduce various key permissions
- New API for getting KN rate historical data
- New API for getting trade history on cexs

## Bugfixes:
- Handle lost transactions

## Improvements:
- Using multiple nodes to broadcast tx
- Avoid storing redundant rate data
- More code refactoring

## Compatability:
- This version only works with KyberNetwork smart contracts version 0.3.0


