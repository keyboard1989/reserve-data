package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	ihttp "github.com/KyberNetwork/reserve-data/http"
	"github.com/KyberNetwork/reserve-data/signer"
)

func SortByKey(params map[string]string) map[string]string {
	newParams := make(map[string]string, len(params))
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		newParams[key] = params[key]
	}
	return newParams
}

func MakeSign(req *http.Request, message string, nonce string) {

	//fmt.Println(message)
	fileSigner, _ := signer.NewFileSigner("/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_config.json")

	hmac512auth := ihttp.KNAuthentication{
		fileSigner.KNSecret,
		REQ_SESCRET,
		fileSigner.KNConfiguration,
		fileSigner.KNConfirmConf,
	}
	signed := hmac512auth.KNReadonlySign(message)
	req.Header.Add("nonce", nonce)
	req.Header.Add("signed", signed)
}

func GetResponse(method string, url string,
	params map[string]string, signNeeded bool, timepoint uint64) ([]byte, error) {
	params = SortByKey(params)
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	//create request
	req, ok := http.NewRequest(method, url, nil)
	if ok != nil {
		fmt.Println("can't establish request", ok)
	}
	// Add header
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// Create raw query
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	if signNeeded {
		nonce, ok := params["nonce"]
		if !ok {
			log.Printf("there was no nonce")
		} else {
			MakeSign(req, q.Encode(), nonce)
		}
	}
	//do the request and return the reply
	var err error
	var resp_body []byte
	resp, err := client.Do(req)
	if err != nil {
		return resp_body, err
	} else {
		defer resp.Body.Close()
		switch resp.StatusCode {
		case 200:
			resp_body, err = ioutil.ReadAll(resp.Body)
		default:
			log.Printf("The reply code %v was unexpected", resp.StatusCode)
			resp_body, err = ioutil.ReadAll(resp.Body)
		}
		log.Printf("\n request to %s, got response: \n %s \n\n", req.URL, common.TruncStr(resp_body))
		return resp_body, err
	}
}
