package cmd

import (
	"github.com/spf13/cobra"
)

func configCmdRun(cmd *cobra.Command, args []string) {
	out, err := execCmd("echo", getExampleCfg())
	if err != nil {
		panic(err)
	}
	if out != "" {
		println(out)
	}

}
