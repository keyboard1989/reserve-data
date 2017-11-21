package core

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type ReserveCore struct {
	blockchain      Blockchain
	activityStorage ActivityStorage
	rm              ethereum.Address
}

func NewReserveCore(
	blockchain Blockchain,
	storage ActivityStorage,
	rm ethereum.Address) *ReserveCore {
	return &ReserveCore{
		blockchain,
		storage,
		rm,
	}
}

func (self ReserveCore) Trade(
	exchange common.Exchange,
	tradeType string,
	base common.Token,
	quote common.Token,
	rate float64,
	amount float64,
	timepoint uint64) (id string, done float64, remaining float64, finished bool, err error) {

	id, done, remaining, finished, err = exchange.Trade(tradeType, base, quote, rate, amount, timepoint)
	go self.activityStorage.Record(
		"trade", map[string]interface{}{
			"exchange":  exchange,
			"type":      tradeType,
			"base":      base,
			"quote":     quote,
			"rate":      rate,
			"amount":    amount,
			"timepoint": timepoint,
		}, map[string]interface{}{
			"done":      done,
			"remaining": remaining,
			"finished":  finished,
			"error":     err,
		},
	)
	log.Printf(
		"Core ----------> %s on %s: base: %s, quote: %s, rate: %s, amount: %s, timestamp: %d ==> Result: done: %s, remaining: %s, finished: %t, error: %s",
		tradeType, exchange.ID(), base.ID, quote.ID,
		strconv.FormatFloat(rate, 'f', -1, 64),
		strconv.FormatFloat(amount, 'f', -1, 64), timepoint,
		strconv.FormatFloat(done, 'f', -1, 64),
		strconv.FormatFloat(remaining, 'f', -1, 64),
		finished, err,
	)
	return id, done, remaining, finished, err
}

func (self ReserveCore) Deposit(
	exchange common.Exchange,
	token common.Token,
	amount *big.Int,
	timepoint uint64) (ethereum.Hash, error) {

	address, supported := exchange.Address(token)
	tx := ethereum.Hash{}
	var err error
	if !supported {
		tx = ethereum.Hash{}
		err = errors.New(fmt.Sprintf("Exchange %s doesn't support token %s", exchange.ID(), token.ID))
	} else {
		tx, err = self.blockchain.Send(token, amount, address)
	}
	go self.activityStorage.Record(
		"deposit", map[string]interface{}{
			"exchange":  exchange,
			"token":     token,
			"amount":    common.BigToFloat(amount, token.Decimal),
			"timepoint": timepoint,
		}, map[string]interface{}{
			"tx":    tx,
			"error": err,
		},
	)
	log.Printf(
		"Core ----------> Deposit to %s: token: %s, amount: %d, timestamp: %d ==> Result: tx: %s, error: %s",
		exchange.ID(), token.ID, amount.Uint64(), timepoint, tx.Hex(), err,
	)
	return tx, err
}

func (self ReserveCore) Withdraw(
	exchange common.Exchange, token common.Token,
	amount *big.Int, timepoint uint64) (ethereum.Hash, error) {

	_, supported := exchange.Address(token)
	var err error
	var txHash ethereum.Hash
	if !supported {
		err = errors.New(fmt.Sprintf("Exchange %s doesn't support token %s", exchange.ID(), token.ID))
	} else {
		txHash, err = exchange.Withdraw(token, amount, self.rm, timepoint)
	}
	go self.activityStorage.Record(
		"withdraw", map[string]interface{}{
			"exchange":  exchange,
			"token":     token,
			"amount":    common.BigToFloat(amount, token.Decimal),
			"timepoint": timepoint,
		}, map[string]interface{}{
			"error":  err,
			"txhash": txHash,
		},
	)
	log.Printf(
		"Core ----------> Withdraw from %s: token: %s, amount: %d, timestamp: %d ==> Result: txhash: %s, error: %s",
		exchange.ID(), token.ID, amount.Uint64(), timepoint, txHash.Hex(), err,
	)
	return txHash, err
}

func (self ReserveCore) SetRates(
	sources []common.Token,
	dests []common.Token,
	rates []*big.Int,
	expiryBlocks []*big.Int) (ethereum.Hash, error) {

	lensources := len(sources)
	lendests := len(dests)
	lenrates := len(rates)
	lenblocks := len(expiryBlocks)
	tx := ethereum.Hash{}
	var err error
	if lensources != lendests || lensources != lenrates || lensources != lenblocks {
		err = errors.New("Sources, dests, rates and expiryBlocks must have the same length")
	} else {
		sourceAddrs := []ethereum.Address{}
		for _, source := range sources {
			sourceAddrs = append(sourceAddrs, ethereum.HexToAddress(source.Address))
		}
		destAddrs := []ethereum.Address{}
		for _, dest := range dests {
			destAddrs = append(destAddrs, ethereum.HexToAddress(dest.Address))
		}
		tx, err = self.blockchain.SetRates(sourceAddrs, destAddrs, rates, expiryBlocks)
	}
	go self.activityStorage.Record(
		"set_rates", map[string]interface{}{
			"sources":      sources,
			"dests":        dests,
			"rates":        rates,
			"expiryBlocks": expiryBlocks,
		}, map[string]interface{}{
			"tx":    tx,
			"error": err,
		},
	)
	log.Printf(
		"Core ----------> Set rates: ==> Result: tx: %s, error: %s",
		tx.Hex(), err,
	)
	return tx, err
}

func (self ReserveCore) GetRecords() ([]common.ActivityRecord, error) {
	return self.activityStorage.GetAllRecords()
}
