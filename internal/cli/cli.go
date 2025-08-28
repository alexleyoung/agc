package cli

import (
	"fmt"
	"os"

	"github.com/alexleyoung/agc/internal/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/agcli/config.toml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.AddCommand(config.Cmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search in $HOME/.config/agc for agcli.conf
		path := home + "/.config/agcli"
		viper.AddConfigPath(path)
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	var err error
	if err = viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		return
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// create file for user
		err = viper.SafeWriteConfig()
		if err != nil {
			cobra.CheckErr(err)
		}
		// read in new config
		viper.ReadInConfig()
		fmt.Printf("Config file created at: %s\n", viper.ConfigFileUsed())
	}
	cobra.CheckErr(err)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
