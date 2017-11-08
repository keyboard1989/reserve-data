package bittrex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/exchange"
	ethereum "github.com/ethereum/go-ethereum/common"
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

func (self *BittrexEndpoint) Depth(tokens string, timepoint uint64) (exchange.Liqresp, error) {
	result := exchange.Liqresp{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	u, err := url.Parse(self.interf.PublicEndpoint(timepoint))
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Set("ignore_invalid", "1")
	u.RawQuery = q.Encode()
	u.Path = path.Join(
		u.Path,
		"depth",
		tokens,
	)
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Accept", "application/json")
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

func (self *BittrexEndpoint) Trade(tradeType string, base, quote common.Token, rate, amount float64, timepoint uint64) (done float64, remaining float64, finished bool, err error) {
	result := exchange.Liqtrade{}
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
		self.interf.AuthenticatedEndpoint(timepoint),
		bytes.NewBufferString(params),
	)
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetLiquiKey())
	req.Header.Add("Sign", self.signer.LiquiSign(params))
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
		return 0, 0, false, errors.New("Trade rejected by Liqui")
	}
}

func (self *BittrexEndpoint) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) error {
	// ignoring timepoint because it's only relevant in simulation
	result := exchange.Liqwithdraw{}
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
		self.interf.AuthenticatedEndpoint(timepoint),
		bytes.NewBufferString(params),
	)
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetLiquiKey())
	req.Header.Add("Sign", self.signer.LiquiSign(params))
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
		return errors.New("withdraw rejected by Liqui")
	}
}

func (self *BittrexEndpoint) GetInfo(timepoint uint64) (exchange.Liqinfo, error) {
	result := exchange.Liqinfo{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	data := url.Values{}
	data.Set("method", "getInfo")
	data.Add("nonce", nonce())
	params := data.Encode()
	fmt.Printf("endpoint: %v\n", self.interf.AuthenticatedEndpoint(timepoint))
	req, _ := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(timepoint),
		bytes.NewBufferString(params),
	)
	fmt.Printf("params: %v\n", params)
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetLiquiKey())
	req.Header.Add("Sign", self.signer.LiquiSign(params))
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

func NewLiquiEndpoint(signer Signer, interf Interface) *BittrexEndpoint {
	return &BittrexEndpoint{signer, interf}
}

func NewRealLiquiEndpoint(signer Signer) *BittrexEndpoint {
	return &BittrexEndpoint{signer, NewRealInterface()}
}

func NewSimulatedLiquiEndpoint(signer Signer) *BittrexEndpoint {
	return &BittrexEndpoint{signer, NewSimulatedInterface()}
}
