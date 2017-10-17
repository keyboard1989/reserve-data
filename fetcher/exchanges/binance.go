package exchanges

import (
	"encoding/json"
	"fmt"
	"github.com/KyberNetwork/reserve-data/fetcher"
	"github.com/KyberNetwork/reserve-data/market"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Binance struct {
}

func (self Binance) ID() string {
	return "binance"
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

func (self Binance) FetchOnePairData(wg *sync.WaitGroup, pair fetcher.TokenPair, data *ConcurrentResultMap) {
	defer wg.Done()
	result := market.ExchangeData{}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.binance.com/api/v1/depth", nil)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("symbol", fmt.Sprintf("%s%s", pair.Base.Symbol(), pair.Quote.Symbol()))
	q.Add("limit", "50")
	req.URL.RawQuery = q.Encode()

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
						market.Price{
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

// https://www.binance.com/api/v1/depth?symbol=OMGETH&limit=50
func (self Binance) FetchData(pairs []fetcher.TokenPair) (string, map[string]*market.ExchangeData) {
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
