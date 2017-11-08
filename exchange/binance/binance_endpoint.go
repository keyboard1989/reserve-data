package binance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/exchange"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type BinanceEndpoint struct {
	signer Signer
	interf Interface
}

func nonce() string {
	epsilon := 30 * time.Millisecond
	anchor := int64(50299954901)
	timestamp := time.Now().UnixNano()/int64(epsilon) - anchor
	return strconv.Itoa(int(timestamp))
}

func (self *BinanceEndpoint) Depth(tokens string, timepoint uint64) (exchange.Binaresp, error) {
	result := exchange.Binaresp{}
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

func (self *BinanceEndpoint) Trade(tradeType string, base, quote common.Token, rate, amount float64, timepoint uint64) (done float64, remaining float64, finished bool, err error) {
	return
}

func (self *BinanceEndpoint) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) error {
	return nil
}

func (self *BinanceEndpoint) GetInfo(timepoint uint64) (exchange.Binainfo, error) {
	result := exchange.Binainfo{}
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
	req.Header.Add("Key", self.signer.GetBinanceKey())
	req.Header.Add("Sign", self.signer.BinanceSign(params))
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

func NewBinanceEndpoint(signer Signer, interf Interface) *BinanceEndpoint {
	return &BinanceEndpoint{signer, interf}
}

func NewRealBinanceEndpoint(signer Signer) *BinanceEndpoint {
	return &BinanceEndpoint{signer, NewRealInterface()}
}

func NewSimulatedBinanceEndpoint(signer Signer) *BinanceEndpoint {
	return &BinanceEndpoint{signer, NewSimulatedInterface()}
}
