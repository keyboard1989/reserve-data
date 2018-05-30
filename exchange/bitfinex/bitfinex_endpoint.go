package bitfinex

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/exchange"
	ethereum "github.com/ethereum/go-ethereum/common"
)

//BitfinexEndpoint object
type BitfinexEndpoint struct {
	signer Signer
	interf Interface
}

func nonce() string {
	epsilon := 30 * time.Millisecond
	anchor := int64(50299954901)
	timestamp := time.Now().UnixNano()/int64(epsilon) - anchor
	return strconv.Itoa(int(timestamp))
}

func (self *BitfinexEndpoint) fillRequest(req *http.Request, signNeeded bool, timepoint uint64) {
	if req.Method == "POST" || req.Method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Add("Accept", "application/json")
	if signNeeded == true {
		payload := map[string]interface{}{
			"request": req.URL.Path,
			"nonce":   fmt.Sprintf("%v", timepoint),
		}
		payloadJSON, _ := json.Marshal(payload)
		payloadEnc := base64.StdEncoding.EncodeToString(payloadJSON)
		req.Header.Add("X-BFX-APIKEY", self.signer.GetKey())
		req.Header.Add("X-BFX-PAYLOAD", payloadEnc)
		req.Header.Add("X-BFX-SIGNATURE", self.signer.Sign(req.URL.String()))
	}
}

//FetchOnePairData return one pair data
func (self *BitfinexEndpoint) FetchOnePairData(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map,
	timepoint uint64) {

	defer wg.Done()
	result := common.ExchangePrice{}

	client := &http.Client{}
	url := self.interf.PublicEndpoint() + fmt.Sprintf(
		"/book/%s%s",
		strings.ToLower(pair.Base.ID),
		strings.ToLower(pair.Quote.ID))
	req, _ := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Set("group", "1")
	q.Set("limit_bids", "50")
	q.Set("limit_asks", "50")
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, false, timepoint)

	resp, err := client.Do(req)
	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Response body close error: %s", err.Error())
			}
		}()
		respBody, err := ioutil.ReadAll(resp.Body)
		returnTime := common.GetTimestamp()
		result.ReturnTime = returnTime
		if err != nil {
			result.Valid = false
			result.Error = err.Error()
		} else {
			respData := exchange.Bitfresp{}
			if err := json.Unmarshal(respBody, &respData); err != nil {
				log.Printf("Unmarshal response error: %s", err.Error())
			}
			if len(respData.Asks) == 0 && len(respData.Bids) == 0 {
				result.Valid = false
			} else {
				for _, buy := range respData.Bids {
					quantity, _ := strconv.ParseFloat(buy["amount"], 64)
					rate, _ := strconv.ParseFloat(buy["price"], 64)
					result.Bids = append(
						result.Bids,
						common.NewPriceEntry(
							quantity,
							rate,
						),
					)
				}
				for _, sell := range respData.Asks {
					quantity, _ := strconv.ParseFloat(sell["amount"], 64)
					rate, _ := strconv.ParseFloat(sell["price"], 64)
					result.Asks = append(
						result.Asks,
						common.NewPriceEntry(
							quantity,
							rate,
						),
					)
				}
			}
		}
	}
	data.Store(pair.PairID(), result)
}

func (self *BitfinexEndpoint) Trade(tradeType string, base, quote common.Token, rate, amount float64) (done float64, remaining float64, finished bool, err error) {
	result := exchange.Bitftrade{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	data := url.Values{}
	data.Set("method", "Trade")
	data.Set("pair", fmt.Sprintf("%s_%s", strings.ToLower(base.ID), strings.ToLower(quote.ID)))
	data.Set("type", tradeType)
	data.Set("rate", strconv.FormatFloat(rate, 'f', -1, 64))
	data.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	params := data.Encode()
	req, _ := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(),
		bytes.NewBufferString(params),
	)
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetKey())
	req.Header.Add("Sign", self.signer.Sign(params))
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Response body close error: %s", err.Error())
			}
		}()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			err = json.Unmarshal(respBody, &result)
		}
	} else {
		log.Printf("Error: %v, Code: %v\n", err, resp)
	}
	return
}

func (self *BitfinexEndpoint) Withdraw(token common.Token, amount *big.Int, address ethereum.Address) error {
	// ignoring timepoint because it's only relevant in simulation
	result := exchange.Bitfwithdraw{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	data := url.Values{}
	data.Set("method", "WithdrawCoin")
	data.Set("coinName", token.ID)
	data.Set("amount", strconv.FormatFloat(common.BigToFloat(amount, token.Decimal), 'f', -1, 64))
	data.Set("address", address.Hex())
	params := data.Encode()
	req, _ := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(),
		bytes.NewBufferString(params),
	)
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetKey())
	req.Header.Add("Sign", self.signer.Sign(params))
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Response body close error: %s", err.Error())
			}
		}()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			err = json.Unmarshal(respBody, &result)
		}
		if err != nil {
			return err
		}
		if result.Error != "" {
			return errors.New(result.Error)
		}
		return nil
	}
	return errors.New("withdraw rejected by Bitfinex")
}

func (self *BitfinexEndpoint) GetInfo() (exchange.Bitfinfo, error) {
	result := exchange.Bitfinfo{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	data := url.Values{}
	data.Set("method", "getInfo")
	data.Add("nonce", nonce())
	params := data.Encode()
	log.Printf("endpoint: %v\n", self.interf.AuthenticatedEndpoint())
	req, _ := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(),
		bytes.NewBufferString(params),
	)
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetKey())
	req.Header.Add("Sign", self.signer.Sign(params))
	resp, err := client.Do(req)
	if err == nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Response body close error: %s", err.Error())
			}
		}()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			err = json.Unmarshal(respBody, &result)
		}
	}
	return result, err
}

//NewBitfinexEndpoint return bitfinex endpoint instance
func NewBitfinexEndpoint(signer Signer, interf Interface) *BitfinexEndpoint {
	return &BitfinexEndpoint{signer, interf}
}

//NewRealBitfinexEndpoint return real endpoint instance
func NewRealBitfinexEndpoint(signer Signer) *BitfinexEndpoint {
	return &BitfinexEndpoint{signer, NewRealInterface()}
}

//NewSimulatedBitfinexEndpoint return simulated endpoint
func NewSimulatedBitfinexEndpoint(signer Signer) *BitfinexEndpoint {
	return &BitfinexEndpoint{signer, NewSimulatedInterface()}
}
