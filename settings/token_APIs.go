package settings

import (
	"log"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

func GetAllTokens() ([]common.Token, error) {
	return setting.Tokens.Storage.GetAllTokens()
}

func GetActiveTokens() ([]common.Token, error) {
	return setting.Tokens.Storage.GetActiveTokens()
}

func GetInternalTokens() ([]common.Token, error) {
	return setting.Tokens.Storage.GetInternalTokens()
}

func GetExternalTokens() ([]common.Token, error) {
	return setting.Tokens.Storage.GetExternalTokens()
}

func GetTokenByID(id string) (common.Token, error) {
	return setting.Tokens.Storage.GetTokenByID(id)
}

func GetActiveTokenByID(id string) (common.Token, error) {
	return setting.Tokens.Storage.GetActiveTokenByID(id)
}

func GetInternalTokenByID(id string) (common.Token, error) {
	return setting.Tokens.Storage.GetInternalTokenByID(id)
}

func GetExternalTokenByID(id string) (common.Token, error) {
	return setting.Tokens.Storage.GetExternalTokenByID(id)
}

func GetTokenByAddress(addr ethereum.Address) (common.Token, error) {
	return setting.Tokens.Storage.GetTokenByAddress(addr)
}

func GetActiveTokenByAddress(addr ethereum.Address) (common.Token, error) {
	return setting.Tokens.Storage.GetActiveTokenByAddress(addr)
}

func GetInternalTokenByAddress(addr ethereum.Address) (common.Token, error) {
	return setting.Tokens.Storage.GetInternalTokenByAddress(addr)
}

func GetExternalTokenByAddress(addr ethereum.Address) (common.Token, error) {
	return setting.Tokens.Storage.GetExternalTokenByAddress(addr)
}

func ETHToken() common.Token {
	eth, err := setting.Tokens.Storage.GetInternalTokenByID("ETH")
	if err != nil {
		log.Panicf("There is no ETH token in token DB, this should not happen (%s)", err)
	}
	return eth
}
