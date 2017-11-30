package core

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

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

func timebasedID(id string) string {
	return strconv.Itoa(int(time.Now().UnixNano())) + "|" + id
}

func (self ReserveCore) CancelOrder(base, quote common.Token, id string, exchange common.Exchange) error {
	return exchange.CancelOrder(base, quote, id)
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
	var status string
	if err != nil {
		status = "failed"
	} else {
		if finished {
			status = "done"
		} else {
			status = "submitted"
		}
	}
	go self.activityStorage.Record(
		"trade",
		timebasedID(id),
		string(exchange.ID()),
		map[string]interface{}{
			"exchange":  exchange,
			"type":      tradeType,
			"base":      base,
			"quote":     quote,
			"rate":      rate,
			"amount":    amount,
			"timepoint": timepoint,
		}, map[string]interface{}{
			"id":        id,
			"done":      done,
			"remaining": remaining,
			"finished":  finished,
			"error":     err,
		},
		status,
		timepoint,
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
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "submitted"
	}
	go self.activityStorage.Record(
		"deposit",
		timebasedID(tx.Hex()),
		string(exchange.ID()),
		map[string]interface{}{
			"exchange":  exchange,
			"token":     token,
			"amount":    common.BigToFloat(amount, token.Decimal),
			"timepoint": timepoint,
		}, map[string]interface{}{
			"tx":    tx.Hex(),
			"error": err,
		},
		status,
		timepoint,
	)
	log.Printf(
		"Core ----------> Deposit to %s: token: %s, amount: %d, timestamp: %d ==> Result: tx: %s, error: %s",
		exchange.ID(), token.ID, amount.Uint64(), timepoint, tx.Hex(), err,
	)
	return tx, err
}

func (self ReserveCore) Withdraw(
	exchange common.Exchange, token common.Token,
	amount *big.Int, timepoint uint64) (string, error) {

	_, supported := exchange.Address(token)
	var err error
	var id string
	if !supported {
		err = errors.New(fmt.Sprintf("Exchange %s doesn't support token %s", exchange.ID(), token.ID))
	} else {
		id, err = exchange.Withdraw(token, amount, self.rm, timepoint)
	}
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "submitted"
	}
	go self.activityStorage.Record(
		"withdraw",
		timebasedID(id),
		string(exchange.ID()),
		map[string]interface{}{
			"exchange":  exchange,
			"token":     token,
			"amount":    common.BigToFloat(amount, token.Decimal),
			"timepoint": timepoint,
		}, map[string]interface{}{
			"error": err,
			"id":    id,
		},
		status,
		timepoint,
	)
	log.Printf(
		"Core ----------> Withdraw from %s: token: %s, amount: %d, timestamp: %d ==> Result: id: %s, error: %s",
		exchange.ID(), token.ID, amount.Uint64(), timepoint, id, err,
	)
	return id, err
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
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "submitted"
	}
	go self.activityStorage.Record(
		"set_rates",
		timebasedID(tx.Hex()),
		"blockchain",
		map[string]interface{}{
			"sources":      sources,
			"dests":        dests,
			"rates":        rates,
			"expiryBlocks": expiryBlocks,
		}, map[string]interface{}{
			"tx":    tx.Hex(),
			"error": err,
		},
		status,
		common.GetTimepoint(),
	)
	log.Printf(
		"Core ----------> Set rates: ==> Result: tx: %s, error: %s",
		tx.Hex(), err,
	)
	return tx, err
}
