package exchange

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

const BITTREX_EPSILON float64 = 0.000001

type Bittrex struct {
	interf       BittrexInterface
	pairs        []common.TokenPair
	tokens       []common.Token
	addresses    *common.ExchangeAddresses
	storage      BittrexStorage
	exchangeInfo *common.ExchangeInfo
	fees         common.ExchangeFees
	minDeposit   common.ExchangesMinDeposit
}

func (self *Bittrex) TokenAddresses() map[string]ethereum.Address {
	return self.addresses.GetData()
}

func (self *Bittrex) MarshalText() (text []byte, err error) {
	return []byte(self.ID()), nil
}

func (self *Bittrex) Address(token common.Token) (ethereum.Address, bool) {
	addr, supported := self.addresses.Get(token.ID)
	return addr, supported
}

func (self *Bittrex) GetFee() common.ExchangeFees {
	return self.fees
}

func (self *Bittrex) GetMinDeposit() common.ExchangesMinDeposit {
	return self.minDeposit
}

func (self *Bittrex) UpdateAllDepositAddresses(address string) {
	data := self.addresses.GetData()
	for k := range data {
		self.addresses.Update(k, ethereum.HexToAddress(address))
	}
}

func (self *Bittrex) UpdateDepositAddress(token common.Token, address string) {
	liveAddress, _ := self.interf.GetDepositAddress(token.ID)
	if liveAddress.Result.Address != "" {
		self.addresses.Update(token.ID, ethereum.HexToAddress(liveAddress.Result.Address))
	} else {
		self.addresses.Update(token.ID, ethereum.HexToAddress(address))
	}
}

func (self *Bittrex) UpdatePrecisionLimit(pair common.TokenPair, symbols []BittPairInfo) {
	pairName := strings.ToUpper(pair.Base.ID) + strings.ToUpper(pair.Quote.ID)
	for _, symbol := range symbols {
		symbolName := strings.ToUpper(symbol.Base) + strings.ToUpper(symbol.Quote)
		if symbolName == pairName {
			exchangePrecisionLimit := common.ExchangePrecisionLimit{}
			//update precision
			exchangePrecisionLimit.Precision.Amount = 8
			exchangePrecisionLimit.Precision.Price = 8
			// update limit
			exchangePrecisionLimit.AmountLimit.Min = symbol.MinAmount
			exchangePrecisionLimit.MinNotional = 0.02
			self.exchangeInfo.Update(pair.PairID(), exchangePrecisionLimit)
			break
		}
	}
}

func (self *Bittrex) GetExchangeInfo(pair common.TokenPairID) (common.ExchangePrecisionLimit, error) {
	pairInfo, err := self.exchangeInfo.Get(pair)
	return pairInfo, err
}

func (self *Bittrex) GetInfo() (*common.ExchangeInfo, error) {
	return self.exchangeInfo, nil
}

func (self *Bittrex) UpdatePairsPrecision() {
	exchangeInfo, err := self.interf.GetExchangeInfo()
	if err == nil {
		symbols := exchangeInfo.Pairs
		for _, pair := range self.pairs {
			self.UpdatePrecisionLimit(pair, symbols)
		}
	} else {
		log.Printf("Get exchange info failed: %s\n", err)
	}
}

func (self *Bittrex) ID() common.ExchangeID {
	return common.ExchangeID("bittrex")
}

func (self *Bittrex) Name() string {
	return "bittrex"
}

func (self *Bittrex) QueryOrder(uuid string, timepoint uint64) (float64, float64, bool, error) {
	result, err := self.interf.OrderStatus(uuid)
	if err != nil {
		return 0, 0, false, err
	} else {
		remaining := result.Result.QuantityRemaining
		done := result.Result.Quantity - remaining
		return done, remaining, remaining < BITTREX_EPSILON, nil
	}
}

func (self *Bittrex) Trade(tradeType string, base common.Token, quote common.Token, rate float64, amount float64, timepoint uint64) (string, float64, float64, bool, error) {
	result, err := self.interf.Trade(tradeType, base, quote, rate, amount)

	if err != nil {
		return "", 0, 0, false, errors.New("Trade rejected by Bittrex")
	} else {
		if result.Success {
			uuid := result.Result["uuid"]
			done, remaining, finished, err := self.QueryOrder(
				uuid, timepoint+20)
			return uuid, done, remaining, finished, err
		} else {
			return "", 0, 0, false, errors.New(result.Error)
		}
	}
}

