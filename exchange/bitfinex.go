package exchange

import (
	"encoding/json"
	"math/big"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Bitfinex struct {
	interf    BitfinexInterface
	pairs     []common.TokenPair
	addresses map[string]ethereum.Address
}

func (self *Bitfinex) MarshalText() (text []byte, err error) {
	return []byte(self.ID()), nil
}

func (self *Bitfinex) Address(token common.Token) (ethereum.Address, bool) {
	addr, supported := self.addresses[token.ID]
	return addr, supported
}

func (self *Bitfinex) ID() common.ExchangeID {
	return common.ExchangeID("bitfinex")
}

func (self *Bitfinex) TokenPairs() []common.TokenPair {
	return self.pairs
}

func (self *Bitfinex) Name() string {
	return "bitfinex"
}

func (self *Bitfinex) Trade(tradeType string, base common.Token, quote common.Token, rate float64, amount float64, timepoint uint64) (done float64, remaining float64, finished bool, err error) {
	return self.interf.Trade(tradeType, base, quote, rate, amount, timepoint)
}

func (self *Bitfinex) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) error {
	return self.interf.Withdraw(token, amount, address, timepoint)
}

func (self *Bitfinex) FetchOnePairData(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map) {

	defer wg.Done()
	result := common.ExchangePrice{}

	client := &http.Client{}
	url := fmt.Sprintf(
		"https://api.bitfinex.com/v1/book/%s%s?group=1&limit_bids=50&limit_asks=50",
		strings.ToLower(pair.Base.ID),
		strings.ToLower(pair.Quote.ID),
	)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")

	timestamp := common.GetTimestamp()
	resp, err := client.Do(req)
	result.Timestamp = timestamp
	result.Valid = true
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		returnTime := common.GetTimestamp()
		result.ReturnTime = returnTime
		if err != nil {
			result.Valid = false
			result.Error = err.Error()
		} else {
			resp_data := Bitfresp{}
			json.Unmarshal(resp_body, &resp_data)
			if len(resp_data.Asks) == 0 && len(resp_data.Bids) == 0 {
				result.Valid = false
			} else {
				for _, buy := range resp_data.Bids {
					quantity, _ := strconv.ParseFloat(buy["amount"], 64)
					rate, _ := strconv.ParseFloat(buy["price"], 64)
					result.Bids = append(
						result.Bids,
						common.PriceEntry{
							quantity,
							rate,
						},
					)
				}
				for _, sell := range resp_data.Asks {
					quantity, _ := strconv.ParseFloat(sell["amount"], 64)
					rate, _ := strconv.ParseFloat(sell["price"], 64)
					result.Asks = append(
						result.Asks,
						common.PriceEntry{
							quantity,
							rate,
						},
					)
				}
			}
		}
	}
	data.Store(pair.PairID(), result)
}

func (self *Bitfinex) FetchPriceData(timepoint uint64) (map[common.TokenPairID]common.ExchangePrice, error) {
	wait := sync.WaitGroup{}
	data := sync.Map{}
	pairs := self.pairs
	for _, pair := range pairs {
		wait.Add(1)
		go self.FetchOnePairData(&wait, pair, &data)
	}
	wait.Wait()
	result := map[common.TokenPairID]common.ExchangePrice{}
	data.Range(func(key, value interface{}) bool {
		result[key.(common.TokenPairID)] = value.(common.ExchangePrice)
		return true
	})
	// fmt.Printf("result: %v\n", result)
	return result, nil
}

func (self *Bitfinex) FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
	return result, nil
}

func NewBitfinex(interf BitfinexInterface) *Bitfinex {
	return &Bitfinex{
		interf,
		[]common.TokenPair{
			common.MustCreateTokenPair("OMG", "ETH"),
			common.MustCreateTokenPair("EOS", "ETH"),
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
