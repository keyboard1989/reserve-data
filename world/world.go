package world

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
)

type TheWorld struct {
	endpoint Endpoint
}

func (self *TheWorld) GetGoldInfo() (common.GoldData, error) {
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	url := self.endpoint.GoldDataEndpoint()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	var err error
	var resp_body []byte
	log.Printf("request to gold feed endpoint: %s", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return common.GoldData{}, err
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			err = errors.New(fmt.Sprintf("Gold feed returned with code: %d", resp.StatusCode))
			return common.GoldData{}, err
		}
		resp_body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return common.GoldData{}, err
		}
		log.Printf("request to %s, got response from gold feed %s", req.URL, resp_body)
		result := common.GoldData{}
		err = json.Unmarshal(resp_body, &result)
		return result, err
	}
}

func NewTheWorld(env string) *TheWorld {
	switch env {
	case "dev":
		return &TheWorld{RealEndpoint{}}
	case "kovan":
		return &TheWorld{RealEndpoint{}}
	case "mainnet":
		return &TheWorld{RealEndpoint{}}
	case "staging":
		return &TheWorld{RealEndpoint{}}
	case "simulation":
		return &TheWorld{SimulatedEndpoint{}}
	case "ropsten":
		return &TheWorld{RealEndpoint{}}
	}
	panic("unsupported environment")
}
