package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkDependenciesCmd = &cobra.Command{
	Use:   "checkDependencies",
	Short: "Check dependencies and suggest how to install them",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Check dependencies")
		return nil
	},
}
