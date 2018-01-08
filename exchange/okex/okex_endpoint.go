package okex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/exchange"
)

type OkexEndpoint struct {
	signer Signer
	interf Interface
}

func (self *OkexEndpoint) fillRequest(req *http.Request, signNeeded bool, timepoint uint64) {
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Add("Acceept", "application/json")
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
