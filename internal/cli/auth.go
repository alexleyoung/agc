package cli

import "github.com/spf13/cobra"

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Login to Google Calendar",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
