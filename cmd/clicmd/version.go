package cmd

import (
	"fmt"

	"github.com/KyberNetwork/reserve-data"
	"github.com/spf13/cobra"
)

func init() {

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of the application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(reserve.VERSION)
		},
	}
	RootCmd.AddCommand(versionCmd)
}
