package exchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
)

type Liqui struct {
	signer Signer
	pairs  []common.TokenPair
}

func (self *Liqui) ID() common.ExchangeID {
	return common.ExchangeID("liqui")
}

func (self *Liqui) Name() string {
	return "liqui"
}

// map of token pair to map of asks/bids to array of [rate, amount]
type liqresp map[string]map[string][][]float64

func nonce() string {
	epsilon := 30 * time.Millisecond
	anchor := int64(50299954901)
	timestamp := time.Now().UnixNano()/int64(epsilon) - anchor
	return strconv.Itoa(int(timestamp))
}

type liqinfo struct {
	Success int `json:"success"`
	Return  map[string]map[string]float64
	Error   string `json:"error"`
}

func (self *Liqui) FetchEBalanceData() (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
	client := &http.Client{}
	data := url.Values{}
	data.Set("method", "getInfo")
	data.Add("nonce", nonce())
	params := data.Encode()
	url := "https://api.liqui.io/tapi"
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(params))
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Key", self.signer.GetLiquiKey())
	signed := self.signer.LiquiSign(params)
	req.Header.Add("Sign", signed)
	timestamp := common.GetTimestamp()
	resp, err := client.Do(req)
	result.Timestamp = timestamp
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
			resp_data := liqinfo{}
			json.Unmarshal(resp_body, &resp_data)
			if resp_data.Success == 1 {
				balances := resp_data.Return["funds"]
				result.Balance = map[string]float64{}
				for tokenID, _ := range common.SupportedTokens {
					result.Balance[tokenID] = balances[strings.ToLower(tokenID)]
				}
			} else {
				result.Valid = false
				result.Error = resp_data.Error
			}
		}
	}
	return result, nil
}

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
	if err != nil {
		for _, pair := range pairs {
			one_pair_result := common.ExchangePrice{}
			one_pair_result.Timestamp = timestamp
			one_pair_result.Valid = false
			one_pair_result.Error = err.Error()
			result[pair.PairID()] = one_pair_result
		}
	} else {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		returnTime := common.GetTimestamp()
		if err != nil {
			for _, pair := range pairs {
				one_pair_result := common.ExchangePrice{}
				one_pair_result.Timestamp = timestamp
				one_pair_result.ReturnTime = returnTime
				one_pair_result.Valid = false
				one_pair_result.Error = err.Error()
				result[pair.PairID()] = one_pair_result
			}
		} else {
			resp_data := liqresp{}
			json.Unmarshal(resp_body, &resp_data)
			for _, pair := range pairs {
				one_pair_result := common.ExchangePrice{}
				one_pair_result.Timestamp = timestamp
				one_pair_result.ReturnTime = returnTime
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
							buy[1],
							buy[0],
						},
					)
				}
				for _, sell := range one_data["asks"] {
					one_pair_result.SellPrices = append(
						one_pair_result.SellPrices,
						common.PriceEntry{
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

func NewLiqui(signer Signer) *Liqui {
	return &Liqui{
		signer,
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
