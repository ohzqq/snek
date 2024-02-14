package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmdCobra = &cobra.Command{
	Use:   "snek",
	Short: "generate a cli for cobra",
	Long: `A Cobra companion for quickly generating a cli with a config file.
`,
}

func Execute() {
	err := rootCmdCobra.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.snek.yaml)")
}
