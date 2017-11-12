package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	"errors"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/exchange"
	ethereum "github.com/ethereum/go-ethereum/common"
	"strconv"
	"sync"
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

func (self *BinanceEndpoint) FetchOnePairData(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map,
	timepoint uint64) {

	defer wg.Done()
	result := common.ExchangePrice{}

	client := &http.Client{}
	req, _ := http.NewRequest(
		"GET",
		self.interf.PublicEndpoint()+"/api/v1/depth",
		nil)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("symbol", fmt.Sprintf("%s%s", pair.Base.ID, pair.Quote.ID))
	q.Add("limit", "50")
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, false, timepoint)

	timestamp := common.Timestamp(fmt.Sprintf("%d", timepoint))
	resp, err := client.Do(req)
	result.Timestamp = timestamp
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
			resp_data := exchange.Binaresp{}
			json.Unmarshal(resp_body, &resp_data)
			if resp_data.Code != 0 || resp_data.Msg != "" {
				result.Valid = false
			} else {
				for _, buy := range resp_data.Bids {
					quantity, _ := strconv.ParseFloat(buy[1], 64)
					rate, _ := strconv.ParseFloat(buy[0], 64)
					result.Bids = append(
						result.Bids,
						common.PriceEntry{
							quantity,
							rate,
						},
					)
				}
				for _, sell := range resp_data.Asks {
					quantity, _ := strconv.ParseFloat(sell[1], 64)
					rate, _ := strconv.ParseFloat(sell[0], 64)
					result.Asks = append(
						result.Asks,
						common.PriceEntry{
							quantity,
							rate,
						},
					)
				}
			}
		}
	}
	data.Store(pair.PairID(), result)
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
		log.Printf("response: %s\n", resp_body)
		if err == nil {
			err = json.Unmarshal(resp_body, &result)
		}
	} else {
		log.Printf("Error: %v, Code: %v\n", err, resp)
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
		log.Printf("response: %s\n", resp_body)
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
		log.Printf("Error: %v, Code: %v\n", err, resp)
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
