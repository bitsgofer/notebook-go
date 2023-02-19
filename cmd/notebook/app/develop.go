package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var developCmd = &cobra.Command{
	Use:   "develop",
	Short: "Run render on file changes and serve over HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Run development web server")
		return nil
	},
}
