package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/KyberNetwork/reserve-data/common"
)

type token struct {
	Address          string `json:"address"`
	Name             string `json:"name"`
	Decimals         int64  `json:"decimals"`
	KNReserveSupport bool   `json:"internal use"`
	Active           bool   `json:"listed"`
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

func AddToken(t common.Token, active bool, knSupported bool) error {
	if err := setting.Tokens.Storage.AddTokenByID(t); err != nil {
		return err
	}
	if err := setting.Tokens.Storage.AddTokenByAddress(t); err != nil {
		return err
	}
	if active {
		if err := setting.Tokens.Storage.AddActiveTokenByID(t); err != nil {
			return err
		}
		if err := setting.Tokens.Storage.AddActiveTokenByAddress(t); err != nil {
			return err
		}
		if knSupported {
			if err := setting.Tokens.Storage.AddInternalTokenByID(t); err != nil {
				return err
			}
			if err := setting.Tokens.Storage.AddInternalTokenByAddress(t); err != nil {
				return err
			}
		} else {
			if err := setting.Tokens.Storage.AddExternalTokenByID(t); err != nil {
				return err
			}
			if err := setting.Tokens.Storage.AddExternalTokenByAddress(t); err != nil {
				return err
			}
		}
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
		tok := common.Token{id, t.Address, t.Decimals}
		err = AddToken(tok, t.Active, t.KNReserveSupport)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewTokenPair(base, quote string) (common.TokenPair, error) {
	bToken, err1 := GetInternalTokenByID(base)
	qToken, err2 := GetInternalTokenByID(quote)
	if err1 != nil || err2 != nil {
		return common.TokenPair{}, errors.New(fmt.Sprintf("%s or %s is not supported", base, quote))
	} else {
		return common.TokenPair{bToken, qToken}, nil
	}
}

func MustCreateTokenPair(base, quote string) common.TokenPair {
	pair, err := NewTokenPair(base, quote)
	if err != nil {
		panic(err)
	} else {
		return pair
	}
}
