package cmd

import "text/template"

var tmpl = template.Must(template.New("cobra").Parse(commandTmpl))

const commandTmpl = `
{{- define "cobra"}}
// Code generated by snek; DO NOT EDIT.

package cmd

import (
	"github.com/spf13/cobra"
{{- if .Viper }}
	"github.com/spf13/viper"{{ end }}
)

{{range .Commands}}
	{{- $cmd := .Name -}}

type {{.Use}}Command struct{*Cmd}

var (
	{{.Use}}Cmd = {{.Use}}Command{
		Cmd: &Cmd {
			{{with .Use}}Fuse: "{{.}}",{{end}}
			{{with .Short}}Fshort: "{{.}}",{{end}}
			{{with .Long}}Flong: "{{.}}",{{end}}
			{{- with .Aliases}}
			Faliases: []string{
				{{range .}}"{{.}}",{{end}}
			},
			{{end}}
		},
	}
	{{$cmd}}Cobra *cobra.Command
)
{{end}}

func init() {
{{- range .Commands -}}
	{{- $cmd := .Name -}}
	
	{{$cmd}}Cobra = NewCobraCmd({{$cmd}})

	{{- with .Parent}}
		{{.}}Cobra.AddCommand({{$cmd}}Cobra)
	{{- end}}

	{{range .Flags -}}
		{{$cmd}}Cobra.
		{{- with .Persistent}}Persistent{{end -}}
		Flags().{{.Type}}("{{.Name}}",
		{{- with .Shorthand}}"{{.}}",{{end}}
		{{- .Value -}}
		, "{{.Usage}}")

		{{if .Viper}}
			viper.BindPFlag("{{.Name}}", {{$cmd}}Cobra.Flags().Lookup("{{.Name}}"))
		{{end -}}
	{{end -}}
{{end -}}
}
{{end}}

{{define "run" -}}
package cmd

import (
	"github.com/spf13/cobra"
)

func (c {{.Fuse}}Command) Runner() func(cmd *cobra.Command, args []string) {
	return {{.Fname}}Run
}

func {{.Fname}}Run(cmd *cobra.Command, args []string) {
	println(cmd.Name())
}
{{end}}
`
