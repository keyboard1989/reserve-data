package bitfinex

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
	"github.com/KyberNetwork/reserve-data/exchange"
	ethereum "github.com/ethereum/go-ethereum/common"
)

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

func (self *BitfinexEndpoint) fillRequest(req *http.Request, signNeeded bool) {
	if req.Method == "POST" || req.Method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Add("Accept", "application/json")
	if signNeeded == true {
		q := req.URL.Query()
		q.Set("apiKey", self.signer.GetBitfinexKey())
		q.Set("nonce", nonce())
		req.URL.RawQuery = q.Encode()
		req.Header.Add("apisign", self.signer.BitfinexSign(req.URL.String()))
	}
}

func (self *BitfinexEndpoint) Trade(tradeType string, base, quote common.Token, rate, amount float64, timepoint uint64) (done float64, remaining float64, finished bool, err error) {
	result := exchange.Bitftrade{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	data := url.Values{}
	data.Set("method", "Trade")
	data.Set("pair", fmt.Sprintf("%s_%s", strings.ToLower(base.ID), strings.ToLower(quote.ID)))
	data.Set("type", tradeType)
	data.Set("rate", fmt.Sprintf("%f", rate))
	data.Set("amount", fmt.Sprintf("%f", amount))
	params := data.Encode()
	req, _ := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(),
		bytes.NewBufferString(params),
	)
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetBitfinexKey())
	req.Header.Add("Sign", self.signer.BitfinexSign(params))
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("response: %s\n", resp_body)
		if err == nil {
			err = json.Unmarshal(resp_body, &result)
		}
	} else {
		fmt.Printf("Error: %v, Code: %v\n", err, resp)
	}
	return
}

func (self *BitfinexEndpoint) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) error {
	// ignoring timepoint because it's only relevant in simulation
	result := exchange.Bitfwithdraw{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	data := url.Values{}
	data.Set("method", "WithdrawCoin")
	data.Set("coinName", token.ID)
	data.Set("amount", fmt.Sprintf("%f", common.BigToFloat(amount, token.Decimal)))
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
	req.Header.Add("Key", self.signer.GetBitfinexKey())
	req.Header.Add("Sign", self.signer.BitfinexSign(params))
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
		return errors.New("withdraw rejected by Bitfinex")
	}
}

func (self *BitfinexEndpoint) GetInfo(timepoint uint64) (exchange.Bitfinfo, error) {
	result := exchange.Bitfinfo{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	data := url.Values{}
	data.Set("method", "getInfo")
	data.Add("nonce", nonce())
	params := data.Encode()
	fmt.Printf("endpoint: %v\n", self.interf.AuthenticatedEndpoint())
	req, _ := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(),
		bytes.NewBufferString(params),
	)
	fmt.Printf("params: %v\n", params)
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetBitfinexKey())
	req.Header.Add("Sign", self.signer.BitfinexSign(params))
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

func NewBitfinexEndpoint(signer Signer, interf Interface) *BitfinexEndpoint {
	return &BitfinexEndpoint{signer, interf}
}

func NewRealBitfinexEndpoint(signer Signer) *BitfinexEndpoint {
	return &BitfinexEndpoint{signer, NewRealInterface()}
}

func NewSimulatedBitfinexEndpoint(signer Signer) *BitfinexEndpoint {
	return &BitfinexEndpoint{signer, NewSimulatedInterface()}
}
