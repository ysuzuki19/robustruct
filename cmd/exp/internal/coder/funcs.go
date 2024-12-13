package coder

import (
	"fmt"
	"strings"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func Bracket(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("[%s]", strings.TrimSpace(s))
}

func csvConnect(s string) string {
	if s == "" {
		return ""
	}
	return ", " + s
}