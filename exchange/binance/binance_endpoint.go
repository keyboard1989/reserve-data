package binance

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
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/exchange"
	ethereum "github.com/ethereum/go-ethereum/common"
)

const EPSILON float64 = 0.0000000001 // 10e-10

type BinanceEndpoint struct {
	signer Signer
	interf Interface
}

func (self *BinanceEndpoint) fillRequest(req *http.Request, signNeeded bool, timepoint uint64) {
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	if signNeeded {
		req.Header.Set("X-MBX-APIKEY", self.signer.GetBinanceKey())
		q.Set("timestamp", fmt.Sprintf("%d", timepoint))
		q.Set("recvWindow", "5000")
		q.Set("signature", self.signer.BinanceSign(q.Encode()))
	}
	req.URL.RawQuery = q.Encode()
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
				result.Error = fmt.Sprintf("Code: %d, Msg: %s", resp_data.Code, resp_data.Msg)
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

// Relevant params:
// symbol ("%s%s", base, quote)
// side (BUY/SELL)
// type (LIMIT/MARKET)
// timeInForce (GTC/IOC)
// quantity
// price
//
// In this version, we only support LIMIT order which means only buy/sell with acceptable price,
// and GTC time in force which means that the order will be active until it's implicitly canceled
func (self *BinanceEndpoint) Trade(tradeType string, base, quote common.Token, rate, amount float64, timepoint uint64) (id string, done float64, remaining float64, finished bool, err error) {
	result := exchange.Binatrade{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, _ := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint()+"/api/v3/order",
		nil,
	)
	q := req.URL.Query()
	q.Add("symbol", base.ID+quote.ID)
	q.Add("side", strings.ToUpper(tradeType))
	q.Add("type", "LIMIT")
	q.Add("timeInForce", "GTC")
	q.Add("quantity", strconv.FormatFloat(amount, 'f', -1, 64))
	q.Add("price", strconv.FormatFloat(rate, 'f', -1, 64))
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, true, timepoint)
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		log.Printf("response: %s\n", resp_body)
		if err == nil {
			err = json.Unmarshal(resp_body, &result)
		}
	} else {
		log.Printf("Error: %v, Code: %v\n", err, resp)
	}
	done, remaining, finished, err = self.QueryOrder(
		base.ID+quote.ID,
		result.OrderID,
		timepoint+20,
	)
	return result.ClientOrderID, done, remaining, finished, err
}

func (self *BinanceEndpoint) CancelOrder(base, quote common.Token, id uint64) (exchange.Binacancel, error) {
	result := exchange.Binacancel{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, _ := http.NewRequest(
		"DELETE",
		self.interf.AuthenticatedEndpoint()+"/api/v3/order",
		nil,
	)
	q := req.URL.Query()
	q.Add("symbol", base.ID+quote.ID)
	q.Add("orderId", fmt.Sprintf("%d", id))
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, true, common.GetTimepoint())
	var resp_body []byte
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		resp_body, err = ioutil.ReadAll(resp.Body)
		log.Printf("response: %s\n", resp_body)
		if err == nil {
			err = json.Unmarshal(resp_body, &result)
		}
	}
	return result, err
}

func (self *BinanceEndpoint) QueryOrder(symbol string, id uint64, timepoint uint64) (done float64, remaining float64, finished bool, err error) {
	result := exchange.Binaorder{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, _ := http.NewRequest(
		"GET",
		self.interf.AuthenticatedEndpoint()+"/api/v3/order",
		nil,
	)
	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("orderId", fmt.Sprintf("%d", id))
	req.URL.RawQuery = q.Encode()
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
			return 0, 0, false, err
		}
		if result.Code != 0 {
			return 0, 0, false, errors.New(result.Message)
		}
		done, _ := strconv.ParseFloat(result.ExecutedQty, 64)
		total, _ := strconv.ParseFloat(result.OrigQty, 64)
		return done, total - done, total-done < EPSILON, nil
	} else {
		log.Printf("Error: %v, Code: %v\n", err, resp)
		return 0, 0, false, errors.New("withdraw rejected by Binnace")
	}
}

