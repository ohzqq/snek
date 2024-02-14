package cmd

import (
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

// genCmd represents the genCmd command
var genCmd = &cobra.Command{
	Use:   "gen file",
	Short: "generate your cli",
	Long:  `scaffold your cli application using a yaml config file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f := args[0]
		err := readConfig(f)
		if err != nil {
			log.Fatal(err)
		}

		err = cfg.GenCmds()
		if err != nil {
			log.Fatal(err)
		}

		for _, c := range cfg.Commands {
			println(c.Cobra())
			for _, f := range c.Flags {
				println(f.Gen(c.Name))
			}
		}

		err = cfg.GenCmdFuncs()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func fmtCommands(file string) error {
	cmd := exec.Command("go", "fmt", file)
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(genCmd)
}
