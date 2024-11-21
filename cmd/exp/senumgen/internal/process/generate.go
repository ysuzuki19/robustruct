package process

import (
	"fmt"
	"strings"

	"github.com/ysuzuki19/robustruct/cmd/exp/senumgen/internal/process/coder"
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
		EnumDefName:   fmt.Sprintf("%sEnum%s", strings.ToLower(args.Name), coder.Bracket(defTypeParams)),
		EnumUseName:   fmt.Sprintf("%sEnum%s", strings.ToLower(args.Name), coder.Bracket(useTypeParams)),
		Variants:      args.AnalyzeResult.Variants,
	}

	c := coder.New().Globals(map[string]interface{}{
		"Package":       templateData.Package,
		"DefTypeParams": templateData.DefTypeParams,
		"UseTypeParams": templateData.UseTypeParams,
		"EnumDefName":   templateData.EnumDefName,
		"EnumUseName":   templateData.EnumUseName,
	})

	c.
		Str(`
// Code generated by senum; DO NOT EDIT.`).
		Format(`
package %s`, templateData.Package).LF().
		Str(`
type tag int`).LF().
		Block("const (", ")", func(c *coder.Coder) *coder.Coder {
			for _, variant := range templateData.Variants {
				c.Format("tag%s tag = iota", coder.Capitalize(variant.Name)).LF()
			}
			return c
		}).
		Format(`
type %s struct {
    %s%s
    tag tag
}`, templateData.EnumDefName, templateData.Package, coder.Bracket(templateData.UseTypeParams)).LF().
		Func(func(c *coder.Coder) *coder.Coder {
			for _, variant := range templateData.Variants {
				c.Tmpl(`
func New{{ .FieldName | capitalize }}{{ .DefTypeParams | bracket }}({{ if .HasData }}v {{ .TypeName }}{{ end }}) {{ .EnumUseName }} {
	return {{ .EnumUseName }}{
			{{ .Package }}: {{ .Package }}{{ .UseTypeParams | bracket }}{
					{{ .FieldName }}: {{ if .HasData }}v{{ else }}nil{{ end }},
			},
			tag: tag{{ .FieldName | capitalize }},
	}
}`, map[string]interface{}{
					"FieldName": variant.FieldName,
					"TypeName":  variant.TypeName,
					"HasData":   variant.HasData,
				})
			}
			return c
		}).LF().
		Func(func(c *coder.Coder) *coder.Coder {
			for _, variant := range templateData.Variants {
				c.Tmpl(`
func (e *{{ .EnumUseName }}) Is{{ .Name }}() bool {
	return e.tag == tag{{ .Name }}
}`, map[string]interface{}{
					"Name": coder.Capitalize(variant.Name),
				})
			}
			return c
		}).LF().
		Func(func(c *coder.Coder) *coder.Coder {
			for _, variant := range templateData.Variants {
				if variant.HasData {
					c.Tmpl(`
func (e *{{ .EnumUseName }}) As{{ .FieldName | capitalize }}() ({{ .TypeName }}, bool) {
	if e.Is{{ .FieldName | capitalize }}() {
		return e.{{ .Package }}.{{ .FieldName }}, true
	}
	return nil, false
}`, map[string]interface{}{
						"FieldName": variant.FieldName,
						"TypeName":  variant.TypeName,
					})
				}
			}
			return c
		}).LF().
		Format(`type Switcher%s struct`, coder.Bracket(templateData.DefTypeParams)).
		Block("{", "}", func(c *coder.Coder) *coder.Coder {
			for _, variant := range templateData.Variants {
				c.Capitalize(variant.Name).Space().Str("func").Block("(", ")", func(c *coder.Coder) *coder.Coder {
					if variant.HasData {
						c.Str("v").Space().Str(variant.TypeName)
					}
					return c
				}).LF()
			}
			return c
		}).LF().
		Tmpl(`
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
}`, map[string]interface{}{
			"Variants": templateData.Variants,
		}).LF().
		Tmpl(`
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
}`, map[string]interface{}{
			"Variants": templateData.Variants,
		}).LF()

	generated, err := c.Export()
	if err != nil {
		return nil, err
	}

	return generated, nil
}
