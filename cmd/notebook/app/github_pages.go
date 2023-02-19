package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupBranchForGithubPages = &cobra.Command{
	Use:   "setupGithubPages",
	Short: "Setup branch to publish to Github Pages",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Setup branch to publish to Github Pages")
		return nil
	},
}

var publishToGithubPagesCmd = &cobra.Command{
	Use:   "publishToGithubPages",
	Short: "Publish to Github Pages",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Publish to Github Pages")
		return nil
	},
}
