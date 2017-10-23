package exchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
)

type Binance struct {
	pairs []common.TokenPair
}

func (self Binance) ID() common.ExchangeID {
	return common.ExchangeID("binance")
}

func (self Binance) Name() string {
	return "binance"
}

type binresp struct {
	LastUpdatedId int64      `json:"lastUpdateId"`
	Code          int        `json:"code"`
	Msg           string     `json:"msg"`
	Bids          [][]string `json:"bids"`
	Asks          [][]string `json:"asks"`
}

func (self Binance) FetchOnePairData(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map) {

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
			resp_data := binresp{}
			json.Unmarshal(resp_body, &resp_data)
			if resp_data.Code != 0 || resp_data.Msg != "" {
				result.Valid = false
			} else {
				for _, buy := range resp_data.Bids {
					quantity, _ := strconv.ParseFloat(buy[1], 64)
					rate, _ := strconv.ParseFloat(buy[0], 64)
					result.BuyPrices = append(
						result.BuyPrices,
						common.PriceEntry{
							returnTime,
							quantity,
							rate,
						},
					)
				}
				for _, sell := range resp_data.Asks {
					quantity, _ := strconv.ParseFloat(sell[1], 64)
					rate, _ := strconv.ParseFloat(sell[0], 64)
					result.SellPrices = append(
						result.SellPrices,
						common.PriceEntry{
							returnTime,
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
func (self Binance) FetchPriceData() (map[common.TokenPairID]common.ExchangePrice, error) {
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
	fmt.Printf("result: %v\n", result)
	return result, nil
}

func NewBinance() *Binance {
	return &Binance{
		[]common.TokenPair{
			common.MustCreateTokenPair("FUN", "ETH"),
			common.MustCreateTokenPair("MCO", "ETH"),
			common.MustCreateTokenPair("OMG", "ETH"),
			common.MustCreateTokenPair("EOS", "ETH"),
			common.MustCreateTokenPair("KNC", "ETH"),
			common.MustCreateTokenPair("LINK", "ETH"),
		},
	}
}
