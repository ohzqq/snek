package cmd

import (
	"fmt"
	"strings"
)

type Cmd struct {
	Aliases   []string `yaml:"Aliases"`
	Flags     []Flag   `yaml:"Flags"`
	Long      string   `yaml:"Long"`
	Parent    string   `yaml:"Parent"`
	NoCommand bool     `yaml:"NoCommand"`
	Short     string   `yaml:"Short"`
	Use       string   `yaml:"Use"`
}

func (cmd Cmd) Cobra() string {
	return fmtCobraVar(cmd.Name(), cmd.fields())
}

func (cmd Cmd) fields() []string {
	var snek []string

	if cmd.Use != "" {
		snek = append(snek, Use.fmtField(cmd.Use))
	}

	if len(cmd.Aliases) > 0 {
		snek = append(snek, Aliases.fmtField(cmd.Aliases...))
	}

	if cmd.Short != "" {
		snek = append(snek, Short.fmtField(cmd.Short))
	}

	if cmd.Long != "" {
		snek = append(snek, Long.fmtField(cmd.Long))
	}

	if !cmd.NoCommand {
		snek = append(snek, cmd.run())
	}

	return snek
}

func fmtCobraVar(name string, snek []string) string {
	return fmt.Sprintf(
		"var %s = &cobra.Command{\n%s,\n}",
		name,
		strings.Join(snek, ",\n"),
	)
}

func (cmd Cmd) Name() string {
	u, _, _ := strings.Cut(cmd.Use, " ")
	return u + "Cmd"
}

func (cmd Cmd) run() string {
	return fmt.Sprintf("Run: %sRun", cmd.Name())
}

func (cmd Cmd) add() string {
	p := "rootCmd"
	if cmd.Parent != "" {
		p = cmd.Parent
	}
	return fmt.Sprintf("%s.AddCommand(%s)", p, cmd.Name())
}
