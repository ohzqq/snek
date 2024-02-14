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

	if f.Type == "" {
		f.Type = "StringP"
	}
	p += f.Type

	return fmt.Sprintf("%s(%s)", p, f.args())
}

func (f Flag) args() string {
	var pos []string

	pos = append(pos, quote(f.Name))

	if s := f.short(); s != "" {
		pos = append(pos, quote(s))
	}

	if f.Value == "" {
		f.Value = quote("")
	}

	pos = append(pos, f.Value)

	pos = append(pos, quote(f.Usage))

	return strings.Join(pos, ",")
}

func (f Flag) short() string {
	if f.Shorthand != "" {
		return f.Shorthand
	}

	if strings.HasSuffix(f.Type, "P") {
		return string(f.Name[0])
	}

	return ""
}

func (f Flag) Changed() string {
	var cmd strings.Builder

	i := fmt.Sprintf(`if cmd.Flags().Changed("%s") {`, f.Name)
	cmd.WriteString(i)
	cmd.WriteString("\n// todo\n}\n")

	return cmd.String()
}
