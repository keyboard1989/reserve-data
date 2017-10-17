package exchanges

import (
	"encoding/json"
	"fmt"
	"github.com/KyberNetwork/reserve-data/fetcher"
	"github.com/KyberNetwork/reserve-data/market"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Bitfinex struct {
}

func (self Bitfinex) ID() string {
	return "bitfinex"
}

func (self Bitfinex) Name() string {
	return "bitfinex"
}

type bitfresp struct {
	Asks []map[string]string `json:"asks"`
	Bids []map[string]string `json:"bids"`
}

func (self Bitfinex) FetchOnePairData(wg *sync.WaitGroup, pair fetcher.TokenPair, data *ConcurrentResultMap) {
	defer wg.Done()
	result := market.ExchangeData{}

	client := &http.Client{}
	url := fmt.Sprintf(
		"https://api.bitfinex.com/v1/book/%s%s?group=1&limit_bids=50&limit_asks=50",
		strings.ToLower(pair.Base.Symbol()),
		strings.ToLower(pair.Quote.Symbol()),
	)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	result.Timestamp = strconv.Itoa(int(timestamp))
	result.Valid = true
	if err != nil {
		result.Valid = false
	} else {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			result.Valid = false
		} else {
			resp_data := bitfresp{}
			json.Unmarshal(resp_body, &resp_data)
			if len(resp_data.Asks) == 0 && len(resp_data.Bids) == 0 {
				result.Valid = false
			} else {
				for _, buy := range resp_data.Bids {
					quantity, _ := strconv.ParseFloat(buy["amount"], 64)
					rate, _ := strconv.ParseFloat(buy["price"], 64)
					result.BuyPrices = append(
						result.BuyPrices,
						market.Price{
							quantity,
							rate,
						},
					)
				}
				for _, sell := range resp_data.Asks {
					quantity, _ := strconv.ParseFloat(sell["amount"], 64)
					rate, _ := strconv.ParseFloat(sell["price"], 64)
					result.SellPrices = append(
						result.SellPrices,
						market.Price{
							quantity,
							rate,
						},
					)
				}
			}
		}
	}
	data.Set(pair.PairString(), &result)
}

func (self Bitfinex) FetchData(pairs []fetcher.TokenPair) (string, map[string]*market.ExchangeData) {
	wait := sync.WaitGroup{}
	result := NewConcurrentResultMap()
	for _, pair := range pairs {
		wait.Add(1)
		go self.FetchOnePairData(&wait, pair, result)
	}
	wait.Wait()
	fmt.Printf("result: %v\n", result.GetData())
	return "timestamp", result.GetData()
}
