package okex

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
)

type OkexEndpoint struct {
	signer Signer
	interf Interface
}

func (self *OkexEndpoint) fillRequest(req *http.Request, signNeeded bool, timepoint uint64) {
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Add("Accept", "application/json")
}

func (self *OkexEndpoint) GetResponse(
	method string, url string,
	params map[string]string, signNeeded bool, timepoint uint64) ([]byte, error) {

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	self.fillRequest(req, signNeeded, timepoint)
	var err error
	var resp_body []byte
	log.Printf("request to Okex: %s\n", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return resp_body, err
	} else {
		defer resp.Body.Close()
		resp_body, err = ioutil.ReadAll(resp.Body)
		log.Printf("request to %s, got response from okex: %s\n", req.URL, common.TruncStr(resp_body))
		return resp_body, err
	}
}

func (self *OkexEndpoint) GetDepthOnePair(
	pair common.TokenPair, timepoint uint64) (exchange.OkexDepth, error) {

	resp_body, err := self.GetResponse(
		"GET", self.interf.PublicEnpoint()+"/api/v1/depth.do",
		map[string]string{
			"symbol": fmt.Sprintf("%s_%s", strings.ToLower(pair.Base.ID), strings.ToLower(pair.Quote.ID)),
			"size":   "50",
		},
		false,
		timepoint,
	)
	resp_data := exchange.OkexDepth{}
	if err != nil {
		return resp_data, err
	} else {
		json.Unmarshal(resp_body, &resp_data)
		return resp_data, nil
	}
}

func (self *OkexEndpoint) GetExchangeInfo() (exchange.OkexInfo, error) {
	result := exchange.OkexInfo{}
	timepoint := common.GetTimepoint()
	resp_body, err := self.GetResponse(
		"GET",
		"https://www.okex.com/v2/markets/products",
		map[string]string{},
		false,
		timepoint,
	)
	if err != nil {
		return result, err
	} else {
		json.Unmarshal(resp_body, &result)
		return result, nil
	}
}

func (self *OkexEndpoint) Trade(base, quote common.Token, rate, amount float64, timepoint uint64) (exchange.OkexTrade, error) {
	result := exchange.OkexTrade{}
	symbol := fmt.Sprintf("%s_%s", base.ID, quote.ID)
	orderType := "limit"
	params := map[string]string{
		"amount":  strconv.FormatFloat(amount, 'f', -1, 64),
		"api_key": "", //TODO:
		"symbol":  symbol,
		"type":    orderType,
	}
	if orderType == "limit" {
		params["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	}
	resp_body, err := self.GetResponse(
		"POST",
		self.interf.AuthenticatedEndpoint()+"/api/v1/trade.do",
		params,
		true,
		timepoint,
	)
	if err != nil {
		return result, err
	} else {
		json.Unmarshal(resp_body, &result)
		return result, nil
	}
}

func (self *OkexEndpoint) CancelOrder(symbol string, id uint64, timepoint uint64) (exchange.OkexCancel, error) {
	result := exchange.OkexCancel{}
	resp_body, err := self.GetResponse(
		"GET",
		self.interf.AuthenticatedEndpoint()+"/api/v1/cancel_order.do",
		map[string]string{
			"api_key":  "", //TODO:
			"symbol":   symbol,
			"order_id": fmt.Sprintf("%d", id),
		},
		true,
		timepoint,
	)
	if err != nil {
		return result, err
	} else {
		json.Unmarshal(resp_body, &result)
		if result.Result == false {
			return result, errors.New(fmt.Sprintf("Order %d cancel failed", id))
		} else {
			return result, nil
		}
	}
}

func (self *OkexEndpoint) OrderStatus(symbol string, id uint64, timepoint uint64) (exchange.OkexOrderStatus, error) {
	result := exchange.OkexOrderStatus{}
	resp_body, err := self.GetResponse(
		"POST",
		self.interf.AuthenticatedEndpoint()+"/api/v1/order_info.do",
		map[string]string{
			"api_key":  "", //TODO:
			"symbol":   symbol,
			"order_id": fmt.Sprintf("%d", id),
		},
		true,
		timepoint,
	)
	if err != nil {
		return result, err
	} else {
		json.Unmarshal(resp_body, &result)
		if result.Result == false {
			return result, errors.New(fmt.Sprintf("Get order status failed: %d\n", id))
		} else {
			return result, nil
		}
	}
}

func (self *OkexEndpoint) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) (string, error) {
	result := exchange.Binawithdraw{}
	resp_body, err := self.GetResponse(
		"POST",
		self.interf.AuthenticatedEndpoint()+"/wapi/v3/withdraw.html",
		map[string]string{
			"asset":   token.ID,
			"address": address.Hex(),
			"name":    "reserve",
			"amount":  strconv.FormatFloat(common.BigToFloat(amount, token.Decimal), 'f', -1, 64),
		},
		true,
		timepoint,
	)
	if err == nil {
		json.Unmarshal(resp_body, &result)
		if result.Success == false {
			return "", errors.New(result.Message)
		}
		return result.ID, nil
	} else {
		log.Printf("Error: %v", err)
		return "", errors.New("withdraw rejected by Binnace")
	}
}

func (self *OkexEndpoint) GetInfo(timepoint uint64) (exchange.OkexAccountInfo, error) {
	result := exchange.OkexAccountInfo{}
	resp_body, err := self.GetResponse(
		"POST",
		self.interf.AuthenticatedEndpoint()+"/api/v1/userinfo.do",
		map[string]string{
			"api_key": "", // TODO:
		},
		true,
		timepoint,
	)
	if err == nil {
		json.Unmarshal(resp_body, &result)
	}
	return result, err
}

func NewOkexEndpoint(signer Signer, interf Interface) *OkexEndpoint {
	return &OkexEndpoint{signer, interf}
}

func NewOkexRealEndpoint(signer Signer) *OkexEndpoint {
	return &OkexEndpoint{signer, NewRealInterface()}
}

func NewSimulatedOkexEndpoint(signer Signer) *OkexEndpoint {
	return &OkexEndpoint{signer, NewSimulatedInterface()}
}

func NewDevOkexEndpoint(signer Signer) *OkexEndpoint {
	return &OkexEndpoint{signer, NewDevInterface()}
}
