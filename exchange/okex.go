package exchange

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Okex struct {
	interf       OkexInterface
	pairs        []common.TokenPair
	addresses    map[string]ethereum.Address
	exchangeInfo *common.ExchangeInfo
	fees         common.ExchangeFees
}

func (self *Okex) MarshalText() (text []byte, err error) {
	return []byte(self.ID()), nil
}

func (self *Okex) Address(token common.Token) (ethereum.Address, bool) {
	addr, supported := self.addresses[token.ID]
	return addr, supported
}

func (self *Okex) UpdateAllDepositAddresses(address string) {
	for k, _ := range self.addresses {
		self.addresses[k] = ethereum.HexToAddress(address)
	}
}

func (self *Okex) UpdateDepositAddress(token common.Token, address string) {
	self.addresses[token.ID] = ethereum.HexToAddress(address)
}

func (self *Okex) UpdatePrecisionLimit(pair common.TokenPair, exchangeInfo OkexInfo) {
	pairName := fmt.Sprintf("%s_%s", strings.ToLower(pair.Base.ID), strings.ToLower(pair.Quote.ID))
	for _, symbol := range exchangeInfo.Data {
		if symbol.Symbol == pairName {
			exchangePrecisionLimit := common.ExchangePrecisionLimit{}
			exchangePrecisionLimit.Precision.Amount = symbol.AmounPrecision
			exchangePrecisionLimit.Precision.Price = symbol.PricePrecision
			exchangePrecisionLimit.AmountLimit.Min = symbol.MinAmount
			self.exchangeInfo.Update(pair.PairID(), exchangePrecisionLimit)
			break
		}
	}
}

func (self *Okex) UpdatePairsPrecision() {
	exchangeInfo, err := self.interf.GetExchangeInfo()
	if err != nil {
		log.Printf("Get exchange info failed: %s\n", err)
	} else {
		for _, pair := range self.pairs {
			self.UpdatePrecisionLimit(pair, exchangeInfo)
		}
	}
}

func (self *Okex) FetchOnePairData(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map,
	timepoint uint64) {

	defer wg.Done()
	result := common.ExchangePrice{}

	timestamp := common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Timestamp = timestamp
	result.Valid = true
	resp_data, err := self.interf.GetDepthOnePair(pair, timepoint)
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		for _, buy := range resp_data.Bids {
			quantity := buy[1]
			rate := buy[0]
			result.Bids = append(
				result.Bids,
				common.PriceEntry{
					quantity,
					rate,
				},
			)
		}
		for _, sell := range resp_data.Asks {
			quantity := sell[1]
			rate := sell[0]
			result.Asks = append(
				result.Asks,
				common.PriceEntry{
					quantity,
					rate,
				},
			)
		}
	}
	data.Store(pair.PairID(), result)
}

func (self *Okex) FetchPriceData(timepoint uint64) (map[common.TokenPairID]common.ExchangePrice, error) {
	wait := sync.WaitGroup{}
	data := sync.Map{}
	pairs := self.pairs
	for _, pair := range pairs {
		wait.Add(1)
		go self.FetchOnePairData(&wait, pair, &data, timepoint)
	}
	wait.Wait()
	result := map[common.TokenPairID]common.ExchangePrice{}
	data.Range(func(key, value interface{}) bool {
		result[key.(common.TokenPairID)] = value.(common.ExchangePrice)
		return true
	})
	return result, nil
}

func (self *Okex) GetInfo() (common.ExchangeInfo, error) {
	return *self.exchangeInfo, nil
}

func (self *Okex) GetExchangeInfo(pair common.TokenPairID) (common.ExchangePrecisionLimit, error) {
	data, err := self.exchangeInfo.Get(pair)
	return data, err
}

func (self *Okex) GetFee() common.ExchangeFees {
	return self.fees
}

func (self *Okex) ID() common.ExchangeID {
	return common.ExchangeID("okex")
}

func (self *Okex) TokenPairs() []common.TokenPair {
	return self.pairs
}

