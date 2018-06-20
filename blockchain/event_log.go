package blockchain

import (
	"errors"
	"math"
	"math/big"
	"strings"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ethRateResolver func(timepoint uint64) float64

func updateTradeLogs(allLogs []common.KNLog, logItem types.Log, ts uint64, ethFn ethRateResolver) ([]common.KNLog, error) {
	var (
		tradeLog      common.TradeLog
		updateLastLog = false
	)

	if len(logItem.Topics) == 0 {
		return allLogs, errors.New("log item has no topic")
	}

	// if the transaction hash is the same and last log is a TradeLog, update it,
	// otherwise append new one
	if len(allLogs) > 0 && allLogs[len(allLogs)-1].TxHash() == logItem.TxHash {
		var ok bool
		if tradeLog, ok = allLogs[len(allLogs)-1].(common.TradeLog); !ok {
			updateLastLog = true
		}
	}

	tradeLog = common.TradeLog{
		BlockNumber:     logItem.BlockNumber,
		TransactionHash: logItem.TxHash,
		Index:           logItem.Index,
		Timestamp:       ts,
	}

	switch logItem.Topics[0].Hex() {
	case feeToWalletEvent:
		reserveAddr, walletAddr, walletFee := logDataToFeeWalletParams(logItem.Data)
		tradeLog.ReserveAddress = reserveAddr
		tradeLog.WalletAddress = walletAddr
		tradeLog.WalletFee = walletFee.Big()
	case burnFeeEvent:
		reserveAddr, burnFees := logDataToBurnFeeParams(logItem.Data)
		tradeLog.ReserveAddress = reserveAddr
		tradeLog.BurnFee = burnFees.Big()
	case tradeEvent:
		srcAddr, destAddr, srcAmount, destAmount := logDataToTradeParams(logItem.Data)
		tradeLog.SrcAddress = srcAddr
		tradeLog.DestAddress = destAddr
		tradeLog.SrcAmount = srcAmount.Big()
		tradeLog.DestAmount = destAmount.Big()
		tradeLog.UserAddress = ethereum.BytesToAddress(logItem.Topics[1].Bytes())

		if ethRate := ethFn(tradeLog.Timestamp / 1000000); ethRate != 0 {
			eth := common.ETHToken()
			ethAmount := new(big.Float)
			// TODO: update the ETH amount calculation
			if strings.ToLower(eth.Address) == strings.ToLower(srcAddr.String()) {
				ethAmount.SetInt(tradeLog.SrcAmount)
			} else {
				ethAmount.SetInt(tradeLog.DestAmount)
			}

			// fiat amount = ETH amount * rate
			ethAmount = ethAmount.Mul(ethAmount, new(big.Float).SetFloat64(ethRate))
			ethAmount.Quo(ethAmount, new(big.Float).SetFloat64(math.Pow10(18)))
			tradeLog.FiatAmount, _ = ethAmount.Float64()
		}
	}

	if updateLastLog {
		allLogs[len(allLogs)-1] = tradeLog
	} else {
		allLogs = append(allLogs, tradeLog)
	}

	return allLogs, nil
}

func logDataToTradeParams(data []byte) (ethereum.Address, ethereum.Address, ethereum.Hash, ethereum.Hash) {
	srcAddr := ethereum.BytesToAddress(data[0:32])
	desAddr := ethereum.BytesToAddress(data[32:64])
	srcAmount := ethereum.BytesToHash(data[64:96])
	desAmount := ethereum.BytesToHash(data[96:128])
	return srcAddr, desAddr, srcAmount, desAmount
}

func logDataToFeeWalletParams(data []byte) (ethereum.Address, ethereum.Address, ethereum.Hash) {
	reserveAddr := ethereum.BytesToAddress(data[0:32])
	walletAddr := ethereum.BytesToAddress(data[32:64])
	walletFee := ethereum.BytesToHash(data[64:96])
	return reserveAddr, walletAddr, walletFee
}

func logDataToBurnFeeParams(data []byte) (ethereum.Address, ethereum.Hash) {
	reserveAddr := ethereum.BytesToAddress(data[0:32])
	burnFees := ethereum.BytesToHash(data[32:64])
	return reserveAddr, burnFees
}

func logDataToCatLog(data []byte) (ethereum.Address, string) {
	address := ethereum.BytesToAddress(data[0:32])
	cat := ethereum.BytesToHash(data[32:64])
	return address, cat.Hex()
}
