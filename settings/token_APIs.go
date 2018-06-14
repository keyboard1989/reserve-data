package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

func (setting *Settings) GetAllTokens() ([]common.Token, error) {
	return setting.Tokens.Storage.GetAllTokens()
}

func (setting *Settings) GetActiveTokens() ([]common.Token, error) {
	return setting.Tokens.Storage.GetActiveTokens()
}

func (setting *Settings) GetInternalTokens() ([]common.Token, error) {
	return setting.Tokens.Storage.GetInternalTokens()
}

func (setting *Settings) GetExternalTokens() ([]common.Token, error) {
	return setting.Tokens.Storage.GetExternalTokens()
}

func (setting *Settings) GetTokenByID(id string) (common.Token, error) {
	return setting.Tokens.Storage.GetTokenByID(id)
}

func (setting *Settings) GetActiveTokenByID(id string) (common.Token, error) {
	return setting.Tokens.Storage.GetActiveTokenByID(id)
}

func (setting *Settings) GetInternalTokenByID(id string) (common.Token, error) {
	return setting.Tokens.Storage.GetInternalTokenByID(id)
}

func (setting *Settings) GetExternalTokenByID(id string) (common.Token, error) {
	return setting.Tokens.Storage.GetExternalTokenByID(id)
}

func (setting *Settings) GetTokenByAddress(addr ethereum.Address) (common.Token, error) {
	return setting.Tokens.Storage.GetTokenByAddress(addr)
}

func (setting *Settings) GetActiveTokenByAddress(addr ethereum.Address) (common.Token, error) {
	return setting.Tokens.Storage.GetActiveTokenByAddress(addr)
}

func (setting *Settings) GetInternalTokenByAddress(addr ethereum.Address) (common.Token, error) {
	return setting.Tokens.Storage.GetInternalTokenByAddress(addr)
}

func (setting *Settings) GetExternalTokenByAddress(addr ethereum.Address) (common.Token, error) {
	return setting.Tokens.Storage.GetExternalTokenByAddress(addr)
}

func (setting *Settings) ETHToken() common.Token {
	eth, err := setting.Tokens.Storage.GetInternalTokenByID("ETH")
	if err != nil {
		log.Panicf("There is no ETH token in token DB, this should not happen (%s)", err)
	}
	return eth
}

func (setting *Settings) NewTokenPairFromID(base, quote string) (common.TokenPair, error) {
	bToken, err1 := setting.GetInternalTokenByID(base)
	qToken, err2 := setting.GetInternalTokenByID(quote)
	if err1 != nil || err2 != nil {
		return common.TokenPair{}, errors.New(fmt.Sprintf("%s or %s is not supported", base, quote))
	} else {
		return common.TokenPair{Base: bToken, Quote: qToken}, nil
	}
}

func (setting *Settings) MustCreateTokenPair(base, quote string) common.TokenPair {
	pair, err := setting.NewTokenPairFromID(base, quote)
	if err != nil {
		panic(err)
	}
	return pair
}

func (setting *Settings) UpdateToken(t common.Token) error {
	return setting.Tokens.Storage.UpdateToken(t)
}

func (setting *Settings) LoadTokenFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	tokens := TokenConfig{}
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &tokens); err != nil {
		return err
	}
	for id, t := range tokens.Tokens {
		token := common.NewToken(id, t.Name, t.Address, t.Decimal, t.Active, t.Internal, t.MinimalRecordResolution, t.MaxPerBlockImbalance, t.MaxTotalImbalance)
		if err = setting.UpdateToken(token); err != nil {
			return err
		}
	}
	return nil
}
