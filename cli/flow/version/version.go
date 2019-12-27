package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dapperlabs/flow-go-sdk/utils/build"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "View version and commit information",
	Run: func(cmd *cobra.Command, args []string) {
		semver := build.Semver()
		commit := build.Commit()

		// Print version/commit strings if they are known
		if build.IsDefined(semver) {
			fmt.Printf("Version: %s\n", semver)
		}
		if build.IsDefined(commit) {
			fmt.Printf("Commit: %s\n", commit)
		}
		// If no version info is known print a message to indicate this.
		if !build.IsDefined(semver) && !build.IsDefined(commit) {
			fmt.Printf("Version information unknown!\n")
		}
	},
}