func (self *Bittrex) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) (string, error) {
	resp, err := self.interf.Withdraw(token, amount, address)
	if err != nil {
		return "", err
	} else {
		if resp.Success {
			return resp.Result["uuid"], nil
		} else {
			return "", errors.New(resp.Error)
		}
	}
}

func bitttimestampToUint64(input string) uint64 {
	var t time.Time
	var err error
	len := len(input)
	if len == 23 {
		t, err = time.Parse("2006-01-02T15:04:05.000", input)
	} else if len == 22 {
		t, err = time.Parse("2006-01-02T15:04:05.00", input)
	} else if len == 21 {
		t, err = time.Parse("2006-01-02T15:04:05.0", input)
	}
	if err != nil {
		panic(err)
	}
	return uint64(t.UnixNano() / int64(time.Millisecond))
}

func (self *Bittrex) DepositStatus(
	id common.ActivityID, txHash, currency string, amount float64, timepoint uint64) (string, error) {
	timestamp := id.Timepoint
	idParts := strings.Split(id.EID, "|")
	if len(idParts) != 3 {
		// here, the exchange id part in id is malformed
		// 1. because analytic didn't pass original ID
		// 2. id is not constructed correctly in a form of uuid + "|" + token + "|" + amount
		return "", errors.New("Invalid deposit id")
	}
	amount, err := strconv.ParseFloat(idParts[2], 64)
	if err != nil {
		panic(err)
	}
	histories, err := self.interf.DepositHistory(currency)
	if err != nil {
		return "", err
	} else {
		for _, deposit := range histories.Result {
			log.Printf("Bittrex deposit history check: %v %v %v %v",
				deposit.Currency == currency,
				deposit.Amount-amount < BITTREX_EPSILON,
				bitttimestampToUint64(deposit.LastUpdated) > timestamp/uint64(time.Millisecond),
				self.storage.IsNewBittrexDeposit(deposit.Id, id),
			)
			log.Printf("deposit.Currency: %s", deposit.Currency)
			log.Printf("currency: %s", currency)
			log.Printf("deposit.Amount: %f", deposit.Amount)
			log.Printf("amount: %f", amount)
			log.Printf("deposit.LastUpdated: %d", bitttimestampToUint64(deposit.LastUpdated))
			log.Printf("timestamp: %d", timestamp/uint64(time.Millisecond))
			log.Printf("is new deposit: %t", self.storage.IsNewBittrexDeposit(deposit.Id, id))
			if deposit.Currency == currency &&
				deposit.Amount-amount < BITTREX_EPSILON &&
				bitttimestampToUint64(deposit.LastUpdated) > timestamp/uint64(time.Millisecond) &&
				self.storage.IsNewBittrexDeposit(deposit.Id, id) {
				if err := self.storage.RegisterBittrexDeposit(deposit.Id, id); err != nil {
					log.Printf("Register bittrex deposit error: %s", err.Error())
				}
				return "done", nil
			}
		}
		return "", nil
	}
}

