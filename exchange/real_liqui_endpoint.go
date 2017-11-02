package exchange

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
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type RealLiquiEndpoint struct {
	AuthenticatedEndpoint string
	PublicEndpoint        string
}

func nonce() string {
	epsilon := 30 * time.Millisecond
	anchor := int64(50299954901)
	timestamp := time.Now().UnixNano()/int64(epsilon) - anchor
	return strconv.Itoa(int(timestamp))
}

func (self *RealLiquiEndpoint) Depth(tokens string) (liqresp, error) {
	result := liqresp{}
	client := &http.Client{}
	url := fmt.Sprintf(
		"%s/depth/%s?ignore_invalid=1",
		self.PublicEndpoint,
		tokens,
	)
	req, _ := http.NewRequest("GET", url, nil)
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

type liqwithdraw struct {
	Success int `json:"success"`
	Return  map[string]interface{}
	Error   string `json:"error"`
}

func (self *RealLiquiEndpoint) Withdraw(key string, token common.Token, amount *big.Int, address ethereum.Address, signer Signer) error {
	result := liqwithdraw{}
	client := &http.Client{}
	data := url.Values{}
	data.Set("method", "WithdrawCoin")
	data.Set("coinName", token.ID)
	data.Set("amount", fmt.Sprintf("%f", common.BigToFloat(amount, token.Decimal)))
	data.Set("address", address.Hex())
	params := data.Encode()
	req, _ := http.NewRequest(
		"POST",
		self.AuthenticatedEndpoint,
		bytes.NewBufferString(params),
	)
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", key)
	req.Header.Add("Sign", signer.LiquiSign(params))
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
		return errors.New("withdraw rejected by Liqui")
	}
}

func (self *RealLiquiEndpoint) GetInfo(key string, signer Signer) (liqinfo, error) {
	result := liqinfo{}
	client := &http.Client{}
	data := url.Values{}
	data.Set("method", "getInfo")
	data.Add("nonce", nonce())
	params := data.Encode()
	req, _ := http.NewRequest(
		"POST",
		self.AuthenticatedEndpoint,
		bytes.NewBufferString(params),
	)
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", key)
	req.Header.Add("Sign", signer.LiquiSign(params))
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

func NewRealLiquiEndpoint() *RealLiquiEndpoint {
	return &RealLiquiEndpoint{
		"https://api.liqui.io/tapi",
		"https://api.liqui.io/api/3",
	}
}
