package process

import (
	"fmt"
	"strings"

	"github.com/ysuzuki19/robustruct/cmd/exp/senumgen/internal/process/value_type"
)

type Field struct {
	Name  string
	Value string
}

type AnalyzeResult struct {
	Name      string
	ValueType value_type.ValueType
	Fields    []Field
	ZeroField Field
}

func Analyze(commentLines []string) (AnalyzeResult, error) {
	analyzeResult := AnalyzeResult{
		Name:      "",
		ValueType: value_type.Default(),
		Fields:    []Field{},
		ZeroField: Field{}, // ignore:fields_require
	}

	var fieldNames []string
	fieldValues := make(map[string]string)

	for _, line := range commentLines {
		line = strings.ReplaceAll(line, " ", "")
		{
			prefix := "$.name="
			fmt.Println(line)
			if strings.HasPrefix(line, prefix) {
				analyzeResult.Name = strings.TrimPrefix(line, prefix)
			}
		}
		{
			prefix := "$.value_type="
			if strings.HasPrefix(line, prefix) {
				trimmed := strings.TrimPrefix(line, prefix)
				vt, err := value_type.FromString(trimmed)
				if err != nil {
					return AnalyzeResult{}, fmt.Errorf("invalid value type: %s", trimmed) //ignore:fields_require
				}
				analyzeResult.ValueType = vt
			}
		}
		{
			prefix := "$.fields="
			if strings.HasPrefix(line, prefix) {
				body := strings.TrimPrefix(line, prefix)
				splitted := strings.Split(body, ",")
				fieldNames = append(fieldNames, splitted...)
			}
		}
		{
			prefix := "$.zero="
			if strings.HasPrefix(line, prefix) {
				analyzeResult.ZeroField.Name = strings.TrimPrefix(line, prefix)
			}
		}
		{
			prefix := "$.field."
			if strings.HasPrefix(line, prefix) {
				body := strings.TrimPrefix(line, prefix)
				splitted := strings.Split(body, "=")
				if len(splitted) != 2 {
					return AnalyzeResult{}, fmt.Errorf("invalid field line: %s", line) // ignore:fields_require
				}
				fieldValues[splitted[0]] = splitted[1]
			}
		}
	}

	if analyzeResult.ValueType.IsInt() {
		if len(fieldValues) == 0 {
			// iota like fields 0, 1, 2, ...
			value := 1 // 0 is (PACKAGE).zero
			for _, fieldName := range fieldNames {
				fmt.Println(fieldName)
				fieldValues[fieldName] = fmt.Sprint(value)
				value++
			}
		}
	}
	for _, name := range fieldNames {
		value, ok := fieldValues[name]
		if !ok {
			return AnalyzeResult{}, fmt.Errorf("field value not found: %s", name) //ignore:fields_require
		}
		analyzeResult.Fields = append(analyzeResult.Fields, Field{
			Name:  name,
			Value: value,
		})
	}

	if analyzeResult.Name == "" {
		return AnalyzeResult{}, fmt.Errorf("name is not found") //ignore:fields_require
	}

	if analyzeResult.ZeroField.Name == "" {
		return AnalyzeResult{}, fmt.Errorf("zero is not found") //ignore:fields_require
	}

	if len(analyzeResult.Fields) == 0 {
		return AnalyzeResult{}, fmt.Errorf("fields has no element") //ignore:fields_require
	}

	if analyzeResult.ValueType.IsString() {
		if len(analyzeResult.Fields) != len(fieldValues) {
			return AnalyzeResult{}, fmt.Errorf("fields are invalid") //ignore:fields_require
		}
	}
	if analyzeResult.ValueType.IsInt() {
		if len(analyzeResult.Fields) != 0 && len(analyzeResult.Fields) != len(fieldValues) {
			return AnalyzeResult{}, fmt.Errorf("fields are invalid") //ignore:fields_require
		}
	}

	if v, ok := fieldValues[analyzeResult.ZeroField.Name]; ok {
		analyzeResult.ZeroField.Value = v
	} else {
		return AnalyzeResult{}, fmt.Errorf("zero value is not found") //ignore:fields_require
	}

	return analyzeResult, nil
}