func (self *BinanceEndpoint) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) (ethereum.Hash, error) {
	result := exchange.Binawithdraw{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, _ := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint()+"/wapi/v3/withdraw.html",
		nil,
	)
	q := req.URL.Query()
	q.Add("asset", token.ID)
	q.Add("address", address.Hex())
	q.Add("amount", strconv.FormatFloat(common.BigToFloat(amount, token.Decimal), 'f', -1, 64))
	req.URL.RawQuery = q.Encode()
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
			return ethereum.Hash{}, err
		}
		if result.Success == false {
			return ethereum.Hash{}, errors.New(result.Message)
		}
		return ethereum.HexToHash("0x" + result.TxHash), nil
	} else {
		log.Printf("Error: %v, Code: %v\n", err, resp)
		return ethereum.Hash{}, errors.New("withdraw rejected by Binnace")
	}
}

func (self *BinanceEndpoint) GetInfo(timepoint uint64) (exchange.Binainfo, error) {
	result := exchange.Binainfo{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	req, _ := http.NewRequest(
		"GET",
		self.interf.AuthenticatedEndpoint()+"/api/v3/account",
		nil)
	self.fillRequest(req, true, timepoint)
	resp, err := client.Do(req)
	if err == nil {
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			resp_body, err := ioutil.ReadAll(resp.Body)
			log.Printf("Binance get balances: %s", string(resp_body))
			if err == nil {
				json.Unmarshal(resp_body, &result)
			}
		} else {
			err = errors.New("Unsuccessful response from Binance: Status " + resp.Status)
		}
	}
	return result, err
}

func (self *BinanceEndpoint) OpenOrdersForOnePair(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map,
	timepoint uint64) {

	defer wg.Done()
	result := exchange.Binaorders{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	req, _ := http.NewRequest(
		"GET",
		self.interf.AuthenticatedEndpoint()+"/api/v3/openOrders",
		nil)
	q := req.URL.Query()
	q.Add("symbol", pair.Base.ID+pair.Quote.ID)
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, true, timepoint)
	resp, err := client.Do(req)
	if err == nil {
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			resp_body, err := ioutil.ReadAll(resp.Body)
			log.Printf("Binance get open orders for %s: %s", pair.PairID(), string(resp_body))
			if err == nil {
				json.Unmarshal(resp_body, &result)
				orders := []common.Order{}
				for _, order := range result {
					price, _ := strconv.ParseFloat(order.Price, 64)
					orgQty, _ := strconv.ParseFloat(order.OrigQty, 64)
					executedQty, _ := strconv.ParseFloat(order.ExecutedQty, 64)
					orders = append(orders, common.Order{
						Base:        strings.ToUpper(pair.Base.ID),
						Quote:       strings.ToUpper(pair.Quote.ID),
						OrderId:     fmt.Sprintf("%d", order.OrderId),
						Price:       price,
						OrigQty:     orgQty,
						ExecutedQty: executedQty,
						TimeInForce: order.TimeInForce,
						Type:        order.Type,
						Side:        order.Side,
						StopPrice:   order.StopPrice,
						IcebergQty:  order.IcebergQty,
						Time:        order.Time,
					})
				}
				data.Store(pair.PairID(), orders)
			}
		} else {
			err = errors.New("Unsuccessful response from Binance: Status " + resp.Status)
		}
	}
}

func (self *BinanceEndpoint) GetDepositAddress(asset string) (exchange.Binadepositaddress, error) {
	result := exchange.Binadepositaddress{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	req, _ := http.NewRequest(
		"GET",
		self.interf.AuthenticatedEndpoint()+"/wapi/v3/depositAddress.html",
		nil)
	timepoint := common.GetTimepoint()
	q := req.URL.Query()
	q.Add("asset", asset)
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, true, timepoint)
	log.Printf("Request to binance: %s", req.URL)
	log.Printf("Header for Request to binance: %s", req.Header)
	resp, err := client.Do(req)
	var resp_body []byte
	if err == nil {
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			resp_body, err = ioutil.ReadAll(resp.Body)
			log.Printf("Binance get balances: %s", string(resp_body))
			if err == nil {
				err = json.Unmarshal(resp_body, &result)
				if !result.Success {
					err = errors.New(result.Msg)
				}
			}
		} else {
			err = errors.New("Unsuccessful response from Binance: Status " + resp.Status)
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

func NewKovanBinanceEndpoint(signer Signer) *BinanceEndpoint {
	return &BinanceEndpoint{signer, NewKovanInterface()}
}

func NewDevBinanceEndpoint(signer Signer) *BinanceEndpoint {
	return &BinanceEndpoint{signer, NewDevInterface()}
}
