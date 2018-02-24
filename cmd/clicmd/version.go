package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of the application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Kyber core 0.4.1 -- HEAD")
		},
	}
	RootCmd.AddCommand(versionCmd)
}
