package exchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/KyberNetwork/reserve-data/common"
)

type Liqui struct {
	pairs []common.TokenPair
}

func (self *Liqui) ID() common.ExchangeID {
	return common.ExchangeID("liqui")
}

func (self *Liqui) Name() string {
	return "liqui"
}

// map of token pair to map of asks/bids to array of [rate, amount]
type liqresp map[string]map[string][][]float64

func (self *Liqui) FetchPriceData() (map[common.TokenPairID]common.ExchangePrice, error) {
	result := map[common.TokenPairID]common.ExchangePrice{}
	pairs := self.pairs
	client := &http.Client{}
	pairs_str := []string{}
	for _, pair := range pairs {
		pairs_str = append(pairs_str, fmt.Sprintf("%s_%s", pair.Base.ID, pair.Quote.ID))
	}
	url := fmt.Sprintf(
		"https://api.liqui.io/api/3/depth/%s?ignore_invalid=1",
		strings.ToLower(strings.Join(pairs_str, "-")),
	)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	timestamp := common.GetTimestamp()
	resp, err := client.Do(req)
	one_pair_result := common.ExchangePrice{}
	one_pair_result.Timestamp = timestamp
	if err != nil {
		for _, pair := range pairs {
			one_pair_result.Valid = false
			one_pair_result.Error = err.Error()
			result[pair.PairID()] = one_pair_result
		}
	} else {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		returnTime := common.GetTimestamp()
		one_pair_result.ReturnTime = returnTime
		if err != nil {
			for _, pair := range pairs {
				one_pair_result.Valid = false
				one_pair_result.Error = err.Error()
				result[pair.PairID()] = one_pair_result
			}
		} else {
			resp_data := liqresp{}
			json.Unmarshal(resp_body, &resp_data)
			for _, pair := range pairs {
				one_pair_result.Valid = true
				one_data := resp_data[fmt.Sprintf(
					"%s_%s",
					strings.ToLower(pair.Base.ID),
					strings.ToLower(pair.Quote.ID),
				)]
				for _, buy := range one_data["bids"] {
					one_pair_result.BuyPrices = append(
						one_pair_result.BuyPrices,
						common.PriceEntry{
							returnTime,
							buy[1],
							buy[0],
						},
					)
				}
				for _, sell := range one_data["asks"] {
					one_pair_result.SellPrices = append(
						one_pair_result.SellPrices,
						common.PriceEntry{
							returnTime,
							sell[1],
							sell[0],
						},
					)
				}
				result[pair.PairID()] = one_pair_result
			}
		}
	}
	return result, err
}

func NewLiqui() *Liqui {
	return &Liqui{
		[]common.TokenPair{
			common.MustCreateTokenPair("OMG", "ETH"),
			common.MustCreateTokenPair("DGD", "ETH"),
			common.MustCreateTokenPair("CVC", "ETH"),
			common.MustCreateTokenPair("MCO", "ETH"),
			common.MustCreateTokenPair("GNT", "ETH"),
			common.MustCreateTokenPair("ADX", "ETH"),
			common.MustCreateTokenPair("EOS", "ETH"),
			common.MustCreateTokenPair("PAY", "ETH"),
			common.MustCreateTokenPair("BAT", "ETH"),
			common.MustCreateTokenPair("KNC", "ETH"),
		},
	}
}
