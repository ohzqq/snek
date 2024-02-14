package cmd

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Cfg struct {
	Viper    bool  `yaml:"Viper"`
	Root     Cmd   `yaml:"Root"`
	Commands []Cmd `yaml:"Commands"`
}

var cfg = &Cfg{}

const (
	genCom   = `// Code generated by snek; DO NOT EDIT.`
	pkgCom   = "package cmd\n\n"
	pkgCobra = `"github.com/spf13/cobra"`
	pkgViper = `"github.com/spf13/viper"`
)

func (cfg *Cfg) Cmds() string {
	var cmd strings.Builder
	cmd.WriteString(genCom)
	cmd.WriteByte('\n')
	cmd.WriteString(pkgCom)

	imps := []string{pkgCobra}
	if cfg.Viper {
		imps = append(imps, pkgViper)
	}
	cmd.WriteString(imports(imps...))

	for _, c := range cfg.Commands {
		cmd.WriteString(c.Cobra())
		cmd.WriteByte('\n')
	}

	cmd.WriteString(cfg.Init())

	return cmd.String()
}

func (cfg *Cfg) RunFuncs() map[string]string {
	funcs := make(map[string]string)
	for _, cmd := range cfg.Commands {
		var fn strings.Builder
		fn.WriteString(pkgCom)
		fn.WriteString(imports(pkgCobra))
		fn.WriteString(cmd.runFunc())
		funcs[cmd.Name()] = fn.String()
	}
	return funcs
}

func (cfg *Cfg) root() string {
	return fmtCobraVar("rootCmd", cfg.Root.fields())
}

func (cfg *Cfg) Init() string {
	snek := []string{"func init() {"}
	for _, cmd := range cfg.Commands {
		snek = append(snek, cmd.add())
		for _, flag := range cmd.Flags {
			snek = append(snek, flag.Gen(cmd.Name()))
		}
	}
	snek = append(snek, "}")
	return strings.Join(snek, "\n")
}

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

func imports(imp ...string) string {
	var cmd strings.Builder
	cmd.WriteString("import (\n")
	for _, i := range imp {
		cmd.WriteString(i)
		cmd.WriteByte('\n')
	}
	cmd.WriteString(")\n\n")
	return cmd.String()
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
