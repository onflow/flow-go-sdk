package cadence

import (
	"github.com/dapperlabs/cadence/runtime/cmd/execute"
	"github.com/spf13/cobra"

	"github.com/dapperlabs/flow-go-sdk/cli/flow/cadence/abi"
	"github.com/dapperlabs/flow-go-sdk/cli/flow/cadence/languageserver"
	"github.com/dapperlabs/flow-go-sdk/cli/flow/cadence/vscode"
)

var Cmd = &cobra.Command{
	Use:   "cadence",
	Short: "Execute Cadence code",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			execute.Execute(args)
		} else {
			execute.RunREPL()
		}
	},
}

func init() {
	Cmd.AddCommand(languageserver.Cmd)
	Cmd.AddCommand(vscode.Cmd)
	Cmd.AddCommand(abi.Cmd)
}
