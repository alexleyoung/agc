package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var defaultPath = "/.config/agc/"
var configType = "toml"
var configName = "config"

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "agc",
		Short: "Interact with Google Calendar through natural language",
		Long: `agc is a command line tool that allows you to interact with 
Google Calendar via through natural language.`,
		Run: func(cmd *cobra.Command, args []string) {},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/"+defaultPath+configName+"."+configType+")")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(authCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
