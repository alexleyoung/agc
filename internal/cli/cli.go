package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexleyoung/agc/internal/ai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/genai"
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
		Run: func(cmd *cobra.Command, args []string) {
			query := strings.Join(args, " ")
			history := make([]*genai.Content, 0)
			resp, err := ai.Chat(cmd.Context(), "gemini-2.5-flash", history, query)
			if err != nil {
				cobra.CheckErr(err)
				return
			}
			fmt.Println(resp.Text())
		},
	}

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Configure agc",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	configSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set a configuration value",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cobra.CheckErr(fmt.Errorf("missing value"))
				return
			}
			viper.Set(args[0], args[1])
			err := viper.WriteConfig()
			cobra.CheckErr(err)
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/"+defaultPath+configName+"."+configType+")")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
