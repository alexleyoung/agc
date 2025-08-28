package cli

import (
	"fmt"
	"os"

	"github.com/alexleyoung/agc/internal/cli/config"
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

	rootCmd.AddCommand(config.Cmd)
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	defaultPath = home + defaultPath

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Search in $HOME/.config/agc for agc.conf
		viper.AddConfigPath(defaultPath)
		viper.SetConfigType(configType)
		viper.SetConfigName(configName)
	}

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// mkdir
			if err := os.MkdirAll(defaultPath, os.ModePerm); err != nil {
				cobra.CheckErr(err)
			}

			// create file for user
			err = viper.SafeWriteConfig()
			if err != nil {
				cobra.CheckErr(err)
			}
			// read in new config
			viper.ReadInConfig()
		} else {
			cobra.CheckErr(err)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
