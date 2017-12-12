package exchange

import (
	"errors"
	"strings"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"

	"fmt"
	"math/big"
)

type Bittrex struct {
	interf    BittrexInterface
	pairs     []common.TokenPair
	addresses map[string]ethereum.Address
}

func (self *Bittrex) MarshalText() (text []byte, err error) {
	return []byte(self.ID()), nil
}

func (self *Bittrex) Address(token common.Token) (ethereum.Address, bool) {
	addr, supported := self.addresses[token.ID]
	return addr, supported
}

func (self *Bittrex) UpdateAllDepositAddresses(address string) {
	for k, _ := range self.addresses {
		self.addresses[k] = ethereum.HexToAddress(address)
	}
}

func (self *Bittrex) UpdateDepositAddress(token common.Token, address string) {
	self.addresses[token.ID] = ethereum.HexToAddress(address)
}

func (self *Bittrex) ID() common.ExchangeID {
	return common.ExchangeID("bittrex")
}

func (self *Bittrex) TokenPairs() []common.TokenPair {
	return self.pairs
}

func (self *Bittrex) Name() string {
	return "bittrex"
}

func (self *Bittrex) Trade(tradeType string, base common.Token, quote common.Token, rate float64, amount float64, timepoint uint64) (string, float64, float64, bool, error) {
	return self.interf.Trade(tradeType, base, quote, rate, amount, timepoint)
}

func (self *Bittrex) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) (string, error) {
	resp, err := self.interf.Withdraw(token, amount, address, timepoint)
	if err != nil {
		return "", err
	} else {
		if resp.Success {
			return resp.Result["uuid"] + "|" + token.ID, nil
		} else {
			return "", errors.New(resp.Error)
		}
	}
}

func (self *Bittrex) DepositStatus(id common.ActivityID, timepoint uint64) (string, error) {
	idParts := strings.Split(id.EID, "|")
	if len(idParts) != 2 {
		// here, the exchange id part in id is malformed
		// 1. because analytic didn't pass original ID
		// 2. id is not constructed correctly in a form of uuid + "|" + token
		return "", errors.New("Invalid deposit id")
	}
	txID := idParts[0]
	currency := idParts[1]
	histories, err := self.interf.DepositHistory(currency, timepoint)
	if err != nil {
		return "", err
	} else {
		for _, deposit := range histories.Result {
			if deposit.TxId == txID {
				if deposit.PendingPayment {
					return "", nil
				} else {
					return "done", nil
				}
			}
		}
		return "", errors.New("Deposit with tx " + txID + " of currency " + currency + " is not found on bittrex")
	}
}

func (self *Bittrex) CancelOrder(id common.ActivityID) error {
	uuid := id.EID
	resp, err := self.interf.CancelOrder(uuid, common.GetTimepoint())
	if err != nil {
		return err
	} else {
		if resp.Success {
			return nil
		} else {
			return errors.New(resp.Error)
		}
	}
}

func (self *Bittrex) WithdrawStatus(id common.ActivityID, timepoint uint64) (string, error) {
	idParts := strings.Split(id.EID, "|")
	if len(idParts) != 2 {
		// here, the exchange id part in id is malformed
		// 1. because analytic didn't pass original ID
		// 2. id is not constructed correctly in a form of uuid + "|" + token
		return "", errors.New("Invalid deposit id")
	}
	txID := idParts[0]
	currency := idParts[1]
	histories, err := self.interf.WithdrawHistory(currency, timepoint)
	if err != nil {
		return "", err
	} else {
		for _, withdraw := range histories.Result {
			if withdraw.TxId == txID {
				if withdraw.PendingPayment {
					return "", nil
				} else {
					return "done", nil
				}
			}
		}
		return "", errors.New("Withdraw with tx " + txID + " of currency " + currency + " is not found on bittrex")
	}
}

func (self *Bittrex) OrderStatus(id common.ActivityID, timepoint uint64) (string, error) {
	uuid := id.EID
	resp_data, err := self.interf.OrderStatus(uuid, timepoint)
	if err != nil {
		return "", err
	} else {
		if resp_data.Result.IsOpen {
			return "", nil
		} else {
			return "done", nil
		}
	}
}

func (self *Bittrex) FetchPriceData(timepoint uint64) (map[common.TokenPairID]common.ExchangePrice, error) {
	wait := sync.WaitGroup{}
	data := sync.Map{}
	pairs := self.pairs
	for _, pair := range pairs {
		wait.Add(1)
		go self.interf.FetchOnePairData(&wait, pair, &data, timepoint)
	}
	wait.Wait()
	result := map[common.TokenPairID]common.ExchangePrice{}
	data.Range(func(key, value interface{}) bool {
		result[key.(common.TokenPairID)] = value.(common.ExchangePrice)
		return true
	})
	return result, nil
}

func (self *Bittrex) FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	resp_data, err := self.interf.GetInfo(timepoint)
	result.ReturnTime = common.GetTimestamp()
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		result.AvailableBalance = map[string]float64{}
		result.LockedBalance = map[string]float64{}
		result.DepositBalance = map[string]float64{}
		if resp_data.Success {
			for _, b := range resp_data.Result {
				tokenID := b.Currency
				_, exist := common.SupportedTokens[tokenID]
				if exist {
					result.AvailableBalance[tokenID] = b.Available
					result.LockedBalance[tokenID] = b.Pending
					result.DepositBalance[tokenID] = 0
				}
			}
		} else {
			result.Valid = false
			result.Error = resp_data.Error
		}
	}
	return result, nil
}

func NewBittrex(interf BittrexInterface) *Bittrex {
	return &Bittrex{
		interf,
		[]common.TokenPair{
			common.MustCreateTokenPair("OMG", "ETH"),
			common.MustCreateTokenPair("DGD", "ETH"),
			common.MustCreateTokenPair("CVC", "ETH"),
			common.MustCreateTokenPair("FUN", "ETH"),
			common.MustCreateTokenPair("MCO", "ETH"),
			common.MustCreateTokenPair("GNT", "ETH"),
			common.MustCreateTokenPair("ADX", "ETH"),
			common.MustCreateTokenPair("PAY", "ETH"),
			common.MustCreateTokenPair("BAT", "ETH"),
		},
		map[string]ethereum.Address{},
	}
}
