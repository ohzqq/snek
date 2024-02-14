package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Cmd struct {
	Aliases   []string `yaml:"Aliases"`
	Exec      []string `yaml:"Exec"`
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

func (c Cmd) runFunc() string {
	var cmd strings.Builder
	fn := fmt.Sprintf("func %sRun(cmd *cobra.Command, args []string) {\n", c.Name())
	cmd.WriteString(fn)

	for _, flag := range c.Flags {
		cmd.WriteString(flag.Changed())
	}

	if c.shellout() {
		cmd.WriteString("out, err := execCmd(")
		cmd.WriteString(fmtSlice(c.Exec))
		cmd.WriteString(")\n")
		cmd.WriteString("if err != nil {\n")
		cmd.WriteString("panic(err)\n")
		cmd.WriteString("}\n")
		cmd.WriteString(`if out != "" {`)
		cmd.WriteByte('\n')
		cmd.WriteString("println(out)")
		cmd.WriteString("}\n")
	} else {
		cmd.WriteString("println(cmd.Name())")
	}

	cmd.WriteString("\n}\n")

	return cmd.String()
}

func (c Cmd) shellout() bool {
	return len(c.Exec) > 0
}

func execCmd(args ...string) (string, error) {
	if len(args) < 1 {
		return "", nil
	}

	cmd := exec.Command(args[0], args[0:]...)

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("finished with error: %v\n", stderr.String())
	}

	var cmdErr error
	if len(stderr.Bytes()) > 0 {
		cmdErr = errors.New(stderr.String())
	}

	var output string
	if len(stdout.Bytes()) > 0 {
		output = stdout.String()
	}

	return output, cmdErr
}

func fmtCobraVar(name string, snek []string) string {
	return fmt.Sprintf(
		"var %s = &cobra.Command{\n%s,\n}",
		name,
		strings.Join(snek, ",\n"),
	)
}
