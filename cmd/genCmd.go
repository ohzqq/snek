package cmd

import (
	"github.com/spf13/cobra"
)

func genCmdRun(cmd *cobra.Command, args []string) {
	println(cmd.Name())
}
