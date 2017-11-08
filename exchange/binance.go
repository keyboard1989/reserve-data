package exchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

func (self *Binance) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) error {
	return self.interf.Withdraw(token, amount, address, timepoint)
}

func (self *Binance) FetchOnePairData(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map,
	timepoint uint64) {

	defer wg.Done()
	result := common.ExchangePrice{}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.binance.com/api/v1/depth", nil)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("symbol", fmt.Sprintf("%s%s", pair.Base.ID, pair.Quote.ID))
	q.Add("limit", "50")
	req.URL.RawQuery = q.Encode()

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
			resp_data := Binaresp{}
			json.Unmarshal(resp_body, &resp_data)
			if resp_data.Code != 0 || resp_data.Msg != "" {
				result.Valid = false
			} else {
				for _, buy := range resp_data.Bids {
					quantity, _ := strconv.ParseFloat(buy[1], 64)
					rate, _ := strconv.ParseFloat(buy[0], 64)
					result.Bids = append(
						result.Bids,
						common.PriceEntry{
							quantity,
							rate,
						},
					)
				}
				for _, sell := range resp_data.Asks {
					quantity, _ := strconv.ParseFloat(sell[1], 64)
					rate, _ := strconv.ParseFloat(sell[0], 64)
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

// https://www.binance.com/api/v1/depth?symbol=OMGETH&limit=50
func (self Binance) FetchPriceData(timepoint uint64) (map[common.TokenPairID]common.ExchangePrice, error) {
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
	// fmt.Printf("result: %v\n", result)
	return result, nil
}

func (self *Binance) FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
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
