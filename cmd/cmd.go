package cmd

import (
	"fmt"
	"strings"
)

type Cmd struct {
	Aliases []string `yaml:"Aliases"`
	Flags   []Flag   `yaml:"Flags"`
	Long    string   `yaml:"Long"`
	Name    string   `yaml:"Name"`
	Parent  string   `yaml:"Parent"`
	Run     string   `yaml:"Run"`
	Short   string   `yaml:"Short"`
	Use     string   `yaml:"Use"`
}

func (cmd Cmd) Cobra() string {
	var snek []string

	snek = append(snek, fmt.Sprintf("var %s = &cobra.Command{", cmd.Name))

	snek = append(snek, cmd.run())

	if cmd.Use != "" {
		snek = append(snek, Use.fmtField(cmd.Use))
	}

	if cmd.Short != "" {
		snek = append(snek, Short.fmtField(cmd.Short))
	}

	if cmd.Long != "" {
		snek = append(snek, Long.fmtField(cmd.Long))
	}

	if len(cmd.Aliases) > 0 {
		snek = append(snek, Aliases.fmtField(cmd.Aliases...))
	}

	snek = append(snek, "}")

	return strings.Join(snek, ",\n")
}

func (cmd Cmd) run() string {
	return fmt.Sprintf("Run: %sRun", cmd.Name)
}
