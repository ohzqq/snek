package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snek",
	Short: "generate a cli for cobra",
	Long: `A Cobra companion for quickly generating a cli app with a config file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := genCLI()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func genCLI() error {
	f := viper.GetString("config")
	err := readConfig(f)
	if err != nil {
		return err
	}

	d, err := os.Create("cmd/commands.go")
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = d.WriteString(cfg.Cmds())

	err = fmtCommands("cmd/commands.go")
	if err != nil {
		return err
	}

	for name, fn := range cfg.RunFuncs() {
		n := filepath.Join("cmd", name+".go")
		f, err := os.Create(n)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.WriteString(fn)
		if err != nil {
			return err
		}

		err = fmtCommands(n)
		if err != nil {
			return fmt.Errorf("gofmt err %s\n%s", err, fn)
		}
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("config", "c", "./snek.yaml", "config file (default is ./snek.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().Bool("example", false, "print example config to stdout")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgFile := viper.GetString("config")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func fmtCommands(file string) error {
	cmd := exec.Command("go", "fmt", file)
	return cmd.Run()
}
