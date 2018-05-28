package settings

import (
	"log"
	"os"
)

const (
	TOKEN_DB_FILE_PATH          string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/token.db"
	TOKEN_DEFAULT_JSON_PATH     string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_setting.json"
	TOKEN_DEFAULT_JSON_SIM_PATH string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/shared/deployment_dev.json"
)

type Settings struct {
	Tokens *TokenSetting
}

var setting Settings

func NewSetting() *Settings {
	tokensSetting := NewTokenSetting()
	setting = Settings{tokensSetting}
	allToks, err := GetAllTokens()
	if err != nil || len(allToks) < 1 {
		if err != nil {
			log.Printf("Setting Init: Token DB is faulty (%s), attempt to load token from file", err)
		} else {
			log.Printf("Setting Init: Token DB is empty, attempt to load token from file")
		}
		tokenPath := TOKEN_DEFAULT_JSON_PATH
		if os.Getenv("KYBER_ENV") == "simulation" {
			tokenPath = TOKEN_DEFAULT_JSON_SIM_PATH
		}

		if err = LoadTokenFromFile(tokenPath); err != nil {
			log.Printf("Setting Init: Can not load Token from file: %s, Token DB is needed to be updated manually", err)
		}
	}
	overalSetting := Settings{tokensSetting}
	return &overalSetting
}
