package app

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"
)

var dependentBinaries = map[string]struct {
	recommendedVersion string
	instructions       string
}{
	"git": {
		recommendedVersion: "2.34+",
		instructions:       "Install from https://git-scm.com/download/linux",
	},
	// TODO: bundle pandoc (might not be possible since it's a Haskell program)
	"pandoc": {
		recommendedVersion: "3.1",
		instructions:       "Install release from https://github.com/jgm/pandoc/releases",
	},
	// TODO: bundle minify (it can be imported as a Go library)
	"minify": {
		recommendedVersion: "v2.12.4",
		instructions:       "Build from https://github.com/tdewolff/minify/releases/tag/v2.12.4",
	},
}

var checkDependenciesCmd = &cobra.Command{
	Use:   "checkDependencies",
	Short: "Check dependencies and suggest how to install them",
	RunE: func(cmd *cobra.Command, args []string) error {
		allInstalled := true

		fmt.Printf("Checking binary dependencies:\n")
		for bin, info := range dependentBinaries {
			if !isBinaryInPath(bin) {
				allInstalled = false

				fmt.Println("")
				fmt.Printf("    %s is not found in $PATH:\n", bin)
				fmt.Printf("      - Install instruction: %s.\n", info.instructions)
				fmt.Printf("      - Recommended version: %s.\n", info.recommendedVersion)
			}
		}

		fmt.Println("")
		if !allInstalled {
			fmt.Printf("=> One or more binarie(s) are not installed.\n")
			return nil
		}
		fmt.Printf("=> All dependent binaries are installed.\n")
		return nil
	},
}

// isBinaryInPath checks if the given binary can be executed (included in $PATH)
func isBinaryInPath(name string) bool {
	_, err := exec.LookPath(name)
	if err != nil {
		klog.V(2).ErrorS(err, "cannot find dependent binary", "name", name)
		return false
	}

	return true
}
