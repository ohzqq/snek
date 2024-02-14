package cmd

import (
	"fmt"
	"strings"
)

type Field int

//go:generate stringer -type Field
const (
	Aliases Field = iota
	Flags
	Long
	Name
	Parent
	Run
	Short
	Use
)

func (f Field) fmtField(v ...string) string {
	switch f {
	case Aliases:
		return fmtSlice(f, v)
	default:
		return fmtField(f, v[0])
	}
}

func quote(v string) string {
	return `"` + v + `"`
}

func fmtSlice(n Field, vals []string) string {
	s := make([]string, len(vals))
	for i, v := range vals {
		s[i] = `"` + v + `"`
	}
	return fmt.Sprintf("%s: []string{%s}", n, strings.Join(s, ","))
}

func fmtField(n Field, v string) string {
	return fmt.Sprintf("%s: \"%s\"", n, v)
}
