package exchange

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"sync"
)

type RealBinanceEndpoint struct {
	AuthenticatedEndpoint string
	PublicEndpoint        string
}

func (self *RealBinanceEndpoint) nonce() string {
	epsilon := 30 * time.Millisecond
	anchor := int64(50299954901)
	timestamp := time.Now().UnixNano()/int64(epsilon) - anchor
	return strconv.Itoa(int(timestamp))
}

// sign the request if signer passed
func (self *RealBinanceEndpoint) fillRequest(req *http.Request, signer Signer) {
	if req.Method == "POST" || req.Method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Add("Accept", "application/json")
	if signer != nil {
		q := req.URL.Query()
		timestamp := time.Now().UnixNano()/int64(time.Millisecond)
		q.Set("timestamp", strconv.FormatInt(timestamp, 10))
		q.Set("recvWindow", "5000")
		q.Set("signature", signer.BinanceSign(""))
		req.Header.Add("X-MBX-APIKEY", signer.GetBinanceKey())
	}
}

type binanceResponse struct {
	LastUpdatedId int64      `json:"lastUpdateId"`
	Code          int        `json:"code"`
	Msg           string     `json:"msg"`
	Bids          [][]string `json:"bids"`
	Asks          [][]string `json:"asks"`
}

func (self RealBinanceEndpoint) FetchOnePairData(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map) {

	defer wg.Done()
	result := common.ExchangePrice{}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", self.PublicEndpoint+"/v1/depth", nil)
	q := req.URL.Query()
	q.Add("symbol", fmt.Sprintf("%s%s", pair.Base.ID, pair.Quote.ID))
	q.Add("limit", "50")
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, nil)

	resp, err := client.Do(req)
	result.Timestamp = common.GetTimestamp()
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
			resp_data := binanceResponse{}
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
							quantity,
							rate,
						},
					)
				}
			}
		}
	}
	fmt.Sprintf("%v\n", result)
	data.Store(pair.PairID(), result)
}

type binanceTradeResult struct {
	Done      float64 `json:"received"`
	Remaining float64 `json:"remains"`
	OrderID   int64   `json:"uuid"`
}

type binanceTrade struct {
	Success int                `json:"success"`
	Return  binanceTradeResult `json:"result"`
	Error   string             `json:"message"`
}

func (self *RealBinanceEndpoint) Trade(key string, tradeType string, base, quote common.Token, rate, amount float64, signer Signer) (done float64, remaining float64, finished bool, err error) {
	result := binanceTrade{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	req, _ := http.NewRequest(
		"GET",
		self.AuthenticatedEndpoint+"/selllimit",
		nil,
	)
	q := req.URL.Query()
	q.Set("market", fmt.Sprintf("%s-%s", strings.ToUpper(base.ID), strings.ToUpper(quote.ID)))
	q.Set("quantity", fmt.Sprintf("%f", amount))
	q.Set("rate", fmt.Sprintf("%f", rate))
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, signer)
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("response: %s\n", resp_body)
		if err == nil {
			err = json.Unmarshal(resp_body, &result)
		}
		if err != nil {
			return 0, 0, false, err
		}
		if result.Error != "" {
			return 0, 0, false, errors.New(result.Error)
		}
		return result.Return.Done, result.Return.Remaining, result.Return.OrderID == 0, nil
	} else {
		fmt.Printf("Error: %v, Code: %v\n", err, resp)
		return 0, 0, false, errors.New("Trade rejected by Binance")
	}
}

type binanceWithdraw struct {
	Success int    `json:"success"`
	Return  map[string]interface{}
	Error   string `json:"message"`
}

func (self *RealBinanceEndpoint) Withdraw(key string, token common.Token, amount *big.Int, address ethereum.Address, signer Signer) error {
	result := binanceWithdraw{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, _ := http.NewRequest(
		"GET",
		self.AuthenticatedEndpoint+"/withdraw",
		nil,
	)
	q := req.URL.Query()
	q.Set("currency", strings.ToUpper(token.ID))
	q.Set("quantity", fmt.Sprintf("%f", common.BigToFloat(amount, token.Decimal)))
	q.Set("address", address.Hex())
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, signer)
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("response: %s\n", resp_body)
		if err == nil {
			err = json.Unmarshal(resp_body, &result)
		}
		if err != nil {
			return err
		}
		if result.Error != "" {
			return errors.New(result.Error)
		}
		return nil
	} else {
		fmt.Printf("Error: %v, Code: %v\n", err, resp)
		return errors.New("withdraw rejected by Binance")
	}
}

type binanceInfo map[string]map[string][][]float64

func (self *RealBinanceEndpoint) GetInfo(key string, signer Signer) (binanceInfo, error) {
	result := binanceInfo{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	data := url.Values{}
	data.Set("method", "getInfo")
	data.Add("nonce", nonce())
	params := data.Encode()
	req, _ := http.NewRequest(
		"POST",
		self.AuthenticatedEndpoint,
		bytes.NewBufferString(params),
	)
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", key)
	req.Header.Add("Sign", signer.LiquiSign(params))
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			json.Unmarshal(resp_body, &result)
		}
	}
	return result, err
}

func NewRealBinanceEndpoint() *RealBinanceEndpoint {
	return &RealBinanceEndpoint{
		"https://www.binance.com/api",
		"https://www.binance.com/api",
	}
}
