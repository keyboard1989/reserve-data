package liqui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

// LiquiEndpoint endpoint for liqui
// including signer for api authentication
// interf for different env interfacw
type LiquiEndpoint struct {
	signer Signer
	interf Interface
}

func nonce() string {
	epsilon := 30 * time.Millisecond
	anchor := int64(50299954901)
	timestamp := time.Now().UnixNano()/int64(epsilon) - anchor
	return strconv.Itoa(int(timestamp))
}

func (self *LiquiEndpoint) Depth(tokens string, timepoint uint64) (exchange.Liqresp, error) {
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
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return result, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err == nil {
		if resp.StatusCode == 200 {
			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Printf("Response body close error: %s", err.Error())
				}
			}()
			respBody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				if err := json.Unmarshal(respBody, &result); err != nil {
					log.Printf("Unmarshal response error: %s", err.Error())
				}
			}
		} else {
			err = errors.New("Unsuccessful response from Liqui: Status " + resp.Status)
		}
	}
	return result, err
}

func (self *LiquiEndpoint) CancelOrder(id string) (exchange.Liqcancel, error) {
	result := exchange.Liqcancel{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	data := url.Values{}
	data.Set("method", "CancelOrder")
	data.Set("order_id", id)
	data.Add("nonce", nonce())
	params := data.Encode()
	req, err := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(common.GetTimepoint()),
		bytes.NewBufferString(params),
	)
	if err != nil {
		return result, err
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetKey())
	req.Header.Add("Sign", self.signer.Sign(params))
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Fatalf("Response body close error: %s", err.Error())
			}
		}()
		respBody, err := ioutil.ReadAll(resp.Body)
		log.Printf("response: %s\n", respBody)
		if err == nil {
			err = json.Unmarshal(respBody, &result)
		}
		return result, err
	}
	return result, errors.New("Cancel rejected by Liqui")
}

func (self *LiquiEndpoint) Trade(tradeType string, base, quote common.Token, rate, amount float64, timepoint uint64) (id string, done float64, remaining float64, finished bool, err error) {
	result := exchange.Liqtrade{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	data := url.Values{}
	data.Set("method", "Trade")
	data.Set("pair", fmt.Sprintf("%s_%s", strings.ToLower(base.ID), strings.ToLower(quote.ID)))
	data.Set("type", tradeType)
	data.Set("rate", strconv.FormatFloat(rate, 'f', -1, 64))
	data.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	data.Add("nonce", nonce())
	params := data.Encode()
	req, err := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(timepoint),
		bytes.NewBufferString(params),
	)
	if err != nil {
		return "", 0, 0, false, err
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetKey())
	req.Header.Add("Sign", self.signer.Sign(params))
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Fatalf("Response body close error: %s", err.Error())
			}
		}()
		respBody, err := ioutil.ReadAll(resp.Body)
		log.Printf("response: %s\n", respBody)
		if err == nil {
			err = json.Unmarshal(respBody, &result)
		}
		if err != nil {
			return "", 0, 0, false, err
		}
		if result.Error != "" {
			return "", 0, 0, false, errors.New(result.Error)
		}
		return strconv.FormatUint(result.Return.OrderID, 10), result.Return.Done, result.Return.Remaining, result.Return.OrderID == 0, nil
	}
	return "", 0, 0, false, errors.New("Trade rejected by Liqui")

}

func (self *LiquiEndpoint) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) error {
	// ignoring timepoint because it's only relevant in simulation
	result := exchange.Liqwithdraw{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	data := url.Values{}
	data.Set("method", "WithdrawCoin")
	data.Set("coinName", token.ID)
	data.Set("amount", strconv.FormatFloat(common.BigToFloat(amount, token.Decimal), 'f', -1, 64))
	data.Set("address", address.Hex())
	data.Add("nonce", nonce())
	params := data.Encode()
	req, err := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(timepoint),
		bytes.NewBufferString(params),
	)
	if err != nil {
		return err
	}
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
		log.Printf("response: %s\n", respBody)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(respBody, &result); err != nil {
			return err
		}
		if result.Error != "" {
			return errors.New(result.Error)
		}
		return nil
	}
	return errors.New("withdraw rejected by Liqui")
}

