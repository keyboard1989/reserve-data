package main

import (
	"github.com/KyberNetwork/reserve-data/cmd/clicmd"
	"github.com/KyberNetwork/reserve-data/settings"
)

func main() {
	settings.NewSetting()
	cmd.Execute()
}
