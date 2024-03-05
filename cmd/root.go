package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snek",
	Short: "generate a cli for cobra",
	Long: `A Cobra companion for quickly generating a cli app with a config file.
`,
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
	},
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./snek.yaml", "config file (default is ./snek.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".snek" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".snek")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
