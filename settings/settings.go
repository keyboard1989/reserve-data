package settings

import (
	"log"

	settingstorage "github.com/KyberNetwork/reserve-data/settings/storage"
)

const (
	TOKEN_DB_FILE_PATH      string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/token.db"
	TOKEN_DEFAULT_JSON_PATH string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/token.json"
)

type Settings struct {
	Tokens *TokenSetting
}

var setting Settings

func NewToken() *TokenSetting {
	BoltTokenStorage, err := settingstorage.NewBoltTokenStorage(TOKEN_DB_FILE_PATH)
	if err != nil {
		log.Panicf("Setting Init: Can not create bolt token storage", err)
	}
	tokenSetting := TokenSetting{BoltTokenStorage}
	return &tokenSetting
}

func NewSetting() *Settings {
	tokensSetting := NewToken()
	allToks, err := GetAllTokens()
	if err != nil || len(allToks) < 1 {
		log.Printf("Setting Init: Token DB is empty, attempt to load token from file")
		err := LoadTokenFromFile(TOKEN_DEFAULT_JSON_PATH)
		if err != nil {
			log.Printf("Setting Init: Can not load Token from file, Token DB is needed to be updated manually")
		}
	}
	overalSetting := Settings{tokensSetting}
	return &overalSetting
}
