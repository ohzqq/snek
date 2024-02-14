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
	var str strings.Builder
	str.WriteString("var ")
	str.WriteString(cmd.Name)
	str.WriteString(" = &cobra.Command{")
	str.WriteString("\n")

	if cmd.Use != "" {
		str.WriteString(fmtField(Use, cmd.Use))
	}

	if cmd.Short != "" {
		str.WriteString(fmtField(Short, cmd.Short))
	}

	if cmd.Long != "" {
		str.WriteString(fmtField(Long, cmd.Long))
	}

	str.WriteByte('}')

	return str.String()
}

func fmtField(n Field, v string) string {
	return fmt.Sprintf("%s: \"%s\",\n", n, v)
}
