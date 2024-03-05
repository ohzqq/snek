//go:build ignore

package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

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

		d, err := os.Create("cmd/commands.go")
		if err != nil {
			log.Fatal(err)
		}
		defer d.Close()

		_, err = d.WriteString(cfg.Cmds())

		err = fmtCommands("cmd/commands.go")
		if err != nil {
			log.Fatal(err)
		}

		for name, fn := range cfg.RunFuncs() {
			n := filepath.Join("cmd", name+".go")
			f, err := os.Create(n)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			_, err = f.WriteString(fn)
			if err != nil {
				log.Fatal(err)
			}

			err = fmtCommands(n)
			if err != nil {
				log.Fatalf("gofmt err %s\n%s", err, fn)
			}
		}
		//fmt.Printf("%+v\n", cfg.RunFuncs())

		//err = cfg.GenCmdFuncs()
		//if err != nil {
		//log.Fatal(err)
		//}
	},
}

func fmtCommands(file string) error {
	cmd := exec.Command("go", "fmt", file)
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(genCmd)
}
