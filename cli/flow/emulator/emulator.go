package emulator

import (
	"github.com/spf13/cobra"

	emu "github.com/dapperlabs/flow-emulator/cmd"
)

var Cmd = &cobra.Command{
	Use:              "emulator",
	Short:            "Flow emulator server",
	TraverseChildren: true,
}

func init() {
	Cmd.AddCommand(emu.Cmd)
}
