package exchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
)

type Bitfinex struct {
	pairs []common.TokenPair
}

func (self Bitfinex) ID() common.ExchangeID {
	return common.ExchangeID("bitfinex")
}

func (self Bitfinex) Name() string {
	return "bitfinex"
}

type bitfresp struct {
	Asks []map[string]string `json:"asks"`
	Bids []map[string]string `json:"bids"`
}

func (self Bitfinex) FetchOnePairData(
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
			resp_data := bitfresp{}
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

func (self Bitfinex) FetchPriceData() (map[common.TokenPairID]common.ExchangePrice, error) {
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

func NewBitfinex() *Bitfinex {
	return &Bitfinex{
		[]common.TokenPair{
			common.MustCreateTokenPair("OMG", "ETH"),
			common.MustCreateTokenPair("EOS", "ETH"),
		},
	}
}
