package exchange

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Dgx struct {
	pairs        []common.TokenPair
	addresses    map[string]ethereum.Address
	exchangeInfo *common.ExchangeInfo
	fees         common.ExchangeFees
}

func (self *Dgx) TokenAddresses() map[string]ethereum.Address {
	return self.addresses.GetData()
}

func (self *Dgx) MarshalText() (text []byte, err error) {
	return []byte(self.ID()), nil
}

func (self *Dgx) Address(token common.Token) (ethereum.Address, bool) {
	addr, supported := self.addresses.Get(token.ID)
	return addr, supported
}

func (self *Dgx) UpdateAllDepositAddresses(address string) {
	panic("dgx doesn't support update deposit addresses")
}

func (self *Dgx) UpdateDepositAddress(token common.Token, address string) {
	panic("dgx doesn't support update deposit addresses")
}

func (self *Dgx) GetInfo() (common.ExchangeInfo, error) {
	return *self.exchangeInfo, nil
}

func (self *Dgx) GetExchangeInfo(pair common.TokenPairID) (common.ExchangePrecisionLimit, error) {
	data, err := self.exchangeInfo.Get(pair)
	return data, err
}

func (self *Dgx) GetFee() common.ExchangeFees {
	return self.fees
}

func (self *Dgx) ID() common.ExchangeID {
	return common.ExchangeID("dgx")
}

func (self *Dgx) TokenPairs() []common.TokenPair {
	return self.pairs
}

func (self *Dgx) Name() string {
	return "digix gold"
}

func (self *Dgx) QueryOrder(symbol string, id uint64) (done float64, remaining float64, finished bool, err error) {
	// TODO: see if trade order (a tx to dgx contract) is successful or not
	// - successful: done = order amount, remaining = 0, finished = true, err = nil
	// - failed: done = 0, remaining = order amount, finished = false, err = some error
	// - pending: done = 0, remaining = order amount, finished = false, err = nil
}

func (self *Dgx) Trade(tradeType string, base common.Token, quote common.Token, rate float64, amount float64, timepoint uint64) (id string, done float64, remaining float64, finished bool, err error) {
	// TODO: communicate with dgx connector to do the trade
}

func (self *Dgx) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) (string, error) {
	// TODO: communicate with dgx connector to withdraw
}

func (self *Dgx) CancelOrder(id common.ActivityID) error {
	return errors.New("Dgx doesn't support trade cancelling")
}

func (self *Dgx) FetchPriceData(timepoint uint64) (map[common.TokenPairID]common.ExchangePrice, error) {
	result := map[common.TokenPairID]common.ExchangePrice{}
	// TODO: Get price data from dgx connector and construct valid orderbooks
	return result, nil
}

func (self *Dgx) FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	// TODO: Get balance data from dgx connector
	result.ReturnTime = common.GetTimestamp()
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		result.AvailableBalance = map[string]float64{}
		result.LockedBalance = map[string]float64{}
		result.DepositBalance = map[string]float64{}
		// TODO: set proper balance to result
	}
	return result, nil
}

func (self *Dgx) FetchTradeHistory(timepoint uint64) (map[common.TokenPairID][]common.TradeHistory, error) {
	result := map[common.TokenPairID][]common.TradeHistory{}
	// TODO: get trade history
	return result, nil
}

func (self *Dgx) DepositStatus(id common.ActivityID, txHash, currency string, amount float64, timepoint uint64) (string, error) {
	// TODO: checking txHash status
	return "", nil
}

func (self *Dgx) WithdrawStatus(id, currency string, amount float64, timepoint uint64) (string, string, error) {
	// TODO: checking id (id is the txhash) status
	return "", nil
}

func (self *Dgx) OrderStatus(id common.ActivityID, timepoint uint64) (string, error) {
	// TODO: checking id (id is the txhash) status
	return "", nil
}

func NewDgx(addressConfig map[string]string, feeConfig common.ExchangeFees) *Dgx {
	pairs, fees := getExchangePairsAndFeesFromConfig(addressConfig, feeConfig, "dgx")
	depositAddrs := map[string]ethereum.Address{}
	for tok, addr := range addressConfig {
		depositAddrs[tok] = ethereum.HexToAddress(addr)
	}
	return &Dgx{
		pairs,
		depositAddrs,
		common.NewExchangeInfo(),
		fees,
	}
}
