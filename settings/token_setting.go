package settings

import (
	"encoding/json"
	"io/ioutil"

	"github.com/KyberNetwork/reserve-data/common"
)

type token struct {
	Address                 string `json:"address"`
	Name                    string `json:"name"`
	Decimal                 int64  `json:"decimals"`
	Active                  bool   `json:"internal use"`
	Internal                bool   `json:"listed"`
	MinimalRecordResolution string `json:"minimalRecordResolution"`
	MaxTotalImbalance       string `json:"maxPerBlockImbalance"`
	MaxPerBlockImbalance    string `json:"maxTotalImbalance"`
}

type TokenConfig struct {
	Tokens map[string]token `json:"tokens"`
}

type TokenSetting struct {
	Storage TokenStorage
}

func NewTokenSetting(storage TokenStorage) *TokenSetting {
	return &TokenSetting{storage}
}

func UpdateToken(t common.Token) error {
	if err := setting.Tokens.Storage.AddTokenByID(t); err != nil {
		return err
	}
	if err := setting.Tokens.Storage.AddTokenByAddress(t); err != nil {
		return err
	}
	return nil
}

func LoadTokenFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	tokens := TokenConfig{}
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return err
	}
	for id, t := range tokens.Tokens {
		token := common.NewToken(id, t.Name, t.Address, t.Decimal, t.Active, t.Internal, t.MinimalRecordResolution, t.MaxPerBlockImbalance, t.MaxTotalImbalance)
		err = UpdateToken(token)
		if err != nil {
			return err
		}
	}
	return nil
}
