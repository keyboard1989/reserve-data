package exchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
)

type Bittrex struct {
	pairs []common.TokenPair
}

func (self Bittrex) ID() common.ExchangeID {
	return common.ExchangeID("bittrex")
}

func (self Bittrex) Name() string {
	return "bittrex"
}

type bitresp struct {
	Success bool                            `json:"success"`
	Msg     string                          `json:"message"`
	Result  map[string][]map[string]float64 `json:"result"`
}

func (self Bittrex) FetchOnePairData(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map) {

	defer wg.Done()
	result := common.ExchangePrice{}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://bittrex.com/api/v1.1/public/getorderbook", nil)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("market", fmt.Sprintf("%s-%s", pair.Quote.ID, pair.Base.ID))
	q.Add("type", "both")
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
			resp_data := bitresp{}
			json.Unmarshal(resp_body, &resp_data)
			if !resp_data.Success {
				result.Valid = false
			} else {
				for _, buy := range resp_data.Result["buy"] {
					result.BuyPrices = append(
						result.BuyPrices,
						common.PriceEntry{
							returnTime,
							buy["Quantity"],
							buy["Rate"],
						},
					)
				}
				for _, sell := range resp_data.Result["sell"] {
					result.SellPrices = append(
						result.SellPrices,
						common.PriceEntry{
							returnTime,
							sell["Quantity"],
							sell["Rate"],
						},
					)
				}
			}
		}
	}
	data.Store(pair.PairID(), result)
}

func (self Bittrex) FetchPriceData() (map[common.TokenPairID]common.ExchangePrice, error) {
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

func NewBittrex() *Bittrex {
	return &Bittrex{
		[]common.TokenPair{
			common.MustCreateTokenPair("OMG", "ETH"),
			common.MustCreateTokenPair("DGD", "ETH"),
			common.MustCreateTokenPair("CVC", "ETH"),
			common.MustCreateTokenPair("FUN", "ETH"),
			common.MustCreateTokenPair("MCO", "ETH"),
			common.MustCreateTokenPair("GNT", "ETH"),
			common.MustCreateTokenPair("ADX", "ETH"),
			common.MustCreateTokenPair("PAY", "ETH"),
			common.MustCreateTokenPair("BAT", "ETH"),
		},
	}
}
