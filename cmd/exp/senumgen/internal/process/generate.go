package process

import (
	"fmt"
	"strings"

	tmpl "github.com/ysuzuki19/robustruct/cmd/exp/senumgen/internal/process/code_collector"
)

type GenerateArgs struct {
	DirPath       string
	Name          string
	AnalyzeResult AnalyzeResult
}

func Generate(args GenerateArgs) ([]byte, error) {
	defTypeParams := strings.Join(
		args.AnalyzeResult.TypeParams.Map(func(tp TypeParam) string {
			return fmt.Sprintf("%s %s", tp.Name, tp.TypeName)
		}),
		", ",
	)

	useTypeParams := strings.Join(
		args.AnalyzeResult.TypeParams.Map(func(tp TypeParam) string {
			return tp.Name
		}),
		", ",
	)

	templateData := struct {
		Package       string
		DefTypeParams string
		UseTypeParams string
		EnumDefName   string
		EnumUseName   string
		Variants      []Variant
	}{
		Package:       args.Name,
		DefTypeParams: defTypeParams,
		UseTypeParams: useTypeParams,
		EnumDefName:   fmt.Sprintf("%sEnum%s", strings.ToLower(args.Name), tmpl.Bracket(defTypeParams)),
		EnumUseName:   fmt.Sprintf("%sEnum%s", strings.ToLower(args.Name), tmpl.Bracket(useTypeParams)),
		Variants:      args.AnalyzeResult.Variants,
	}

	cc := tmpl.NewCodeCollector()

	cc.Merge(`
// Code generated by senum; DO NOT EDIT.
package {{ .Package }}
	`, struct{ Package string }{Package: templateData.Package}).
		Merge(`
type tag int
const (
{{- range $variant := .Variants }}
    tag{{ $variant.Name | capitalize }} tag = iota
{{- end }}
)
		`, struct{ Variants []Variant }{Variants: templateData.Variants}).
		Merge(`
type {{ .EnumDefName }} struct {
    {{ .Package }}{{.UseTypeParams | bracket}}
    tag tag
}
		`, struct {
			Package       string
			EnumDefName   string
			UseTypeParams string
		}{Package: templateData.Package, EnumDefName: templateData.EnumDefName, UseTypeParams: templateData.UseTypeParams}).
		Merge(`
{{- range $variant := .Variants }}
func New{{ $variant.Name | capitalize }}{{$.DefTypeParams | bracket}}({{ if $variant.HasData }}v {{ $variant.TypeName }}{{ end }}) {{ $.EnumUseName }} {
    return {{ $.EnumUseName }}{
        {{ $.Package }}: {{ $.Package }}{{ $.UseTypeParams | bracket }}{
            {{ $variant.FieldName }}: {{ if $variant.HasData }}v{{ else }}nil{{ end }},
        },
        tag: tag{{ $variant.Name | capitalize }},
    }
}
{{- end }}
		`, struct {
			Package       string
			Variants      []Variant
			EnumUseName   string
			EnumDefName   string
			DefTypeParams string
			UseTypeParams string
		}{
			Package:       templateData.Package,
			Variants:      templateData.Variants,
			EnumUseName:   templateData.EnumUseName,
			EnumDefName:   templateData.EnumDefName,
			DefTypeParams: templateData.DefTypeParams,
			UseTypeParams: templateData.UseTypeParams,
		}).
		Merge(`
{{- range $variant := .Variants }}
func (e *{{ $.EnumUseName }}) Is{{ $variant.Name | capitalize }}() bool {
    return e.tag == tag{{ $variant.Name | capitalize }}
}
{{- end }}
		`, struct {
			EnumUseName string
			Variants    []Variant
		}{EnumUseName: templateData.EnumUseName, Variants: templateData.Variants}).
		Merge(`
{{- range $variant := .Variants }}
{{- if $variant.HasData }}
func (e *{{ $.EnumUseName }}) As{{ $variant.Name | capitalize }}() ({{ $variant.TypeName }}, bool) {
    if e.Is{{ $variant.Name | capitalize }}() {
        return e.{{ $.Package }}.{{ $variant.FieldName }}, true
    }
    return nil, false
}
{{- end }}
{{- end }}
		`, struct {
			Package     string
			EnumUseName string
			Variants    []Variant
		}{Package: templateData.Package, EnumUseName: templateData.EnumUseName, Variants: templateData.Variants}).
		Merge(`
type Switcher{{.DefTypeParams | bracket}} struct {
{{- range $variant := .Variants }}
    {{ $variant.Name | capitalize }} func({{ if $variant.HasData }}v {{ $variant.TypeName }}{{ end }})
{{- end }}
}

func (e *{{ $.EnumUseName }}) Switch(s Switcher{{.UseTypeParams | bracket}}) {
    switch e.tag {
    {{- range $variant := .Variants }}
    case tag{{ $variant.Name | capitalize }}:
        {{- if $variant.HasData }}
        s.{{ $variant.Name | capitalize }}(e.{{ $.Package }}.{{ $variant.FieldName }})
        {{- else }}
        s.{{ $variant.Name | capitalize }}()
        {{- end }}
    {{- end }}
    }
}
		`, struct {
			Package       string
			EnumUseName   string
			Variants      []Variant
			DefTypeParams string
			UseTypeParams string
		}{Package: templateData.Package, EnumUseName: templateData.EnumUseName, Variants: templateData.Variants, DefTypeParams: templateData.DefTypeParams, UseTypeParams: templateData.UseTypeParams}).
		Merge(`
type Matcher[MatchResult any {{.DefTypeParams | csvConnect}}] struct {
{{- range $variant := .Variants }}
    {{ $variant.Name | capitalize }} func({{ if $variant.HasData }}v {{ $variant.TypeName }}{{ end }}) MatchResult
{{- end }}
}

func Match[MatchResult any {{.DefTypeParams | csvConnect}}](e *{{ .EnumUseName }}, m Matcher[MatchResult {{.UseTypeParams | csvConnect}}]) MatchResult {
    switch e.tag {
    {{- range $variant := .Variants }}
    case tag{{ $variant.Name | capitalize }}:
        {{- if $variant.HasData }}
        return m.{{ $variant.Name | capitalize }}(e.{{ $.Package }}.{{ $variant.FieldName }})
        {{- else }}
        return m.{{ $variant.Name | capitalize }}()
        {{- end }}
    {{- end }}
    }
    panic("unreachable: invalid tag")
}
		`, struct {
			Package       string
			EnumUseName   string
			Variants      []Variant
			DefTypeParams string
			UseTypeParams string
		}{Package: templateData.Package, EnumUseName: templateData.EnumUseName, Variants: templateData.Variants, DefTypeParams: templateData.DefTypeParams, UseTypeParams: templateData.UseTypeParams})

	generated, err := cc.Export()
	if err != nil {
		return nil, err
	}

	return generated, nil
}
