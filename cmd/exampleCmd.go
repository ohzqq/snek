package cmd

import (
	"github.com/spf13/cobra"
)

func exampleCmdRun(cmd *cobra.Command, args []string) {
	if cmd.Flags().Changed("flag") {
		// todo
	}
	if cmd.Flags().Changed("persistent") {
		// todo
	}

}
