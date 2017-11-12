package bittrex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/exchange"
	ethereum "github.com/ethereum/go-ethereum/common"
	"sync"
)

type BittrexEndpoint struct {
	signer Signer
	interf Interface
}

func nonce() string {
	epsilon := 30 * time.Millisecond
	anchor := int64(50299954901)
	timestamp := time.Now().UnixNano()/int64(epsilon) - anchor
	return strconv.Itoa(int(timestamp))
}

func (self *BittrexEndpoint) fillRequest(req *http.Request, signNeeded bool) {
	if req.Method == "POST" || req.Method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Add("Accept", "application/json")
	if signNeeded == true {
		q := req.URL.Query()
		q.Set("apiKey", self.signer.GetBittrexKey())
		q.Set("nonce", nonce())
		req.URL.RawQuery = q.Encode()
		req.Header.Add("apisign", self.signer.BittrexSign(req.URL.String()))
	}
}

func (self *BittrexEndpoint) FetchOnePairData(wq *sync.WaitGroup, pair common.TokenPair, data *sync.Map, timepoint uint64) {
	defer wq.Done()
	result := common.ExchangePrice{}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", self.interf.PublicEndpoint(timepoint)+"/getorderbook", nil)
	q := req.URL.Query()
	q.Set("market", fmt.Sprintf("%s-%s", pair.Quote.ID, pair.Base.ID))
	q.Set("type", "both")
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, false)
	res, err := client.Do(req)

	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		result.ReturnTime = common.GetTimestamp()
		if err != nil {
			result.Valid = false
			result.Error = err.Error()
		} else {
			data := exchange.Bittresp{}
			json.Unmarshal(body, &data)
			if !data.Success {
				result.Valid = false
			} else {
				for _, buy := range data.Result["buy"] {
					result.Bids = append(
						result.Bids,
						common.PriceEntry{
							buy["Quantity"],
							buy["Rate"],
						},
					)
				}
				for _, sell := range data.Result["sell"] {
					result.Asks = append(
						result.Asks,
						common.PriceEntry{
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

func (self *BittrexEndpoint) Trade(tradeType string, base, quote common.Token, rate, amount float64, timepoint uint64) (done float64, remaining float64, finished bool, err error) {
	result := exchange.Bitttrade{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	req, _ := http.NewRequest(
		"GET",
		self.interf.MarketEndpoint(timepoint)+"/selllimit",
		nil,
	)
	q := req.URL.Query()
	q.Set("market", fmt.Sprintf("%s-%s", strings.ToUpper(base.ID), strings.ToUpper(quote.ID)))
	q.Set("quantity", fmt.Sprintf("%f", amount))
	q.Set("rate", fmt.Sprintf("%f", rate))
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, true)
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		log.Printf("response: %s\n", resp_body)
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
		log.Printf("Error: %v, Code: %v\n", err, resp)
		return 0, 0, false, errors.New("Trade rejected by Bittrex")
	}
}

func (self *BittrexEndpoint) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) error {
	// ignoring timepoint because it's only relevant in simulation
	result := exchange.Bittwithdraw{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, _ := http.NewRequest(
		"GET",
		self.interf.MarketEndpoint(timepoint)+"/withdraw",
		nil,
	)
	q := req.URL.Query()
	q.Set("currency", strings.ToUpper(token.ID))
	q.Set("quantity", fmt.Sprintf("%f", common.BigToFloat(amount, token.Decimal)))
	q.Set("address", address.Hex())
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, true)
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		log.Printf("response: %s\n", resp_body)
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
		log.Printf("Error: %v, Code: %v\n", err, resp)
		return errors.New("withdraw rejected by Bittrex")
	}
}

func (self *BittrexEndpoint) GetInfo(timepoint uint64) (exchange.Bittinfo, error) {
	result := exchange.Bittinfo{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	req, _ := http.NewRequest(
		"GET",
		self.interf.AccountEndpoint(timepoint)+"/getbalances",
		nil,
	)
	self.fillRequest(req, true)
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

func NewBittrexEndpoint(signer Signer, interf Interface) *BittrexEndpoint {
	return &BittrexEndpoint{signer, interf}
}

func NewRealBittrexEndpoint(signer Signer) *BittrexEndpoint {
	return &BittrexEndpoint{signer, NewRealInterface()}
}

func NewSimulatedBittrexEndpoint(signer Signer) *BittrexEndpoint {
	return &BittrexEndpoint{signer, NewSimulatedInterface()}
}
