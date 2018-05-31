package settings

import (
	"log"

	settingstorage "github.com/KyberNetwork/reserve-data/settings/storage"
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

func NewTokenSetting(dbPath string) *TokenSetting {
	BoltTokenStorage, err := settingstorage.NewBoltTokenStorage(dbPath)
	if err != nil {
		log.Panicf("Setting Init: Can not create bolt token storage (%s)", err)
	}
	tokenSetting := TokenSetting{BoltTokenStorage}
	return &tokenSetting

}
