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

var {{$cmd}} = &cobra.Command{
	{{with .Use}}Use: "{{.}}",{{end}}
	Run: {{$cmd}}Run,
	{{with .Aliases}}Aliases: []string{
		{{- range .}}
		"{{.}}",
	},
	{{- end}}
	{{- end}}
	{{with .Short}}Short: "{{.}}",{{end}}
	{{with .Long}}Long: "{{.}}",{{end}}
}
{{end}}

func init() {
{{- range .Commands -}}
	{{- $cmd := .Name -}}

	{{- with .Parent}}
		{{.}}.AddCommand({{$cmd}})
	{{- end}}

	{{range .Flags -}}
		{{$cmd}}.
		{{- with .Persistent}}Persistent{{end -}}
		Flags().{{.Type}}("{{.Name}}",
		{{- with .Shorthand}}"{{.}}",{{end}}
		{{- .Value -}}
		, "{{.Usage}}")

		{{if .Viper}}
			viper.BindPFlag("{{.Name}}", {{$cmd}}.Flags().Lookup("{{.Name}}"))
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

func {{.Name}}Run(cmd *cobra.Command, args []string) {
	println(cmd.Name())
}
{{end}}
`
