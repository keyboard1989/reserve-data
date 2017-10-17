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

type Bittrex struct {
}

func (self Bittrex) ID() string {
	return "bittrex"
}

func (self Bittrex) Name() string {
	return "bittrex"
}

type bitresp struct {
	Success bool                            `json:"success"`
	Msg     string                          `json:"message"`
	Result  map[string][]map[string]float64 `json:"result"`
}

func (self Bittrex) FetchOnePairData(wg *sync.WaitGroup, pair fetcher.TokenPair, data *ConcurrentResultMap) {
	defer wg.Done()
	result := market.ExchangeData{}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://bittrex.com/api/v1.1/public/getorderbook", nil)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("market", fmt.Sprintf("%s-%s", pair.Quote.Symbol(), pair.Base.Symbol()))
	q.Add("type", "both")
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
			resp_data := bitresp{}
			json.Unmarshal(resp_body, &resp_data)
			if !resp_data.Success {
				result.Valid = false
			} else {
				for _, buy := range resp_data.Result["buy"] {
					result.BuyPrices = append(
						result.BuyPrices,
						market.Price{
							buy["Quantity"],
							buy["Rate"],
						},
					)
				}
				for _, sell := range resp_data.Result["sell"] {
					result.SellPrices = append(
						result.SellPrices,
						market.Price{
							sell["Quantity"],
							sell["Rate"],
						},
					)
				}
			}
		}
	}
	data.Set(pair.PairString(), &result)
}

func (self Bittrex) FetchData(pairs []fetcher.TokenPair) (string, map[string]*market.ExchangeData) {
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
