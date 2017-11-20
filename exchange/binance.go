package exchange

import (
	"fmt"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Binance struct {
	interf    BinanceInterface
	pairs     []common.TokenPair
	addresses map[string]ethereum.Address
}

func (self *Binance) MarshalText() (text []byte, err error) {
	return []byte(self.ID()), nil
}

func (self *Binance) Address(token common.Token) (ethereum.Address, bool) {
	addr, supported := self.addresses[token.ID]
	return addr, supported
}

func (self *Binance) UpdateAllDepositAddresses(address string) {
	for k, _ := range self.addresses {
		self.addresses[k] = ethereum.HexToAddress(address)
	}
}

func (self *Binance) UpdateDepositAddress(token common.Token, address string) {
	self.addresses[token.ID] = ethereum.HexToAddress(address)
}

func (self *Binance) ID() common.ExchangeID {
	return common.ExchangeID("binance")
}

func (self *Binance) TokenPairs() []common.TokenPair {
	return self.pairs
}

func (self *Binance) Name() string {
	return "binance"
}

func (self *Binance) Trade(tradeType string, base common.Token, quote common.Token, rate float64, amount float64, timepoint uint64) (done float64, remaining float64, finished bool, err error) {
	return self.interf.Trade(tradeType, base, quote, rate, amount, timepoint)
}

func (self *Binance) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) (ethereum.Hash, error) {
	return self.interf.Withdraw(token, amount, address, timepoint)
}

func (self Binance) FetchPriceData(timepoint uint64) (map[common.TokenPairID]common.ExchangePrice, error) {
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

func (self *Binance) FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error) {
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
		if resp_data.Code != 0 {
			result.Valid = false
			result.Error = fmt.Sprintf("Code: %s, Msg: %s", resp_data.Code, resp_data.Msg)
		} else {
			for _, b := range resp_data.Balances {
				tokenID := b.Asset
				_, exist := common.SupportedTokens[tokenID]
				if exist {
					result.AvailableBalance[tokenID] = b.Free
					result.LockedBalance[tokenID] = b.Locked
					result.DepositBalance[tokenID] = 0
				}
			}
		}
	}
	return result, nil
}

func NewBinance(interf BinanceInterface) *Binance {
	return &Binance{
		interf,
		[]common.TokenPair{
			common.MustCreateTokenPair("FUN", "ETH"),
			common.MustCreateTokenPair("MCO", "ETH"),
			common.MustCreateTokenPair("OMG", "ETH"),
			common.MustCreateTokenPair("EOS", "ETH"),
			common.MustCreateTokenPair("KNC", "ETH"),
			common.MustCreateTokenPair("LINK", "ETH"),
		},
		map[string]ethereum.Address{
			"ETH": ethereum.HexToAddress("0x5b082bc7928e1fd5b757426fe7225cc7a6a75c55"),
			"OMG": ethereum.HexToAddress("0x5b082bc7928e1fd5b757426fe7225cc7a6a75c55"),
			"DGD": ethereum.HexToAddress("0x5b082bc7928e1fd5b757426fe7225cc7a6a75c55"),
			"CVC": ethereum.HexToAddress("0x5b082bc7928e1fd5b757426fe7225cc7a6a75c55"),
			"MCO": ethereum.HexToAddress("0x5b082bc7928e1fd5b757426fe7225cc7a6a75c55"),
			"GNT": ethereum.HexToAddress("0x5b082bc7928e1fd5b757426fe7225cc7a6a75c55"),
			"ADX": ethereum.HexToAddress("0x5b082bc7928e1fd5b757426fe7225cc7a6a75c55"),
			"EOS": ethereum.HexToAddress("0x5b082bc7928e1fd5b757426fe7225cc7a6a75c55"),
			"PAY": ethereum.HexToAddress("0x5b082bc7928e1fd5b757426fe7225cc7a6a75c55"),
			"BAT": ethereum.HexToAddress("0x5b082bc7928e1fd5b757426fe7225cc7a6a75c55"),
			"KNC": ethereum.HexToAddress("0x5b082bc7928e1fd5b757426fe7225cc7a6a75c55"),
		},
	}
}
