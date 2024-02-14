package cmd

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Cfg struct {
	Viper    bool  `yaml:"Viper"`
	Commands []Cmd `yaml:"Commands"`
}

var cfg = &Cfg{}

func readConfig(f string) error {
	d, err := os.ReadFile(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(d, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Cfg) GenCmds() error {
	return genCommands(cfg)
}

func (cfg *Cfg) GenCmdFuncs() error {
	return genCommandFuncs(cfg)
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
