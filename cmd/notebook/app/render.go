package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render articles, drafts and themes",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Render to HTML")
		return nil
	},
}
