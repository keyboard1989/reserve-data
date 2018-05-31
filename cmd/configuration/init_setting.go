package configuration

import (
	"path/filepath"

	"github.com/KyberNetwork/reserve-data/cmd/configuration/mode"
	"github.com/KyberNetwork/reserve-data/settings"
)

const (
	tokenDBFileName         string = "token.db"
	tokenSettingDefaultName string = "mainnet_setting.json"
	tokenSettingSimName     string = "shared/deployment_dev.json"
)

func createSetting() *settings.Settings {
	tokensSetting := settings.NewTokenSetting(filepath.Join(mode.CmdDirLocation(), tokenDBFileName))
	setting := settings.NewSetting(tokensSetting)
	setting.HandleEmptyToken(filepath.Join(mode.CmdDirLocation(), tokenSettingDefaultName), filepath.Join(mode.CmdDirLocation(), tokenSettingSimName))
	return setting
}
