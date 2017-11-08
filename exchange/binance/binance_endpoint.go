package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/exchange"
	ethereum "github.com/ethereum/go-ethereum/common"
	"errors"
)

type BinanceEndpoint struct {
	signer Signer
	interf Interface
}

func (self *BinanceEndpoint) fillRequest(req *http.Request, signNeeded bool, timepoint uint64) {
	if req.Method == "POST" || req.Method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Add("Accept", "application/json")
	if signNeeded {
		req.Header.Add("X-MBX-APIKEY", self.signer.GetBinanceKey())
		q := req.URL.Query()
		q.Set("timestamp", fmt.Sprintf("%d", timepoint))
		q.Set("recvWindow", "5000")
		q.Set("signature", self.signer.BinanceSign(q.Encode()))
		req.URL.RawQuery = q.Encode()
	}
}

func (self *BinanceEndpoint) Depth(tokens string, timepoint uint64) (exchange.Binaresp, error) {
	result := exchange.Binaresp{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	req, _ := http.NewRequest(
		"GET",
		self.interf.PublicEndpoint()+"/api/v1/depth",
		nil)
	self.fillRequest(req, false, timepoint)
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
	result := exchange.Binatrade{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, _ := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint()+"/api/v3/order/test",
		nil,
	)
	self.fillRequest(req, true, timepoint)
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

func (self *BinanceEndpoint) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) error {
	result := exchange.Binawithdraw{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, _ := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint()+"/wapi/v1/withdraw.html",
		nil,
	)
	self.fillRequest(req, true, timepoint)
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
		if result.Success == false {
			return errors.New(result.Message)
		}
		return nil
	} else {
		fmt.Printf("Error: %v, Code: %v\n", err, resp)
		return errors.New("withdraw rejected by Binnace")
	}
}

func (self *BinanceEndpoint) GetInfo(timepoint uint64) (exchange.Binainfo, error) {
	result := exchange.Binainfo{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	req, _ := http.NewRequest(
		"GET",
		self.interf.AuthenticatedEndpoint()+"/api/account",
		nil)
	self.fillRequest(req, true, timepoint)
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
