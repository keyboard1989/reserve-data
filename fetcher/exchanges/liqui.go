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
	"time"
)

type Liqui struct {
}

func (self Liqui) ID() string {
	return "liqui"
}

func (self Liqui) Name() string {
	return "liqui"
}

// map of token pair to map of asks/bids to array of [rate, amount]
type liqresp map[string]map[string][][]float64

func (self Liqui) FetchData(pairs []fetcher.TokenPair) (string, map[string]*market.ExchangeData) {
	result := NewConcurrentResultMap()

	client := &http.Client{}
	pairs_str := []string{}
	for _, pair := range pairs {
		pairs_str = append(pairs_str, fmt.Sprintf("%s_%s", pair.Base.Symbol(), pair.Quote.Symbol()))
	}
	url := fmt.Sprintf(
		"https://api.liqui.io/api/3/depth/%s?ignore_invalid=1",
		strings.ToLower(strings.Join(pairs_str, "-")),
	)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	if err != nil {
		for _, pair := range pairs {
			one_pair_result := market.ExchangeData{}
			one_pair_result.Valid = false
			one_pair_result.Timestamp = strconv.Itoa(int(timestamp))
			result.Set(pair.PairString(), &one_pair_result)
		}
	} else {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			for _, pair := range pairs {
				one_pair_result := market.ExchangeData{}
				one_pair_result.Valid = false
				one_pair_result.Timestamp = strconv.Itoa(int(timestamp))
				result.Set(pair.PairString(), &one_pair_result)
			}
		} else {
			resp_data := liqresp{}
			json.Unmarshal(resp_body, &resp_data)
			for _, pair := range pairs {
				one_pair_result := market.ExchangeData{}
				one_pair_result.Valid = true
				one_pair_result.Timestamp = strconv.Itoa(int(timestamp))
				one_data := resp_data[fmt.Sprintf(
					"%s_%s",
					strings.ToLower(pair.Base.Symbol()),
					strings.ToLower(pair.Quote.Symbol()),
				)]
				for _, buy := range one_data["bids"] {
					one_pair_result.BuyPrices = append(
						one_pair_result.BuyPrices,
						market.Price{
							buy[1],
							buy[0],
						},
					)
				}
				for _, sell := range one_data["asks"] {
					one_pair_result.SellPrices = append(
						one_pair_result.SellPrices,
						market.Price{
							sell[1],
							sell[0],
						},
					)
				}
				result.Set(pair.PairString(), &one_pair_result)
			}
		}
	}
	fmt.Printf("result: %v\n", result.GetData())
	return "timestamp", result.GetData()
}
