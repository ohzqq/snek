package cmd

import (
	"fmt"
	"strings"
)

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

func (f Flag) Gen(cmd string) string {
	var flags []string

	fl := fmt.Sprintf("%s.%s", cmd, f.flag())
	flags = append(flags, fl)

	if f.Viper {
		v := fmt.Sprintf("viper.BindPFlag(\"%s\", %s.Flags().Lookup(\"%s\"))", f.Name, cmd, f.Name)
		flags = append(flags, v)
	}

	return strings.Join(flags, "\n")
}

func (f Flag) flag() string {
	p := "Flags()."
	if f.Persistent {
		p = "Persistent" + p
	}
	p += f.Type

	var pos []string
	pos = append(pos, quote(f.Name))
	if f.Shorthand != "" {
		pos = append(pos, quote(f.Shorthand))
	}
	pos = append(pos, f.Value)
	pos = append(pos, quote(f.Usage))

	return fmt.Sprintf("%s(%s)", p, strings.Join(pos, ","))
}
