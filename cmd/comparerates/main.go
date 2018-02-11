package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	ihttp "github.com/KyberNetwork/reserve-data/http"
	"github.com/KyberNetwork/reserve-data/signer"
)

const (
	BASE_URL    string = "https://internal-mainnet-core.kyber.network"
	REQ_SESCRET string = "vtHpz1l0kxLyGc4R1qJBkFlQre5352xGJU9h8UQTwUTz5p6VrxcEslF4KnDI21s1"
	CONFIG_PATH string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_config.json"
)

type AllRateHTTPReply struct {
	Data    []common.AllRateResponse
	Success bool
}

type AllActionHTTPReply struct {
	Data    []common.ActivityRecord
	Success bool
}

func SortByKey(params map[string]string) map[string]string {
	//to be implement
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

	fmt.Println(message)
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
			log.Fatal("there was no nonce")
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
		case 429:
			err = errors.New("breaking a request rate limit.")
		case 418:
			err = errors.New("IP has been auto-banned for continuing to send requests after receiving 429 codes.")
		case 500:
			err = errors.New("500 from Binance, its fault.")
		case 200:
			resp_body, err = ioutil.ReadAll(resp.Body)
		}
		log.Printf("\n request to %s, got response: \n %s \n\n", req.URL, common.TruncStr(resp_body))
		return resp_body, err
	}
}

func GetActivitiesResponse(params map[string]string) (AllActionHTTPReply, error) {
	timepoint := (time.Now().UnixNano() / int64(time.Millisecond))
	nonce := strconv.FormatInt(timepoint, 10)
	var allActionRep AllActionHTTPReply
	params["nonce"] = nonce
	data, err := GetResponse("GET", fmt.Sprintf("%s/%s", BASE_URL, "activities"), params, true, uint64(timepoint))

	if err != nil {
		fmt.Println("can't get response", err)
	} else {
		if err := json.Unmarshal(data, &allActionRep); err != nil {
			fmt.Println("can't decode the reply", err)
			return allActionRep, err
		}
	}
	return allActionRep, nil
}

func GetAllRateResponse(params map[string]string) (AllRateHTTPReply, error) {
	timepoint := (time.Now().UnixNano() / int64(time.Millisecond))
	var allRateRep AllRateHTTPReply
	data, err := GetResponse("GET", fmt.Sprintf("%s/%s", BASE_URL, "get-all-rates"), params, false, uint64(timepoint))

	if err != nil {
		fmt.Println("can't get response", err)
	} else {
		if err := json.Unmarshal(data, &allRateRep); err != nil {
			fmt.Println("can't decode the reply", err)
			return allRateRep, err
		}
	}
	return allRateRep, nil
}

func RateDifference(r1, r2 float64) float64 {
	return (math.Abs(r1-r2) / r1)
}

func Compare(oneAct common.ActivityRecord, oneRate common.AllRateResponse, blockID uint64) {
	tokenIDs, asrt := oneAct.Params["tokens"].([]interface{})
	buys, asrt1 := oneAct.Params["buys"].([]interface{})
	sells, asrt2 := oneAct.Params["sells"].([]interface{})
	if asrt && asrt1 && asrt2 {
		for idx, tokenID := range tokenIDs {
			tokenid, _ := tokenID.(string)
			val, ok := oneRate.Data[tokenid]
			if ok {
				differ := RateDifference(val.BaseBuy, buys[idx].(float64)/1000000000000000000)
				if differ > 0.001 {
					fmt.Printf("block %d set a buy rate differ %.5f than get rate at token %s \n", blockID, differ, tokenid)
				}
				differ = RateDifference(val.BaseSell, sells[idx].(float64)/1000000000000000000)
				if differ > 0.001 {
					fmt.Printf("block %d set a sell rate differ %.5f than get rate at token %s \n", blockID, differ, tokenid)
				}
			}
		}
	}
}

func CompareRate(acts []common.ActivityRecord, rates []common.AllRateResponse) {
	idx := 0
	for _, oneAct := range acts {
		if oneAct.Action == "set_rates" {
			_, ok := oneAct.Params["block"]
			if ok {
				curBlock := uint64(oneAct.Params["block"].(float64))
				for (curBlock < rates[idx].ToBlockNumber) && (idx < len(rates)) {
					idx += 1
				}
				if (curBlock <= rates[idx].BlockNumber) && (curBlock >= rates[idx].ToBlockNumber) {
					fmt.Printf("\n Block %d is found between block %d to block %d \n", curBlock, rates[idx].BlockNumber, rates[idx].ToBlockNumber)
					Compare(oneAct, rates[idx], curBlock)
				} else {
					fmt.Printf("\n Block %d is not found\n", curBlock)
				}
			}
		}
	}
}

func main() {

	params := make(map[string]string)
	params["fromTime"] = "1518240610536"
	params["toTime"] = "1518247610536"
	allActionRep, err := GetActivitiesResponse(params)
	if err != nil {
		log.Fatal("couldn't get activites: ", err)
	}
	allRateRep, err := GetAllRateResponse(params)
	if err != nil {
		log.Fatal("couldn't get all rates: ", err)
	}
	CompareRate(allActionRep.Data, allRateRep.Data)
}
