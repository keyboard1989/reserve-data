package settings

import (
	"encoding/json"
	"io/ioutil"

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

type ExchangeFeesConfig struct {
	Exchanges map[string]common.ExchangeFees `json:"exchanges"`
}

var exchangeNameValue = map[string]ExchangeName{
	"binance":         Binance,
	"bittrex":         Bittrex,
	"huobi":           Huobi,
	"stable_exchange": StableExchange,
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
	var exConfig ExchangeFeesConfig
	if err = json.Unmarshal(data, &exConfig); err != nil {
		return err
	}
	for ex, exFee := range exConfig.Exchanges {
		if err = setting.Exchange.Storage.StoreFee(exchangeNameValue[ex], exFee); err != nil {
			return err
		}
	}
	return nil
}

func (setting *Settings) LoadDepositAddressFromFile(path string) error {
	return nil
}

func (setting *Settings) LoadMinDepositAddressFromFile(path string) error {
	return nil
}
