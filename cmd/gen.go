package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Cfg struct {
	Viper    bool  `yaml:"Viper"`
	Commands []Cmd `yaml:"Commands"`
}

type Cmd struct {
	Args    []Arg    `yaml:"Args"`
	Aliases []string `yaml:"Aliases"`
	Flags   []Flag   `yaml:"Flags"`
	Long    string   `yaml:"Long"`
	Name    string   `yaml:"Name"`
	Parent  string   `yaml:"Parent"`
	Run     string   `yaml:"Run"`
	Short   string   `yaml:"Short"`
	Use     string   `yaml:"Use"`
}

type Flag struct {
	Name       string `yaml:"Name"`
	Shorthand  string `yaml:"Shorthand"`
	Usage      string `yaml:"Usage"`
	Type       string `yaml:"Type"`
	Var        string `yaml:"Var"`
	Value      string `yaml:"Value"`
	Persistent bool   `yaml:"Persistent"`
	Viper      bool   `yaml:"Viper"`
}

type Arg struct {
	Required bool   `yaml:"Required"`
	Name     string `yaml:"Name"`
}

// genCmd represents the genCmd command
var genCmd = &cobra.Command{
	Use:   "gen file",
	Short: "generate your cli",
	Long:  `scaffold your cli application using a yaml config file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f := args[0]
		cfg, err := readConfig(f)
		if err != nil {
			log.Fatal(err)
		}

		err = genCommands(cfg)
		if err != nil {
			log.Fatal(err)
		}

		err = genCommandFuncs(cfg)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func readConfig(f string) (*Cfg, error) {
	cfg := &Cfg{}

	d, err := os.ReadFile(f)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(d, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func genCommands(cfg *Cfg) error {
	f, err := os.Create("cmd/commands.go")
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.ExecuteTemplate(f, "cobra", cfg)
	if err != nil {
		return err
	}

	err = fmtCommands("cmd/commands.go")
	if err != nil {
		return err
	}
	return nil
}

func fmtCommands(file string) error {
	cmd := exec.Command("go", "fmt", file)
	return cmd.Run()
}

func genCommandFuncs(cfg *Cfg) error {
	for _, c := range cfg.Commands {
		file := "cmd/" + c.Name + ".go"
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		defer f.Close()

		err = tmpl.ExecuteTemplate(f, "run", c)
		if err != nil {
			return err
		}

		//err = fmtCommands(file)
		//if err != nil {
		//  return err
		//}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(genCmd)
}
