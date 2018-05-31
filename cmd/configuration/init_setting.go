package configuration

import (
	"path/filepath"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/settings"
)

const (
	tokenDBFileName         string = "token.db"
	tokenSettingDefaultName string = "mainnet_setting.json"
	tokenSettingSimName     string = "shared/deployment_dev.json"
)

func createSetting() *settings.Settings {
	tokensSetting := settings.NewTokenSetting(filepath.Join(common.CmdDirLocation(), tokenDBFileName))
	setting := settings.NewSetting(tokensSetting)
	setting.HandleEmptyToken(filepath.Join(common.CmdDirLocation(), tokenSettingDefaultName), filepath.Join(common.CmdDirLocation(), tokenSettingSimName))
	return setting
}
