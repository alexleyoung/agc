package cli

import "github.com/spf13/cobra"

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query AGC with natural language",
}
