package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure agc",
	Run:   func(cmd *cobra.Command, args []string) {},
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
		// if configFile not set and default not found, create default
		if _, ok := err.(viper.ConfigFileNotFoundError); ok && cfgFile == "" {
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