func (self *LiquiEndpoint) GetInfo(timepoint uint64) (exchange.Liqinfo, error) {
	result := exchange.Liqinfo{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	data := url.Values{}
	data.Set("method", "getInfo")
	data.Add("nonce", nonce())
	params := data.Encode()
	log.Printf("endpoint: %v\n", self.interf.AuthenticatedEndpoint(timepoint))
	req, err := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(timepoint),
		bytes.NewBufferString(params),
	)
	if err != nil {
		return result, err
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetKey())
	req.Header.Add("Sign", self.signer.Sign(params))
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if resp.StatusCode == 200 {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Response body close error: %s", err.Error())
			}
		}()
		respBody, err := ioutil.ReadAll(resp.Body)
		log.Printf("Liqui GetInfo response: %s", string(respBody))
		if err == nil {
			if err := json.Unmarshal(respBody, &result); err != nil {
				return result, err
			}
		}
		log.Printf("Liqui GetInfo data: %v", result)
	} else {
		err = errors.New("Unsuccessful response from Liqui: Status " + resp.Status)
	}
	return result, err
}

func (self *LiquiEndpoint) OrderInfo(orderID string, timepoint uint64) (exchange.Liqorderinfo, error) {
	result := exchange.Liqorderinfo{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	data := url.Values{}
	data.Set("method", "OrderInfo")
	data.Set("order_id", orderID)
	data.Add("nonce", nonce())
	params := data.Encode()
	req, err := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(timepoint),
		bytes.NewBufferString(params),
	)
	if err != nil {
		return result, err
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetKey())
	req.Header.Add("Sign", self.signer.Sign(params))
	resp, err := client.Do(req)
	if err == nil {
		if resp.StatusCode == 200 {
			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Printf("Response body close error: %s", err.Error())
				}
			}()
			respBody, err := ioutil.ReadAll(resp.Body)
			log.Printf("Liqui Order info response: %s", string(respBody))
			if err == nil {
				if err := json.Unmarshal(respBody, &result); err != nil {
					return result, err
				}
			}
			log.Printf("Liqui Order info data: %v", result)
		} else {
			err = errors.New("Unsuccessful response from Liqui: Status " + resp.Status)
		}
	}
	return result, err
}

func (self *LiquiEndpoint) ActiveOrders(timepoint uint64) (exchange.Liqorders, error) {
	result := exchange.Liqorders{}
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second)}
	data := url.Values{}
	data.Set("method", "ActiveOrders")
	data.Set("pair", "") // all pairs
	data.Add("nonce", nonce())
	params := data.Encode()
	req, err := http.NewRequest(
		"POST",
		self.interf.AuthenticatedEndpoint(timepoint),
		bytes.NewBufferString(params),
	)
	if err != nil {
		return result, err
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", self.signer.GetKey())
	req.Header.Add("Sign", self.signer.Sign(params))
	resp, err := client.Do(req)
	if err == nil {
		if resp.StatusCode == 200 {
			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Printf("Response body close error: %s", err.Error())
				}
			}()
			respBody, err := ioutil.ReadAll(resp.Body)
			log.Printf("Liqui ActiveOrders response: %s", string(respBody))
			if err == nil {
				if err := json.Unmarshal(respBody, &result); err != nil {
					return result, err
				}
			}
			log.Printf("Liqui ActiveOrders data: %v", result)
		} else {
			err = errors.New("Unsuccessful response from Liqui: Status " + resp.Status)
		}
	}
	return result, err
}

//NewLiquiEndpoint return new endpoint instance
func NewLiquiEndpoint(signer Signer, interf Interface) *LiquiEndpoint {
	return &LiquiEndpoint{signer, interf}
}

//NewRealLiquiEndpoint return real endpoint instance
func NewRealLiquiEndpoint(signer Signer) *LiquiEndpoint {
	return &LiquiEndpoint{signer, NewRealInterface()}
}

//NewSimulatedLiquiEndpoint return simulated endpoint instance
func NewSimulatedLiquiEndpoint(signer Signer) *LiquiEndpoint {
	return &LiquiEndpoint{signer, NewSimulatedInterface()}
}

//NewKovanLiquiEndpoint return kovan endpoint instance
func NewKovanLiquiEndpoint(signer Signer) *LiquiEndpoint {
	return &LiquiEndpoint{signer, NewKovanInterface()}
}

//NewDevLiquiEndpoint return dev endpoint instance
func NewDevLiquiEndpoint(signer Signer) *LiquiEndpoint {
	return &LiquiEndpoint{signer, NewDevInterface()}
}