func (self *Okex) Name() string {
	return "okex"
}

func (self *Okex) QueryOrder(symbol string, id uint64, timepoint uint64) (done float64, remaining float64, finished bool, err error) {
	// result, err := self.interf.OrderStatus(symbol, id, timepoint)
	// if err != nil {
	// 	return 0, 0, false, err
	// } else {
	// 	// done, _ := strconv.ParseFloat()
	// 	// TODO: complete query order
	// 	return 0, 0, false, err
	// }
	return 0, 0, false, err
}

func (self *Okex) Trade(
	tradeType string,
	base common.Token,
	quote common.Token,
	rate float64,
	amount float64,
	timepoint uint64) (id string, done float64, remaining float64, finished bool, err error) {

	result, err := self.interf.Trade(base, quote, rate, amount, timepoint)
	// symbol := base.ID + quote.ID
	// TODO: complete Trade
	if err != nil {
		return "", 0, 0, false, err
	} else {
		done, remaining, finished, err := self.QueryOrder(
			base.ID+quote.ID,
			result.OrderID,
			timepoint+20,
		)
		return id, done, remaining, finished, err
	}
}

func (self *Okex) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) (string, error) {
	tx, err := self.interf.Withdraw(token, amount, address, timepoint)
	return tx, err
}

func (self *Okex) CancelOrder(id common.ActivityID) error {
	// TODO: complete cancel order later
	return nil
}

func (self *Okex) FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error) {
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
		if resp_data.Result != true {
			result.Valid = false
			result.Error = fmt.Sprintf("Fetch Ebalance failed")
		} else {
			for tokenID, value := range resp_data.Info.Funds.Free {
				_, exist := common.SupportedTokens[strings.ToUpper(tokenID)]
				if exist {
					avai, _ := strconv.ParseFloat(value, 64)
					result.AvailableBalance[tokenID] = avai
				}
			}
			for tokenID, value := range resp_data.Info.Funds.Freezed {
				_, exist := common.SupportedTokens[strings.ToUpper(tokenID)]
				if exist {
					locked, _ := strconv.ParseFloat(value, 64)
					result.LockedBalance[tokenID] = locked
				}
			}
			for tokenID, value := range resp_data.Info.Funds.Borrow {
				_, exist := common.SupportedTokens[strings.ToUpper(tokenID)]
				if exist {
					borrow, _ := strconv.ParseFloat(value, 64)
					result.DepositBalance[tokenID] = borrow
				}
			}
		}
	}
	return result, nil
}

func (self *Okex) DepositStatus(id common.ActivityID, timepoint uint64) (string, error) {
	// TODO: Complete get deposit status
	return "", errors.New("Deposit doesn't exist. This shouldn't happen unless tx returned from binance and activity ID are not consistently designed")
}

func (self *Okex) WithdrawStatus(id common.ActivityID, timepoint uint64) (string, string, error) {
	// TODO: Complete get withdraw status
	return "", "", errors.New("Withdrawal doesn't exist. This shouldn't happen unless tx returned from withdrawal from binance and activity ID are not consistently designed")
}

func (self *Okex) OrderStatus(id common.ActivityID, timepoint uint64) (string, error) {
	// TODO: Complete get order status
	return "done", nil
}

func NewOkex(interf OkexInterface) *Okex {
	return &Okex{
		interf,
		[]common.TokenPair{
			common.MustCreateTokenPair("FUN", "ETH"),
			common.MustCreateTokenPair("MCO", "ETH"),
			common.MustCreateTokenPair("OMG", "ETH"),
			common.MustCreateTokenPair("EOS", "ETH"),
			common.MustCreateTokenPair("KNC", "ETH"),
			common.MustCreateTokenPair("LINK", "ETH"),
		},
		map[string]ethereum.Address{},
		common.NewExchangeInfo(),
		common.NewExchangeFee(
			common.TradingFee{
				"taker": 0.001,
				"maker": 0.001,
			},
			common.NewFundingFee(
				map[string]float32{},
				map[string]float32{},
			),
		),
	}
}
