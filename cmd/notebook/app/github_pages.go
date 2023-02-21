package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bitsgofer/notebook-go/internal/fileutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	klog "k8s.io/klog/v2"
)

var setupBranchForGithubPages = &cobra.Command{
	Use:   "setupGithubPages",
	Short: "Setup branch to publish to Github Pages",
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := viper.GetString(cfgKeyPublishBranch)
		dir := viper.GetString(cfgKeyPublishDir)

		// clean up publishDir and publishBranch if they exist
		_, err := run("git", "show-ref", "--quiet", "refs/heads/"+branch)
		branchExist := err == nil
		_, err = os.Lstat(dir)
		dirExist := err == nil

		switch { // NOTE: actually harder to write w/o pattern matching like this
		case branchExist && dirExist:
			klog.V(2).InfoS("publishBranch and publishDir both exist => clean up")

			// TODO: maybe prompt if there is data inside `dir`
			mustRun("git", "worktree", "remove", dir, "--force") // also deletes `dir`
			mustRun("git", "branch", "-D", branch)
			cobra.CheckErr(fileutil.EnsureDir(dir))
		case branchExist && !dirExist:
			klog.V(2).InfoS("only publishBranch exist => clean up")

			mustRun("git", "branch", "-D", branch)
			cobra.CheckErr(fileutil.EnsureDir(dir))
		case !branchExist && dirExist:
			klog.V(2).InfoS("only publishDir exist => clean up")

			cobra.CheckErr(os.RemoveAll(dir))
			cobra.CheckErr(fileutil.EnsureDir(dir))
		default:
			// NOP
		}

		// setup publishDir as a worktree for publishBranch, then return to starting branch
		currentBranch, err := run("git", "rev-parse", "--abbrev-ref", "HEAD")
		cobra.CheckErr(err) // -> exit on errr
		mustRun("git", "checkout", "--orphan", branch)
		mustRun("git", "reset", "--hard")
		mustRun("git", "commit", "--allow-empty", "-m", "Init")
		mustRun("git", "checkout", currentBranch)
		mustRun("git", "worktree", "add", dir, branch)

		// done => suggest some optional steps
		fmt.Printf("Setup directory %q to publish to branch: %q\n", dir, branch)
		fmt.Println("")
		fmt.Println("Suggestions:")
		fmt.Printf("- Force-push the branch %q, e.g: `cd %s && git push --set-upstream origin %s --force`\n", branch, dir, branch)
		fmt.Printf("- Add folder %s to .gitignore\n", dir)

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

// mustRun wraps the code for executing a command and checking for exec err
// => make running lots of commands less tedious
func mustRun(cmd string, args ...string) {
	err := exec.Command(cmd, args...).Run()
	klog.V(4).InfoS("run command && exit on error", "cmd", cmd, "args", args, "err", err)

	cobra.CheckErr(err) // -> exit on err
}

// run wraps the code for executing a command and looking at its output
// => must call cobra.CheckErr() separately
func run(cmd string, args ...string) (string, error) {
	stdErrAndStdOut, err := exec.Command(cmd, args...).CombinedOutput()
	klog.V(4).InfoS("run command", "cmd", cmd, "args", args, "stdErrAndStdOut", stdErrAndStdOut, "err", err)

	// strip special chars from output
	output := strings.Trim(string(stdErrAndStdOut), "\n")

	return output, err
}
