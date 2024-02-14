package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Cfg struct {
	Viper    bool   `yaml:"Viper"`
	Commands []*Cmd `yaml:"Commands"`
	Cmds     []Command
}

type Command interface {
	Args() []Arg
	Aliases() []string
	Flags() []Flag
	Long() string
	Name() string
	Parent() string
	Run() string
	Short() string
	Use() string
}

type Cmd struct {
	Fargs    []Arg    `yaml:"Args"`
	Faliases []string `yaml:"Aliases"`
	Fflags   []Flag   `yaml:"Flags"`
	Flong    string   `yaml:"Long"`
	Fname    string   `yaml:"Name"`
	Fparent  string   `yaml:"Parent"`
	Frun     string   `yaml:"Run"`
	Fshort   string   `yaml:"Short"`
	Fuse     string   `yaml:"Use"`
}

func (c *Cmd) Args() []Arg {
	return c.Fargs
}

func (c *Cmd) Aliases() []string {
	return c.Faliases
}

func (c *Cmd) Flags() []Flag {
	return c.Fflags
}

func (c *Cmd) Long() string {
	return c.Flong
}

func (c *Cmd) Name() string {
	return c.Fname
}

func (c *Cmd) Parent() string {
	return c.Fparent
}

func (c *Cmd) Run() string {
	return c.Frun
}

func (c *Cmd) Short() string {
	return c.Fshort
}

func (c *Cmd) Use() string {
	return c.Fuse
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

var cfg = &Cfg{}

func readConfig(f string) (*Cfg, error) {

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

func genInit(cfg *Cfg) error {
	f, err := os.Open("cmd/commands.go")
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.ExecuteTemplate(f, "init", cfg)
	if err != nil {
		return err
	}

	err = fmtCommands("cmd/commands.go")
	if err != nil {
		return err
	}
	return nil
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

	for _, cmd := range cfg.Commands {
		cfg.Cmds = append(cfg.Cmds, cmd)
	}

	err = tmpl.ExecuteTemplate(f, "init", cfg)
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
		file := "cmd/" + c.Name() + ".go"
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

func NewCobraCmd(cmd Command) *cobra.Command {
	return &cobra.Command{
		Use:     cmd.Use(),
		Aliases: cmd.Aliases(),
		Short:   cmd.Short(),
		Long:    cmd.Long(),
		//Run:     cmd.Runner(),
	}
}

func init() {
	rootCmdCobra.AddCommand(genCmd)
}