func (self *Bittrex) CancelOrder(id, base, quote string) error {
	resp, err := self.interf.CancelOrder(id)
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

func (self *Bittrex) WithdrawStatus(id, currency string, amount float64, timepoint uint64) (string, string, error) {
	histories, err := self.interf.WithdrawHistory(currency)
	if err != nil {
		return "", "", err
	} else {
		for _, withdraw := range histories.Result {
			if withdraw.PaymentUuid == id {
				if withdraw.PendingPayment {
					return "", withdraw.TxId, nil
				} else {
					return "done", withdraw.TxId, nil
				}
			}
		}
		log.Printf("Withdraw with uuid " + id + " of currency " + currency + " is not found on bittrex")
		return "", "", nil
	}
}

func (self *Bittrex) OrderStatus(uuid string, base, quote string) (string, error) {
	resp_data, err := self.interf.OrderStatus(uuid)
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

func (self *Bittrex) FetchOnePairData(wq *sync.WaitGroup, pair common.TokenPair, data *sync.Map, timepoint uint64) {
	defer wq.Done()
	result := common.ExchangePrice{}
	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	onePairData, err := self.interf.FetchOnePairData(pair)
	returnTime := common.GetTimestamp()
	result.ReturnTime = returnTime
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		if !onePairData.Success {
			result.Valid = false
			result.Error = onePairData.Msg
		} else {
			for _, buy := range onePairData.Result["buy"] {
				result.Bids = append(
					result.Bids,
					common.NewPriceEntry(
						buy["Quantity"],
						buy["Rate"],
					),
				)
			}
			for _, sell := range onePairData.Result["sell"] {
				result.Asks = append(
					result.Asks,
					common.NewPriceEntry(
						sell["Quantity"],
						sell["Rate"],
					),
				)
			}
		}
	}
	data.Store(pair.PairID(), result)
}

func (self *Bittrex) FetchPriceData(timepoint uint64) (map[common.TokenPairID]common.ExchangePrice, error) {
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

func (self *Bittrex) FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	resp_data, err := self.interf.GetInfo()
	result.ReturnTime = common.GetTimestamp()
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
		result.Status = false
	} else {
		result.AvailableBalance = map[string]float64{}
		result.LockedBalance = map[string]float64{}
		result.DepositBalance = map[string]float64{}
		result.Status = true
		if resp_data.Success {
			for _, b := range resp_data.Result {
				tokenID := b.Currency
				_, err := common.GetInternalToken(tokenID)
				if err == nil {
					result.AvailableBalance[tokenID] = b.Available
					result.DepositBalance[tokenID] = b.Pending
					result.LockedBalance[tokenID] = 0
				}
			}
			// check if bittrex returned balance for all of the
			// supported token.
			// If it didn't, it is considered invalid
			if len(result.AvailableBalance) != len(self.tokens) {
				result.Valid = false
				result.Error = "Bittrex didn't return balance for all supported tokens"
			}
		} else {
			result.Valid = false
			result.Error = resp_data.Error
		}
	}
	return result, nil
}

func (self *Bittrex) FetchOnePairTradeHistory(
	wait *sync.WaitGroup,
	data *sync.Map,
	pair common.TokenPair,
	timepoint uint64) {

	defer wait.Done()
	result := []common.TradeHistory{}
	resp, err := self.interf.GetAccountTradeHistory(pair.Base, pair.Quote)
	if err != nil {
		log.Printf("Cannot fetch data for pair %s%s: %s", pair.Base.ID, pair.Quote.ID, err.Error())
	}
	for _, trade := range resp.Result {
		t, _ := time.Parse("2014-07-09T04:01:00.667", trade.TimeStamp)
		historyType := "sell"
		if trade.OrderType == "LIMIT_BUY" {
			historyType = "buy"
		}
		tradeHistory := common.NewTradeHistory(
			trade.OrderUuid,
			trade.Price,
			trade.Quantity,
			historyType,
			common.TimeToTimepoint(t),
		)
		result = append(result, tradeHistory)
	}
	pairString := pair.PairID()
	data.Store(pairString, result)
}

func (self *Bittrex) FetchTradeHistory() {
	t := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			result := map[common.TokenPairID][]common.TradeHistory{}
			timepoint := common.GetTimepoint()
			data := sync.Map{}
			pairs := self.pairs
			wait := sync.WaitGroup{}
			for _, pair := range pairs {
				wait.Add(1)
				go self.FetchOnePairTradeHistory(&wait, &data, pair, timepoint)
			}
			wait.Wait()
			data.Range(func(key, value interface{}) bool {
				result[key.(common.TokenPairID)] = value.([]common.TradeHistory)
				return true
			})
			if err := self.storage.StoreTradeHistory(result); err != nil {
				log.Printf("Bittrex store trade history error: %s", err.Error())
			}
			<-t.C
		}
	}()
}

func (self *Bittrex) GetTradeHistory(fromTime, toTime uint64) (common.ExchangeTradeHistory, error) {
	return self.storage.GetTradeHistory(fromTime, toTime)
}

func NewBittrex(addressConfig map[string]string,
	feeConfig common.ExchangeFees,
	interf BittrexInterface,
	storage BittrexStorage,
	minDepositConfig common.ExchangesMinDeposit) *Bittrex {
	tokens, pairs, fees, minDeposit := getExchangePairsAndFeesFromConfig(addressConfig, feeConfig, minDepositConfig, "bittrex")
	bittrex := &Bittrex{
		interf,
		pairs,
		tokens,
		common.NewExchangeAddresses(),
		storage,
		common.NewExchangeInfo(),
		fees,
		minDeposit,
	}
	bittrex.FetchTradeHistory()
	return bittrex
}
