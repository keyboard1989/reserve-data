package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/KyberNetwork/reserve-data/common"
)

// ExchangeName is the name of exchanges of which core will use to rebalance.
//go:generate stringer -type=ExchangeName
type ExchangeName int

const (
	Binance ExchangeName = iota
	Bittrex
	Huobi
	StableExchange
)
const exchange_env string = "KYBER_EXCHANGES"

type ExchangeFeesConfig struct {
	Exchanges map[string]common.ExchangeFees `json:"exchanges"`
}

var exchangeNameValue = map[string]ExchangeName{
	"binance":         Binance,
	"bittrex":         Bittrex,
	"huobi":           Huobi,
	"stable_exchange": StableExchange,
}

func RunningExchange() []string {
	exchangesStr, ok := os.LookupEnv(exchange_env)
	if !ok {
		log.Print("WARNING: core is running without exchange")
		return nil
	}
	exchanges := strings.Split(exchangesStr, ",")
	log.Printf("exchanges %v", exchanges)
	return exchanges
}

func ExchangTypeValues() map[string]ExchangeName {
	return exchangeNameValue
}

type ExchangeSetting struct {
	Storage ExchangeStorage
}

func NewExchangeSetting(exchangeStorage ExchangeStorage) (*ExchangeSetting, error) {
	return &ExchangeSetting{exchangeStorage}, nil
}

func (setting *Settings) LoadFeeFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var exFeeConfig ExchangeFeesConfig
	if err = json.Unmarshal(data, &exFeeConfig); err != nil {
		return err
	}
	runningExs := RunningExchange()

	for _, ex := range runningExs {
		//Check if the exchange is in current code deployment.
		exName, ok := exchangeNameValue[ex]
		if !ok {
			return fmt.Errorf("Exchange %s is in KYBER_EXCHANGES, but not avail in current deployment", ex)
		}
		//Check if the current database has a record for such exchange
		_, err := setting.Exchange.Storage.GetFee(exName)
		if err != nil {
			log.Printf("Exchange %s is in KYBER_EXCHANGES but can't load fee in Database (%s). atempt to load it from config file", string(exName), err.Error())
			//Check if the config file has config for such exchange
			exFee, ok := exFeeConfig.Exchanges[ex]
			if !ok {
				log.Printf("Warning: Exchange %s is in KYBER_EXCHANGES, but not avail in Fee config file.")
				continue
			}
			if err = setting.Exchange.Storage.StoreFee(exName, exFee); err != nil {
				return err
			}
		}
	}
	return nil
}

type ExchangesMinDepositConfig struct {
	Exchanges map[string]common.ExchangesMinDeposit `json:"exchanges"`
}

func (setting *Settings) LoadMinDepositFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var exConfig ExchangesMinDepositConfig
	err = json.Unmarshal(data, &exConfig)
	if err != nil {
		return err
	}
	runningExs := RunningExchange()
	for _, ex := range runningExs {
		//Check if the exchange is in current code deployment.
		exName, ok := exchangeNameValue[ex]
		if !ok {
			return fmt.Errorf("Exchange %s is in KYBER_EXCHANGES, but not avail in current deployment", ex)
		}
		//Check if the current database has a record for such exchange
		_, err := setting.Exchange.Storage.GetMinDeposit(exName)
		if err != nil {
			log.Printf("Exchange %s is in KYBER_EXCHANGES but can't load fee in Database (%s). atempt to load it from config file", string(exName), err.Error())
			//Check if the config file has config for such exchange
			minDepo, ok := exConfig.Exchanges[ex]
			if !ok {
				log.Printf("Warning: Exchange %s is in KYBER_EXCHANGES, but not avail in config file")
				continue
			}
			if err = setting.Exchange.Storage.StoreMinDeposit(exName, minDepo); err != nil {
				return err
			}
		}
	}
	return nil
}

func (setting *Settings) LoadMinDepositAddressFromFile(path string) error {
	return nil
}
