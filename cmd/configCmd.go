package cmd

import (
	"github.com/spf13/cobra"
)

func configCmdRun(cmd *cobra.Command, args []string) {
	out, err := execCmd("echo", "this is a shell command")
	if err != nil {
		panic(err)
	}
	if out != "" {
		println(out)
	}

}
