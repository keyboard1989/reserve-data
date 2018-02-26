package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "./cmd ",
	Short: "entry point to the application, required KYBER_ENV and KYBER_EXCHANGES as environment variables",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.Flags().BoolP("verbose", "v", false, "verbose mode enable")
}

// initConfig reads in config file and ENV variables if set.
// currently due to the fact that all configuration files are read seperatedly,
// Viper is not a good choice for this current development. Hence initConfig is empty.
func initConfig() {
}
