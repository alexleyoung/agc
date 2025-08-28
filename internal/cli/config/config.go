package config

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Configure agc",
	Long:  `Configure agc with your Google account and preferences.`,
}
