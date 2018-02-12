package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
)

const (
	BASE_URL    string        = "https://internal-mainnet-core.kyber.network"
	REQ_SESCRET string        = "vtHpz1l0kxLyGc4R1qJBkFlQre5352xGJU9h8UQTwUTz5p6VrxcEslF4KnDI21s1"
	CONFIG_PATH string        = "/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_config.json"
	TWEI_ADJUST float64       = 1000000000000000000
	SLEEP_TIME  time.Duration = 60 //sleep time for forever run mode
	DIFFER_RATE float64       = 0.001
)

type AllRateHTTPReply struct {
	Data    []common.AllRateResponse
	Success bool
}

type AllActionHTTPReply struct {
	Data    []common.ActivityRecord
	Success bool
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
	return ((r2 - r1) / r1)
}

func printInterfaceMap(inf map[string]interface{}) {
	for k, v := range inf {
		switch v.(type) {
		case float64:
			log.Printf("\t\t %s %.5f \n", k, v.(float64))
		case []interface{}:
			log.Printf("\t\t %s \n", k)
			for _, vv := range v.([]interface{}) {
				log.Printf("\t\t\t %v", vv)
			}
		}
	}
}

func printAction(oneAct common.ActivityRecord) {
	i := int64(oneAct.Timestamp.ToUint64()) / 1000

	log.Printf("\t Time: %v", time.Unix(i, 0))
	log.Printf("\t Activity : %v\n", oneAct.Action)
	log.Printf("\t Destination : %v\n", oneAct.Destination)
	log.Printf("\t ExchangeStatus : %v\n", oneAct.ExchangeStatus)
	log.Printf("\t ID : %v\n", oneAct.ID)
	log.Printf("\t MiningStatus : %v\n", oneAct.MiningStatus)
	log.Printf("\t Params : \n")
	printInterfaceMap(oneAct.Params)
	log.Printf("\t Result : \n")
	printInterfaceMap(oneAct.Params)
}

func CompareRate(oneAct common.ActivityRecord, oneRate common.AllRateResponse, blockID uint64) {
	tokenIDs, asrt := oneAct.Params["tokens"].([]interface{})
	buys, asrt1 := oneAct.Params["buys"].([]interface{})
	sells, asrt2 := oneAct.Params["sells"].([]interface{})
	warning := false
	if asrt && asrt1 && asrt2 {
		for idx, tokenID := range tokenIDs {
			tokenid, _ := tokenID.(string)
			val, ok := oneRate.Data[tokenid]
			if ok {
				differ := RateDifference(val.BaseBuy*(1+float64(val.CompactBuy)/1000.0)*TWEI_ADJUST, buys[idx].(float64))
				if math.Abs(differ) > DIFFER_RATE {
					defer log.Printf("block %d set a buys rate differ %.5f%% than get rate at token %s \n", blockID, differ*100, tokenid)
					warning = true
				}
				differ = RateDifference(val.BaseSell*(1+float64(val.CompactSell)/1000.0)*TWEI_ADJUST, sells[idx].(float64))
				if math.Abs(differ) > DIFFER_RATE {
					defer log.Printf("block %d set a sell rate differ %.5f%% than get rate at token %s \n", blockID, differ*100, tokenid)
					warning = true
				}
			}
		}
	}
	if warning {
		log.Printf("There was different in set rate at block %d \n", blockID)
		printAction(oneAct)
	}
}

func CompareRates(acts []common.ActivityRecord, rates []common.AllRateResponse) {
	idx := 0
	for _, oneAct := range acts {
		if (oneAct.Action == "set_rates") && (oneAct.MiningStatus == "mined") {
			_, ok := oneAct.Result["blockNumber"]
			if ok {
				curBlock := uint64(oneAct.Result["blockNumber"].(float64))
				for (idx < len(rates)) && (curBlock < rates[idx].ToBlockNumber) {
					idx += 1
				}
				if (idx < len(rates)) && (curBlock <= rates[idx].BlockNumber) && (curBlock >= rates[idx].ToBlockNumber) {
					fmt.Printf("\n Block %d is found between block %d to block %d \n", curBlock, rates[idx].BlockNumber, rates[idx].ToBlockNumber)
					CompareRate(oneAct, rates[idx], curBlock)
				} else {
					fmt.Printf("\n Block %d is not found\n", curBlock)
				}
			}
		}
	}
}

func doQuery(params map[string]string) {
	allActionRep, err := GetActivitiesResponse(params)
	if err != nil {
		log.Printf("couldn't get activites: ", err)
		return
	}
	allRateRep, err := GetAllRateResponse(params)
	if err != nil {
		log.Printf("couldn't get all rates: ", err)
		return
	}
	if (len(allActionRep.Data) < 1) || (len(allRateRep.Data) < 1) {
		log.Printf("One of the reply was empty")
		return
	}
	CompareRates(allActionRep.Data, allRateRep.Data)
}

func main() {
	params := make(map[string]string)
	params["fromTime"] = os.Getenv("FROMTIME")
	params["toTime"] = os.Getenv("TOTIME")
	if len(params["fromTime"]) < 1 {
		log.Fatal("Wrong usage \n FROMTIME=<timestamp> [TOTIME=<totime>] ./compareRates")
	}
	if len(params["toTime"]) < 1 {
		log.Printf("There was no end time, go to foverer run mode...")
		for {
			params["toTime"] = strconv.FormatInt((time.Now().UnixNano() / int64(time.Millisecond)), 10)
			doQuery(params)
			time.Sleep(SLEEP_TIME * time.Second)
			params["fromTime"] = params["toTime"]
		}

	} else {
		log.Printf("Go to single query returning mode")
		doQuery(params)
	}

}
